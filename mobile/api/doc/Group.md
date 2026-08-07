# openapi.model.Group

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**createdAt** | **String** |  | [optional] 
**createdBy** | **int** |  | [optional] 
**groupSettings** | [**GroupSettings**](GroupSettings.md) |  | [optional] 
**groupReceiptSettings** | [**GroupReceiptSettings**](GroupReceiptSettings.md) |  | 
**groupMembers** | [**BuiltList&lt;GroupMember&gt;**](GroupMember.md) | Members of the group | 
**id** | **int** |  | 
**isDefault** | **bool** | Is default group (not used yet) | [optional] 
**name** | **String** | Name of the group | 
**isAllGroup** | **bool** | Is all group for user | 
**isolateMembers** | **bool** | Whether member-presence isolation is enabled for the group. When on, members cannot discover other members unless they hold a group role flagged seesAllMembers. Defaults to false. | [optional] 
**status** | [**GroupStatus**](GroupStatus.md) |  | 
**updatedAt** | **String** |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


