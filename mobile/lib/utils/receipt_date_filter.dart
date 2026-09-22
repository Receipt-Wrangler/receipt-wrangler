import 'package:intl/intl.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

/// A calendar month.
///
/// [month] is **1-based**, matching `DateTime.month`. Desktop's `FilterMonth`
/// (`desktop/src/utils/receipt-date-filter.ts`) is zero-based because it is
/// built on JavaScript's `Date`; keeping that here would mean an off-by-one at
/// every `DateTime` constructor for no gain.
class FilterMonth {
  const FilterMonth(this.year, this.month);

  final int year;

  /// 1 (January) through 12 (December).
  final int month;

  @override
  bool operator ==(Object other) =>
      other is FilterMonth && other.year == year && other.month == month;

  @override
  int get hashCode => Object.hash(year, month);

  @override
  String toString() => "FilterMonth($year, $month)";
}

/// The condition describing a whole calendar month.
///
/// A month is expressed as `BETWEEN [first day, last day]` -- the one operation
/// that can describe *any* month, which is why the quick date control needs no
/// API change. `WITHIN_CURRENT_MONTH` is not equivalent: the server pins it to
/// month-start through *today*, so it can only ever mean the current month.
///
/// The boundaries are plain dates rather than day-bounded instants on purpose:
/// [buildReceiptPagedRequestFilter] re-applies `startOfDay` / `endOfDay` to
/// whatever pair it is handed, so narrowing them here would only duplicate a
/// rule that already has one home.
ReceiptFilterCondition monthFilterCondition(FilterMonth month) {
  return ReceiptFilterCondition(
    operation: api.FilterOperation.BETWEEN,
    // Day 0 of the following month is the last day of this one, and month 13
    // rolls into the next January -- so February and December need no special
    // case.
    value: [
      DateTime(month.year, month.month, 1),
      DateTime(month.year, month.month + 1, 0),
    ],
  );
}

/// The month [condition] describes, or null when it describes anything else --
/// a partial range, a `GREATER_THAN`, `WITHIN_CURRENT_MONTH`, nothing at all.
///
/// The comparison is on calendar fields rather than on instants, so a pair the
/// advanced filter's date-range picker wrote (both at local midnight) still
/// reads as the month it spans.
FilterMonth? monthFromCondition(ReceiptFilterCondition? condition) {
  if (condition?.operation != api.FilterOperation.BETWEEN) {
    return null;
  }

  final value = condition!.value;
  if (value is! List || value.length != 2) {
    return null;
  }

  final start = value[0];
  final end = value[1];
  if (start is! DateTime || end is! DateTime) {
    return null;
  }

  final spansWholeMonth = start.day == 1 &&
      start.year == end.year &&
      start.month == end.month &&
      end.day == _daysInMonth(start.year, start.month);

  return spansWholeMonth ? FilterMonth(start.year, start.month) : null;
}

/// Steps [month] by [delta] months, rolling the year at either boundary.
FilterMonth shiftMonth(FilterMonth month, int delta) {
  final shifted = DateTime(month.year, month.month + delta, 1);

  return FilterMonth(shifted.year, shifted.month);
}

/// The calendar month [date] falls in.
FilterMonth monthOfDate(DateTime date) => FilterMonth(date.year, date.month);

/// How many days [month] of [year] has, leap years included. Private because
/// nothing outside this file needs it; if something does, it belongs in
/// `utils/date.dart` beside `startOfDay` / `endOfDay`.
int _daysInMonth(int year, int month) => DateTime(year, month + 1, 0).day;

/// How a month reads on the stepper, e.g. "September 2026".
String filterMonthLabel(FilterMonth month) {
  return DateFormat("MMMM y").format(DateTime(month.year, month.month, 1));
}
