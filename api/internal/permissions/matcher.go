package permissions

import "strings"

// matches reports whether a single granted permission string satisfies a single
// required permission key. The granted string may contain '*' wildcard segments
// (e.g. as stored on a role); the required key is expected to be a concrete
// permission.
//
// Segments are split on '.'. A '*' as the final granted segment matches one or
// more remaining required segments (e.g. "group.*" satisfies
// "group.receipts.read"). A '*' in a non-final position matches exactly one
// required segment (e.g. "group.*.read" satisfies "group.receipts.read"). "*"
// alone matches any permission.
//
// The optional ":sub-scope" suffix (e.g. "read:any") is treated as part of the
// final segment for now and compared literally; superset semantics for :any are
// intentionally not implemented yet.
func matches(granted string, required string) bool {
	if granted == required {
		return true
	}
	if granted == "*" {
		return true
	}

	grantedSegments := strings.Split(granted, ".")
	requiredSegments := strings.Split(required, ".")

	for i, segment := range grantedSegments {
		if segment == "*" {
			if i == len(grantedSegments)-1 {
				// Trailing wildcard: matches one or more remaining segments.
				return i < len(requiredSegments)
			}
			if i >= len(requiredSegments) {
				return false
			}
			continue
		}

		if i >= len(requiredSegments) || segment != requiredSegments[i] {
			return false
		}
	}

	// No trailing wildcard consumed the remainder, so the segment counts must
	// line up exactly (a more specific grant must not match a broader key).
	return len(grantedSegments) == len(requiredSegments)
}

// satisfies reports whether any of the granted permissions matches the required key.
func satisfies(granted []string, required string) bool {
	for _, grant := range granted {
		if matches(grant, required) {
			return true
		}
	}
	return false
}

// HasAll reports whether the granted permissions satisfy ALL of the required
// permissions (logical AND). This is the default matcher; called with a single
// required permission it answers a simple "does the user have this permission?"
// check. Returns false when no required permissions are supplied (deny by
// default).
func HasAll(granted []string, required ...string) bool {
	if len(required) == 0 {
		return false
	}

	for _, req := range required {
		if !satisfies(granted, req) {
			return false
		}
	}
	return true
}

// HasAny reports whether the granted permissions satisfy AT LEAST ONE of the
// required permissions (logical OR). Returns false when no required permissions
// are supplied (deny by default).
func HasAny(granted []string, required ...string) bool {
	for _, req := range required {
		if satisfies(granted, req) {
			return true
		}
	}
	return false
}
