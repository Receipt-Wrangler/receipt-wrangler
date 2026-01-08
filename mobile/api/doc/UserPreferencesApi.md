# openapi.api.UserPreferencesApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getUserPreferences**](UserPreferencesApi.md#getuserpreferences) | **GET** /userPreferences | Get user preferences
[**updateUserPreferences**](UserPreferencesApi.md#updateuserpreferences) | **PUT** /userPreferences | Update user preferences


# **getUserPreferences**
> UserPreferences getUserPreferences()

Get user preferences

This will return the user's preferences for the currently logged in user [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getUserPreferencesApi();

try {
    final response = api.getUserPreferences();
    print(response);
} catch on DioException (e) {
    print('Exception when calling UserPreferencesApi->getUserPreferences: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**UserPreferences**](UserPreferences.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateUserPreferences**
> UserPreferences updateUserPreferences(userPreferences)

Update user preferences

This will update the user's preferences for the currently logged in user [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getUserPreferencesApi();
final UserPreferences userPreferences = ; // UserPreferences | User preferences to update

try {
    final response = api.updateUserPreferences(userPreferences);
    print(response);
} catch on DioException (e) {
    print('Exception when calling UserPreferencesApi->updateUserPreferences: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userPreferences** | [**UserPreferences**](UserPreferences.md)| User preferences to update | 

### Return type

[**UserPreferences**](UserPreferences.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

