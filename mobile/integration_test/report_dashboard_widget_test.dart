import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/report_widget.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/group_dashboard_wrapper.dart';
import 'package:webview_flutter/webview_flutter.dart';

import 'helpers/api.dart';
import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

/// E2E for the view-only Report dashboard widget: seeds a saved report template
/// and a dashboard that pins it in a `REPORT` widget, then drives the app to the
/// group's dashboards route and asserts the widget actually renders (WebView
/// success branch, no error placeholder) with its server-gated Download button.
///
/// **iOS/Android only** (`skip: Platform.isLinux`): `webview_flutter` has no Linux
/// desktop implementation, so the widget can't mount under the Linux runner.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets(
    'report dashboard widget renders the report and shows the gated download',
    (tester) async {
      final adminJwt = await apiLogin();

      // A user that can both open the dashboard and fully render a report: the
      // Legacy Owner group role grants group.dashboards.read/create +
      // group.reports.read; the app perms grant the app-scoped read/generate the
      // render endpoint's allowedActions resolution requires.
      final fixture = await provisionUserWithAppPermissions(
        const [
          'app.reports.read',
          'app.reports.readAll',
          'app.reports.generate',
          'app.reports.generateAll',
        ],
        groupRoleName: 'Legacy Owner',
      );

      final suffix = DateTime.now().microsecondsSinceEpoch;

      // Template creation needs app.reports.create (admin-only); ownership does
      // not affect the fixture user's render access (they read via *All).
      final templateId = await createReportTemplate(
        groupIds: [fixture.groupId!],
        jwt: adminJwt,
        name: 'e2e-report-$suffix',
      );
      addTearDown(
          () async => deleteReportTemplate(templateId, jwt: await apiLogin()));

      // The dashboard must be OWNED by the viewing user — getDashboardsForUserByGroup
      // filters on user_id — so seed it with the fixture user's own jwt.
      final userJwt = await apiLoginAs(fixture.username, fixture.password);
      final widgetName = 'e2e-report-widget-$suffix';
      final dashboardId = await createDashboard(
        groupId: fixture.groupId!,
        reportTemplateId: templateId,
        jwt: userJwt,
        name: 'e2e-dashboard-$suffix',
        widgetName: widgetName,
      );
      addTearDown(() async => deleteDashboard(
            dashboardId,
            jwt: await apiLoginAs(fixture.username, fixture.password),
          ));

      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );

      // Legacy Owner holds group.dashboards.read, so entering the group lands on
      // the dashboards route.
      await enterGroup(tester, fixture.groupName!);
      await pumpUntilFound(tester, find.byType(GroupDashboardWrapper));
      await pumpUntilFound(tester, find.byType(ReportWidget));

      // The widget's own spinner clears once the render call resolves and the
      // WebView finishes painting the HTML.
      final spinner = find.descendant(
        of: find.byType(ReportWidget),
        matching: find.byType(CircularProgressIndicator),
      );
      await pumpUntilGone(
        tester,
        spinner,
        timeout: const Duration(seconds: 25),
      );

      // Success branch — the WebView is mounted (not the error placeholder) and
      // the report rendered without error.
      expect(find.byType(WebViewWidget), findsOneWidget);
      expect(find.text("Couldn't load this report."), findsNothing);
      expect(find.text(widgetName), findsOneWidget);

      // The server-gated Download button is shown, proving allowedActions.generate
      // flowed through from the render response.
      expect(
        find.byKey(const ValueKey('report-widget-download')),
        findsOneWidget,
      );
    },
    skip: Platform.isLinux, // webview_flutter has no Linux desktop implementation
  );
}
