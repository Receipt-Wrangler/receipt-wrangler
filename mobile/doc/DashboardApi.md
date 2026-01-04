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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = DashboardApi();
final upsertDashboardCommand = UpsertDashboardCommand(); // UpsertDashboardCommand | Dashboard

try {
    final result = api_instance.createDashboard(upsertDashboardCommand);
    print(result);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = DashboardApi();
final dashboardId = 56; // int | Id of dashboard to operate on

try {
    final result = api_instance.deleteDashboard(dashboardId);
    print(result);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getDashboardsForUserByGroupId**
> List<Dashboard> getDashboardsForUserByGroupId(groupId)

Get dashboards for a user by group id

This will get a dashboards for a user by group id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = DashboardApi();
final groupId = groupId_example; // String | Id of group to get dashboard for

try {
    final result = api_instance.getDashboardsForUserByGroupId(groupId);
    print(result);
} catch (e) {
    print('Exception when calling DashboardApi->getDashboardsForUserByGroupId: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **String**| Id of group to get dashboard for | 

### Return type

[**List<Dashboard>**](Dashboard.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = DashboardApi();
final dashboardId = 56; // int | Id of dashboard to operate on
final upsertDashboardCommand = UpsertDashboardCommand(); // UpsertDashboardCommand | Dashboard to update

try {
    final result = api_instance.updateDashboard(dashboardId, upsertDashboardCommand);
    print(result);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

