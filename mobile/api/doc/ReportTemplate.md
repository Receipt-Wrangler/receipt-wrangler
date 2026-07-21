# openapi.model.ReportTemplate

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
**name** | **String** | The template name (mirrors the saved report's name). | 
**configuration** | [**ReportRequestCommand**](ReportRequestCommand.md) |  | 
**configurationVersion** | **int** | Schema version the stored configuration was written under. | 
**allowedActions** | **BuiltList&lt;String&gt;** | The actions the requesting user may perform on this template (read, generate, update, delete, duplicate), resolved per user and populated only on the list response. Drives the row action buttons. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


