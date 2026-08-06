import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/auth/set-homeserver-url/screens/set_homeserver_url.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/persistence/global_shared_preferences.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  const qrButtonKey = ValueKey('qr-scan-button');
  const invalidMessage = "That QR code doesn't contain a valid server URL";

  setUp(() async {
    // The screen reads AuthModel.basePath -> GlobalSharedPreferences on build.
    SharedPreferences.setMockInitialValues({});
    await GlobalSharedPreferences.initialize();
  });

  // The injected scanQrCode seam lets us drive the post-scan flow without
  // constructing MobileScanner (which throws under flutter_test).
  Future<void> pumpScreen(
    WidgetTester tester,
    Future<String?> Function(BuildContext context) scanQrCode,
  ) async {
    await tester.pumpWidget(
      ChangeNotifierProvider<AuthModel>(
        create: (_) => AuthModel(),
        child: MaterialApp(
          home: Scaffold(
            body: SetHomeserverUrl(scanQrCode: scanQrCode),
          ),
        ),
      ),
    );
    await tester.pump();
  }

  Future<void> tapScanAndSettle(WidgetTester tester) async {
    await tester.tap(find.byKey(qrButtonKey));
    await tester.pump(); // resume the async scan handler
    await tester.pump(const Duration(milliseconds: 400)); // render patch/snackbar
  }

  testWidgets('renders the QR scan button', (tester) async {
    await pumpScreen(tester, (_) async => null);

    expect(find.byKey(qrButtonKey), findsOneWidget);
  });

  testWidgets('a valid scan populates the URL field (trimmed)', (tester) async {
    await pumpScreen(tester, (_) async => '  https://scanned.example.io/api  ');

    await tapScanAndSettle(tester);

    expect(find.text('https://scanned.example.io/api'), findsOneWidget);
    expect(find.text(invalidMessage), findsNothing);
  });

  testWidgets('a deep-link QR scan fills the field with the decoded server URL',
      (tester) async {
    // The same receiptwrangler.io/app/setup#url=... link the OS uses to open the
    // app is also readable by the in-app scanner: _onScanPressed resolves it via
    // extractDeepLinkServerUrl before falling back to normalizeServerUrl.
    await pumpScreen(
      tester,
      (_) async =>
          'https://receiptwrangler.io/app/setup#url=https%3A%2F%2Fscanned.example.io%2Fapi',
    );

    await tapScanAndSettle(tester);

    expect(find.text('https://scanned.example.io/api'), findsOneWidget);
    expect(find.text(invalidMessage), findsNothing);
  });

  testWidgets('an invalid scan shows an error and leaves the field empty',
      (tester) async {
    await pumpScreen(tester, (_) async => 'not a url');

    await tapScanAndSettle(tester);

    expect(find.text(invalidMessage), findsOneWidget);
    expect(find.text('not a url'), findsNothing);
  });

  testWidgets('a cancelled scan (null) is a no-op', (tester) async {
    await pumpScreen(tester, (_) async => null);

    await tapScanAndSettle(tester);

    expect(find.text(invalidMessage), findsNothing);
  });

  testWidgets(
      'AuthModel.pendingServerUrl pre-fills the field without auto-connecting',
      (tester) async {
    // Use a test-owned model so we can push the pending URL onto it, mimicking
    // the deep-link handler in main.dart (warm case: the screen is already
    // mounted when the value arrives).
    final authModel = AuthModel();
    await tester.pumpWidget(
      ChangeNotifierProvider<AuthModel>.value(
        value: authModel,
        child: const MaterialApp(
          home: Scaffold(body: SetHomeserverUrl()),
        ),
      ),
    );
    await tester.pump();

    authModel.setPendingServerUrl('https://deeplink.example.io/api');
    await tester.pump(); // rebuild (listener) -> schedule post-frame patch
    await tester.pump(const Duration(milliseconds: 400)); // run patch + clear

    // Field is pre-filled from the pending URL.
    expect(find.text('https://deeplink.example.io/api'), findsOneWidget);
    // Consumed exactly once (cleared), so a later rebuild won't re-apply it.
    expect(authModel.pendingServerUrl, isNull);
    // No auto-connect: still on the Connect screen (never navigated away).
    expect(find.text('Connect to Server'), findsOneWidget);
  });
}
