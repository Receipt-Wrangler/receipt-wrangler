# openapi.api.SystemEmailApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**checkSystemEmailConnectivity**](SystemEmailApi.md#checksystememailconnectivity) | **POST** /systemEmail/checkConnectivity | Check system email connectivity
[**createSystemEmail**](SystemEmailApi.md#createsystememail) | **POST** /systemEmail/ | Create system email
[**deleteSystemEmailById**](SystemEmailApi.md#deletesystememailbyid) | **DELETE** /systemEmail/{id} | Delete system email by id
[**getPagedSystemEmails**](SystemEmailApi.md#getpagedsystememails) | **POST** /systemEmail/getSystemEmails | Gets paged system emails
[**getSystemEmailById**](SystemEmailApi.md#getsystememailbyid) | **GET** /systemEmail/{id} | Get system email by id
[**updateSystemEmailById**](SystemEmailApi.md#updatesystememailbyid) | **PUT** /systemEmail/{id} | Update system email by id


# **checkSystemEmailConnectivity**
> SystemTask checkSystemEmailConnectivity(checkEmailConnectivityCommand)

Check system email connectivity

This will check system email connectivity

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getSystemEmailApi();
final CheckEmailConnectivityCommand checkEmailConnectivityCommand = ; // CheckEmailConnectivityCommand | System email to check connectivity

try {
    final response = api.checkSystemEmailConnectivity(checkEmailConnectivityCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling SystemEmailApi->checkSystemEmailConnectivity: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **checkEmailConnectivityCommand** | [**CheckEmailConnectivityCommand**](CheckEmailConnectivityCommand.md)| System email to check connectivity | 

### Return type

[**SystemTask**](SystemTask.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createSystemEmail**
> SystemEmail createSystemEmail(upsertSystemEmailCommand)

Create system email

This will create a system email

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getSystemEmailApi();
final UpsertSystemEmailCommand upsertSystemEmailCommand = ; // UpsertSystemEmailCommand | System email to create

try {
    final response = api.createSystemEmail(upsertSystemEmailCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling SystemEmailApi->createSystemEmail: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **upsertSystemEmailCommand** | [**UpsertSystemEmailCommand**](UpsertSystemEmailCommand.md)| System email to create | 

### Return type

[**SystemEmail**](SystemEmail.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteSystemEmailById**
> deleteSystemEmailById(id)

Delete system email by id

This will delete a system email by id

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getSystemEmailApi();
final int id = 56; // int | Id of system email to delete

try {
    api.deleteSystemEmailById(id);
} catch on DioException (e) {
    print('Exception when calling SystemEmailApi->deleteSystemEmailById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of system email to delete | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getPagedSystemEmails**
> PagedData getPagedSystemEmails(pagedRequestCommand)

Gets paged system emails

This will return paged and sorted system emails

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getSystemEmailApi();
final PagedRequestCommand pagedRequestCommand = ; // PagedRequestCommand | Paging and sorting data

try {
    final response = api.getPagedSystemEmails(pagedRequestCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling SystemEmailApi->getPagedSystemEmails: $e\n');
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

# **getSystemEmailById**
> SystemEmail getSystemEmailById(id)

Get system email by id

This will get a system email by id

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getSystemEmailApi();
final int id = 56; // int | Id of system email to get

try {
    final response = api.getSystemEmailById(id);
    print(response);
} catch on DioException (e) {
    print('Exception when calling SystemEmailApi->getSystemEmailById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of system email to get | 

### Return type

[**SystemEmail**](SystemEmail.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateSystemEmailById**
> SystemEmail updateSystemEmailById(id, updatePassword, upsertSystemEmailCommand)

Update system email by id

This will update a system email by id

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getSystemEmailApi();
final int id = 56; // int | Id of system email to update
final bool updatePassword = true; // bool | Whether or not to update the password
final UpsertSystemEmailCommand upsertSystemEmailCommand = ; // UpsertSystemEmailCommand | System email to update

try {
    final response = api.updateSystemEmailById(id, updatePassword, upsertSystemEmailCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling SystemEmailApi->updateSystemEmailById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **int**| Id of system email to update | 
 **updatePassword** | **bool**| Whether or not to update the password | 
 **upsertSystemEmailCommand** | [**UpsertSystemEmailCommand**](UpsertSystemEmailCommand.md)| System email to update | 

### Return type

[**SystemEmail**](SystemEmail.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

