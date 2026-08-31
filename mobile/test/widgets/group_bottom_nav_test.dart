import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/groups/nav/group/group_bottom_nav.dart';

import '../helpers/permission_test_helpers.dart';
import '../helpers/receipt_entry_test_helpers.dart';
import '../helpers/receipt_form_test_helpers.dart';

/// The group-context bottom nav. Two of its destinations are permission-gated --
/// the scan slot and the trailing Search -- so positions are not stable, and the
/// nav resolves taps by destination id rather than by index. These cover that
/// the shared scan slot lands here correctly and that a hidden slot doesn't
/// misroute the destinations after it.
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
    Widget screen(String label) => Scaffold(
          body: Text(label),
          bottomNavigationBar: const GroupBottomNav(),
        );

    GoRoute record(String path, String label) => GoRoute(
          path: path,
          builder: (context, state) {
            visited.add(path);
            return screen(label);
          },
        );

    final router = GoRouter(
      initialLocation: '/groups/5/receipts',
      routes: [
        record('/groups/:groupId/receipts', 'Receipts screen'),
        record('/groups/:groupId/dashboards', 'Dashboards screen'),
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

  testWidgets('mounts the shared scan slot as "Scan" when Quick Scan can run',
      (tester) async {
    await pumpNav(tester, group: [quickScan, create]);

    expect(find.text('Scan'), findsOneWidget);
    expect(find.byIcon(Icons.document_scanner), findsOneWidget);
  });

  testWidgets('mounts it as "Add" when Quick Scan cannot', (tester) async {
    await pumpNav(tester, aiEnabled: false, group: [quickScan, create]);

    expect(find.text('Add'), findsOneWidget);
    expect(find.byIcon(Icons.add), findsOneWidget);
  });

  testWidgets('omits the slot when the user can neither scan nor create',
      (tester) async {
    await pumpNav(tester,
        group: [api.Permission.groupPeriodReceiptsPeriodRead]);

    expect(find.text('Scan'), findsNothing);
    expect(find.text('Add'), findsNothing);
    // The unconditional destinations remain.
    expect(find.byIcon(Icons.dashboard), findsOneWidget);
    expect(find.byIcon(Icons.receipt), findsOneWidget);
  });

  testWidgets('hides the Search destination without app.receipts.search',
      (tester) async {
    await pumpNav(tester, group: [create]);

    expect(find.byIcon(Icons.search), findsNothing);
    expect(find.text('Search'), findsNothing);
  });

  testWidgets('shows the Search destination with app.receipts.search',
      (tester) async {
    await pumpNav(tester, group: [create], app: [search]);

    expect(find.text('Search'), findsOneWidget);
  });

  testWidgets('routes by identity when the scan slot is present',
      (tester) async {
    await pumpNav(tester, group: [create], app: [search]);

    await tester.tap(find.text('Search'));
    await tester.pumpAndSettle();

    expect(visited, ['/search']);
  });

  testWidgets('a hidden scan slot does not misroute the destinations after it',
      (tester) async {
    // Without the id lookup, Search would sit at the index the scan slot used
    // to occupy and a tap would land on the wrong screen.
    await pumpNav(tester,
        group: [api.Permission.groupPeriodReceiptsPeriodRead], app: [search]);

    await tester.tap(find.text('Search'));
    await tester.pumpAndSettle();
    expect(visited, ['/search']);

    visited.clear();
    await tester.tap(find.text('Dashboards'));
    await tester.pumpAndSettle();
    expect(visited, ['/groups/:groupId/dashboards']);
  });

  testWidgets('highlights the destination matching the current route',
      (tester) async {
    await pumpNav(tester, group: [create], app: [search]);

    final bar = tester.widget<NavigationBar>(find.byType(NavigationBar));
    expect(bar.destinations[bar.selectedIndex].key, isNull);
    expect(
      (bar.destinations[bar.selectedIndex] as NavigationDestination).label,
      'Receipts',
      reason: 'the initial location is the receipts route',
    );
  });
}
