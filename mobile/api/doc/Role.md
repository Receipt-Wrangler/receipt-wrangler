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
**categoryGrants** | **BuiltList&lt;int&gt;** | Category ids a GROUP role restricts its members to. Empty means unrestricted (members may use every category). Always empty for app roles. | [optional] 
**tagGrants** | **BuiltList&lt;int&gt;** | Tag ids a GROUP role restricts its members to. Empty means unrestricted (members may use every tag). Always empty for app roles. | [optional] 
**paidByUserGrants** | **BuiltList&lt;int&gt;** | User ids whose receipts a GROUP role lets its members see (by the receipt's \"paid by\" user). Empty with includeOwnPaidReceipts false means unrestricted (members see every payer's receipts). Always empty for app roles. | [optional] 
**includeOwnPaidReceipts** | **bool** | Whether a GROUP role lets each member see receipts they paid for. Part of the paid-by visibility filter; always false for app roles. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


