import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;

import '../../utils/currency.dart';
import '../../utils/receipt_summary.dart';
import '../../utils/receipts.dart';

/// Width of the frozen label column. Wide enough for "Needs Attention Receipts" over
/// two lines; the label ellipsizes rather than overflowing, which is the trap
/// `ListItemTrailingStatus` already fell into with a fixed 100px box.
const _labelColumnWidth = 148.0;
const _rowHeight = 40.0;
const _headerHeight = 24.0;
const _figureColumnWidth = 116.0;

/// The fraction of the viewport the bar may take before it scrolls internally. Without
/// a cap, a group breaking out five statuses on a short phone starves the Expanded list
/// beside it and the Column overflows.
const _maxHeightFraction = 0.35;

/// The block of totals pinned above or below the receipts list: a receipt count and
/// amount total over the whole current filter, then the same figures per configured
/// status, plus a column per configured CURRENCY custom field.
///
/// Presentational -- it fetches nothing and knows nothing about the filter.
/// `GroupReceiptsList` owns the request and hands the result down, exactly as
/// `receipts-table` does for the desktop's `app-receipt-totals`.
///
/// Named *Bar*, not *Summary*: `api.ReceiptSummary` is the model it takes, and
/// `group_summary.dart` already owns the word on the dashboards tab -- the same
/// collision the desktop dodged by calling its component `app-receipt-totals`.
class ReceiptSummaryBar extends StatefulWidget {
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

  @override
  State<ReceiptSummaryBar> createState() => _ReceiptSummaryBarState();
}

class _ReceiptSummaryBarState extends State<ReceiptSummaryBar> {
  /// ONE horizontal controller for the whole figure grid.
  ///
  /// Load-bearing: a SingleChildScrollView per row would let the rows desync under a
  /// drag and park a value under the wrong column heading. The labels sit outside this
  /// view, so they stay put while the figures move.
  final ScrollController _figuresController = ScrollController();

  @override
  void dispose() {
    _figuresController.dispose();
    super.dispose();
  }

  /// The overall row first, then one per configured status. Flattened here rather than
  /// rendered through two loops so the row markup exists once.
  List<api.ReceiptSummaryRow> get _rows => [
        widget.summary.overall,
        ...widget.summary.statuses,
      ];

  /// The column headings: the amount total, then one per configured currency field.
  /// Read off the OVERALL row, which always carries every configured field -- a status
  /// row that matched nothing still carries them zeroed, but the overall row is the one
  /// guaranteed to exist.
  List<String> get _columnNames => [
        'Total',
        ...widget.summary.overall.customFieldTotals.map((total) => total.name),
      ];

  /// `formatCurrency` parses the wire string with `double.parse`, which THROWS on
  /// anything non-numeric -- and from inside this bar that would take down the whole
  /// receipts screen, not just the block. Fall back to the raw text instead.
  String _money(String amount) {
    if (double.tryParse(amount) == null) {
      return amount;
    }
    return formatCurrency(context, amount) ?? amount;
  }

  /// A row that matched nothing is muted rather than dropped, so the block keeps its
  /// shape as the filter narrows and a legitimate zero does not read as a bug.
  ///
  /// The mute is a COLOUR, not an Opacity widget: the frozen-label split would need two
  /// of them, and two layers to keep in step is a drift waiting to happen.
  /// `onSurfaceVariant` is the theme's own muted role.
  TextStyle _rowStyle(ThemeData theme, {required bool isOverall, required bool isEmpty}) {
    final base = theme.textTheme.bodyMedium ?? const TextStyle();
    return base.copyWith(
      fontWeight: isOverall ? FontWeight.w600 : FontWeight.w400,
      color: isEmpty ? theme.colorScheme.onSurfaceVariant : theme.colorScheme.onSurface,
      // Desktop's font-variant-numeric: tabular-nums, so the columns line up.
      fontFeatures: const [FontFeature.tabularFigures()],
    );
  }

  Widget _buildConfigRow(ThemeData theme) {
    // A chip row with one option is not a choice; name the group instead. Desktop's rule.
    if (widget.configGroups.length == 1) {
      return Padding(
        key: const ValueKey('receipt-summary-config-note'),
        padding: const EdgeInsets.only(bottom: 4),
        child: Text(
          'Using the summary configuration from ${widget.configGroups.first.name}.',
          style: theme.textTheme.bodySmall
              ?.copyWith(color: theme.colorScheme.onSurfaceVariant),
        ),
      );
    }

    return SizedBox(
      key: const ValueKey('receipt-summary-config-chips'),
      height: 40,
      child: ListView.separated(
        scrollDirection: Axis.horizontal,
        itemCount: widget.configGroups.length,
        separatorBuilder: (_, __) => const SizedBox(width: 6),
        itemBuilder: (context, index) {
          final group = widget.configGroups[index];
          return ChoiceChip(
            key: ValueKey('receipt-summary-config-group-${group.id}'),
            label: Text(group.name),
            selected: group.id == widget.selectedConfigGroupId,
            onSelected: (_) => widget.onConfigGroupSelected?.call(group.id),
          );
        },
      ),
    );
  }

