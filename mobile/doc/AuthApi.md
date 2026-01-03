# openapi.api.AuthApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getNewRefreshToken**](AuthApi.md#getnewrefreshtoken) | **POST** /token/ | Get fresh tokens
[**login**](AuthApi.md#login) | **POST** /login/ | Login
[**logout**](AuthApi.md#logout) | **POST** /logout/ | Logout
[**signUp**](AuthApi.md#signup) | **POST** /signUp | Signs up


# **getNewRefreshToken**
> GetNewRefreshToken200Response getNewRefreshToken(logoutCommand)

Get fresh tokens

This will get a fresh token pair for the user

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = AuthApi();
final logoutCommand = LogoutCommand(); // LogoutCommand | Refresh token body for clients that don't use cookies

try {
    final result = api_instance.getNewRefreshToken(logoutCommand);
    print(result);
} catch (e) {
    print('Exception when calling AuthApi->getNewRefreshToken: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **logoutCommand** | [**LogoutCommand**](LogoutCommand.md)| Refresh token body for clients that don't use cookies | [optional] 

### Return type

[**GetNewRefreshToken200Response**](GetNewRefreshToken200Response.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **login**
> AppData login(loginCommand)

Login

This will log a user into the system

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = AuthApi();
final loginCommand = LoginCommand(); // LoginCommand | Login data

try {
    final result = api_instance.login(loginCommand);
    print(result);
} catch (e) {
    print('Exception when calling AuthApi->login: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **loginCommand** | [**LoginCommand**](LoginCommand.md)| Login data | 

### Return type

[**AppData**](AppData.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **logout**
> logout(logoutCommand)

Logout

This will log a user out of the system and revoke their token [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = AuthApi();
final logoutCommand = LogoutCommand(); // LogoutCommand | Refresh token

try {
    api_instance.logout(logoutCommand);
} catch (e) {
    print('Exception when calling AuthApi->logout: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **logoutCommand** | [**LogoutCommand**](LogoutCommand.md)| Refresh token | [optional] 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **signUp**
> signUp(signUpCommand)

Signs up

This will sign a user up for the system

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = AuthApi();
final signUpCommand = SignUpCommand(); // SignUpCommand | Sign up data

try {
    api_instance.signUp(signUpCommand);
} catch (e) {
    print('Exception when calling AuthApi->signUp: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **signUpCommand** | [**SignUpCommand**](SignUpCommand.md)| Sign up data | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

