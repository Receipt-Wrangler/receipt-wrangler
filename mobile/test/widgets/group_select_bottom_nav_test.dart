import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/groups/nav/group_select/group_select_bottom_nav.dart';

import '../helpers/permission_test_helpers.dart';
import '../helpers/receipt_entry_test_helpers.dart';
import '../helpers/receipt_form_test_helpers.dart';

/// The group-select bottom nav. There is no current group here, so the scan slot
/// falls back to "held in any group" -- and unlike the group nav, only one of
/// its three destinations is unconditional, so it can be gated down to a single
/// entry.
void main() {
  const quickScan = api.Permission.groupPeriodReceiptsPeriodQuickScan;
  const create = api.Permission.groupPeriodReceiptsPeriodCreate;
  const search = api.Permission.appPeriodReceiptsPeriodSearch;

  late List<String> visited;

  setUp(() => visited = []);

  Future<void> pumpNav(
    WidgetTester tester, {
    bool aiEnabled = true,
    List<api.Permission> group = const [],
    List<api.Permission> app = const [],
  }) async {
    GoRoute record(String path, String label) => GoRoute(
          path: path,
          builder: (context, state) {
            visited.add(path);
            return Scaffold(
              body: Text(label),
              bottomNavigationBar: const GroupSelectBottomNav(),
            );
          },
        );

    final router = GoRouter(
      initialLocation: '/groups',
      routes: [
        record('/groups', 'Groups screen'),
        record('/search', 'Search screen'),
        record('/receipts/add', 'Add screen'),
      ],
    );

    await tester.pumpWidget(pumpReceiptEntryApp(
      router: router,
      aiEnabled: aiEnabled,
      permissions: seededPermissions(app: app, group: {5: group}),
      groups: [buildGroup(id: 5, name: 'Household')],
    ));
    await tester.pump();
    visited.clear();
  }

  testWidgets('the scan slot appears when the permission is held in any group',
      (tester) async {
    await pumpNav(tester, group: [quickScan, create]);

    expect(find.text('Scan'), findsOneWidget,
        reason: 'there is no current group here, so the check falls back to '
            '"held in any group"');
  });

  testWidgets('falls back to "Add" with the same rules as the group nav',
      (tester) async {
    await pumpNav(tester, aiEnabled: false, group: [quickScan, create]);

    expect(find.text('Add'), findsOneWidget);
  });

  testWidgets('omits the slot when the permission is held nowhere',
      (tester) async {
    await pumpNav(tester,
        group: [api.Permission.groupPeriodReceiptsPeriodRead], app: [search]);

    expect(find.text('Scan'), findsNothing);
    expect(find.text('Add'), findsNothing);
    expect(find.text('Groups'), findsOneWidget);
    expect(find.text('Search'), findsOneWidget);
  });

  testWidgets('routes Search correctly with the scan slot hidden',
      (tester) async {
    await pumpNav(tester,
        group: [api.Permission.groupPeriodReceiptsPeriodRead], app: [search]);

    await tester.tap(find.text('Search'));
    await tester.pumpAndSettle();

    expect(visited, ['/search']);
  });

  testWidgets('renders no bar at all when only one destination survives',
      (tester) async {
    // Groups is the sole unconditional destination, and this user can neither
    // add receipts anywhere nor search. Material's NavigationBar asserts on a
    // single destination, and a nav pointing only at the screen you are on is
    // useless anyway.
    await pumpNav(tester,
        group: [api.Permission.groupPeriodReceiptsPeriodRead]);

    expect(find.byType(NavigationBar), findsNothing);
    expect(tester.takeException(), isNull);
  });
}
