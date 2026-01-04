# openapi.api.CategoryApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createCategory**](CategoryApi.md#createcategory) | **POST** /category/ | Create category
[**deleteCategory**](CategoryApi.md#deletecategory) | **DELETE** /category/{categoryId} | Delete category
[**getAllCategories**](CategoryApi.md#getallcategories) | **GET** /category/ | Get all categories
[**getCategoryCountByName**](CategoryApi.md#getcategorycountbyname) | **GET** /category/{categoryName} | Get category count by name
[**getPagedCategories**](CategoryApi.md#getpagedcategories) | **POST** /category/getPagedCategories | Get paged categories
[**updateCategory**](CategoryApi.md#updatecategory) | **PUT** /category/{categoryId} | Update category


# **createCategory**
> Category createCategory(category)

Create category

This will create a category

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = CategoryApi();
final category = Category(); // Category | Category to create

try {
    final result = api_instance.createCategory(category);
    print(result);
} catch (e) {
    print('Exception when calling CategoryApi->createCategory: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **category** | [**Category**](Category.md)| Category to create | 

### Return type

[**Category**](Category.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteCategory**
> deleteCategory(categoryId)

Delete category

This will delete a category by id

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = CategoryApi();
final categoryId = 56; // int | Category Id to get

try {
    api_instance.deleteCategory(categoryId);
} catch (e) {
    print('Exception when calling CategoryApi->deleteCategory: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **categoryId** | **int**| Category Id to get | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getAllCategories**
> List<Category> getAllCategories()

Get all categories

This will return all categories in the system

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = CategoryApi();

try {
    final result = api_instance.getAllCategories();
    print(result);
} catch (e) {
    print('Exception when calling CategoryApi->getAllCategories: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**List<Category>**](Category.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getCategoryCountByName**
> int getCategoryCountByName(categoryName)

Get category count by name

This will return a count of categories with the same name

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = CategoryApi();
final categoryName = categoryName_example; // String | Category name to get count of

try {
    final result = api_instance.getCategoryCountByName(categoryName);
    print(result);
} catch (e) {
    print('Exception when calling CategoryApi->getCategoryCountByName: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **categoryName** | **String**| Category name to get count of | 

### Return type

**int**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: text/plain

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPagedCategories**
> PagedData getPagedCategories(pagedRequestCommand)

Get paged categories

This will return paged categories

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = CategoryApi();
final pagedRequestCommand = PagedRequestCommand(); // PagedRequestCommand | Paging and sorting data

try {
    final result = api_instance.getPagedCategories(pagedRequestCommand);
    print(result);
} catch (e) {
    print('Exception when calling CategoryApi->getPagedCategories: $e\n');
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

# **updateCategory**
> updateCategory(categoryId, category)

Update category

This will update a category

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = CategoryApi();
final categoryId = 56; // int | Category Id to get
final category = Category(); // Category | Category to update

try {
    api_instance.updateCategory(categoryId, category);
} catch (e) {
    print('Exception when calling CategoryApi->updateCategory: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **categoryId** | **int**| Category Id to get | 
 **category** | [**Category**](Category.md)| Category to update | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

