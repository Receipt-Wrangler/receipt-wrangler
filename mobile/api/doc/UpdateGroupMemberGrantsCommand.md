# openapi.model.UpdateGroupMemberGrantsCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**categoryGrants** | **BuiltList&lt;int&gt;** | Category ids to assign to this member. Every id must sit within the ceiling set by the member's group role, or the request is rejected with 400; duplicate ids are rejected with 400 as well. An empty array clears the member's category restriction — handing them back to their role's set, or, when the role sets requiresIndividualCategoryGrants, leaving them with no categories at all (fail closed). | [optional] 
**tagGrants** | **BuiltList&lt;int&gt;** | Tag counterpart of categoryGrants. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


