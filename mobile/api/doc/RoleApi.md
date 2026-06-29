# openapi.api.RoleApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createRole**](RoleApi.md#createrole) | **POST** /role | Create role
[**deleteRole**](RoleApi.md#deleterole) | **DELETE** /role/{roleId} | Delete role
[**getRoles**](RoleApi.md#getroles) | **GET** /role | List all roles
[**setDefaultRole**](RoleApi.md#setdefaultrole) | **PUT** /role/{roleId}/default | Set default role
[**updateRole**](RoleApi.md#updaterole) | **PUT** /role/{roleId} | Update role


# **createRole**
> Role createRole(upsertRoleCommand)

Create role

This will create an app-scoped or group-scoped role.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getRoleApi();
final UpsertRoleCommand upsertRoleCommand = ; // UpsertRoleCommand | Role to create

try {
    final response = api.createRole(upsertRoleCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling RoleApi->createRole: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **upsertRoleCommand** | [**UpsertRoleCommand**](UpsertRoleCommand.md)| Role to create | 

### Return type

[**Role**](Role.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteRole**
> deleteRole(scope, roleId)

Delete role

Deletes an app-scoped or group-scoped role. The scope query parameter disambiguates the role id, since app and group role ids overlap. A role cannot be deleted while it is assigned to any user or group member, and system roles cannot be deleted.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getRoleApi();
final PermissionScope scope = ; // PermissionScope | Scope of the role to delete
final int roleId = 56; // int | Id of the role to update

try {
    api.deleteRole(scope, roleId);
} catch on DioException (e) {
    print('Exception when calling RoleApi->deleteRole: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **scope** | [**PermissionScope**](.md)| Scope of the role to delete | 
 **roleId** | **int**| Id of the role to update | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getRoles**
> BuiltList<Role> getRoles()

List all roles

Returns the full pool of roles, both app-scoped and group-scoped.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getRoleApi();

try {
    final response = api.getRoles();
    print(response);
} catch on DioException (e) {
    print('Exception when calling RoleApi->getRoles: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**BuiltList&lt;Role&gt;**](Role.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **setDefaultRole**
> Role setDefaultRole(scope, roleId)

Set default role

Makes the role the default for its scope — the role assigned to new accounts (APP) or to group creators (GROUP) — clearing the previous default. The scope query parameter disambiguates the role id, since app and group role ids overlap. System roles may be made default.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getRoleApi();
final PermissionScope scope = ; // PermissionScope | Scope of the role to make default
final int roleId = 56; // int | Id of the role to make the default for its scope

try {
    final response = api.setDefaultRole(scope, roleId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling RoleApi->setDefaultRole: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **scope** | [**PermissionScope**](.md)| Scope of the role to make default | 
 **roleId** | **int**| Id of the role to make the default for its scope | 

### Return type

[**Role**](Role.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateRole**
> Role updateRole(roleId, upsertRoleCommand)

Update role

Updates an existing app-scoped or group-scoped role. A role's type cannot be changed: the scope in the request must match the existing role's scope. System roles cannot be modified.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getRoleApi();
final int roleId = 56; // int | Id of the role to update
final UpsertRoleCommand upsertRoleCommand = ; // UpsertRoleCommand | Role to update

try {
    final response = api.updateRole(roleId, upsertRoleCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling RoleApi->updateRole: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **roleId** | **int**| Id of the role to update | 
 **upsertRoleCommand** | [**UpsertRoleCommand**](UpsertRoleCommand.md)| Role to update | 

### Return type

[**Role**](Role.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

