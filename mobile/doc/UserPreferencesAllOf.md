# openapi.model.UserPreferencesAllOf

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **int** | User preferences id | 
**userId** | **int** | User foreign key | 
**quickScanDefaultGroupId** | **int** | Group foreign key | [optional] [default to 0]
**quickScanDefaultPaidById** | **int** | User foreign key | [optional] [default to 0]
**quickScanDefaultStatus** | [**ReceiptStatus**](ReceiptStatus.md) |  | [optional] 
**showLargeImagePreviews** | **bool** | Whether to show large image previews | [optional] [default to false]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


