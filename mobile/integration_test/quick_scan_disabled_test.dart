// Verifies that the bottom-nav Add menu hides "Quick Scan" when
// authModel.featureConfig.aiPoweredReceipts is false.
//
// Regression for show_add_menu.dart: previously the Quick Scan
// PopupMenuItem was always rendered, regardless of the AI feature
// flag, and only the in-action showQuickScanBottomSheet enforced the
// gate (by firing an error snackbar instead of opening the sheet).
// The fix moves the gate up to the menu itself.
//
// Strategy: log in normally, then mutate AuthModel.featureConfig
// directly via Provider before tapping Add. This is closer to the
// real failure mode than mocking /featureConfig because the menu
// reads featureConfig synchronously at tap time -- which IS the
// decision point. After login, nothing else writes featureConfig
// until the next storeAppData call, so the mutation sticks.

import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';

import 'helpers/login.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('Add menu hides Quick Scan when aiPoweredReceipts is false',
      (tester) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    await loginAsAdmin(tester);

    // Force the feature flag off via the live AuthModel.
    final ctx = tester.element(find.byType(Scaffold).first);
    final authModel = Provider.of<AuthModel>(ctx, listen: false);
    final originalEnableLocalSignUp =
        authModel.featureConfig.enableLocalSignUp;

    authModel.setFeatureConfig(
      (api.FeatureConfigBuilder()
            ..aiPoweredReceipts = false
            ..enableLocalSignUp = originalEnableLocalSignUp)
          .build(),
    );
    await tester.pump();

    // Open the bottom-nav Add menu.
    await tester.tap(find.text('Add'));
    // Manual receipt entry is unconditional; finding it confirms the
    // menu actually opened (vs. the tap missing).
    await pumpUntilFound(tester, find.text('Add Manual Receipt'));

    expect(find.text('Quick Scan'), findsNothing,
        reason: 'Quick Scan must not appear when aiPoweredReceipts=false');

    // Dismiss the popup before reopening with the flag flipped on.
    await tester.tapAt(const Offset(10, 10));
    await tester.pumpAndSettle(const Duration(milliseconds: 300));

    authModel.setFeatureConfig(
      (api.FeatureConfigBuilder()
            ..aiPoweredReceipts = true
            ..enableLocalSignUp = originalEnableLocalSignUp)
          .build(),
    );
    await tester.pump();

    await tester.tap(find.text('Add'));
    await pumpUntilFound(tester, find.text('Add Manual Receipt'));
    expect(find.text('Quick Scan'), findsOneWidget,
        reason: 'Quick Scan must reappear when aiPoweredReceipts=true');
  });
}
