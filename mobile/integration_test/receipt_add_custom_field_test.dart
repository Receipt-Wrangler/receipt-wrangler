// Flow #9 -- add a receipt with a custom field value.
//
// The custom-fields feature lets group admins define optional per-
// receipt fields (TEXT, DATE, SELECT, CURRENCY, BOOLEAN). The receipt
// form has an "Add Custom Field" button that opens a modal listing
// the available fields; selecting one adds it to the form, and the
// user fills the value before saving.
//
// The "E2E Notes" TEXT field is auto-provisioned via `ensureCustomField`
// at test setup -- no manual seeding required. Once the field exists
// on the backend, subsequent runs reuse it.
//
// Own file for the GoRouter-persistence reason (see other test files).

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

const _testFieldName = 'E2E Notes';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('admin can add a receipt with a custom field value',
      (tester) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    // Auto-provision the "E2E Notes" TEXT field. Idempotent -- reuses
    // an existing one of the same name, otherwise creates it.
    final jwt = await apiLogin();
    final notesField = await ensureCustomField(
      jwt: jwt,
      name: _testFieldName,
      type: 'TEXT',
    );
    final fieldId = notesField['id'] as int;

    await loginAsAdmin(tester);

    final receiptName =
        'e2e-cf-${DateTime.now().millisecondsSinceEpoch}';
    final fieldValue =
        'note-${DateTime.now().millisecondsSinceEpoch}';

    await tester.tap(find.text('Add'));
    await pumpUntilFound(tester, find.text('Add Manual Receipt'));
    await tester.tap(find.text('Add Manual Receipt'));
    await pumpUntilFound(tester, find.text('Name'));

    // Fill required fields first -- the custom field add UI is below
    // the standard fields in the scroll view.
    await tester.enterText(formField('name'), receiptName);
    await tester.enterText(formField('amount'), '12.34');
    await selectDropdown(tester, 'groupId', 'My Receipts');
    await selectDropdown(tester, 'paidByUserId', adminDisplayName(tester));

    // Drain dropdown overlay teardown before tapping the Add Custom
    // Field button (the modal-bottom-sheet open / close interacts
    // with the same overlay routes the dropdowns used).
    await tester.pumpAndSettle(const Duration(seconds: 3));

    // Open the modal of available fields. The button is a TextButton.icon
    // (NOT ElevatedButton -- buildAddCustomFieldButton at
    // mobile/lib/receipts/widgets/receipt_form.dart:274). It's only
    // rendered when `customFieldModel.customFields` minus already-added
    // ones is non-empty -- our setUp's `ensureCustomField` makes sure of
    // that. The TextButton.icon's label is the Text 'Add Custom Field'.
    final addCustomFieldBtn =
        find.widgetWithText(TextButton, 'Add Custom Field');
    await tester.ensureVisible(addCustomFieldBtn);
    await tester.tap(addCustomFieldBtn);
    await pumpUntilFound(tester, find.text(_testFieldName));
    await tester.tap(find.text(_testFieldName));

    // The custom field widget mounts with name `customField_<id>`
    // (mobile/lib/shared/widgets/custom_field_widget.dart line ~47 for
    // TEXT type). Wait for it to render and fill it.
    await pumpUntilFound(tester, formField('customField_$fieldId'));
    await tester.enterText(
      formField('customField_$fieldId'),
      fieldValue,
    );

    // Save.
    await tester.pumpAndSettle(const Duration(seconds: 3));
    await tester.tap(find.byType(BottomSubmitButton));
    // /view shell mounted -> ReceiptEditPopupMenu is in the tree.
    await pumpUntilFound(tester, find.byType(ReceiptEditPopupMenu));
    final receiptId = receiptIdFromUrl(currentUrl(tester));
    scheduleReceiptCleanup(receiptId);

    // Verify via API: the receipt has a customFieldValue with our id
    // and value.
    final receipt = await getReceipt(receiptId, jwt: jwt);
    final customFieldValues = (receipt['customFields'] as List?) ?? const [];
    final match = customFieldValues.cast<Map<String, dynamic>>().firstWhere(
          (cf) =>
              cf['customFieldId'] == fieldId &&
              cf['stringValue'] == fieldValue,
          orElse: () => <String, dynamic>{},
        );
    expect(
      match.isNotEmpty,
      isTrue,
      reason:
          'Receipt should have a custom-field value for "$_testFieldName" '
          '(id=$fieldId) equal to "$fieldValue". Got: $customFieldValues',
    );
  });
}
