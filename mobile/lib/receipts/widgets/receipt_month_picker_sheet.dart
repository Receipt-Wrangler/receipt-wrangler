import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_date_filter.dart';

const _shortMonthLabels = [
  "Jan", "Feb", "Mar", "Apr", "May", "Jun", //
  "Jul", "Aug", "Sep", "Oct", "Nov", "Dec",
];

/// Arbitrary, purely so the year pager is not an infinite hole. Mirrors
/// desktop's `MIN_YEAR` / `MAX_YEAR_OFFSET`.
const _minYear = 1970;
const _maxYearOffset = 5;

/// What the month sheet came back with.
///
/// "All time" is a distinct outcome from "a month", and both are distinct from
/// a dismissal (a null future) -- the same three the desktop panel expresses as
/// its two separate outputs plus a backdrop click.
class MonthPickerResult {
  const MonthPickerResult.month(FilterMonth this.month) : isAllTime = false;

  const MonthPickerResult.allTime()
      : month = null,
        isAllTime = true;

  final FilterMonth? month;

  final bool isAllTime;
}

/// Opens the month sheet, returning null when it is dismissed unchanged.
///
/// A plain compact sheet rather than `showFullscreenBottomSheet`: the content is
/// a fixed-height grid, and that helper mounts a `TopAppBar` over a full-height
/// body for children that scroll themselves. `receipt_form.dart` is the
/// in-repo precedent for a compact one.
Future<MonthPickerResult?> showReceiptMonthPickerSheet(
  BuildContext context, {
  FilterMonth? selected,
}) {
  final sheetContext =
      Provider.of<ContextModel>(context, listen: false).resolveSheetContext(context);

  return showModalBottomSheet<MonthPickerResult>(
    context: sheetContext,
    useSafeArea: true,
    showDragHandle: true,
    // The default sheet caps itself at 9/16 of the screen, which the grid plus
    // the pager and shortcuts exceeds on a short phone -- and a bottom sheet
    // overflows rather than scrolling. Scroll-controlled it sizes to its
    // content instead, and the body scrolls if even that does not fit.
    isScrollControlled: true,
    builder: (context) => ReceiptMonthPicker(selected: selected),
  );
}

/// The sheet's body: a year pager, a grid of the twelve months, and the three
/// shortcuts desktop's panel offers.
class ReceiptMonthPicker extends StatefulWidget {
  const ReceiptMonthPicker({super.key, this.selected});

  final FilterMonth? selected;

  @override
  State<ReceiptMonthPicker> createState() => _ReceiptMonthPickerState();
}

class _ReceiptMonthPickerState extends State<ReceiptMonthPicker> {
  /// The year the grid is paging through. Seeded from the selection, then owned
  /// by the user -- paging is a view concern and selects nothing, so it emits
  /// nothing and leaves the sheet open.
  late int _pagedYear = widget.selected?.year ?? DateTime.now().year;

  bool get _canPageBack => _pagedYear > _minYear;

  bool get _canPageForward => _pagedYear < DateTime.now().year + _maxYearOffset;

  void _pageYear(int delta) => setState(() => _pagedYear += delta);

  void _pick(FilterMonth month) {
    Navigator.of(context).pop(MonthPickerResult.month(month));
  }

  /// `0` is this month, `-1` last month.
  void _pickRelativeMonth(int delta) {
    _pick(shiftMonth(monthOfDate(DateTime.now()), delta));
  }

  @override
  Widget build(BuildContext context) {
    final colors = Theme.of(context).colorScheme;

    return SingleChildScrollView(
      padding: const EdgeInsets.fromLTRB(16, 0, 16, 24),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              IconButton(
                key: const ValueKey("month-picker-year-prev"),
                icon: const Icon(Icons.chevron_left),
                tooltip: "Previous year",
                onPressed: _canPageBack ? () => _pageYear(-1) : null,
              ),
              Text(
                "$_pagedYear",
                key: const ValueKey("month-picker-year"),
                style: Theme.of(context).textTheme.titleMedium,
              ),
              IconButton(
                key: const ValueKey("month-picker-year-next"),
                icon: const Icon(Icons.chevron_right),
                tooltip: "Next year",
                onPressed: _canPageForward ? () => _pageYear(1) : null,
              ),
            ],
          ),
          const SizedBox(height: 8),
          GridView.count(
            crossAxisCount: 4,
            shrinkWrap: true,
            physics: const NeverScrollableScrollPhysics(),
            childAspectRatio: 2,
            mainAxisSpacing: 8,
            crossAxisSpacing: 8,
            children: List.generate(12, (index) {
              final month = index + 1;
              // Only within its own year: paging to 2025 must not leave
              // September 2026 looking picked.
              final isSelected = widget.selected?.year == _pagedYear &&
                  widget.selected?.month == month;

              return TextButton(
                key: ValueKey("month-picker-month-$month"),
                style: TextButton.styleFrom(
                  shape: const StadiumBorder(),
                  backgroundColor:
                      isSelected ? colors.primaryContainer : colors.surfaceDim,
                  foregroundColor: isSelected
                      ? colors.onPrimaryContainer
                      : colors.onSurface,
                ),
                onPressed: () => _pick(FilterMonth(_pagedYear, month)),
                child: Text(_shortMonthLabels[index]),
              );
            }),
          ),
          const SizedBox(height: 16),
          Row(
            children: [
              Expanded(
                child: TextButton(
                  key: const ValueKey("month-picker-this-month"),
                  onPressed: () => _pickRelativeMonth(0),
                  child: const Text("This month"),
                ),
              ),
              Expanded(
                child: TextButton(
                  key: const ValueKey("month-picker-last-month"),
                  onPressed: () => _pickRelativeMonth(-1),
                  child: const Text("Last month"),
                ),
              ),
              Expanded(
                child: TextButton(
                  key: const ValueKey("month-picker-all-time"),
                  style: TextButton.styleFrom(
                      foregroundColor: colors.onSurfaceVariant),
                  onPressed: () => Navigator.of(context)
                      .pop(const MonthPickerResult.allTime()),
                  child: const Text("All time"),
                ),
              ),
            ],
          ),
        ],
      ),
    );
  }
}
