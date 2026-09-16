import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/receipts/screens/receipt_filter_screen.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_filter_condition_card.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

import '../helpers/receipt_filter_widget_helpers.dart';
import '../helpers/widget_test_helpers.dart';

/// The screen edits a DRAFT copy and commits only on Apply. That split is the
/// easiest thing here to get wrong -- committing on every edit would refetch the
/// list mid-authoring, and never committing would make the button do nothing.
void main() {
  setUpAll(registerCustomCurrencyForTests);

  const applyKey = ValueKey("receipt-filter-apply");
  const resetKey = ValueKey("receipt-filter-reset");
  const addKey = ValueKey("receipt-filter-add");
  const emptyKey = ValueKey("receipt-filter-empty");
  const countKey = ValueKey("receipt-filter-count");

  const nameCondition = ReceiptFilterCondition(
      operation: api.FilterOperation.CONTAINS, value: "Costco");
  const amountCondition = ReceiptFilterCondition(
      operation: api.FilterOperation.GREATER_THAN, value: 50.0);

  Future<ReceiptFilterHarness> pumpScreen(
    WidgetTester tester, {
    Map<String, ReceiptFilterCondition> applied = const {},
    String groupId = "${ReceiptFilterHarness.householdId}",
  }) async {
    final harness = buildReceiptFilterHarness();
    harness.receiptListModel.setFilter(applied, false);

    await pumpWithFilterHarness(
        tester, harness, ReceiptFilterScreen(groupId: groupId));

    return harness;
  }

  group("the condition list", () {
    testWidgets("says so when there is nothing to show", (tester) async {
      await pumpScreen(tester);

      expect(find.byKey(emptyKey), findsOneWidget);
      expect(find.text("No conditions yet — everything is showing."),
          findsOneWidget);
      expect(find.byType(ReceiptFilterConditionCard), findsNothing);
    });

    testWidgets("offers Add filter even when empty", (tester) async {
      await pumpScreen(tester);

      expect(find.byKey(addKey), findsOneWidget);
    });

    testWidgets("renders one card per applied condition", (tester) async {
      await pumpScreen(
          tester, applied: {"name": nameCondition, "amount": amountCondition});

      expect(find.byType(ReceiptFilterConditionCard), findsNWidgets(2));
      expect(find.byKey(emptyKey), findsNothing);
    });

    testWidgets("a card reads its field, operation and value", (tester) async {
      await pumpScreen(tester, applied: {"name": nameCondition});

      expect(find.text("Name"), findsOneWidget);
      expect(find.text("Contains"), findsOneWidget);
      expect(find.text("Costco"), findsOneWidget);
    });

    testWidgets("cards render in the shared table's order, not insertion order",
        (tester) async {
      // "name" is declared before "amount", so it renders first however the map
      // was built -- the screen looks the same however the user got there.
      final harness = buildReceiptFilterHarness();
      harness.receiptListModel.setFilter(
          {"amount": amountCondition, "name": nameCondition}, false);

      await pumpWithFilterHarness(tester, harness,
          const ReceiptFilterScreen(groupId: "${ReceiptFilterHarness.householdId}"));

      final cards = tester
          .widgetList<ReceiptFilterConditionCard>(
              find.byType(ReceiptFilterConditionCard))
          .toList();

      expect(cards.map((card) => card.field.key).toList(), ["name", "amount"]);
    });
  });

  group("the condition count header", () {
    testWidgets("is absent with nothing to count, where the empty state speaks",
        (tester) async {
      await pumpScreen(tester);

      expect(find.byKey(countKey), findsNothing);
      expect(find.byKey(emptyKey), findsOneWidget);
    });

    testWidgets("is singular for one condition", (tester) async {
      await pumpScreen(tester, applied: {"name": nameCondition});

      expect(find.text("1 CONDITION"), findsOneWidget);
      expect(find.byKey(emptyKey), findsNothing);
    });

    testWidgets("counts the draft, not the applied filter", (tester) async {
      await pumpScreen(tester,
          applied: {"name": nameCondition, "amount": amountCondition});
      expect(find.text("2 CONDITIONS"), findsOneWidget);

      await tester.tap(find.byKey(const ValueKey("receipt-filter-remove-name")));
      await tester.pump();

      expect(find.text("1 CONDITION"), findsOneWidget);
    });
  });

  group("the draft is not the applied filter", () {
    testWidgets("removing a card leaves the applied filter alone",
        (tester) async {
      final harness = await pumpScreen(tester, applied: {"name": nameCondition});

      await tester.tap(find.byKey(const ValueKey("receipt-filter-remove-name")));
      await tester.pump();

      expect(find.byType(ReceiptFilterConditionCard), findsNothing);
      expect(harness.receiptListModel.activeFilterCount, 1,
          reason: "the list must not refetch until Apply");
    });

    testWidgets("reset clears the cards but not the applied filter",
        (tester) async {
      final harness = await pumpScreen(
          tester, applied: {"name": nameCondition, "amount": amountCondition});

      await tester.tap(find.byKey(resetKey));
      await tester.pump();

      expect(find.byKey(emptyKey), findsOneWidget);
      expect(harness.receiptListModel.activeFilterCount, 2);
    });

    testWidgets("leaving without applying discards the edits", (tester) async {
      final harness = buildReceiptFilterHarness();
      harness.receiptListModel.setFilter({"name": nameCondition}, false);

      await pumpWithFilterHarness(
        tester,
        harness,
        Builder(
          builder: (context) => Scaffold(
            body: ElevatedButton(
              onPressed: () => Navigator.of(context).push(MaterialPageRoute(
                  builder: (_) => const ReceiptFilterScreen(
                      groupId: "${ReceiptFilterHarness.householdId}"))),
              child: const Text("open"),
            ),
          ),
        ),
      );

      await tester.tap(find.text("open"));
      await tester.pumpAndSettle();

      await tester.tap(find.byKey(const ValueKey("receipt-filter-remove-name")));
      await tester.pump();
      await tester.tap(find.byKey(const ValueKey("receipt-filter-close")));
      await tester.pumpAndSettle();

      expect(harness.receiptListModel.activeFilterCount, 1);
      expect(harness.receiptListModel.filter.containsKey("name"), isTrue);
    });
  });

  group("Apply", () {
    testWidgets("commits the draft and notifies once", (tester) async {
      final harness = await pumpScreen(tester, applied: {"name": nameCondition});

      var notifications = 0;
      harness.receiptListModel.addListener(() => notifications++);

      await tester.tap(find.byKey(const ValueKey("receipt-filter-remove-name")));
      await tester.pump();
      await tester.tap(find.byKey(applyKey));
      await tester.pump();

      expect(harness.receiptListModel.activeFilterCount, 0);
      expect(notifications, 1,
          reason: "the list refetches on each notification");
    });

    testWidgets("reads 'Apply Filter', not a match count", (tester) async {
      // The count would need a speculative request for a filter that has not
      // been applied yet.
      await pumpScreen(tester);

      expect(find.text("Apply Filter"), findsOneWidget);
    });

    testWidgets("applying an empty draft clears the filter", (tester) async {
      final harness = await pumpScreen(tester, applied: {"name": nameCondition});

      await tester.tap(find.byKey(resetKey));
      await tester.pump();
      await tester.tap(find.byKey(applyKey));
      await tester.pump();

      expect(harness.receiptListModel.hasActiveFilter, isFalse);
    });
  });

  group("reset", () {
    testWidgets("is disabled with nothing to reset", (tester) async {
      await pumpScreen(tester);

      final button = tester.widget<IconButton>(find.byKey(resetKey));
      expect(button.onPressed, isNull);
    });

    testWidgets("is enabled once there is a condition", (tester) async {
      await pumpScreen(tester, applied: {"name": nameCondition});

      final button = tester.widget<IconButton>(find.byKey(resetKey));
      expect(button.onPressed, isNotNull);
    });
  });

  group("the Add filter sheet", () {
    testWidgets("offers only the fields not already used", (tester) async {
      await pumpScreen(tester, applied: {"name": nameCondition});

      await tester.tap(find.byKey(addKey));
      await tester.pumpAndSettle();

      expect(find.byKey(const ValueKey("add-receipt-filter-name")), findsNothing);
      expect(find.byKey(const ValueKey("add-receipt-filter-amount")),
          findsOneWidget);
    });

    testWidgets("hides Group inside a real group", (tester) async {
      // The receipts endpoint already scopes every query to that group, so a
      // group condition there is redundant at best.
      await pumpScreen(tester, groupId: "${ReceiptFilterHarness.householdId}");

      await tester.tap(find.byKey(addKey));
      await tester.pumpAndSettle();

      expect(find.byKey(const ValueKey("add-receipt-filter-group")),
          findsNothing);
    });

    testWidgets("offers Group on the All group", (tester) async {
      await pumpScreen(tester, groupId: "${ReceiptFilterHarness.allGroupId}");

      await tester.tap(find.byKey(addKey));
      await tester.pumpAndSettle();

      expect(find.byKey(const ValueKey("add-receipt-filter-group")),
          findsOneWidget);
    });

    testWidgets("picking a field opens its editor", (tester) async {
      await pumpScreen(tester);

      await tester.tap(find.byKey(addKey));
      await tester.pumpAndSettle();
      await tester.tap(find.byKey(const ValueKey("add-receipt-filter-name")));
      await tester.pumpAndSettle();

      expect(find.byKey(const ValueKey("receipt-filter-condition-save")),
          findsOneWidget);
    });
  });
}
