# openapi.model.Claims

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**userId** | **int** | User foreign key | [default to 0]
**userRole** | [**UserRole**](UserRole.md) | User's role | [default to 'USER']
**displayName** | **String** | Display name | [default to '']
**defaultAvatarColor** | **String** | Default avatar color | [default to '']
**username** | **String** | User's username used to login | [default to '']
**iss** | **String** | Issuer | [default to '']
**sub** | **String** | Subject | [optional] [default to '']
**aud** | **BuiltList&lt;String&gt;** | Audience | [optional] [default to ListBuilder()]
**exp** | **int** | Expiration time | [default to 0]
**nbf** | **int** | Not before | [optional] [default to 0]
**iat** | **int** | Issued at | [optional] [default to 0]
**jti** | **String** | JWT ID | [optional] [default to '']

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


