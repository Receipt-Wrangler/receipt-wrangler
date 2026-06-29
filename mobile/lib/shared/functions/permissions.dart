import 'package:openapi/openapi.dart';

import '../../models/permissions_model.dart';

/// Whether the current user may edit receipts in [groupId]. Gates on the modern
/// `group.receipts.update` permission — the legacy GroupRole EDITOR and OWNER
/// tiers both held it, so this preserves the previous behavior.
bool canEditReceipt(PermissionsModel permissionsModel, int groupId) {
  return permissionsModel.hasGroupPermission(
      groupId, Permission.groupPeriodReceiptsPeriodUpdate);
}

/// Whether the current user may add comments to receipts in [groupId]. Gates on
/// `group.comments.create`, mirroring the desktop comment gate and the backend
/// enforcement.
bool canCommentCreate(PermissionsModel permissionsModel, int groupId) {
  return permissionsModel.hasGroupPermission(
      groupId, Permission.groupPeriodCommentsPeriodCreate);
}

/// Whether the current user may delete comments on receipts in [groupId]. Gates
/// on `group.comments.delete`, mirroring the desktop comment gate and the
/// backend enforcement.
bool canCommentDelete(PermissionsModel permissionsModel, int groupId) {
  return permissionsModel.hasGroupPermission(
      groupId, Permission.groupPeriodCommentsPeriodDelete);
}
