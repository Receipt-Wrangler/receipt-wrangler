import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for UserPreferencesApi
void main() {
  final instance = Openapi().getUserPreferencesApi();

  group(UserPreferencesApi, () {
    // Get user preferences
    //
    // This will return the user's preferences for the currently logged in user [SYSTEM USER]
    //
    //Future<UserPreferences> getUserPreferences() async
    test('test getUserPreferences', () async {
      // TODO
    });

    // Update user preferences
    //
    // This will update the user's preferences for the currently logged in user [SYSTEM USER]
    //
    //Future<UserPreferences> updateUserPreferences(UserPreferences userPreferences) async
    test('test updateUserPreferences', () async {
      // TODO
    });

  });
}
