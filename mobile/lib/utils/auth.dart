import 'package:dart_jsonwebtoken/dart_jsonwebtoken.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/category_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/models/tag_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/models/user_preferences_model.dart';

import '../models/system_settings_model.dart';

Future<void> getAndSetTokens(AuthModel authModel) async {
  var refreshToken = await authModel.getRefreshToken() ?? "";
  var logoutCommand =
      (LogoutCommandBuilder()..refreshToken = refreshToken).build();
  var tokenPairResponse = await OpenApiClient.client
      .getAuthApi()
      .getNewRefreshToken(logoutCommand: logoutCommand);

  var tokenPair = tokenPairResponse.data?.anyOf.values[0] as TokenPair;

  await authModel.setTokens(tokenPair.jwt, tokenPair.refreshToken);
}

bool isTokenValid(String? token) {
  if (token == null || token.isEmpty) {
    return false;
  }

  try {
    var claims = JWT.decode(token);
    DateTime expiration = DateTime.fromMillisecondsSinceEpoch(
        claims.payload["exp"] * 1000,
        isUtc: false);

    return expiration.isAfter(DateTime.now());
  } catch (_) {
    return false;
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
    PermissionsModel permissionsModel,
    AppData appData) async {
  if ((appData.jwt?.isNotEmpty ?? false) &&
      (appData.refreshToken?.isNotEmpty ?? false)) {
    await authModel.setTokens(appData.jwt, appData.refreshToken);
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
  permissionsModel.setPermissions(
      appData.appPermissions, appData.groupPermissions);
}
