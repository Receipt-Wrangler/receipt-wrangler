import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_date_filter.dart';

/// The quick date control's stepper: previous / next arrows either side of a
/// label that opens the month sheet, plus a clear back to all time.
///
/// ```
/// ◄   📅 September 2026   ✕   ►
/// ```
///
/// Deliberately presentational, like desktop's `app-month-stepper`: it is told
/// which month to show and what to call it, and emits. It knows nothing about
/// filters, which date field it is pointed at, or how a month is encoded.
class ReceiptMonthStepper extends StatelessWidget {
  const ReceiptMonthStepper({
    super.key,
    required this.month,
    required this.label,
    required this.hasCondition,
    required this.onMonthSelected,
    required this.onAllTimeSelected,
    required this.onLabelPressed,
  });

  /// The month being shown, or null when the label is saying something a month
  /// cannot express ("All time", "Custom").
  final FilterMonth? month;

  /// What the label reads. The caller owns the wording because only it knows
  /// whether "no month" means nothing is set or something it cannot describe.
  final String label;

  /// Whether the target field carries a condition at all, which is what the
  /// clear button is for. Not the same as [month] being set: a range the
  /// stepper cannot describe is still a filter the user needs a way out of.
  final bool hasCondition;

  final ValueChanged<FilterMonth> onMonthSelected;

  final VoidCallback onAllTimeSelected;

  final VoidCallback onLabelPressed;

  /// Steps by whole months. With no month showing it steps from the current
  /// month, so the first press always lands somewhere useful rather than
  /// no-oping -- and `◄` and `►` never do the same thing.
  void _step(int delta) {
    onMonthSelected(shiftMonth(month ?? monthOfDate(DateTime.now()), delta));
  }

  @override
  Widget build(BuildContext context) {
    final colors = Theme.of(context).colorScheme;

    return Row(
      children: [
        IconButton(
          key: const ValueKey("receipt-month-prev"),
          icon: const Icon(Icons.chevron_left),
          tooltip: "Previous month",
          onPressed: () => _step(-1),
        ),
        Expanded(
          child: TextButton.icon(
            key: const ValueKey("receipt-month-label"),
            icon: const Icon(Icons.calendar_month, size: 19),
            label: Text(
              label,
              overflow: TextOverflow.ellipsis,
              style: const TextStyle(fontWeight: FontWeight.w600),
            ),
            style: TextButton.styleFrom(
              minimumSize: const Size.fromHeight(44),
              shape: const StadiumBorder(),
              // Tinted once it is narrowing anything, so an active quick filter
              // is visible without reading the label.
              backgroundColor:
                  hasCondition ? colors.primaryContainer : Colors.transparent,
              foregroundColor: hasCondition
                  ? colors.onPrimaryContainer
                  : colors.onSurfaceVariant,
            ),
            onPressed: onLabelPressed,
          ),
        ),
        if (hasCondition)
          IconButton(
            key: const ValueKey("receipt-month-clear"),
            icon: const Icon(Icons.close, size: 20),
            tooltip: "Back to all time",
            color: colors.onSurfaceVariant,
            onPressed: onAllTimeSelected,
          ),
        IconButton(
          key: const ValueKey("receipt-month-next"),
          icon: const Icon(Icons.chevron_right),
          tooltip: "Next month",
          onPressed: () => _step(1),
        ),
      ],
    );
  }
}
