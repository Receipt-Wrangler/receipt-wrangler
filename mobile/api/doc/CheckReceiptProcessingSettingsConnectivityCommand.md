# openapi.model.CheckReceiptProcessingSettingsConnectivityCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **int** | Receipt processing settings id | [optional] 
**name** | **String** | Name of the settings | [optional] 
**aiType** | [**AiType**](AiType.md) |  | [optional] 
**url** | **String** | URL for custom endpoints | [optional] 
**key** | **String** | Key for endpoints that require authentication | [optional] 
**model** | **String** | LLM model | [optional] 
**enforceJsonResponseFormat** | **bool** | Enforce JSON response format on the LLM provider. Disable if the provider does not support this flag. | [optional] 
**numWorkers** | **int** | Number of workers to use | [optional] 
**ocrEngine** | [**OcrEngine**](OcrEngine.md) |  | [optional] 
**promptId** | **int** | Prompt foreign key | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


