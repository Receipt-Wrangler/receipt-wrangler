import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for Claim
void main() {
  final instance = ClaimBuilder();
  // TODO add properties to the builder and call build()

  group(Claim, () {
    // int id
    test('to test the property `id`', () async {
      // TODO
    });

    // The group the claim and all its receipts belong to
    // int groupId
    test('to test the property `groupId`', () async {
      // TODO
    });

    // String name
    test('to test the property `name`', () async {
      // TODO
    });

    // A travel authorization or voucher number. Absent when not set.
    // String referenceNumber
    test('to test the property `referenceNumber`', () async {
      // TODO
    });

    // Free text, e.g. \"Travel voucher\" or \"TQSA\". Absent when not set.
    // String type
    test('to test the property `type`', () async {
      // TODO
    });

    // One of DRAFT, SUBMITTED, PARTIALLY_PAID, PAID, DISPUTED, CLOSED. Set by hand; any status may follow any other. A plain string, not an enum, so a client must tolerate values it does not know.
    // String status
    test('to test the property `status`', () async {
      // TODO
    });

    // Whether the status is one the server treats as waiting on the other party (SUBMITTED, PARTIALLY_PAID, DISPUTED). Drives daysWaiting and the list summary, so a client need not hard-code the set.
    // bool awaitingSettlement
    test('to test the property `awaitingSettlement`', () async {
      // TODO
    });

    // Calendar day the claim was submitted (yyyy-MM-dd). Absent when not set.
    // String submittedDate
    test('to test the property `submittedDate`', () async {
      // TODO
    });

    // Calendar day a response or payment is due (yyyy-MM-dd). Absent when not set.
    // String dueDate
    test('to test the property `dueDate`', () async {
      // TODO
    });

    // The amount explicitly entered as claimed, as a decimal string. Absent when not set, in which case effectiveClaimedAmount falls back to receiptsTotal.
    // String claimedAmount
    test('to test the property `claimedAmount`', () async {
      // TODO
    });

    // Free-text notes. Absent when not set.
    // String notes
    test('to test the property `notes`', () async {
      // TODO
    });

    // RFC 3339 instant
    // String createdAt
    test('to test the property `createdAt`', () async {
      // TODO
    });

    // RFC 3339 instant
    // String updatedAt
    test('to test the property `updatedAt`', () async {
      // TODO
    });

    // The user who created the claim
    // int createdBy
    test('to test the property `createdBy`', () async {
      // TODO
    });

    // Ids of the linked receipts the caller may see, ascending. Fetch the receipts themselves through the receipt endpoints. Always present; empty when none.
    // BuiltList<int> receiptIds
    test('to test the property `receiptIds`', () async {
      // TODO
    });

    // Number of entries in receiptIds
    // int receiptCount
    test('to test the property `receiptCount`', () async {
      // TODO
    });

    // Sum of the amounts of the receipts in receiptIds
    // String receiptsTotal
    test('to test the property `receiptsTotal`', () async {
      // TODO
    });

    // claimedAmount when set, otherwise receiptsTotal
    // String effectiveClaimedAmount
    test('to test the property `effectiveClaimedAmount`', () async {
      // TODO
    });

    // Sum of the amounts of the PAYMENT ledger entries
    // String paymentsTotal
    test('to test the property `paymentsTotal`', () async {
      // TODO
    });

    // Sum of the amounts of the RECOUPMENT ledger entries
    // String recoupmentsTotal
    test('to test the property `recoupmentsTotal`', () async {
      // TODO
    });

    // paymentsTotal minus recoupmentsTotal; negative when more was taken back than paid
    // String paidAmount
    test('to test the property `paidAmount`', () async {
      // TODO
    });

    // effectiveClaimedAmount minus paidAmount. Negative when more was received than claimed.
    // String outstandingAmount
    test('to test the property `outstandingAmount`', () async {
      // TODO
    });

    // Calendar days from submittedDate to today in the app time zone, when awaitingSettlement is true and submittedDate is set; otherwise absent. Never negative.
    // int daysWaiting
    test('to test the property `daysWaiting`', () async {
      // TODO
    });

    // The claim's ledger, oldest first (by date, then id). Always present; empty when none.
    // BuiltList<ClaimLedgerEntry> ledgerEntries
    test('to test the property `ledgerEntries`', () async {
      // TODO
    });

  });
}
