import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/top_app_bar.dart';

import '../helpers/auth_test_helpers.dart';
import '../helpers/permission_test_helpers.dart';

/// Every avatar-menu entry must act on a tap ANYWHERE on its row, not only on
/// its label.
///
/// The entries used to be TextButtons nested in PopupMenuItems that had no
/// `onTap` of their own. Only the button's box (roughly the label) ran the
/// action; the rest of the row was the PopupMenuItem's InkWell, which pops the
/// menu and does nothing else -- so Logout routinely took two taps. Each case
/// here taps the row's far edge, well clear of the label, which is exactly the
/// tap the old layout dropped.
void main() {
  late MockAuthModel authModel;
  late MockAuthApi authApi;

  setUpAll(() {
    registerFallbackValue(FakeLogoutCommand());
  });

  setUp(() {
    authModel = MockAuthModel();
    when(() => authModel.claims)
        .thenReturn(api.Claims((b) => b..displayName = 'Admin'));
    when(() => authModel.getRefreshToken())
        .thenAnswer((_) async => 'refresh-token');
    when(() => authModel.purgeTokens()).thenAnswer((_) async {});

    authApi = MockAuthApi();
    when(() => authApi.logout(logoutCommand: any(named: 'logoutCommand')))
        .thenAnswer((_) async => Response<void>(
              requestOptions: RequestOptions(path: '/logout'),
              statusCode: 200,
            ));

    // The default client builds a real Openapi with a relative baseUrl that
    // dio rejects off-Web, so install the mock rather than saving the default.
    final client = MockOpenapi();
    when(() => client.getAuthApi()).thenReturn(authApi);
    OpenApiClient.client = client;
  });

  Future<void> pumpApp(WidgetTester tester) async {
    final router = GoRouter(
      routes: [
        GoRoute(
          path: '/',
          builder: (_, __) =>
              const Scaffold(appBar: TopAppBar(titleText: 'Home')),
        ),
        GoRoute(
          path: '/profile',
          builder: (_, __) => const Scaffold(body: Text('profile page')),
        ),
        GoRoute(
          path: '/reports',
          builder: (_, __) => const Scaffold(body: Text('reports page')),
        ),
        GoRoute(
          path: '/login',
          builder: (_, __) => const Scaffold(body: Text('login page')),
        ),
      ],
    );

    await tester.pumpWidget(
      MultiProvider(
        providers: [
          ChangeNotifierProvider<AuthModel>(create: (_) => authModel),
          ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
          ChangeNotifierProvider<UserModel>(create: (_) => UserModel()),
          ChangeNotifierProvider<PermissionsModel>(
            create: (_) =>
                seededPermissions(app: [Permission.appPeriodReportsPeriodRead]),
          ),
        ],
        child: MaterialApp.router(routerConfig: router),
      ),
    );
    await tester.pump();
  }

  /// Opens the avatar menu and taps the [label] entry's row at its far edge,
  /// asserting first that the point really is off the label -- otherwise the
  /// test would pass against the old layout for the wrong reason.
  Future<void> tapRowAwayFromLabel(WidgetTester tester, String label) async {
    await tester.tap(find.byKey(const ValueKey('user-avatar-menu')));
    await tester.pump();
    await tester.pump(const Duration(milliseconds: 400));

    final labelFinder = find.text(label);
    final row = find.ancestor(
      of: labelFinder,
      matching: find.byWidgetPredicate((w) => w is PopupMenuItem),
    );
    final rowRect = tester.getRect(row);
    final point = rowRect.centerRight - const Offset(12, 0);

    expect(rowRect.contains(point), isTrue);
    expect(tester.getRect(labelFinder).contains(point), isFalse,
        reason: 'the tap must land on the row, clear of the "$label" label');

    await tester.tapAt(point);
    // Menu close + route push, then any awaited work in the action.
    for (var i = 0; i < 5; i++) {
      await tester.pump(const Duration(milliseconds: 200));
    }
  }

  testWidgets('User Profile opens from anywhere on its row', (tester) async {
    await pumpApp(tester);
    await tapRowAwayFromLabel(tester, 'User Profile');

    expect(find.text('profile page'), findsOneWidget);
  });

  testWidgets('Reports opens from anywhere on its row', (tester) async {
    await pumpApp(tester);
    await tapRowAwayFromLabel(tester, 'Reports');

    expect(find.text('reports page'), findsOneWidget);
  });

  testWidgets('Logout logs out from anywhere on its row', (tester) async {
    await pumpApp(tester);
    await tapRowAwayFromLabel(tester, 'Logout');

    verify(() => authApi.logout(logoutCommand: any(named: 'logoutCommand')))
        .called(1);
    verify(() => authModel.purgeTokens()).called(1);
    expect(find.text('login page'), findsOneWidget);
  });
}
