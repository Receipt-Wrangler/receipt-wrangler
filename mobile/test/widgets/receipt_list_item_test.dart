import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/receipt_list_item.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/system_settings_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/slidable_widget.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';

import '../helpers/permission_test_helpers.dart';
import '../helpers/widget_test_helpers.dart';

/// Widget-level coverage for the receipt-list swipe-to-edit gate
/// (`receipt_list_item.dart` → `SlidableWidget.slideEnabled` =
/// `canEditReceipt`, i.e. `group.receipts.update`). Complements the e2e in
/// `integration_test/permission_receipt_edit_test.dart` with a fast,
/// deterministic check on the exact production input.
void main() {
  const groupId = 7;

  // The list item renders the receipt amount via formatCurrency, which needs the
  // app's custom currency registered in money2's process-wide registry.
  setUpAll(registerCustomCurrencyForTests);

  api.Receipt receiptInGroup() => getDefaultReceipt().rebuild((b) => b
    ..id = 1
    ..groupId = groupId
    ..paidByUserId = 1
    ..name = 'Test Receipt'
    ..amount = '12.34'
    ..createdAt = DateTime.now().toIso8601String());

  UserModel userModelWithTester() {
    final model = UserModel();
    model.setUsers([
      (api.UserViewBuilder()
            ..id = 1
            ..username = 'tester'
            ..displayName = 'Tester'
            ..isDummyUser = false)
          .build(),
    ]);
    return model;
  }

  Widget wrap(List<Permission> groupPermissions) {
    return MultiProvider(
      providers: [
        ChangeNotifierProvider.value(
          value: seededPermissions(group: {groupId: groupPermissions}),
        ),
        ChangeNotifierProvider.value(value: GroupModel()),
        ChangeNotifierProvider.value(value: userModelWithTester()),
        ChangeNotifierProvider.value(value: SystemSettingsModel()),
      ],
      child: MaterialApp(
        home: Scaffold(body: ReceiptListItem(receipt: receiptInGroup())),
      ),
    );
  }

  Future<void> pumpItem(
      WidgetTester tester, List<Permission> groupPermissions) async {
    await tester.pumpWidget(wrap(groupPermissions));
    // ListItemLead packs the formatted date into a fixed 50px lead, which
    // overflows the headless test canvas by a few px (it renders fine on a real
    // device). That's incidental to the swipe-to-edit gate under test, so
    // consume the layout exception rather than coupling this test to that
    // sub-widget's pixel layout.
    tester.takeException();
  }

  SlidableWidget slidable(WidgetTester tester) =>
      tester.widget<SlidableWidget>(find.byType(SlidableWidget));

  group('ReceiptListItem swipe-to-edit gate', () {
    testWidgets('enabled with group.receipts.update', (tester) async {
      await pumpItem(tester, [Permission.groupPeriodReceiptsPeriodUpdate]);

      expect(slidable(tester).slideEnabled, isTrue);
    });

    testWidgets('disabled without group.receipts.update', (tester) async {
      await pumpItem(tester, [Permission.groupPeriodReceiptsPeriodRead]);

      expect(slidable(tester).slideEnabled, isFalse);
    });
  });
}
