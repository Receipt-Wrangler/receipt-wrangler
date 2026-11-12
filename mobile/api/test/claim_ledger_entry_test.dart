import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ClaimLedgerEntry
void main() {
  final instance = ClaimLedgerEntryBuilder();
  // TODO add properties to the builder and call build()

  group(ClaimLedgerEntry, () {
    // int id
    test('to test the property `id`', () async {
      // TODO
    });

    // int claimId
    test('to test the property `claimId`', () async {
      // TODO
    });

    // One of PAYMENT (money received; adds to paidAmount), RECOUPMENT (money taken back; subtracts from paidAmount) or NOTE (no effect on the figures). A plain string, not an enum; an entry of a type the server does not recognize counts toward no figure.
    // String type
    test('to test the property `type`', () async {
      // TODO
    });

    // Calendar day the event happened (yyyy-MM-dd)
    // String date
    test('to test the property `date`', () async {
      // TODO
    });

    // Decimal string, always positive; the type gives the direction. Present on every PAYMENT and RECOUPMENT; optional, and never counted, on a NOTE.
    // String amount
    test('to test the property `amount`', () async {
      // TODO
    });

    // Free text. Absent when not set.
    // String note
    test('to test the property `note`', () async {
      // TODO
    });

    // RFC 3339 instant the entry was recorded
    // String createdAt
    test('to test the property `createdAt`', () async {
      // TODO
    });

    // RFC 3339 instant the entry was last changed
    // String updatedAt
    test('to test the property `updatedAt`', () async {
      // TODO
    });

    // The user who recorded the entry
    // int createdBy
    test('to test the property `createdBy`', () async {
      // TODO
    });

  });
}
