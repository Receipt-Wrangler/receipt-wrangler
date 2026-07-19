import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/top_app_bar.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/user_avatar.dart';

import 'helpers/api.dart';
import 'helpers/login.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

/// E2E coverage for the Reports **list** page (`/reports`, reached from the avatar
/// menu). It verifies navigation + list rendering and the app-scoped menu gate:
///
///  - a user with `app.reports.readAll` sees the Reports menu, navigates, and a
///    seeded template row renders. This is the regression guard for the
///    dimension-column deserialization fix: the API used to emit `"aggFunc":""`
///    for dimension columns, which the generated Dart enum could not deserialize,
///    so the row collapsed to an invisible blank `SizedBox`. If that regresses,
///    the row key never appears and this test fails.
///  - a user with only base `app.reports.read` and no visible templates gets the
///    empty state ("No reports found").
///  - a Legacy User (no report permissions) never sees the Reports menu entry.
///
/// Permission enforcement is server-side; these tests assert the UI reflects it.
/// Scope is the list view only — preview/generate/delete are not exercised here.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  // Opens the avatar popup menu from the post-login GroupSelect landing. The
  // popup mounts its items on the animation's first frame, so we wait on a stable
  // item (User Profile) being hit-testable and drain a few frames before acting —
  // the popup-menu tap-flake pattern documented in mobile/CLAUDE.md.
  Future<void> openAvatarMenu(WidgetTester tester) async {
    final avatar = find.descendant(
      of: find.byType(TopAppBar),
      matching: find.byType(UserAvatar),
    );
    await pumpUntilFound(tester, avatar);
    await tester.tap(avatar);
    await pumpUntilFound(tester, find.text('User Profile').hitTestable());
    for (var i = 0; i < 5; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }
  }

  Future<void> tapReports(WidgetTester tester) async {
    final reports = find.text('Reports').hitTestable();
    await pumpUntilFound(tester, reports);
    for (var i = 0; i < 5; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }
    await tester.tap(reports);
  }

  testWidgets('list renders a seeded report template row', (tester) async {
    final jwt = await apiLogin();
    final fixture = await provisionUserWithAppPermissions(
      const ['app.reports.read', 'app.reports.readAll'],
      groupRoleName: 'Legacy Owner',
    );
    final templateName = 'e2e-report-${DateTime.now().microsecondsSinceEpoch}';
    final templateId = await createReportTemplate(
      groupIds: [fixture.groupId!],
      jwt: jwt,
      name: templateName,
    );
    addTearDown(
        () async => deleteReportTemplate(templateId, jwt: await apiLogin()));

    await loginAs(tester,
        username: fixture.username, password: fixture.password);

    await openAvatarMenu(tester);
    await tapReports(tester);

    // The seeded row renders — proves the anyOf unwrap and the aggFunc
    // omitempty fix (a dimension column no longer breaks deserialization).
    await pumpUntilFound(tester, find.byKey(ValueKey('report-$templateId')));
    expect(find.text(templateName), findsOneWidget);
    // readAll grants every action, so allowedActions includes 'read' → the
    // preview (eye) action is present on the row.
    expect(find.byKey(ValueKey('report-preview-$templateId')), findsOneWidget);
  });

  testWidgets('empty state when no templates are visible', (tester) async {
    final fixture =
        await provisionUserWithAppPermissions(const ['app.reports.read']);

    await loginAs(tester,
        username: fixture.username, password: fixture.password);

    await openAvatarMenu(tester);
    await tapReports(tester);

    // Base app.reports.read (no readAll, no group.reports.read, no group) → the
    // page loads but no template is visible.
    await pumpUntilFound(tester, find.text('No reports found'));
    expect(find.text('No reports found'), findsOneWidget);
  });

  testWidgets('Reports menu entry hidden without app.reports.read',
      (tester) async {
    final fixture = await provisionPermUser(); // Legacy User, no report perms

    await loginAs(tester,
        username: fixture.username, password: fixture.password);

    await openAvatarMenu(tester);

    expect(find.text('Reports'), findsNothing);
    // The menu did open — the ungated items are present.
    expect(find.text('User Profile'), findsOneWidget);
    expect(find.text('Logout'), findsOneWidget);
  });
}
