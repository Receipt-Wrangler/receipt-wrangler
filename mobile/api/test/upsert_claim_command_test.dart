import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for UpsertClaimCommand
void main() {
  final instance = UpsertClaimCommandBuilder();
  // TODO add properties to the builder and call build()

  group(UpsertClaimCommand, () {
    // Required, non-blank after trimming, at most 255 characters
    // String name
    test('to test the property `name`', () async {
      // TODO
    });

    // At most 100 characters
    // String referenceNumber
    test('to test the property `referenceNumber`', () async {
      // TODO
    });

    // Free text, at most 100 characters
    // String type
    test('to test the property `type`', () async {
      // TODO
    });

    // One of DRAFT, SUBMITTED, PARTIALLY_PAID, PAID, DISPUTED, CLOSED; anything else is a 400 keyed status
    // String status
    test('to test the property `status`', () async {
      // TODO
    });

    // Calendar day (yyyy-MM-dd)
    // String submittedDate
    test('to test the property `submittedDate`', () async {
      // TODO
    });

    // Calendar day (yyyy-MM-dd)
    // String dueDate
    test('to test the property `dueDate`', () async {
      // TODO
    });

    // Decimal string, zero or more, at most two decimal places
    // String claimedAmount
    test('to test the property `claimedAmount`', () async {
      // TODO
    });

    // At most 10000 characters
    // String notes
    test('to test the property `notes`', () async {
      // TODO
    });

  });
}
