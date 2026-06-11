# openapi.model.Role

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **int** |  | 
**name** | **String** |  | 
**description** | **String** |  | [optional] 
**scope** | [**PermissionScope**](PermissionScope.md) |  | 
**isDefault** | **bool** | Whether this role is the default for its scope — assigned to new accounts (APP) or to group creators (GROUP). Exactly one role per scope is the default. | 
**isSystem** | **bool** |  | 
**permissions** | [**BuiltList&lt;Permission&gt;**](Permission.md) |  | 
**assignedCount** | **int** | Number of users or group members currently assigned this role | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


