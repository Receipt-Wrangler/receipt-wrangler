import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/constants/receipt_filter_fields.dart';
import 'package:receipt_wrangler_mobile/constants/receipts.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

import '../helpers/receipt_filter_test_helpers.dart';

/// The field table is the single definition the "Add filter" sheet, the
/// condition cards and the editor all read. These cases pin the things that
/// would otherwise drift silently -- a key that no longer names a wire field, an
/// operation with no label, a label that disagrees with the sort chip for the
/// same column.
void main() {
  test("declares the ten filterable fields, in the UI's order", () {
    expect(receiptFilterFields.map((field) => field.key).toList(), [
      "date",
      "name",
      "paidBy",
      "group",
      "amount",
      "categories",
      "tags",
      "status",
      "resolvedDate",
      "createdAt",
    ]);
  });

  test("every key is unique", () {
    final keys = receiptFilterFields.map((field) => field.key).toList();
    expect(keys.toSet().length, keys.length);
  });

  test("every key names a real ReceiptPagedRequestFilter property", () {
    // setReceiptFilterField dispatches on a string, so a typo would not fail to
    // compile -- the condition would simply never reach the API. Round-trip each
    // key through the builder and assert the property actually landed.
    for (final field in receiptFilterFields) {
      final condition = _validConditionFor(field);

      expect(serializedField({field.key: condition}, field.key), isNotNull,
          reason: "${field.key} did not reach the filter");
    }
  });

  test("every field type offers at least one operation", () {
    for (final type in ReceiptFilterFieldType.values) {
      expect(filterOperationsByType[type], isNotNull,
          reason: "$type has no operation list");
      expect(filterOperationsByType[type], isNotEmpty);
    }
  });

  test("every offered operation has a display label", () {
    for (final field in receiptFilterFields) {
      for (final operation in operationsForFilterField(field)) {
        expect(filterOperationLabels[operation], isNotNull,
            reason: "${operation.name} would render blank on ${field.key}");
        expect(filterOperationLabels[operation], isNotEmpty);
      }
    }
  });

  test("no field offers the empty operation", () {
    // FilterOperation.empty is the deserialization fallback, not something a
    // user picks -- and its wire name is "" while its Dart name is "empty".
    for (final field in receiptFilterFields) {
      expect(operationsForFilterField(field),
          isNot(contains(api.FilterOperation.empty)));
    }
  });

  group("operation lists match the desktop client", () {
    test("date excludes CONTAINS", () {
      expect(filterOperationsByType[ReceiptFilterFieldType.date], [
        api.FilterOperation.EQUALS,
        api.FilterOperation.GREATER_THAN,
        api.FilterOperation.LESS_THAN,
        api.FilterOperation.BETWEEN,
        api.FilterOperation.WITHIN_CURRENT_MONTH,
      ]);
    });

    test("number excludes CONTAINS and WITHIN_CURRENT_MONTH", () {
      expect(filterOperationsByType[ReceiptFilterFieldType.number], [
        api.FilterOperation.EQUALS,
        api.FilterOperation.GREATER_THAN,
        api.FilterOperation.LESS_THAN,
        api.FilterOperation.BETWEEN,
      ]);
    });

    test("text is CONTAINS and EQUALS", () {
      expect(filterOperationsByType[ReceiptFilterFieldType.text],
          [api.FilterOperation.CONTAINS, api.FilterOperation.EQUALS]);
    });

    test("list and users are CONTAINS only", () {
      expect(filterOperationsByType[ReceiptFilterFieldType.list],
          [api.FilterOperation.CONTAINS]);
      expect(filterOperationsByType[ReceiptFilterFieldType.users],
          [api.FilterOperation.CONTAINS]);
    });
  });

  test("a column that can also be sorted on is named the same either way", () {
    // A filter card and a sort chip describing the same column must not read
    // differently -- "Receipt Date" in one and "Date" in the other leaves the
    // user guessing which of the three date columns is meant.
    const sharedColumns = {
      "date": "date",
      "name": "name",
      "amount": "amount",
      "paidBy": "paid_by_user_id",
      "status": "status",
      "resolvedDate": "resolved_date",
      "createdAt": "created_at",
    };

    for (final field in receiptFilterFields) {
      final sortColumn = sharedColumns[field.key];
      if (sortColumn == null) continue;

      // Matching on the COLUMN as well as the label is the point: filtering on
      // displayLabel alone passes whenever any sort option happens to share the
      // name, even if this field is paired with the wrong column. It was --
      // "date" was mapped to created_at here, which this assertion now catches.
      final sortOption = receiptSortOptions.where((option) =>
          option.columnName == sortColumn &&
          option.displayLabel == field.label);

      expect(sortOption, isNotEmpty,
          reason: '"${field.label}" has no sort option on $sortColumn with '
              'the same label');
    }
  });

  test("receiptFilterFieldByKey resolves a known key and rejects the rest", () {
    expect(receiptFilterFieldByKey("paidBy")?.label, "Paid By");
    expect(receiptFilterFieldByKey("nonsense"), isNull);
  });

  group("the quick date control's fields", () {
    test("offers exactly the three date fields, in the table's order", () {
      expect(receiptDateFilterFields.map((field) => field.key).toList(),
          ["date", "resolvedDate", "createdAt"]);
    });

    test("names them as the rest of the filter UI does", () {
      // The picker, the condition card and the "Add filter" sheet all read this
      // one table, so a field cannot be called one thing in the picker and
      // another in the condition it produces.
      expect(receiptDateFilterFields.map((field) => field.label).toList(),
          ["Receipt Date", "Resolved Date", "Added At"]);
    });

    test("every one of them is a real filter field", () {
      for (final field in receiptDateFilterFields) {
        expect(receiptFilterFieldByKey(field.key), same(field));
      }
    });

    test("the default field is one it offers", () {
      expect(receiptDateFilterFields.map((field) => field.key),
          contains(defaultQuickDateFieldKey));
    });
  });
}

/// A minimal valid condition for [field], so the round-trip case above exercises
/// the encoder's skip-invalid guard rather than tripping it.
ReceiptFilterCondition _validConditionFor(ReceiptFilterField field) {
  final operation = operationsForFilterField(field).first;

  switch (field.type) {
    case ReceiptFilterFieldType.text:
      return ReceiptFilterCondition(operation: operation, value: "whole foods");
    case ReceiptFilterFieldType.number:
      return ReceiptFilterCondition(operation: operation, value: 12.5);
    case ReceiptFilterFieldType.date:
      return ReceiptFilterCondition(
          operation: operation, value: DateTime(2026, 9, 16));
    case ReceiptFilterFieldType.list:
    case ReceiptFilterFieldType.users:
      return ReceiptFilterCondition(
          operation: operation, value: [api.ReceiptStatus.OPEN]);
  }
}
