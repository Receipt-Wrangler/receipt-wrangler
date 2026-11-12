import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ClaimReceiptsFailure
void main() {
  final instance = ClaimReceiptsFailureBuilder();
  // TODO add properties to the builder and call build()

  group(ClaimReceiptsFailure, () {
    // Requested receipt id that was not added (or removed)
    // int receiptId
    test('to test the property `receiptId`', () async {
      // TODO
    });

    // Why not. A receipt that does not exist, belongs to another group, or is hidden from the caller reads the same, so this never reveals that a receipt exists.
    // String message
    test('to test the property `message`', () async {
      // TODO
    });

  });
}
