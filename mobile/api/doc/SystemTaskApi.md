# openapi.api.SystemTaskApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**downloadSystemTaskSourceFile**](SystemTaskApi.md#downloadsystemtasksourcefile) | **GET** /systemTask/{id}/sourceFile/download | Download an activity&#39;s source file
[**getPagedActivities**](SystemTaskApi.md#getpagedactivities) | **POST** /systemTask/getPagedActivities | Gets paged activities
[**getPagedSystemTasks**](SystemTaskApi.md#getpagedsystemtasks) | **POST** /systemTask/getPagedSystemTasks | Gets paged system tasks
[**getSystemTaskSourceFile**](SystemTaskApi.md#getsystemtasksourcefile) | **GET** /systemTask/{id}/sourceFile | Get an activity&#39;s source file
[**rerunActivity**](SystemTaskApi.md#rerunactivity) | **POST** /systemTask/rerunActivity/{id} | Attempts to rerun activity


# **downloadSystemTaskSourceFile**
> Uint8List downloadSystemTaskSourceFile(id)

Download an activity's source file

Returns the upload behind a quick scan or email upload activity verbatim, so a user can enter a failed receipt by hand.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getSystemTaskApi();
final int id = 56; // int | Id of the system task

try {
    final response = api.downloadSystemTaskSourceFile(id);
    print(response);
} catch on DioException (e) {
    print('Exception when calling SystemTaskApi->downloadSystemTaskSourceFile: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of the system task | 

### Return type

[**Uint8List**](Uint8List.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/octet-stream, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPagedActivities**
> PagedData getPagedActivities(pagedActivityRequestCommand)

Gets paged activities

This will return paged activities for a list of groups

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getSystemTaskApi();
final PagedActivityRequestCommand pagedActivityRequestCommand = ; // PagedActivityRequestCommand | Paging and sorting data

try {
    final response = api.getPagedActivities(pagedActivityRequestCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling SystemTaskApi->getPagedActivities: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pagedActivityRequestCommand** | [**PagedActivityRequestCommand**](PagedActivityRequestCommand.md)| Paging and sorting data | 

### Return type

[**PagedData**](PagedData.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPagedSystemTasks**
> PagedData getPagedSystemTasks(getSystemTaskCommand)

Gets paged system tasks

This will return paged system tasks

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getSystemTaskApi();
final GetSystemTaskCommand getSystemTaskCommand = ; // GetSystemTaskCommand | Paging and sorting data

try {
    final response = api.getPagedSystemTasks(getSystemTaskCommand);
    print(response);
} catch on DioException (e) {
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

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getSystemTaskSourceFile**
> SystemTaskSourceFileView getSystemTaskSourceFile(id)

Get an activity's source file

Returns the upload behind a quick scan or email upload activity, converted for display, so a user can see what a failed activity was working from.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getSystemTaskApi();
final int id = 56; // int | Id of the system task

try {
    final response = api.getSystemTaskSourceFile(id);
    print(response);
} catch on DioException (e) {
    print('Exception when calling SystemTaskApi->getSystemTaskSourceFile: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of the system task | 

### Return type

[**SystemTaskSourceFileView**](SystemTaskSourceFileView.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **rerunActivity**
> rerunActivity(id)

Attempts to rerun activity

This will rerun a failed activity

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getSystemTaskApi();
final int id = 56; // int | Id of activity to restart

try {
    api.rerunActivity(id);
} catch on DioException (e) {
    print('Exception when calling SystemTaskApi->rerunActivity: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of activity to restart | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

