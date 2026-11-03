# openapi.model.ReceiptHistoryEntry

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**systemTaskId** | **int** | The system task the entry was built from | 
**action** | **String** | What the entry did. One of CREATED, UPDATED, STATUS_CHANGED, RESOLVED, IMAGES_COMMENTS_UPDATED, FAILED. A plain string, not an enum, so a client must tolerate values it does not know. | 
**recordedAt** | **String** | When the entry was recorded (RFC 3339 instant) | 
**actor** | [**ReceiptHistoryActor**](ReceiptHistoryActor.md) |  | 
**hasSnapshots** | **bool** | Whether getReceiptHistorySnapshots has receipt copies for this entry | 
**errorMessage** | **String** | The error a FAILED entry recorded | [optional] 
**beforeState** | **String** | How far the compared \"before\" can be trusted, on update entries. One of COMPLETE, REBUILT (an older row compared against an earlier complete copy) or INCOMPLETE (an older row with no such copy). Plain string. | [optional] 
**beforeSourceType** | **String** | On a REBUILT entry, the system task type the \"before\" copy came from | [optional] 
**beforeSourceRecordedAt** | **String** | On a REBUILT entry, when the \"before\" copy was recorded (RFC 3339 instant) | [optional] 
**changes** | [**BuiltList&lt;ReceiptHistoryChange&gt;**](ReceiptHistoryChange.md) |  | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


