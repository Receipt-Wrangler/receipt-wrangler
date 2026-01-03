# openapi.api.SystemTaskApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getPagedSystemTasks**](SystemTaskApi.md#getpagedsystemtasks) | **POST** /systemTask/getPagedSystemTasks | Gets paged system tasks


# **getPagedSystemTasks**
> PagedData getPagedSystemTasks(getSystemTaskCommand)

Gets paged system tasks

This will return paged system tasks

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = SystemTaskApi();
final getSystemTaskCommand = GetSystemTaskCommand(); // GetSystemTaskCommand | Paging and sorting data

try {
    final result = api_instance.getPagedSystemTasks(getSystemTaskCommand);
    print(result);
} catch (e) {
    print('Exception when calling SystemTaskApi->getPagedSystemTasks: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **getSystemTaskCommand** | [**GetSystemTaskCommand**](GetSystemTaskCommand.md)| Paging and sorting data | 

### Return type

[**PagedData**](PagedData.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

