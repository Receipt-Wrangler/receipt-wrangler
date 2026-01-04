# openapi.model.GroupSettings

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **int** | Group settings id | 
**groupId** | **int** | Group foreign key | 
**emailIntegrationEnabled** | **bool** | Whether email integration is enabled | [optional] 
**systemEmailId** | **int** | System email foreign key | [optional] 
**systemEmail** | [**SystemEmail**](SystemEmail.md) |  | [optional] 
**emailToRead** | **String** | Email to read | [optional] 
**subjectLineRegexes** | [**List<SubjectLineRegex>**](SubjectLineRegex.md) | Subject line regexes | [optional] [default to const []]
**emailWhiteList** | [**List<GroupSettingsWhiteListEmail>**](GroupSettingsWhiteListEmail.md) | Email white list | [optional] [default to const []]
**emailDefaultReceiptStatus** | [**ReceiptStatus**](ReceiptStatus.md) |  | [optional] 
**emailDefaultReceiptPaidById** | **int** | User foreign key | [optional] 
**prompt** | [**Prompt**](Prompt.md) |  | [optional] 
**promptId** | **int** | Prompt foreign key | [optional] 
**fallbackPrompt** | [**Prompt**](Prompt.md) |  | [optional] 
**fallbackPromptId** | **int** | Fallback prompt foreign key | [optional] 
**createdAt** | **String** |  | [optional] 
**createdBy** | **int** |  | [optional] 
**updatedAt** | **String** |  | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


