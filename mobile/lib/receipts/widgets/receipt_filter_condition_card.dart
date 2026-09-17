import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/constants/colors.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_filter_fields.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

/// One authored condition on the filter screen: what it narrows, how, and to
/// what. Tapping anywhere reopens it; the X drops it.
///
/// A hand-built container rather than a [Card]: the app's `Card` is Material's
/// elevated one, whose tinted surface and drop shadow are far heavier than the
/// near-flat row this needs. Deliberately carries no field icon -- the icons
/// belong to the add sheet, where they help pick a field, and repeating them
/// here only competes with the label.
class ReceiptFilterConditionCard extends StatelessWidget {
  const ReceiptFilterConditionCard({
    super.key,
    required this.field,
    required this.condition,
    required this.onTap,
    required this.onRemove,
  });

  final ReceiptFilterField field;

  final ReceiptFilterCondition condition;

  final VoidCallback onTap;

  final VoidCallback onRemove;

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Padding(
      padding: const EdgeInsets.only(bottom: 10),
      child: Material(
        color: theme.colorScheme.surface,
        borderRadius: BorderRadius.circular(14),
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(14),
          child: Ink(
            decoration: BoxDecoration(
              // The fill belongs on the decoration, not only on the Material:
              // a BoxDecoration paints its boxShadow as a silhouette of the
              // whole shape, so without a colour here the shadow shows through
              // the card's interior and greys it out.
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
            padding: const EdgeInsets.fromLTRB(14, 8, 8, 12),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Row(
                  children: [
                    Expanded(
                      child: Text(
                        field.label,
                        style: const TextStyle(
                            fontSize: 14, fontWeight: FontWeight.w600),
                        overflow: TextOverflow.ellipsis,
                      ),
                    ),
                    SizedBox.square(
                      dimension: 32,
                      child: IconButton(
                        key: ValueKey("receipt-filter-remove-${field.key}"),
                        icon: const Icon(Icons.close, size: 18),
                        color: theme.colorScheme.onSurfaceVariant,
                        padding: EdgeInsets.zero,
                        visualDensity: VisualDensity.compact,
                        tooltip: "Remove ${field.label} condition",
                        onPressed: onRemove,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 4),
                Row(
                  children: [
                    Container(
                      height: 30,
                      alignment: Alignment.center,
                      padding: const EdgeInsets.symmetric(horizontal: 12),
                      decoration: BoxDecoration(
                        color: theme.colorScheme.surfaceContainer,
                        borderRadius: BorderRadius.circular(15),
                      ),
                      child: Text(
                        filterOperationLabels[condition.operation] ??
                            condition.operation.name,
                        style: const TextStyle(
                          fontSize: 12,
                          fontWeight: FontWeight.w600,
                          color: slate700,
                        ),
                      ),
                    ),
                    const SizedBox(width: 8),
                    Expanded(
                      child: Text(
                        receiptFilterValueLabel(field, condition),
                        overflow: TextOverflow.ellipsis,
                        style: const TextStyle(fontSize: 14),
                      ),
                    ),
                    Icon(Icons.chevron_right,
                        size: 20, color: theme.colorScheme.onSurfaceVariant),
                  ],
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }
}
