import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/top_app_bar.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';

import '../../../models/group_model.dart';

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
    var groupId = getGroupId(context) ?? "";
    var group =
        Provider.of<GroupModel>(context, listen: false).getGroupById(groupId);

    return TopAppBar(
      titleText: getGroupTitleText(group as api.Group),
      leadingArrowRedirect: "/groups",
    );
  }
}
