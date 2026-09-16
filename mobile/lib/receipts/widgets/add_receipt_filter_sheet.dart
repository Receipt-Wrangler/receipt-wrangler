import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/constants/colors.dart';
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
    // The list scrolls itself, so the sheet must hand it bounded constraints
    // rather than wrapping it in a SingleChildScrollView -- see
    // showFullscreenBottomSheet's doc comment.
    bodyFillsSheet: true,
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

    return ListView(
      shrinkWrap: false,
      children: fields
          .map((field) => ListTile(
                key: ValueKey("add-receipt-filter-${field.key}"),
                // accentBlueDark, not the theme's primary: at 22px on white
                // the primary is 2.2:1 and the glyph washes out.
                leading: Icon(field.icon, color: accentBlueDark),
                title: Text(field.label),
                subtitle: Text(field.hint),
                trailing: const Icon(Icons.add, size: 20, color: slate400),
                onTap: () => Navigator.of(context).pop(field.key),
              ))
          .toList(),
    );
  }
}
