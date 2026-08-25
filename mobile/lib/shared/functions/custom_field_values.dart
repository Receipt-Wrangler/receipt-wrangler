import 'package:openapi/openapi.dart' as api;

import '../../utils/date.dart';

/// The FormBuilder field name a custom field's input registers itself under.
/// Shared so the widget that mounts the field and the code that reads its
/// value back can never disagree about the key.
String customFieldFormFieldName(int customFieldId) =>
    "customField_$customFieldId";

/// The type-specific value columns of a single custom field value.
///
/// A custom field value is one row with five nullable, mutually exclusive
/// value columns; which one is populated depends on the custom field's type.
/// All-null is a valid, meaningful state -- an *empty* value (see
/// [buildCustomFieldValueUpsertCommands]).
class CustomFieldValueColumns {
  const CustomFieldValueColumns({
    this.stringValue,
    this.dateValue,
    this.selectValue,
    this.currencyValue,
    this.booleanValue,
  });

  final String? stringValue;

  final String? dateValue;

  final int? selectValue;

  final String? currencyValue;

  final bool? booleanValue;
}

String? _nonEmptyString(dynamic fieldValue) {
  if (fieldValue == null) {
    return null;
  }

  var asString = fieldValue.toString();
  return asString.isEmpty ? null : asString;
}

/// Reads the value columns for [customField] out of a saved FormBuilder value
/// map. A missing or blank form value yields all-null columns rather than no
/// columns at all -- an empty custom field value is representable on purpose.
CustomFieldValueColumns customFieldValueColumnsFromForm(
  api.CustomField customField,
  Map<String, dynamic> form,
) {
  var fieldValue = form[customFieldFormFieldName(customField.id)];

  switch (customField.type) {
    case api.CustomFieldType.TEXT:
      return CustomFieldValueColumns(stringValue: _nonEmptyString(fieldValue));
    case api.CustomFieldType.DATE:
      return CustomFieldValueColumns(
        dateValue: fieldValue is DateTime
            ? formatDate(zuluDateFormat, fieldValue)
            : _nonEmptyString(fieldValue),
      );
    case api.CustomFieldType.SELECT:
      return CustomFieldValueColumns(
          selectValue: fieldValue is int ? fieldValue : null);
    case api.CustomFieldType.CURRENCY:
      return CustomFieldValueColumns(
          currencyValue: _nonEmptyString(fieldValue));
    case api.CustomFieldType.BOOLEAN:
      return CustomFieldValueColumns(
          booleanValue: fieldValue is bool ? fieldValue : null);
  }

  // CustomFieldType is a built_value EnumClass, so the switch above is not
  // statically exhaustive: an enum value added by a future client regen lands
  // here rather than failing to compile.
  return const CustomFieldValueColumns();
}

/// The columns already stored on [value], used as the fallback when the form
/// never rendered an input for it.
CustomFieldValueColumns customFieldValueColumnsFromValue(
        api.CustomFieldValue value) =>
    CustomFieldValueColumns(
      stringValue: value.stringValue,
      dateValue: value.dateValue,
      selectValue: value.selectValue,
      currencyValue: value.currencyValue,
      booleanValue: value.booleanValue,
    );

api.CustomField? _findCustomField(
    List<api.CustomField> customFields, int customFieldId) {
  for (var customField in customFields) {
    if (customField.id == customFieldId) {
      return customField;
    }
  }

  return null;
}

/// Resolves the columns to submit for one custom field value attached to the
/// receipt: the live form value when its template is known, otherwise the
/// value already stored on the receipt.
CustomFieldValueColumns _resolveColumns(
  api.CustomFieldValue attachedValue,
  List<api.CustomField> customFields,
  Map<String, dynamic> form,
) {
  var customField =
      _findCustomField(customFields, attachedValue.customFieldId);

  // The template is missing when the caller lacks `app.custom-fields.read`
  // (CustomFieldModel swallows that 403 into an empty catalog) or the catalog
  // simply hasn't loaded yet. Either way the form never mounted an input for
  // this value, so there is nothing to read -- carry the stored value through
  // untouched instead of dropping the entry.
  if (customField == null) {
    return customFieldValueColumnsFromValue(attachedValue);
  }

  return customFieldValueColumnsFromForm(customField, form);
}

/// Builds the `customFields` payload of an `UpsertReceiptCommand` from the
/// receipt's currently attached custom field values plus the live form.
///
/// One entry is emitted for **every** attached value, including empty ones --
/// their type-specific columns are simply left null. Empty values are
/// meaningful, not noise: a group can declare default custom fields, and a
/// field left blank still means "this field belongs on this receipt".
///
/// Dropping empty entries also breaks saving outright, because the backend
/// REPLACES the whole association on update (`UpdateReceipt`,
/// `api/internal/repositories/receipts.go`) and rejects a submitted id set
/// that does not match the stored one (`enforceReceiptCustomFieldSelection`,
/// HTTP 403). The desktop client sends the same full array with null columns;
/// the two clients must agree about what a receipt looks like.
List<api.UpsertCustomFieldValueCommand> buildCustomFieldValueUpsertCommands({
  required Iterable<api.CustomFieldValue> attachedValues,
  required List<api.CustomField> customFields,
  required Map<String, dynamic> form,
  required int receiptId,
}) {
  return attachedValues.map((attachedValue) {
    var columns = _resolveColumns(attachedValue, customFields, form);

    return (api.UpsertCustomFieldValueCommandBuilder()
          ..receiptId = receiptId
          ..customFieldId = attachedValue.customFieldId
          ..stringValue = columns.stringValue
          ..dateValue = columns.dateValue
          ..selectValue = columns.selectValue
          ..currencyValue = columns.currencyValue
          ..booleanValue = columns.booleanValue)
        .build();
  }).toList();
}
