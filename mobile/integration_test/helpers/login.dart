import 'dart:io' show Platform;

import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/groups/screens/group_select.dart';
import 'package:receipt_wrangler_mobile/main.dart' show buildApp;
import 'package:receipt_wrangler_mobile/persistence/global_shared_preferences.dart';

import 'env.dart';
import 'form_actions.dart';
import 'pump.dart';

/// Pumps a fresh app widget tree and walks the SetHomeserverUrl + Login
/// screens using the admin credentials from `--dart-define=E2E_*`. Returns
/// when `GroupSelect` is on screen (= logged-in landing).
///
/// Pre-conditions:
/// - `IntegrationTestWidgetsFlutterBinding.ensureInitialized()` was
///   called by the caller's `main()`.
/// - On Linux, `installLinuxDesktopMocks()` was called by the caller's
///   `main()` (mobile-only plugins are stubbed there). On Android/iOS
///   the call should be skipped via `Platform.isLinux` so real plugin
///   channels run.
///
/// `buildApp()` returns a fresh widget tree (providers + ReceiptWrangler)
/// with a per-`State` `late final GoRouter`, so previous tests' router
/// location and provider state never leak across `testWidgets`.
Future<void> loginAsAdmin(WidgetTester tester) async {
  E2eEnv.assertAdmin();
  await loginAs(
    tester,
    username: E2eEnv.adminUsername,
    password: E2eEnv.adminPassword,
  );
}

/// Same flow as [loginAsAdmin] but uses the regular e2e-user credentials.
/// Use for permission-gating specs that need a non-admin caller — e.g. a
/// Legacy Viewer / Legacy Editor group member (provisioned over the API)
/// or a user that belongs to no group at all.
Future<void> loginAsUser(WidgetTester tester) async {
  E2eEnv.assertUser();
  await loginAs(
    tester,
    username: E2eEnv.userUsername,
    password: E2eEnv.userPassword,
  );
}

/// Generic login: pumps a fresh app tree and walks the SetHomeserverUrl +
/// Login screens with [username]/[password], returning once `GroupSelect`
/// is on screen. [loginAsAdmin] / [loginAsUser] are thin wrappers.
Future<void> loginAs(
  WidgetTester tester, {
  required String username,
  required String password,
}) async {
  await resetPersistedAppState();

  await tester.pumpWidget(buildApp());

  await pumpUntilFound(tester, find.text('Server URL'));

  await tester.enterText(formField('url'), E2eEnv.baseUrl);
  await tester.tap(filledButton('Connect'));

  await loginFromLoginScreen(tester, username: username, password: password);
}

/// Wipes the persistent state a previous run may have left behind, so the
/// next `pumpWidget(buildApp())` always starts on the SetHomeserverUrl
/// ("Connect to Server") screen. Call before pumping the app tree.
///
/// flutter_secure_storage (JWT): iOS keychain entries are scoped to the
/// bundle id and survive `simctl uninstall` -- that's documented Apple
/// behavior, not a CI quirk. Without this wipe, a JWT written by a
/// prior `flutter drive` invocation auto-logs the app in and the test
/// never sees the login screen. Linux uses installLinuxDesktopMocks for
/// isolation so the channel call is skipped there.
///
/// SharedPreferences (basePath = homeserver URL): wiping it ensures we
/// always land on the SetHomeserverUrl screen first, even when multiple
/// `testWidgets` run inside one `flutter drive` (now possible since the
/// GoRouter is built per-State -- see `lib/main.dart:_ReceiptWrangler._router`).
Future<void> resetPersistedAppState() async {
  if (!Platform.isLinux) {
    const secureChannel =
        MethodChannel('plugins.it_nomads.com/flutter_secure_storage');
    try {
      await secureChannel.invokeMethod('deleteAll', <String, dynamic>{
        'options': const <String, dynamic>{},
      });
    } catch (_) {
      // Best-effort: empty storage or unwired channel both fall through.
    }
  }
  await GlobalSharedPreferences.initialize();
  await GlobalSharedPreferences.instance.remove('basePath');
}

/// Waits for the login screen, enters [username]/[password], submits, and
/// returns once `GroupSelect` is on screen (= logged-in landing).
///
/// Split out of [loginAs] so specs that reach the login screen some other way
/// -- e.g. `login_qr_deep_link_test.dart`, which gets its server URL from a
/// deep link rather than typing it -- can reuse the credential half without
/// duplicating the locators or the landing assertion.
Future<void> loginFromLoginScreen(
  WidgetTester tester, {
  required String username,
  required String password,
}) async {
  await pumpUntilFound(tester, find.text('Log In'));

  await tester.enterText(formField('username'), username);
  await tester.enterText(formField('password'), password);
  await tester.tap(filledButton('Log In'));

  await pumpUntilFound(
    tester,
    find.byType(GroupSelect),
    timeout: const Duration(seconds: 15),
  );
}
