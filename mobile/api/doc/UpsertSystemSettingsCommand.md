# openapi.model.UpsertSystemSettingsCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**enableLocalSignUp** | **bool** | Whether local sign up is enabled | [optional] 
**currencyDisplay** | **String** | Currency display | [optional] 
**currencyThousandthsSeparator** | [**CurrencySeparator**](CurrencySeparator.md) |  | 
**currencyDecimalSeparator** | [**CurrencySeparator**](CurrencySeparator.md) |  | 
**currencySymbolPosition** | [**CurrencySymbolPosition**](CurrencySymbolPosition.md) |  | 
**currencyHideDecimalPlaces** | **bool** | Whether to hide decimal places | 
**debugOcr** | **bool** |  | [optional] 
**numWorkers** | **int** | Number of workers to use | [optional] [default to 1]
**emailPollingInterval** | **int** | Email polling interval | [optional] 
**receiptProcessingSettingsId** | **int** | Receipt processing settings foreign key | [optional] 
**fallbackReceiptProcessingSettingsId** | **int** | Fallback receipt processing settings foreign key | [optional] 
**taskConcurrency** | **int** | Concurrency for task worker | 
**taskQueueConfigurations** | [**BuiltList&lt;UpsertTaskQueueConfiguration&gt;**](UpsertTaskQueueConfiguration.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


