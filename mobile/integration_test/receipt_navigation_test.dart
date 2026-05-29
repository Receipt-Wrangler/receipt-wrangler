import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/receipt_list_item.dart';

import 'helpers/login.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/receipt_test_helpers.dart';

/// Switches between two receipts via the receipt list and verifies each
/// /view screen shows the *current* receipt's data with no stale residue
/// from the previous one.
///
/// Specifically pins down the post-fix behavior of
/// `ReceiptModel.setReceipt()` (`receipt_model.dart:69`): when the receipt
/// id actually changes the form key MUST regenerate so the FormBuilder
/// remounts with the new initialValues, and `resetModel()` (called from
/// the app-bar back arrow at `receipt_app_bar.dart:44-46`) MUST clear the
/// previous receipt's data so it doesn't bleed into the list-flash or the
/// next /view.
///
/// Catches:
///   - A `setReceipt` regression that fails to swap form values when the id changes.
///   - A `resetModel` regression that leaves stale name/amount in the model.
///   - Any `late final`-captured field in the form tree that doesn't refresh.
void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets(
      'admin can switch between two receipts via the list without stale data',
      (tester) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    await loginAsAdmin(tester);

    // Distinct names per receipt -- timestamped + A/B suffix so a fast
    // double-create within the same millisecond can't accidentally produce
    // identical strings.
    final ts = DateTime.now().millisecondsSinceEpoch;
    final nameA = 'e2e-nav-A-$ts';
    final nameB = 'e2e-nav-B-$ts';

    // Receipt A: created via UI, lands us on /receipts/<id1>/view.
    // `addManualReceiptViaUI` returns as soon as the URL flips to /view,
    // but the form's FutureBuilder still has to resolve the receipt
    // before the name renders (ReceiptFormScreen shows a spinner in the
    // interim). pumpUntilFound the name so the assertion doesn't race.
    final id1 =
        await addManualReceiptViaUI(tester, nameA, amount: '10.00');
    scheduleReceiptCleanup(id1);
    await pumpUntilFound(tester, find.text(nameA));

    // Back to the list. The back arrow on /view is the TopAppBar's leading
    // IconButton; it triggers `receiptModel.resetModel()` and pushes
    // /groups/<groupId>/receipts (receipt_app_bar.dart:42-46). Wait for
    // the bottom-nav "Add" entry as the destination-mounted marker
    // (the same target addManualReceiptViaUI taps next).
    await _tapBackArrow(tester);
    await pumpUntilFound(tester, find.text('Add'));

    // Receipt B: created from the list. The bottom-nav "Add" entry is
    // present on both /groups and /groups/<id>/receipts (group_select_bottom_nav.dart
    // and group_bottom_nav.dart both declare it), so `addManualReceiptViaUI`'s
    // pre-condition is satisfied from the list too.
    final id2 =
        await addManualReceiptViaUI(tester, nameB, amount: '99.99');
    scheduleReceiptCleanup(id2);

    // Positive find first to make pumpUntilFound wait for the new data,
    // then the negative for receipt A. Doing them in this order avoids a
    // false-clean from asserting absence before the screen has redrawn.
    await pumpUntilFound(tester, find.text(nameB));
    expect(find.text(nameA).evaluate(), isEmpty,
        reason: 'receipt B\'s view must not still be rendering A\'s name '
            '-- that\'d indicate setReceipt didn\'t propagate or some '
            'late-final cached the prior receipt');

    // Back to the list, then switch from B's view back to A. The
    // widgetWithText finder on ReceiptListItem doubles as the
    // destination-mounted check -- if it's there, the list shell has
    // built and the item is hit-testable.
    await _tapBackArrow(tester);
    await pumpUntilFound(
        tester, find.widgetWithText(ReceiptListItem, nameA));

    await tester.tap(find.widgetWithText(ReceiptListItem, nameA));
    await pumpUntilFound(tester, find.text(nameA));
    expect(find.text(nameB).evaluate(), isEmpty,
        reason: 'after navigating list -> A, B\'s name must not linger');

    // Final hop: A -> list -> B.
    await _tapBackArrow(tester);
    await pumpUntilFound(
        tester, find.widgetWithText(ReceiptListItem, nameB));

    await tester.tap(find.widgetWithText(ReceiptListItem, nameB));
    await pumpUntilFound(tester, find.text(nameB));
    expect(find.text(nameA).evaluate(), isEmpty,
        reason: 'after the final list -> B switch, A\'s name must not linger');
  });
}

/// Taps the AppBar's back arrow on the receipt /view screen. Scoped to a
/// descendant of `AppBar` to dodge the (unlikely) case where another
/// Icons.arrow_back lives elsewhere in the tree.
Future<void> _tapBackArrow(WidgetTester tester) async {
  final back = find.descendant(
    of: find.byType(AppBar),
    matching: find.byIcon(Icons.arrow_back),
  );
  await pumpUntilFound(tester, back);
  await tester.tap(back);
}
