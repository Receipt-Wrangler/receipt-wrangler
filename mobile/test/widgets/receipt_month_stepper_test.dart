import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_month_stepper.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_date_filter.dart';

/// The stepper is presentational: it is told a month and a label, and emits.
/// These pin the one rule it owns -- stepping from nothing seeds today -- plus
/// the clear button's visibility, which follows the *condition*, not the month.
void main() {
  late List<FilterMonth> selected;
  late int allTimeCount;
  late int labelPressCount;

  setUp(() {
    selected = [];
    allTimeCount = 0;
    labelPressCount = 0;
  });

  Future<void> pumpStepper(
    WidgetTester tester, {
    FilterMonth? month,
    String label = "All time",
    bool hasCondition = false,
  }) {
    return tester.pumpWidget(MaterialApp(
      home: Scaffold(
        body: ReceiptMonthStepper(
          month: month,
          label: label,
          hasCondition: hasCondition,
          onMonthSelected: selected.add,
          onAllTimeSelected: () => allTimeCount++,
          onLabelPressed: () => labelPressCount++,
        ),
      ),
    ));
  }

  testWidgets("renders the label it is given", (tester) async {
    await pumpStepper(tester,
        month: const FilterMonth(2026, 9), label: "September 2026");

    expect(find.text("September 2026"), findsOneWidget);
  });

  testWidgets("steps back and forward from the bound month", (tester) async {
    await pumpStepper(tester, month: const FilterMonth(2026, 9));

    await tester.tap(find.byKey(const ValueKey("receipt-month-prev")));
    await tester.tap(find.byKey(const ValueKey("receipt-month-next")));

    expect(selected, [const FilterMonth(2026, 8), const FilterMonth(2026, 10)]);
  });

  testWidgets("steps from today when no month is showing", (tester) async {
    // Without this the arrows no-op on "All time" and "Custom", which is every
    // first press. Both arrows must also land somewhere different.
    await pumpStepper(tester, month: null);

    await tester.tap(find.byKey(const ValueKey("receipt-month-prev")));
    await tester.tap(find.byKey(const ValueKey("receipt-month-next")));

    final thisMonth = monthOfDate(DateTime.now());
    expect(selected,
        [shiftMonth(thisMonth, -1), shiftMonth(thisMonth, 1)]);
    expect(selected.first, isNot(selected.last));
  });

  testWidgets("rolls the year at either boundary", (tester) async {
    await pumpStepper(tester, month: const FilterMonth(2026, 1));
    await tester.tap(find.byKey(const ValueKey("receipt-month-prev")));

    await pumpStepper(tester, month: const FilterMonth(2026, 12));
    await tester.tap(find.byKey(const ValueKey("receipt-month-next")));

    expect(selected, [const FilterMonth(2025, 12), const FilterMonth(2027, 1)]);
  });

  testWidgets("offers no clear while nothing is filtered", (tester) async {
    await pumpStepper(tester, hasCondition: false);

    expect(find.byKey(const ValueKey("receipt-month-clear")), findsNothing);
  });

  testWidgets("clears to all time when a condition is set", (tester) async {
    await pumpStepper(tester, hasCondition: true, label: "September 2026");

    await tester.tap(find.byKey(const ValueKey("receipt-month-clear")));

    expect(allTimeCount, 1);
    expect(selected, isEmpty);
  });

  testWidgets("offers the clear for a condition it cannot describe",
      (tester) async {
    // "Custom" carries no month, but it is still a filter the user needs a way
    // out of -- so the button follows the condition, not the month.
    await pumpStepper(tester, month: null, label: "Custom", hasCondition: true);

    expect(find.byKey(const ValueKey("receipt-month-clear")), findsOneWidget);
  });

  testWidgets("the label opens the month picker rather than selecting",
      (tester) async {
    await pumpStepper(tester, month: const FilterMonth(2026, 9));

    await tester.tap(find.byKey(const ValueKey("receipt-month-label")));

    expect(labelPressCount, 1);
    expect(selected, isEmpty);
  });
}
