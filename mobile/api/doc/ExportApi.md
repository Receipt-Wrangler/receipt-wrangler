# openapi.api.ExportApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**exportReceiptsById**](ExportApi.md#exportreceiptsbyid) | **POST** /export | Exports receipts
[**exportReceiptsForGroup**](ExportApi.md#exportreceiptsforgroup) | **POST** /export/{groupId} | Exports receipts


# **exportReceiptsById**
> Uint8List exportReceiptsById(format, receiptIds)

Exports receipts

This will export individual receipts [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getExportApi();
final ExportFormat format = ; // ExportFormat | 
final BuiltList<int> receiptIds = ; // BuiltList<int> | 

try {
    final response = api.exportReceiptsById(format, receiptIds);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ExportApi->exportReceiptsById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **format** | [**ExportFormat**](.md)|  | 
 **receiptIds** | [**BuiltList&lt;int&gt;**](int.md)|  | [optional] 

### Return type

[**Uint8List**](Uint8List.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/octet-stream, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **exportReceiptsForGroup**
> Uint8List exportReceiptsForGroup(format, groupId, receiptPagedRequestCommand)

Exports receipts

This will export all receipts that belong to a group based on a filter [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getExportApi();
final ExportFormat format = ; // ExportFormat | 
final int groupId = 56; // int | Get all receipts that belong to groupId
final ReceiptPagedRequestCommand receiptPagedRequestCommand = ; // ReceiptPagedRequestCommand | 

try {
    final response = api.exportReceiptsForGroup(format, groupId, receiptPagedRequestCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ExportApi->exportReceiptsForGroup: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **format** | [**ExportFormat**](.md)|  | 
 **groupId** | **int**| Get all receipts that belong to groupId | 
 **receiptPagedRequestCommand** | [**ReceiptPagedRequestCommand**](ReceiptPagedRequestCommand.md)|  | 

### Return type

[**Uint8List**](Uint8List.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/octet-stream, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

