import 'package:flutter/material.dart';
import 'package:flutter_secure_storage/flutter_secure_storage.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/persistence/global_shared_preferences.dart';

import '../client/client.dart';

class AuthModel extends ChangeNotifier {
  api.Claims? _claims;

  api.Claims? get claims => _claims;

  final FlutterSecureStorage _storage = const FlutterSecureStorage(
      aOptions: AndroidOptions(
        encryptedSharedPreferences: true,
      ),
      iOptions: IOSOptions(accessibility: KeychainAccessibility.first_unlock));

  final String _refreshTokenKey = "refreshToken";

  final String _jwtKey = "jwt";

  final _basePathKey = "basePath";

  String get basePath =>
      GlobalSharedPreferences.instance.getString(_basePathKey) ?? "";

  api.FeatureConfig _featureConfig = (api.FeatureConfigBuilder()
        ..aiPoweredReceipts = false
        ..enableLocalSignUp = false)
      .build();

  api.FeatureConfig get featureConfig => _featureConfig;

  void initializeAuth() {
    _updateDefaultApiClient();
  }

  void setClaims(api.Claims claims) {
    _claims = claims;

    notifyListeners();
  }

  Future<void> setJwt(String? jwt) async {
    await _storage.write(key: _jwtKey, value: jwt ?? null);
    await _updateDefaultApiClient();

    notifyListeners();
  }

  Future<void> setRefreshToken(String? refreshToken) async {
    await _storage.write(key: _refreshTokenKey, value: refreshToken ?? null);

    await _updateDefaultApiClient();

    notifyListeners();
  }

  Future<void> purgeTokens() async {
    await _storage.delete(key: _jwtKey);
    await _storage.delete(key: _refreshTokenKey);

    await _updateDefaultApiClient();

    notifyListeners();
  }

  Future<String?> getJwt() async {
    return await _storage.read(key: _jwtKey);
  }

  Future<String?> getRefreshToken() async {
    return await _storage.read(key: _refreshTokenKey);
  }

  Future<void> setBasePath(String basePath) async {
    GlobalSharedPreferences.instance.setString(_basePathKey, basePath);

    await _updateDefaultApiClient();

    notifyListeners();
  }

  void setFeatureConfig(api.FeatureConfig? featureConfig) {
    if (featureConfig == null) {
      return;
    } else {
      _featureConfig = featureConfig;
      notifyListeners();
    }
  }

  Future<void> _updateDefaultApiClient() async {
    var jwt = await getJwt();
    var newClient = api.Openapi(basePathOverride: basePath);
    if (jwt != null) {
      newClient.setBearerAuth("bearerAuth", jwt);
    }

    newClient.dio.options.receiveTimeout = Duration(minutes: 5);
    OpenApiClient.client = newClient;
    return;
  }
}
