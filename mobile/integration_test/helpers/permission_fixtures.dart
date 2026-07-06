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

/// Fetches the full `RoleView` map for the system role [name] within [scope]
/// (`APP` or `GROUP`) via `GET /role`.
Future<Map<String, dynamic>> _findRole(
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
  return roles.firstWhere(
    (r) => r['name'] == name && r['scope'] == scope,
    orElse: () => throw StateError(
      'No $scope role named "$name". Available: '
      '${roles.where((r) => r['scope'] == scope).map((r) => r['name']).toList()}',
    ),
  );
}

/// Resolves the id of a system role by [name] within [scope] (`APP` or `GROUP`)
/// via `GET /role`.
Future<int> _roleIdByName(String name, String scope, {required String jwt}) async =>
    (await _findRole(name, scope, jwt: jwt))['id'] as int;

/// Returns the permission strings held by the system role [name] within [scope].
/// Used by the negative-permission specs to derive a "Legacy role minus one
/// permission" set — the custom role is then identical to the Legacy baseline
/// except the single permission under test (robust against registry drift).
Future<List<String>> rolePermissionsByName(
  String name,
  String scope, {
  required String jwt,
}) async =>
    ((await _findRole(name, scope, jwt: jwt))['permissions'] as List)
        .cast<String>();

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

/// Creates a regular account and returns its id. The admin create endpoint
/// requires an `appRoleId`; [appRoleId] assigns a specific app role (e.g. a
/// custom one minus a permission), defaulting to Legacy User (the default app
/// role, which carries no group permissions — the user only gets the group
/// permissions a fixture grants it).
Future<int> createUser({
  required String username,
  required String password,
  required String displayName,
  required String jwt,
  int? appRoleId,
}) async {
  final resolvedAppRoleId =
      appRoleId ?? await appRoleIdByName('Legacy User', jwt: jwt);
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/user/'),
        headers: _jsonAuth(jwt),
        body: jsonEncode({
          'username': username,
          'password': password,
          'displayName': displayName,
          'appRoleId': resolvedAppRoleId,
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

/// Creates a custom role via `POST /role/` and returns its id. [scope] is `APP`
/// or `GROUP`; every entry in [permissions] must belong to that scope (validated
/// server-side against the registry). Used by the negative-permission specs to
/// mint a role mirroring a Legacy role minus the single permission under test.
Future<int> createRole({
  required String name,
  required String scope,
  required List<String> permissions,
  required String jwt,
  String description = 'e2e custom role',
  bool includeOwnPaidReceipts = false,
  List<int> paidByUserGrants = const [],
}) async {
  final body = <String, dynamic>{
    'name': name,
    'description': description,
    'scope': scope,
    'permissions': permissions,
  };
  // Group-role paid-by visibility (group scope only; empty/false = unrestricted,
  // i.e. members see every payer's receipts). Mirrors the desktop createRole
  // `paidByOwn` / `paidByUsers` options. Only sent when restricting, so APP-role
  // creates are unaffected.
  if (includeOwnPaidReceipts || paidByUserGrants.isNotEmpty) {
    body['includeOwnPaidReceipts'] = includeOwnPaidReceipts;
    body['paidByUserGrants'] = paidByUserGrants;
  }
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/role/'),
        headers: _jsonAuth(jwt),
        body: jsonEncode(body),
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError('createRole($name) failed: '
        'HTTP ${res.statusCode}: ${res.body}');
  }
  return (jsonDecode(res.body) as Map<String, dynamic>)['id'] as int;
}

/// Best-effort `DELETE /role/{roleId}?scope=…`. The backend refuses to delete an
/// *assigned* role, so this must run only after the user/group that reference it
/// are gone — see the teardown ordering in [provisionUserWithoutAppPermission] /
/// [provisionGroupMemberWithoutPermission]. Swallows errors like [deleteUser].
Future<void> deleteRole(
  int roleId, {
  required String scope,
  required String jwt,
}) async {
  try {
    await http
        .delete(
          Uri.parse('${E2eEnv.baseUrl}/role/$roleId?scope=$scope'),
          headers: _auth(jwt),
        )
        .timeout(const Duration(seconds: 5));
  } catch (_) {
    // best-effort
  }
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

/// Creates a global category (categories are app-wide; group visibility is via
/// group-role grants) and returns its id. A group role with no category grants
/// is unrestricted, so a freshly-created category shows up in every member's
/// per-group catalog. Used to prove non-admins receive categories via the
/// per-group `groupCategories` map (not the admin-only flat list).
Future<int> createCategory({
  required String name,
  required String jwt,
  String description = 'e2e category',
}) async {
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/category/'),
        headers: _jsonAuth(jwt),
        body: jsonEncode({'name': name, 'description': description}),
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError('createCategory($name) failed: '
        'HTTP ${res.statusCode}: ${res.body}');
  }
  return (jsonDecode(res.body) as Map<String, dynamic>)['id'] as int;
}

/// Creates a comment on [receiptId] (as the caller behind [jwt]) and returns its
/// id. Seeds an existing comment so the comment swipe-to-delete gate
/// (`group.comments.delete`) has something to render.
Future<int> createComment({
  required int receiptId,
  required String jwt,
  String comment = 'e2e seeded comment',
}) async {
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/comment/'),
        headers: _jsonAuth(jwt),
        body: jsonEncode({'comment': comment, 'receiptId': receiptId}),
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError('createComment($receiptId) failed: '
        'HTTP ${res.statusCode}: ${res.body}');
  }
  return (jsonDecode(res.body) as Map<String, dynamic>)['id'] as int;
}

/// Best-effort `DELETE /category/{id}`. Swallows errors like the other cleanups.
Future<void> deleteCategory(int categoryId, {required String jwt}) async {
  try {
    await http
        .delete(Uri.parse('${E2eEnv.baseUrl}/category/$categoryId'),
            headers: _auth(jwt))
        .timeout(const Duration(seconds: 5));
  } catch (_) {
    // best-effort
  }
}

/// Creates a global tag (the tag analogue of [createCategory]) and returns its
/// id. Used to prove non-admins receive tags via the per-group `groupTags` map.
Future<int> createTag({
  required String name,
  required String jwt,
  String description = 'e2e tag',
}) async {
  final res = await http
      .post(
        Uri.parse('${E2eEnv.baseUrl}/tag/'),
        headers: _jsonAuth(jwt),
        body: jsonEncode({'name': name, 'description': description}),
      )
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError('createTag($name) failed: '
        'HTTP ${res.statusCode}: ${res.body}');
  }
  return (jsonDecode(res.body) as Map<String, dynamic>)['id'] as int;
}

