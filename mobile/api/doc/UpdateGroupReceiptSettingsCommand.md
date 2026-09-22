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
**quickScanCommentEnabled** | **bool** | Show the comment field in quick scan | [optional] 
**quickScanCommentRequired** | **bool** | Require the comment field in quick scan | [optional] 
**defaultCustomFieldIds** | **BuiltList&lt;int&gt;** | Custom field ids to pre-add to every receipt created for this group. OMIT the key to leave the configured set unchanged (clients that hide this section, e.g. for a user without app.custom-fields.read, must omit it); send an empty array to clear it. Requires app.custom-fields.read - a caller without it gets a 403. | [optional] 
**applyDefaultCustomFieldsOnIngest** | **bool** | Also attach the group's default custom fields to receipts the SERVER creates (quick scan, email integration). OMIT the key to leave the stored value unchanged. | [optional] 
**receiptSummaryEnabled** | **bool** | Show the block of totals under this group's receipts table. OMIT the key to leave the stored value unchanged - a client that does not render this section must omit it rather than send false, or it switches the summary off for the whole group. | [optional] 
**receiptSummaryCustomFieldIds** | **BuiltList&lt;int&gt;** | CURRENCY custom field ids to total in the receipt summary. OMIT the key to leave the configured set unchanged (clients that hide this control, e.g. for a user without app.custom-fields.read, must omit it); send an empty array to clear it. Requires app.custom-fields.read - a caller without it gets a 403. Every id must be an existing CURRENCY custom field - anything else is a 400, because only a currency value can be summed. | [optional] 
**receiptSummaryStatuses** | [**BuiltList&lt;ReceiptStatus&gt;**](ReceiptStatus.md) | Receipt statuses to break out as their own row in the receipt summary. OMIT the key to leave the configured set unchanged; send an empty array to clear it. Unlike the custom field ids this needs NO extra permission - gating it would lock an admin without app.custom-fields.read out of the feature entirely. An unrecognized status is a 400. | [optional] 
**receiptSummaryPosition** | [**ReceiptSummaryPosition**](ReceiptSummaryPosition.md) | Where the receipt summary renders relative to the receipts list. OMIT the key to leave the stored value unchanged - a client that does not render this control must omit it rather than send a zero value, or it moves the block for the whole group. Needs no extra permission, for the same reason the statuses do not. Anything but TOP or BOTTOM is a 400. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


