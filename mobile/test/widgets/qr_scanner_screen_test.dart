import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/auth/set-homeserver-url/screens/qr_scanner_screen.dart';

void main() {
  // These cover the non-camera fallback states only. The live-camera path
  // constructs MobileScanner / starts the controller, which hits platform
  // channels unavailable under flutter_test — it's exercised on device /
  // integration tests. The lazy controller + debug seams let the fallback
  // states render here without any camera or channel mocks.

  Future<void> pumpScreen(WidgetTester tester, Widget screen) async {
    await tester.pumpWidget(MaterialApp(home: screen));
    await tester.pump();
  }

  testWidgets('shows an unsupported message when the scanner is unavailable',
      (tester) async {
    await pumpScreen(tester, const QrScannerScreen(debugScannerSupported: false));

    expect(
      find.text("QR scanning isn't supported on this device."),
      findsOneWidget,
    );
  });

  testWidgets('permission-denied state shows the message and Open Settings',
      (tester) async {
    await pumpScreen(
      tester,
      const QrScannerScreen(debugForcePermissionDenied: true),
    );

    expect(find.textContaining('Camera access is required'), findsOneWidget);
    expect(
      find.widgetWithText(ElevatedButton, 'Open Settings'),
      findsOneWidget,
    );
  });

  testWidgets('camera-error state shows the message and Retry', (tester) async {
    await pumpScreen(
      tester,
      const QrScannerScreen(debugForceCameraError: true),
    );

    expect(find.textContaining("Couldn't start the camera"), findsOneWidget);
    expect(find.widgetWithText(ElevatedButton, 'Retry'), findsOneWidget);
  });

  testWidgets('close button pops the screen', (tester) async {
    await tester.pumpWidget(
      MaterialApp(
        home: Builder(
          builder: (context) => Scaffold(
            body: Center(
              child: ElevatedButton(
                onPressed: () => Navigator.of(context).push(
                  MaterialPageRoute<void>(
                    builder: (_) =>
                        const QrScannerScreen(debugScannerSupported: false),
                  ),
                ),
                child: const Text('open scanner'),
              ),
            ),
          ),
        ),
      ),
    );

    await tester.tap(find.text('open scanner'));
    await tester.pumpAndSettle();
    expect(find.byKey(const ValueKey('qr-scanner-close')), findsOneWidget);

    await tester.tap(find.byKey(const ValueKey('qr-scanner-close')));
    await tester.pumpAndSettle();

    // Back on the launcher screen; the scanner is gone.
    expect(find.byKey(const ValueKey('qr-scanner-close')), findsNothing);
    expect(find.text('open scanner'), findsOneWidget);
  });
}