/// Best-effort `DELETE /tag/{id}`. Swallows errors like the other cleanups.
Future<void> deleteTag(int tagId, {required String jwt}) async {
  try {
    await http
        .delete(Uri.parse('${E2eEnv.baseUrl}/tag/$tagId'), headers: _auth(jwt))
        .timeout(const Duration(seconds: 5));
  } catch (_) {
    // best-effort
  }
}

/// The admin's first non-"All" group (id + name), read via the API. Use to
/// target a group for config persistence + dropdown selection before login.
/// Every group `GET /group/` returns for the admin has the admin as a member, so
/// the group's paid-by dropdown always includes the admin.
Future<({int id, String name})> firstNonAllGroup(String jwt) async {
  final groups = await _adminGroups(jwt);
  final g = groups.firstWhere((x) => x['isAllGroup'] != true,
      orElse: () => throw StateError('no non-all group for the admin'));
  return (id: g['id'] as int, name: g['name'] as String);
}

Future<List<Map<String, dynamic>>> _adminGroups(String jwt) async {
  final res = await http
      .get(Uri.parse('${E2eEnv.baseUrl}/group/'), headers: _auth(jwt))
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError('GET /group/ failed: HTTP ${res.statusCode}: ${res.body}');
  }
  return (jsonDecode(res.body) as List).cast<Map<String, dynamic>>();
}

/// Builds an `UpdateGroupReceiptSettingsCommand` from a settings map, applying
/// [overrides]. Only the hide* + quick-scan enabled/required flags are sent -- the
/// default enum fields (`quickScanDefaultPaidByType` / `...Status`) are omitted
/// because the backend keeps them and rejects an empty enum, so we never echo one
/// back. (This is why persisted configs keep paid-by/status *required*: making
/// them optional would need a persisted default the backend enforces.)
Map<String, dynamic> _settingsToCommand(
  Map<String, dynamic> s, {
  Map<String, dynamic> overrides = const {},
}) =>
    {
      for (final k in const [
        'hideImages', 'hideReceiptCategories', 'hideReceiptTags',
        'hideItemCategories', 'hideItemTags', 'hideComments',
        'hideShareCategories', 'hideShareTags',
        'quickScanPaidByEnabled', 'quickScanPaidByRequired',
        'quickScanStatusEnabled', 'quickScanStatusRequired',
        'quickScanCategoriesEnabled', 'quickScanCategoriesRequired',
        'quickScanTagsEnabled', 'quickScanTagsRequired',
      ])
        k: s[k] ?? false,
      ...overrides,
    };

