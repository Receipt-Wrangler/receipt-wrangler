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

/// E2E for the view-only Report dashboard widget. Seeds a saved report template
/// and a dashboard that pins it in a `REPORT` widget, drives the app to the
/// group's dashboards route, and asserts the widget renders — plus the two deny
/// paths that prove rejection works for viewers without report perms.
///
/// The render endpoint (`POST /report/template/{id}/render`) **never 403s**: a
/// caller lacking report access gets restricted-notice HTML at 200 with empty
/// `allowedActions`. So the widget always renders *some* HTML, and denial shows
/// up as the Download button being withheld — never an error branch.
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

  /// Seeds a report template (admin jwt) + a dashboard holding a `REPORT` widget
  /// owned by [fixture]'s user (its jwt — `getDashboardsForUserByGroup` filters
  /// on `user_id`), logs that user in, enters the group, and waits for the widget
  /// to finish rendering. Returns the widget's name. Registers teardown for both
  /// seeded rows.
  Future<String> seedAndOpenReportDashboard(
    WidgetTester tester,
    PermFixture fixture,
  ) async {
    final adminJwt = await apiLogin();
    final suffix = DateTime.now().microsecondsSinceEpoch;

    // Template creation needs app.reports.create (admin-only); ownership does not
    // affect the viewer's render access (that's resolved per-request).
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

    // The group role holds group.dashboards.read, so entering the group lands on
    // the dashboards route (no redirect).
    await enterGroup(tester, fixture.groupName!);
    await pumpUntilFound(tester, find.byType(GroupDashboardWrapper));
    await pumpUntilFound(tester, find.byType(ReportWidget));

    // The widget's own spinner clears once the render call resolves and the
    // WebView finishes painting the HTML.
    final spinner = find.descendant(
      of: find.byType(ReportWidget),
      matching: find.byType(CircularProgressIndicator),
    );
    await pumpUntilGone(tester, spinner, timeout: const Duration(seconds: 25));

    return widgetName;
  }

  final downloadButton = find.byKey(const ValueKey('report-widget-download'));

  testWidgets(
    'renders the report and shows the gated download for a full-access user',
    (tester) async {
      // Legacy Owner group role → group.dashboards.read/create + group.reports.read;
      // the app perms grant the app-scoped read + generate the render endpoint's
      // allowedActions resolution requires.
      final fixture = await provisionUserWithAppPermissions(
        const [
          'app.reports.read',
          'app.reports.readAll',
          'app.reports.generate',
          'app.reports.generateAll',
        ],
        groupRoleName: 'Legacy Owner',
      );

      final widgetName = await seedAndOpenReportDashboard(tester, fixture);

      // Success branch: WebView mounted (not the error placeholder), report rendered.
      expect(find.byType(WebViewWidget), findsOneWidget);
      expect(find.text("Couldn't load this report."), findsNothing);
      expect(find.text(widgetName), findsOneWidget);
      // Download shown — allowedActions.generate flowed through from the render.
      expect(downloadButton, findsOneWidget);
    },
    skip: Platform.isLinux, // webview_flutter has no Linux desktop implementation
  );

  testWidgets(
    'hides the download when the user can view but not generate',
    (tester) async {
      // Read (+ Legacy Owner's group.reports.read) but NO generate: the report
      // renders fully, yet the download action must be withheld.
      final fixture = await provisionUserWithAppPermissions(
        const ['app.reports.read', 'app.reports.readAll'],
        groupRoleName: 'Legacy Owner',
      );

      await seedAndOpenReportDashboard(tester, fixture);

      // The real report still renders (not restricted) — read is granted...
      expect(find.byType(WebViewWidget), findsOneWidget);
      expect(find.text("Couldn't load this report."), findsNothing);
      // ...but the Download button is rejected: no generate in allowedActions.
      expect(downloadButton, findsNothing);
    },
    skip: Platform.isLinux,
  );

  testWidgets(
    'renders the restricted notice and no download for a user without report perms',
    (tester) async {
      // Legacy Viewer group role holds group.dashboards.read/create (so the
      // dashboard opens and the user can seed their own) but no group.reports.read,
      // and the default Legacy User app role has no app.reports.* at all.
      final fixture = await provisionPermUser(roleName: 'Legacy Viewer');

      await seedAndOpenReportDashboard(tester, fixture);

      // The render returns restricted-notice HTML at 200, so the widget still
      // renders gracefully (no error branch, no crash)...
      expect(find.byType(WebViewWidget), findsOneWidget);
      expect(find.text("Couldn't load this report."), findsNothing);
      // ...with the Download button hidden (empty allowedActions).
      expect(downloadButton, findsNothing);
    },
    skip: Platform.isLinux,
  );
}
