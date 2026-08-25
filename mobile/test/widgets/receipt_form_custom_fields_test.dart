import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/shared/widgets/custom_field_widget.dart';

import '../helpers/receipt_form_test_helpers.dart';

// Widget coverage for the receipt form's manual custom field management: the
// "Add Custom Field" sheet, the per-field remove button, and the guarantee that
// a removed field comes back blank rather than carrying its old value
// (FormBuilder keeps the value of an unregistered field unless it is cleared
// first -- `clearValueOnUnregister` defaults to false).

const _groupId = 1;
const _textFieldId = 11;
const _checkboxFieldId = 12;

final _textCustomField = buildCustomField(id: _textFieldId, name: 'PO Number');

final _booleanCustomField = buildCustomField(
  id: _checkboxFieldId,
  name: 'Reimbursed',
  type: api.CustomFieldType.BOOLEAN,
);

Finder _textInput(int customFieldId) => find.byWidgetPredicate(
      (w) => w is FormBuilderTextField && w.name == 'customField_$customFieldId',
    );

Finder _addCustomFieldButton() => find.text('Add Custom Field');

/// The sheet's row for [name]. Scoped to the ListTile so it can never match a
/// mounted field's own label, which renders the same text.
Finder _customFieldSheetRow(String name) => find.descendant(
      of: find.byType(ListTile),
      matching: find.text(name),
    );

/// Opens the "Add Custom Field" sheet and picks [name]. The button sits below
/// the fold once a few fields are mounted, so scroll it into view first.
Future<void> _addCustomFieldNamed(WidgetTester tester, String name) async {
  await tester.ensureVisible(_addCustomFieldButton());
  await tester.pumpAndSettle();
  await tester.tap(_addCustomFieldButton());
  await tester.pumpAndSettle();

  expect(find.text('Select Custom Field'), findsOneWidget);

  await tester.tap(_customFieldSheetRow(name));
  await tester.pumpAndSettle();
}

void main() {
  testWidgets('adds a custom field through the Add Custom Field sheet',
      (tester) async {
    final harness = await pumpReceiptForm(
      tester,
      groups: [buildGroup(id: _groupId)],
      customFields: [_textCustomField],
    );

    expect(find.byType(CustomFieldWidget), findsNothing);

    await _addCustomFieldNamed(tester, 'PO Number');

    expect(find.byType(CustomFieldWidget), findsOneWidget);
    expect(_textInput(_textFieldId), findsOneWidget);
    expect(
      harness.receiptModel.modifiedReceipt.customFields
          .map((cfv) => cfv.customFieldId),
      [_textFieldId],
      reason: 'the value is attached to the receipt being edited',
    );
    expect(tester.takeException(), isNull);
  });

  testWidgets('offers only custom fields that are not already attached',
      (tester) async {
    await pumpReceiptForm(
      tester,
      groups: [buildGroup(id: _groupId)],
      customFields: [_textCustomField, _booleanCustomField],
    );

    await _addCustomFieldNamed(tester, 'PO Number');
    await tester.ensureVisible(_addCustomFieldButton());
    await tester.pumpAndSettle();
    await tester.tap(_addCustomFieldButton());
    await tester.pumpAndSettle();

    expect(_customFieldSheetRow('Reimbursed'), findsOneWidget);
    expect(
      _customFieldSheetRow('PO Number'),
      findsNothing,
      reason: 'an already-attached field is not offered again',
    );
    expect(tester.takeException(), isNull);
  });

  testWidgets('hides the Add Custom Field button once the catalog is exhausted',
      (tester) async {
    await pumpReceiptForm(
      tester,
      groups: [buildGroup(id: _groupId)],
      customFields: [_textCustomField],
    );

    await _addCustomFieldNamed(tester, 'PO Number');

    expect(_addCustomFieldButton(), findsNothing);
    expect(tester.takeException(), isNull);
  });

  testWidgets('removes a custom field through its remove button',
      (tester) async {
    final harness = await pumpReceiptForm(
      tester,
      groups: [buildGroup(id: _groupId)],
      customFields: [_textCustomField],
    );

    await _addCustomFieldNamed(tester, 'PO Number');
    expect(find.byType(CustomFieldWidget), findsOneWidget);

    await tester.ensureVisible(find.byIcon(Icons.remove_circle_outline));
    await tester.pumpAndSettle();
    await tester.tap(find.byIcon(Icons.remove_circle_outline));
    await tester.pumpAndSettle();

    expect(find.byType(CustomFieldWidget), findsNothing);
    expect(_textInput(_textFieldId), findsNothing);
    expect(harness.receiptModel.modifiedReceipt.customFields, isEmpty);
    expect(tester.takeException(), isNull);
  });

  testWidgets('re-adding a removed custom field gives a blank field',
      (tester) async {
    // FormBuilder keeps an unregistered field's value (clearValueOnUnregister
    // defaults to false) and hands it straight back to the next field
    // registering under the same name. Without the explicit clear in
    // `_removeCustomField`, the re-added field would resurrect "PO-1234".
    final harness = await pumpReceiptForm(
      tester,
      groups: [buildGroup(id: _groupId)],
      customFields: [_textCustomField],
    );

    await _addCustomFieldNamed(tester, 'PO Number');
    await tester.enterText(_textInput(_textFieldId), 'PO-1234');
    await tester.pumpAndSettle();
    expect(harness.customFieldValue(_textFieldId), 'PO-1234');

    await tester.ensureVisible(find.byIcon(Icons.remove_circle_outline));
    await tester.pumpAndSettle();
    await tester.tap(find.byIcon(Icons.remove_circle_outline));
    await tester.pumpAndSettle();

    await _addCustomFieldNamed(tester, 'PO Number');

    expect(_textInput(_textFieldId), findsOneWidget);
    expect(
      harness.customFieldValue(_textFieldId),
      isNull,
      reason: 'the stale form value did not come back',
    );
    expect(find.text('PO-1234'), findsNothing);
    expect(tester.takeException(), isNull);
  });
}
