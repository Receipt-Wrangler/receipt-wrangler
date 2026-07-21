package utils

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestSanitizeFileName(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected string
	}{
		{"plain name is unchanged", "receipt.jpg", "receipt.jpg"},
		{"relative traversal reduced to base", "../../../../tmp/evil.txt", "evil.txt"},
		{"absolute path reduced to base", "/etc/cron.d/evil", "evil"},
		{"nested path reduced to base", "a/b/c.jpg", "c.jpg"},
		{"module overwrite reduced to base", "../imap-client/utils.py", "utils.py"},
		{"dotdot is rejected", "..", ""},
		{"dot is rejected", ".", ""},
		{"empty is rejected", "", ""},
		{"root is rejected", "/", ""},
		{"trailing slash reduced to base", "foo/bar/", "bar"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := SanitizeFileName(tc.input)
			if got != tc.expected {
				t.Errorf("SanitizeFileName(%q) = %q, want %q", tc.input, got, tc.expected)
			}
		})
	}
}

// TestSanitizeFileNameStaysWithinDirectory is the security property: for any
// malicious input, joining the sanitized result onto a base directory must not
// escape that directory.
func TestSanitizeFileNameStaysWithinDirectory(t *testing.T) {
	dir := "/srv/app/temp"
	malicious := []string{
		"../../../../tmp/evil.txt",
		"../imap-client/utils.py",
		"/etc/passwd",
		"..",
		"foo/../../bar",
	}

	for _, input := range malicious {
		safe := SanitizeFileName(input)
		if safe == "" {
			continue // rejected outright — nothing is joined
		}
		joined := filepath.Join(dir, safe)
		if joined != filepath.Join(dir, filepath.Base(joined)) || !strings.HasPrefix(joined, dir+string(filepath.Separator)) {
			t.Errorf("SanitizeFileName(%q)=%q joined to %q escaped the directory", input, safe, joined)
		}
	}
}
