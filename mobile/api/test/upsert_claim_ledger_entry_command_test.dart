import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for UpsertClaimLedgerEntryCommand
void main() {
  final instance = UpsertClaimLedgerEntryCommandBuilder();
  // TODO add properties to the builder and call build()

  group(UpsertClaimLedgerEntryCommand, () {
    // One of PAYMENT, RECOUPMENT, NOTE; anything else is a 400 keyed type
    // String type
    test('to test the property `type`', () async {
      // TODO
    });

    // Calendar day (yyyy-MM-dd)
    // String date
    test('to test the property `date`', () async {
      // TODO
    });

    // Decimal string, greater than zero, at most two decimal places. Required for PAYMENT and RECOUPMENT; optional for NOTE.
    // String amount
    test('to test the property `amount`', () async {
      // TODO
    });

    // At most 2000 characters. Required (non-blank) for NOTE.
    // String note
    test('to test the property `note`', () async {
      // TODO
    });

  });
}
