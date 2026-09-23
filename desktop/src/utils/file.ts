/**
 * Reads the file name out of a `Content-Disposition` header, for a download whose
 * name only the server knows.
 *
 * Handles the three forms the server can emit, in the order RFC 6266 prescribes:
 * the RFC 5987 `filename*` (percent-encoded, with a charset and optional language
 * tag), the quoted `filename="…"` including backslash-escaped characters, and the
 * bare `filename=…`. `filename*` is preferred when both are present, because that
 * is the one that can carry non-ASCII.
 *
 * `filename*` is not hypothetical: the server builds this header with Go's
 * `mime.FormatMediaType`, which switches to that form for any name outside
 * us-ascii and emits it INSTEAD of the plain parameter. A parser that only reads
 * `filename=` would fall through to the generic fallback for every accented or
 * non-Latin attachment name.
 *
 * The fallback is deliberately generic rather than an error: the bytes are the
 * point, and a download that lands under an unhelpful name still beats one that
 * does not happen.
 */
export function filenameFromContentDisposition(
  header: string | null,
  fallback: string = "download"
): string {
  if (!header) {
    return fallback;
  }

  return (
    extendedFilename(header) ?? quotedFilename(header) ?? bareFilename(header) ?? fallback
  );
}

// filename*=UTF-8''receipt%20final.pdf — charset'language'percent-encoded-value.
// A value that will not decode (a stray "%" is enough) falls through to the other
// forms rather than throwing.
function extendedFilename(header: string): string | null {
  const match = /filename\*\s*=\s*[^']*'[^']*'([^;]+)/i.exec(header);
  if (!match) {
    return null;
  }

  try {
    return nonEmpty(decodeURIComponent(match[1].trim()));
  } catch {
    return null;
  }
}

// filename="receipt\"final.pdf" — a backslash escapes the character after it, so
// the quoted value ends at the first UNescaped quote.
function quotedFilename(header: string): string | null {
  const match = /filename\s*=\s*"((?:[^"\\]|\\.)*)"/i.exec(header);
  if (!match) {
    return null;
  }

  return nonEmpty(match[1].replace(/\\(.)/g, "$1").trim());
}

// filename=receipt.png — an unquoted token. A value that OPENS with a quote is
// the quoted form and belongs to the branch above, even when that branch declined
// it: otherwise `filename="   "` falls through to here and yields the quote
// characters themselves as the name.
function bareFilename(header: string): string | null {
  const match = /filename\s*=\s*([^";][^;]*)/i.exec(header);

  return match ? nonEmpty(match[1].trim()) : null;
}

function nonEmpty(value: string): string | null {
  return value.length > 0 ? value : null;
}

export function downloadFile(blob: Blob, filename: string): void {
  const url = window.URL.createObjectURL(blob);
  const link = document.createElement("a");
  link.href = url;
  link.setAttribute("download", filename);
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
}
