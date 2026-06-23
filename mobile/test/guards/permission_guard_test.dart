import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/search.dart';
import 'package:receipt_wrangler_mobile/guards/permission-guard.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';

import '../helpers/permission_test_helpers.dart';

/// Unit coverage for the go_router permission redirects. Each test mounts a
/// minimal router (provider-backed) at a neutral start route, then navigates so
/// the redirect runs with the navigator mounted — exactly as it does in the app.
void main() {
  Widget harness(PermissionsModel permissions, GoRouter router) {
    return ChangeNotifierProvider<PermissionsModel>.value(
      value: permissions,
      child: MaterialApp.router(routerConfig: router),
    );
  }

  GoRouter dashboardRouter() => GoRouter(
        initialLocation: '/start',
        routes: [
          GoRoute(path: '/start', builder: (c, s) => const Text('START')),
          GoRoute(
            path: '/groups/:groupId/dashboards',
            redirect: groupDashboardReadRedirect,
            builder: (c, s) => const Text('DASHBOARDS'),
          ),
          GoRoute(
            path: '/groups/:groupId/receipts',
            builder: (c, s) => const Text('RECEIPTS'),
          ),
        ],
      );

  GoRouter searchRouter() => GoRouter(
        initialLocation: '/start',
        routes: [
          GoRoute(path: '/start', builder: (c, s) => const Text('START')),
          GoRoute(
            path: '/search',
            redirect: receiptsSearchRedirect,
            builder: (c, s) => const Text('SEARCH'),
          ),
          GoRoute(path: '/groups', builder: (c, s) => const Text('GROUPS')),
          GoRoute(
            path: '/groups/:groupId/receipts',
            builder: (c, s) => const Text('RECEIPTS'),
          ),
        ],
      );

  group('groupDashboardReadRedirect', () {
    testWidgets('redirects to receipts without group.dashboards.read',
        (tester) async {
      final router = dashboardRouter();
      await tester.pumpWidget(harness(seededPermissions(), router));
      await tester.pumpAndSettle();

      router.go('/groups/5/dashboards');
      await tester.pumpAndSettle();

      expect(find.text('RECEIPTS'), findsOneWidget);
      expect(find.text('DASHBOARDS'), findsNothing);
    });

    testWidgets('allows dashboards with group.dashboards.read',
        (tester) async {
      final router = dashboardRouter();
      await tester.pumpWidget(harness(
        seededPermissions(
          group: {
            5: [Permission.groupPeriodDashboardsPeriodRead]
          },
        ),
        router,
      ));
      await tester.pumpAndSettle();

      router.go('/groups/5/dashboards');
      await tester.pumpAndSettle();

      expect(find.text('DASHBOARDS'), findsOneWidget);
      expect(find.text('RECEIPTS'), findsNothing);
    });

    testWidgets('redirect is group-scoped — read in another group does not allow',
        (tester) async {
      final router = dashboardRouter();
      await tester.pumpWidget(harness(
        seededPermissions(
          group: {
            9: [Permission.groupPeriodDashboardsPeriodRead]
          },
        ),
        router,
      ));
      await tester.pumpAndSettle();

      router.go('/groups/5/dashboards');
      await tester.pumpAndSettle();

      expect(find.text('RECEIPTS'), findsOneWidget);
      expect(find.text('DASHBOARDS'), findsNothing);
    });
  });

  group('receiptsSearchRedirect', () {
    testWidgets('redirects to /groups without app.receipts.search',
        (tester) async {
      final router = searchRouter();
      await tester.pumpWidget(harness(seededPermissions(), router));
      await tester.pumpAndSettle();

      router.go('/search');
      await tester.pumpAndSettle();

      expect(find.text('GROUPS'), findsOneWidget);
      expect(find.text('SEARCH'), findsNothing);
    });

    testWidgets('redirects to the group receipts when extra carries a groupId',
        (tester) async {
      final router = searchRouter();
      await tester.pumpWidget(harness(seededPermissions(), router));
      await tester.pumpAndSettle();

      router.go('/search',
          extra: {'from': fromGroupBottomNav, 'groupId': '3'});
      await tester.pumpAndSettle();

      expect(find.text('RECEIPTS'), findsOneWidget);
      expect(find.text('SEARCH'), findsNothing);
    });

    testWidgets('allows search with app.receipts.search', (tester) async {
      final router = searchRouter();
      await tester.pumpWidget(harness(
        seededPermissions(app: [Permission.appPeriodReceiptsPeriodSearch]),
        router,
      ));
      await tester.pumpAndSettle();

      router.go('/search');
      await tester.pumpAndSettle();

      expect(find.text('SEARCH'), findsOneWidget);
      expect(find.text('GROUPS'), findsNothing);
    });
  });
}
