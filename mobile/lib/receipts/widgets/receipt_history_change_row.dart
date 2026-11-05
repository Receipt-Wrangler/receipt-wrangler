import 'package:flutter/material.dart';

import '../../constants/colors.dart';
import '../../utils/receipt_history.dart';

/// One field-level change: its label above the old value (red, struck through),
/// an arrow, and the new value (green). ADDED and REMOVED rows show one side.
class ReceiptHistoryChangeRow extends StatelessWidget {
  const ReceiptHistoryChangeRow({super.key, required this.row});

  final HistoryChangeRow row;

  static const fromKey = ValueKey('receipt-history-from');

  static const toKey = ValueKey('receipt-history-to');

  @override
  Widget build(BuildContext context) {
    return Container(
      width: double.infinity,
      padding: const EdgeInsets.symmetric(vertical: 8),
      decoration: const BoxDecoration(
        border: Border(top: BorderSide(color: slate100)),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            row.label,
            style: const TextStyle(
              fontSize: 12,
              fontWeight: FontWeight.w500,
              color: slate500,
            ),
          ),
          const SizedBox(height: 4),
          Wrap(
            spacing: 6,
            runSpacing: 4,
            crossAxisAlignment: WrapCrossAlignment.center,
            children: [
              if (row.from != null)
                _Chip(
                  key: fromKey,
                  text: row.from!,
                  surface: historyRemovedSurface,
                  foreground: historyRemovedText,
                  struckThrough: true,
                ),
              if (row.arrow)
                const Icon(Icons.arrow_forward, size: 16, color: slate400),
              if (row.to != null)
                _Chip(
                  key: toKey,
                  text: row.to!,
                  surface: historyAddedSurface,
                  foreground: historyAddedText,
                ),
            ],
          ),
        ],
      ),
    );
  }
}

class _Chip extends StatelessWidget {
  const _Chip({
    super.key,
    required this.text,
    required this.surface,
    required this.foreground,
    this.struckThrough = false,
  });

  final String text;
  final Color surface;
  final Color foreground;
  final bool struckThrough;

  @override
  Widget build(BuildContext context) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: surface,
        borderRadius: BorderRadius.circular(6),
      ),
      child: Text(
        text,
        style: TextStyle(
          fontSize: 14,
          height: 1.45,
          color: foreground,
          decoration: struckThrough ? TextDecoration.lineThrough : null,
          decorationColor: foreground.withValues(alpha: 0.5),
        ),
      ),
    );
  }
}
