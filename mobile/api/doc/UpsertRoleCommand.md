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
**seesAllMembers** | **bool** | Whether a GROUP role exempts its members from member-presence isolation (they see, and are seen by, every member of an isolated group). Only valid on group roles. | [optional] 
**skipDefaultGroupCreation** | **bool** | Whether users created with this APP role skip the automatic personal \"My Receipts\" group (the virtual \"All\" group is always created). Only valid on app roles; applies at user-creation time only. | [optional] 
**requiresIndividualCategoryGrants** | **bool** | Whether this GROUP role requires per-member category assignment. When true, a member with no individual category grants sees NO categories instead of falling back to the role's set, so an unassigned member fails closed. Only valid on group roles. | [optional] 
**requiresIndividualTagGrants** | **bool** | Tag counterpart of requiresIndividualCategoryGrants. Only valid on group roles. | [optional] 
**reportTemplateGrants** | [**BuiltList&lt;ReportTemplateGrant&gt;**](ReportTemplateGrant.md) | Per-template action grants for a GROUP role, restricting which report templates its members may act on. Only valid on group roles; omit or leave empty for unrestricted access. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


