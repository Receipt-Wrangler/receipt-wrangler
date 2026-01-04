import 'package:dart_jsonwebtoken/dart_jsonwebtoken.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/category_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/tag_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/models/user_preferences_model.dart';

import '../models/system_settings_model.dart';

Future<bool> refreshTokens(
    AuthModel authModel,
    GroupModel groupModel,
    UserModel userModel,
    UserPreferencesModel userPreferencesModel,
    CategoryModel categoryModel,
    TagModel tagModel,
    SystemSettingsModel systemSettingsModel,
    {bool force = false}) async {
  var jwt = await authModel.getJwt();
  var refreshToken = await authModel.getRefreshToken();
  var isAuthenticated = false;

  if (force == true) {
    try {
      await getAndSetTokens(authModel);
      return true;
    } catch (e) {
      print(e);
      print("failed");
      authModel.purgeTokens();
      return false;
    }
  }

  if (!force) {
    // If token is valid, then continue on
    if (isTokenValid(jwt)) {
      isAuthenticated = true;
    } else {
      // If token is invalid, but refresh token is valid, then get a new token pair
      if (isTokenValid(refreshToken)) {
        try {
          await getAndSetTokens(authModel);
          isAuthenticated = true;
        } catch (e) {
          // If the refresh fails, redirect to redirect path and consider it a failure
          authModel.purgeTokens();
          isAuthenticated = false;
        }
      } else {
        // purge old tokens
        authModel.purgeTokens();
        isAuthenticated = false;
      }
    }
  }

  // If user is authenticated, but data does not exist yet
  if (isAuthenticated && groupModel.groups.isEmpty) {
    try {
      var appDataResponse =
          await OpenApiClient.client.getUserApi().getAppData();
      await storeAppData(
          authModel,
          groupModel,
          userModel,
          userPreferencesModel,
          categoryModel,
          tagModel,
          systemSettingsModel,
          appDataResponse.data as AppData);
    } catch (e) {
      print(e);
      print("failed to set token");
      isAuthenticated = false;
    }
  }

  return isAuthenticated;
}

Future<void> getAndSetTokens(AuthModel authModel) async {
  var refreshToken = await authModel.getRefreshToken() ?? "";
  var logoutCommand =
      (LogoutCommandBuilder()..refreshToken = refreshToken).build();
  var tokenPairResponse = await OpenApiClient.client
      .getAuthApi()
      .getNewRefreshToken(logoutCommand: logoutCommand);

  var tokenPair = tokenPairResponse.data?.anyOf.values[0] as TokenPair;

  authModel.setJwt(tokenPair.jwt);
  authModel.setRefreshToken(tokenPair.refreshToken);

  return;
}

bool isTokenValid(String? token) {
  if (token == null || token.isEmpty) {
    return false;
  } else {
    var claims = JWT.decode(token);
    DateTime expiration = DateTime.fromMillisecondsSinceEpoch(
        claims.payload["exp"] * 1000,
        isUtc: false);

    return expiration.isAfter(DateTime.now());
  }
}

Future<void> storeAppData(
    AuthModel authModel,
    GroupModel groupModel,
    UserModel userModel,
    UserPreferencesModel userPreferencesModel,
    CategoryModel categoryModel,
    TagModel tagModel,
    SystemSettingsModel systemSettingsModel,
    AppData appData) async {
  if (appData!.jwt!.isNotEmpty) {
    await authModel.setJwt(appData.jwt);
  }

  if (appData.refreshToken!.isNotEmpty) {
    await authModel.setRefreshToken(appData.refreshToken);
  }

  authModel.setClaims(appData.claims);
  authModel.setFeatureConfig(appData.featureConfig);
  groupModel.setGroups(appData.groups.toList());
  userModel.setUsers(appData.users.toList());
  userPreferencesModel.setUserPreferences(appData.userPreferences);
  categoryModel.setCategories(appData.categories.toList());
  tagModel.setTags(appData.tags.toList());
  systemSettingsModel.setCurrencyDisplay(appData.currencyDisplay);
  systemSettingsModel.setCurrencyDecimalSeparator(
      appData?.currencyDecimalSeparator ?? CurrencySeparator.period);
  systemSettingsModel.setCurrencyThousandSeparator(
      appData?.currencyThousandthsSeparator ?? CurrencySeparator.comma);
  systemSettingsModel.setCurrencySymbolPosition(
      appData?.currencySymbolPosition ?? CurrencySymbolPosition.END);
  systemSettingsModel.setCurrencyHideDecimalPlaces(
      appData?.currencyHideDecimalPlaces ?? false);
}
