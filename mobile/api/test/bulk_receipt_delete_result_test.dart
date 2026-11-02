import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for BulkReceiptDeleteResult
void main() {
  final instance = BulkReceiptDeleteResultBuilder();
  // TODO add properties to the builder and call build()

  group(BulkReceiptDeleteResult, () {
    // Receipt ids that were deleted
    // BuiltList<int> deletedReceiptIds
    test('to test the property `deletedReceiptIds`', () async {
      // TODO
    });

    // One entry per requested receipt that was not deleted
    // BuiltList<BulkReceiptDeleteFailure> failures
    test('to test the property `failures`', () async {
      // TODO
    });

  });
}
