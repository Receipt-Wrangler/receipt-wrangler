// Flow A -- verifies bug-fix #1 (double-submit guard).
//
// The bug: rapid Save taps used to fire onPressed twice, producing
// duplicate receipts. The fix wraps the body in
// LoadingModel.setIsLoading(true) at the start, which causes
// BottomSubmitButton's Consumer<LoadingModel> to rebuild with
// onPressed: null. The second tap then has nothing to fire.
//
// This test fires two `tester.tap` calls back-to-back. The internal
// pump between them is enough for the LoadingModel rebuild to swap
// onPressed to null; the second tap "misses" with a warning, which
// we suppress via warnIfMissed: false. Server-side, exactly one
// receipt with the unique name should exist.
//
// In its own file so run-e2e.sh's per-file invocation gives it a
// fresh process (GoRouter URL persistence avoidance, same reason as
// receipt_add_share_test.dart).

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

  testWidgets('rapid double-tap of Submit creates exactly one receipt',
      (tester) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    await loginAsAdmin(tester);

    final receiptName =
        'e2e-double-${DateTime.now().millisecondsSinceEpoch}';

    // Replicates flow 1's add-receipt fill sequence inline. Can't reuse
    // addManualReceiptViaUI because that helper taps Submit only once
    // and waits for navigation; this test needs to inject a SECOND tap
    // before the wait.
    await tester.tap(find.text('Add'));
    await pumpUntilFound(tester, find.text('Add Manual Receipt'));
    await tester.tap(find.text('Add Manual Receipt'));
    await pumpUntilFound(tester, find.text('Name'));

    await tester.enterText(formField('name'), receiptName);
    await tester.enterText(formField('amount'), '12.34');
    await selectDropdown(tester, 'groupId', 'My Receipts');
    await selectDropdown(tester, 'paidByUserId', adminDisplayName(tester));

    // Drain the dropdown overlay teardown so BottomSubmitButton is
    // hit-testable.
    await tester.pumpAndSettle(const Duration(seconds: 3));

    final submitFinder = find.byType(BottomSubmitButton);

    // First tap: fires onPressed -> setIsLoading(true) ->
    // notifyListeners -> Consumer rebuilds button with onPressed: null.
    await tester.tap(submitFinder);

    // Second tap immediately: button is now disabled. This tap should
    // miss; warnIfMissed: false because the miss is the *desired*
    // behavior, not a test failure.
    await tester.tap(submitFinder, warnIfMissed: false);

    // One navigation should fire (from the first tap's addReceipt).
    // /view shell mounted -> ReceiptEditPopupMenu is in the tree.
    await pumpUntilFound(tester, find.byType(ReceiptEditPopupMenu));
    final firstId = receiptIdFromUrl(currentUrl(tester));
    // Schedule cleanup BEFORE assertions so it fires even if the
    // assertion below blows up (otherwise leftover receipts pile up
    // on the demo backend across failed runs).
    scheduleReceiptCleanup(firstId);

    // Server-side check: exactly one receipt with our unique name.
    // GET the just-created receipt to find its group, then list that
    // group's recent receipts and filter by name.
    final jwt = await apiLogin();
    final receipt = await getReceipt(firstId, jwt: jwt);
    final groupId = receipt['groupId'] as int;
    final recent = await listReceiptsForGroup(groupId, jwt: jwt);
    final matching =
        recent.where((r) => r['name'] == receiptName).toList();

    // Defensive: if the bug regressed and a duplicate exists, schedule
    // cleanup for it before the assertion blows up. addTearDown is
    // LIFO so these fire before the firstId cleanup.
    for (final dup in matching) {
      final id = dup['id'] as int;
      if (id != firstId) scheduleReceiptCleanup(id);
    }
    expect(
      matching.length,
      1,
      reason: 'Double-tap should produce exactly one receipt named '
          '"$receiptName"; found ${matching.length}: '
          '${matching.map((r) => r['id']).toList()}',
    );
  });
}
