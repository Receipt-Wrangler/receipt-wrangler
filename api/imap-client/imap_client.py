import datetime
import email
import logging
import os
import re
import ssl
from html.parser import HTMLParser
from io import StringIO
from mailbox import Message

from imapclient import IMAPClient

from utils import valid_from_email, valid_subject

base_path = os.environ.get("BASE_PATH", "")


class _HTMLTextExtractor(HTMLParser):
    """Simple HTML parser that extracts visible text content."""

    def __init__(self):
        super().__init__()
        self._result = StringIO()

    def handle_data(self, data):
        self._result.write(data)

    def get_text(self):
        return self._result.getvalue()


def strip_html_tags(html):
    """Strip HTML tags and return plain text."""
    extractor = _HTMLTextExtractor()
    extractor.feed(html)
    return extractor.get_text()


class ImapClient:
    host = None
    port = None
    username = None
    password = None
    subject_line_regexes = None
    email_whitelist = None
    client = None

    def __init__(self, host, port, username, password, use_starttls, subject_line_regexes, email_whitelist):
        self.host = host
        self.port = port
        self.username = username
        self.password = password
        self.use_starttls = use_starttls
        self.subject_line_regexes = subject_line_regexes or []
        self.email_whitelist = email_whitelist or []

    def connect(self):
        # When connecting to a self-signed mail server (e.g. an internal
        # homelab Dovecot), certificate verification will fail with an
        # "unknown ca" alert. IMAP_TLS_SKIP_VERIFY=true disables verification
        # for these trusted internal hosts. It defaults to off so public
        # servers are still verified normally.
        skip_verify = os.environ.get("IMAP_TLS_SKIP_VERIFY", "").lower() in ("1", "true", "yes")
        ssl_context = None
        if skip_verify:
            ssl_context = ssl.create_default_context()
            ssl_context.check_hostname = False
            ssl_context.verify_mode = ssl.CERT_NONE

        if self.use_starttls:
            self.client = IMAPClient(self.host, self.port, ssl=False)
            self.client.starttls(ssl_context=ssl_context)
        else:
            self.client = IMAPClient(self.host, self.port, ssl_context=ssl_context)

        self.client.login(self.username, self.password)

    def get_unread_email_metadata(self):
        unread_emails = self._get_unread_emails()
        return self._messages_to_email_metadata(unread_emails)

    def _get_unread_emails(self):
        if self.client is None:
            self.connect()

        self.client.select_folder('INBOX')

        messages = self.client.search(['UNSEEN'])
        response = self.client.fetch(messages, ['FLAGS', 'RFC822'])
        return response or None

    def _messages_to_email_metadata(self, response):
        if response is None:
            return []

        results = []
        for message_id, data in response.items():
            # Process each message independently: a single malformed message
            # must not abort the whole poll (which would drop every fetched
            # message and re-raise to the Go caller as exit status 1).
            try:
                formatted_data = self._get_formatted_message_data(data)
            except Exception as e:
                logging.error(f"Skipping message {message_id}: {e}")
                continue
            if len(formatted_data) > 0:
                formatted_data[message_id] = message_id
                results.append(formatted_data)

        return results

    def _get_formatted_message_data(self, data):
        raw_rfc822 = data.get(b"RFC822") if data else None
        if not raw_rfc822:
            logging.warning("Skipping message with empty RFC822 payload")
            return {}
        message_data = email.message_from_bytes(raw_rfc822)

        from_data = self._get_formatted_to_or_from_data(message_data, "From")
        if from_data["email"] is None:
            return {}

        to_data = self._get_formatted_to_or_from_data(message_data, "To")
        if to_data["email"] is None:
            return {}

        subject = message_data.get("Subject")
        should_process = self._valid_from_email(from_data["email"]) and self._valid_subject(subject)

        # TODO: V5 - Could we set this message to unread if not process?
        if not should_process:
            return {}

        formatted_date = self.get_formatted_date(message_data.get("Date"))

        attachments = self._get_attachments(message_data)
        body, body_html = self._get_body_text(message_data)

        result = {
            "date": formatted_date,
            "subject": subject,
            "to": to_data["email"],
            "fromName": from_data["name"],
            "fromEmail": from_data["email"],
            "attachments": attachments,
            "body": body,
            "bodyHtml": body_html,
            "groupSettingsIds": [],
        }

        if len(attachments) == 0 and not body and not body_html:
            return {}

        logging.info(f"Formatted message data: {result}")
        return result

    def _get_formatted_to_or_from_data(self, message_data: Message, key: str):
        result = {
            "name": None,
            "email": None
        }

        from_data = message_data.get(key).split("<")
        if len(from_data) == 2:
            result["name"] = from_data[0]
            result["email"] = from_data[1].replace("<", "").replace(">", "")

        if len(from_data) == 1:
            result["email"] = from_data[0]

        logging.info(f"Formatted from data: {result}")

        return result

    def get_formatted_date(self, date):
        date_parts = date.split(
            "(")  # Fixes case when date comes back in utc , so (UTC) is appended
        logging.info(date[0])
        date = datetime.datetime.strptime(
            date_parts[0].strip(), "%a, %d %b %Y %H:%M:%S %z")
        utc_date = date.replace(tzinfo=datetime.timezone.utc)
        formatted_date = utc_date.strftime("%Y-%m-%dT%H:%M:%S.%fZ")

        return formatted_date

    def _valid_from_email(self, from_email):
        return valid_from_email(from_email, self.email_whitelist)

    def _valid_subject(self, subject):
        return valid_subject(subject, self.subject_line_regexes)

    def _get_attachments(self, message_data: Message):
        result = []
        for part in message_data.walk():
            if part.get_content_maintype() == 'multipart':
                continue
            if part.get('Content-Disposition') is None:
                continue

            filename = part.get_filename()
            mime_type = part.get_content_type()

            logging.info(f"Filename: {filename} mime_type: {mime_type}")

            if len(filename) > 0 and self.valid_mime_type(mime_type):
                payload = part.get_payload(decode=True)
                if payload is None:
                    logging.warning(f"Skipping attachment with undecodable payload: {filename}")
                    continue
                filePath = os.path.join(base_path, "temp", filename)
                with open(filePath, 'wb') as f:
                    f.write(payload)

                size = os.path.getsize(filePath)
                data = {
                    "filename": filename,
                    "fileType": mime_type,
                    "size": size,
                }

                result.append(data)

        return result

    def _get_body_text(self, message_data: Message):
        """Extract email body. Returns (text, html) tuple where text is the
        plain/stripped representation used for fallback prompting and html is
        the raw HTML used downstream for chromedp PDF rendering.
        """
        plain_parts = []
        html_parts = []

        for part in message_data.walk():
            if part.get_content_maintype() == 'multipart':
                continue

            content_disposition = str(part.get('Content-Disposition') or '')
            if 'attachment' in content_disposition:
                continue

            content_type = part.get_content_type()
            charset = part.get_content_charset() or 'utf-8'

            if content_type == 'text/plain':
                payload = part.get_payload(decode=True)
                if payload:
                    plain_parts.append(payload.decode(charset, errors='replace'))
            elif content_type == 'text/html':
                payload = part.get_payload(decode=True)
                if payload:
                    html_parts.append(payload.decode(charset, errors='replace'))

        html_body = '\n'.join(html_parts) if html_parts else ""

        if plain_parts:
            text = '\n'.join(plain_parts)
        elif html_parts:
            text = '\n'.join(strip_html_tags(html) for html in html_parts)
        else:
            return "", html_body

        # Collapse excessive whitespace
        text = re.sub(r'[ \t]+', ' ', text)
        text = re.sub(r'\n{3,}', '\n\n', text)
        return text.strip(), html_body

    def valid_mime_type(self, mime_type):
        image_mime_types_regex = r"^(image\/(jpeg|png|heic|bmp|webp|tiff)|application\/pdf)$"
        match = re.search(image_mime_types_regex, mime_type or "")
        return match is not None