Future<void> _putGroupReceiptSettings(
    int groupId, String jwt, Map<String, dynamic> command) async {
  final res = await http
      .put(Uri.parse('${E2eEnv.baseUrl}/group/$groupId/groupReceiptSettings'),
          headers: _jsonAuth(jwt), body: jsonEncode(command))
      .timeout(const Duration(seconds: 10));
  if (res.statusCode != 200) {
    throw StateError('PUT groupReceiptSettings($groupId) failed: '
        'HTTP ${res.statusCode}: ${res.body}');
  }
}

/// Persists quick-scan [overrides] (merged over the group's current settings) on
/// [groupId], restoring the original settings on teardown.
///
/// Submit tests must persist -- not client-mutate GroupModel -- because the
/// backend's `resolveQuickScanFields` validates the submit against the group's
/// *persisted* config, so client (via AppData at login) and server must agree, as
/// they do in production. Keep paid-by/status required in [overrides] (see
/// `_settingsToCommand`).
Future<void> setGroupQuickScanConfig({
  required int groupId,
  required String jwt,
  required Map<String, dynamic> overrides,
}) async {
  final groups = await _adminGroups(jwt);
  final original = ((groups.firstWhere((x) => x['id'] == groupId,
              orElse: () => throw StateError('group $groupId not found'))[
          'groupReceiptSettings']) as Map)
      .cast<String, dynamic>();
  await _putGroupReceiptSettings(
      groupId, jwt, _settingsToCommand(original, overrides: overrides));
  addTearDown(() async {
    final j = await apiLogin();
    await _putGroupReceiptSettings(groupId, j, _settingsToCommand(original));
  });
}

