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

final api = Openapi().getReceiptImageApi();
final MultipartFile file = BINARY_DATA_HERE; // MultipartFile | Base64 encoded image

try {
    final response = api.convertToJpg(file);
    print(response);
} catch on DioException (e) {
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

final api = Openapi().getReceiptImageApi();
final int receiptImageId = 56; // int | Id of receipt image to get

try {
    api.deleteReceiptImageById(receiptImageId);
} catch on DioException (e) {
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
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **downloadReceiptImageById**
> Uint8List downloadReceiptImageById(receiptImageId)

Download receipt image

This will download a receipt image by id, [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getReceiptImageApi();
final int receiptImageId = 56; // int | Id of receipt image to download

try {
    final response = api.downloadReceiptImageById(receiptImageId);
    print(response);
} catch on DioException (e) {
    print('Exception when calling ReceiptImageApi->downloadReceiptImageById: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **receiptImageId** | **int**| Id of receipt image to download | 

### Return type

[**Uint8List**](Uint8List.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/octet-stream, application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

# **getReceiptImageById**
> FileDataView getReceiptImageById(receiptImageId)

Get receipt image

This will get a receipt image by id, [SYSTEM USER]

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getReceiptImageApi();
final int receiptImageId = 56; // int | Id of receipt image to get

try {
    final response = api.getReceiptImageById(receiptImageId);
    print(response);
} catch on DioException (e) {
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

final api = Openapi().getReceiptImageApi();
final int receiptImageId = 56; // int | Id of receipt image to perform magic fill on
final MultipartFile file = BINARY_DATA_HERE; // MultipartFile | 

try {
    final response = api.magicFillReceipt(receiptImageId, file);
    print(response);
} catch on DioException (e) {
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

final api = Openapi().getReceiptImageApi();
final MultipartFile file = BINARY_DATA_HERE; // MultipartFile | 
final int receiptId = 56; // int | Receipt foreign key
final String encodedImage = encodedImage_example; // String | Base64 encoded image for file types that aren't viewable natively in the browser, such as PDFs

try {
    final response = api.uploadReceiptImage(file, receiptId, encodedImage);
    print(response);
} catch on DioException (e) {
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

