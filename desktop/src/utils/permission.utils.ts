/**
 * TypeScript port of the backend permission matcher
 * (`api/internal/permissions/matcher.go`). Kept byte-for-byte consistent with
 * the Go implementation so the desktop UI gates exactly as the server enforces.
 *
 * Segments are split on '.'. A '*' as the final granted segment matches one or
 * more remaining required segments (e.g. "group.*" satisfies
 * "group.receipts.read"). A '*' in a non-final position matches exactly one
 * required segment (e.g. "group.*.read" satisfies "group.receipts.read"). "*"
 * alone matches any permission. The optional ":sub-scope" suffix (e.g.
 * "read:any") is treated as part of the final segment and compared literally.
 */

/**
 * Reports whether a single granted permission string satisfies a single
 * required permission key. The granted string may contain '*' wildcard segments
 * (e.g. as stored on a role); the required key is expected to be concrete.
 */
export function matches(granted: string, required: string): boolean {
  if (granted === required) {
    return true;
  }
  if (granted === "*") {
    return true;
  }

  const grantedSegments = granted.split(".");
  const requiredSegments = required.split(".");

  for (let i = 0; i < grantedSegments.length; i++) {
    const segment = grantedSegments[i];

    if (segment === "*") {
      if (i === grantedSegments.length - 1) {
        // Trailing wildcard: matches one or more remaining segments.
        return i < requiredSegments.length;
      }
      if (i >= requiredSegments.length) {
        return false;
      }
      continue;
    }

    if (i >= requiredSegments.length || segment !== requiredSegments[i]) {
      return false;
    }
  }

  // No trailing wildcard consumed the remainder, so the segment counts must
  // line up exactly (a more specific grant must not match a broader key).
  return grantedSegments.length === requiredSegments.length;
}

/** Reports whether any of the granted permissions matches the required key. */
function satisfies(granted: string[], required: string): boolean {
  return granted.some((grant) => matches(grant, required));
}

/**
 * Reports whether the granted permissions satisfy ALL of the required
 * permissions (logical AND). Returns false when no required permissions are
 * supplied (deny by default), matching the Go implementation.
 */
export function hasAll(granted: string[], ...required: string[]): boolean {
  if (required.length === 0) {
    return false;
  }

  return required.every((req) => satisfies(granted, req));
}

/**
 * Reports whether the granted permissions satisfy AT LEAST ONE of the required
 * permissions (logical OR). Returns false when no required permissions are
 * supplied (deny by default).
 */
export function hasAny(granted: string[], ...required: string[]): boolean {
  return required.some((req) => satisfies(granted, req));
}
