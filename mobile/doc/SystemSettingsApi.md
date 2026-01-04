# openapi.api.SystemSettingsApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getSystemSettings**](SystemSettingsApi.md#getsystemsettings) | **GET** /systemSettings | Get system settings
[**updateSystemSettings**](SystemSettingsApi.md#updatesystemsettings) | **PUT** /systemSettings | Update system settings


# **getSystemSettings**
> SystemSettings getSystemSettings()

Get system settings

This will get system settings

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = SystemSettingsApi();

try {
    final result = api_instance.getSystemSettings();
    print(result);
} catch (e) {
    print('Exception when calling SystemSettingsApi->getSystemSettings: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**SystemSettings**](SystemSettings.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = SystemSettingsApi();
final upsertSystemSettingsCommand = UpsertSystemSettingsCommand(); // UpsertSystemSettingsCommand | System settings to update

try {
    final result = api_instance.updateSystemSettings(upsertSystemSettingsCommand);
    print(result);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