  Widget _buildLabelColumn(ThemeData theme) {
    // Hoisted: _rows SYNTHESISES a list on every read, so reading it in a loop condition
    // and again per field allocates one list per access.
    final rows = _rows;

    return SizedBox(
      width: _labelColumnWidth,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Aligns the labels with the figure rows, past the column headings.
          const SizedBox(height: _headerHeight),
          for (var i = 0; i < rows.length; i++)
            SizedBox(
              key: ValueKey(
                  'receipt-summary-row-${summaryRowKey(rows[i], isOverall: i == 0)}'),
              height: _rowHeight,
              child: Align(
                alignment: Alignment.centerLeft,
                child: Text(
                  '${receiptSummaryRowLabel(rows[i], receiptStatusLabel, isOverall: i == 0)}'
                  '\n(${receiptSummaryCountLabel(rows[i])})',
                  maxLines: 2,
                  overflow: TextOverflow.ellipsis,
                  style: _rowStyle(
                    theme,
                    isOverall: i == 0,
                    isEmpty: rows[i].receiptCount == 0,
                  ),
                ),
              ),
            ),
        ],
      ),
    );
  }

  Widget _buildFigureGrid(ThemeData theme) {
    final names = _columnNames;
    final rows = _rows;

    return Scrollbar(
      controller: _figuresController,
      // Nothing else hints there is more to the right on a phone.
      thumbVisibility: true,
      child: SingleChildScrollView(
        key: const ValueKey('receipt-summary-figures-scroll'),
        controller: _figuresController,
        scrollDirection: Axis.horizontal,
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            SizedBox(
              height: _headerHeight,
              child: Row(
                children: [
                  for (final name in names)
                    SizedBox(
                      width: _figureColumnWidth,
                      child: Text(
                        name,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        textAlign: TextAlign.right,
                        style: theme.textTheme.bodySmall
                            ?.copyWith(color: theme.colorScheme.onSurfaceVariant),
                      ),
                    ),
                ],
              ),
            ),
            for (var i = 0; i < rows.length; i++)
              SizedBox(
                height: _rowHeight,
                child: Row(
                  children: _buildFigureCells(theme, rows[i], isOverall: i == 0),
                ),
              ),
          ],
        ),
      ),
    );
  }

  List<Widget> _buildFigureCells(
    ThemeData theme,
    api.ReceiptSummaryRow row, {
    required bool isOverall,
  }) {
    final rowKey = summaryRowKey(row, isOverall: isOverall);
    final style = _rowStyle(theme, isOverall: isOverall, isEmpty: row.receiptCount == 0);

    Widget cell(String key, String text) => SizedBox(
          key: ValueKey(key),
          width: _figureColumnWidth,
          child: Align(
            alignment: Alignment.centerRight,
            child: Text(text, maxLines: 1, overflow: TextOverflow.ellipsis, style: style),
          ),
        );

    return [
      cell('receipt-summary-figure-$rowKey-total', _money(row.total)),
      for (final customFieldTotal in row.customFieldTotals)
        cell(
          'receipt-summary-figure-$rowKey-cf-${customFieldTotal.customFieldId}',
          _money(customFieldTotal.total),
        ),
    ];
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final divider = BorderSide(color: theme.colorScheme.outlineVariant);

    return Container(
      key: const ValueKey('receipt-summary'),
      width: double.infinity,
      decoration: BoxDecoration(
        color: theme.colorScheme.surfaceContainer,
        // The divider goes on whichever edge faces the list.
        border: Border(
          bottom: widget.atTop ? divider : BorderSide.none,
          top: widget.atTop ? BorderSide.none : divider,
        ),
      ),
      // Its own horizontal inset: the receipts route zeroes the ScreenWrapper body
      // padding so rows sit edge to edge, and the bar's background must reach the edges
      // too or it reads as a list row rather than a pinned bar.
      padding: const EdgeInsets.fromLTRB(16, 8, 16, 8),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          if (widget.configGroups.isNotEmpty) _buildConfigRow(theme),
          ConstrainedBox(
            constraints: BoxConstraints(
              maxHeight: MediaQuery.sizeOf(context).height * _maxHeightFraction,
            ),
            child: SingleChildScrollView(
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  _buildLabelColumn(theme),
                  Expanded(child: _buildFigureGrid(theme)),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }
}
