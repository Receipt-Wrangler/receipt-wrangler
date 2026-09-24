# openapi.model.Activity

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **int** |  | 
**type** | [**SystemTaskType**](SystemTaskType.md) |  | 
**status** | [**SystemTaskStatus**](SystemTaskStatus.md) |  | 
**startedAt** | **String** |  | 
**endedAt** | **String** |  | 
**ranByUserId** | **int** |  | [optional] 
**receiptId** | **int** |  | [optional] 
**groupId** | **int** |  | [optional] 
**canBeRestarted** | **bool** |  | [optional] 
**hasSourceFile** | **bool** | Whether the upload behind this activity is still on disk, so it can be previewed or downloaded. False once the temp-file retention window has passed, or for an activity that never had an upload. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


