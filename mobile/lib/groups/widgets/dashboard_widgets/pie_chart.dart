import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/shared/widgets/pie_chart_widget.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';

import '../../../client/client.dart';
import '../constants/text_styles.dart';

class DashboardPieChart extends StatefulWidget {
  const DashboardPieChart({super.key, required this.dashboardWidget});

  final api.Widget dashboardWidget;

  @override
  State<DashboardPieChart> createState() => _DashboardPieChartState();
}

class _DashboardPieChartState extends State<DashboardPieChart> {
  late Future<api.Response<api.PieChartData>> _pieChartFuture;
  bool _isInitialized = false;

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();
    if (!_isInitialized) {
      _loadData();
      _isInitialized = true;
    }
  }

  void _loadData() {
    var groupId = int.tryParse(getGroupId(context) ?? "") ?? 0;
    var config = widget.dashboardWidget.configuration;

    // Extract chart grouping from configuration
    api.ChartGrouping chartGrouping = api.ChartGrouping.CATEGORIES;
    if (config != null && config.containsKey('chartGrouping')) {
      var chartGroupingValue = config['chartGrouping'];
      if (chartGroupingValue != null) {
        try {
          // JsonObject wraps the value, use asString to extract it
          var groupingString = chartGroupingValue.asString;
          chartGrouping = api.ChartGrouping.valueOf(groupingString);
        } catch (_) {
          // Default to categories if parsing fails
        }
      }
    }

    // Build the command
    var command = api.PieChartDataCommand((b) => b
      ..chartGrouping = chartGrouping
    );

    _pieChartFuture = OpenApiClient.client.getWidgetApi().getPieChartData(
      groupId: groupId,
      pieChartDataCommand: command,
    );
  }

  String _getChartGroupingLabel() {
    var config = widget.dashboardWidget.configuration;
    if (config != null && config.containsKey('chartGrouping')) {
      var chartGroupingValue = config['chartGrouping'];
      if (chartGroupingValue != null) {
        try {
          var groupingString = chartGroupingValue.asString;
          switch (groupingString) {
            case 'CATEGORIES':
              return 'By Categories';
            case 'TAGS':
              return 'By Tags';
            case 'PAIDBY':
              return 'By Paid By';
          }
        } catch (_) {
          // Ignore parsing errors
        }
      }
    }
    return 'By Categories';
  }

  @override
  Widget build(BuildContext context) {
    return FutureBuilder<api.Response<api.PieChartData>>(
      future: _pieChartFuture,
      builder: (context, snapshot) {
        bool isLoading = snapshot.connectionState != ConnectionState.done;
        List<PieChartDataPoint> data = [];

        if (snapshot.hasData && snapshot.data?.data != null) {
          data = snapshot.data!.data!.data.map((point) {
            return PieChartDataPoint(
              label: point.label,
              value: point.value,
            );
          }).toList();
        }

        return Column(
          mainAxisAlignment: MainAxisAlignment.start,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            const SizedBox(height: 10),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Text(
                  widget.dashboardWidget.name ?? 'Pie Chart',
                  style: dashboardWidgetNameStyle,
                ),
                Container(
                  padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: Colors.grey[200],
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    _getChartGroupingLabel(),
                    style: const TextStyle(fontSize: 12),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 10),
            Expanded(
              child: PieChartWidget(
                data: data,
                isLoading: isLoading,
                height: 250,
              ),
            ),
          ],
        );
      },
    );
  }
}
