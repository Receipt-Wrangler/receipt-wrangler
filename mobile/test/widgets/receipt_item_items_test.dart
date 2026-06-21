import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/interfaces/form_item.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/models/system_settings_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_item_items.dart';

import '../helpers/widget_test_helpers.dart';

class _MockUserView extends Mock implements api.UserView {}

class _MockUserModel extends Mock implements UserModel {}

api.UserView _user({required int id, required String displayName}) {
  // UserView has many non-nullable fields. Mock the accessors we read so the
  // test isn't coupled to the full builder graph (mirrors group_app_bar_test).
  final _MockUserView mock = _MockUserView();
  when(() => mock.id).thenReturn(id);
  when(() => mock.displayName).thenReturn(displayName);
  return mock;
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
  api.UserView? Function(String id)? userLookup,
}) async {
  final _MockUserModel userModel = _MockUserModel();
  when(() => userModel.getUserById(any())).thenAnswer(
    (invocation) => userLookup == null
        ? null
        : userLookup(invocation.positionalArguments.first as String),
  );

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
      // Plain Provider.value avoids ChangeNotifierProvider trying to call
      // addListener on a mocktail Mock (matches group_app_bar_test).
      Provider<UserModel>.value(value: userModel),
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
    registerFallbackValue('');
    registerCustomCurrencyForTests();
    // Allow Provider<T>.value with mocktail Mocks of ChangeNotifier subclasses
    // (matches group_app_bar_test setup).
    Provider.debugCheckInvalidValueType = null;
  });

  testWidgets(
    'happy path: groups items by chargedToUserId with one panel per user',
    (tester) async {
      final List<FormItem> items = <FormItem>[
        _formItem(name: 'Burger', amount: '10.00', chargedToUserId: 1),
        _formItem(name: 'Fries', amount: '5.00', chargedToUserId: 1),
        _formItem(name: 'Salad', amount: '7.50', chargedToUserId: 2),
      ];
      final Map<int, api.UserView> users = <int, api.UserView>{
        1: _user(id: 1, displayName: 'Alice'),
        2: _user(id: 2, displayName: 'Bob'),
      };

      await _pumpItems(
        tester,
        items: items,
        userLookup: (String id) => users[int.tryParse(id)],
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
        userLookup: (String id) =>
            id == '1' ? _user(id: 1, displayName: 'Alice') : null,
      );

      expect(tester.takeException(), isNull);
      expect(find.text('(unnamed item)'), findsNWidgets(2));
    },
  );
}
