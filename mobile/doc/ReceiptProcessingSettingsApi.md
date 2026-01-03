# openapi.api.ReceiptProcessingSettingsApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**checkReceiptProcessingSettingsConnectivity**](ReceiptProcessingSettingsApi.md#checkreceiptprocessingsettingsconnectivity) | **POST** /receiptProcessingSettings/checkConnectivity | Check receipt processing settings connectivity
[**createReceiptProcessingSettings**](ReceiptProcessingSettingsApi.md#createreceiptprocessingsettings) | **POST** /receiptProcessingSettings | Create receipt processing settings
[**deleteReceiptProcessingSettingsById**](ReceiptProcessingSettingsApi.md#deletereceiptprocessingsettingsbyid) | **DELETE** /receiptProcessingSettings/{id} | Delete receipt processing settings by id
[**getPagedProcessingSettings**](ReceiptProcessingSettingsApi.md#getpagedprocessingsettings) | **POST** /receiptProcessingSettings/getPagedProcessingSettings | Gets paged processing settings
[**getReceiptProcessingSettingsById**](ReceiptProcessingSettingsApi.md#getreceiptprocessingsettingsbyid) | **GET** /receiptProcessingSettings/{id} | Get receipt processing settings by id
[**updateReceiptProcessingSettingsById**](ReceiptProcessingSettingsApi.md#updatereceiptprocessingsettingsbyid) | **PUT** /receiptProcessingSettings/{id} | Update receipt processing settings by id


# **checkReceiptProcessingSettingsConnectivity**
> SystemTask checkReceiptProcessingSettingsConnectivity(checkReceiptProcessingSettingsConnectivityCommand)

Check receipt processing settings connectivity

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptProcessingSettingsApi();
final checkReceiptProcessingSettingsConnectivityCommand = CheckReceiptProcessingSettingsConnectivityCommand(); // CheckReceiptProcessingSettingsConnectivityCommand | Receipt processing settings to check connectivity

try {
    final result = api_instance.checkReceiptProcessingSettingsConnectivity(checkReceiptProcessingSettingsConnectivityCommand);
    print(result);
} catch (e) {
    print('Exception when calling ReceiptProcessingSettingsApi->checkReceiptProcessingSettingsConnectivity: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **checkReceiptProcessingSettingsConnectivityCommand** | [**CheckReceiptProcessingSettingsConnectivityCommand**](CheckReceiptProcessingSettingsConnectivityCommand.md)| Receipt processing settings to check connectivity | 

### Return type

[**SystemTask**](SystemTask.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createReceiptProcessingSettings**
> ReceiptProcessingSettings createReceiptProcessingSettings(upsertReceiptProcessingSettingsCommand)

Create receipt processing settings

This will create receipt processing settings

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptProcessingSettingsApi();
final upsertReceiptProcessingSettingsCommand = UpsertReceiptProcessingSettingsCommand(); // UpsertReceiptProcessingSettingsCommand | Receipt processing settings to create

try {
    final result = api_instance.createReceiptProcessingSettings(upsertReceiptProcessingSettingsCommand);
    print(result);
} catch (e) {
    print('Exception when calling ReceiptProcessingSettingsApi->createReceiptProcessingSettings: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **upsertReceiptProcessingSettingsCommand** | [**UpsertReceiptProcessingSettingsCommand**](UpsertReceiptProcessingSettingsCommand.md)| Receipt processing settings to create | 

### Return type

[**ReceiptProcessingSettings**](ReceiptProcessingSettings.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteReceiptProcessingSettingsById**
> deleteReceiptProcessingSettingsById(id)

Delete receipt processing settings by id

This will delete receipt processing settings by id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptProcessingSettingsApi();
final id = 56; // int | Id of receipt processing settings to delete

try {
    api_instance.deleteReceiptProcessingSettingsById(id);
} catch (e) {
    print('Exception when calling ReceiptProcessingSettingsApi->deleteReceiptProcessingSettingsById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of receipt processing settings to delete | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPagedProcessingSettings**
> PagedData getPagedProcessingSettings(pagedRequestCommand)

Gets paged processing settings

This will return paged processing settings

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptProcessingSettingsApi();
final pagedRequestCommand = PagedRequestCommand(); // PagedRequestCommand | Paging and sorting data

try {
    final result = api_instance.getPagedProcessingSettings(pagedRequestCommand);
    print(result);
} catch (e) {
    print('Exception when calling ReceiptProcessingSettingsApi->getPagedProcessingSettings: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pagedRequestCommand** | [**PagedRequestCommand**](PagedRequestCommand.md)| Paging and sorting data | 

### Return type

[**PagedData**](PagedData.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReceiptProcessingSettingsById**
> ReceiptProcessingSettings getReceiptProcessingSettingsById(id)

Get receipt processing settings by id

This will get receipt processing settings by id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptProcessingSettingsApi();
final id = 56; // int | Id of receipt processing settings to get

try {
    final result = api_instance.getReceiptProcessingSettingsById(id);
    print(result);
} catch (e) {
    print('Exception when calling ReceiptProcessingSettingsApi->getReceiptProcessingSettingsById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of receipt processing settings to get | 

### Return type

[**ReceiptProcessingSettings**](ReceiptProcessingSettings.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateReceiptProcessingSettingsById**
> ReceiptProcessingSettings updateReceiptProcessingSettingsById(id, updateKey, upsertReceiptProcessingSettingsCommand)

Update receipt processing settings by id

This will update receipt processing settings by id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptProcessingSettingsApi();
final id = 56; // int | Id of receipt processing settings to update
final updateKey = true; // bool | Whether or not to update the key
final upsertReceiptProcessingSettingsCommand = UpsertReceiptProcessingSettingsCommand(); // UpsertReceiptProcessingSettingsCommand | Receipt processing settings to update

try {
    final result = api_instance.updateReceiptProcessingSettingsById(id, updateKey, upsertReceiptProcessingSettingsCommand);
    print(result);
} catch (e) {
    print('Exception when calling ReceiptProcessingSettingsApi->updateReceiptProcessingSettingsById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of receipt processing settings to update | 
 **updateKey** | **bool**| Whether or not to update the key | 
 **upsertReceiptProcessingSettingsCommand** | [**UpsertReceiptProcessingSettingsCommand**](UpsertReceiptProcessingSettingsCommand.md)| Receipt processing settings to update | 

### Return type

[**ReceiptProcessingSettings**](ReceiptProcessingSettings.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

