import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';

/// Builds a [PermissionsModel] seeded with the given app- and group-scoped
/// permissions, mirroring how `setPermissions` hydrates from `AppData` (wire
/// strings, keyed by the string group id). Takes [Permission] values for
/// call-site readability and converts them the way the server would. Use in
/// widget/guard tests that need a caller with a specific permission set without
/// standing up the backend.
PermissionsModel seededPermissions({
  List<Permission> app = const [],
  Map<int, List<Permission>> group = const {},
}) {
  final model = PermissionsModel();
  model.setPermissions(
    app.map(permissionWireName).toList(),
    {
      for (final entry in group.entries)
        entry.key.toString(): entry.value.map(permissionWireName).toList(),
    },
  );
  return model;
}
