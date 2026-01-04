import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/group_dashboard_wrapper.dart';

class GroupDashboards extends StatefulWidget {
  const GroupDashboards({super.key});

  @override
  State<GroupDashboards> createState() => _GroupDashboards();
}

class _GroupDashboards extends State<GroupDashboards> {
  @override
  Widget build(BuildContext context) {
    return const GroupDashboardWrapper();
  }
}
