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
**document** | [**ReportDocument**](ReportDocument.md) |  | [optional] 
**formats** | **BuiltList&lt;String&gt;** | One or more output formats; multiple are zipped together | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


