import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/shared/functions/quick_scan.dart';

import '../../models/auth_model.dart';
import '../../models/group_model.dart';
import '../../models/permissions_model.dart';
import '../../utils/group.dart';
import '../../utils/snackbar.dart';

void showAddMenu(BuildContext context, GlobalKey addButtonKey) {
  final RenderBox renderBox =
      addButtonKey.currentContext?.findRenderObject() as RenderBox;
  final Offset offset = renderBox.localToGlobal(Offset.zero);
  final Size size = renderBox.size;

  final RelativeRect position = RelativeRect.fromLTRB(
    offset.dx,
    offset.dy,
    offset.dx + size.width,
    offset.dy + size.height,
  );

  final authModel = Provider.of<AuthModel>(context, listen: false);
  final groupModel = Provider.of<GroupModel>(context, listen: false);
  final permissionsModel =
      Provider.of<PermissionsModel>(context, listen: false);

  // Resolve the group being acted on. On the group-select screen and the
  // all-groups view there is no single group, so fall back to "held in any
  // group" — the receipt add form lets the user pick a group they can create in.
  final group = groupModel.getGroupById(getGroupId(context));
  bool can(Permission permission) {
    if (group == null || group.isAllGroup) {
      return permissionsModel.hasGroupPermissionInAnyGroup(permission);
    }
    return permissionsModel.hasGroupPermission(group.id, permission);
  }

  final items = <PopupMenuItem>[
    if (can(Permission.groupPeriodReceiptsPeriodCreate))
      PopupMenuItem(
        value: 0,
        child: const Text("Add Manual Receipt"),
        onTap: () => context.go("/receipts/add"),
      ),
    if (authModel.featureConfig.aiPoweredReceipts &&
        can(Permission.groupPeriodReceiptsPeriodQuickScan))
      PopupMenuItem(
        value: 1,
        child: const Text("Quick Scan"),
        onTap: () => showQuickScanBottomSheet(context),
      ),
  ];

  if (items.isEmpty) {
    showErrorSnackbar(
        context, "You don't have permission to add receipts here.");
    return;
  }

  showMenu(context: context, position: position, items: items);
}
