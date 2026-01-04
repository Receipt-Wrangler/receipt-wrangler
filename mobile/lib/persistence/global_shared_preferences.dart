import 'package:shared_preferences/shared_preferences.dart';

class GlobalSharedPreferences {
  static SharedPreferences? _instance;

  static Future<void> initialize() async {
    if (_instance == null) {
      _instance = await SharedPreferences.getInstance();
    }
  }

  static SharedPreferences get instance {
    if (_instance == null) {
      throw Exception('GlobalPreferences is not initialized');
    }
    return _instance as SharedPreferences;
  }
}
