import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_filter_fields.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

/// One authored condition on the filter screen: what it narrows, how, and to
/// what. Tapping anywhere reopens it; the X drops it.
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

    return Card(
      margin: const EdgeInsets.only(bottom: 10),
      child: InkWell(
        onTap: onTap,
        child: Padding(
          padding: const EdgeInsets.fromLTRB(14, 10, 6, 12),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Row(
                children: [
                  Icon(field.icon, size: 20, color: theme.primaryColor),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(field.label,
                        style: theme.textTheme.titleSmall,
                        overflow: TextOverflow.ellipsis),
                  ),
                  IconButton(
                    key: ValueKey("receipt-filter-remove-${field.key}"),
                    icon: const Icon(Icons.close, size: 18),
                    tooltip: "Remove ${field.label} condition",
                    onPressed: onRemove,
                  ),
                ],
              ),
              Row(
                children: [
                  Chip(
                    label: Text(
                      filterOperationLabels[condition.operation] ??
                          condition.operation.name,
                      style: theme.textTheme.labelSmall,
                    ),
                    visualDensity: VisualDensity.compact,
                    materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
                  ),
                  const SizedBox(width: 8),
                  Expanded(
                    child: Text(
                      receiptFilterValueLabel(field, condition),
                      overflow: TextOverflow.ellipsis,
                      style: theme.textTheme.bodyMedium,
                    ),
                  ),
                  const Icon(Icons.chevron_right, size: 20),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }
}
