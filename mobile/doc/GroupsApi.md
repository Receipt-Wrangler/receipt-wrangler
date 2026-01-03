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
[**getOcrTextForGroup**](GroupsApi.md#getocrtextforgroup) | **GET** /group/{groupId}/ocrText | Reads each image in a group and returns the zipped read text
[**getPagedGroups**](GroupsApi.md#getpagedgroups) | **POST** /group/getPagedGroups | Get paged groups
[**pollGroupEmail**](GroupsApi.md#pollgroupemail) | **POST** /group/{groupId}/pollGroupEmail | Poll group email
[**updateGroup**](GroupsApi.md#updategroup) | **PUT** /group/{groupId} | Update a group
[**updateGroupSettings**](GroupsApi.md#updategroupsettings) | **PUT** /group/{groupId}/groupSettings | Update group settings


# **createGroup**
> createGroup(group)

Create group

This will create a group

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = GroupsApi();
final group = Group(); // Group | Group to create

try {
    api_instance.createGroup(group);
} catch (e) {
    print('Exception when calling GroupsApi->createGroup: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **group** | [**Group**](Group.md)| Group to create | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteGroup**
> deleteGroup(groupId)

Delete group

This will delete a group by id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = GroupsApi();
final groupId = 56; // int | Group Id to get

try {
    api_instance.deleteGroup(groupId);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getGroupById**
> getGroupById(groupId)

Gets a group by Id

This will get a group by Id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = GroupsApi();
final groupId = 56; // int | Group Id to get

try {
    api_instance.getGroupById(groupId);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getGroupsForuser**
> List<Group> getGroupsForuser()

Get groups for user

This will get groups for the currently logged in user

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = GroupsApi();

try {
    final result = api_instance.getGroupsForuser();
    print(result);
} catch (e) {
    print('Exception when calling GroupsApi->getGroupsForuser: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**List<Group>**](Group.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getOcrTextForGroup**
> getOcrTextForGroup(groupId)

Reads each image in a group and returns the zipped read text

This will get the ocr text, zipped, for each image in a group and one text file per image

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = GroupsApi();
final groupId = 56; // int | Group Id to get ocr text for

try {
    api_instance.getOcrTextForGroup(groupId);
} catch (e) {
    print('Exception when calling GroupsApi->getOcrTextForGroup: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Group Id to get ocr text for | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPagedGroups**
> PagedData getPagedGroups(pagedGroupRequestCommand)

Get paged groups

This will return paged groups

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = GroupsApi();
final pagedGroupRequestCommand = PagedGroupRequestCommand(); // PagedGroupRequestCommand | Paging and sorting data

try {
    final result = api_instance.getPagedGroups(pagedGroupRequestCommand);
    print(result);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = GroupsApi();
final groupId = 56; // int | Group Id to poll

try {
    api_instance.pollGroupEmail(groupId);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateGroup**
> updateGroup(groupId, group)

Update a group

This will update a group

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = GroupsApi();
final groupId = 56; // int | Group Id to get
final group = Group(); // Group | Group to update

try {
    api_instance.updateGroup(groupId, group);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateGroupSettings**
> GroupSettings updateGroupSettings(groupId, updateGroupSettingsCommand)

Update group settings

This will update the group settings for a group

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = GroupsApi();
final groupId = 56; // int | Group Id to update
final updateGroupSettingsCommand = UpdateGroupSettingsCommand(); // UpdateGroupSettingsCommand | Group settings to update

try {
    final result = api_instance.updateGroupSettings(groupId, updateGroupSettingsCommand);
    print(result);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

