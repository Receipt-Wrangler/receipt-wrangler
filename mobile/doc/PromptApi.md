# openapi.api.PromptApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createDefaultPrompt**](PromptApi.md#createdefaultprompt) | **POST** /prompt/createDefaultPrompt | Create default prompt
[**createPrompt**](PromptApi.md#createprompt) | **POST** /prompt/ | Create prompt
[**deletePromptById**](PromptApi.md#deletepromptbyid) | **DELETE** /prompt/{id} | Delete prompt by id
[**getPagedPrompts**](PromptApi.md#getpagedprompts) | **POST** /prompt/getPagedPrompts | Gets paged prompts
[**getPromptById**](PromptApi.md#getpromptbyid) | **GET** /prompt/{id} | Get prompt by id
[**updatePromptById**](PromptApi.md#updatepromptbyid) | **PUT** /prompt/{id} | Update prompt by id


# **createDefaultPrompt**
> Prompt createDefaultPrompt()

Create default prompt

This will create a default prompt

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = PromptApi();

try {
    final result = api_instance.createDefaultPrompt();
    print(result);
} catch (e) {
    print('Exception when calling PromptApi->createDefaultPrompt: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**Prompt**](Prompt.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createPrompt**
> Prompt createPrompt(upsertPromptCommand)

Create prompt

This will create a prompt

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = PromptApi();
final upsertPromptCommand = UpsertPromptCommand(); // UpsertPromptCommand | Prompt to create

try {
    final result = api_instance.createPrompt(upsertPromptCommand);
    print(result);
} catch (e) {
    print('Exception when calling PromptApi->createPrompt: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **upsertPromptCommand** | [**UpsertPromptCommand**](UpsertPromptCommand.md)| Prompt to create | 

### Return type

[**Prompt**](Prompt.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deletePromptById**
> deletePromptById(id)

Delete prompt by id

This will delete a prompt by id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = PromptApi();
final id = 56; // int | Id of prompt to delete

try {
    api_instance.deletePromptById(id);
} catch (e) {
    print('Exception when calling PromptApi->deletePromptById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of prompt to delete | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPagedPrompts**
> PagedData getPagedPrompts(pagedRequestCommand)

Gets paged prompts

This will return paged prompts

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = PromptApi();
final pagedRequestCommand = PagedRequestCommand(); // PagedRequestCommand | Paging and sorting data

try {
    final result = api_instance.getPagedPrompts(pagedRequestCommand);
    print(result);
} catch (e) {
    print('Exception when calling PromptApi->getPagedPrompts: $e\n');
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

# **getPromptById**
> Prompt getPromptById(id)

Get prompt by id

This will get a prompt by id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = PromptApi();
final id = 56; // int | Id of prompt to get

try {
    final result = api_instance.getPromptById(id);
    print(result);
} catch (e) {
    print('Exception when calling PromptApi->getPromptById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of prompt to get | 

### Return type

[**Prompt**](Prompt.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updatePromptById**
> Prompt updatePromptById(id, upsertPromptCommand)

Update prompt by id

This will update a prompt by id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = PromptApi();
final id = 56; // int | Id of prompt to update
final upsertPromptCommand = UpsertPromptCommand(); // UpsertPromptCommand | Prompt to update

try {
    final result = api_instance.updatePromptById(id, upsertPromptCommand);
    print(result);
} catch (e) {
    print('Exception when calling PromptApi->updatePromptById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of prompt to update | 
 **upsertPromptCommand** | [**UpsertPromptCommand**](UpsertPromptCommand.md)| Prompt to update | 

### Return type

[**Prompt**](Prompt.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

