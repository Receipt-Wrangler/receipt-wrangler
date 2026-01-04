import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/group_list_card.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';

class GroupList extends StatefulWidget {
  const GroupList({super.key});

  @override
  State<GroupList> createState() => _GroupList();
}

class _GroupList extends State<GroupList> {
  void onGroupTap(api.Group group) {
    context.go("/groups/${group.id}/dashboards");
  }

  @override
  Widget build(BuildContext context) {
    late final groupProvider = Provider.of<GroupModel>(context, listen: true);

    List<Widget> buildGroupCards() {
      var widgets = <Widget>[];
      var activeGroups = groupProvider.groups
          .where((group) => group.status == api.GroupStatus.ACTIVE)
          .toList();

      for (int i = 0; i < activeGroups.length; i++) {
        var group = activeGroups[i];

        widgets.add(GroupListCard(group: group, onGroupTap: onGroupTap));

        if (i < activeGroups.length - 1) {
          widgets.add(const Divider());
        }
      }

      return widgets;
    }

    return SizedBox(child: ListView(children: buildGroupCards()));
  }
}
