import 'package:built_collection/built_collection.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';

/// Builds a [PermissionsModel] seeded with the given app- and group-scoped
/// permissions, mirroring how `setPermissions` hydrates from `AppData` (group
/// keys arrive as the string group ids). Use in widget/guard tests that need a
/// caller with a specific permission set without standing up the backend.
PermissionsModel seededPermissions({
  List<Permission> app = const [],
  Map<int, List<Permission>> group = const {},
}) {
  final model = PermissionsModel();
  model.setPermissions(
    BuiltList<Permission>(app),
    BuiltMap<String, BuiltList<Permission>>({
      for (final entry in group.entries)
        entry.key.toString(): BuiltList<Permission>(entry.value),
    }),
  );
  return model;
}
