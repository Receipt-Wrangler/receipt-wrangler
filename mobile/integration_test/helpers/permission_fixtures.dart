import 'dart:convert';

import 'package:flutter_test/flutter_test.dart';
import 'package:http/http.dart' as http;

import 'api.dart';
import 'env.dart';

/// Admin-API fixtures for the permission-gating e2e specs.
///
/// The permission gates (`receipt_edit_popup_menu.dart`, `receipt_list_item.dart`,
/// `group_activity_list_item.dart`, `show_add_menu.dart`) read the caller's
/// *group-scoped* permissions, so to exercise them we need a logged-in user whose
/// group membership/role we control.
///
/// We provision a **fresh, uniquely-named user per spec** rather than reusing the
/// shared `e2e-user`: that account accumulates group memberships across runs
/// (other specs add it to groups and don't always clean up), which would poison
/// the "held in any group" add-menu fallback and the "user in no group" negative
/// case. A brand-new user is deterministically a member of zero groups until a
/// fixture adds it to exactly one.
///
/// All provisioning runs over the admin API (`e2e-admin` holds every app
/// permission and becomes Legacy Owner of any group it creates). Everything is
/// torn down via `addTearDown`.
class PermFixture {
  PermFixture({
    required this.username,
    required this.password,
    required this.userId,
    this.groupId,
    this.groupName,
    this.receiptId,
    this.receiptName,
  });

  /// Credentials for the provisioned user — pass to
  /// `loginAs(tester, username: f.username, password: f.password)`.
  final String username;
  final String password;
  final int userId;

  /// The fixture group the user belongs to (null when provisioned with no group,
  /// e.g. the "user in no group" add-menu negative case).
  final int? groupId;
  final String? groupName;

  /// A receipt seeded into [groupId] (null unless `withReceipt: true`).
  final int? receiptId;
  final String? receiptName;
}

const _password = 'perm-user-password';

Map<String, String> _jsonAuth(String jwt) => {
      'Content-Type': 'application/json',
      'Cookie': 'jwt=$jwt',
    };

Map<String, String> _auth(String jwt) => {'Cookie': 'jwt=$jwt'};

String _unique() => DateTime.now().microsecondsSinceEpoch.toString();

/// Resolves the id of a system role by [name] within [scope] (`APP` or `GROUP`)
/// via `GET /role`.
Future<int> _roleIdByName(
  String name,
  String scope, {
  required String jwt,
}) async {
  final res = await http
      .get(Uri.parse('${E2eEnv.baseUrl}/role'), headers: _auth(jwt))
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError('GET /role failed: HTTP ${res.statusCode}: ${res.body}');
  }
  final roles = (jsonDecode(res.body) as List).cast<Map<String, dynamic>>();
  final match = roles.firstWhere(
    (r) => r['name'] == name && r['scope'] == scope,
    orElse: () => throw StateError(
      'No $scope role named "$name". Available: '
      '${roles.where((r) => r['scope'] == scope).map((r) => r['name']).toList()}',
    ),
  );
  return match['id'] as int;
}

/// Resolves the id of a system group role by name (e.g. "Legacy Viewer",
/// "Legacy Editor", "Legacy Owner") via `GET /role`.
Future<int> groupRoleIdByName(String name, {required String jwt}) =>
    _roleIdByName(name, 'GROUP', jwt: jwt);

/// Resolves the id of a system app role by name (e.g. "Legacy User",
/// "Legacy Admin") via `GET /role`.
Future<int> appRoleIdByName(String name, {required String jwt}) =>
    _roleIdByName(name, 'APP', jwt: jwt);

/// Resolves a user's id by username via `GET /user/` (admin only).
Future<int> userIdByUsername(String username, {required String jwt}) async {
  final res = await http
      .get(Uri.parse('${E2eEnv.baseUrl}/user/'), headers: _auth(jwt))
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError('GET /user/ failed: HTTP ${res.statusCode}: ${res.body}');
  }
  final users = (jsonDecode(res.body) as List).cast<Map<String, dynamic>>();
  final match = users.firstWhere(
    (u) => u['username'] == username,
    orElse: () =>
        throw StateError('No user with username "$username" in GET /user/'),
  );
  return match['id'] as int;
}

/// Creates a regular account assigned the Legacy User app role and returns its
/// id. Legacy User carries no group permissions, so the user only gets the
/// group permissions a fixture grants it. The admin create endpoint requires an
/// `appRoleId`, so we resolve Legacy User (the default app role) by name.
Future<int> createUser({
  required String username,
  required String password,
  required String displayName,
  required String jwt,
}) async {
  final appRoleId = await appRoleIdByName('Legacy User', jwt: jwt);
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/user/'),
        headers: _jsonAuth(jwt),
        body: jsonEncode({
          'username': username,
          'password': password,
          'displayName': displayName,
          'appRoleId': appRoleId,
          'isDummyUser': false,
        }),
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError('createUser($username) failed: '
        'HTTP ${res.statusCode}: ${res.body}');
  }
  return (jsonDecode(res.body) as Map<String, dynamic>)['id'] as int;
}

