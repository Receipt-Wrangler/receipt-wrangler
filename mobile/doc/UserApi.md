# openapi.api.UserApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**convertDummyUserById**](UserApi.md#convertdummyuserbyid) | **POST** /user/{userId}/convertDummyUserToNormalUser | Converts dummy user
[**createUser**](UserApi.md#createuser) | **POST** /user | Create user
[**deleteUserById**](UserApi.md#deleteuserbyid) | **DELETE** /user/{userId} | Delete user
[**getAmountOwedForUser**](UserApi.md#getamountowedforuser) | **GET** /user/amountOwedForUser | Get amount owed for user
[**getAppData**](UserApi.md#getappdata) | **GET** /user/appData | Get app data
[**getUserClaims**](UserApi.md#getuserclaims) | **GET** /user/getUserClaims | Get claims for logged in user
[**getUsernameCount**](UserApi.md#getusernamecount) | **GET** /user/{username} | Get username count
[**getUsers**](UserApi.md#getusers) | **GET** /user | Get users
[**resetPasswordById**](UserApi.md#resetpasswordbyid) | **POST** /user/{userId}/resetPassword | Reset password
[**updateUserById**](UserApi.md#updateuserbyid) | **PUT** /user/{userId} | Update user by id
[**updateUserProfile**](UserApi.md#updateuserprofile) | **PUT** /user/updateUserProfile | Update user profile


# **convertDummyUserById**
> convertDummyUserById(userId, resetPasswordCommand)

Converts dummy user

This will convert a dummy user to a normal system user, [SYSTEM ADMIN]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = UserApi();
final userId = 56; // int | Id of user to convert to normal system user
final resetPasswordCommand = ResetPasswordCommand(); // ResetPasswordCommand | Login credentials for new user

try {
    api_instance.convertDummyUserById(userId, resetPasswordCommand);
} catch (e) {
    print('Exception when calling UserApi->convertDummyUserById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **int**| Id of user to convert to normal system user | 
 **resetPasswordCommand** | [**ResetPasswordCommand**](ResetPasswordCommand.md)| Login credentials for new user | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **createUser**
> createUser(user)

Create user

This will to create a user, [SYSTEM ADMIN]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = UserApi();
final user = User(); // User | User to create

try {
    api_instance.createUser(user);
} catch (e) {
    print('Exception when calling UserApi->createUser: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **user** | [**User**](User.md)| User to create | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteUserById**
> deleteUserById(userId)

Delete user

This will delete a system user by id [SYSTEM ADMIN]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = UserApi();
final userId = 56; // int | Id of user to update

try {
    api_instance.deleteUserById(userId);
} catch (e) {
    print('Exception when calling UserApi->deleteUserById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **int**| Id of user to update | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getAmountOwedForUser**
> Map<String, String> getAmountOwedForUser(groupId, receiptIds)

Get amount owed for user

This will return the amount owed for the logged in user, in the specified group, [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = UserApi();
final groupId = 56; // int | The Id of the group to get amount owed for
final receiptIds = []; // List<int> | The Id of the receipts to get amount owed for

try {
    final result = api_instance.getAmountOwedForUser(groupId, receiptIds);
    print(result);
} catch (e) {
    print('Exception when calling UserApi->getAmountOwedForUser: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| The Id of the group to get amount owed for | [optional] 
 **receiptIds** | [**List<int>**](int.md)| The Id of the receipts to get amount owed for | [optional] [default to const []]

### Return type

**Map<String, String>**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getAppData**
> AppData getAppData()

Get app data

This will return the user's app data for the currently logged in user [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = UserApi();

try {
    final result = api_instance.getAppData();
    print(result);
} catch (e) {
    print('Exception when calling UserApi->getAppData: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**AppData**](AppData.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getUserClaims**
> Claims getUserClaims()

Get claims for logged in user

This will return the user's token claims for the currently logged in user [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = UserApi();

try {
    final result = api_instance.getUserClaims();
    print(result);
} catch (e) {
    print('Exception when calling UserApi->getUserClaims: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**Claims**](Claims.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getUsernameCount**
> int getUsernameCount(username)

Get username count

This will return the number of users in the system with the same username

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = UserApi();
final username = username_example; // String | Username to get the count of

try {
    final result = api_instance.getUsernameCount(username);
    print(result);
} catch (e) {
    print('Exception when calling UserApi->getUsernameCount: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **String**| Username to get the count of | 

### Return type

**int**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getUsers**
> List<UserView> getUsers()

Get users

This will get all the users in the system and return a view without sensative information

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = UserApi();

try {
    final result = api_instance.getUsers();
    print(result);
} catch (e) {
    print('Exception when calling UserApi->getUsers: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**List<UserView>**](UserView.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **resetPasswordById**
> resetPasswordById(userId, resetPasswordCommand)

Reset password

This will reset a password for a user, [SYSTEM ADMIN]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = UserApi();
final userId = 56; // int | Id of user to reset password
final resetPasswordCommand = ResetPasswordCommand(); // ResetPasswordCommand | Login credentials for new user

try {
    api_instance.resetPasswordById(userId, resetPasswordCommand);
} catch (e) {
    print('Exception when calling UserApi->resetPasswordById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **int**| Id of user to reset password | 
 **resetPasswordCommand** | [**ResetPasswordCommand**](ResetPasswordCommand.md)| Login credentials for new user | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateUserById**
> updateUserById(userId, user)

Update user by id

This will update a user by id, [SYSTEM ADMIN]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = UserApi();
final userId = 56; // int | Id of user to update
final user = User(); // User | User to update

try {
    api_instance.updateUserById(userId, user);
} catch (e) {
    print('Exception when calling UserApi->updateUserById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **userId** | **int**| Id of user to update | 
 **user** | [**User**](User.md)| User to update | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **updateUserProfile**
> updateUserProfile(updateProfileCommand)

Update user profile

This will update the logged in user's user profile

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = UserApi();
final updateProfileCommand = UpdateProfileCommand(); // UpdateProfileCommand | User profile to update

try {
    api_instance.updateUserProfile(updateProfileCommand);
} catch (e) {
    print('Exception when calling UserApi->updateUserProfile: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **updateProfileCommand** | [**UpdateProfileCommand**](UpdateProfileCommand.md)| User profile to update | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

