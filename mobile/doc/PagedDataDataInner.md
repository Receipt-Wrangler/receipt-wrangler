# openapi.model.PagedDataDataInner

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**amount** | **String** | Receipt total amount | 
**categories** | [**List<Category>**](Category.md) | Categories associated to receipt | [optional] [default to const []]
**comments** | [**List<Comment>**](Comment.md) | Comments associated to receipt | [optional] [default to const []]
**createdAt** | **String** |  | [optional] 
**createdBy** | **int** |  | [optional] 
**date** | **String** | Receipt date | 
**groupId** | **int** | Group foreign key | 
**id** | **int** |  | 
**imageFiles** | [**List<FileData>**](FileData.md) | Files associated to receipt | [optional] [default to const []]
**name** | **String** | Tag name | 
**paidByUserId** | **int** | User paid foreign key | 
**receiptItems** | [**List<Item>**](Item.md) | Items associated to receipt | [optional] [default to const []]
**resolvedDate** | **String** | Date resolved | [optional] 
**status** | [**ReceiptStatus**](ReceiptStatus.md) |  | 
**tags** | [**List<Tag>**](Tag.md) | Tags associated to receipt | [optional] [default to const []]
**updatedAt** | **String** |  | [optional] 
**createdByString** | **String** | Created by string, which is anything that is not a user | [optional] 
**description** | **String** | Tag description | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


