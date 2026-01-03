import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for UserApi
void main() {
  final instance = Openapi().getUserApi();

  group(UserApi, () {
    // Converts dummy user
    //
    // This will convert a dummy user to a normal system user, [SYSTEM ADMIN]
    //
    //Future convertDummyUserById(int userId, ResetPasswordCommand resetPasswordCommand) async
    test('test convertDummyUserById', () async {
      // TODO
    });

    // Create user
    //
    // This will to create a user, [SYSTEM ADMIN]
    //
    //Future createUser(User user) async
    test('test createUser', () async {
      // TODO
    });

    // Delete user
    //
    // This will delete a system user by id [SYSTEM ADMIN]
    //
    //Future deleteUserById(int userId) async
    test('test deleteUserById', () async {
      // TODO
    });

    // Get amount owed for user
    //
    // This will return the amount owed for the logged in user, in the specified group, [SYSTEM USER]
    //
    //Future<BuiltMap<String, String>> getAmountOwedForUser({ int groupId, BuiltList<int> receiptIds }) async
    test('test getAmountOwedForUser', () async {
      // TODO
    });

    // Get app data
    //
    // This will return the user's app data for the currently logged in user [SYSTEM USER]
    //
    //Future<AppData> getAppData() async
    test('test getAppData', () async {
      // TODO
    });

    // Get claims for logged in user
    //
    // This will return the user's token claims for the currently logged in user [SYSTEM USER]
    //
    //Future<Claims> getUserClaims() async
    test('test getUserClaims', () async {
      // TODO
    });

    // Get username count
    //
    // This will return the number of users in the system with the same username
    //
    //Future<int> getUsernameCount(String username) async
    test('test getUsernameCount', () async {
      // TODO
    });

    // Get users
    //
    // This will get all the users in the system and return a view without sensative information
    //
    //Future<BuiltList<UserView>> getUsers() async
    test('test getUsers', () async {
      // TODO
    });

    // Reset password
    //
    // This will reset a password for a user, [SYSTEM ADMIN]
    //
    //Future resetPasswordById(int userId, ResetPasswordCommand resetPasswordCommand) async
    test('test resetPasswordById', () async {
      // TODO
    });

    // Update user by id
    //
    // This will update a user by id, [SYSTEM ADMIN]
    //
    //Future updateUserById(int userId, User user) async
    test('test updateUserById', () async {
      // TODO
    });

    // Update user profile
    //
    // This will update the logged in user's user profile
    //
    //Future updateUserProfile(UpdateProfileCommand updateProfileCommand) async
    test('test updateUserProfile', () async {
      // TODO
    });

  });
}
