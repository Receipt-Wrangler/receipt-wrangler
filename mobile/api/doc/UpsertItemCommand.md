# openapi.model.UpsertItemCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**amount** | **String** | Amount the item costs | 
**chargedToUserId** | **int** | User foreign key | 
**name** | **String** | Item name | 
**receiptId** | **int** | Receipt foreign key | 
**status** | [**ItemStatus**](ItemStatus.md) |  | 
**categories** | [**BuiltList&lt;UpsertCategoryCommand&gt;**](UpsertCategoryCommand.md) | Categories associated to item | [optional] 
**tags** | [**BuiltList&lt;UpsertTagCommand&gt;**](UpsertTagCommand.md) | Tags associated to item | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


