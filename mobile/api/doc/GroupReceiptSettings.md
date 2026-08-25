# openapi.model.GroupReceiptSettings

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
**groupId** | **int** | Group foreign key | 
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
**quickScanCommentEnabled** | **bool** | Show the comment field in quick scan | [optional] 
**quickScanCommentRequired** | **bool** | Require the comment field in quick scan | [optional] 
**defaultCustomFieldIds** | **BuiltList&lt;int&gt;** | Custom field ids that are pre-added to every receipt created for this group. Always present; an empty array means the group has configured none. Read only here - write via UpdateGroupReceiptSettingsCommand.defaultCustomFieldIds. | [optional] 
**applyDefaultCustomFieldsOnIngest** | **bool** | Also attach the group's default custom fields to receipts the SERVER creates (quick scan, email integration). Off by default. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


