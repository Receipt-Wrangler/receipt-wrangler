import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/group_dashboard.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/circular_loading_progress.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';

import '../../client/client.dart';

class GroupDashboardWrapper extends StatefulWidget {
  const GroupDashboardWrapper({super.key});

  @override
  State<GroupDashboardWrapper> createState() => _GroupDashboardWrapper();
}

class _GroupDashboardWrapper extends State<GroupDashboardWrapper> {
  @override
  Widget build(BuildContext context) {
    var groupId = getGroupId(context);
    var dashboardFuture = OpenApiClient.client
        .getDashboardApi()
        .getDashboardsForUserByGroupId(groupId: groupId ?? "");

    return FutureBuilder(
        future: dashboardFuture,
        builder: (context, snapshot) {
          if (snapshot.connectionState == ConnectionState.done) {
            return GroupDashboard(
              dashboards: snapshot.data?.data?.toList() ?? [],
            );
          }

          return const CircularLoadingProgress();
        });
  }
}
