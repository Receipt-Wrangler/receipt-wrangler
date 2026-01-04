# openapi.api.ReceiptApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**bulkReceiptStatusUpdate**](ReceiptApi.md#bulkreceiptstatusupdate) | **POST** /receipt/bulkStatusUpdate | Bulk receipt status update
[**createReceipt**](ReceiptApi.md#createreceipt) | **POST** /receipt/ | Create receipt
[**deleteReceiptById**](ReceiptApi.md#deletereceiptbyid) | **DELETE** /receipt/{receiptId} | Delete receipt
[**duplicateReceipt**](ReceiptApi.md#duplicatereceipt) | **POST** /receipt/{receiptId}/duplicate | Duplicate receipt
[**getReceiptById**](ReceiptApi.md#getreceiptbyid) | **GET** /receipt/{receiptId} | Get receipt
[**getReceiptsForGroup**](ReceiptApi.md#getreceiptsforgroup) | **POST** /receipt/group/{groupId} | Gets receipts
[**hasAccessToReceipt**](ReceiptApi.md#hasaccesstoreceipt) | **GET** /receipt/hasAccess | Has access to receipt
[**quickScanReceipt**](ReceiptApi.md#quickscanreceipt) | **POST** /receipt/quickScan | Quick scan a receipt
[**updateReceipt**](ReceiptApi.md#updatereceipt) | **PUT** /receipt/{receiptId} | Update receipt


# **bulkReceiptStatusUpdate**
> List<Receipt> bulkReceiptStatusUpdate(bulkStatusUpdateCommand)

Bulk receipt status update

This will bulk update receipt statuses with the option of adding a comment to each [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptApi();
final bulkStatusUpdateCommand = BulkStatusUpdateCommand(); // BulkStatusUpdateCommand | Bulk status data

try {
    final result = api_instance.bulkReceiptStatusUpdate(bulkStatusUpdateCommand);
    print(result);
} catch (e) {
    print('Exception when calling ReceiptApi->bulkReceiptStatusUpdate: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **bulkStatusUpdateCommand** | [**BulkStatusUpdateCommand**](BulkStatusUpdateCommand.md)| Bulk status data | 

### Return type

[**List<Receipt>**](Receipt.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createReceipt**
> Receipt createReceipt(upsertReceiptCommand)

Create receipt

This will create a receipt [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptApi();
final upsertReceiptCommand = UpsertReceiptCommand(); // UpsertReceiptCommand | Receipt to create

try {
    final result = api_instance.createReceipt(upsertReceiptCommand);
    print(result);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteReceiptById**
> deleteReceiptById(receiptId)

Delete receipt

This will delete a receipt by id [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptApi();
final receiptId = 56; // int | Id of receipt to get

try {
    api_instance.deleteReceiptById(receiptId);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **duplicateReceipt**
> duplicateReceipt(receiptId)

Duplicate receipt

This will duplicate a receipt [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptApi();
final receiptId = 56; // int | Id of receipt to duplicate

try {
    api_instance.duplicateReceipt(receiptId);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReceiptById**
> Receipt getReceiptById(receiptId)

Get receipt

This will get a receipt by receipt id [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptApi();
final receiptId = 56; // int | Id of receipt to get

try {
    final result = api_instance.getReceiptById(receiptId);
    print(result);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReceiptsForGroup**
> PagedData getReceiptsForGroup(groupId, receiptPagedRequestCommand)

Gets receipts

This will return receipts with the option to sort and filter [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptApi();
final groupId = 56; // int | Get all receipts that belong to groupId
final receiptPagedRequestCommand = ReceiptPagedRequestCommand(); // ReceiptPagedRequestCommand | 

try {
    final result = api_instance.getReceiptsForGroup(groupId, receiptPagedRequestCommand);
    print(result);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **hasAccessToReceipt**
> hasAccessToReceipt(receiptId, groupRole)

Has access to receipt

This will return whether or not the currently logged in user has access to the receipt

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptApi();
final receiptId = 56; // int | 
final groupRole = groupRole_example; // String | Role required to have access to receipt

try {
    api_instance.hasAccessToReceipt(receiptId, groupRole);
} catch (e) {
    print('Exception when calling ReceiptApi->hasAccessToReceipt: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receiptId** | **int**|  | 
 **groupRole** | **String**| Role required to have access to receipt | [optional] 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **quickScanReceipt**
> List<Receipt> quickScanReceipt(files, groupIds, paidByUserIds, statuses)

Quick scan a receipt

This take an image and use magic fill to fill and save the receipt [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptApi();
final files = [/path/to/file.txt]; // List<MultipartFile> | 
final groupIds = []; // List<int> | 
final paidByUserIds = []; // List<int> | 
final statuses = []; // List<ReceiptStatus> | 

try {
    final result = api_instance.quickScanReceipt(files, groupIds, paidByUserIds, statuses);
    print(result);
} catch (e) {
    print('Exception when calling ReceiptApi->quickScanReceipt: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **files** | [**List<MultipartFile>**](MultipartFile.md)|  | 
 **groupIds** | [**List<int>**](int.md)|  | 
 **paidByUserIds** | [**List<int>**](int.md)|  | 
 **statuses** | [**List<ReceiptStatus>**](ReceiptStatus.md)|  | 

### Return type

[**List<Receipt>**](Receipt.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptApi();
final receiptId = 56; // int | Id of receipt to get
final upsertReceiptCommand = UpsertReceiptCommand(); // UpsertReceiptCommand | Receipt to update

try {
    final result = api_instance.updateReceipt(receiptId, upsertReceiptCommand);
    print(result);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

