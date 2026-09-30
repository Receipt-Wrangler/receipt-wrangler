import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';

/// Builds a [PermissionsModel] seeded with the given app- and group-scoped
/// permissions, mirroring how `setPermissions` hydrates from `AppData` (wire
/// strings, keyed by the string group id). Takes [Permission] values for
/// call-site readability and converts them the way the server would. Use in
/// widget/guard tests that need a caller with a specific permission set without
/// standing up the backend.
///
/// [receiptRequirements] seeds `AppData.groupReceiptRequirements` the same way:
/// only the groups where something is required, keyed by string group id.
PermissionsModel seededPermissions({
  List<Permission> app = const [],
  Map<int, List<Permission>> group = const {},
  Map<int, ReceiptRequirements> receiptRequirements = const {},
}) {
  final model = PermissionsModel();
  model.setPermissions(
    app.map(permissionWireName).toList(),
    {
      for (final entry in group.entries)
        entry.key.toString(): entry.value.map(permissionWireName).toList(),
    },
  );
  model.setReceiptRequirements({
    for (final entry in receiptRequirements.entries)
      entry.key.toString(): entry.value,
  });
  return model;
}

/// Builds the [ReceiptRequirements] a group role resolves to.
ReceiptRequirements receiptRequirements({
  bool comment = false,
  bool image = false,
}) =>
    ReceiptRequirements((b) => b
      ..commentRequired = comment
      ..imageRequired = image);
