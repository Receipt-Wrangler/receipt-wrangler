import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ReceiptSummary
void main() {
  final instance = ReceiptSummaryBuilder();
  // TODO add properties to the builder and call build()

  group(ReceiptSummary, () {
    // Whether the configuration group has the summary turned on. False comes back at a normal 200 with zeroed rows, so a client with stale group settings renders nothing rather than surfacing an error.
    // bool enabled
    test('to test the property `enabled`', () async {
      // TODO
    });

    // The group whose settings produced this breakdown
    // int configurationGroupId
    test('to test the property `configurationGroupId`', () async {
      // TODO
    });

    // ReceiptSummaryRow overall
    test('to test the property `overall`', () async {
      // TODO
    });

    // One row per configured status, in ReceiptStatus declaration order. A configured status matching no receipt is still present, with zeroed figures. Always present; empty when no status is configured.
    // BuiltList<ReceiptSummaryRow> statuses
    test('to test the property `statuses`', () async {
      // TODO
    });

  });
}
