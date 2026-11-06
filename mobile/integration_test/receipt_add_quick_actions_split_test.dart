import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_form.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_edit_popup_menu.dart';

import 'helpers/api.dart';
import 'helpers/form_actions.dart';
import 'helpers/login.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/quick_actions_actions.dart';
import 'helpers/receipt_test_helpers.dart';
import 'helpers/users.dart';

/// Quick Actions -> Split Evenly on the ADD form, against the real API.
///
/// `receipt_cost_split_test.dart` covers the same split on the EDIT form only.
void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('add form: fill, Split Evenly between 2, submit creates it', (
    tester,
  ) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    final group = await provisionSharedSplitGroup();
    await loginAsAdmin(tester);

    final receiptName =
        'e2e-add-split-${DateTime.now().millisecondsSinceEpoch}';

    await openManualReceiptForm(tester);
    await tester.enterText(formField('name'), receiptName);
    await tester.enterText(formField('amount'), '100.00');
    await selectDropdown(tester, 'groupId', group.name);
    await selectDropdown(tester, 'paidByUserId', adminDisplayName(tester));
    await tester.pumpAndSettle(const Duration(seconds: 3));

    await openQuickActionsSheet(tester);
    await selectQuickActionsUsers(tester, [
      adminDisplayName(tester),
      userDisplayName(tester),
    ]);
    await pumpUntilFound(tester, find.textContaining('2 users'));

    final splitButton = find.widgetWithText(BottomSubmitButton, 'Split');
    await pumpUntilFound(tester, splitButton.hitTestable());
    await tester.tap(splitButton.hitTestable());
    await pumpUntilFound(tester, find.text('Name'));
    await tester.pumpAndSettle(const Duration(seconds: 2));

    final saveButton = find.byType(BottomSubmitButton);
    await pumpUntilFound(tester, saveButton.hitTestable());
    await tester.tap(saveButton.hitTestable());

    // The /view shell mounting is what a successful create looks like.
    var landed = false;
    for (var i = 0; i < 100 && !landed; i++) {
      await tester.pump(const Duration(milliseconds: 100));
      landed = find.byType(ReceiptEditPopupMenu).evaluate().isNotEmpty;
    }

    if (!landed) {
      await _dumpDiagnostics(tester, group.id, receiptName);
    }
    expect(
      landed,
      isTrue,
      reason: 'submit did not navigate to the new receipt (see DIAG lines)',
    );

    final jwt = await apiLogin();
    final created =
        (await listReceiptsForGroup(
          group.id,
          jwt: jwt,
        )).where((r) => r['name'] == receiptName).toList();
    expect(created.length, 1);
    final receipt = await getReceipt(created.single['id'] as int, jwt: jwt);
    expect((receipt['receiptItems'] as List?)?.length, 2);
  });
}

Future<void> _dumpDiagnostics(
  WidgetTester tester,
  int groupId,
  String receiptName,
) async {
  final snackbars = find.descendant(
    of: find.byType(SnackBar),
    matching: find.byType(Text),
  );
  for (final e in snackbars.evaluate()) {
    debugPrint('DIAG snackbar: ${(e.widget as Text).data}');
  }
  final forms = find.byType(ReceiptForm).evaluate();
  if (forms.isNotEmpty) {
    final model = Provider.of<ReceiptModel>(forms.first, listen: false);
    final state = model.receiptFormKey.currentState;
    debugPrint('DIAG form attached: ${state != null}');
    debugPrint('DIAG form errors: ${state?.errors}');
    debugPrint(
      'DIAG items: ${model.items.map((i) => '${i.name}=${i.amount}').toList()}',
    );
  }
  final jwt = await apiLogin();
  final matches =
      (await listReceiptsForGroup(
        groupId,
        jwt: jwt,
      )).where((r) => r['name'] == receiptName).length;
  debugPrint('DIAG receipts created on server with this name: $matches');
}
