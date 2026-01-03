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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = TagApi();
final upsertTagCommand = UpsertTagCommand(); // UpsertTagCommand | Tag to create

try {
    final result = api_instance.createTag(upsertTagCommand);
    print(result);
} catch (e) {
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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = TagApi();
final tagId = 56; // int | Id of tag to get

try {
    api_instance.deleteTag(tagId);
} catch (e) {
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
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getAllTags**
> List<Tag> getAllTags()

Get all tags

This will return all tags in the system

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = TagApi();

try {
    final result = api_instance.getAllTags();
    print(result);
} catch (e) {
    print('Exception when calling TagApi->getAllTags: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**List<Tag>**](Tag.md)

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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = TagApi();
final pagedRequestCommand = PagedRequestCommand(); // PagedRequestCommand | Paging and sorting data

try {
    final result = api_instance.getPagedTags(pagedRequestCommand);
    print(result);
} catch (e) {
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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = TagApi();
final tagName = tagName_example; // String | Tag name to get count of

try {
    final result = api_instance.getTagCountByName(tagName);
    print(result);
} catch (e) {
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
 - **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateTag**
> Tag updateTag(tagId, upsertTagCommand)

Update tag

This will update a tag

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = TagApi();
final tagId = 56; // int | Id of tag to get
final upsertTagCommand = UpsertTagCommand(); // UpsertTagCommand | Tag to update

try {
    final result = api_instance.updateTag(tagId, upsertTagCommand);
    print(result);
} catch (e) {
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

