import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;

import '../../constants/colors.dart';
import '../../utils/receipt_history.dart';
import 'receipt_history_actor_avatar.dart';
import 'receipt_history_change_row.dart';

/// One history entry: who and when, a failed-save or older-version box when
/// either applies, and its change rows, collapsed after [historyCollapseAfter].
class ReceiptHistoryEntryCard extends StatefulWidget {
  const ReceiptHistoryEntryCard({
    super.key,
    required this.entry,
    required this.formatters,
  });

  final api.ReceiptHistoryEntry entry;

  final HistoryValueFormatters formatters;

  static const failedKey = ValueKey('receipt-history-failed');

  static const noticeKey = ValueKey('receipt-history-notice');

  static const showMoreKey = ValueKey('receipt-history-show-more');

  @override
  State<ReceiptHistoryEntryCard> createState() => _ReceiptHistoryEntryCardState();
}

class _ReceiptHistoryEntryCardState extends State<ReceiptHistoryEntryCard> {
  bool _expanded = false;

  @override
  Widget build(BuildContext context) {
    final entry = widget.entry;
    final actor = historyActorDisplay(entry.actor, widget.formatters.findUser);
    final rows = entry.changes
        .map((change) => buildHistoryChangeRow(change, widget.formatters))
        .toList();
    final collapsible = rows.length > historyCollapseAfter;
    final visibleRows =
        _expanded || !collapsible ? rows : rows.take(historyCollapseAfter).toList();
    final notice = historyBeforeStateNotice(entry);
    final errorMessage = entry.errorMessage;
    final onSurface = Theme.of(context).colorScheme.onSurface;

    return Container(
      padding: EdgeInsets.fromLTRB(16, 14, 16, collapsible ? 0 : 14),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(14),
        border: Border.all(color: Colors.black.withValues(alpha: 0.06)),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 2,
            offset: const Offset(0, 1),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            children: [
              ReceiptHistoryActorAvatar(actor: actor),
              const SizedBox(width: 12),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text.rich(
                      TextSpan(children: [
                        TextSpan(
                          text: actor.name,
                          style: const TextStyle(fontWeight: FontWeight.w600),
                        ),
                        TextSpan(text: ' ${historyActionVerb(entry.action)}'),
                      ]),
                      style: TextStyle(fontSize: 15, height: 1.3, color: onSurface),
                    ),
                    Text(
                      historyTime(entry.recordedAt),
                      style: const TextStyle(fontSize: 12, color: slate500),
                    ),
                  ],
                ),
              ),
            ],
          ),
          if (errorMessage != null && errorMessage.isNotEmpty)
            _Box(
              key: ReceiptHistoryEntryCard.failedKey,
              icon: Icons.error_outline,
              text: 'Save failed: $errorMessage',
              surface: historyRemovedSurface,
              foreground: historyRemovedText,
              fontSize: 13,
            ),
          if (notice != null)
            _Box(
              key: ReceiptHistoryEntryCard.noticeKey,
              icon: Icons.info,
              text: notice,
              surface: accentTint,
              foreground: historyNoticeText,
              fontSize: 12,
            ),
          if (visibleRows.isNotEmpty) const SizedBox(height: 6),
          for (var i = 0; i < visibleRows.length; i++)
            ReceiptHistoryChangeRow(
              key: ValueKey('receipt-history-change-$i'),
              row: visibleRows[i],
            ),
          if (collapsible)
            Container(
              width: double.infinity,
              height: 44,
              margin: const EdgeInsets.only(top: 4),
              decoration: const BoxDecoration(
                border: Border(top: BorderSide(color: slate100)),
              ),
              child: TextButton(
                key: ReceiptHistoryEntryCard.showMoreKey,
                onPressed: () => setState(() => _expanded = !_expanded),
                child: Text(
                  _expanded
                      ? 'Show less'
                      : 'Show ${rows.length - historyCollapseAfter} more',
                  style: const TextStyle(
                    fontSize: 14,
                    fontWeight: FontWeight.w600,
                    color: accentBlueDark,
                  ),
                ),
              ),
            ),
        ],
      ),
    );
  }
}

class _Box extends StatelessWidget {
  const _Box({
    super.key,
    required this.icon,
    required this.text,
    required this.surface,
    required this.foreground,
    required this.fontSize,
  });

  final IconData icon;
  final String text;
  final Color surface;
  final Color foreground;
  final double fontSize;

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      margin: const EdgeInsets.only(top: 12),
      padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 10),
      decoration: BoxDecoration(
        color: surface,
        borderRadius: BorderRadius.circular(10),
      ),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Icon(icon, size: 18, color: foreground),
          const SizedBox(width: 8),
          Expanded(
            child: Text(
              text,
              style: TextStyle(fontSize: fontSize, height: 1.45, color: foreground),
            ),
          ),
        ],
      ),
    );
  }
}
