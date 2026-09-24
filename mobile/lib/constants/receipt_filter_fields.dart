import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;

/// The operation-option bucket a filter field belongs to. Mirrors the desktop
/// client's `ReceiptFilterFieldType` so a field offers the same operations on
/// both clients.
enum ReceiptFilterFieldType { date, text, number, list, users }

/// One filterable receipt field: how it is named on the wire, how it reads, and
/// what kind of value it holds.
class ReceiptFilterField {
  const ReceiptFilterField({
    required this.key,
    required this.label,
    required this.icon,
    required this.hint,
    required this.type,
  });

  /// The `ReceiptPagedRequestFilter` property this field writes to.
  final String key;

  final String label;

  final IconData icon;

  /// The one-line description shown beside the label in the "Add filter" sheet.
  final String hint;

  final ReceiptFilterFieldType type;
}

/// The ten filterable receipt fields, in the order the filter UI offers them.
///
/// This is the single definition of a field's label, icon and type -- the "Add
/// filter" sheet, the condition cards and the condition editor all read it, so a
/// card can never disagree with the row that produced it. It mirrors desktop's
/// `RECEIPT_FILTER_FIELDS` (`desktop/src/constants/receipt-filter-fields.constant.ts`).
///
/// The labels match [receiptSortOptions] (`constants/receipts.dart`) wherever the
/// two overlap: `date` is "Receipt Date" rather than a bare "Date", because the
/// list can also be filtered on `resolvedDate` and `createdAt` and "Date" leaves
/// the user guessing which of the three is meant.
const List<ReceiptFilterField> receiptFilterFields = [
  ReceiptFilterField(
    key: "date",
    label: "Receipt Date",
    icon: Icons.calendar_month,
    hint: "Receipt date",
    type: ReceiptFilterFieldType.date,
  ),
  ReceiptFilterField(
    key: "name",
    label: "Name",
    icon: Icons.title,
    hint: "Merchant or receipt name",
    type: ReceiptFilterFieldType.text,
  ),
  ReceiptFilterField(
    key: "paidBy",
    label: "Paid By",
    icon: Icons.person,
    hint: "One or more people",
    type: ReceiptFilterFieldType.users,
  ),
  ReceiptFilterField(
    key: "group",
    label: "Group",
    icon: Icons.group,
    hint: "All groups view only",
    type: ReceiptFilterFieldType.list,
  ),
  ReceiptFilterField(
    key: "amount",
    label: "Amount",
    icon: Icons.attach_money,
    hint: "Currency value",
    type: ReceiptFilterFieldType.number,
  ),
  ReceiptFilterField(
    key: "categories",
    label: "Categories",
    icon: Icons.sell,
    hint: "Any of the selected",
    type: ReceiptFilterFieldType.list,
  ),
  ReceiptFilterField(
    key: "tags",
    label: "Tags",
    icon: Icons.label,
    hint: "Any of the selected",
    type: ReceiptFilterFieldType.list,
  ),
  ReceiptFilterField(
    key: "status",
    label: "Status",
    icon: Icons.check_circle,
    hint: "Open, needs attention, resolved",
    type: ReceiptFilterFieldType.list,
  ),
  ReceiptFilterField(
    key: "resolvedDate",
    label: "Resolved Date",
    icon: Icons.event_available,
    hint: "When it got resolved",
    type: ReceiptFilterFieldType.date,
  ),
  ReceiptFilterField(
    key: "createdAt",
    label: "Added At",
    icon: Icons.schedule,
    hint: "When it was uploaded",
    type: ReceiptFilterFieldType.date,
  ),
];

/// The date fields the quick date control can be pointed at, in the order
/// [receiptFilterFields] declares them: Receipt Date, Resolved Date, Added At.
///
/// Derived rather than restated, exactly as desktop's
/// `RECEIPT_DATE_FILTER_FIELDS` is, so the stepper's field picker, the "Add
/// filter" sheet and the condition card can never name a field differently.
/// Desktop narrows on the *key* because TypeScript needs the runtime check to
/// justify the narrower type it claims; Dart needs no such narrowing, so
/// filtering on the type is both honest and self-maintaining -- a new
/// `ReceiptFilterFieldType.date` field reaches the picker with no second edit.
final List<ReceiptFilterField> receiptDateFilterFields = receiptFilterFields
    .where((field) => field.type == ReceiptFilterFieldType.date)
    .toList(growable: false);

/// The date field the quick date control starts on.
const String defaultQuickDateFieldKey = "date";

/// The field named [key], or null when nothing declares it.
ReceiptFilterField? receiptFilterFieldByKey(String key) {
  for (final field in receiptFilterFields) {
    if (field.key == key) {
      return field;
    }
  }
  return null;
}

/// Which operations each field type offers, mirroring desktop's
/// `filter-operations-options.constant.ts`.
///
/// `CONTAINS` is a substring match for text and an `IN` for a list, which is why
/// it is the only operation the list/users buckets offer and is absent from the
/// date and number ones. `WITHIN_CURRENT_MONTH` is date-only -- it carries no
/// value and the API pins it to month-start through today.
const Map<ReceiptFilterFieldType, List<api.FilterOperation>>
    filterOperationsByType = {
  ReceiptFilterFieldType.date: [
    api.FilterOperation.EQUALS,
    api.FilterOperation.GREATER_THAN,
    api.FilterOperation.LESS_THAN,
    api.FilterOperation.BETWEEN,
    api.FilterOperation.WITHIN_CURRENT_MONTH,
  ],
  ReceiptFilterFieldType.text: [
    api.FilterOperation.CONTAINS,
    api.FilterOperation.EQUALS,
  ],
  ReceiptFilterFieldType.number: [
    api.FilterOperation.EQUALS,
    api.FilterOperation.GREATER_THAN,
    api.FilterOperation.LESS_THAN,
    api.FilterOperation.BETWEEN,
  ],
  ReceiptFilterFieldType.list: [api.FilterOperation.CONTAINS],
  ReceiptFilterFieldType.users: [api.FilterOperation.CONTAINS],
};

/// Human-readable label per operation, shared by the operation chips and the
/// condition cards so both read a condition the same way. Mirrors desktop's
/// `FILTER_OPERATION_DISPLAY_VALUES`.
const Map<api.FilterOperation, String> filterOperationLabels = {
  api.FilterOperation.CONTAINS: "Contains",
  api.FilterOperation.EQUALS: "Equals",
  api.FilterOperation.GREATER_THAN: "Greater than",
  api.FilterOperation.LESS_THAN: "Less than",
  api.FilterOperation.BETWEEN: "Between",
  api.FilterOperation.WITHIN_CURRENT_MONTH: "Within current month",
};

/// The operations [field] offers, empty for a type with no entry.
List<api.FilterOperation> operationsForFilterField(ReceiptFilterField field) {
  return filterOperationsByType[field.type] ?? const [];
}
