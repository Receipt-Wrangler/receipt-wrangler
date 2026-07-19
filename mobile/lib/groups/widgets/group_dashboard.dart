import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/group_activities.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/group_summary.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/pie_chart.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/report_widget.dart';

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
    var effectiveIndex = selectedDashboardIndex ?? 0;

    for (int i = 0; i < dashboards.length; i++) {
      var dashboard = dashboards[i];
      var selected = i == effectiveIndex;
      var theme = Theme.of(context);

      widgets.add(ChoiceChip(
        key: Key(dashboard.id.toString()),
        label: Text(dashboards[i].name),
        selected: selected,
        selectedColor: theme.primaryColor,
        onSelected: (value) => setSelectedDashboardIndex(i),
      ));
      widgets.add(const SizedBox(width: 10));
    }

    return SizedBox(
        height: 50,
        child: ListView(
          scrollDirection: Axis.horizontal,
          children: widgets,
        ));
  }

  List<Widget> buildDashboardWidgets(
      api.Dashboard? dashboard, double widgetHeight) {
    var widgets = <Widget>[];

    if (dashboard != null) {
      for (var widget in (dashboard.widgets)?.toList() ?? []) {
        switch (widget.widgetType) {
          case api.WidgetType.FILTERED_RECEIPTS:
            widgets.add(SizedBox(
              height: widgetHeight,
              child: FilteredReceipts(
                dashboardWidget: widget,
              ),
            ));
            break;
          case api.WidgetType.GROUP_SUMMARY:
            widgets.add(GroupSummary(
              dashboardWidget: widget,
            ));
            break;
          case api.WidgetType.GROUP_ACTIVITY:
            widgets.add(SizedBox(
              height: widgetHeight,
              child: GroupActivities(
                dashboardWidget: widget,
              ),
            ));
            break;
          case api.WidgetType.PIE_CHART:
            widgets.add(SizedBox(
              height: widgetHeight,
              child: DashboardPieChart(
                dashboardWidget: widget,
              ),
            ));
            break;
          case api.WidgetType.REPORT:
            widgets.add(SizedBox(
              height: widgetHeight,
              child: ReportWidget(
                // Key by widget id so switching to a dashboard with a different
                // report in this slot remounts (re-fetches) instead of reusing
                // the previous report's cached HTML + download gate.
                key: ValueKey(widget.id),
                dashboardWidget: widget,
              ),
            ));
            break;
        }
      }
    }

    return widgets;
  }

  api.Dashboard? getSelectedDashboard(List<api.Dashboard>? dashboards) {
    if (dashboards == null || dashboards.isEmpty) {
      return null;
    }
    var index = (selectedDashboardIndex ?? 0).clamp(0, dashboards.length - 1);
    return dashboards[index];
  }

  @override
  Widget build(BuildContext context) {
    if (widget.dashboards.isEmpty) {
      return const Center(child: Text("No dashboards found"));
    }

    var chipList = buildChoiceChipList(widget.dashboards);
    api.Dashboard? selectedDashboard = getSelectedDashboard(widget.dashboards);
    var widgetHeight = MediaQuery.of(context).size.height * 0.6;
    List<Widget> children =
        buildDashboardWidgets(selectedDashboard, widgetHeight);

    return Column(
        mainAxisAlignment: MainAxisAlignment.start,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          chipList,
          Expanded(child: ListView(children: children))
        ]);
  }
}