/// Best-effort `DELETE /user/{id}`. Swallows errors so a cleanup failure (e.g.
/// the user was already removed) doesn't mask the test result.
Future<void> deleteUser(int userId, {required String jwt}) async {
  try {
    await http
        .delete(Uri.parse('${E2eEnv.baseUrl}/user/$userId'), headers: _auth(jwt))
        .timeout(const Duration(seconds: 5));
  } catch (_) {
    // best-effort
  }
}

/// Creates a group whose creator (the admin behind [jwt]) is auto-added as
/// Owner, plus [memberUserId] with [groupRoleId]. Returns the created group id.
///
/// `CreateGroup` persists the command's `groupMembers` AND separately adds the
/// token user as Owner, so one POST yields admin-owner + the member-with-role —
/// no `PUT /group/{id}` member replacement needed.
Future<int> createGroupWithMember({
  required String name,
  required int memberUserId,
  required int groupRoleId,
  required String jwt,
}) async {
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/group/'),
        headers: _jsonAuth(jwt),
        body: jsonEncode({
          'name': name,
          'status': 'ACTIVE',
          'groupMembers': [
            {'userId': memberUserId, 'groupRoleId': groupRoleId},
          ],
        }),
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError('createGroupWithMember($name) failed: '
        'HTTP ${res.statusCode}: ${res.body}');
  }
  return (jsonDecode(res.body) as Map<String, dynamic>)['id'] as int;
}

/// Best-effort `DELETE /group/{id}` (cascades the group's receipts).
Future<void> deleteGroup(int groupId, {required String jwt}) async {
  try {
    await http
        .delete(Uri.parse('${E2eEnv.baseUrl}/group/$groupId'),
            headers: _auth(jwt))
        .timeout(const Duration(seconds: 5));
  } catch (_) {
    // best-effort
  }
}

/// Creates a receipt in [groupId] paid by [paidByUserId]. Returns its id.
Future<int> createReceipt({
  required int groupId,
  required int paidByUserId,
  required String jwt,
  required String name,
  String amount = '12.34',
}) async {
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/receipt/'),
        headers: _jsonAuth(jwt),
        body: jsonEncode({
          'name': name,
          'amount': amount,
          'date': '2026-06-11T00:00:00Z',
          'groupId': groupId,
          'paidByUserId': paidByUserId,
          'status': 'OPEN',
        }),
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError('createReceipt($name) failed: '
        'HTTP ${res.statusCode}: ${res.body}');
  }
  return (jsonDecode(res.body) as Map<String, dynamic>)['id'] as int;
}

/// Provisions a fresh user and (optionally) a fixture group it belongs to,
/// registering `addTearDown` to delete the group and user afterwards.
///
/// - [roleName] null → user belongs to no group (the add-menu "no permission"
///   negative case). Non-null → a group is created with the user holding that
///   system role ("Legacy Viewer" / "Legacy Editor" / "Legacy Owner").
/// - [withReceipt] seeds one receipt into the group (for the receipt-edit /
///   swipe gates). Ignored when [roleName] is null.
///
/// Provision BEFORE `loginAs` so the user's permissions hydrate from AppData at
/// login. The returned [PermFixture] carries the login credentials.
Future<PermFixture> provisionPermUser({
  String? roleName,
  bool withReceipt = false,
}) async {
  final jwt = await apiLogin(); // admin
  final suffix = _unique();
  final username = 'perm-$suffix';

  final adminId = await userIdByUsername(E2eEnv.adminUsername, jwt: jwt);
  final userId = await createUser(
    username: username,
    password: _password,
    displayName: 'Perm $suffix',
    jwt: jwt,
  );
  addTearDown(() async => deleteUser(userId, jwt: await apiLogin()));

  int? groupId;
  String? groupName;
  int? receiptId;
  String? receiptName;

  if (roleName != null) {
    final roleId = await groupRoleIdByName(roleName, jwt: jwt);
    groupName = 'e2e-perm-$suffix';
    groupId = await createGroupWithMember(
      name: groupName,
      memberUserId: userId,
      groupRoleId: roleId,
      jwt: jwt,
    );
    addTearDown(() async => deleteGroup(groupId!, jwt: await apiLogin()));

    if (withReceipt) {
      receiptName = 'e2e-perm-rcpt-$suffix';
      receiptId = await createReceipt(
        groupId: groupId,
        paidByUserId: adminId,
        jwt: jwt,
        name: receiptName,
      );
    }
  }

  return PermFixture(
    username: username,
    password: _password,
    userId: userId,
    groupId: groupId,
    groupName: groupName,
    receiptId: receiptId,
    receiptName: receiptName,
  );
}
