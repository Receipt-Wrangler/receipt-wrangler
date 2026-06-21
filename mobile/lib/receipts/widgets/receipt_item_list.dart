import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/spacing.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_item_items.dart';

import '../../models/receipt_model.dart';
import '../../shared/functions/forms.dart';

class ReceiptItemField extends StatefulWidget {
  const ReceiptItemField({super.key, required this.groupId});

  final int groupId;

  @override
  State<ReceiptItemField> createState() => _ReceiptItemFieldState();
}

class _ReceiptItemFieldState extends State<ReceiptItemField> {
  @override
  Widget build(BuildContext context) {
    return Consumer<ReceiptModel>(
      builder: (context, consumerModel, child) {
        return Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            buildHeaderText("Items"),
            textFieldSpacing,
            ReceiptItemItems(
              items: consumerModel.items,
              groupId: widget.groupId,
            ),
          ],
        );
      },
    );
  }
}
