import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for BulkReceiptDeleteFailure
void main() {
  final instance = BulkReceiptDeleteFailureBuilder();
  // TODO add properties to the builder and call build()

  group(BulkReceiptDeleteFailure, () {
    // Requested receipt id that was not deleted
    // int receiptId
    test('to test the property `receiptId`', () async {
      // TODO
    });

    // Why it was not deleted. A missing receipt and one the caller may not delete read the same, so this never reveals that a receipt exists
    // String message
    test('to test the property `message`', () async {
      // TODO
    });

  });
}
