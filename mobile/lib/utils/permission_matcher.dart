// Dart port of the backend permission matcher
// (api/internal/permissions/matcher.go) and its desktop twin
// (desktop/src/utils/permission.utils.ts). Kept consistent with both so the
// mobile UI gates exactly as the server enforces.
//
// Segments are split on '.'. A '*' as the final granted segment matches one or
// more remaining required segments (e.g. "group.*" satisfies
// "group.receipts.read"). A '*' in a non-final position matches exactly one
// required segment (e.g. "group.*.read" satisfies "group.receipts.read"). "*"
// alone matches any permission. The optional ":sub-scope" suffix (e.g.
// "read:any") is treated as part of the final segment and compared literally.

/// Reports whether a single granted permission string satisfies a single
/// required permission key. The granted string may contain '*' wildcard
/// segments (e.g. as stored on a role); the required key is expected to be
/// concrete.
bool permissionMatches(String granted, String required) {
  if (granted == required) {
    return true;
  }
  if (granted == '*') {
    return true;
  }

  final grantedSegments = granted.split('.');
  final requiredSegments = required.split('.');

  for (var i = 0; i < grantedSegments.length; i++) {
    final segment = grantedSegments[i];

    if (segment == '*') {
      if (i == grantedSegments.length - 1) {
        // Trailing wildcard: matches one or more remaining segments.
        return i < requiredSegments.length;
      }
      if (i >= requiredSegments.length) {
        return false;
      }
      continue;
    }

    if (i >= requiredSegments.length || segment != requiredSegments[i]) {
      return false;
    }
  }

  // No trailing wildcard consumed the remainder, so the segment counts must
  // line up exactly (a more specific grant must not match a broader key).
  return grantedSegments.length == requiredSegments.length;
}

/// Reports whether any of the granted permissions matches the required key.
bool _satisfies(List<String> granted, String required) {
  return granted.any((grant) => permissionMatches(grant, required));
}

/// Reports whether the granted permissions satisfy ALL of the required
/// permissions (logical AND). Returns false when no required permissions are
/// supplied (deny by default), matching the Go/TS implementations.
bool hasAll(List<String> granted, List<String> required) {
  if (required.isEmpty) {
    return false;
  }
  return required.every((req) => _satisfies(granted, req));
}

/// Reports whether the granted permissions satisfy AT LEAST ONE of the required
/// permissions (logical OR). Returns false when no required permissions are
/// supplied (deny by default).
bool hasAny(List<String> granted, List<String> required) {
  return required.any((req) => _satisfies(granted, req));
}
