# openapi.model.Item

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**isTaxed** | **bool** | Is taxed (not used) | [optional] 
**amount** | **String** | Amount the item costs | 
**chargedToUserId** | **int** | User foreign key | 
**createdAt** | **String** |  | [optional] 
**createdBy** | **int** |  | [optional] 
**id** | **int** |  | [optional] 
**name** | **String** | Item name | 
**receiptId** | **int** | Receipt foreign key | 
**status** | [**ItemStatus**](ItemStatus.md) |  | 
**categories** | [**BuiltList&lt;Category&gt;**](Category.md) | Categories associated to the item | [optional] 
**tags** | [**BuiltList&lt;Tag&gt;**](Tag.md) | Tags associated to the item | [optional] 
**updatedAt** | **String** |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


