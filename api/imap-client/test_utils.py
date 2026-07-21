import unittest

from utils import sanitize_filename, valid_from_email, valid_subject


class TestSanitizeFilename(unittest.TestCase):

    def test_plain_name_unchanged(self):
        self.assertEqual(sanitize_filename("receipt.jpg"), "receipt.jpg")

    def test_relative_traversal_reduced_to_base(self):
        self.assertEqual(sanitize_filename("../../../../tmp/evil.txt"), "evil.txt")

    def test_absolute_path_reduced_to_base(self):
        self.assertEqual(sanitize_filename("/etc/cron.d/evil"), "evil")

    def test_nested_path_reduced_to_base(self):
        self.assertEqual(sanitize_filename("a/b/c.jpg"), "c.jpg")

    def test_module_overwrite_reduced_to_base(self):
        self.assertEqual(sanitize_filename("../imap-client/utils.py"), "utils.py")

    def test_dotdot_rejected(self):
        self.assertEqual(sanitize_filename(".."), "")

    def test_dot_rejected(self):
        self.assertEqual(sanitize_filename("."), "")

    def test_empty_rejected(self):
        self.assertEqual(sanitize_filename(""), "")

    def test_none_rejected(self):
        self.assertEqual(sanitize_filename(None), "")


class TestEmailUtils(unittest.TestCase):

    def test_valid_from_email_empty_whitelist(self):
        self.assertTrue(valid_from_email("test@example.com", []))

    def test_valid_from_email_in_whitelist(self):
        whitelist = [{"email": "test@example.com"}, {"email": "admin@example.com"}]
        self.assertTrue(valid_from_email("test@example.com", whitelist))

    def test_valid_from_email_not_in_whitelist(self):
        whitelist = [{"email": "admin@example.com"}]
        self.assertFalse(valid_from_email("test@example.com", whitelist))

    def test_valid_subject_empty_regex_list(self):
        self.assertTrue(valid_subject("Subject", []))

    def test_valid_subject_match_found(self):
        subject_line_regexes = [{"regex": r"^Test"}]
        self.assertTrue(valid_subject("Test subject", subject_line_regexes))

    def test_valid_subject_no_match_found(self):
        subject_line_regexes = [{"regex": r"^Hello"}]
        self.assertFalse(valid_subject("Test subject", subject_line_regexes))


if __name__ == "__main__":
    unittest.main()
