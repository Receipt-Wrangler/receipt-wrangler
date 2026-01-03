import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';

class UserPreferencesModel extends ChangeNotifier {
  UserPreferences _userPreferences = (UserPreferencesBuilder()
        ..id = 0
        ..userId = 0
        ..createdAt = "")
      .build();

  UserPreferences get userPreferences => _userPreferences;

  void setUserPreferences(UserPreferences userPreferences) {
    _userPreferences = userPreferences;

    notifyListeners();
  }
}
