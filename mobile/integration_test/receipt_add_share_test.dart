// Flow 3 lives in its own test file because the top-level GoRouter in
// `mobile/lib/main.dart` is a final global -- its current location
// persists across testWidgets in the same `flutter test` invocation. A
// previous test that ended at /receipts/<id>/view triggers a 403 fetch
// when the next test boots, which surfaces as `showSnackBar() during
// build` (production code at receipt_form_screen.dart:91 calls the
// snackbar in the FutureBuilder.builder). Splitting per file gives
// this test a fresh process via `mobile/run-e2e.sh`'s per-file loop.

import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_edit_popup_menu.dart';

import 'helpers/api.dart';
import 'helpers/form_actions.dart';
import 'helpers/login.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/receipt_test_helpers.dart';
import 'helpers/users.dart';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('admin can add a receipt with an item-split share',
      (tester) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    await loginAsAdmin(tester);

    final receiptName =
        'e2e-share-${DateTime.now().millisecondsSinceEpoch}';

    await tester.tap(find.text('Add'));
    await pumpUntilFound(tester, find.text('Add Manual Receipt'));
    await tester.tap(find.text('Add Manual Receipt'));
    await pumpUntilFound(tester, find.text('Name'));

    // Fill required fields. groupId MUST be set before the Add Share
    // icon button enables (mobile/lib/receipts/widgets/receipt_form.dart:465
    // disables the button when groupId == 0).
    await tester.enterText(formField('name'), receiptName);
    await tester.enterText(formField('amount'), '20.00');
    await selectDropdown(tester, 'groupId', 'My Receipts');
    await selectDropdown(tester, 'paidByUserId', adminDisplayName(tester));

    // Drain the dropdown overlay teardown so the Add Share IconButton
    // is hit-testable.
    await tester.pumpAndSettle(const Duration(seconds: 3));

    // The receipt form's "Shares" header has exactly one
    // IconButton<Icons.add>. The Details header's "Images" / "Comments"
    // compact actions use different icons.
    final addShareIcon = find.byWidgetPredicate(
      (w) =>
          w is IconButton &&
          w.icon is Icon &&
          (w.icon as Icon).icon == Icons.add,
    );
    expect(addShareIcon, findsOneWidget,
        reason: 'Exactly one IconButton<Icons.add> on the receipt form '
            '(the Add Share button)');
    // The form is taller than a Pixel-class screen at this point (Name,
    // Amount, Date, Group, Paid By, Status, Custom Fields, Categories,
    // Tags, Shares header...). On Pixel 6 Pro (411x868 logical) the
    // Shares header lands below the fold; without this scroll the tap
    // misses on a hit-test region outside the IconButton.
    await tester.ensureVisible(addShareIcon);
    await tester.pump();
    await tester.tap(addShareIcon);
    await pumpUntilFound(tester, find.text('Shared With'));

    // Share-card form. Field names ('name', 'amount') collide with the
    // main form's; both copies are in the tree, so .last picks the
    // share card's (further down in the scroll order).
    await selectDropdown(tester, 'chargedToUserId', adminDisplayName(tester));
    await tester.enterText(formField('name').last, 'Test Share');
    await tester.enterText(formField('amount').last, '5.00');

    // Card's confirm button is an ElevatedButton with "Add Share" text,
    // distinct from the IconButton that opened the card. Settle the
    // dropdown overlay and scroll the button into view before tapping.
    await tester.pumpAndSettle(const Duration(seconds: 2));
    final shareConfirm = find.widgetWithText(ElevatedButton, 'Add Share');
    await tester.ensureVisible(shareConfirm);
    await tester.tap(shareConfirm);
    await pumpUntilFound(tester, find.text('Test Share'));

    await tester.pumpAndSettle(const Duration(seconds: 3));
    await tester.tap(find.byType(BottomSubmitButton));
    // /view shell mounted -> ReceiptEditPopupMenu is in the tree.
    await pumpUntilFound(tester, find.byType(ReceiptEditPopupMenu));
    final receiptId = receiptIdFromUrl(currentUrl(tester));

    // Server-side verification: the receipt has one receiptItem named
    // "Test Share".
    final jwt = await apiLogin();
    final receipt = await getReceipt(receiptId, jwt: jwt);
    final items = receipt['receiptItems'] as List?;
    expect(items?.length ?? 0, 1,
        reason: 'Receipt should have exactly 1 receiptItem after save');
    expect(items?.first['name'], 'Test Share');

    scheduleReceiptCleanup(receiptId);
  });
}
