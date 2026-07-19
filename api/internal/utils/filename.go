package utils

import (
	"os"
	"path/filepath"
)

// SanitizeFileName reduces an untrusted filename (for example an email
// attachment's MIME Content-Disposition name) to a safe single path component,
// so joining it onto a directory can never escape that directory
// (path traversal, CWE-22). It strips any directory portion with filepath.Base
// and returns "" when the name is missing or cannot be reduced to a usable
// component ("." / ".."), letting the caller skip it.
func SanitizeFileName(name string) string {
	base := filepath.Base(name)
	if base == "." || base == ".." || base == string(os.PathSeparator) {
		return ""
	}
	return base
}
