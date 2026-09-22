import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_month_picker_sheet.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_date_filter.dart';

/// The sheet behind the stepper's label: a year pager, the twelve months, and
/// the three shortcuts. Pumped directly rather than through
/// `showReceiptMonthPickerSheet`, so the body's own rules are isolated from the
/// sheet plumbing.
void main() {
  late List<MonthPickerResult> results;

  setUp(() => results = []);

  /// Reads only. A second `pumpPicker` in one test would reuse the existing
  /// `State` -- and with it the year it had already paged to -- so every case
  /// that needs a different seed gets its own test.
  Future<void> pumpPicker(WidgetTester tester, {FilterMonth? selected}) {
    return tester.pumpWidget(MaterialApp(
      home: Scaffold(body: ReceiptMonthPicker(selected: selected)),
    ));
  }

  /// The value the body popped with.
  MonthPickerResult popped() => results.last;

  /// Opens the body in a real modal sheet, so a pick can actually pop.
  Future<void> pumpPickerCapturingPop(
    WidgetTester tester, {
    FilterMonth? selected,
  }) async {
    await tester.pumpWidget(MaterialApp(
      home: Builder(
        builder: (context) => Scaffold(
          body: ElevatedButton(
            child: const Text("open"),
            onPressed: () async {
              final result = await showModalBottomSheet<MonthPickerResult>(
                context: context,
                isScrollControlled: true,
                builder: (_) => ReceiptMonthPicker(selected: selected),
              );
              if (result != null) {
                results.add(result);
              }
            },
          ),
        ),
      ),
    ));

    await tester.tap(find.text("open"));
    await tester.pumpAndSettle();
  }

  group("the year pager", () {
    testWidgets("seeds from the selected month", (tester) async {
      await pumpPicker(tester, selected: const FilterMonth(2019, 4));

      expect(find.text("2019"), findsOneWidget);
    });

    testWidgets("falls back to the current year with no selection",
        (tester) async {
      await pumpPicker(tester);

      expect(find.text("${DateTime.now().year}"), findsOneWidget);
    });

    testWidgets("paging selects nothing and leaves the sheet open",
        (tester) async {
      // Paging is a view concern: it must not write a filter, and it must not
      // close the sheet out from under the user mid-browse.
      await pumpPickerCapturingPop(tester,
          selected: const FilterMonth(2026, 9));

      await tester.tap(find.byKey(const ValueKey("month-picker-year-prev")));
      await tester.pump();

      expect(find.text("2025"), findsOneWidget);
      expect(results, isEmpty);
      expect(find.byType(ReceiptMonthPicker), findsOneWidget);
    });

    testWidgets("stops paging back at 1970", (tester) async {
      await pumpPicker(tester, selected: const FilterMonth(1970, 1));

      expect(
          tester
              .widget<IconButton>(
                  find.byKey(const ValueKey("month-picker-year-prev")))
              .onPressed,
          isNull);
    });

    testWidgets("stops paging forward five years out", (tester) async {
      await pumpPicker(tester,
          selected: FilterMonth(DateTime.now().year + 5, 1));

      expect(
          tester
              .widget<IconButton>(
                  find.byKey(const ValueKey("month-picker-year-next")))
              .onPressed,
          isNull);
    });
  });

  group("the month grid", () {
    testWidgets("offers all twelve months", (tester) async {
      await pumpPicker(tester);

      for (final short in ["Jan", "Jun", "Dec"]) {
        expect(find.text(short), findsOneWidget);
      }
    });

    testWidgets("picks a month out of the year being paged", (tester) async {
      await pumpPickerCapturingPop(tester,
          selected: const FilterMonth(2026, 9));

      await tester.tap(find.byKey(const ValueKey("month-picker-year-prev")));
      await tester.pump();
      await tester.tap(find.byKey(const ValueKey("month-picker-month-3")));
      await tester.pumpAndSettle();

      expect(popped().month, const FilterMonth(2025, 3));
      expect(popped().isAllTime, isFalse);
    });
  });

  group("the shortcuts", () {
    testWidgets("This month picks the current one", (tester) async {
      await pumpPickerCapturingPop(tester);

      await tester.tap(find.byKey(const ValueKey("month-picker-this-month")));
      await tester.pumpAndSettle();

      expect(popped().month, monthOfDate(DateTime.now()));
    });

    testWidgets("Last month steps back one", (tester) async {
      await pumpPickerCapturingPop(tester);

      await tester.tap(find.byKey(const ValueKey("month-picker-last-month")));
      await tester.pumpAndSettle();

      expect(popped().month,
          shiftMonth(monthOfDate(DateTime.now()), -1));
    });

    testWidgets("All time is its own outcome, carrying no month",
        (tester) async {
      // Distinct from a month and from a dismissal -- the caller has to be able
      // to tell "clear the filter" from "leave it alone".
      await pumpPickerCapturingPop(tester,
          selected: const FilterMonth(2026, 9));

      await tester.tap(find.byKey(const ValueKey("month-picker-all-time")));
      await tester.pumpAndSettle();

      expect(popped().isAllTime, isTrue);
      expect(popped().month, isNull);
    });
  });
}
