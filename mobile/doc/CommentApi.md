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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = CommentApi();
final upsertCommentCommand = UpsertCommentCommand(); // UpsertCommentCommand | Comment to create

try {
    final result = api_instance.addComment(upsertCommentCommand);
    print(result);
} catch (e) {
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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = CommentApi();
final commentId = 56; // int | Comment Id to delete

try {
    api_instance.deleteComment(commentId);
} catch (e) {
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
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

