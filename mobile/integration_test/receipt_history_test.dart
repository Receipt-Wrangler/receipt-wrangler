import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/receipts/nav/receipt_app_bar_action_builder.dart';
import 'package:receipt_wrangler_mobile/receipts/screens/receipt_history_screen.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_edit_popup_menu.dart';

import 'helpers/api.dart';
import 'helpers/form_actions.dart';
import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

/// The receipt history screen against a real API.
///
/// The widget suite mocks the history response, so only this proves the menu
/// entry reaches the real endpoint and renders what the server computed. It
/// covers the three things a unit test cannot:
///   - a read-only member holding `app.receipt-history.read` gets a view menu
///     with only History, and the screen lists the server's entries, including
///     a comment (which recorded nothing before the history feature);
///   - History opened from the edit form leaves the form's unsaved input in
///     place -- it is pushed, not routed, because the form reads view/edit from
///     the URL;
///   - without the permission there is no History entry.
///
/// The permission is granted by giving the provisioned user the Legacy Admin
/// app role, which holds every app permission; their group role decides
/// whether they can also edit.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  final historyItem = find.byKey(ReceiptAppBarActionBuilder.historyMenuItemKey);

  final viewMenuButton = find.descendant(
    of: find.byType(ReceiptEditPopupMenu),
    matching: find.byType(PopupMenuButton<dynamic>),
  );

  Future<int> legacyAdminAppRoleId() async =>
      appRoleIdByName('Legacy Admin', jwt: await apiLogin());

  // A popup mounts its items on the open animation's first frame, where a tap
  // computed from an item's centre misses: wait for hittability, then drain.
  Future<void> openMenu(WidgetTester tester, Finder button, Finder item) async {
    await pumpUntilFound(tester, button);
    await tester.tap(button);
    await pumpUntilFound(tester, item.hitTestable());
    for (var i = 0; i < 5; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }
  }

  testWidgets(
    'a read-only member with the permission opens History from the view menu',
    (tester) async {
      final fixture = await provisionPermUser(
        roleName: 'Legacy Viewer',
        appRoleId: await legacyAdminAppRoleId(),
        withReceipt: true,
      );
      const comment = 'e2e history comment';
      await createComment(
        receiptId: fixture.receiptId!,
        jwt: await apiLoginAs(fixture.username, fixture.password),
        comment: comment,
      );

      await loginAs(tester, username: fixture.username, password: fixture.password);
      await openGroupReceipts(tester, fixture.groupName!, fixture.receiptName!);
      await tester.tap(find.text(fixture.receiptName!));
      await pumpUntilFound(tester, find.byType(ReceiptEditPopupMenu));

      await openMenu(tester, viewMenuButton, historyItem);
      // A Viewer cannot edit, so History is the menu's only entry.
      expect(find.text('Edit'), findsNothing);

      await tester.tap(historyItem.hitTestable());
      await pumpUntilFound(tester, find.byType(ReceiptHistoryScreen));
      await pumpUntilFound(
          tester, find.textContaining('created this receipt', findRichText: true));
      expect(
        find.textContaining('updated images and comments', findRichText: true),
        findsOneWidget,
      );
      expect(find.textContaining(comment), findsOneWidget);
    },
  );

  testWidgets(
    'History opened from the edit form leaves the unsaved edit in place',
    (tester) async {
      final fixture = await provisionPermUser(
        roleName: 'Legacy Editor',
        appRoleId: await legacyAdminAppRoleId(),
        withReceipt: true,
      );

      await loginAs(tester, username: fixture.username, password: fixture.password);
      await openReceiptEditForm(tester, fixture.groupName!, fixture.receiptName!);

      final unsavedName = 'e2e-history-unsaved-${DateTime.now().millisecondsSinceEpoch}';
      await tester.enterText(formField('name'), unsavedName);
      await tester.pump();

      await openMenu(
        tester,
        find.byKey(ReceiptAppBarActionBuilder.editOverflowMenuKey),
        historyItem,
      );
      await tester.tap(historyItem.hitTestable());
      await pumpUntilFound(
          tester, find.textContaining('created this receipt', findRichText: true));

      await tester.tap(find.byTooltip('Back'));
      await pumpUntilGone(tester, find.byType(ReceiptHistoryScreen));

      // Still the edit form -- not flipped to view -- and still holding the edit.
      await pumpUntilFound(tester, find.byType(BottomSubmitButton));
      expect(find.byType(ReceiptEditPopupMenu), findsNothing);
      expect(find.text(unsavedName), findsOneWidget);
    },
  );

  testWidgets('without the permission the menu offers no History', (tester) async {
    // Default app role: Legacy User, which does not hold the permission.
    final fixture = await provisionPermUser(
      roleName: 'Legacy Editor',
      withReceipt: true,
    );

    await loginAs(tester, username: fixture.username, password: fixture.password);
    await openGroupReceipts(tester, fixture.groupName!, fixture.receiptName!);
    await tester.tap(find.text(fixture.receiptName!));
    await pumpUntilFound(tester, find.byType(ReceiptEditPopupMenu));

    await openMenu(tester, viewMenuButton, find.text('Edit'));
    expect(historyItem, findsNothing);
  });
}
