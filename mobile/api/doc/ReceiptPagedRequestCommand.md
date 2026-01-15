# openapi.model.ReceiptPagedRequestCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**page** | **int** | Page number | 
**pageSize** | **int** | Number of records per page | 
**orderBy** | **String** | field to order on | [optional] 
**sortDirection** | [**SortDirection**](SortDirection.md) |  | [optional] 
**filter** | [**ReceiptPagedRequestFilter**](ReceiptPagedRequestFilter.md) |  | [optional] 
**fullReceipts** | **bool** | Whether to include all receipt associations (receiptItems, comments, customFields, imageFiles, etc.) | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


