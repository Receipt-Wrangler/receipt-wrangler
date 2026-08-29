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

/// Whether the current user may create receipts in [groupId]. Gates on
/// `group.receipts.create`, mirroring `handlers.CreateReceipt`'s enforcement.
bool canCreateReceipt(PermissionsModel permissionsModel, int groupId) {
  return permissionsModel.hasGroupPermission(
      groupId, Permission.groupPeriodReceiptsPeriodCreate);
}

/// Whether the current user may quick-scan receipts into [groupId]. Gates on
/// `group.receipts.quick-scan`, mirroring `handlers.QuickScan`'s enforcement.
/// Deliberately independent of [canCreateReceipt] — the backend enforces the two
/// permissions separately, so a role may hold either without the other.
bool canQuickScanReceipt(PermissionsModel permissionsModel, int groupId) {
  return permissionsModel.hasGroupPermission(
      groupId, Permission.groupPeriodReceiptsPeriodQuickScan);
}
