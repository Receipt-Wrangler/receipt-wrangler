# openapi.model.AppData

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**about** | [**About**](About.md) |  | 
**claims** | [**Claims**](Claims.md) |  | 
**groups** | [**BuiltList&lt;Group&gt;**](Group.md) | Groups in the system | 
**users** | [**BuiltList&lt;UserView&gt;**](UserView.md) | Users in the system | 
**userPreferences** | [**UserPreferences**](UserPreferences.md) |  | 
**featureConfig** | [**FeatureConfig**](FeatureConfig.md) |  | 
**categories** | [**BuiltList&lt;Category&gt;**](Category.md) | Categories in the system | 
**tags** | [**BuiltList&lt;Tag&gt;**](Tag.md) | Tags in the system | 
**jwt** | **String** | JWT token | [optional] 
**refreshToken** | **String** | Refresh token | [optional] 
**currencyDisplay** | **String** | Currency display | 
**currencyThousandthsSeparator** | [**CurrencySeparator**](CurrencySeparator.md) |  | [optional] 
**currencyDecimalSeparator** | [**CurrencySeparator**](CurrencySeparator.md) |  | [optional] 
**currencySymbolPosition** | [**CurrencySymbolPosition**](CurrencySymbolPosition.md) |  | [optional] 
**currencyHideDecimalPlaces** | **bool** | Whether to hide decimal places | [optional] 
**icons** | [**BuiltList&lt;Icon&gt;**](Icon.md) | Icons in the system | 
**appPermissions** | [**BuiltList&lt;Permission&gt;**](Permission.md) | The calling user's effective app-level permissions. | 
**groupPermissions** | [**BuiltMap&lt;String, BuiltList&lt;Permission&gt;&gt;**](BuiltList.md) | The calling user's effective group-level permissions, keyed by group id. | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


