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
**appPermissions** | **BuiltList&lt;String&gt;** | The calling user's effective app-level permissions. Deliberately typed as plain strings rather than the Permission enum: this is server-resolved data, not a contract. A granted entry may be a wildcard (e.g. \"app.*\"), which is not an enum member, and a client built before a newly added permission must still be able to parse the payload. Clients match these with the wildcard matcher. | 
**groupPermissions** | [**BuiltMap&lt;String, BuiltList&lt;String&gt;&gt;**](BuiltList.md) | The calling user's effective group-level permissions, keyed by group id. Plain strings for the same reason as appPermissions. | 
**groupCategories** | [**BuiltMap&lt;String, BuiltList&lt;Category&gt;&gt;**](BuiltList.md) | The categories the calling user may use in each group, keyed by group id. Filtered to the user's group-role grants (the full pool when unrestricted). Non-admins receive categories only through this map. | [optional] 
**groupTags** | [**BuiltMap&lt;String, BuiltList&lt;Tag&gt;&gt;**](BuiltList.md) | The tags the calling user may use in each group, keyed by group id. Filtered to the user's group-role grants (the full pool when unrestricted). Non-admins receive tags only through this map. | [optional] 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


