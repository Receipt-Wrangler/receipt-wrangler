import 'package:flutter/foundation.dart';
import 'package:openapi/openapi.dart';

import '../utils/permission_matcher.dart' as matcher;

/// The wire string for a generated [Permission] value. Exposed so tests and
/// helpers can seed a model the same way `AppData` does.
String permissionWireName(Permission permission) =>
    serializers.serializeWith(Permission.serializer, permission) as String;

/// Holds the calling user's effective permissions (delivered on `AppData`) and
/// answers permission checks used to gate the mobile UI. Mirrors the desktop
/// `AuthState` permission slice (`desktop/src/store/auth.state.ts`): the stored
/// set is a UI hint refreshed on login / app-init via [setPermissions]; the
/// server re-checks the real permissions on every request, so a stale action at
/// worst returns 403.
class PermissionsModel extends ChangeNotifier {
  List<String> _appPermissions = const [];
  Map<int, List<String>> _groupPermissions = const {};

  /// Replaces the stored permissions from `AppData`.
  ///
  /// These arrive — and are stored — as plain wire **strings**, deliberately
  /// not as the generated `Permission` enum. "Which permissions does this user
  /// hold" is server-authoritative data the client only pattern-matches; the
  /// enum is the catalog of which permissions *exist*. Typing the payload as a
  /// closed enum made every shipped build fail on the next backend permission
  /// (built_value's `_$valueOf` throws on an unknown value, which fails the
  /// entire `AppData` parse and hard-fails login) — that shipped twice. It also
  /// could never represent a wildcard grant, which the matcher supports.
  ///
  /// The matcher (a faithful port of the backend matcher) applies wildcard
  /// semantics over these strings, exactly like the desktop client. Group keys
  /// arrive as strings (JSON object keys) and are parsed to ints.
  void setPermissions(
    Iterable<String> appPermissions,
    Map<String, Iterable<String>> groupPermissions,
  ) {
    _appPermissions = appPermissions.toList(growable: false);
    _groupPermissions = {
      for (final entry in groupPermissions.entries)
        if (int.tryParse(entry.key) != null)
          int.parse(entry.key): entry.value.toList(growable: false),
    };
    notifyListeners();
  }

  /// True when the user holds [permission] at the app scope.
  bool hasAppPermission(Permission permission) =>
      matcher.hasAll(_appPermissions, [permissionWireName(permission)]);

  /// True when the user holds at least one of [permissions] at the app scope.
  bool hasAnyAppPermission(List<Permission> permissions) =>
      matcher.hasAny(_appPermissions, permissions.map(permissionWireName).toList());

  /// True when the user holds [permission] in [groupId]. If any of [orApp] is
  /// held at the app scope, the group check is bypassed — mirroring the backend
  /// `OrAppPermissions` (admin-not-a-member) pattern.
  bool hasGroupPermission(
    int groupId,
    Permission permission, {
    List<Permission> orApp = const [],
  }) {
    if (orApp.isNotEmpty &&
        matcher.hasAny(_appPermissions, orApp.map(permissionWireName).toList())) {
      return true;
    }
    return matcher.hasAll(
      _groupPermissions[groupId] ?? const [],
      [permissionWireName(permission)],
    );
  }

  /// True when [permission] is held in ANY group the user belongs to. Used where
  /// there is no single current group (e.g. the group-select / all-groups add
  /// menu).
  bool hasGroupPermissionInAnyGroup(Permission permission) {
    final required = [permissionWireName(permission)];
    return _groupPermissions.values
        .any((granted) => matcher.hasAll(granted, required));
  }
}
