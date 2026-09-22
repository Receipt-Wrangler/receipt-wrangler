import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;

import '../../utils/currency.dart';
import '../../utils/receipt_summary.dart';
import '../../utils/receipts.dart';

/// The fraction of the viewport the bar may take before it scrolls internally. Without
/// a cap, a group breaking out five statuses on a short phone starves the Expanded list
/// beside it and the Column overflows.
const _maxHeightFraction = 0.35;

/// Diameter of the status dot, and the gutter its column occupies.
const _dotSize = 10.0;
const _dotGutter = _dotSize + 8;

/// The block of totals pinned above or below the receipts list: a receipt count and
/// amount total over the whole current filter, then the same figures per configured
/// status, plus a figure per configured CURRENCY custom field.
///
/// Presentational -- it fetches nothing and knows nothing about the filter.
/// `GroupReceiptsList` owns the request and hands the result down, exactly as
/// `receipts-table` does for the desktop's `app-receipt-totals`.
///
/// Named *Bar*, not *Summary*: `api.ReceiptSummary` is the model it takes, and
/// `group_summary.dart` already owns the word on the dashboards tab -- the same
/// collision the desktop dodged by calling its component `app-receipt-totals`.
///
/// **Rows WRAP; nothing is a column.** The first version of this widget was a frozen
/// label column beside a horizontally scrolling grid of fixed-width figure columns, and
/// fixed widths cannot fit an unknown number of configured currency fields into 390pt:
/// with one field the second column was already chopped mid-number, and a long status
/// name ate its own receipt count through the label's ellipsis. So a row is a label line
/// plus `name value` pairs that wrap -- the same flex-wrap the desktop's `receipt-totals`
/// uses, and for the same reason. Anything that does not fit moves to the next line
/// instead of disappearing, and a row whose group configures no currency fields -- the
/// common case -- stays on one line.
class ReceiptSummaryBar extends StatelessWidget {
  const ReceiptSummaryBar({
    super.key,
    required this.summary,
    required this.atTop,
    this.configGroups = const [],
    this.selectedConfigGroupId,
    this.onConfigGroupSelected,
  });

  final api.ReceiptSummary summary;

  /// Which edge faces the list, so the divider lands between the two.
  final bool atTop;

  /// Groups whose configuration the viewer may pick between. Non-empty only on the
  /// synthetic "All" group, which spans several groups and has no configuration of its
  /// own; elsewhere the group being viewed decides.
  final List<api.Group> configGroups;

  final int? selectedConfigGroupId;

  final void Function(int groupId)? onConfigGroupSelected;

  /// The overall row first, then one per configured status.
  List<api.ReceiptSummaryRow> get _rows => [summary.overall, ...summary.statuses];

  /// `formatCurrency` parses the wire string with `double.parse`, which THROWS on
  /// anything non-numeric -- and from inside this bar that would take down the whole
  /// receipts screen, not just the block. Fall back to the raw text instead.
  String _money(BuildContext context, String amount) {
    if (double.tryParse(amount) == null) {
      return amount;
    }
    return formatCurrency(context, amount) ?? amount;
  }

  /// A row that matched nothing is muted rather than dropped, so the block keeps its
  /// shape as the filter narrows and a legitimate zero does not read as a bug.
  ///
  /// The mute is a COLOUR, not an Opacity widget: it has to reach the row's label, its
  /// count and every figure, and keeping several Opacity layers in step is a drift
  /// waiting to happen. `onSurfaceVariant` is the theme's own muted role.
  Color _primaryInk(ThemeData theme, {required bool isEmpty}) =>
      isEmpty ? theme.colorScheme.onSurfaceVariant : theme.colorScheme.onSurface;

  Widget _buildConfigRow(BuildContext context, ThemeData theme) {
    // A chip row with one option is not a choice; name the group instead. Desktop's rule.
    if (configGroups.length == 1) {
      return Padding(
        key: const ValueKey('receipt-summary-config-note'),
        padding: const EdgeInsets.only(bottom: 8),
        child: Text(
          'Using the summary configuration from ${configGroups.first.name}.',
          style: theme.textTheme.bodySmall?.copyWith(
            color: theme.colorScheme.onSurfaceVariant,
          ),
        ),
      );
    }

    return Padding(
      key: const ValueKey('receipt-summary-config-chips'),
      padding: const EdgeInsets.only(bottom: 4),
      child: SizedBox(
        height: 40,
        child: ListView.separated(
          scrollDirection: Axis.horizontal,
          itemCount: configGroups.length,
          separatorBuilder: (_, __) => const SizedBox(width: 6),
          itemBuilder: (context, index) {
            final group = configGroups[index];
            return ChoiceChip(
              key: ValueKey('receipt-summary-config-group-${group.id}'),
              label: Text(group.name),
              selected: group.id == selectedConfigGroupId,
              onSelected: (_) => onConfigGroupSelected?.call(group.id),
            );
          },
        ),
      ),
    );
  }

