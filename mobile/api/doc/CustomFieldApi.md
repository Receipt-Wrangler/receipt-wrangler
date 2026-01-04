# openapi.api.CustomFieldApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**createCustomField**](CustomFieldApi.md#createcustomfield) | **POST** /customField/ | Create custom field
[**deleteCustomField**](CustomFieldApi.md#deletecustomfield) | **DELETE** /customField/{customFieldId} | Delete custom field
[**getCustomFieldById**](CustomFieldApi.md#getcustomfieldbyid) | **GET** /customField/{customFieldId} | Get custom field
[**getPagedCustomFields**](CustomFieldApi.md#getpagedcustomfields) | **POST** /customField/getPagedCustomFields | Get paged custom fields


# **createCustomField**
> CustomField createCustomField(upsertCustomFieldCommand)

Create custom field

This will create a custom field

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getCustomFieldApi();
final UpsertCustomFieldCommand upsertCustomFieldCommand = ; // UpsertCustomFieldCommand | Custom field to create

try {
    final response = api.createCustomField(upsertCustomFieldCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling CustomFieldApi->createCustomField: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **upsertCustomFieldCommand** | [**UpsertCustomFieldCommand**](UpsertCustomFieldCommand.md)| Custom field to create | 

### Return type

[**CustomField**](CustomField.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteCustomField**
> deleteCustomField(customFieldId)

Delete custom field

This will delete a custom field by id

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getCustomFieldApi();
final int customFieldId = 56; // int | Custom field Id to get

try {
    api.deleteCustomField(customFieldId);
} catch on DioException (e) {
    print('Exception when calling CustomFieldApi->deleteCustomField: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **customFieldId** | **int**| Custom field Id to get | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getCustomFieldById**
> CustomField getCustomFieldById(customFieldId)

Get custom field

This will get a custom field by id

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getCustomFieldApi();
final int customFieldId = 56; // int | Custom field Id to get

try {
    final response = api.getCustomFieldById(customFieldId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling CustomFieldApi->getCustomFieldById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **customFieldId** | **int**| Custom field Id to get | 

### Return type

[**CustomField**](CustomField.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPagedCustomFields**
> PagedData getPagedCustomFields(pagedRequestCommand)

Get paged custom fields

This will return paged custom fields

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getCustomFieldApi();
final PagedRequestCommand pagedRequestCommand = ; // PagedRequestCommand | Paging and sorting data

try {
    final response = api.getPagedCustomFields(pagedRequestCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling CustomFieldApi->getPagedCustomFields: $e\n');
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

