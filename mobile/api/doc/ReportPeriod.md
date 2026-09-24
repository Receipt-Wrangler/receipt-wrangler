# openapi.model.ReportPeriod

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**preset** | **String** |  | 
**startDate** | **String** | Start date (YYYY-MM-DD), read only when preset is custom | [optional] 
**endDate** | **String** | End date (YYYY-MM-DD), read only when preset is custom | [optional] 
**dateField** | **String** | Which receipt date the period covers, as a ReceiptPagedRequestFilter date key: date, resolvedDate or createdAt. Omitted means date. A plain string rather than an enum on purpose: this rides inside ReportTemplate.configuration, and a value added to a closed enum would fail that whole payload on already-released mobile builds. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


