class E2eEnv {
  static const String baseUrl = String.fromEnvironment('E2E_BASE_URL');
  static const String adminUsername =
      String.fromEnvironment('E2E_ADMIN_USERNAME');
  static const String adminPassword =
      String.fromEnvironment('E2E_ADMIN_PASSWORD');
  static const String userUsername =
      String.fromEnvironment('E2E_USER_USERNAME');
  static const String userPassword =
      String.fromEnvironment('E2E_USER_PASSWORD');

  static void assertAdmin() {
    _require('E2E_BASE_URL', baseUrl);
    _require('E2E_ADMIN_USERNAME', adminUsername);
    _require('E2E_ADMIN_PASSWORD', adminPassword);
  }

  static void assertUser() {
    _require('E2E_BASE_URL', baseUrl);
    _require('E2E_USER_USERNAME', userUsername);
    _require('E2E_USER_PASSWORD', userPassword);
  }

  static void _require(String key, String value) {
    if (value.isEmpty) {
      throw StateError(
        'Missing $key. Source api/dev/switch-to-sqlite.sh and run via '
        'mobile/run-e2e.sh, or pass --dart-define=$key=... directly.',
      );
    }
  }
}
