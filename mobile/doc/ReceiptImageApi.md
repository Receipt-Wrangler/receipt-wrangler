# openapi.api.ReceiptImageApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**convertToJpg**](ReceiptImageApi.md#converttojpg) | **POST** /receiptImage/convertToJpg | Converts a receipt image to jpg
[**deleteReceiptImageById**](ReceiptImageApi.md#deletereceiptimagebyid) | **DELETE** /receiptImage/{receiptImageId} | Delete receipt image
[**downloadReceiptImageById**](ReceiptImageApi.md#downloadreceiptimagebyid) | **GET** /receiptImage/{receiptImageId}/download | Download receipt image
[**getReceiptImageById**](ReceiptImageApi.md#getreceiptimagebyid) | **GET** /receiptImage/{receiptImageId} | Get receipt image
[**magicFillReceipt**](ReceiptImageApi.md#magicfillreceipt) | **POST** /receiptImage/magicFill | Reads a receipt image and returns the parsed results
[**uploadReceiptImage**](ReceiptImageApi.md#uploadreceiptimage) | **POST** /receiptImage/ | Uploads a receipt image


# **convertToJpg**
> EncodedImage convertToJpg(file)

Converts a receipt image to jpg

This will convert a receipt image to jpg, [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptImageApi();
final file = BINARY_DATA_HERE; // MultipartFile | Base64 encoded image

try {
    final result = api_instance.convertToJpg(file);
    print(result);
} catch (e) {
    print('Exception when calling ReceiptImageApi->convertToJpg: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **file** | **MultipartFile**| Base64 encoded image | 

### Return type

[**EncodedImage**](EncodedImage.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **deleteReceiptImageById**
> deleteReceiptImageById(receiptImageId)

Delete receipt image

This will delete a receipt image by id [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptImageApi();
final receiptImageId = 56; // int | Id of receipt image to get

try {
    api_instance.deleteReceiptImageById(receiptImageId);
} catch (e) {
    print('Exception when calling ReceiptImageApi->deleteReceiptImageById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receiptImageId** | **int**| Id of receipt image to get | 

### Return type

void (empty response body)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **downloadReceiptImageById**
> MultipartFile downloadReceiptImageById(receiptImageId)

Download receipt image

This will download a receipt image by id, [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptImageApi();
final receiptImageId = 56; // int | Id of receipt image to download

try {
    final result = api_instance.downloadReceiptImageById(receiptImageId);
    print(result);
} catch (e) {
    print('Exception when calling ReceiptImageApi->downloadReceiptImageById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receiptImageId** | **int**| Id of receipt image to download | 

### Return type

[**MultipartFile**](MultipartFile.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/octet-stream

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReceiptImageById**
> FileDataView getReceiptImageById(receiptImageId)

Get receipt image

This will get a receipt image by id, [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptImageApi();
final receiptImageId = 56; // int | Id of receipt image to get

try {
    final result = api_instance.getReceiptImageById(receiptImageId);
    print(result);
} catch (e) {
    print('Exception when calling ReceiptImageApi->getReceiptImageById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receiptImageId** | **int**| Id of receipt image to get | 

### Return type

[**FileDataView**](FileDataView.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **magicFillReceipt**
> Receipt magicFillReceipt(receiptImageId, file)

Reads a receipt image and returns the parsed results

This will parse and read a receipt image, [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptImageApi();
final receiptImageId = 56; // int | Id of receipt image to perform magic fill on
final file = BINARY_DATA_HERE; // MultipartFile | 

try {
    final result = api_instance.magicFillReceipt(receiptImageId, file);
    print(result);
} catch (e) {
    print('Exception when calling ReceiptImageApi->magicFillReceipt: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receiptImageId** | **int**| Id of receipt image to perform magic fill on | [optional] 
 **file** | **MultipartFile**|  | [optional] 

### Return type

[**Receipt**](Receipt.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **uploadReceiptImage**
> FileDataView uploadReceiptImage(file, receiptId, encodedImage)

Uploads a receipt image

This will upload a receipt image, [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = ReceiptImageApi();
final file = BINARY_DATA_HERE; // MultipartFile | 
final receiptId = 56; // int | Receipt foreign key
final encodedImage = encodedImage_example; // String | Base64 encoded image for file types that aren't viewable natively in the browser, such as PDFs

try {
    final result = api_instance.uploadReceiptImage(file, receiptId, encodedImage);
    print(result);
} catch (e) {
    print('Exception when calling ReceiptImageApi->uploadReceiptImage: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **file** | **MultipartFile**|  | 
 **receiptId** | **int**| Receipt foreign key | 
 **encodedImage** | **String**| Base64 encoded image for file types that aren't viewable natively in the browser, such as PDFs | [optional] 

### Return type

[**FileDataView**](FileDataView.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: multipart/form-data
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

