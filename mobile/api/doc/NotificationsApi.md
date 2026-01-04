# openapi.api.NotificationsApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**deleteAllNotificationsForUser**](NotificationsApi.md#deleteallnotificationsforuser) | **DELETE** /notifications/ | Delete all notifications for user
[**deleteNotificationById**](NotificationsApi.md#deletenotificationbyid) | **DELETE** /notifications/{notificationId} | Delete notification by id
[**getNotificationCount**](NotificationsApi.md#getnotificationcount) | **GET** /notifications/notificationCount | Notification count
[**getNotificationsForuser**](NotificationsApi.md#getnotificationsforuser) | **GET** /notifications/ | Get all user notifications


# **deleteAllNotificationsForUser**
> deleteAllNotificationsForUser()

Delete all notifications for user

This deletes all notifications for a user

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getNotificationsApi();

try {
    api.deleteAllNotificationsForUser();
} catch on DioException (e) {
    print('Exception when calling NotificationsApi->deleteAllNotificationsForUser: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteNotificationById**
> deleteNotificationById(notificationId)

Delete notification by id

This deletes a notification by id

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getNotificationsApi();
final int notificationId = 56; // int | Notification Id to delete

try {
    api.deleteNotificationById(notificationId);
} catch on DioException (e) {
    print('Exception when calling NotificationsApi->deleteNotificationById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **notificationId** | **int**| Notification Id to delete | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getNotificationCount**
> int getNotificationCount()

Notification count

This will get the notification count for the currently logged in user

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getNotificationsApi();

try {
    final response = api.getNotificationCount();
    print(response);
} catch on DioException (e) {
    print('Exception when calling NotificationsApi->getNotificationCount: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

**int**

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getNotificationsForuser**
> BuiltList<Notification> getNotificationsForuser()

Get all user notifications

This will get all the notifications for the currently logged in user

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getNotificationsApi();

try {
    final response = api.getNotificationsForuser();
    print(response);
} catch on DioException (e) {
    print('Exception when calling NotificationsApi->getNotificationsForuser: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**BuiltList&lt;Notification&gt;**](Notification.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

