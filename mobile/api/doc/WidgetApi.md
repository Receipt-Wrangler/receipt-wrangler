# openapi.api.WidgetApi

## Load the API package
```dart
import 'package:openapi/api.dart';
```

All URIs are relative to */api*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getPieChartData**](WidgetApi.md#getpiechartdata) | **POST** /widget/pieChart/{groupId} | Get pie chart data


# **getPieChartData**
> PieChartData getPieChartData(groupId, pieChartDataCommand)

Get pie chart data

This will get pie chart data for a group based on the specified grouping

### Example
```dart
import 'package:openapi/api.dart';
// TODO Configure API key authorization: apiKeyAuth
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKey = 'YOUR_API_KEY';
// uncomment below to setup prefix (e.g. Bearer) for API key, if needed
//defaultApiClient.getAuthentication<ApiKeyAuth>('apiKeyAuth').apiKeyPrefix = 'Bearer';

final api = Openapi().getWidgetApi();
final int groupId = 56; // int | Id of group to get pie chart data for
final PieChartDataCommand pieChartDataCommand = ; // PieChartDataCommand | Pie chart data request

try {
    final response = api.getPieChartData(groupId, pieChartDataCommand);
    print(response);
} catch on DioException (e) {
    print('Exception when calling WidgetApi->getPieChartData: $e\n');
}
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **groupId** | **int**| Id of group to get pie chart data for | 
 **pieChartDataCommand** | [**PieChartDataCommand**](PieChartDataCommand.md)| Pie chart data request | 

### Return type

[**PieChartData**](PieChartData.md)

### Authorization

[apiKeyAuth](../README.md#apiKeyAuth), [bearerAuth](../README.md#bearerAuth)

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to Model list]](../README.md#documentation-for-models) [[Back to README]](../README.md)

