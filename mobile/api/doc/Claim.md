# openapi.model.Claim

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **int** |  | 
**groupId** | **int** | The group the claim and all its receipts belong to | 
**name** | **String** |  | 
**referenceNumber** | **String** | A travel authorization or voucher number. Absent when not set. | [optional] 
**type** | **String** | Free text, e.g. \"Travel voucher\" or \"TQSA\". Absent when not set. | [optional] 
**status** | **String** | One of DRAFT, SUBMITTED, PARTIALLY_PAID, PAID, DISPUTED, CLOSED. Set by hand; any status may follow any other. A plain string, not an enum, so a client must tolerate values it does not know. | 
**awaitingSettlement** | **bool** | Whether the status is one the server treats as waiting on the other party (SUBMITTED, PARTIALLY_PAID, DISPUTED). Drives daysWaiting and the list summary, so a client need not hard-code the set. | 
**submittedDate** | **String** | Calendar day the claim was submitted (yyyy-MM-dd). Absent when not set. | [optional] 
**dueDate** | **String** | Calendar day a response or payment is due (yyyy-MM-dd). Absent when not set. | [optional] 
**claimedAmount** | **String** | The amount explicitly entered as claimed, as a decimal string. Absent when not set, in which case effectiveClaimedAmount falls back to receiptsTotal. | [optional] 
**notes** | **String** | Free-text notes. Absent when not set. | [optional] 
**createdAt** | **String** | RFC 3339 instant | 
**updatedAt** | **String** | RFC 3339 instant | 
**createdBy** | **int** | The user who created the claim | [optional] 
**receiptIds** | **BuiltList&lt;int&gt;** | Ids of the linked receipts the caller may see, ascending. Fetch the receipts themselves through the receipt endpoints. Always present; empty when none. | 
**receiptCount** | **int** | Number of entries in receiptIds | 
**receiptsTotal** | **String** | Sum of the amounts of the receipts in receiptIds | 
**effectiveClaimedAmount** | **String** | claimedAmount when set, otherwise receiptsTotal | 
**paymentsTotal** | **String** | Sum of the amounts of the PAYMENT ledger entries | 
**recoupmentsTotal** | **String** | Sum of the amounts of the RECOUPMENT ledger entries | 
**paidAmount** | **String** | paymentsTotal minus recoupmentsTotal; negative when more was taken back than paid | 
**outstandingAmount** | **String** | effectiveClaimedAmount minus paidAmount. Negative when more was received than claimed. | 
**daysWaiting** | **int** | Calendar days from submittedDate to today in the app time zone, when awaitingSettlement is true and submittedDate is set; otherwise absent. Never negative. | [optional] 
**ledgerEntries** | [**BuiltList&lt;ClaimLedgerEntry&gt;**](ClaimLedgerEntry.md) | The claim's ledger, oldest first (by date, then id). Always present; empty when none. | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


