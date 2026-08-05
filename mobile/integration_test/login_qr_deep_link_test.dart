import 'dart:async';
import 'dart:io' show Platform;

import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/main.dart' show buildApp;

import 'helpers/env.dart';
import 'helpers/form_actions.dart';
import 'helpers/login.dart';
import 'helpers/login_qr_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

/// End-to-end coverage for the login-QR deep link
/// (`https://receiptwrangler.io/app/setup#url=<server url>`).
///
/// The point of these tests is the CROSS-COMPONENT contract: the Go API composes
/// the link (`BuildLoginQrUrl`), the desktop login page renders it as a QR, and
/// this app parses it (`extractDeepLinkServerUrl`) to pre-fill the Connect
/// screen. Go tests and Dart unit tests each assert against their own hardcoded
/// strings, so nothing pinned the two together until this spec: here the link is
/// read back off the live `/featureConfig` and fed to the app verbatim.
///
/// Deep links are injected through the `buildApp` test seams
/// (`initialDeepLink` / `deepLinkStream`) rather than the real `app_links`
/// plugin -- the test process runs ON the device and can't ask the OS to open a
/// URL, and mocking the plugin's private channels would couple the spec to its
/// internals. The seams stand in exactly where `AppLinks.getInitialLink()` /
/// `AppLinks.uriLinkStream` would be read (`lib/main.dart:_initDeepLinks`).
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  // Per-test install so each gets a fresh in-memory secure-storage map.
  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  const deepLinkPrefix = 'https://receiptwrangler.io/app/setup#url=';

  /// A single-subscription (NOT broadcast) controller on purpose: it buffers
  /// events until `_initDeepLinks` attaches its listener, which happens after an
  /// `await` on the initial-link source. A broadcast controller would silently
  /// drop an event added before that await completed.
  StreamController<Uri> deepLinkController() {
    final controller = StreamController<Uri>();
    addTearDown(controller.close);
    return controller;
  }

  testWidgets('a login-QR deep link pre-fills the server URL and connects',
      (tester) async {
    E2eEnv.assertAdmin();

    // The exact string the desktop QR would encode for this backend.
    final loginQrUrl = await enableLoginQrForTest(E2eEnv.baseUrl);
    expect(loginQrUrl, startsWith(deepLinkPrefix));

    await resetPersistedAppState();
    await tester.pumpWidget(buildApp(initialDeepLink: Uri.parse(loginQrUrl)));

    await pumpUntilFound(tester, find.text('Server URL'));

    // THE contract assertion: Go's url.QueryEscape'd fragment, parsed by Dart,
    // round-trips back to the exact URL an admin typed into System Settings.
    await pumpUntilFound(tester, find.text(E2eEnv.baseUrl));

    // Never auto-connect -- the user reviews the URL and taps Connect themselves
    // (phishing mitigation: a malicious QR must not silently point the app at
    // an attacker's server that would then harvest credentials).
    expect(find.text('Connect to Server'), findsOneWidget);
    expect(find.text('Log In'), findsNothing);

    // And the pre-filled URL is actually usable: connect + log in with it.
    await tester.tap(filledButton('Connect'));
    await loginFromLoginScreen(
      tester,
      username: E2eEnv.adminUsername,
      password: E2eEnv.adminPassword,
    );
  });

  testWidgets('a deep link delivered while the app is running pre-fills it',
      (tester) async {
    E2eEnv.assertAdmin();

    final loginQrUrl = await enableLoginQrForTest(E2eEnv.baseUrl);
    final controller = deepLinkController();

    await resetPersistedAppState();
    await tester.pumpWidget(buildApp(deepLinkStream: controller.stream));
    await pumpUntilFound(tester, find.text('Server URL'));

    // Warm delivery (app already open -- e.g. the link is tapped from another
    // app): AuthModel.pendingServerUrl notifies and the mounted Connect screen
    // patches the field.
    controller.add(Uri.parse(loginQrUrl));

    await pumpUntilFound(tester, find.text(E2eEnv.baseUrl));
    expect(find.text('Connect to Server'), findsOneWidget);
  });

  testWidgets('a deep link from another host or path is ignored',
      (tester) async {
    final controller = deepLinkController();

    await resetPersistedAppState();
    await tester.pumpWidget(buildApp(deepLinkStream: controller.stream));
    await pumpUntilFound(tester, find.text('Server URL'));

    // Only receiptwrangler.io/app/setup may pre-fill the field. Anything else is
    // dropped by _handleDeepLink, so an attacker who can get the OS to hand us a
    // look-alike link can't seed the server URL.
    controller.add(Uri.parse('https://evil.example.com/app/setup#url='
        'https%3A%2F%2Fevil.example.com%2Fapi'));
    controller.add(Uri.parse('https://receiptwrangler.io/app/other#url='
        'https%3A%2F%2Fevil.example.com%2Fapi'));
    for (var i = 0; i < 10; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }

    expect(find.text('https://evil.example.com/api'), findsNothing);
    expect(find.text('Connect to Server'), findsOneWidget);

    // Positive control: the stream IS live, so the rejections above were real
    // decisions and not a dead subscription.
    controller.add(Uri.parse(
        '${deepLinkPrefix}https%3A%2F%2Fpositive.example.com%2Fapi'));
    await pumpUntilFound(tester, find.text('https://positive.example.com/api'));
  });
}
