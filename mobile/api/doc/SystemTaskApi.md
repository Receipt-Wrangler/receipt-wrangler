# openapi.api.SystemTaskApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getPagedActivities**](SystemTaskApi.md#getpagedactivities) | **POST** /systemTask/getPagedActivities | Gets paged activities
[**getPagedSystemTasks**](SystemTaskApi.md#getpagedsystemtasks) | **POST** /systemTask/getPagedSystemTasks | Gets paged system tasks
[**rerunActivity**](SystemTaskApi.md#rerunactivity) | **POST** /systemTask/rerunActivity/{id} | Attempts to rerun activity


# **getPagedActivities**
> PagedData getPagedActivities(pagedActivityRequestCommand)

Gets paged activities

This will return paged activities for a list of groups

### Example
```dart
import 'package:openapi/api.dart';

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

[bearerAuth](../README.md#bearerAuth)

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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **rerunActivity**
> rerunActivity(id)

Attempts to rerun activity

This will rerun a failed activity

### Example
```dart
import 'package:openapi/api.dart';

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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

