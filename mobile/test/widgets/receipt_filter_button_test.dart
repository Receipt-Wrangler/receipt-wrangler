import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/category_model.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt-list-model.dart';
import 'package:receipt_wrangler_mobile/models/system_settings_model.dart';
import 'package:receipt_wrangler_mobile/models/tag_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_filter_button.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

import '../helpers/receipt_filter_widget_helpers.dart';
import '../helpers/widget_test_helpers.dart';

/// The badge is the only thing on the receipts screen that says a filter is
/// narrowing the list, so it has to track the applied filter rather than a
/// snapshot taken when the app bar was built.
void main() {
  setUpAll(registerCustomCurrencyForTests);

  const buttonKey = ValueKey("receipt-filter-button");

  const nameCondition = ReceiptFilterCondition(
      operation: api.FilterOperation.CONTAINS, value: "Costco");
  const amountCondition = ReceiptFilterCondition(
      operation: api.FilterOperation.GREATER_THAN, value: 50.0);

  /// Pumps the button under a router, because it reads the browsed group from
  /// the route to scope the filter screen's pickers.
  Future<ReceiptFilterHarness> pumpButton(WidgetTester tester) async {
    final harness = buildReceiptFilterHarness();

    final router = GoRouter(
      initialLocation: "/groups/${ReceiptFilterHarness.householdId}/receipts",
      routes: [
        GoRoute(
          path: "/groups/:groupId/receipts",
          builder: (context, state) => const Scaffold(
            body: Row(children: [ReceiptFilterButton()]),
          ),
        ),
      ],
    );

    await tester.pumpWidget(MultiProvider(
      providers: [
        ChangeNotifierProvider<ReceiptListModel>.value(
            value: harness.receiptListModel),
        ChangeNotifierProvider<GroupModel>.value(value: harness.groupModel),
        ChangeNotifierProvider<CategoryModel>.value(
            value: harness.categoryModel),
        ChangeNotifierProvider<TagModel>.value(value: harness.tagModel),
        ChangeNotifierProvider<UserModel>.value(value: harness.userModel),
        ChangeNotifierProvider<ContextModel>(create: (_) => ContextModel()),
        ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
        ChangeNotifierProvider<SystemSettingsModel>(
            create: (_) => SystemSettingsModel()),
      ],
      child: MaterialApp.router(routerConfig: router),
    ));
    await tester.pump();

    return harness;
  }

  Badge badgeOf(WidgetTester tester) =>
      tester.widget<Badge>(find.ancestor(
          of: find.byKey(buttonKey), matching: find.byType(Badge)));

  testWidgets("renders no count when nothing is filtered", (tester) async {
    await pumpButton(tester);

    expect(find.byKey(buttonKey), findsOneWidget);
    expect(badgeOf(tester).isLabelVisible, isFalse);
  });

  testWidgets("shows how many conditions are applied", (tester) async {
    final harness = await pumpButton(tester);

    harness.receiptListModel
        .setFilter({"name": nameCondition, "amount": amountCondition}, true, groupId: "${ReceiptFilterHarness.householdId}");
    await tester.pump();

    expect(badgeOf(tester).isLabelVisible, isTrue);
    expect(find.text("2"), findsOneWidget);
  });

  testWidgets("follows the applied filter down as well as up", (tester) async {
    final harness = await pumpButton(tester);

    harness.receiptListModel.setFilter({"name": nameCondition}, true, groupId: "${ReceiptFilterHarness.householdId}");
    await tester.pump();
    expect(find.text("1"), findsOneWidget);

    harness.receiptListModel.clearFilter(true);
    await tester.pump();

    expect(badgeOf(tester).isLabelVisible, isFalse);
  });

  testWidgets("opens the filter screen", (tester) async {
    await pumpButton(tester);

    await tester.tap(find.byKey(buttonKey));
    await tester.pumpAndSettle();

    expect(find.byKey(const ValueKey("receipt-filter-apply")), findsOneWidget);
    expect(find.text("Filter"), findsOneWidget);
  });

  testWidgets("scopes the screen to the group being browsed", (tester) async {
    // Household is a real group, so the Group condition must not be offered.
    await pumpButton(tester);

    await tester.tap(find.byKey(buttonKey));
    await tester.pumpAndSettle();
    await tester.tap(find.byKey(const ValueKey("receipt-filter-add")));
    await tester.pumpAndSettle();

    expect(find.byKey(const ValueKey("add-receipt-filter-group")), findsNothing);
    expect(find.byKey(const ValueKey("add-receipt-filter-categories")),
        findsOneWidget);
  });
}
