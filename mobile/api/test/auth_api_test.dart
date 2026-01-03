import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for AuthApi
void main() {
  final instance = Openapi().getAuthApi();

  group(AuthApi, () {
    // Get fresh tokens
    //
    // This will get a fresh token pair for the user
    //
    //Future<GetNewRefreshToken200Response> getNewRefreshToken({ LogoutCommand logoutCommand }) async
    test('test getNewRefreshToken', () async {
      // TODO
    });

    // Login
    //
    // This will log a user into the system
    //
    //Future<AppData> login(LoginCommand loginCommand) async
    test('test login', () async {
      // TODO
    });

    // Logout
    //
    // This will log a user out of the system and revoke their token [SYSTEM USER]
    //
    //Future logout({ LogoutCommand logoutCommand }) async
    test('test logout', () async {
      // TODO
    });

    // Signs up
    //
    // This will sign a user up for the system
    //
    //Future signUp(SignUpCommand signUpCommand) async
    test('test signUp', () async {
      // TODO
    });

  });
}
