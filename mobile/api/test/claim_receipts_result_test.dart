import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ClaimReceiptsResult
void main() {
  final instance = ClaimReceiptsResultBuilder();
  // TODO add properties to the builder and call build()

  group(ClaimReceiptsResult, () {
    // Requested receipt ids that were added (or removed), in request order
    // BuiltList<int> receiptIds
    test('to test the property `receiptIds`', () async {
      // TODO
    });

    // One entry per requested receipt that was not added (or removed)
    // BuiltList<ClaimReceiptsFailure> failures
    test('to test the property `failures`', () async {
      // TODO
    });

  });
}
