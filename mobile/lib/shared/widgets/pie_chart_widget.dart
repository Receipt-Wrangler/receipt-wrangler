import 'package:fl_chart/fl_chart.dart';
import 'package:flutter/material.dart';

/// A reusable pie chart widget that can display data with customizable styling.
class PieChartWidget extends StatelessWidget {
  const PieChartWidget({
    super.key,
    required this.data,
    this.height = 300,
    this.isLoading = false,
    this.noDataMessage = 'No data available',
    this.loadingMessage = 'Loading...',
    this.formatValue,
  });

  /// List of data points to display in the chart
  final List<PieChartDataPoint> data;

  /// Height of the chart container
  final double height;

  /// Whether the chart is in loading state
  final bool isLoading;

  /// Message to display when there is no data
  final String noDataMessage;

  /// Message to display while loading
  final String loadingMessage;

  /// Optional callback to format the value for display in the legend
  final String Function(double value)? formatValue;

  /// Default colors for the pie chart slices
  static const List<Color> defaultColors = [
    Color(0xFF2196F3), // Blue
    Color(0xFF4CAF50), // Green
    Color(0xFFF44336), // Red
    Color(0xFFFF9800), // Orange
    Color(0xFF9C27B0), // Purple
    Color(0xFF00BCD4), // Cyan
    Color(0xFFFFEB3B), // Yellow
    Color(0xFF795548), // Brown
    Color(0xFF607D8B), // Blue Grey
    Color(0xFFE91E63), // Pink
    Color(0xFF3F51B5), // Indigo
    Color(0xFF009688), // Teal
  ];

  @override
  Widget build(BuildContext context) {
    if (isLoading) {
      return SizedBox(
        height: height,
        child: Center(
          child: Column(
            mainAxisAlignment: MainAxisAlignment.center,
            children: [
              const CircularProgressIndicator(),
              const SizedBox(height: 16),
              Text(loadingMessage),
            ],
          ),
        ),
      );
    }

    if (data.isEmpty) {
      return SizedBox(
        height: height,
        child: Center(
          child: Text(
            noDataMessage,
            style: Theme.of(context).textTheme.bodyLarge?.copyWith(
                  color: Colors.grey,
                ),
          ),
        ),
      );
    }

    return SizedBox(
      height: height,
      child: Row(
        children: [
          Expanded(
            flex: 2,
            child: PieChart(
              PieChartData(
                sections: _buildSections(),
                sectionsSpace: 2,
                centerSpaceRadius: 40,
                pieTouchData: PieTouchData(enabled: false),
              ),
            ),
          ),
          const SizedBox(width: 16),
          Expanded(
            flex: 1,
            child: _buildLegend(context),
          ),
        ],
      ),
    );
  }

  List<PieChartSectionData> _buildSections() {
    final total = data.fold<double>(0, (sum, item) => sum + item.value);

    return data.asMap().entries.map((entry) {
      final index = entry.key;
      final item = entry.value;
      final percentage = total > 0 ? (item.value / total * 100) : 0;
      final color = defaultColors[index % defaultColors.length];

      return PieChartSectionData(
        value: item.value,
        title: '${percentage.toStringAsFixed(1)}%',
        color: color,
        radius: 80,
        titleStyle: const TextStyle(
          fontSize: 12,
          fontWeight: FontWeight.bold,
          color: Colors.white,
        ),
      );
    }).toList();
  }

  Widget _buildLegend(BuildContext context) {
    return SingleChildScrollView(
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisAlignment: MainAxisAlignment.center,
        children: data.asMap().entries.map((entry) {
          final index = entry.key;
          final item = entry.value;
          final color = defaultColors[index % defaultColors.length];
          final formattedValue = formatValue != null ? formatValue!(item.value) : null;

          return Padding(
            padding: const EdgeInsets.symmetric(vertical: 4),
            child: Row(
              children: [
                Container(
                  width: 12,
                  height: 12,
                  decoration: BoxDecoration(
                    color: color,
                    shape: BoxShape.circle,
                  ),
                ),
                const SizedBox(width: 8),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        item.label,
                        style: Theme.of(context).textTheme.bodySmall,
                        overflow: TextOverflow.ellipsis,
                      ),
                      if (formattedValue != null)
                        Text(
                          formattedValue,
                          style: Theme.of(context).textTheme.bodySmall?.copyWith(
                                color: Colors.grey[600],
                                fontSize: 10,
                              ),
                        ),
                    ],
                  ),
                ),
              ],
            ),
          );
        }).toList(),
      ),
    );
  }
}

/// A data point for the pie chart
class PieChartDataPoint {
  const PieChartDataPoint({
    required this.label,
    required this.value,
  });

  final String label;
  final double value;
}