/// Provisions a fresh user and (optionally) a fixture group it belongs to,
/// registering `addTearDown` to delete the group and user afterwards.
///
/// - [groupRoleId] (wins over [roleName]) / [roleName] → a group is created with
///   the user holding that group role. [groupRoleId] takes an explicit role id
///   (e.g. a custom role from [createRole]); [roleName] resolves a system role
///   ("Legacy Viewer" / "Legacy Editor" / "Legacy Owner") by name. Both null →
///   the user belongs to no group (the add-menu "no permission" negative case).
/// - [appRoleId] assigns a specific app role (default: Legacy User).
/// - [withReceipt] seeds one receipt into the group (for the receipt-edit /
///   swipe gates). Ignored when the user belongs to no group.
///
/// Provision BEFORE `loginAs` so the user's permissions hydrate from AppData at
/// login. The returned [PermFixture] carries the login credentials.
Future<PermFixture> provisionPermUser({
  String? roleName,
  int? groupRoleId,
  int? appRoleId,
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
    appRoleId: appRoleId,
  );
  addTearDown(() async => deleteUser(userId, jwt: await apiLogin()));

  // An explicit group role id wins; otherwise resolve a system role by name.
  int? resolvedGroupRoleId = groupRoleId;
  if (resolvedGroupRoleId == null && roleName != null) {
    resolvedGroupRoleId = await groupRoleIdByName(roleName, jwt: jwt);
  }

  int? groupId;
  String? groupName;
  int? receiptId;
  String? receiptName;

  if (resolvedGroupRoleId != null) {
    groupName = 'e2e-perm-$suffix';
    groupId = await createGroupWithMember(
      name: groupName,
      memberUserId: userId,
      groupRoleId: resolvedGroupRoleId,
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

/// Provisions a user whose APP role is "Legacy User" **minus** [permission],
/// registering teardown for both the user and the custom role. Use for negative
/// app-permission specs (e.g. `app.receipts.search`). The user belongs to no
/// fixture group — the custom app role is the only thing under test.
///
/// The custom role is the Legacy User permission set with exactly [permission]
/// removed, so the only behavioral difference from a Legacy User is that single
/// missing permission. The backend refuses to delete an assigned role, so the
/// role-delete teardown is registered **before** [provisionPermUser] — with
/// LIFO teardown that makes it run **after** the user delete, when the role is
/// unassigned and deletable.
Future<PermFixture> provisionUserWithoutAppPermission(String permission) async {
  final jwt = await apiLogin(); // admin
  final perms = (await rolePermissionsByName('Legacy User', 'APP', jwt: jwt))
      .where((p) => p != permission)
      .toList();
  final roleId = await createRole(
    name: 'e2e-app-${_unique()}',
    scope: 'APP',
    permissions: perms,
    jwt: jwt,
  );
  addTearDown(() async => deleteRole(roleId, scope: 'APP', jwt: await apiLogin()));

  return provisionPermUser(appRoleId: roleId);
}

/// Provisions a user in a fixture group whose GROUP role is [baselineRole]
/// **minus** [permission], registering teardown for the user, group and custom
/// role. Use for negative group-permission specs (e.g. `group.dashboards.read`).
///
/// [baselineRole] defaults to "Legacy Viewer". Pass "Legacy Editor" when the
/// scenario needs a permission the Viewer lacks — e.g. the comment-gate specs
/// reach the edit-state comment screen, which requires `group.receipts.update`
/// (Viewer is read-only).
///
/// Same minus-one rationale and LIFO teardown ordering as
/// [provisionUserWithoutAppPermission]: the role-delete is registered first so
/// it runs after the group (and its membership) and the user are gone.
Future<PermFixture> provisionGroupMemberWithoutPermission(
  String permission, {
  bool withReceipt = false,
  String baselineRole = 'Legacy Viewer',
}) async {
  final jwt = await apiLogin(); // admin
  final perms = (await rolePermissionsByName(baselineRole, 'GROUP', jwt: jwt))
      .where((p) => p != permission)
      .toList();
  final roleId = await createRole(
    name: 'e2e-group-${_unique()}',
    scope: 'GROUP',
    permissions: perms,
    jwt: jwt,
  );
  addTearDown(
      () async => deleteRole(roleId, scope: 'GROUP', jwt: await apiLogin()));

  return provisionPermUser(groupRoleId: roleId, withReceipt: withReceipt);
}

/// A [PermFixture] plus the two receipts seeded for a paid-by-visibility spec:
/// [ownReceiptName] is paid by the member (visible to them) and
/// [hiddenReceiptName] is paid by the admin (filtered out by the role's
/// "their own receipts" paid-by grant).
class PaidByFixture {
  PaidByFixture({
    required this.fixture,
    required this.ownReceiptId,
    required this.ownReceiptName,
    required this.hiddenReceiptId,
    required this.hiddenReceiptName,
  });

  final PermFixture fixture;
  final int ownReceiptId;
  final String ownReceiptName;
  final int hiddenReceiptId;
  final String hiddenReceiptName;
}

/// Provisions a group member whose GROUP role holds the full Legacy Viewer
/// permission set (so `group.receipts.read` is held — the denial is paid-by, not
/// a missing permission) but is restricted to **their own receipts** on the
/// paid-by axis. Seeds two receipts in the group: one paid by the member (own,
/// visible) and one paid by the admin (hidden). Mirrors the desktop
/// `paid-by-visibility.spec.ts` provisioning.
///
/// Teardown follows the same LIFO ordering as the other negative-permission
/// fixtures: the role delete is registered before [provisionPermUser]'s
/// user/group deletes, so it runs last (once the assignment is gone). The seeded
/// receipts are cascade-deleted with the group.
Future<PaidByFixture> provisionPaidByOwnMember() async {
  final jwt = await apiLogin(); // admin
  final adminId = await userIdByUsername(E2eEnv.adminUsername, jwt: jwt);
  final viewerPerms =
      await rolePermissionsByName('Legacy Viewer', 'GROUP', jwt: jwt);
  final roleId = await createRole(
    name: 'e2e-paidby-${_unique()}',
    scope: 'GROUP',
    permissions: viewerPerms,
    jwt: jwt,
    includeOwnPaidReceipts: true,
  );
  addTearDown(
      () async => deleteRole(roleId, scope: 'GROUP', jwt: await apiLogin()));

  final fixture = await provisionPermUser(groupRoleId: roleId);

  final suffix = _unique();
  final ownReceiptName = 'e2e-paidby-own-$suffix';
  final hiddenReceiptName = 'e2e-paidby-hidden-$suffix';
  final seedJwt = await apiLogin();
  final hiddenReceiptId = await createReceipt(
    groupId: fixture.groupId!,
    paidByUserId: adminId,
    jwt: seedJwt,
    name: hiddenReceiptName,
  );
  final ownReceiptId = await createReceipt(
    groupId: fixture.groupId!,
    paidByUserId: fixture.userId,
    jwt: seedJwt,
    name: ownReceiptName,
  );

  return PaidByFixture(
    fixture: fixture,
    ownReceiptId: ownReceiptId,
    ownReceiptName: ownReceiptName,
    hiddenReceiptId: hiddenReceiptId,
    hiddenReceiptName: hiddenReceiptName,
  );
}
