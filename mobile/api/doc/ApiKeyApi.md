# openapi.api.ApiKeyApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createApiKey**](ApiKeyApi.md#createapikey) | **POST** /apiKey/ | Create API key
[**deleteApiKey**](ApiKeyApi.md#deleteapikey) | **DELETE** /apiKey/{id} | Delete API key
[**getPagedApiKeys**](ApiKeyApi.md#getpagedapikeys) | **POST** /apiKey/paged | Get paged API keys
[**updateApiKey**](ApiKeyApi.md#updateapikey) | **PUT** /apiKey/{id} | Update API key


# **createApiKey**
> ApiKeyResult createApiKey(upsertApiKeyCommand)

Create API key

Create a new API key for the authenticated user

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getApiKeyApi();
final UpsertApiKeyCommand upsertApiKeyCommand = ; // UpsertApiKeyCommand | API key details

try {
    final response = api.createApiKey(upsertApiKeyCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ApiKeyApi->createApiKey: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **upsertApiKeyCommand** | [**UpsertApiKeyCommand**](UpsertApiKeyCommand.md)| API key details | 

### Return type

[**ApiKeyResult**](ApiKeyResult.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteApiKey**
> deleteApiKey(id)

Delete API key

Delete an API key. Admins can delete any API key, non-admins can only delete their own API keys.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getApiKeyApi();
final String id = id_example; // String | API key ID to update

try {
    api.deleteApiKey(id);
} catch on DioException (e) {
    print('Exception when calling ApiKeyApi->deleteApiKey: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| API key ID to update | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPagedApiKeys**
> PagedData getPagedApiKeys(pagedApiKeyRequestCommand)

Get paged API keys

This will return paged API keys for the authenticated user or all API keys for admins

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getApiKeyApi();
final PagedApiKeyRequestCommand pagedApiKeyRequestCommand = ; // PagedApiKeyRequestCommand | Paging and sorting data

try {
    final response = api.getPagedApiKeys(pagedApiKeyRequestCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ApiKeyApi->getPagedApiKeys: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pagedApiKeyRequestCommand** | [**PagedApiKeyRequestCommand**](PagedApiKeyRequestCommand.md)| Paging and sorting data | 

### Return type

[**PagedData**](PagedData.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateApiKey**
> updateApiKey(id, upsertApiKeyCommand)

Update API key

This will update an API key. Users can only update their own API keys.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getApiKeyApi();
final String id = id_example; // String | API key ID to update
final UpsertApiKeyCommand upsertApiKeyCommand = ; // UpsertApiKeyCommand | API key details to update

try {
    api.updateApiKey(id, upsertApiKeyCommand);
} catch on DioException (e) {
    print('Exception when calling ApiKeyApi->updateApiKey: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**| API key ID to update | 
 **upsertApiKeyCommand** | [**UpsertApiKeyCommand**](UpsertApiKeyCommand.md)| API key details to update | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

