import 'package:flutter/widgets.dart';
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';

import '../../models/auth_model.dart';
import '../../models/group_model.dart';
import '../../models/permissions_model.dart';
import '../../utils/group.dart';

/// Why Quick Scan cannot run for the current user in the current group.
enum QuickScanBlockedReason {
  /// No Receipt Processing Setting is configured server-side, so the
  /// `aiPoweredReceipts` feature flag is off for the whole install.
  aiDisabled,

  /// The user does not hold `group.receipts.quick-scan` here.
  noPermission,
}

/// What the receipt-entry affordances (the Scan/Add bottom-nav slot, its
/// long-press menu, and the receipts-screen overflow menu) may offer the current
/// user, in the current group.
///
/// This is the single source of truth for that decision: every entry point reads
/// it rather than re-deriving the gates, so the slot's icon, the menu's items and
/// the action a tap performs can never disagree.
///
/// The three gates are independent, and the backend enforces the two permissions
/// separately (`handlers.QuickScan` → `group.receipts.quick-scan`,
/// `handlers.CreateReceipt` → `group.receipts.create`), so a role may hold either
/// without the other.
@immutable
class ReceiptEntryAvailability {
  const ReceiptEntryAvailability({
    required this.canQuickScan,
    required this.canCreateManual,
    required this.blockedReason,
    required this.groupName,
  });

  /// `aiPoweredReceipts` AND `group.receipts.quick-scan`.
  final bool canQuickScan;

  /// `group.receipts.create`.
  final bool canCreateManual;

  /// Why Quick Scan is unavailable; `null` exactly when [canQuickScan].
  final QuickScanBlockedReason? blockedReason;

  /// The name of the single group being acted on, or `null` when there isn't one
  /// (group-select / the all-groups view). Used only for the banner copy.
  final String? groupName;

  /// Whether the entry point should be rendered at all. With neither permission
  /// there is no action the user could take, so the affordance is hidden rather
  /// than shown and then refused.
  bool get isVisible => canQuickScan || canCreateManual;
}

/// Resolves [ReceiptEntryAvailability] for the group the caller is currently in.
///
/// Scoping mirrors what the add menu has always done: inside a group the
/// permissions are checked against THAT group; on the group-select screen and the
/// all-groups view there is no single group, so it falls back to "held in any
/// group" — the receipt form then lets the user pick a group they can act in.
ReceiptEntryAvailability resolveReceiptEntryAvailability(BuildContext context) {
  final authModel = Provider.of<AuthModel>(context, listen: false);
  final groupModel = Provider.of<GroupModel>(context, listen: false);
  final permissionsModel =
      Provider.of<PermissionsModel>(context, listen: false);

  final group = groupModel.getGroupById(getGroupId(context));
  final hasSingleGroup = group != null && !group.isAllGroup;

  bool can(Permission permission) {
    if (!hasSingleGroup) {
      return permissionsModel.hasGroupPermissionInAnyGroup(permission);
    }
    return permissionsModel.hasGroupPermission(group.id, permission);
  }

  final aiEnabled = authModel.featureConfig.aiPoweredReceipts;
  final hasQuickScanPermission =
      can(Permission.groupPeriodReceiptsPeriodQuickScan);
  final canQuickScan = aiEnabled && hasQuickScanPermission;

  // The AI flag wins when both are missing: it is the install-wide, administrator
  // level explanation, and telling a user to chase a permission they'd still not
  // be able to use would send them to the wrong person.
  final QuickScanBlockedReason? blockedReason = canQuickScan
      ? null
      : (!aiEnabled
          ? QuickScanBlockedReason.aiDisabled
          : QuickScanBlockedReason.noPermission);

  return ReceiptEntryAvailability(
    canQuickScan: canQuickScan,
    canCreateManual: can(Permission.groupPeriodReceiptsPeriodCreate),
    blockedReason: blockedReason,
    groupName: hasSingleGroup ? group.name : null,
  );
}
