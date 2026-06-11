import 'dart:async';

import 'package:flutter/foundation.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/category_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/models/system_settings_model.dart';
import 'package:receipt_wrangler_mobile/models/tag_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/models/user_preferences_model.dart';
import 'package:receipt_wrangler_mobile/utils/auth.dart';

/// Serializes all token refresh calls so that only one HTTP request
/// is in-flight at a time. This prevents race conditions with the
/// backend's one-time-use refresh tokens.
///
/// Dart equivalent of the desktop's TokenRefreshService (shareReplay(1)).
class TokenRefreshService {
  static final TokenRefreshService _instance = TokenRefreshService._internal();

  factory TokenRefreshService() => _instance;

  TokenRefreshService._internal();

  Completer<bool>? _refreshCompleter;

  late AuthModel _authModel;
  late GroupModel _groupModel;
  late UserModel _userModel;
  late UserPreferencesModel _userPreferencesModel;
  late CategoryModel _categoryModel;
  late TagModel _tagModel;
  late SystemSettingsModel _systemSettingsModel;
  late PermissionsModel _permissionsModel;

  bool _initialized = false;

  @visibleForTesting
  void resetForTesting() {
    _refreshCompleter = null;
    _initialized = false;
  }

  void initialize({
    required AuthModel authModel,
    required GroupModel groupModel,
    required UserModel userModel,
    required UserPreferencesModel userPreferencesModel,
    required CategoryModel categoryModel,
    required TagModel tagModel,
    required SystemSettingsModel systemSettingsModel,
    required PermissionsModel permissionsModel,
  }) {
    _authModel = authModel;
    _groupModel = groupModel;
    _userModel = userModel;
    _userPreferencesModel = userPreferencesModel;
    _categoryModel = categoryModel;
    _tagModel = tagModel;
    _systemSettingsModel = systemSettingsModel;
    _permissionsModel = permissionsModel;
    _initialized = true;
  }

  /// Returns the current JWT for use by the auth interceptor.
  Future<String?> getCurrentJwt() => _authModel.getJwt();

  /// Serialized token refresh. If a refresh is already in-flight,
  /// all callers share the same Future (and thus the same HTTP request).
  Future<bool> refreshTokens({bool force = false}) async {
    if (!_initialized) return false;

    if (_refreshCompleter != null) {
      return _refreshCompleter!.future;
    }

    _refreshCompleter = Completer<bool>();

    try {
      final result = await _doRefresh(force: force);
      _refreshCompleter!.complete(result);
      return result;
    } catch (e) {
      _refreshCompleter!.complete(false);
      return false;
    } finally {
      _refreshCompleter = null;
    }
  }

  Future<bool> _doRefresh({bool force = false}) async {
    var jwt = await _authModel.getJwt();
    var refreshToken = await _authModel.getRefreshToken();

    bool needsRefresh = force || !isTokenValid(jwt);

    if (!needsRefresh) {
      await _tryLoadAppData();
      return true;
    }

    if (!isTokenValid(refreshToken)) {
      await _authModel.purgeTokens();
      return false;
    }

    try {
      await getAndSetTokens(_authModel);
    } catch (e) {
      print(e);
      await _authModel.purgeTokens();
      return false;
    }

    await _tryLoadAppData();
    return true;
  }

  /// Attempts to load app data, swallowing errors so that a transient
  /// getAppData failure never causes refreshTokens() to report false
  /// (which would trigger token purge and logout).
  Future<void> _tryLoadAppData() async {
    try {
      await _loadAppDataIfNeeded();
    } catch (e) {
      print(e);
    }
  }

  Future<void> _loadAppDataIfNeeded() async {
    if (_groupModel.groups.isEmpty) {
      var appDataResponse =
          await OpenApiClient.client.getUserApi().getAppData();
      await storeAppData(
        _authModel,
        _groupModel,
        _userModel,
        _userPreferencesModel,
        _categoryModel,
        _tagModel,
        _systemSettingsModel,
        _permissionsModel,
        appDataResponse.data as AppData,
      );
    }
  }
}
