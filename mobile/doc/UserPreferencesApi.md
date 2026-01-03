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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = UserPreferencesApi();

try {
    final result = api_instance.getUserPreferences();
    print(result);
} catch (e) {
    print('Exception when calling UserPreferencesApi->getUserPreferences: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**UserPreferences**](UserPreferences.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = UserPreferencesApi();
final userPreferences = UserPreferences(); // UserPreferences | User preferences to update

try {
    final result = api_instance.updateUserPreferences(userPreferences);
    print(result);
} catch (e) {
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

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

