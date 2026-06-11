// Flow #4 -- edit an existing receipt.
//
// Login -> UI-add a receipt (returns id from URL) -> tap the view
// screen's 3-dot menu -> "Edit" -> modify the name -> save -> assert
// /view -> API GET to confirm the name actually changed -> cleanup.
//
// Own file for the GoRouter-persistence reason (see
// receipt_add_share_test.dart's header comment).

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

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('admin can edit an existing receipt', (tester) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    await loginAsAdmin(tester);

    final originalName =
        'e2e-edit-orig-${DateTime.now().millisecondsSinceEpoch}';
    final receiptId = await addManualReceiptViaUI(tester, originalName);
    scheduleReceiptCleanup(receiptId);

    // We're now on /receipts/<id>/view. The view app bar's
    // ReceiptEditPopupMenu (mobile/lib/shared/widgets/receipt_edit_popup_menu.dart:32)
    // is gated on canEdit, which is FALSE until the FutureBuilder
    // resolves and ReceiptModel is populated with the loaded
    // receipt's groupId. So pump until the PopupMenuButton actually
    // shows up before tapping.
    await pumpUntilFound(
      tester,
      find.byType(PopupMenuButton<dynamic>),
      timeout: const Duration(seconds: 10),
    );
    await tester.tap(find.byType(PopupMenuButton<dynamic>));
    // The popup scales in; the "Edit" item mounts on the animation's first
    // frame where a tap computed from its center misses (deterministic on
    // iOS: "Offset(1184.8, 99.0) ... would not hit test"). Wait for
    // hittability and drain the open animation -- same hardening as
    // receipt_cost_split_test._navigateToEdit.
    await pumpUntilFound(tester, find.text('Edit').hitTestable());
    for (int i = 0; i < 5; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }
    await tester.tap(find.text('Edit').hitTestable());
    // /edit's destination-mounted marker: the bottom save button, which only
    // renders on edit/add paths -- find.text('Name') matches on /view too,
    // so it cannot prove the navigation happened.
    await pumpUntilFound(tester, find.byType(BottomSubmitButton));

    // The submit button is built by the screen's bottomSheetBuilder even
    // while the form body is still the loading spinner (showChild false), so
    // it proves the ROUTE mounted, not the FORM. Wait for the name field
    // itself before typing -- on iOS the receipt/customFields load loses the
    // race and enterText otherwise throws "Bad state: No element".
    await pumpUntilFound(tester, formField('name'));

    // Modify the name field. enterText replaces existing content.
    final newName =
        'e2e-edit-new-${DateTime.now().millisecondsSinceEpoch}';
    await tester.enterText(formField('name'), newName);

    // No dropdown was opened on this screen so the overlay-teardown
    // settle isn't strictly needed, but keep it for symmetry with the
    // other tests -- cheap insurance.
    await tester.pumpAndSettle(const Duration(seconds: 2));

    await tester.tap(find.byType(BottomSubmitButton));
    // /view shell has mounted once ReceiptEditPopupMenu renders.
    await pumpUntilFound(tester, find.byType(ReceiptEditPopupMenu));

    // API check: the receipt's name is now the new value.
    final jwt = await apiLogin();
    final receipt = await getReceipt(receiptId, jwt: jwt);
    expect(receipt['name'], newName,
        reason: 'Receipt name should be updated to "$newName"');
  });
}
