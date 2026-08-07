# openapi.model.UpsertGroupCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**groupMembers** | [**BuiltList&lt;UpsertGroupMemberCommand&gt;**](UpsertGroupMemberCommand.md) | Members of the group | 
**isDefault** | **bool** | Is default group (not used yet) | [optional] 
**name** | **String** | Name of the group | 
**isAllGroup** | **bool** | Is all group for user | [optional] 
**isolateMembers** | **bool** | Whether to enable member-presence isolation for the group. When on, members cannot discover other members unless they hold a group role flagged seesAllMembers. Defaults to false. | [optional] 
**status** | [**GroupStatus**](GroupStatus.md) |  | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


