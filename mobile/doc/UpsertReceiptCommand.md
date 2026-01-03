# openapi.model.UpsertReceiptCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **String** | Receipt name | 
**amount** | **String** | Receipt total amount | 
**date** | **String** | Receipt date | 
**groupId** | **int** | Group foreign key | 
**paidByUserId** | **int** | User paid foreign key | 
**status** | [**ReceiptStatus**](ReceiptStatus.md) |  | 
**categories** | [**List<UpsertCategoryCommand>**](UpsertCategoryCommand.md) | Categories associated to receipt | [optional] [default to const []]
**tags** | [**List<UpsertTagCommand>**](UpsertTagCommand.md) | Tags associated to receipt | [optional] [default to const []]
**receiptItems** | [**List<UpsertItemCommand>**](UpsertItemCommand.md) | Items associated to receipt | [optional] [default to const []]
**comments** | [**List<UpsertCommentCommand>**](UpsertCommentCommand.md) | Comments associated to receipt | [optional] [default to const []]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


