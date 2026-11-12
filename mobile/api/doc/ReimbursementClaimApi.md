# openapi.api.ReimbursementClaimApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**addReceiptsToClaim**](ReimbursementClaimApi.md#addreceiptstoclaim) | **POST** /claim/{claimId}/addReceipts | Add receipts to a claim
[**createClaim**](ReimbursementClaimApi.md#createclaim) | **POST** /claim/group/{groupId} | Create claim
[**createClaimLedgerEntry**](ReimbursementClaimApi.md#createclaimledgerentry) | **POST** /claim/{claimId}/ledgerEntry | Add a ledger entry to a claim
[**deleteClaim**](ReimbursementClaimApi.md#deleteclaim) | **DELETE** /claim/{claimId} | Delete claim
[**deleteClaimLedgerEntry**](ReimbursementClaimApi.md#deleteclaimledgerentry) | **DELETE** /claim/{claimId}/ledgerEntry/{ledgerEntryId} | Delete a ledger entry
[**getClaimById**](ReimbursementClaimApi.md#getclaimbyid) | **GET** /claim/{claimId} | Get claim
[**getClaimsForGroup**](ReimbursementClaimApi.md#getclaimsforgroup) | **GET** /claim/group/{groupId} | Get claims for a group
[**removeReceiptsFromClaim**](ReimbursementClaimApi.md#removereceiptsfromclaim) | **POST** /claim/{claimId}/removeReceipts | Remove receipts from a claim
[**updateClaim**](ReimbursementClaimApi.md#updateclaim) | **PUT** /claim/{claimId} | Update claim
[**updateClaimLedgerEntry**](ReimbursementClaimApi.md#updateclaimledgerentry) | **PUT** /claim/{claimId}/ledgerEntry/{ledgerEntryId} | Update a ledger entry


# **addReceiptsToClaim**
> ClaimReceiptsResult addReceiptsToClaim(claimId, claimReceiptsCommand)

Add receipts to a claim

Links each listed receipt to the claim and reports each one it could not link. A receipt must belong to the claim's group, be visible to the caller, and not already be in another claim (remove it from that claim first; nothing is moved silently). A receipt already in this claim counts as added. Receipts are judged independently, so one refusal does not stop the others. Gated on group.receipts.update in the claim's group. The response carries only ids the caller sent; fetch the claim again for its new figures [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReimbursementClaimApi();
final int claimId = 56; // int | Id of the claim to add receipts to
final ClaimReceiptsCommand claimReceiptsCommand = ; // ClaimReceiptsCommand | Receipt ids to add

try {
    final response = api.addReceiptsToClaim(claimId, claimReceiptsCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReimbursementClaimApi->addReceiptsToClaim: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **claimId** | **int**| Id of the claim to add receipts to | 
 **claimReceiptsCommand** | [**ClaimReceiptsCommand**](ClaimReceiptsCommand.md)| Receipt ids to add | 

### Return type

[**ClaimReceiptsResult**](ClaimReceiptsResult.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createClaim**
> Claim createClaim(groupId, upsertClaimCommand)

Create claim

Creates a claim in the group, with no receipts and no ledger entries; add receipts with addReceiptsToClaim. Gated on group.receipts.update. The synthetic \"All\" group is never a destination (400) [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReimbursementClaimApi();
final int groupId = 56; // int | Id of the group. On GET this may be the caller's synthetic \"All\" group, which lists the claims of every group where the caller holds group.receipts.read. POST rejects the All group with a 400.
final UpsertClaimCommand upsertClaimCommand = ; // UpsertClaimCommand | Claim to create

try {
    final response = api.createClaim(groupId, upsertClaimCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReimbursementClaimApi->createClaim: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Id of the group. On GET this may be the caller's synthetic \"All\" group, which lists the claims of every group where the caller holds group.receipts.read. POST rejects the All group with a 400. | 
 **upsertClaimCommand** | [**UpsertClaimCommand**](UpsertClaimCommand.md)| Claim to create | 

### Return type

[**Claim**](Claim.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createClaimLedgerEntry**
> ClaimLedgerEntry createClaimLedgerEntry(claimId, upsertClaimLedgerEntryCommand)

Add a ledger entry to a claim

Records a payment received, a recoupment or a note on the claim. Gated on group.receipts.update in the claim's group. Fetch the claim again for its new figures [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReimbursementClaimApi();
final int claimId = 56; // int | Id of the claim
final UpsertClaimLedgerEntryCommand upsertClaimLedgerEntryCommand = ; // UpsertClaimLedgerEntryCommand | Entry to add

try {
    final response = api.createClaimLedgerEntry(claimId, upsertClaimLedgerEntryCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReimbursementClaimApi->createClaimLedgerEntry: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **claimId** | **int**| Id of the claim | 
 **upsertClaimLedgerEntryCommand** | [**UpsertClaimLedgerEntryCommand**](UpsertClaimLedgerEntryCommand.md)| Entry to add | 

### Return type

[**ClaimLedgerEntry**](ClaimLedgerEntry.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteClaim**
> deleteClaim(claimId)

Delete claim

Deletes the claim with its ledger entries and receipt links. The receipts themselves are not deleted. Gated on group.receipts.update in the claim's group [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReimbursementClaimApi();
final int claimId = 56; // int | Id of the claim

try {
    api.deleteClaim(claimId);
} catch on DioException (e) {
    print('Exception when calling ReimbursementClaimApi->deleteClaim: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **claimId** | **int**| Id of the claim | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteClaimLedgerEntry**
> deleteClaimLedgerEntry(claimId, ledgerEntryId)

Delete a ledger entry

Deletes one entry. An entry that is not on this claim, or that the caller may not see, is a 403. Gated on group.receipts.update in the claim's group [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReimbursementClaimApi();
final int claimId = 56; // int | Id of the claim the entry belongs to
final int ledgerEntryId = 56; // int | Id of the ledger entry

try {
    api.deleteClaimLedgerEntry(claimId, ledgerEntryId);
} catch on DioException (e) {
    print('Exception when calling ReimbursementClaimApi->deleteClaimLedgerEntry: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **claimId** | **int**| Id of the claim the entry belongs to | 
 **ledgerEntryId** | **int**| Id of the ledger entry | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getClaimById**
> Claim getClaimById(claimId)

Get claim

Returns one claim with its computed figures, linked receipt ids and ledger entries. Gated on group.receipts.read in the claim's group. A missing claim and one the caller may not see are the same 403 [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReimbursementClaimApi();
final int claimId = 56; // int | Id of the claim

try {
    final response = api.getClaimById(claimId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReimbursementClaimApi->getClaimById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **claimId** | **int**| Id of the claim | 

### Return type

[**Claim**](Claim.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getClaimsForGroup**
> ClaimList getClaimsForGroup(groupId)

Get claims for a group

Returns the group's claims, newest first, each with its computed figures, linked receipt ids and ledger entries, plus a summary for an \"outstanding and ageing\" header. Every figure covers only what the caller may see (see Claim). Gated on group.receipts.read [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReimbursementClaimApi();
final int groupId = 56; // int | Id of the group. On GET this may be the caller's synthetic \"All\" group, which lists the claims of every group where the caller holds group.receipts.read. POST rejects the All group with a 400.

try {
    final response = api.getClaimsForGroup(groupId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReimbursementClaimApi->getClaimsForGroup: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Id of the group. On GET this may be the caller's synthetic \"All\" group, which lists the claims of every group where the caller holds group.receipts.read. POST rejects the All group with a 400. | 

### Return type

[**ClaimList**](ClaimList.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **removeReceiptsFromClaim**
> ClaimReceiptsResult removeReceiptsFromClaim(claimId, claimReceiptsCommand)

Remove receipts from a claim

Unlinks each listed receipt from the claim and reports each one it could not unlink. A receipt that is not in the claim and one the caller may not see read the same. The receipts themselves are not changed. Gated on group.receipts.update in the claim's group [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReimbursementClaimApi();
final int claimId = 56; // int | Id of the claim to remove receipts from
final ClaimReceiptsCommand claimReceiptsCommand = ; // ClaimReceiptsCommand | Receipt ids to remove

try {
    final response = api.removeReceiptsFromClaim(claimId, claimReceiptsCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReimbursementClaimApi->removeReceiptsFromClaim: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **claimId** | **int**| Id of the claim to remove receipts from | 
 **claimReceiptsCommand** | [**ClaimReceiptsCommand**](ClaimReceiptsCommand.md)| Receipt ids to remove | 

### Return type

[**ClaimReceiptsResult**](ClaimReceiptsResult.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateClaim**
> Claim updateClaim(claimId, upsertClaimCommand)

Update claim

Replaces the claim's own fields (name, reference number, type, status, dates, claimed amount, notes) with the body: an optional field that is omitted is cleared. Linked receipts and ledger entries are not part of the body and are never touched here; they have their own endpoints. A claim never moves between groups. Gated on group.receipts.update in the claim's group [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReimbursementClaimApi();
final int claimId = 56; // int | Id of the claim
final UpsertClaimCommand upsertClaimCommand = ; // UpsertClaimCommand | The claim's new fields

try {
    final response = api.updateClaim(claimId, upsertClaimCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReimbursementClaimApi->updateClaim: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **claimId** | **int**| Id of the claim | 
 **upsertClaimCommand** | [**UpsertClaimCommand**](UpsertClaimCommand.md)| The claim's new fields | 

### Return type

[**Claim**](Claim.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateClaimLedgerEntry**
> ClaimLedgerEntry updateClaimLedgerEntry(claimId, ledgerEntryId, upsertClaimLedgerEntryCommand)

Update a ledger entry

Replaces the entry's type, date, amount and note with the body (an omitted optional field is cleared). An entry that is not on this claim, or that the caller may not see, is a 403. Gated on group.receipts.update in the claim's group [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReimbursementClaimApi();
final int claimId = 56; // int | Id of the claim the entry belongs to
final int ledgerEntryId = 56; // int | Id of the ledger entry
final UpsertClaimLedgerEntryCommand upsertClaimLedgerEntryCommand = ; // UpsertClaimLedgerEntryCommand | The entry's new fields

try {
    final response = api.updateClaimLedgerEntry(claimId, ledgerEntryId, upsertClaimLedgerEntryCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReimbursementClaimApi->updateClaimLedgerEntry: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **claimId** | **int**| Id of the claim the entry belongs to | 
 **ledgerEntryId** | **int**| Id of the ledger entry | 
 **upsertClaimLedgerEntryCommand** | [**UpsertClaimLedgerEntryCommand**](UpsertClaimLedgerEntryCommand.md)| The entry's new fields | 

### Return type

[**ClaimLedgerEntry**](ClaimLedgerEntry.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

