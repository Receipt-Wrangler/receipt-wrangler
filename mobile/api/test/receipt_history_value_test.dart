import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ReceiptHistoryValue
void main() {
  final instance = ReceiptHistoryValueBuilder();
  // TODO add properties to the builder and call build()

  group(ReceiptHistoryValue, () {
    // How to format the value. One of TEXT, CURRENCY (decimal string), DATE (calendar day, yyyy-MM-dd), INSTANT (RFC 3339), USER (see userId), STATUS (a receipt status key), BOOLEAN (\"true\" / \"false\"). Plain string.
    // String type
    test('to test the property `type`', () async {
      // TODO
    });

    // String value
    test('to test the property `value`', () async {
      // TODO
    });

    // int userId
    test('to test the property `userId`', () async {
      // TODO
    });

  });
}
