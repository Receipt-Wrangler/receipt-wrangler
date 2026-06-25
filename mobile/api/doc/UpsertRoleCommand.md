# openapi.model.UpsertRoleCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **String** |  | 
**description** | **String** |  | [optional] 
**scope** | [**PermissionScope**](PermissionScope.md) |  | 
**permissions** | [**BuiltList&lt;Permission&gt;**](Permission.md) |  | 
**categoryGrants** | **BuiltList&lt;int&gt;** | Category ids to restrict a GROUP role's members to. Only valid on group roles; omit or leave empty for unrestricted access. | [optional] 
**tagGrants** | **BuiltList&lt;int&gt;** | Tag ids to restrict a GROUP role's members to. Only valid on group roles; omit or leave empty for unrestricted access. | [optional] 
**paidByUserGrants** | **BuiltList&lt;int&gt;** | User ids whose receipts a GROUP role's members may see (by the receipt's \"paid by\" user). Only valid on group roles; omit or leave empty (with includeOwnPaidReceipts false) for unrestricted access. | [optional] 
**includeOwnPaidReceipts** | **bool** | Whether to also let each member see receipts they paid for. Only valid on group roles. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


