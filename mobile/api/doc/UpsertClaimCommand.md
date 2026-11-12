# openapi.model.UpsertClaimCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **String** | Required, non-blank after trimming, at most 255 characters | 
**referenceNumber** | **String** | At most 100 characters | [optional] 
**type** | **String** | Free text, at most 100 characters | [optional] 
**status** | **String** | One of DRAFT, SUBMITTED, PARTIALLY_PAID, PAID, DISPUTED, CLOSED; anything else is a 400 keyed status | 
**submittedDate** | **String** | Calendar day (yyyy-MM-dd) | [optional] 
**dueDate** | **String** | Calendar day (yyyy-MM-dd) | [optional] 
**claimedAmount** | **String** | Decimal string, zero or more, at most two decimal places | [optional] 
**notes** | **String** | At most 10000 characters | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


