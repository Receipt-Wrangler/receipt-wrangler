# openapi.model.ReportPreviewResponse

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**html** | **String** | The rendered report preview as a self-contained HTML document. | 
**receiptCount** | **int** | The number of receipts the current configuration covers. | 
**allowedActions** | **BuiltList&lt;String&gt;** | The actions the requesting user may perform on the template (read, generate, update, delete, duplicate), resolved per user and populated only by the dashboard report-widget render endpoint. Drives the widget's download button. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