  /// The status tint, ringed.
  ///
  /// `receiptStatusColor` returns a pale *background* tint -- it is designed to sit
  /// behind dark text, which is how `ListItemTrailingStatus` and `ListItemColorBlock`
  /// use it. Painted as a bare dot on white, OPEN (#FFFACD) lands near 1.1:1 and is
  /// effectively invisible, so the ring is what makes the marker readable at every
  /// status while staying quieter than a chip.
  Widget _buildStatusDot(
    ThemeData theme,
    api.ReceiptSummaryRow row, {
    required bool isOverall,
  }) {
    // The overall row is not a status, so it takes an empty gutter of the same width --
    // every label in the block then starts on the same x.
    if (isOverall || row.status == api.ReceiptStatus.empty) {
      return const SizedBox(width: _dotSize, height: _dotSize);
    }

    return Container(
      width: _dotSize,
      height: _dotSize,
      decoration: BoxDecoration(
        color: receiptStatusColor(row.status),
        shape: BoxShape.circle,
        border: Border.all(color: theme.colorScheme.outlineVariant),
      ),
    );
  }

  /// One `name value` pair, baseline-aligned like the desktop's `receipt-totals__figure`.
  Widget _buildFigure(
    ThemeData theme, {
    required Key key,
    required String name,
    required String value,
    required bool isOverall,
    required bool isEmpty,
  }) {
    final ink = _primaryInk(theme, isEmpty: isEmpty);

    return Row(
      key: key,
      mainAxisSize: MainAxisSize.min,
      crossAxisAlignment: CrossAxisAlignment.baseline,
      textBaseline: TextBaseline.alphabetic,
      children: [
        Text(
          name,
          style: theme.textTheme.bodySmall?.copyWith(
            color: theme.colorScheme.onSurfaceVariant,
          ),
        ),
        const SizedBox(width: 6),
        Text(
          value,
          style: theme.textTheme.bodyMedium?.copyWith(
            color: ink,
            fontWeight: isOverall ? FontWeight.w600 : FontWeight.w500,
            // Desktop's font-variant-numeric: tabular-nums, so digits keep a constant
            // advance and the figures still line up where a run of rows happens to
            // wrap the same way.
            fontFeatures: const [FontFeature.tabularFigures()],
          ),
        ),
      ],
    );
  }

