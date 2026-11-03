import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ReceiptHistoryActor
void main() {
  final instance = ReceiptHistoryActorBuilder();
  // TODO add properties to the builder and call build()

  group(ReceiptHistoryActor, () {
    // One of USER, QUICK_SCAN, EMAIL, SYSTEM. Plain string.
    // String type
    test('to test the property `type`', () async {
      // TODO
    });

    // The user, for USER; the scanning user (when visible), for QUICK_SCAN
    // int userId
    test('to test the property `userId`', () async {
      // TODO
    });

  });
}
