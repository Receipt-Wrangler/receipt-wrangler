# openapi.model.ReceiptSummary

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**enabled** | **bool** | Whether the configuration group has the summary turned on. False comes back at a normal 200 with zeroed rows, so a client with stale group settings renders nothing rather than surfacing an error. | 
**configurationGroupId** | **int** | The group whose settings produced this breakdown | 
**overall** | [**ReceiptSummaryRow**](ReceiptSummaryRow.md) |  | 
**statuses** | [**BuiltList&lt;ReceiptSummaryRow&gt;**](ReceiptSummaryRow.md) | One row per configured status, in ReceiptStatus declaration order. A configured status matching no receipt is still present, with zeroed figures. Always present; empty when no status is configured. | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


