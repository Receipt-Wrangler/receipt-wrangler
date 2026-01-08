import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/group_activities.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/group_summary.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/pie_chart.dart';

import 'dashboard_widgets/filtered_receipts.dart';

class GroupDashboard extends StatefulWidget {
  GroupDashboard({super.key, required this.dashboards});

  @override
  State<GroupDashboard> createState() => _GroupDashboard();

  List<api.Dashboard> dashboards = [];
}

class _GroupDashboard extends State<GroupDashboard> {
  int? selectedDashboardIndex;

  void onGroupTap(api.Group group) {
    context.go("/groups/${group.id}");
  }

  void setSelectedDashboardIndex(int index) {
    setState(() {
      selectedDashboardIndex = index;
    });
  }

  Widget buildChoiceChipList(List<api.Dashboard> dashboards) {
    var widgets = <Widget>[];
    for (int i = 0; i < dashboards.length; i++) {
      var dashboard = dashboards[i];
      var defaultSelected = i == 0 && selectedDashboardIndex == null;
      var selected = i == selectedDashboardIndex || defaultSelected;
      var theme = Theme.of(context);

      widgets.add(ChoiceChip(
        key: Key(dashboard.id.toString()),
        label: Text(dashboards[i].name),
        selected: selected,
        selectedColor: theme.primaryColor,
        onSelected: (value) => setSelectedDashboardIndex(i),
      ));
      widgets.add(const SizedBox(width: 10));

      if (selected) {
        selectedDashboardIndex = i;
      }
    }

    return SizedBox(
        height: 50,
        child: ListView(
          scrollDirection: Axis.horizontal,
          children: widgets,
        ));
  }

  List<Widget> buildDashboardWidgets(api.Dashboard? dashboard) {
    var widgets = <Widget>[];

    if (dashboard != null) {
      for (var widget in (dashboard.widgets)?.toList() ?? []) {
        switch (widget.widgetType) {
          case api.WidgetType.FILTERED_RECEIPTS:
            widgets.add(FilteredReceipts(
              dashboardWidget: widget,
            ));
            break;
          case api.WidgetType.GROUP_SUMMARY:
            widgets.add(GroupSummary(
              dashboardWidget: widget,
            ));
            break;
          case api.WidgetType.GROUP_ACTIVITY:
            widgets.add(GroupActivities(
              dashboardWidget: widget,
            ));
            break;
          case api.WidgetType.PIE_CHART:
            widgets.add(DashboardPieChart(
              dashboardWidget: widget,
            ));
            break;
        }
      }
    }

    return widgets;
  }

  api.Dashboard? getSelectedDashboard(List<api.Dashboard>? dashboards) {
    if (selectedDashboardIndex == null) {
      return null;
    } else {
      return dashboards![selectedDashboardIndex!];
    }
  }

  @override
  Widget build(BuildContext context) {
    var chipList = buildChoiceChipList(widget.dashboards);
    api.Dashboard? selectedDashboard = getSelectedDashboard(widget.dashboards);

    List<Widget> children = [];

    if (widget.dashboards.isEmpty) {
      return const Center(child: Text("No dashboards found"));
    } else {
      children.addAll(buildDashboardWidgets(selectedDashboard));
    }

    return SizedBox(
        child: Column(
            mainAxisAlignment: MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
          chipList,
          Expanded(child: GridView.count(crossAxisCount: 1, children: children))
        ]));
  }
}
