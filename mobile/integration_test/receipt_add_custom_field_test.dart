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

    await openManualReceiptForm(tester);
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
    // Match the label text directly: the button is a `TextButton.icon`, which
    // builds a private `_TextButtonWithIcon` subclass, so
    // `find.widgetWithText(TextButton, ...)` matches nothing on this Flutter
    // version. The button only renders once customFieldModel has loaded the
    // group's fields (async), so wait for it rather than asserting immediately.
    final addCustomFieldBtn = find.text('Add Custom Field');
    await pumpUntilFound(tester, addCustomFieldBtn);
    await tester.ensureVisible(addCustomFieldBtn);
    // ensureVisible jumps the scroll position WITHOUT a relayout, so the
    // button's global offset is stale until a frame is pumped; tapping
    // immediately computes the pre-scroll (off-screen) center, the tap misses,
    // and the modal never opens (tap-flake pattern #1 in mobile/CLAUDE.md --
    // deterministic on iOS, intermittent on Android). Pump a frame so the tap
    // hits the post-scroll center, then drain a few so the bottom sheet has
    // mounted before we look for its rows.
    await tester.pump(const Duration(milliseconds: 100));
    await tester.tap(addCustomFieldBtn);
    for (int i = 0; i < 3; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }
    // The selection modal (buildCustomFieldSelectionSheet) lists every
    // available field; wait for its header to confirm it opened, then for the
    // "E2E Notes" row. "E2E Notes" can sort last / be off-screen in the sheet,
    // so scroll it into view before tapping or the tap lands on empty space and
    // the field is never added. ensureVisible does a jumpTo on the sheet's
    // ListView WITHOUT pumping -- the item's global offset only updates after a
    // relayout, so tap() would still compute the stale off-screen center. Pump
    // frames before tapping.
    await pumpUntilFound(tester, find.text('Select Custom Field'));
    // The sheet's ListView.builder is lazy and "E2E Notes" sorts last
    // (loadCustomFields fetches every field ordered by name DESC), so its row
    // isn't built until scrolled into view -- a plain `find.text` / ensureVisible
    // can't act on an element that doesn't exist yet, which is why waiting on the
    // row alone times out. Scroll the sheet's list (the topmost Scrollable, mounted
    // above the form) until the row is built and on-screen, then tap it.
    final notesRow = find.text(_testFieldName);
    await tester.scrollUntilVisible(
      notesRow,
      120,
      scrollable: find.byType(Scrollable).last,
      maxScrolls: 80,
    );
    await tester.pump(const Duration(milliseconds: 100));
    await tester.tap(notesRow);

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
