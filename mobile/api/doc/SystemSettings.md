# openapi.model.SystemSettings

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**id** | **int** |  | 
**createdAt** | **String** |  | 
**createdBy** | **int** |  | [optional] [default to 0]
**createdByString** | **String** | Created by entity's name | [optional] [default to '']
**updatedAt** | **String** |  | [optional] [default to '']
**enableLocalSignUp** | **bool** | Whether local sign up is enabled | [optional] [default to false]
**currencyDisplay** | **String** | Currency display | [optional] [default to '$']
**currencyThousandthsSeparator** | [**CurrencySeparator**](CurrencySeparator.md) |  | [optional] 
**currencyDecimalSeparator** | [**CurrencySeparator**](CurrencySeparator.md) |  | [optional] 
**currencySymbolPosition** | [**CurrencySymbolPosition**](CurrencySymbolPosition.md) |  | [optional] 
**currencyHideDecimalPlaces** | **bool** | Whether to hide decimal places | [optional] [default to false]
**debugOcr** | **bool** | Debug OCR | [optional] [default to false]
**numWorkers** | **int** | Number of workers to use | [optional] [default to 1]
**emailPollingInterval** | **int** | Email polling interval | [optional] [default to 1800]
**receiptProcessingSettingsId** | **int** | Receipt processing settings foreign key | [optional] 
**fallbackReceiptProcessingSettingsId** | **int** | Fallback receipt processing settings foreign key | [optional] 
**taskConcurrency** | **int** | Concurrency for task worker | [optional] [default to 10]
**pdfDpi** | **int** | DPI used when rasterizing PDFs for OCR/vision processing | [optional] [default to 300]
**taskQueueConfigurations** | [**BuiltList&lt;TaskQueueConfiguration&gt;**](TaskQueueConfiguration.md) |  | 
**mcpEnabled** | **bool** | Whether the OAuth 2.1-protected MCP server is enabled | [optional] [default to false]
**mcpPublicUrl** | **String** | Externally reachable origin used for MCP OAuth/metadata/redirect URLs and token audience | [optional] 
**showLoginQr** | **bool** | Whether to show the mobile-setup QR code on the desktop login page | [optional] [default to false]
**mobileServerUrl** | **String** | Server/API URL mobile clients connect to; encoded into the login QR's deep link | [optional] 
**refreshTokenValidForHours** | **int** | How long a refresh token stays valid, in hours. Refresh tokens rotate on every use, so this is how long a user can be away and still return signed in, not an absolute session cap. 1-720 (30 days); 0 means unset and falls back to the default. | [optional] [default to 24]
**mcpRefreshTokenValidForHours** | **int** | The same for MCP/OAuth connector refresh tokens, kept separate so a long window chosen for human convenience does not extend third-party client tokens. 1-720 (30 days); 0 means unset and falls back to the default. | [optional] [default to 24]

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


