import 'dart:io' show Platform;

import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/group_dashboard_wrapper.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/group_receipts_list.dart';

import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

/// Permission-gating coverage for the group dashboards route, gated on
/// `group.dashboards.read` by `groupDashboardReadRedirect` (wired on
/// `/groups/:groupId/dashboards` in `main.dart`). Selecting a group always lands
/// on that route; a member without the permission is redirected to the group's
/// receipt list instead, mirroring the desktop `groupDashboardReadGuard`.
///
/// Landing is told apart by screen-unique widgets: `GroupReceiptsList` (receipts
/// screen) vs `GroupDashboardWrapper` (dashboards screen) — `GroupReceiptsList`
/// is used nowhere else. The deny case needs a custom group role (every seeded
/// group role holds `group.dashboards.read`), provisioned + torn down by
/// `provisionGroupMemberWithoutPermission`.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets(
    'dashboards: a member without group.dashboards.read is redirected to receipts',
    (tester) async {
      final fixture =
          await provisionGroupMemberWithoutPermission('group.dashboards.read');
      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );

      // Selecting the group routes to /groups/:id/dashboards, which the redirect
      // bounces to /groups/:id/receipts (the role keeps group.receipts.read, so
      // the receipts list loads; the dashboards screen never builds).
      await enterGroup(tester, fixture.groupName!);
      await pumpUntilFound(tester, find.byType(GroupReceiptsList));

      expect(find.byType(GroupReceiptsList), findsOneWidget);
      expect(find.byType(GroupDashboardWrapper), findsNothing);
    },
  );

  testWidgets(
    'dashboards: a Legacy Viewer (holds group.dashboards.read) sees the dashboard',
    (tester) async {
      final fixture = await provisionPermUser(roleName: 'Legacy Viewer');
      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );

      // Legacy Viewer holds group.dashboards.read, so no redirect fires.
      await enterGroup(tester, fixture.groupName!);
      await pumpUntilFound(tester, find.byType(GroupDashboardWrapper));

      expect(find.byType(GroupDashboardWrapper), findsOneWidget);
      expect(find.byType(GroupReceiptsList), findsNothing);
    },
  );
}
