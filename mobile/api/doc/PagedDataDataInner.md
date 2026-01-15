# openapi.model.PagedDataDataInner

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**amount** | **String** | Receipt total amount | 
**categories** | [**BuiltList&lt;Category&gt;**](Category.md) | Categories associated to receipt | 
**comments** | [**BuiltList&lt;Comment&gt;**](Comment.md) | Comments associated to receipt | 
**customFields** | [**BuiltList&lt;CustomFieldValue&gt;**](CustomFieldValue.md) | Custom fields associated to receipt | 
**createdAt** | **String** |  | 
**createdBy** | **int** |  | [optional] [default to 0]
**date** | **String** | Receipt date | 
**groupId** | **int** |  | 
**id** | **int** |  | 
**imageFiles** | [**BuiltList&lt;FileData&gt;**](FileData.md) | Files associated to receipt | [optional] 
**name** | **String** | Custom Field name | 
**paidByUserId** | **int** | User paid foreign key | 
**receiptItems** | [**BuiltList&lt;Item&gt;**](Item.md) | Items associated to receipt | 
**resolvedDate** | **String** | Date resolved | [optional] 
**status** | [**SystemTaskStatus**](SystemTaskStatus.md) |  | 
**tags** | [**BuiltList&lt;Tag&gt;**](Tag.md) | Tags associated to receipt | 
**updatedAt** | **String** |  | [optional] [default to '']
**createdByString** | **String** | Created by entity's name | [optional] [default to '']
**description** | **String** | Custom Field description | [optional] 
**prompt** | [**Prompt**](Prompt.md) |  | 
**groupSettings** | [**GroupSettings**](GroupSettings.md) |  | [optional] 
**groupReceiptSettings** | [**GroupReceiptSettings**](GroupReceiptSettings.md) |  | 
**groupMembers** | [**BuiltList&lt;GroupMember&gt;**](GroupMember.md) | Members of the group | 
**isDefault** | **bool** | Is default group (not used yet) | [optional] 
**isAllGroup** | **bool** | Is all group for user | 
**numberOfReceipts** | **int** | Number of receipts associated with this tag | 
**type** | [**CustomFieldType**](CustomFieldType.md) |  | 
**startedAt** | **String** |  | 
**endedAt** | **String** |  | 
**associatedEntityId** | **int** |  | [optional] 
**associatedEntityType** | [**AssociatedEntityType**](AssociatedEntityType.md) |  | [optional] 
**ranByUserId** | **int** |  | [optional] 
**receiptId** | **int** |  | [optional] 
**resultDescription** | **String** |  | [optional] 
**apiKeyId** | **String** |  | [optional] 
**childSystemTasks** | [**BuiltList&lt;SystemTask&gt;**](SystemTask.md) |  | [optional] 
**aiType** | [**AiType**](AiType.md) |  | [optional] 
**url** | **String** | URL for custom endpoints | [optional] 
**key** | **String** | Key for endpoints that require authentication | [optional] 
**model** | **String** | LLM model | [optional] 
**isVisionModel** | **bool** | Is vision model | [optional] 
**ocrEngine** | [**OcrEngine**](OcrEngine.md) |  | [optional] 
**promptId** | **int** | Prompt foreign key | [optional] 
**host** | **String** | IMAP host | [optional] 
**port** | **String** | IMAP port | [optional] 
**username** | **String** | IMAP username | [optional] 
**password** | **String** | IMAP password | [optional] 
**useStartTLS** | **bool** | Whether to use STARTTLS | [optional] 
**canBeRestarted** | **bool** |  | [optional] 
**options** | [**BuiltList&lt;CustomFieldOption&gt;**](CustomFieldOption.md) |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


