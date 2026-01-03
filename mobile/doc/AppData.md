# openapi.model.AppData

## Load the model package
```dart
import 'package:openapi/api.dart';
```

## Properties
Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**claims** | [**Claims**](Claims.md) |  | 
**groups** | [**List<Group>**](Group.md) | Groups in the system | [default to const []]
**users** | [**List<UserView>**](UserView.md) | Users in the system | [default to const []]
**userPreferences** | [**UserPreferences**](UserPreferences.md) |  | 
**featureConfig** | [**FeatureConfig**](FeatureConfig.md) |  | 
**categories** | [**List<Category>**](Category.md) | Categories in the system | [default to const []]
**tags** | [**List<Tag>**](Tag.md) | Tags in the system | [default to const []]
**jwt** | **String** | JWT token | [optional] 
**refreshToken** | **String** | Refresh token | [optional] 
**currencyDisplay** | **String** | Currency display | 

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


