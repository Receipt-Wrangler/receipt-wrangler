# openapi.model.ClaimListSummary

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**asOfDate** | **String** | Today in the app time zone (yyyy-MM-dd), the day daysWaiting is counted to | 
**claimCount** | **int** | Number of claims in the list | 
**awaitingCount** | **int** | Number of claims with awaitingSettlement true | 
**outstandingTotal** | **String** | Sum of outstandingAmount over the claims with awaitingSettlement true | 
**oldestDaysWaiting** | **int** | The largest daysWaiting in the list; absent when no claim has one | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


