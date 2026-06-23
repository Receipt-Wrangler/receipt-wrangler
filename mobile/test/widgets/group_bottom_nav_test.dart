import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/groups/nav/group/group_bottom_nav.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';

import '../helpers/permission_test_helpers.dart';

/// The group-context bottom nav hides its trailing Search destination when the
/// caller lacks `app.receipts.search`. A real `GoRouter` ancestor is required
/// because `setIndexSelected` reads `GoRouter.of(context)` during build.
void main() {
  Widget harness(PermissionsModel permissions) {
    final router = GoRouter(
      initialLocation: '/groups/1/receipts',
      routes: [
        GoRoute(
          path: '/groups/:groupId/receipts',
          builder: (context, state) =>
              const Scaffold(bottomNavigationBar: GroupBottomNav()),
        ),
      ],
    );
    return ChangeNotifierProvider<PermissionsModel>.value(
      value: permissions,
      child: MaterialApp.router(routerConfig: router),
    );
  }

  testWidgets('hides the Search destination without app.receipts.search',
      (tester) async {
    await tester.pumpWidget(harness(seededPermissions()));
    await tester.pump();

    expect(find.byIcon(Icons.search), findsNothing);
    expect(find.text('Search'), findsNothing);
    // The other destinations remain.
    expect(find.byIcon(Icons.dashboard), findsOneWidget);
    expect(find.byIcon(Icons.add), findsOneWidget);
    expect(find.byIcon(Icons.receipt), findsOneWidget);
  });

  testWidgets('shows the Search destination with app.receipts.search',
      (tester) async {
    await tester.pumpWidget(harness(
      seededPermissions(app: [Permission.appPeriodReceiptsPeriodSearch]),
    ));
    await tester.pump();

    expect(find.byIcon(Icons.search), findsOneWidget);
    expect(find.text('Search'), findsOneWidget);
  });
}
