# openapi.model.ReceiptHistoryChange

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**field** | **String** | Which field changed. One of NAME, AMOUNT, DATE, GROUP, PAID_BY, STATUS, RESOLVED_DATE, CATEGORIES, TAGS, CUSTOM_FIELD, ITEM, SHARE, ITEM_CATEGORIES, ITEM_TAGS, COMMENT, IMAGE, EMAIL_INBOX, EMAIL_SENDER. Plain string. | 
**kind** | **String** | One of CHANGED, ADDED, REMOVED. Plain string. | 
**label** | **String** | What the change is about, where the field alone does not say: the item's name (ITEM, ITEM_CATEGORIES, ITEM_TAGS) or the custom field's name (CUSTOM_FIELD) | [optional] 
**userId** | **int** | The share's charged-to user (SHARE) or the comment's author (COMMENT) | [optional] 
**from** | [**ReceiptHistoryValue**](ReceiptHistoryValue.md) |  | [optional] 
**to** | [**ReceiptHistoryValue**](ReceiptHistoryValue.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


