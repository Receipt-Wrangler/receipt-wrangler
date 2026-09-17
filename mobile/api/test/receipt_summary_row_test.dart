import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ReceiptSummaryRow
void main() {
  final instance = ReceiptSummaryRowBuilder();
  // TODO add properties to the builder and call build()

  group(ReceiptSummaryRow, () {
    // ReceiptStatus status
    test('to test the property `status`', () async {
      // TODO
    });

    // How many receipts this row covers. On the overall row this matches the table's own total count for the same filter.
    // int receiptCount
    test('to test the property `receiptCount`', () async {
      // TODO
    });

    // Sum of the receipts' amounts
    // String total
    test('to test the property `total`', () async {
      // TODO
    });

    // One entry per configured currency custom field, in the configured order. Always present; empty when the group totals only the receipt amount.
    // BuiltList<ReceiptSummaryCustomFieldTotal> customFieldTotals
    test('to test the property `customFieldTotals`', () async {
      // TODO
    });

  });
}
