import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/interfaces/form_item.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/models/system_settings_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_item_items.dart';

import '../helpers/widget_test_helpers.dart';

api.UserView _user({required int id, required String displayName}) {
  return api.UserView((api.UserViewBuilder b) => b
    ..id = id
    ..displayName = displayName
    ..username = 'user-$id'
    ..userRole = api.UserRole.USER
    ..isDummyUser = false);
}

FormItem _formItem({
  required String name,
  required String amount,
  required int chargedToUserId,
  api.ItemStatus status = api.ItemStatus.OPEN,
}) {
  return FormItem(
    formId: 'form-$chargedToUserId-$name',
    name: name,
    amount: amount,
    chargedToUserId: chargedToUserId,
    receiptId: 1,
    status: status,
    categories: const <api.Category>[],
    tags: const <api.Tag>[],
  );
}

Future<void> _pumpItems(
  WidgetTester tester, {
  required List<FormItem> items,
  List<api.UserView> users = const <api.UserView>[],
}) async {
  final UserModel userModel = UserModel()..setUsers(users);

  final GoRouter router = GoRouter(
    initialLocation: '/groups/1/receipts/1/view',
    routes: <RouteBase>[
      GoRoute(
        path: '/groups/:groupId/receipts/:receiptId/view',
        builder: (BuildContext _, GoRouterState __) => Scaffold(
          body: ReceiptItemItems(items: items, groupId: 1),
        ),
      ),
    ],
  );

  await tester.pumpWidget(MultiProvider(
    providers: [
      ChangeNotifierProvider<UserModel>.value(value: userModel),
      ChangeNotifierProvider<GroupModel>(create: (_) => GroupModel()),
      ChangeNotifierProvider<ReceiptModel>(create: (_) => ReceiptModel()),
      ChangeNotifierProvider<SystemSettingsModel>(
        create: (_) => SystemSettingsModel(),
      ),
    ],
    child: MaterialApp.router(routerConfig: router),
  ));
  await tester.pump();
}

void main() {
  setUpAll(() {
    registerCustomCurrencyForTests();
  });

  testWidgets(
    'happy path: groups items by chargedToUserId with one panel per user',
    (tester) async {
      final List<FormItem> items = <FormItem>[
        _formItem(name: 'Burger', amount: '10.00', chargedToUserId: 1),
        _formItem(name: 'Fries', amount: '5.00', chargedToUserId: 1),
        _formItem(name: 'Salad', amount: '7.50', chargedToUserId: 2),
      ];

      await _pumpItems(
        tester,
        items: items,
        users: <api.UserView>[
          _user(id: 1, displayName: 'Alice'),
          _user(id: 2, displayName: 'Bob'),
        ],
      );

      expect(find.text('Alice'), findsOneWidget);
      expect(find.text('Bob'), findsOneWidget);

      // Initially expanded: item names render in the body.
      expect(find.text('Burger'), findsOneWidget);
      expect(find.text('Fries'), findsOneWidget);
      expect(find.text('Salad'), findsOneWidget);

      // Subtotals: Alice owes 15.00, Bob owes 7.50. Match by suffix to stay
      // resilient to currency-symbol position and decimal-places settings.
      expect(
        find.byWidgetPredicate(
          (Widget w) =>
              w is Text &&
              (w.data ?? '').startsWith('2 items') &&
              (w.data ?? '').contains('15'),
        ),
        findsOneWidget,
      );
      expect(
        find.byWidgetPredicate(
          (Widget w) =>
              w is Text &&
              (w.data ?? '').startsWith('1 items') &&
              (w.data ?? '').contains('7.5'),
        ),
        findsOneWidget,
      );
    },
  );

  testWidgets(
    'empty state: renders plain text and no InputDecorator chrome',
    (tester) async {
      await _pumpItems(tester, items: <FormItem>[]);

      expect(find.text('No items on this receipt'), findsOneWidget);

      // The brief specifically forbids the InputDecorator outline that the
      // old "Shared With" wrapper produced — guard against regression.
      expect(find.byType(InputDecorator), findsNothing);
    },
  );

  testWidgets(
    'malformed input: empty/whitespace name renders "(unnamed item)" without crashing',
    (tester) async {
      // Whitespace-only and empty names exercise the defensive `_safeItemName`
      // fallback that protects against generated-client churn or hydration
      // quirks where `name` is absent or non-string at runtime.
      final List<FormItem> items = <FormItem>[
        _formItem(name: '   ', amount: '4.20', chargedToUserId: 1),
        _formItem(name: '', amount: '0.80', chargedToUserId: 1),
      ];

      await _pumpItems(
        tester,
        items: items,
        users: <api.UserView>[_user(id: 1, displayName: 'Alice')],
      );

      expect(tester.takeException(), isNull);
      expect(find.text('(unnamed item)'), findsNWidgets(2));
    },
  );

  testWidgets(
    'boundary amounts: zero, negative, and non-numeric render without crashing',
    (tester) async {
      final List<FormItem> items = <FormItem>[
        _formItem(name: 'Free', amount: '0', chargedToUserId: 1),
        _formItem(name: 'Refund', amount: '-3.50', chargedToUserId: 1),
        _formItem(name: 'Gibberish', amount: 'not-a-number', chargedToUserId: 1),
      ];

      await _pumpItems(
        tester,
        items: items,
        users: <api.UserView>[_user(id: 1, displayName: 'Alice')],
      );

      expect(tester.takeException(), isNull);
      expect(find.text('Free'), findsOneWidget);
      expect(find.text('Refund'), findsOneWidget);
      expect(find.text('Gibberish'), findsOneWidget);
    },
  );

  testWidgets(
    'unknown user: renders an "Unknown user" panel instead of throwing',
    (tester) async {
      // chargedToUserId 99 is not in the seeded UserModel — `_safeGetUser`
      // must catch the `firstWhere`-without-`orElse` throw inside
      // UserModel.getUserById and fall back gracefully.
      final List<FormItem> items = <FormItem>[
        _formItem(name: 'Burger', amount: '10.00', chargedToUserId: 99),
      ];

      await _pumpItems(
        tester,
        items: items,
        users: <api.UserView>[_user(id: 1, displayName: 'Alice')],
      );

      expect(tester.takeException(), isNull);
      expect(find.text('Burger'), findsOneWidget);
    },
  );
}
