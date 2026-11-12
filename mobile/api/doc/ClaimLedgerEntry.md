# openapi.model.ClaimLedgerEntry

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **int** |  | 
**claimId** | **int** |  | 
**type** | **String** | One of PAYMENT (money received; adds to paidAmount), RECOUPMENT (money taken back; subtracts from paidAmount) or NOTE (no effect on the figures). A plain string, not an enum; an entry of a type the server does not recognize counts toward no figure. | 
**date** | **String** | Calendar day the event happened (yyyy-MM-dd) | 
**amount** | **String** | Decimal string, always positive; the type gives the direction. Present on every PAYMENT and RECOUPMENT; optional, and never counted, on a NOTE. | [optional] 
**note** | **String** | Free text. Absent when not set. | [optional] 
**createdAt** | **String** | RFC 3339 instant the entry was recorded | 
**updatedAt** | **String** | RFC 3339 instant the entry was last changed | 
**createdBy** | **int** | The user who recorded the entry | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


