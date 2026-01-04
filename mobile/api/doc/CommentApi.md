# openapi.api.CommentApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**addComment**](CommentApi.md#addcomment) | **POST** /comment/ | Add comment
[**deleteComment**](CommentApi.md#deletecomment) | **DELETE** /comment/{commentId} | Delete comment


# **addComment**
> Comment addComment(upsertCommentCommand)

Add comment

This will add a comment to a receipt, [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getCommentApi();
final UpsertCommentCommand upsertCommentCommand = ; // UpsertCommentCommand | Comment to create

try {
    final response = api.addComment(upsertCommentCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling CommentApi->addComment: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **upsertCommentCommand** | [**UpsertCommentCommand**](UpsertCommentCommand.md)| Comment to create | 

### Return type

[**Comment**](Comment.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteComment**
> deleteComment(commentId)

Delete comment

This will delete a comment by id [SYSTEM User]

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getCommentApi();
final int commentId = 56; // int | Comment Id to delete

try {
    api.deleteComment(commentId);
} catch on DioException (e) {
    print('Exception when calling CommentApi->deleteComment: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **commentId** | **int**| Comment Id to delete | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

