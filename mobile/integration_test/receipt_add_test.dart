import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_edit_popup_menu.dart';

import 'helpers/form_actions.dart';
import 'helpers/login.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/receipt_test_helpers.dart';
import 'helpers/users.dart';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  // Per-test install so each gets a fresh in-memory secure-storage map
  // (a leaked JWT from the previous test would skip the login screens).
  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('admin can add a manual receipt', (tester) async {
    // The Linux desktop test window defaults to 1280x720 -- too short
    // to render the receipt form's persistent bottom sheet (Submit
    // button) inside the visible viewport. 1280x900 keeps the entire
    // form in view while staying close to a tablet-sized layout.
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    await loginAsAdmin(tester);

    final receiptName =
        'e2e-manual-${DateTime.now().millisecondsSinceEpoch}';

    // Open the bottom-nav Add popup menu and pick "Add Manual Receipt".
    await tester.tap(find.text('Add'));
    await pumpUntilFound(tester, find.text('Add Manual Receipt'));
    await tester.tap(find.text('Add Manual Receipt'));
    await pumpUntilFound(tester, find.text('Name'));

    // Fill required fields. Date defaults to now and status defaults to
    // OPEN via getDefaultReceipt() (mobile/lib/utils/receipts.dart:16),
    // so both pass validation without user interaction.
    await tester.enterText(formField('name'), receiptName);
    await tester.enterText(formField('amount'), '12.34');
    await selectDropdown(tester, 'groupId', 'My Receipts');
    // The admin's displayName from signup is 'ee', not 'e2e-admin'.
    await selectDropdown(tester, 'paidByUserId', adminDisplayName(tester));

    // Drain the dropdown overlay teardown -- the popup-route's overlay
    // entry can otherwise leave the Scaffold's bottom-sheet area in an
    // Offstage state and the BottomSubmitButton tap silently misses.
    await tester.pumpAndSettle(const Duration(seconds: 3));

    final submitFinder = find.byType(BottomSubmitButton);
    expect(submitFinder, findsOneWidget,
        reason: 'BottomSubmitButton should be rendered on /receipts/add');
    await tester.tap(submitFinder);
    // Assert /view via the ReceiptEditPopupMenu (only mounted on /view
    // per receipt_app_bar_action_builder.dart:56-67), then read the
    // id off the URL for cleanup.
    await pumpUntilFound(tester, find.byType(ReceiptEditPopupMenu));
    scheduleReceiptCleanup(receiptIdFromUrl(currentUrl(tester)));
  });
}
