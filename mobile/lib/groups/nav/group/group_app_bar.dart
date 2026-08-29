import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_entry_overflow_menu.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/top_app_bar.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';

import '../../../models/group_model.dart';

/// The group shell's app bar, which covers both the dashboards and the receipts
/// routes.
const _receiptsRoutePath = '/groups/:groupId/receipts';

class GroupAppBar extends StatefulWidget implements PreferredSizeWidget {
  const GroupAppBar({super.key});

  @override
  State<GroupAppBar> createState() => _GroupAppBar();

  @override
  Size get preferredSize => AppBar().preferredSize;
}

class _GroupAppBar extends State<GroupAppBar> {
  String getGroupTitleText(api.Group group) {
    if (group.name.toLowerCase().contains("receipt")) {
      return group.name;
    }

    return "${group.name} Receipts";
  }

  @override
  Widget build(BuildContext context) {
    final groupId = getGroupId(context);
    final group =
        Provider.of<GroupModel>(context, listen: false).getGroupById(groupId);

    // The receipt-entry overflow belongs to the receipts screen only. This app
    // bar is shared with the dashboards route, so gate on the route itself --
    // the same check main.dart uses to pick that route's padding.
    final isReceiptsRoute =
        GoRouterState.of(context).fullPath == _receiptsRoutePath;

    return TopAppBar(
      titleText: group == null ? 'Receipts' : getGroupTitleText(group),
      leadingArrowRedirect: "/groups",
      actions: isReceiptsRoute ? const [ReceiptEntryOverflowMenu()] : null,
    );
  }
}
