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

/// Extracts the server URL from a receiptwrangler.io/app/setup deep link.
/// Returns null if [raw] is not such a link or the inner url is invalid.
///
/// The deep link carries the server URL in its FRAGMENT as
/// `#url=<percent-encoded server url>`, e.g.
/// `https://receiptwrangler.io/app/setup#url=https%3A%2F%2Fdemo.receiptwrangler.io%2Fapi`.
/// The same link is produced by the desktop login QR and works both when the OS
/// opens the app (App Links / Universal Links) and when it's scanned by the
/// in-app QR scanner. The inner value is passed back through [normalizeServerUrl]
/// so the exact same http/https + non-empty-host validation gate is reused, and
/// the URL is never auto-connected — the user still reviews it and taps Connect.
String? extractDeepLinkServerUrl(String raw) {
  final uri = Uri.tryParse(raw.trim());
  if (uri == null) {
    return null;
  }
  if (uri.host != 'receiptwrangler.io' || uri.path != '/app/setup') {
    return null;
  }

  final inner = Uri.splitQueryString(uri.fragment)['url'];
  if (inner == null) {
    return null;
  }

  return normalizeServerUrl(inner);
}
