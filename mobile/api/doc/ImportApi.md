# openapi.api.ImportApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**importConfigJson**](ImportApi.md#importconfigjson) | **POST** /import/importConfigJson | Import config json


# **importConfigJson**
> importConfigJson(file)

Import config json

This will import a config json

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getImportApi();
final MultipartFile file = BINARY_DATA_HERE; // MultipartFile | Files to quick scan

try {
    api.importConfigJson(file);
} catch on DioException (e) {
    print('Exception when calling ImportApi->importConfigJson: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **file** | **MultipartFile**| Files to quick scan | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