  Widget _buildRow(
    BuildContext context,
    ThemeData theme,
    api.ReceiptSummaryRow row, {
    required bool isOverall,
  }) {
    final isEmpty = row.receiptCount == 0;
    final ink = _primaryInk(theme, isEmpty: isEmpty);
    final rowKey = summaryRowKey(row, isOverall: isOverall);

    final label = Text(
      receiptSummaryRowLabel(row, receiptStatusLabel, isOverall: isOverall),
      style: theme.textTheme.bodyMedium?.copyWith(
        color: ink,
        fontWeight: isOverall ? FontWeight.w600 : FontWeight.w500,
      ),
    );

    // The count is its own Text, not glued to the label with a newline: as a separate
    // widget it can take bodySmall and the muted role, which is the work the old
    // parentheses were standing in for -- and it can no longer be eaten by the label's
    // ellipsis, which is how a long status name used to lose its count entirely.
    final count = Text(
      receiptSummaryCountLabel(row),
      style: theme.textTheme.bodySmall?.copyWith(
        color: theme.colorScheme.onSurfaceVariant,
      ),
    );

    final figures = <Widget>[
      _buildFigure(
        theme,
        key: ValueKey('receipt-summary-figure-$rowKey-total'),
        name: 'Total',
        value: _money(context, row.total),
        isOverall: isOverall,
        isEmpty: isEmpty,
      ),
      for (final customFieldTotal in row.customFieldTotals)
        _buildFigure(
          theme,
          key: ValueKey(
            'receipt-summary-figure-$rowKey-cf-${customFieldTotal.customFieldId}',
          ),
          name: customFieldTotal.name,
          value: _money(context, customFieldTotal.total),
          isOverall: isOverall,
          isEmpty: isEmpty,
        ),
    ];

    final hasSingleFigure = figures.length == 1;

    return Padding(
      key: ValueKey('receipt-summary-row-$rowKey'),
      padding: const EdgeInsets.symmetric(vertical: 8),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Nudged onto the first line's optical centre; the Row aligns to the top so
          // the dot stays on the label line rather than drifting to the row's middle.
          Padding(
            padding: const EdgeInsets.only(top: 5),
            child: _buildStatusDot(theme, row, isOverall: isOverall),
          ),
          const SizedBox(width: _dotGutter - _dotSize),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              mainAxisSize: MainAxisSize.min,
              children: [
                // The label line. A Wrap rather than a Row so a lone Total can ride it,
                // pushed right: that is the shape most groups have -- no currency
                // fields -- and giving one number a line of its own wasted half the
                // block's height. An Expanded would look the same until the text does
                // not fit, and then it is a yellow-and-black overflow stripe; with a
                // Wrap the figure drops to the next line at a large text scale or
                // behind a long status name.
                Wrap(
                  alignment: WrapAlignment.spaceBetween,
                  spacing: 12,
                  runSpacing: 2,
                  children: [
                    Row(
                      mainAxisSize: MainAxisSize.min,
                      crossAxisAlignment: CrossAxisAlignment.baseline,
                      textBaseline: TextBaseline.alphabetic,
                      children: [
                        // Flexible, so a long status name wraps inside the label rather
                        // than pushing the count off the end -- which is how "Needs
                        // Attention Receipts" used to lose its count to an ellipsis.
                        Flexible(child: label),
                        const SizedBox(width: 8),
                        count,
                      ],
                    ),
                    if (hasSingleFigure) figures.first,
                  ],
                ),
                // Two or more figures take a line of their own, every row starting at
                // the same x. Deliberately NOT folded into the Wrap above: one Wrap
                // over the label and every figure lets the first figure ride up
                // whenever a label happens to be short, so "Total" begins at a
                // different x on each row and the block reads as ragged. They still
                // wrap among themselves when a group configures enough currency fields
                // to overflow the line -- nothing is ever clipped.
                if (!hasSingleFigure) ...[
                  const SizedBox(height: 2),
                  Wrap(spacing: 16, runSpacing: 4, children: figures),
                ],
              ],
            ),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final divider = BorderSide(color: theme.colorScheme.outlineVariant);
    final rows = _rows;

    return Container(
      key: const ValueKey('receipt-summary'),
      width: double.infinity,
      decoration: BoxDecoration(
        // The band behind the card is the page canvas, so the white card reads as
        // raised against it rather than as another sheet of the same paper.
        color: theme.colorScheme.surfaceDim,
        // The divider goes on whichever edge faces the list.
        border: Border(
          bottom: atTop ? divider : BorderSide.none,
          top: atTop ? BorderSide.none : divider,
        ),
      ),
      // The BAND is full-bleed on purpose: the receipts route zeroes ScreenWrapper's
      // body padding so rows sit edge to edge, and a block inset on all sides reads as
      // a list row rather than as pinned chrome. The card inside it is what carries the
      // inset.
      padding: const EdgeInsets.all(12),
      child: Container(
        // Hand-built rather than a Material Card, matching ReceiptFilterConditionCard:
        // the app's Card is Material's elevated one, whose tinted surface and drop
        // shadow are far heavier than the near-flat surface this wants.
        decoration: BoxDecoration(
          color: theme.colorScheme.surface,
          borderRadius: BorderRadius.circular(14),
          border: Border.all(color: Colors.black.withValues(alpha: 0.06)),
          boxShadow: [
            BoxShadow(
              color: Colors.black.withValues(alpha: 0.05),
              offset: const Offset(0, 1),
              blurRadius: 2,
            ),
          ],
        ),
        padding: const EdgeInsets.symmetric(horizontal: 14, vertical: 2),
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            if (configGroups.isNotEmpty) _buildConfigRow(context, theme),
            ConstrainedBox(
              constraints: BoxConstraints(
                maxHeight: MediaQuery.sizeOf(context).height * _maxHeightFraction,
              ),
              child: _ScrollableRows(
                children: [
                  for (var i = 0; i < rows.length; i++) ...[
                    if (i > 0)
                      Divider(
                        height: 1,
                        thickness: 1,
                        color: theme.colorScheme.outlineVariant,
                      ),
                    _buildRow(context, theme, rows[i], isOverall: i == 0),
                  ],
                ],
              ),
            ),
          ],
        ),
      ),
    );
  }
}

/// The rows, scrolling within the bar's height cap.
///
/// A separate widget purely so the [ScrollController] has somewhere to live and be
/// disposed -- the bar itself is presentational and has no other state.
///
/// The [Scrollbar] earns its place: a group breaking out every status with several
/// currency fields runs past the cap, and a row cut off with no affordance reads as
/// broken rather than as scrollable -- exactly the mistake the horizontally clipped
/// grid this replaced was making. The thumb shows only while the content actually
/// overflows, so the common case shows nothing at all.
///
/// The controller is explicit rather than inherited: `SingleChildScrollView` adopts the
/// `PrimaryScrollController` by default on the vertical axis, and this bar mounts inside
/// a route that already has one.
class _ScrollableRows extends StatefulWidget {
  const _ScrollableRows({required this.children});

  final List<Widget> children;

  @override
  State<_ScrollableRows> createState() => _ScrollableRowsState();
}

class _ScrollableRowsState extends State<_ScrollableRows> {
  final ScrollController _controller = ScrollController();

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scrollbar(
      controller: _controller,
      // Always on rather than fade-on-scroll: the cut-off row is the only other hint
      // that there is more, and a row clipped flush at the card's edge reads as broken.
      // Flutter draws nothing when the content fits, so the common case is unaffected.
      thumbVisibility: true,
      child: SingleChildScrollView(
        controller: _controller,
        child: Column(
          mainAxisSize: MainAxisSize.min,
          crossAxisAlignment: CrossAxisAlignment.start,
          children: widget.children,
        ),
      ),
    );
  }
}
