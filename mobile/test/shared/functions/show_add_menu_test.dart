import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/shared/functions/show_add_menu.dart';

import '../../helpers/permission_test_helpers.dart';

PermissionsModel _permissions(List<api.Permission> group5) =>
    seededPermissions(group: {5: group5});

AuthModel _auth({required bool aiPoweredReceipts}) {
  final model = AuthModel();
  model.setFeatureConfig((api.FeatureConfigBuilder()
        ..aiPoweredReceipts = aiPoweredReceipts
        ..enableLocalSignUp = false)
      .build());
  return model;
}

// A minimal router app whose single screen has an "add" button wired to
// showAddMenu. With no current group (GroupModel empty, route '/'), the menu
// falls back to the "held in any group" check — the group-select behavior.
Widget _app({required AuthModel auth, required PermissionsModel permissions}) {
  final addKey = GlobalKey();
  final router = GoRouter(
    routes: [
      GoRoute(
        path: '/',
        builder: (context, state) => Scaffold(
          body: Center(
            child: ElevatedButton(
              key: addKey,
              onPressed: () => showAddMenu(context, addKey),
              child: const Text('add'),
            ),
          ),
        ),
      ),
    ],
  );

  return MultiProvider(
    providers: [
      ChangeNotifierProvider<AuthModel>.value(value: auth),
      ChangeNotifierProvider<GroupModel>(create: (_) => GroupModel()),
      ChangeNotifierProvider<PermissionsModel>.value(value: permissions),
    ],
    child: MaterialApp.router(routerConfig: router),
  );
}

void main() {
  group('showAddMenu permission gating', () {
    testWidgets('shows Add Manual Receipt when the user can create in a group',
        (tester) async {
      await tester.pumpWidget(_app(
        auth: _auth(aiPoweredReceipts: false),
        permissions:
            _permissions([api.Permission.groupPeriodReceiptsPeriodCreate]),
      ));

      await tester.tap(find.text('add'));
      await tester.pumpAndSettle();

      expect(find.text('Add Manual Receipt'), findsOneWidget);
      // Quick Scan is hidden because the ai feature is disabled.
      expect(find.text('Quick Scan'), findsNothing);
    });

    testWidgets('shows Quick Scan only when ai is enabled and permitted',
        (tester) async {
      await tester.pumpWidget(_app(
        auth: _auth(aiPoweredReceipts: true),
        permissions: _permissions([
          api.Permission.groupPeriodReceiptsPeriodCreate,
          api.Permission.groupPeriodReceiptsPeriodQuickScan,
        ]),
      ));

      await tester.tap(find.text('add'));
      await tester.pumpAndSettle();

      expect(find.text('Add Manual Receipt'), findsOneWidget);
      expect(find.text('Quick Scan'), findsOneWidget);
    });

    testWidgets('shows a snackbar and no menu when the user cannot add',
        (tester) async {
      await tester.pumpWidget(_app(
        auth: _auth(aiPoweredReceipts: true),
        permissions:
            _permissions([api.Permission.groupPeriodReceiptsPeriodRead]),
      ));

      await tester.tap(find.text('add'));
      await tester.pump();

      expect(find.text('Add Manual Receipt'), findsNothing);
      expect(find.text('Quick Scan'), findsNothing);
      expect(
        find.text("You don't have permission to add receipts here."),
        findsOneWidget,
      );
    });
  });
}
