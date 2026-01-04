import 'package:built_collection/built_collection.dart';
import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/constants/text_styles.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/group_activity_list_item.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/paged_data_list.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';

import '../../../client/client.dart';
import '../../../models/group_model.dart';

class GroupActivities extends StatefulWidget {
  const GroupActivities({super.key, required api.Widget this.dashboardWidget});

  final api.Widget dashboardWidget;

  @override
  State<GroupActivities> createState() => _GroupActivities();
}

class _GroupActivities extends State<GroupActivities> {
  late final groupId = getGroupId(context);
  late final groupModel = Provider.of<GroupModel>(context, listen: false);

  api.PagedActivityRequestCommandBuilder buildCommand(int page) {
    var groupIds = ListBuilder<int>([int.parse(groupId)]);
    var group = groupModel.getGroupById(groupId);

    if (group != null && group.isAllGroup) {
      groupIds = ListBuilder(groupModel.groupsWithoutAllGroup.map((g) => g.id));
    }

    return api.PagedActivityRequestCommandBuilder()
      ..page = page
      ..pageSize = 25
      ..orderBy = "started_at"
      ..sortDirection = api.SortDirection.desc
      ..groupIds = groupIds;
  }

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      child: Column(
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            widget.dashboardWidget.name ?? "",
            style: dashboardWidgetNameStyle,
          ),
          PagedDataList(
              noItemsFoundText: "No activity found",
              listItemBuilder: (context, activity, index) {
                return GroupActivityListItem(
                    activity: activity.anyOf.values[9] as api.Activity,
                    groupId: int.parse(groupId));
              },
              getPagedDataFuture: (page) {
                return OpenApiClient.client
                    .getSystemTaskApi()
                    .getPagedActivities(
                      pagedActivityRequestCommand: buildCommand(page).build(),
                    );
              }),
          SizedBox(height: 50),
        ],
      ),
    );
  }
}
