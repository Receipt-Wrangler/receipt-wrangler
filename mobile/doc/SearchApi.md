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
> List<SearchResult> receiptSearch(searchTerm)

Receipt Search

This will search for receipts based on a search term

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = SearchApi();
final searchTerm = searchTerm_example; // String | search term

try {
    final result = api_instance.receiptSearch(searchTerm);
    print(result);
} catch (e) {
    print('Exception when calling SearchApi->receiptSearch: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **searchTerm** | **String**| search term | 

### Return type

[**List<SearchResult>**](SearchResult.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

