import 'package:built_collection/built_collection.dart';
import 'package:flutter/foundation.dart';
import 'package:openapi/openapi.dart';

import '../utils/permission_matcher.dart' as matcher;

/// Holds the calling user's effective permissions (delivered on `AppData`) and
/// answers permission checks used to gate the mobile UI. Mirrors the desktop
/// `AuthState` permission slice (`desktop/src/store/auth.state.ts`): the stored
/// set is a UI hint refreshed on login / app-init via [setPermissions]; the
/// server re-checks the real permissions on every request, so a stale action at
/// worst returns 403.
class PermissionsModel extends ChangeNotifier {
  List<String> _appPermissions = const [];
  Map<int, List<String>> _groupPermissions = const {};

  /// Replaces the stored permissions from `AppData`. The generated client types
  /// these as the `Permission` enum; we convert each to its wire string so the
  /// matcher (a faithful port of the backend matcher) can apply wildcard
  /// semantics over plain strings, exactly like the desktop client. Group keys
  /// arrive as strings (JSON object keys) and are parsed to ints.
  void setPermissions(
    BuiltList<Permission> appPermissions,
    BuiltMap<String, BuiltList<Permission>> groupPermissions,
  ) {
    _appPermissions = appPermissions.map(_toWire).toList(growable: false);
    _groupPermissions = {
      for (final entry in groupPermissions.entries)
        if (int.tryParse(entry.key) != null)
          int.parse(entry.key):
              entry.value.map(_toWire).toList(growable: false),
    };
    notifyListeners();
  }

  /// True when the user holds [permission] at the app scope.
  bool hasAppPermission(Permission permission) =>
      matcher.hasAll(_appPermissions, [_toWire(permission)]);

  /// True when the user holds at least one of [permissions] at the app scope.
  bool hasAnyAppPermission(List<Permission> permissions) =>
      matcher.hasAny(_appPermissions, permissions.map(_toWire).toList());

  /// True when the user holds [permission] in [groupId]. If any of [orApp] is
  /// held at the app scope, the group check is bypassed — mirroring the backend
  /// `OrAppPermissions` (admin-not-a-member) pattern.
  bool hasGroupPermission(
    int groupId,
    Permission permission, {
    List<Permission> orApp = const [],
  }) {
    if (orApp.isNotEmpty &&
        matcher.hasAny(_appPermissions, orApp.map(_toWire).toList())) {
      return true;
    }
    return matcher.hasAll(
      _groupPermissions[groupId] ?? const [],
      [_toWire(permission)],
    );
  }

  /// True when [permission] is held in ANY group the user belongs to. Used where
  /// there is no single current group (e.g. the group-select / all-groups add
  /// menu).
  bool hasGroupPermissionInAnyGroup(Permission permission) {
    final required = [_toWire(permission)];
    return _groupPermissions.values
        .any((granted) => matcher.hasAll(granted, required));
  }

  static String _toWire(Permission permission) =>
      serializers.serializeWith(Permission.serializer, permission) as String;
}
