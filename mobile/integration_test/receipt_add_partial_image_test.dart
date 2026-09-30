// Creating a receipt is ONE atomic call (`POST /receipt/withFiles`) carrying
// the receipt, its comments and its staged images.
//
// It used to be a JSON create followed by one upload per image, so a failed
// upload left a receipt with some or none of its images and the app showed a
// "Receipt added, but one or more images failed to upload" snackbar. That
// state no longer exists: either everything is stored or nothing is.
//
// This spec pins both halves against the real server:
//  1. A file the server rejects (text bytes behind a .png name -- the server
//     sniffs content, not the name) fails the create with nothing written: no
//     receipt, the user stays on the add form with an error.
//  2. A valid image lands with the receipt in the same call.
//
// Runs on every target: it drives the **file** source, which
// `installFileSelectorMock` intercepts by swapping the platform interface.

import 'dart:convert' show utf8;
import 'dart:io' show Platform;
import 'dart:typed_data';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_edit_popup_menu.dart';

import 'helpers/api.dart';
import 'helpers/file_selector_mock.dart';
import 'helpers/form_actions.dart';
import 'helpers/login.dart';
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

  /// Logs in as a fresh Legacy Editor in a fixture group (whose role requires
  /// nothing) and fills the add form up to the images. Returns the fixture.
  Future<PermFixture> openFilledAddForm(
    WidgetTester tester,
    String receiptName,
  ) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    final fixture = await provisionPermUser(roleName: 'Legacy Editor');
    await loginAs(
      tester,
      username: fixture.username,
      password: fixture.password,
    );

    await openManualReceiptForm(tester);
    await attachFileFromReceiptForm(tester);

    await tester.enterText(formField('name'), receiptName);
    await tester.enterText(formField('amount'), '12.34');
    await selectDropdown(tester, 'groupId', fixture.groupName!);
    await selectDropdown(tester, 'paidByUserId', fixture.displayName);
    await tester.pumpAndSettle(const Duration(seconds: 3));
    return fixture;
  }

  testWidgets('a rejected file creates nothing -- the create is atomic', (
    tester,
  ) async {
    await installFileSelectorMock(
      bytes: Uint8List.fromList(utf8.encode('this is not an image')),
      name: 'not-an-image.png',
    );
    final receiptName =
        'e2e-atomic-fail-${DateTime.now().millisecondsSinceEpoch}';

    final fixture = await openFilledAddForm(tester, receiptName);
    await tester.tap(find.byType(BottomSubmitButton));

    // The server's 400 surfaces as an error snackbar...
    await pumpUntilFound(
      tester,
      find.byType(SnackBar),
      timeout: const Duration(seconds: 15),
    );
    // ...and the user is still on the add form, not a receipt view.
    expect(find.byType(BottomSubmitButton), findsOneWidget);
    expect(find.byType(ReceiptEditPopupMenu), findsNothing);
    expect(find.textContaining('Receipt added'), findsNothing);

    // Server-side: no receipt by that name exists in the group.
    final receipts = await listReceiptsForGroup(
      fixture.groupId!,
      jwt: await apiLogin(),
    );
    expect(
      receipts.where((r) => r['name'] == receiptName),
      isEmpty,
      reason: 'a failed create must leave no receipt behind',
    );
  });

  testWidgets('a valid image is stored with the receipt in the same call', (
    tester,
  ) async {
    await installFileSelectorMock();
    final receiptName =
        'e2e-atomic-ok-${DateTime.now().millisecondsSinceEpoch}';

    await openFilledAddForm(tester, receiptName);
    final receiptId = await submitManualReceiptForm(tester);
    scheduleReceiptCleanup(receiptId);

    final receipt = await getReceipt(receiptId, jwt: await apiLogin());
    expect(receipt['name'], receiptName);
    expect((receipt['imageFiles'] as List?)?.length ?? 0, 1);
  });
}
