# openapi.api.FeatureConfigApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getFeatureConfig**](FeatureConfigApi.md#getfeatureconfig) | **GET** /featureConfig | Get feature config


# **getFeatureConfig**
> FeatureConfig getFeatureConfig()

Get feature config

This will get the server's feature config

### Example
```dart
import 'package:openapi/api.dart';

final api = Openapi().getFeatureConfigApi();

try {
    final response = api.getFeatureConfig();
    print(response);
} catch on DioException (e) {
    print('Exception when calling FeatureConfigApi->getFeatureConfig: $e\n');
}
```

### Parameters
This endpoint does not need any parameter.

### Return type

[**FeatureConfig**](FeatureConfig.md)

### Authorization

[bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: Not defined
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

