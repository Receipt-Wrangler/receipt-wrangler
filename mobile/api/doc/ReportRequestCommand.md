# openapi.model.ReportRequestCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**name** | **String** | Report name; used for the download filename | [optional] 
**groupIds** | **BuiltList&lt;String&gt;** | Ids of the groups the report covers | 
**period** | [**ReportPeriod**](ReportPeriod.md) |  | 
**filter** | [**ReceiptPagedRequestFilter**](ReceiptPagedRequestFilter.md) | Which receipts go into the report | [optional] 
**groupBy** | **BuiltList&lt;String&gt;** | Ordered engine field keys to nest the report by | [optional] 
**groupByLabels** | **BuiltMap&lt;String, String&gt;** | Column-heading overrides for the grouping levels, keyed by the groupBy field key. A key that is absent, blank, or not present in groupBy falls back to the field catalog's label. | [optional] 
**detail** | [**ReportDetail**](ReportDetail.md) |  | 
**columns** | [**BuiltList&lt;ReportColumn&gt;**](ReportColumn.md) |  | 
**subtotals** | **bool** | Emit a subtotal row at each grouping level | [optional] 
**grandTotals** | **bool** | Emit one grand-total row across everything | [optional] 
**splitCategoriesEqually** | **bool** | Divide each receipt's money equally across its categories instead of attributing the whole amount to every category. Takes effect only when the report groups or aggregates by category. Omitted means false. | [optional] 
**splitTagsEqually** | **bool** | Divide each receipt's money equally across its tags instead of attributing the whole amount to every tag. Takes effect only when the report groups or aggregates by tag. Omitted means false. | [optional] 
**splitExcludedFields** | **BuiltList&lt;String&gt;** | Currency custom field keys (custom_<id>) that a split leaves whole. An entry naming a field that no longer exists or is not a currency field is ignored. | [optional] 
**document** | [**ReportDocument**](ReportDocument.md) |  | [optional] 
**formats** | **BuiltList&lt;String&gt;** | One or more output formats; multiple are zipped together | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


