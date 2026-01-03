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
**currencyThousandthsSeparator** | [**CurrencySeparator**](CurrencySeparator.md) |  | [optional] 
**currencyDecimalSeparator** | [**CurrencySeparator**](CurrencySeparator.md) |  | [optional] 
**currencySymbolPosition** | [**CurrencySymbolPosition**](CurrencySymbolPosition.md) |  | [optional] 
**currencyHideDecimalPlaces** | **bool** | Whether to hide decimal places | [optional] [default to false]
**debugOcr** | **bool** | Debug OCR | [optional] [default to false]
**numWorkers** | **int** | Number of workers to use | [optional] [default to 1]
**emailPollingInterval** | **int** | Email polling interval | [optional] [default to 1800]
**receiptProcessingSettingsId** | **int** | Receipt processing settings foreign key | [optional] 
**fallbackReceiptProcessingSettingsId** | **int** | Fallback receipt processing settings foreign key | [optional] 
**taskConcurrency** | **int** | Concurrency for task worker | [optional] [default to 10]
**taskQueueConfigurations** | [**BuiltList&lt;TaskQueueConfiguration&gt;**](TaskQueueConfiguration.md) |  | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


