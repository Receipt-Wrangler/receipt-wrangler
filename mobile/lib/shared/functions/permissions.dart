import 'package:openapi/openapi.dart';

import '../../models/permissions_model.dart';

/// Whether the current user may edit receipts in [groupId]. Gates on the modern
/// `group.receipts.update` permission — the legacy GroupRole EDITOR and OWNER
/// tiers both held it, so this preserves the previous behavior.
bool canEditReceipt(PermissionsModel permissionsModel, int groupId) {
  return permissionsModel.hasGroupPermission(
      groupId, Permission.groupPeriodReceiptsPeriodUpdate);
}
