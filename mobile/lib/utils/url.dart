/// Returns the trimmed URL when [raw] is a well-formed absolute http/https URL
/// with a non-empty host; otherwise null.
///
/// Both http and https are accepted on purpose: self-hosted instances are
/// commonly reached over plain http on a LAN or bare IP. The string is returned
/// trimmed but otherwise verbatim (no trailing-slash rewriting) so it matches
/// exactly what hand-typed input produces and what the user reviews in the field
/// before connecting.
String? normalizeServerUrl(String raw) {
  final trimmed = raw.trim();
  if (trimmed.isEmpty) {
    return null;
  }

  final uri = Uri.tryParse(trimmed);
  if (uri == null || !uri.hasScheme) {
    return null;
  }

  final scheme = uri.scheme.toLowerCase();
  if (scheme != 'http' && scheme != 'https') {
    return null;
  }
  if (uri.host.isEmpty) {
    return null;
  }

  return trimmed;
}
