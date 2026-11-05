import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/receipts/nav/receipt_app_bar_action_builder.dart';
import 'package:receipt_wrangler_mobile/receipts/screens/receipt_history_screen.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_edit_popup_menu.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';

import '../helpers/permission_test_helpers.dart';
import '../helpers/receipt_history_test_helpers.dart';

class _MockOpenapi extends Mock implements api.Openapi {}

class _MockReceiptApi extends Mock implements api.ReceiptApi {}

const _receiptId = 5;
const _groupId = 7;

/// Stands in for the receipt form screen: an app bar built the way
/// `ReceiptFormScreen` builds it, from the live route.
class _MenuHost extends StatelessWidget {
  const _MenuHost();

  @override
  Widget build(BuildContext context) {
    final receiptModel = Provider.of<ReceiptModel>(context, listen: false);
    final builder = ReceiptAppBarActionBuilder(context, receiptModel);
    return Scaffold(
      appBar: AppBar(actions: builder.buildAppBarMenu(GoRouterState.of(context))),
    );
  }
}

GoRouter _router(String location) => GoRouter(
      initialLocation: location,
      routes: [
        GoRoute(
          path: '/receipts/:receiptId/view',
          builder: (context, state) => const _MenuHost(),
        ),
        GoRoute(
          path: '/receipts/:receiptId/edit',
          builder: (context, state) => const _MenuHost(),
        ),
      ],
    );

Future<GoRouter> _pump(
  WidgetTester tester, {
  required String mode,
  List<api.Permission> app = const [],
  List<api.Permission> group = const [api.Permission.groupPeriodReceiptsPeriodRead],
  int loadedReceiptId = _receiptId,
}) async {
  final receiptModel = ReceiptModel()
    ..setReceipt(
        getDefaultReceipt().rebuild((b) => b
          ..id = loadedReceiptId
          ..groupId = _groupId),
        false);
  final router = _router('/receipts/$_receiptId/$mode');

  await tester.pumpWidget(MultiProvider(
    providers: [
      ChangeNotifierProvider<PermissionsModel>.value(
          value: seededPermissions(app: app, group: {_groupId: group})),
      ChangeNotifierProvider<ReceiptModel>.value(value: receiptModel),
      ChangeNotifierProvider<UserModel>(create: (_) => UserModel()),
    ],
    child: MaterialApp.router(routerConfig: router),
  ));
  await tester.pumpAndSettle();
  return router;
}

Finder _viewMenuButton() => find.descendant(
      of: find.byType(ReceiptEditPopupMenu),
      matching: find.byType(PopupMenuButton<dynamic>),
    );

Finder _editMenuButton() =>
    find.byKey(ReceiptAppBarActionBuilder.editOverflowMenuKey);

void main() {
  setUp(() {
    final client = _MockOpenapi();
    final receiptApi = _MockReceiptApi();
    when(() => client.getReceiptApi()).thenReturn(receiptApi);
    when(() => receiptApi.getReceiptHistory(receiptId: any(named: 'receiptId')))
        .thenAnswer((_) async => Response<api.ReceiptHistory>(
              requestOptions: RequestOptions(path: ''),
              statusCode: 200,
              data: historyFromJson(historyJson([])),
            ));
    OpenApiClient.client = client;
  });

  group('view mode', () {
    testWidgets('an editor with the history permission gets Edit, then History',
        (tester) async {
      await _pump(tester,
          mode: 'view',
          app: [api.Permission.appPeriodReceiptHistoryPeriodRead],
          group: [api.Permission.groupPeriodReceiptsPeriodUpdate]);

      await tester.tap(_viewMenuButton());
      await tester.pumpAndSettle();

      expect(find.text('Edit'), findsOneWidget);
      expect(find.byKey(ReceiptAppBarActionBuilder.historyMenuItemKey), findsOneWidget);
    });

    testWidgets('a read-only member with the history permission gets only History',
        (tester) async {
      await _pump(tester,
          mode: 'view', app: [api.Permission.appPeriodReceiptHistoryPeriodRead]);

      await tester.tap(_viewMenuButton());
      await tester.pumpAndSettle();

      expect(find.text('History'), findsOneWidget);
      expect(find.text('Edit'), findsNothing);
    });

    // The app bar renders while the receipt is still loading, when the model
    // holds the default receipt (or the previous screen's). History must wait
    // for this route's receipt, or it would open the wrong one's history.
    testWidgets('offers no History until this route\'s receipt has loaded',
        (tester) async {
      await _pump(tester,
          mode: 'view',
          app: [api.Permission.appPeriodReceiptHistoryPeriodRead],
          loadedReceiptId: 0);

      expect(find.byType(ReceiptEditPopupMenu), findsOneWidget);
      expect(_viewMenuButton(), findsNothing);
    });

    testWidgets('without the permission there is no History, and a read-only member no menu',
        (tester) async {
      await _pump(tester, mode: 'view');

      expect(find.byType(ReceiptEditPopupMenu), findsOneWidget);
      expect(_viewMenuButton(), findsNothing);
    });
  });

  group('edit mode', () {
    // Integration specs read ReceiptEditPopupMenu as "the view screen mounted",
    // so the edit screen's menu must be a different widget.
    testWidgets('offers History in a menu that is not the view marker', (tester) async {
      await _pump(tester,
          mode: 'edit',
          app: [api.Permission.appPeriodReceiptHistoryPeriodRead],
          group: [api.Permission.groupPeriodReceiptsPeriodUpdate]);

      expect(_editMenuButton(), findsOneWidget);
      expect(find.byType(ReceiptEditPopupMenu), findsNothing);
    });

    testWidgets('shows no menu without the permission', (tester) async {
      await _pump(tester,
          mode: 'edit', group: [api.Permission.groupPeriodReceiptsPeriodUpdate]);

      expect(_editMenuButton(), findsNothing);
      expect(find.byType(PopupMenuButton<dynamic>), findsNothing);
    });

    // The form reads view/edit from the URL, so opening History must leave it
    // alone or the edit form underneath would turn into a read-only view.
    testWidgets('opens History on top without changing the URL', (tester) async {
      final router = await _pump(tester,
          mode: 'edit',
          app: [api.Permission.appPeriodReceiptHistoryPeriodRead],
          group: [api.Permission.groupPeriodReceiptsPeriodUpdate]);

      await tester.tap(_editMenuButton());
      await tester.pumpAndSettle();
      await tester.tap(find.byKey(ReceiptAppBarActionBuilder.historyMenuItemKey));
      await tester.pumpAndSettle();

      expect(find.byType(ReceiptHistoryScreen), findsOneWidget);
      expect(
        router.routeInformationProvider.value.uri.path,
        '/receipts/$_receiptId/edit',
      );

      await tester.tap(find.byTooltip('Back'));
      await tester.pumpAndSettle();
      expect(find.byType(ReceiptHistoryScreen), findsNothing);
      expect(_editMenuButton(), findsOneWidget);
    });
  });
}
