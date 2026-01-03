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
**categories** | [**BuiltList&lt;UpsertCategoryCommand&gt;**](UpsertCategoryCommand.md) | Categories associated to receipt | [optional] 
**tags** | [**BuiltList&lt;UpsertTagCommand&gt;**](UpsertTagCommand.md) | Tags associated to receipt | [optional] 
**receiptItems** | [**BuiltList&lt;UpsertItemCommand&gt;**](UpsertItemCommand.md) | Items associated to receipt | [optional] 
**comments** | [**BuiltList&lt;UpsertCommentCommand&gt;**](UpsertCommentCommand.md) | Comments associated to receipt | [optional] 
**customFields** | [**BuiltList&lt;UpsertCustomFieldValueCommand&gt;**](UpsertCustomFieldValueCommand.md) | Custom fields associated to receipt | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


