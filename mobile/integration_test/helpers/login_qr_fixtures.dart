import 'dart:convert';

import 'package:http/http.dart' as http;

import 'api.dart';
import 'env.dart';

/// Turns on the desktop login QR for [serverUrl] and returns the deep link the
/// backend composed for it (`FeatureConfig.loginQrUrl`), restoring the previous
/// settings via `addTearDown`.
///
/// `showLoginQr` / `mobileServerUrl` are GLOBAL system settings, so this mutates
/// shared server state -- the same caution the desktop suite documents in
/// `desktop/e2e/login-qr.spec.ts`. Both this and the `aiPoweredReceipts`
/// fixtures in `feature_flags.dart` go through `overrideSystemSettingsForTest`,
/// which restores by patching onto a FRESH read rather than replaying a stale
/// snapshot, so a concurrent change to an unrelated setting isn't clobbered.
///
/// The returned string is the REAL production link
/// (`https://receiptwrangler.io/app/setup#url=<url.QueryEscape(serverUrl)>`,
/// `api/internal/services/system_settings.go`), which is the whole point: tests
/// feed it to the app instead of a hand-written literal, so the Go escaping and
/// the Dart `extractDeepLinkServerUrl` parsing are pinned against each other.
Future<String> enableLoginQrForTest(String serverUrl) async {
  final jwt = await apiLogin();
  final settings = await getSystemSettings(jwt);

  await overrideSystemSettingsForTest(
    jwt,
    settings,
    overrides: {
      'showLoginQr': true,
      'mobileServerUrl': serverUrl,
    },
    restoreTo: {
      'showLoginQr': settings['showLoginQr'] ?? false,
      'mobileServerUrl': settings['mobileServerUrl'] ?? '',
    },
  );

  final loginQrUrl = await getLoginQrUrl();
  if (loginQrUrl == null || loginQrUrl.isEmpty) {
    throw StateError(
      'featureConfig.loginQrUrl is empty after enabling the login QR for '
      '$serverUrl -- the backend did not compose a deep link.',
    );
  }
  return loginQrUrl;
}

/// Reads `loginQrUrl` off `GET /featureConfig`, returning null when the login QR
/// is disabled. The route is UNAUTHENTICATED (`api/internal/routers/feature_config.go`
/// mounts no token middleware) -- exactly how a not-yet-connected phone and the
/// pre-auth desktop login page reach it -- so no jwt is sent.
Future<String?> getLoginQrUrl() async {
  final res = await http
      .get(Uri.parse('${E2eEnv.baseUrl}/featureConfig'))
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError(
      'GET featureConfig failed: HTTP ${res.statusCode}: ${res.body}',
    );
  }
  final body = jsonDecode(res.body) as Map<String, dynamic>;
  return body['loginQrUrl'] as String?;
}
