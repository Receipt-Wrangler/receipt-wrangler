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
[**updateGroupReceiptSettings**](GroupsApi.md#updategroupreceiptsettings) | **PUT** /group/{groupId}/groupReceiptSettings | Update group receipt settings
[**updateGroupSettings**](GroupsApi.md#updategroupsettings) | **PUT** /group/{groupId}/groupSettings | Update group settings


# **createGroup**
> createGroup(upsertGroupCommand)

Create group

This will create a group

### Example
```dart
import 'package:openapi/api.dart';

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

[bearerAuth](../README.md#bearerAuth)

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

[bearerAuth](../README.md#bearerAuth)

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

[bearerAuth](../README.md#bearerAuth)

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

final api = Openapi().getGroupsApi();
final int groupId = 56; // int | Group Id to get ocr text for

try {
    api.getOcrTextForGroup(groupId);
} catch on DioException (e) {
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
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPagedGroups**
> PagedData getPagedGroups(pagedGroupRequestCommand)

Get paged groups

This will return paged groups

### Example
```dart
import 'package:openapi/api.dart';

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

[bearerAuth](../README.md#bearerAuth)

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

[bearerAuth](../README.md#bearerAuth)

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

[bearerAuth](../README.md#bearerAuth)

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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

