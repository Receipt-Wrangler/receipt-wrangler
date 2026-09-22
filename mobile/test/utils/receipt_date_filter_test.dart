import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/utils/receipt_date_filter.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

import '../helpers/receipt_filter_test_helpers.dart';

/// The month <-> condition translation the quick date control is built on, a
/// port of desktop's `receipt-date-filter.spec.ts`.
void main() {
  group("monthFilterCondition", () {
    test("is a BETWEEN over the month's first and last day", () {
      final condition = monthFilterCondition(const FilterMonth(2026, 9));

      expect(condition.operation, api.FilterOperation.BETWEEN);
      expect(condition.value, [DateTime(2026, 9, 1), DateTime(2026, 9, 30)]);
    });

    test("ends on the 28th in a common-year February", () {
      expect(monthFilterCondition(const FilterMonth(2026, 2)).value,
          [DateTime(2026, 2, 1), DateTime(2026, 2, 28)]);
    });

    test("ends on the 29th in a leap-year February", () {
      expect(monthFilterCondition(const FilterMonth(2028, 2)).value,
          [DateTime(2028, 2, 1), DateTime(2028, 2, 29)]);
    });

    test("ends on the 31st in December without rolling into January", () {
      // `DateTime(y, 13, 0)` has to land on Dec 31 of the same year, not on a
      // day in the next one.
      expect(monthFilterCondition(const FilterMonth(2026, 12)).value,
          [DateTime(2026, 12, 1), DateTime(2026, 12, 31)]);
    });
  });

  group("monthFromCondition", () {
    test("round-trips a month written by monthFilterCondition", () {
      const month = FilterMonth(2026, 9);

      expect(monthFromCondition(monthFilterCondition(month)), month);
    });

    test("reads a pair the date-range picker wrote, both at local midnight", () {
      // The advanced filter's `showDateRangePicker` yields plain midnights, so
      // matching on calendar fields rather than on instants is what keeps the
      // stepper naming a month it did not author itself.
      final condition = ReceiptFilterCondition(
        operation: api.FilterOperation.BETWEEN,
        value: [DateTime(2026, 6, 1), DateTime(2026, 6, 30)],
      );

      expect(monthFromCondition(condition), const FilterMonth(2026, 6));
    });

    test("is null for a range that stops short of the month's end", () {
      final condition = ReceiptFilterCondition(
        operation: api.FilterOperation.BETWEEN,
        value: [DateTime(2026, 9, 1), DateTime(2026, 9, 29)],
      );

      expect(monthFromCondition(condition), isNull);
    });

    test("is null for a range starting after the first", () {
      final condition = ReceiptFilterCondition(
        operation: api.FilterOperation.BETWEEN,
        value: [DateTime(2026, 9, 2), DateTime(2026, 9, 30)],
      );

      expect(monthFromCondition(condition), isNull);
    });

    test("is null for a range spanning two months", () {
      final condition = ReceiptFilterCondition(
        operation: api.FilterOperation.BETWEEN,
        value: [DateTime(2026, 9, 1), DateTime(2026, 10, 31)],
      );

      expect(monthFromCondition(condition), isNull);
    });

    test("is null for an operation other than BETWEEN", () {
      final condition = ReceiptFilterCondition(
        operation: api.FilterOperation.GREATER_THAN,
        value: DateTime(2026, 9, 1),
      );

      expect(monthFromCondition(condition), isNull);
    });

    test("is null for WITHIN_CURRENT_MONTH, which carries no value", () {
      // It reads as "Custom" on the stepper: the server pins it to month-start
      // through today, so it is not the whole month it looks like.
      const condition = ReceiptFilterCondition(
          operation: api.FilterOperation.WITHIN_CURRENT_MONTH);

      expect(monthFromCondition(condition), isNull);
    });

    test("is null for a malformed or absent value", () {
      expect(monthFromCondition(null), isNull);
      expect(
          monthFromCondition(const ReceiptFilterCondition(
              operation: api.FilterOperation.BETWEEN)),
          isNull);
      expect(
          monthFromCondition(ReceiptFilterCondition(
              operation: api.FilterOperation.BETWEEN,
              value: [DateTime(2026, 9, 1)])),
          isNull);
      expect(
          monthFromCondition(const ReceiptFilterCondition(
              operation: api.FilterOperation.BETWEEN,
              value: ["2026-09-01", "2026-09-30"])),
          isNull);
    });
  });

  group("shiftMonth", () {
    test("steps within a year", () {
      expect(shiftMonth(const FilterMonth(2026, 9), -1),
          const FilterMonth(2026, 8));
      expect(shiftMonth(const FilterMonth(2026, 9), 1),
          const FilterMonth(2026, 10));
    });

    test("rolls back over January", () {
      expect(shiftMonth(const FilterMonth(2026, 1), -1),
          const FilterMonth(2025, 12));
    });

    test("rolls forward over December", () {
      expect(shiftMonth(const FilterMonth(2026, 12), 1),
          const FilterMonth(2027, 1));
    });
  });

  test("monthOfDate reads the calendar month a date falls in", () {
    expect(monthOfDate(DateTime(2026, 9, 22)), const FilterMonth(2026, 9));
  });

  test("filterMonthLabel spells the month out with its year", () {
    expect(filterMonthLabel(const FilterMonth(2026, 9)), "September 2026");
  });

  group("through the real encoder", () {
    // What actually reaches the server. The boundaries are day-bounded by
    // `buildReceiptPagedRequestFilter`, so a month must arrive as
    // start-of-first-day through end-of-last-day -- a bare midnight upper bound
    // would drop everything recorded on the last day of the month.
    test("a month encodes as a full-month zulu range", () {
      final field = serializedField(
        {"date": monthFilterCondition(const FilterMonth(2026, 9))},
        "date",
      );

      expect(field, isNotNull);
      expect(field!["operation"], "BETWEEN");
      expect(field["value"],
          ["2026-09-01T00:00:00Z", "2026-09-30T23:59:59Z"]);
    });

    test("it can be written to any of the three date fields", () {
      for (final key in ["date", "resolvedDate", "createdAt"]) {
        final field = serializedField(
          {key: monthFilterCondition(const FilterMonth(2026, 2))},
          key,
        );

        expect(field?["value"],
            ["2026-02-01T00:00:00Z", "2026-02-28T23:59:59Z"],
            reason: "$key should carry the whole of February");
      }
    });
  });
}
