# openapi.api.SearchApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**receiptSearch**](SearchApi.md#receiptsearch) | **GET** /search/ | Receipt Search


# **receiptSearch**
> BuiltList<SearchResult> receiptSearch(searchTerm)

Receipt Search

This will search for receipts based on a search term

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getSearchApi();
final String searchTerm = searchTerm_example; // String | search term

try {
    final response = api.receiptSearch(searchTerm);
    print(response);
} catch on DioException (e) {
    print('Exception when calling SearchApi->receiptSearch: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **searchTerm** | **String**| search term | 

### Return type

[**BuiltList&lt;SearchResult&gt;**](SearchResult.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

