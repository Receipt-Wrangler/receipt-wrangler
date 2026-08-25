import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/shared/functions/custom_field_values.dart';
import 'package:receipt_wrangler_mobile/utils/date.dart';

import '../../helpers/receipt_form_test_helpers.dart';

// The mobile custom field payload builder. Its contract is that an *empty*
// custom field value is representable: the backend replaces the whole
// association on update and rejects a submitted id set that doesn't match the
// stored one, so every attached value has to round-trip -- even one nobody
// filled in, and even one whose template the caller isn't allowed to read.

const _receiptId = 5;

const _textId = 1;
const _dateId = 2;
const _selectId = 3;
const _currencyId = 4;
const _booleanId = 5;

final _catalog = [
  buildCustomField(id: _textId, name: 'PO Number'),
  buildCustomField(id: _dateId, name: 'Invoiced', type: api.CustomFieldType.DATE),
  buildCustomField(
    id: _selectId,
    name: 'Cost Centre',
    type: api.CustomFieldType.SELECT,
    options: [
      buildCustomFieldOption(id: 91, customFieldId: _selectId, value: 'Ops'),
    ],
  ),
  buildCustomField(
    id: _currencyId,
    name: 'Tax',
    type: api.CustomFieldType.CURRENCY,
  ),
  buildCustomField(
    id: _booleanId,
    name: 'Reimbursed',
    type: api.CustomFieldType.BOOLEAN,
  ),
];

List<api.UpsertCustomFieldValueCommand> _commands({
  required List<api.CustomFieldValue> attachedValues,
  required List<api.CustomField> customFields,
  Map<String, dynamic> form = const {},
}) =>
    buildCustomFieldValueUpsertCommands(
      attachedValues: attachedValues,
      customFields: customFields,
      form: form,
      receiptId: _receiptId,
    );

void main() {
  group('buildCustomFieldValueUpsertCommands', () {
    test('maps each type onto its own value column', () {
      final invoicedOn = DateTime(2026, 8, 25);
      final commands = _commands(
        attachedValues: [
          buildCustomFieldValue(customFieldId: _textId),
          buildCustomFieldValue(customFieldId: _dateId),
          buildCustomFieldValue(customFieldId: _selectId),
          buildCustomFieldValue(customFieldId: _currencyId),
          buildCustomFieldValue(customFieldId: _booleanId),
        ],
        customFields: _catalog,
        form: {
          'customField_$_textId': 'PO-1234',
          'customField_$_dateId': invoicedOn,
          'customField_$_selectId': 91,
          'customField_$_currencyId': '12.34',
          'customField_$_booleanId': true,
        },
      );

      expect(commands.map((c) => c.customFieldId),
          [_textId, _dateId, _selectId, _currencyId, _booleanId]);
      expect(commands.every((c) => c.receiptId == _receiptId), isTrue);
      expect(commands[0].stringValue, 'PO-1234');
      expect(commands[1].dateValue, formatDate(zuluDateFormat, invoicedOn));
      expect(commands[2].selectValue, 91);
      expect(commands[3].currencyValue, '12.34');
      expect(commands[4].booleanValue, isTrue);
    });

    test('emits an empty entry rather than dropping a blank value', () {
      // A blank field still means "this custom field belongs on this receipt",
      // which is what makes a group's default custom fields work.
      final commands = _commands(
        attachedValues: [
          buildCustomFieldValue(customFieldId: _textId),
          buildCustomFieldValue(customFieldId: _selectId),
        ],
        customFields: _catalog,
        form: const {
          'customField_$_textId': '',
          'customField_$_selectId': null,
        },
      );

      expect(commands.map((c) => c.customFieldId), [_textId, _selectId]);
      expect(commands[0].stringValue, isNull);
      expect(commands[1].selectValue, isNull);
    });

    test('emits an unchecked boolean as false', () {
      final commands = _commands(
        attachedValues: [buildCustomFieldValue(customFieldId: _booleanId)],
        customFields: _catalog,
        form: const {'customField_$_booleanId': false},
      );

      expect(commands.single.booleanValue, isFalse);
    });

    test('clears a value the user emptied out', () {
      final commands = _commands(
        attachedValues: [
          buildCustomFieldValue(customFieldId: _textId, stringValue: 'PO-1234'),
        ],
        customFields: _catalog,
        form: const {'customField_$_textId': ''},
      );

      expect(commands.single.stringValue, isNull);
    });

    test('falls back to the stored value when the template is missing', () {
      // The catalog is empty for a caller without `app.custom-fields.read`
      // (CustomFieldModel swallows the 403), and no input is ever mounted for
      // those values. Dropping them would both lose data and make the submitted
      // id set mismatch the stored one -- a 403 from
      // enforceReceiptCustomFieldSelection.
      final commands = _commands(
        attachedValues: [
          buildCustomFieldValue(customFieldId: _textId, stringValue: 'PO-1234'),
          buildCustomFieldValue(customFieldId: _booleanId, booleanValue: true),
          buildCustomFieldValue(customFieldId: _currencyId),
        ],
        customFields: const [],
      );

      expect(commands.map((c) => c.customFieldId),
          [_textId, _booleanId, _currencyId],
          reason: 'the whole attached set survives');
      expect(commands[0].stringValue, 'PO-1234');
      expect(commands[1].booleanValue, isTrue);
      expect(commands[2].currencyValue, isNull);
    });

    test('falls back per value, not per payload', () {
      // A partially loaded catalog must not make the loaded fields fall back to
      // stale stored values.
      final commands = _commands(
        attachedValues: [
          buildCustomFieldValue(customFieldId: _textId, stringValue: 'stale'),
          buildCustomFieldValue(customFieldId: _dateId, dateValue: 'kept'),
        ],
        customFields: [_catalog.first],
        form: const {'customField_$_textId': 'fresh'},
      );

      expect(commands[0].stringValue, 'fresh');
      expect(commands[1].dateValue, 'kept');
    });

    test('is empty when the receipt carries no custom fields', () {
      expect(_commands(attachedValues: const [], customFields: _catalog),
          isEmpty);
    });
  });
}
