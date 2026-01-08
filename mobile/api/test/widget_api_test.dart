import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for WidgetApi
void main() {
  final instance = Openapi().getWidgetApi();

  group(WidgetApi, () {
    // Get pie chart data
    //
    // This will get pie chart data for a group based on the specified grouping
    //
    //Future<PieChartData> getPieChartData(int groupId, PieChartDataCommand pieChartDataCommand) async
    test('test getPieChartData', () async {
      // TODO
    });

  });
}
