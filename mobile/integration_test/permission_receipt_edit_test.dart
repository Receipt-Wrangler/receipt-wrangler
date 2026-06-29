import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_edit_popup_menu.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/slidable_widget.dart';

import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

/// Permission-gating coverage for editing a receipt, both gated on
/// `group.receipts.update`:
///   - the swipe-to-edit action on the receipts list (`receipt_list_item.dart`,
///     `SlidableWidget.slideEnabled`), and
///   - the edit popup menu on the receipt view (`receipt_edit_popup_menu.dart`,
///     a `PopupMenuButton` mounted by `receipt_app_bar_action_builder.dart`).
///
/// A Legacy Viewer (read but not update) should get neither; a Legacy Editor
/// (update) should get both. We assert `SlidableWidget.slideEnabled` directly
/// rather than performing a drag -- a slidable drag in integration_test is
/// velocity/timing fragile, and the bool is the exact production input.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  Finder editPopupButton() => find.descendant(
        of: find.byType(ReceiptEditPopupMenu),
        matching: find.byType(PopupMenuButton<dynamic>),
      );

  SlidableWidget receiptSlidable(WidgetTester tester, String receiptName) =>
      tester.widget<SlidableWidget>(
        find.ancestor(
          of: find.text(receiptName),
          matching: find.byType(SlidableWidget),
        ),
      );

  testWidgets(
    'Legacy Viewer: receipt swipe-edit disabled and edit popup hidden',
    (tester) async {
      final fixture = await provisionPermUser(
        roleName: 'Legacy Viewer',
        withReceipt: true,
      );
      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );
      await openGroupReceipts(tester, fixture.groupName!, fixture.receiptName!);

      // Gate: swipe-to-edit is disabled for a Viewer.
      expect(receiptSlidable(tester, fixture.receiptName!).slideEnabled, isFalse);

      // Gate: opening the receipt view, the edit popup renders SizedBox.shrink
      // (the ReceiptEditPopupMenu widget mounts, but with no PopupMenuButton).
      await tester.tap(find.text(fixture.receiptName!));
      await pumpUntilFound(tester, find.byType(ReceiptEditPopupMenu));
      expect(editPopupButton(), findsNothing);
    },
  );

  testWidgets(
    'Legacy Editor: receipt swipe-edit enabled and edit popup shown',
    (tester) async {
      final fixture = await provisionPermUser(
        roleName: 'Legacy Editor',
        withReceipt: true,
      );
      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );
      await openGroupReceipts(tester, fixture.groupName!, fixture.receiptName!);

      // Gate: swipe-to-edit is enabled for an Editor.
      expect(receiptSlidable(tester, fixture.receiptName!).slideEnabled, isTrue);

      // Gate: opening the receipt view, the edit popup button is present.
      await tester.tap(find.text(fixture.receiptName!));
      await pumpUntilFound(tester, find.byType(ReceiptEditPopupMenu));
      await pumpUntilFound(tester, editPopupButton());
      expect(editPopupButton(), findsOneWidget);
    },
  );
}
