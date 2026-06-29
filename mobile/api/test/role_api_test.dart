import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for RoleApi
void main() {
  final instance = Openapi().getRoleApi();

  group(RoleApi, () {
    // Create role
    //
    // This will create an app-scoped or group-scoped role.
    //
    //Future<Role> createRole(UpsertRoleCommand upsertRoleCommand) async
    test('test createRole', () async {
      // TODO
    });

    // List all roles
    //
    // Returns the full pool of roles, both app-scoped and group-scoped.
    //
    //Future<BuiltList<Role>> getRoles() async
    test('test getRoles', () async {
      // TODO
    });

    // Update role
    //
    // Updates an existing app-scoped or group-scoped role. A role's type cannot be changed: the scope in the request must match the existing role's scope. System roles cannot be modified.
    //
    //Future<Role> updateRole(int roleId, UpsertRoleCommand upsertRoleCommand) async
    test('test updateRole', () async {
      // TODO
    });

  });
}
