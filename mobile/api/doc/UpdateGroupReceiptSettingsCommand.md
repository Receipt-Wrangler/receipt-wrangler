# openapi.model.UpdateGroupReceiptSettingsCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**hideImages** | **bool** | Hide receipt images | [optional] 
**hideReceiptCategories** | **bool** | Hide receipt categories | [optional] 
**hideReceiptTags** | **bool** | Hide receipt tags | [optional] 
**hideItemCategories** | **bool** | Hide receipt item categories | [optional] 
**hideItemTags** | **bool** | Hide receipt item tags | [optional] 
**hideComments** | **bool** | Hide receipt comments | [optional] 
**hideShareCategories** | **bool** | Hide share categories | [optional] 
**hideShareTags** | **bool** | Hide share tags | [optional] 
**quickScanPaidByEnabled** | **bool** | Show the paid by field in quick scan | [optional] 
**quickScanPaidByRequired** | **bool** | Require the paid by field in quick scan | [optional] 
**quickScanDefaultPaidByType** | [**QuickScanDefaultPaidByType**](QuickScanDefaultPaidByType.md) |  | [optional] 
**quickScanDefaultPaidById** | **int** | Default paid by user id when paid by is optional and type is USER | [optional] 
**quickScanStatusEnabled** | **bool** | Show the status field in quick scan | [optional] 
**quickScanStatusRequired** | **bool** | Require the status field in quick scan | [optional] 
**quickScanDefaultStatus** | [**ReceiptStatus**](ReceiptStatus.md) |  | [optional] 
**quickScanCategoriesEnabled** | **bool** | Show the categories field in quick scan | [optional] 
**quickScanCategoriesRequired** | **bool** | Require the categories field in quick scan | [optional] 
**quickScanTagsEnabled** | **bool** | Show the tags field in quick scan | [optional] 
**quickScanTagsRequired** | **bool** | Require the tags field in quick scan | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


