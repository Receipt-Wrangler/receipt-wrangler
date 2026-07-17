import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/top_app_bar.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/user_avatar.dart';

import '../helpers/auth_test_helpers.dart';
import '../helpers/permission_test_helpers.dart';

api.Claims _claims(String displayName) =>
    api.Claims((b) => b..displayName = displayName);

void main() {
  setUpAll(() {
    // Allow Provider<AuthModel>.value with a mocktail Mock of the ChangeNotifier
    // subclass (mirrors group_app_bar_test); the test relies on no listener
    // rebuilds.
    Provider.debugCheckInvalidValueType = null;
  });

  Future<void> pumpBar(WidgetTester tester, PermissionsModel permissions) async {
    final authModel = MockAuthModel();
    when(() => authModel.claims).thenReturn(_claims('Admin'));

    await tester.pumpWidget(
      MultiProvider(
        providers: [
          Provider<AuthModel>.value(value: authModel),
          ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
          ChangeNotifierProvider<UserModel>(create: (_) => UserModel()),
          ChangeNotifierProvider<PermissionsModel>.value(value: permissions),
        ],
        child: const MaterialApp(
          home: Scaffold(appBar: TopAppBar(titleText: 'Home')),
        ),
      ),
    );
    await tester.pump();
    // Open the avatar popup menu.
    await tester.tap(find.byType(UserAvatar));
    await tester.pumpAndSettle();
  }

  testWidgets('shows Reports with app.reports.read', (tester) async {
    await pumpBar(
      tester,
      seededPermissions(app: [Permission.appPeriodReportsPeriodRead]),
    );
    expect(find.text('Reports'), findsOneWidget);
  });

  testWidgets('shows Reports with app.reports.readAll (base absent)',
      (tester) async {
    await pumpBar(
      tester,
      seededPermissions(app: [Permission.appPeriodReportsPeriodReadAll]),
    );
    expect(find.text('Reports'), findsOneWidget);
  });

  testWidgets('hides Reports without either read permission', (tester) async {
    await pumpBar(
      tester,
      seededPermissions(app: [Permission.appPeriodReceiptsPeriodSearch]),
    );
    expect(find.text('Reports'), findsNothing);
    // The other menu items still render.
    expect(find.text('User Profile'), findsOneWidget);
    expect(find.text('Logout'), findsOneWidget);
  });
}
