# openapi.model.UpsertClaimLedgerEntryCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**type** | **String** | One of PAYMENT, RECOUPMENT, NOTE; anything else is a 400 keyed type | 
**date** | **String** | Calendar day (yyyy-MM-dd) | 
**amount** | **String** | Decimal string, greater than zero, at most two decimal places. Required for PAYMENT and RECOUPMENT; optional for NOTE. | [optional] 
**note** | **String** | At most 2000 characters. Required (non-blank) for NOTE. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


