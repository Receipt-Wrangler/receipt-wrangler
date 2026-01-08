# openapi.api.SystemSettingsApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getSystemSettings**](SystemSettingsApi.md#getsystemsettings) | **GET** /systemSettings | Get system settings
[**restartTaskServer**](SystemSettingsApi.md#restarttaskserver) | **POST** /systemSettings/restartTaskServer | Restart task server
[**updateSystemSettings**](SystemSettingsApi.md#updatesystemsettings) | **PUT** /systemSettings | Update system settings


# **getSystemSettings**
> SystemSettings getSystemSettings()

Get system settings

This will get system settings

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getSystemSettingsApi();

try {
    final response = api.getSystemSettings();
    print(response);
} catch on DioException (e) {
    print('Exception when calling SystemSettingsApi->getSystemSettings: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**SystemSettings**](SystemSettings.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **restartTaskServer**
> restartTaskServer()

Restart task server

This will restart the task server

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getSystemSettingsApi();

try {
    api.restartTaskServer();
} catch on DioException (e) {
    print('Exception when calling SystemSettingsApi->restartTaskServer: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

void (empty response body)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateSystemSettings**
> SystemSettings updateSystemSettings(upsertSystemSettingsCommand)

Update system settings

This will update system settings

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getSystemSettingsApi();
final UpsertSystemSettingsCommand upsertSystemSettingsCommand = ; // UpsertSystemSettingsCommand | System settings to update

try {
    final response = api.updateSystemSettings(upsertSystemSettingsCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling SystemSettingsApi->updateSystemSettings: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **upsertSystemSettingsCommand** | [**UpsertSystemSettingsCommand**](UpsertSystemSettingsCommand.md)| System settings to update | 

### Return type

[**SystemSettings**](SystemSettings.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

