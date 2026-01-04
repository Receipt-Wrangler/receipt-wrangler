# openapi.model.SystemSettings

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
**enableLocalSignUp** | **bool** | Whether local sign up is enabled | [optional] [default to false]
**currencyDisplay** | **String** | Currency display | [optional] [default to '$']
**debugOcr** | **bool** | Debug OCR | [optional] [default to false]
**numWorkers** | **int** | Number of workers to use | [optional] [default to 1]
**emailPollingInterval** | **int** | Email polling interval | [optional] [default to 1800]
**receiptProcessingSettingsId** | **int** | Receipt processing settings foreign key | [optional] 
**fallbackReceiptProcessingSettingsId** | **int** | Fallback receipt processing settings foreign key | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


