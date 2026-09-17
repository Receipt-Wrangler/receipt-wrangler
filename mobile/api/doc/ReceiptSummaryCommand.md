# openapi.model.ReceiptSummaryCommand

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**filter** | [**ReceiptPagedRequestFilter**](ReceiptPagedRequestFilter.md) |  | [optional] 
**configurationGroupId** | **int** | The group whose receipt settings shape the breakdown. Exists for the synthetic \"All\" group, which spans every group the caller belongs to and so has no meaningful settings of its own - the client picks which member group's configuration to apply. OMIT it for a real group, where the group in the path is used. The DATA is always the path group's filtered set; this only chooses the shape of the breakdown. A group the caller cannot read is a 403. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


