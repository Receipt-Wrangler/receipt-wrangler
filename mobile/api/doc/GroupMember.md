# openapi.model.GroupMember

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**createdAt** | **String** |  | [optional] 
**groupId** | **int** | Group compound primary key | 
**groupRoleId** | **int** | Id of the modern group role assigned to the member | [optional] 
**updatedAt** | **String** |  | [optional] 
**userId** | **int** | User compound primary key | 
**categoryGrants** | **BuiltList&lt;int&gt;** | Category ids this individual member may see, narrowing WITHIN whatever their group role allows (the two layers intersect). Empty means the member adds no narrowing of their own: they fall back to the role's set, EXCEPT when the role sets requiresIndividualCategoryGrants, in which case an unassigned member sees no categories at all (fail closed). Read only here — write via PUT /group/{groupId}/member/{userId}/grants. | [optional] 
**tagGrants** | **BuiltList&lt;int&gt;** | Tag counterpart of categoryGrants. Restricted independently of categories. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


