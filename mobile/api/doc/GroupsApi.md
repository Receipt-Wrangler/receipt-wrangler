# openapi.api.GroupsApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createGroup**](GroupsApi.md#creategroup) | **POST** /group | Create group
[**deleteGroup**](GroupsApi.md#deletegroup) | **DELETE** /group/{groupId} | Delete group
[**getGroupById**](GroupsApi.md#getgroupbyid) | **GET** /group/{groupId} | Gets a group by Id
[**getGroupsForuser**](GroupsApi.md#getgroupsforuser) | **GET** /group | Get groups for user
[**getPagedGroups**](GroupsApi.md#getpagedgroups) | **POST** /group/getPagedGroups | Get paged groups
[**pollGroupEmail**](GroupsApi.md#pollgroupemail) | **POST** /group/{groupId}/pollGroupEmail | Poll group email
[**updateGroup**](GroupsApi.md#updategroup) | **PUT** /group/{groupId} | Update a group
[**updateGroupMemberGrants**](GroupsApi.md#updategroupmembergrants) | **PUT** /group/{groupId}/member/{userId}/grants | Update a group member&#39;s category and tag assignment
[**updateGroupReceiptSettings**](GroupsApi.md#updategroupreceiptsettings) | **PUT** /group/{groupId}/groupReceiptSettings | Update group receipt settings
[**updateGroupSettings**](GroupsApi.md#updategroupsettings) | **PUT** /group/{groupId}/groupSettings | Update group settings


# **createGroup**
> createGroup(upsertGroupCommand)

Create group

This will create a group

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getGroupsApi();
final UpsertGroupCommand upsertGroupCommand = ; // UpsertGroupCommand | Group to create

try {
    api.createGroup(upsertGroupCommand);
} catch on DioException (e) {
    print('Exception when calling GroupsApi->createGroup: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **upsertGroupCommand** | [**UpsertGroupCommand**](UpsertGroupCommand.md)| Group to create | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteGroup**
> deleteGroup(groupId)

Delete group

This will delete a group by id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getGroupsApi();
final int groupId = 56; // int | Group Id to get

try {
    api.deleteGroup(groupId);
} catch on DioException (e) {
    print('Exception when calling GroupsApi->deleteGroup: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Group Id to get | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getGroupById**
> getGroupById(groupId)

Gets a group by Id

This will get a group by Id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getGroupsApi();
final int groupId = 56; // int | Group Id to get

try {
    api.getGroupById(groupId);
} catch on DioException (e) {
    print('Exception when calling GroupsApi->getGroupById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Group Id to get | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getGroupsForuser**
> BuiltList<Group> getGroupsForuser()

Get groups for user

This will get groups for the currently logged in user

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getGroupsApi();

try {
    final response = api.getGroupsForuser();
    print(response);
} catch on DioException (e) {
    print('Exception when calling GroupsApi->getGroupsForuser: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**BuiltList&lt;Group&gt;**](Group.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPagedGroups**
> PagedData getPagedGroups(pagedGroupRequestCommand)

Get paged groups

This will return paged groups

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getGroupsApi();
final PagedGroupRequestCommand pagedGroupRequestCommand = ; // PagedGroupRequestCommand | Paging and sorting data

try {
    final response = api.getPagedGroups(pagedGroupRequestCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling GroupsApi->getPagedGroups: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pagedGroupRequestCommand** | [**PagedGroupRequestCommand**](PagedGroupRequestCommand.md)| Paging and sorting data | 

### Return type

[**PagedData**](PagedData.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **pollGroupEmail**
> pollGroupEmail(groupId)

Poll group email

This will poll the group email for new receipts and add them to the group

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getGroupsApi();
final int groupId = 56; // int | Group Id to poll

try {
    api.pollGroupEmail(groupId);
} catch on DioException (e) {
    print('Exception when calling GroupsApi->pollGroupEmail: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Group Id to poll | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateGroup**
> updateGroup(groupId, group)

Update a group

This will update a group

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getGroupsApi();
final int groupId = 56; // int | Group Id to get
final Group group = ; // Group | Group to update

try {
    api.updateGroup(groupId, group);
} catch on DioException (e) {
    print('Exception when calling GroupsApi->updateGroup: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Group Id to get | 
 **group** | [**Group**](Group.md)| Group to update | 

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateGroupMemberGrants**
> UpdateGroupMemberGrantsCommand updateGroupMemberGrants(groupId, userId, updateGroupMemberGrantsCommand)

Update a group member's category and tag assignment

Replaces one member's per-member category/tag grants. These narrow WITHIN the ceiling set by the member's group role — the two layers intersect — so every submitted id must be one the role already allows, or the request is rejected with 400. Duplicate ids within either array are also rejected with 400.  An empty array clears that resource's restriction. What the member then sees depends on their group role: when the role's requiresIndividualCategoryGrants (resp. requiresIndividualTagGrants) is false — the default — they fall back to the role's set; when it is true the assignment is mandatory, so clearing it leaves them seeing NOTHING for that resource. That is deliberate: it fails closed, so forgetting to assign a newly added member cannot silently widen their visibility.  Deliberately a dedicated endpoint rather than a field on the group-member upsert: it carries its own permission (group.members.grants.update), so a member who can manage the roster cannot thereby widen their own visibility.

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getGroupsApi();
final int groupId = 56; // int | Group the membership belongs to
final int userId = 56; // int | Member whose assignment is being replaced
final UpdateGroupMemberGrantsCommand updateGroupMemberGrantsCommand = ; // UpdateGroupMemberGrantsCommand | The category and tag ids to assign

try {
    final response = api.updateGroupMemberGrants(groupId, userId, updateGroupMemberGrantsCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling GroupsApi->updateGroupMemberGrants: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Group the membership belongs to | 
 **userId** | **int**| Member whose assignment is being replaced | 
 **updateGroupMemberGrantsCommand** | [**UpdateGroupMemberGrantsCommand**](UpdateGroupMemberGrantsCommand.md)| The category and tag ids to assign | 

### Return type

[**UpdateGroupMemberGrantsCommand**](UpdateGroupMemberGrantsCommand.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateGroupReceiptSettings**
> GroupReceiptSettings updateGroupReceiptSettings(groupId, updateGroupReceiptSettingsCommand)

Update group receipt settings

This will update the group receipt settings for a group

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getGroupsApi();
final int groupId = 56; // int | Group Id to update
final UpdateGroupReceiptSettingsCommand updateGroupReceiptSettingsCommand = ; // UpdateGroupReceiptSettingsCommand | Group settings to update

try {
    final response = api.updateGroupReceiptSettings(groupId, updateGroupReceiptSettingsCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling GroupsApi->updateGroupReceiptSettings: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Group Id to update | 
 **updateGroupReceiptSettingsCommand** | [**UpdateGroupReceiptSettingsCommand**](UpdateGroupReceiptSettingsCommand.md)| Group settings to update | 

### Return type

[**GroupReceiptSettings**](GroupReceiptSettings.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateGroupSettings**
> GroupSettings updateGroupSettings(groupId, updateGroupSettingsCommand)

Update group settings

This will update the group settings for a group

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getGroupsApi();
final int groupId = 56; // int | Group Id to update
final UpdateGroupSettingsCommand updateGroupSettingsCommand = ; // UpdateGroupSettingsCommand | Group settings to update

try {
    final response = api.updateGroupSettings(groupId, updateGroupSettingsCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling GroupsApi->updateGroupSettings: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Group Id to update | 
 **updateGroupSettingsCommand** | [**UpdateGroupSettingsCommand**](UpdateGroupSettingsCommand.md)| Group settings to update | 

### Return type

[**GroupSettings**](GroupSettings.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

