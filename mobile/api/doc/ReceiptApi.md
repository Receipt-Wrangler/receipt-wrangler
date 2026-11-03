# openapi.api.ReceiptApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**bulkDeleteReceipts**](ReceiptApi.md#bulkdeletereceipts) | **POST** /receipt/bulkDelete | Bulk delete receipts
[**bulkReceiptStatusUpdate**](ReceiptApi.md#bulkreceiptstatusupdate) | **POST** /receipt/bulkStatusUpdate | Bulk receipt status update
[**createReceipt**](ReceiptApi.md#createreceipt) | **POST** /receipt/ | Create receipt
[**createReceiptWithFiles**](ReceiptApi.md#createreceiptwithfiles) | **POST** /receipt/withFiles | Create receipt with files
[**deleteReceiptById**](ReceiptApi.md#deletereceiptbyid) | **DELETE** /receipt/{receiptId} | Delete receipt
[**duplicateReceipt**](ReceiptApi.md#duplicatereceipt) | **POST** /receipt/{receiptId}/duplicate | Duplicate receipt
[**getReceiptById**](ReceiptApi.md#getreceiptbyid) | **GET** /receipt/{receiptId} | Get receipt
[**getReceiptHistory**](ReceiptApi.md#getreceipthistory) | **GET** /receipt/{receiptId}/history | Get receipt change history
[**getReceiptHistorySnapshots**](ReceiptApi.md#getreceipthistorysnapshots) | **GET** /receipt/{receiptId}/history/{systemTaskId}/snapshots | Get the receipt snapshots behind one history entry
[**getReceiptSummaryForGroup**](ReceiptApi.md#getreceiptsummaryforgroup) | **POST** /receipt/group/{groupId}/summary | Gets the receipt summary for a group
[**getReceiptsForGroup**](ReceiptApi.md#getreceiptsforgroup) | **POST** /receipt/group/{groupId} | Gets receipts
[**hasAccessToReceipt**](ReceiptApi.md#hasaccesstoreceipt) | **GET** /receipt/hasAccess | Has access to receipt
[**quickScanReceipt**](ReceiptApi.md#quickscanreceipt) | **POST** /receipt/quickScan | Quick scan a receipt
[**updateReceipt**](ReceiptApi.md#updatereceipt) | **PUT** /receipt/{receiptId} | Update receipt


# **bulkDeleteReceipts**
> BulkReceiptDeleteResult bulkDeleteReceipts(bulkReceiptDeleteCommand)

Bulk delete receipts

This will delete every listed receipt the caller may delete (group.receipts.delete on the receipt's own group) and report each one it did not delete. Receipts are deleted independently, so a refused receipt does not stop the others [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final BulkReceiptDeleteCommand bulkReceiptDeleteCommand = ; // BulkReceiptDeleteCommand | Receipt ids to delete

try {
    final response = api.bulkDeleteReceipts(bulkReceiptDeleteCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->bulkDeleteReceipts: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **bulkReceiptDeleteCommand** | [**BulkReceiptDeleteCommand**](BulkReceiptDeleteCommand.md)| Receipt ids to delete | 

### Return type

[**BulkReceiptDeleteResult**](BulkReceiptDeleteResult.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **bulkReceiptStatusUpdate**
> BuiltList<Receipt> bulkReceiptStatusUpdate(bulkStatusUpdateCommand)

Bulk receipt status update

This will bulk update receipt statuses with the option of adding a comment to each [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final BulkStatusUpdateCommand bulkStatusUpdateCommand = ; // BulkStatusUpdateCommand | Bulk status data

try {
    final response = api.bulkReceiptStatusUpdate(bulkStatusUpdateCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->bulkReceiptStatusUpdate: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **bulkStatusUpdateCommand** | [**BulkStatusUpdateCommand**](BulkStatusUpdateCommand.md)| Bulk status data | 

### Return type

[**BuiltList&lt;Receipt&gt;**](Receipt.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createReceipt**
> Receipt createReceipt(upsertReceiptCommand)

Create receipt

This will create a receipt [SYSTEM USER]. Deprecated in favour of createReceiptWithFiles, which carries the receipt's images in the same call; kept for already-released clients. It enforces the caller's role-required fields, so when the caller's group role requires an image this endpoint always returns 400 (it cannot carry one).

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final UpsertReceiptCommand upsertReceiptCommand = ; // UpsertReceiptCommand | Receipt to create

try {
    final response = api.createReceipt(upsertReceiptCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->createReceipt: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **upsertReceiptCommand** | [**UpsertReceiptCommand**](UpsertReceiptCommand.md)| Receipt to create | 

### Return type

[**Receipt**](Receipt.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createReceiptWithFiles**
> Receipt createReceiptWithFiles(receipt, files)

Create receipt with files

Creates a receipt, its comments and its images in one atomic call: a failure anywhere leaves no receipt, image or file behind. Requires group.receipts.create in the receipt's group. Every file must be an image or PDF (checked before anything is written). Returns 400 when the caller's group role requires a comment (key `comments`) or an image (key `files`) that the request does not carry. [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final UpsertReceiptCommand receipt = ; // UpsertReceiptCommand | 
final BuiltList<MultipartFile> files = /path/to/file.txt; // BuiltList<MultipartFile> | Images (or PDFs) to attach to the new receipt. Together with the receipt they share the server's 50 MB request body limit.

try {
    final response = api.createReceiptWithFiles(receipt, files);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->createReceiptWithFiles: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receipt** | [**UpsertReceiptCommand**](UpsertReceiptCommand.md)|  | 
 **files** | [**BuiltList&lt;MultipartFile&gt;**](MultipartFile.md)| Images (or PDFs) to attach to the new receipt. Together with the receipt they share the server's 50 MB request body limit. | [optional] 

### Return type

[**Receipt**](Receipt.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteReceiptById**
> deleteReceiptById(receiptId)

Delete receipt

This will delete a receipt by id [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final int receiptId = 56; // int | Id of receipt to get

try {
    api.deleteReceiptById(receiptId);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->deleteReceiptById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receiptId** | **int**| Id of receipt to get | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **duplicateReceipt**
> duplicateReceipt(receiptId)

Duplicate receipt

This will duplicate a receipt [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final int receiptId = 56; // int | Id of receipt to duplicate

try {
    api.duplicateReceipt(receiptId);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->duplicateReceipt: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receiptId** | **int**| Id of receipt to duplicate | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReceiptById**
> Receipt getReceiptById(receiptId)

Get receipt

This will get a receipt by receipt id [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final int receiptId = 56; // int | Id of receipt to get

try {
    final response = api.getReceiptById(receiptId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->getReceiptById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receiptId** | **int**| Id of receipt to get | 

### Return type

[**Receipt**](Receipt.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReceiptHistory**
> ReceiptHistory getReceiptHistory(receiptId)

Get receipt change history

Returns the receipt's change history, newest first: its creation, every recorded change and every failed save. Built from the receipt's system tasks and redacted to what the caller may see on the receipt today (category/tag grants, member isolation, paid-by visibility). Requires app.receipt-history.read and group.receipts.read on the receipt's group [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final int receiptId = 56; // int | Id of the receipt

try {
    final response = api.getReceiptHistory(receiptId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->getReceiptHistory: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receiptId** | **int**| Id of the receipt | 

### Return type

[**ReceiptHistory**](ReceiptHistory.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReceiptHistorySnapshots**
> ReceiptHistorySnapshots getReceiptHistorySnapshots(receiptId, systemTaskId)

Get the receipt snapshots behind one history entry

Returns the receipt as it was before and after one history entry, redacted exactly as the history is. Loaded on demand so the history list stays light. Same permissions as getReceiptHistory; an entry the caller cannot see is a 403 [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final int receiptId = 56; // int | Id of the receipt
final int systemTaskId = 56; // int | systemTaskId of the history entry

try {
    final response = api.getReceiptHistorySnapshots(receiptId, systemTaskId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->getReceiptHistorySnapshots: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receiptId** | **int**| Id of the receipt | 
 **systemTaskId** | **int**| systemTaskId of the history entry | 

### Return type

[**ReceiptHistorySnapshots**](ReceiptHistorySnapshots.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReceiptSummaryForGroup**
> ReceiptSummary getReceiptSummaryForGroup(groupId, receiptSummaryCommand)

Gets the receipt summary for a group

Returns the block of totals rendered under the receipts table: a receipt count and amount total over the WHOLE filtered result set (not the current page), then the same figures per configured status. Which statuses break out and which currency custom fields are totalled come from the group's receipt settings, not from the request. Gated on group.receipts.read, the same permission as the receipts it aggregates. [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final int groupId = 56; // int | Summarize the receipts that belong to groupId
final ReceiptSummaryCommand receiptSummaryCommand = ; // ReceiptSummaryCommand | 

try {
    final response = api.getReceiptSummaryForGroup(groupId, receiptSummaryCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->getReceiptSummaryForGroup: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Summarize the receipts that belong to groupId | 
 **receiptSummaryCommand** | [**ReceiptSummaryCommand**](ReceiptSummaryCommand.md)|  | 

### Return type

[**ReceiptSummary**](ReceiptSummary.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReceiptsForGroup**
> PagedData getReceiptsForGroup(groupId, receiptPagedRequestCommand)

Gets receipts

This will return receipts with the option to sort and filter [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final int groupId = 56; // int | Get all receipts that belong to groupId
final ReceiptPagedRequestCommand receiptPagedRequestCommand = ; // ReceiptPagedRequestCommand | 

try {
    final response = api.getReceiptsForGroup(groupId, receiptPagedRequestCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->getReceiptsForGroup: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Get all receipts that belong to groupId | 
 **receiptPagedRequestCommand** | [**ReceiptPagedRequestCommand**](ReceiptPagedRequestCommand.md)|  | 

### Return type

[**PagedData**](PagedData.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **hasAccessToReceipt**
> hasAccessToReceipt(receiptId, permission)

Has access to receipt

This will return whether or not the currently logged in user has access to the receipt

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final int receiptId = 56; // int | 
final Permission permission = ; // Permission | Group permission required to access the receipt

try {
    api.hasAccessToReceipt(receiptId, permission);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->hasAccessToReceipt: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receiptId** | **int**|  | 
 **permission** | [**Permission**](.md)| Group permission required to access the receipt | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **quickScanReceipt**
> quickScanReceipt(files, groupIds, paidByUserIds, statuses, categoryIds, tagIds, comments)

Quick scan a receipt

This take an image and use magic fill to fill and save the receipt [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final BuiltList<MultipartFile> files = /path/to/file.txt; // BuiltList<MultipartFile> | 
final BuiltList<int> groupIds = ; // BuiltList<int> | 
final BuiltList<int> paidByUserIds = ; // BuiltList<int> | 
final BuiltList<ReceiptStatus> statuses = ; // BuiltList<ReceiptStatus> | 
final BuiltList<String> categoryIds = ; // BuiltList<String> | 
final BuiltList<String> tagIds = ; // BuiltList<String> | 
final BuiltList<String> comments = ; // BuiltList<String> | 

try {
    api.quickScanReceipt(files, groupIds, paidByUserIds, statuses, categoryIds, tagIds, comments);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->quickScanReceipt: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **files** | [**BuiltList&lt;MultipartFile&gt;**](MultipartFile.md)|  | 
 **groupIds** | [**BuiltList&lt;int&gt;**](int.md)|  | 
 **paidByUserIds** | [**BuiltList&lt;int&gt;**](int.md)|  | 
 **statuses** | [**BuiltList&lt;ReceiptStatus&gt;**](ReceiptStatus.md)|  | 
 **categoryIds** | [**BuiltList&lt;String&gt;**](String.md)|  | [optional] 
 **tagIds** | [**BuiltList&lt;String&gt;**](String.md)|  | [optional] 
 **comments** | [**BuiltList&lt;String&gt;**](String.md)|  | [optional] 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateReceipt**
> Receipt updateReceipt(receiptId, upsertReceiptCommand)

Update receipt

This will update a receipt by receipt id [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getReceiptApi();
final int receiptId = 56; // int | Id of receipt to get
final UpsertReceiptCommand upsertReceiptCommand = ; // UpsertReceiptCommand | Receipt to update

try {
    final response = api.updateReceipt(receiptId, upsertReceiptCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReceiptApi->updateReceipt: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receiptId** | **int**| Id of receipt to get | 
 **upsertReceiptCommand** | [**UpsertReceiptCommand**](UpsertReceiptCommand.md)| Receipt to update | 

### Return type

[**Receipt**](Receipt.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

