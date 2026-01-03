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
// TODO Configure HTTP Bearer authorization: bearerAuth
// Case 1. Use String Token
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken('YOUR_ACCESS_TOKEN');
// Case 2. Use Function which generate token.
// String yourTokenGeneratorFunction() { ... }
//defaultApiClient.getAuthentication<HttpBearerAuth>('bearerAuth').setAccessToken(yourTokenGeneratorFunction);

final api_instance = FeatureConfigApi();

try {
    final result = api_instance.getFeatureConfig();
    print(result);
} catch (e) {
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

