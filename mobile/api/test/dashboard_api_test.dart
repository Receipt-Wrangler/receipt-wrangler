import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for DashboardApi
void main() {
  final instance = Openapi().getDashboardApi();

  group(DashboardApi, () {
    // Create dashboard
    //
    // This will create a dashboard [SYSTEM USER]
    //
    //Future<Dashboard> createDashboard(UpsertDashboardCommand upsertDashboardCommand) async
    test('test createDashboard', () async {
      // TODO
    });

    // Delete dashboard
    //
    // This will delete a dashboard by id
    //
    //Future<Dashboard> deleteDashboard(int dashboardId) async
    test('test deleteDashboard', () async {
      // TODO
    });

    // Get dashboards for a user by group id
    //
    // This will get a dashboards for a user by group id
    //
    //Future<BuiltList<Dashboard>> getDashboardsForUserByGroupId(String groupId) async
    test('test getDashboardsForUserByGroupId', () async {
      // TODO
    });

    // Update dashboard
    //
    // This will update a dashboard
    //
    //Future<Dashboard> updateDashboard(int dashboardId, UpsertDashboardCommand upsertDashboardCommand) async
    test('test updateDashboard', () async {
      // TODO
    });

  });
}
