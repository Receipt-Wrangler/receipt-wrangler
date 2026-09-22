import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ReceiptSummaryCommand
void main() {
  final instance = ReceiptSummaryCommandBuilder();
  // TODO add properties to the builder and call build()

  group(ReceiptSummaryCommand, () {
    // ReceiptPagedRequestFilter filter
    test('to test the property `filter`', () async {
      // TODO
    });

    // The group whose receipt settings shape the breakdown. Exists for the synthetic \"All\" group, which spans every group the caller belongs to and so has no meaningful settings of its own - the client picks which member group's configuration to apply. OMIT it for a real group, where the group in the path is used. The DATA is always the path group's filtered set; this only chooses the shape of the breakdown. A group the caller cannot read is a 403.
    // int configurationGroupId
    test('to test the property `configurationGroupId`', () async {
      // TODO
    });

  });
}
