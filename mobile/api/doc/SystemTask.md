# openapi.model.SystemTask

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **int** |  | 
**createdAt** | **String** |  | 
**createdBy** | **int** |  | [optional] [default to 0]
**createdByString** | **String** | Created by entity's name | [optional] [default to '']
**updatedAt** | **String** |  | [optional] [default to '']
**type** | [**SystemTaskType**](SystemTaskType.md) |  | [optional] 
**status** | [**SystemTaskStatus**](SystemTaskStatus.md) |  | [optional] 
**startedAt** | **String** |  | [optional] 
**endedAt** | **String** |  | [optional] 
**associatedEntityId** | **int** |  | [optional] 
**associatedEntityType** | [**AssociatedEntityType**](AssociatedEntityType.md) |  | [optional] 
**ranByUserId** | **int** |  | [optional] 
**receiptId** | **int** |  | [optional] 
**groupId** | **int** |  | [optional] 
**resultDescription** | **String** |  | [optional] 
**apiKeyId** | **String** |  | [optional] 
**childSystemTasks** | [**BuiltList&lt;SystemTask&gt;**](SystemTask.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


