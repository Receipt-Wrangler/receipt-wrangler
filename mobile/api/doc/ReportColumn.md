# openapi.model.ReportColumn

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**kind** | **String** |  | 
**name** | **String** | Machine identifier; formula expressions reference this | 
**label** | **String** |  | [optional] 
**field** | **String** | Field key the column displays (dimension columns) | [optional] 
**aggFunc** | **String** | Aggregate function (aggregate columns) | [optional] 
**measure** | **String** | Measure field key (aggregate columns, omitted for COUNT) | [optional] 
**expr** | **String** | Expression over other column names (formula columns) | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


