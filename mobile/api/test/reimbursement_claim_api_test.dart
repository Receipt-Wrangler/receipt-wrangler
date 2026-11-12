import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for ReimbursementClaimApi
void main() {
  final instance = Openapi().getReimbursementClaimApi();

  group(ReimbursementClaimApi, () {
    // Add receipts to a claim
    //
    // Links each listed receipt to the claim and reports each one it could not link. A receipt must belong to the claim's group, be visible to the caller, and not already be in another claim (remove it from that claim first; nothing is moved silently). A receipt already in this claim counts as added. Receipts are judged independently, so one refusal does not stop the others. Gated on group.receipts.update in the claim's group. The response carries only ids the caller sent; fetch the claim again for its new figures [SYSTEM USER]
    //
    //Future<ClaimReceiptsResult> addReceiptsToClaim(int claimId, ClaimReceiptsCommand claimReceiptsCommand) async
    test('test addReceiptsToClaim', () async {
      // TODO
    });

    // Create claim
    //
    // Creates a claim in the group, with no receipts and no ledger entries; add receipts with addReceiptsToClaim. Gated on group.receipts.update. The synthetic \"All\" group is never a destination (400) [SYSTEM USER]
    //
    //Future<Claim> createClaim(int groupId, UpsertClaimCommand upsertClaimCommand) async
    test('test createClaim', () async {
      // TODO
    });

    // Add a ledger entry to a claim
    //
    // Records a payment received, a recoupment or a note on the claim. Gated on group.receipts.update in the claim's group. Fetch the claim again for its new figures [SYSTEM USER]
    //
    //Future<ClaimLedgerEntry> createClaimLedgerEntry(int claimId, UpsertClaimLedgerEntryCommand upsertClaimLedgerEntryCommand) async
    test('test createClaimLedgerEntry', () async {
      // TODO
    });

    // Delete claim
    //
    // Deletes the claim with its ledger entries and receipt links. The receipts themselves are not deleted. Gated on group.receipts.update in the claim's group [SYSTEM USER]
    //
    //Future deleteClaim(int claimId) async
    test('test deleteClaim', () async {
      // TODO
    });

    // Delete a ledger entry
    //
    // Deletes one entry. An entry that is not on this claim, or that the caller may not see, is a 403. Gated on group.receipts.update in the claim's group [SYSTEM USER]
    //
    //Future deleteClaimLedgerEntry(int claimId, int ledgerEntryId) async
    test('test deleteClaimLedgerEntry', () async {
      // TODO
    });

    // Get claim
    //
    // Returns one claim with its computed figures, linked receipt ids and ledger entries. Gated on group.receipts.read in the claim's group. A missing claim and one the caller may not see are the same 403 [SYSTEM USER]
    //
    //Future<Claim> getClaimById(int claimId) async
    test('test getClaimById', () async {
      // TODO
    });

    // Get claims for a group
    //
    // Returns the group's claims, newest first, each with its computed figures, linked receipt ids and ledger entries, plus a summary for an \"outstanding and ageing\" header. Every figure covers only what the caller may see (see Claim). Gated on group.receipts.read [SYSTEM USER]
    //
    //Future<ClaimList> getClaimsForGroup(int groupId) async
    test('test getClaimsForGroup', () async {
      // TODO
    });

    // Remove receipts from a claim
    //
    // Unlinks each listed receipt from the claim and reports each one it could not unlink. A receipt that is not in the claim and one the caller may not see read the same. The receipts themselves are not changed. Gated on group.receipts.update in the claim's group [SYSTEM USER]
    //
    //Future<ClaimReceiptsResult> removeReceiptsFromClaim(int claimId, ClaimReceiptsCommand claimReceiptsCommand) async
    test('test removeReceiptsFromClaim', () async {
      // TODO
    });

    // Update claim
    //
    // Replaces the claim's own fields (name, reference number, type, status, dates, claimed amount, notes) with the body: an optional field that is omitted is cleared. Linked receipts and ledger entries are not part of the body and are never touched here; they have their own endpoints. A claim never moves between groups. Gated on group.receipts.update in the claim's group [SYSTEM USER]
    //
    //Future<Claim> updateClaim(int claimId, UpsertClaimCommand upsertClaimCommand) async
    test('test updateClaim', () async {
      // TODO
    });

    // Update a ledger entry
    //
    // Replaces the entry's type, date, amount and note with the body (an omitted optional field is cleared). An entry that is not on this claim, or that the caller may not see, is a 403. Gated on group.receipts.update in the claim's group [SYSTEM USER]
    //
    //Future<ClaimLedgerEntry> updateClaimLedgerEntry(int claimId, int ledgerEntryId, UpsertClaimLedgerEntryCommand upsertClaimLedgerEntryCommand) async
    test('test updateClaimLedgerEntry', () async {
      // TODO
    });

  });
}
