// The scan slot's direct actions, against a real backend.
//
// The widget suite proves the gate matrix; this proves the wiring: that a tap on
// the slot really reaches the document scanner and lands its pages in the Quick
// Scan sheet, and that a hold really produces the menu with the entries the
// caller's live permissions allow. Neither is observable from a test that
// injects a seeded PermissionsModel.

import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_entry.dart';

import 'helpers/document_scanner_mock.dart';
import 'helpers/feature_flags.dart';
import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/receipt_test_helpers.dart';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('tapping Scan captures and opens the sheet already seeded',
      (tester) async {
    await enableAiPoweredReceiptsForTest();
    // Feeds a fixed on-disk PNG through the scanner's method channel, and
    // grants camera permission on every platform (getPictures re-requests it
    // itself before invoking its native channel).
    await installDocumentScannerMock();
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    final fixture = await provisionPermUser(roleName: 'Legacy Editor');
    await loginAs(
      tester,
      username: fixture.username,
      password: fixture.password,
    );
    await enterGroup(tester, fixture.groupName!);

    expect(find.text('Scan').hitTestable(), findsWidgets,
        reason: 'a Legacy Editor with the flag on can quick-scan');

    await tester.tap(scanNavSlot());

    // No menu in between, and no separate "add a photo" step: the per-image
    // form is on screen because the captured page was carried into the sheet.
    await pumpUntilFound(tester, find.text(quickScanLabel).hitTestable());
    await pumpUntilFound(tester, find.text('Group'));
    expect(find.text('Group'), findsWidgets);
  });

  testWidgets('holding Scan offers every entry the user is permitted',
      (tester) async {
    await enableAiPoweredReceiptsForTest();
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    final fixture = await provisionPermUser(roleName: 'Legacy Editor');
    await loginAs(
      tester,
      username: fixture.username,
      password: fixture.password,
    );
    await enterGroup(tester, fixture.groupName!);

    await pumpUntilFound(tester, scanNavSlot());
    await tester.longPress(scanNavSlot());
    await pumpUntilFound(tester, find.text(quickScanLabel).hitTestable());

    expect(find.text(quickScanLabel), findsOneWidget);
    expect(find.text(addManualReceiptLabel), findsOneWidget);
    expect(find.text(uploadFromGalleryLabel), findsOneWidget);
  });

  testWidgets('the receipts screen overflow carries the same entries',
      (tester) async {
    await enableAiPoweredReceiptsForTest();
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    final fixture = await provisionPermUser(roleName: 'Legacy Editor');
    await loginAs(
      tester,
      username: fixture.username,
      password: fixture.password,
    );
    await enterGroup(tester, fixture.groupName!);
    await tester.tap(find.text('Receipts').hitTestable());

    // The accessible route to what the hold offers, for anyone who never
    // discovers the gesture.
    final overflow = find.byKey(const ValueKey('receipt-entry-overflow-menu'));
    await pumpUntilFound(tester, overflow.hitTestable());
    await tester.tap(overflow.hitTestable());
    await pumpUntilFound(tester, find.text(addManualReceiptLabel).hitTestable());

    expect(find.text(addManualReceiptLabel), findsOneWidget);
    expect(find.text(uploadFromGalleryLabel), findsOneWidget);
  });
}
