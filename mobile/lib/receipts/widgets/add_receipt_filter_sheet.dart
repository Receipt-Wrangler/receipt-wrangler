import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_filter_fields.dart';
import 'package:receipt_wrangler_mobile/utils/bottom_sheet.dart';

/// Offers the fields not already in the draft, and returns the chosen key.
///
/// A field the view cannot use is never listed rather than listed and then
/// ignored -- see [unusedReceiptFilterFields] for the `Group` rule.
Future<String?> showAddReceiptFilterSheet(
  BuildContext context, {
  required Set<String> usedKeys,
  required bool showGroupField,
}) async {
  final available = unusedReceiptFilterFields(usedKeys, showGroupField: showGroupField);

  final result = await showFullscreenBottomSheet(
    context,
    AddReceiptFilterList(fields: available),
    "Add a filter",
  );

  return result is String ? result : null;
}

/// The fields [showAddReceiptFilterSheet] offers, in [receiptFilterFields] order.
///
/// `group` is only offered on the synthetic "All" group, mirroring desktop's
/// `showGroupFilter`. Inside a real group the receipts endpoint already scopes
/// every query to it, so a group condition there is redundant at best and
/// contradictory at worst.
List<ReceiptFilterField> unusedReceiptFilterFields(
  Set<String> usedKeys, {
  required bool showGroupField,
}) {
  return receiptFilterFields
      .where((field) => !usedKeys.contains(field.key))
      .where((field) => field.key != "group" || showGroupField)
      .toList();
}

class AddReceiptFilterList extends StatelessWidget {
  const AddReceiptFilterList({super.key, required this.fields});

  final List<ReceiptFilterField> fields;

  @override
  Widget build(BuildContext context) {
    if (fields.isEmpty) {
      return const Padding(
        key: ValueKey("add-receipt-filter-empty"),
        padding: EdgeInsets.symmetric(vertical: 36, horizontal: 16),
        child: Text("Every field is already filtered.",
            textAlign: TextAlign.center),
      );
    }

    return Column(
      children: fields
          .map((field) => ListTile(
                key: ValueKey("add-receipt-filter-${field.key}"),
                leading:
                    Icon(field.icon, color: Theme.of(context).primaryColor),
                title: Text(field.label),
                subtitle: Text(field.hint),
                trailing: const Icon(Icons.add),
                onTap: () => Navigator.of(context).pop(field.key),
              ))
          .toList(),
    );
  }
}
