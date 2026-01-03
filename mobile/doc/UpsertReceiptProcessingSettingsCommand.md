# openapi.model.UpsertReceiptProcessingSettingsCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **String** | Name of the settings | 
**description** | **String** | Description of the settings | [optional] 
**aiType** | [**AiType**](AiType.md) |  | 
**url** | **String** | URL for custom endpoints | [optional] 
**key** | **String** | Key for endpoints that require authentication | [optional] 
**model** | **String** | LLM model | [optional] 
**isVisionModel** | **bool** | Is vision model | [optional] 
**ocrEngine** | [**OcrEngine**](OcrEngine.md) |  | 
**promptId** | **int** | Prompt foreign key | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


