# openapi.api.TagApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createTag**](TagApi.md#createtag) | **POST** /tag/ | Create tag
[**deleteTag**](TagApi.md#deletetag) | **DELETE** /tag/{tagId} | Delete tag
[**getAllTags**](TagApi.md#getalltags) | **GET** /tag/ | Get all tags
[**getPagedTags**](TagApi.md#getpagedtags) | **POST** /tag/getPagedTags | Get paged tags
[**getTagCountByName**](TagApi.md#gettagcountbyname) | **GET** /tag/{tagName} | Get tag count by name
[**updateTag**](TagApi.md#updatetag) | **PUT** /tag/{tagId} | Update tag


# **createTag**
> Tag createTag(upsertTagCommand)

Create tag

This will create a tag

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getTagApi();
final UpsertTagCommand upsertTagCommand = ; // UpsertTagCommand | Tag to create

try {
    final response = api.createTag(upsertTagCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling TagApi->createTag: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **upsertTagCommand** | [**UpsertTagCommand**](UpsertTagCommand.md)| Tag to create | 

### Return type

[**Tag**](Tag.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteTag**
> deleteTag(tagId)

Delete tag

This will delete a tag by id

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getTagApi();
final int tagId = 56; // int | Id of tag to get

try {
    api.deleteTag(tagId);
} catch on DioException (e) {
    print('Exception when calling TagApi->deleteTag: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tagId** | **int**| Id of tag to get | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getAllTags**
> BuiltList<Tag> getAllTags()

Get all tags

This will return all tags in the system

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getTagApi();

try {
    final response = api.getAllTags();
    print(response);
} catch on DioException (e) {
    print('Exception when calling TagApi->getAllTags: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**BuiltList&lt;Tag&gt;**](Tag.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPagedTags**
> PagedData getPagedTags(pagedRequestCommand)

Get paged tags

This will return paged tags

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getTagApi();
final PagedRequestCommand pagedRequestCommand = ; // PagedRequestCommand | Paging and sorting data

try {
    final response = api.getPagedTags(pagedRequestCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling TagApi->getPagedTags: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **pagedRequestCommand** | [**PagedRequestCommand**](PagedRequestCommand.md)| Paging and sorting data | 

### Return type

[**PagedData**](PagedData.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getTagCountByName**
> int getTagCountByName(tagName)

Get tag count by name

This will count of names with the same name

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getTagApi();
final String tagName = tagName_example; // String | Tag name to get count of

try {
    final response = api.getTagCountByName(tagName);
    print(response);
} catch on DioException (e) {
    print('Exception when calling TagApi->getTagCountByName: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tagName** | **String**| Tag name to get count of | 

### Return type

**int**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateTag**
> Tag updateTag(tagId, upsertTagCommand)

Update tag

This will update a tag

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getTagApi();
final int tagId = 56; // int | Id of tag to get
final UpsertTagCommand upsertTagCommand = ; // UpsertTagCommand | Tag to update

try {
    final response = api.updateTag(tagId, upsertTagCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling TagApi->updateTag: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tagId** | **int**| Id of tag to get | 
 **upsertTagCommand** | [**UpsertTagCommand**](UpsertTagCommand.md)| Tag to update | 

### Return type

[**Tag**](Tag.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

