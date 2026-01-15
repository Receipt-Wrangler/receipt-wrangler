# openapi.api.DashboardApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createDashboard**](DashboardApi.md#createdashboard) | **POST** /dashboard/ | Create dashboard
[**deleteDashboard**](DashboardApi.md#deletedashboard) | **DELETE** /dashboard/{dashboardId} | Delete dashboard
[**getDashboardsForUserByGroupId**](DashboardApi.md#getdashboardsforuserbygroupid) | **GET** /dashboard/{groupId} | Get dashboards for a user by group id
[**updateDashboard**](DashboardApi.md#updatedashboard) | **PUT** /dashboard/{dashboardId} | Update dashboard


# **createDashboard**
> Dashboard createDashboard(upsertDashboardCommand)

Create dashboard

This will create a dashboard [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getDashboardApi();
final UpsertDashboardCommand upsertDashboardCommand = ; // UpsertDashboardCommand | Dashboard

try {
    final response = api.createDashboard(upsertDashboardCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling DashboardApi->createDashboard: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **upsertDashboardCommand** | [**UpsertDashboardCommand**](UpsertDashboardCommand.md)| Dashboard | 

### Return type

[**Dashboard**](Dashboard.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteDashboard**
> Dashboard deleteDashboard(dashboardId)

Delete dashboard

This will delete a dashboard by id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getDashboardApi();
final int dashboardId = 56; // int | Id of dashboard to operate on

try {
    final response = api.deleteDashboard(dashboardId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling DashboardApi->deleteDashboard: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **dashboardId** | **int**| Id of dashboard to operate on | 

### Return type

[**Dashboard**](Dashboard.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getDashboardsForUserByGroupId**
> BuiltList<Dashboard> getDashboardsForUserByGroupId(groupId)

Get dashboards for a user by group id

This will get a dashboards for a user by group id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getDashboardApi();
final String groupId = groupId_example; // String | Id of group to get dashboard for

try {
    final response = api.getDashboardsForUserByGroupId(groupId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling DashboardApi->getDashboardsForUserByGroupId: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **String**| Id of group to get dashboard for | 

### Return type

[**BuiltList&lt;Dashboard&gt;**](Dashboard.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateDashboard**
> Dashboard updateDashboard(dashboardId, upsertDashboardCommand)

Update dashboard

This will update a dashboard

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getDashboardApi();
final int dashboardId = 56; // int | Id of dashboard to operate on
final UpsertDashboardCommand upsertDashboardCommand = ; // UpsertDashboardCommand | Dashboard to update

try {
    final response = api.updateDashboard(dashboardId, upsertDashboardCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling DashboardApi->updateDashboard: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **dashboardId** | **int**| Id of dashboard to operate on | 
 **upsertDashboardCommand** | [**UpsertDashboardCommand**](UpsertDashboardCommand.md)| Dashboard to update | 

### Return type

[**Dashboard**](Dashboard.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

