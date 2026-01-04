import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';

import '../../models/receipt_model.dart';
import '../../shared/widgets/receipt_edit_popup_menu.dart';

class ReceiptCommentAppBar extends StatelessWidget implements PreferredSizeWidget {
  const ReceiptCommentAppBar({
    super.key,
    required this.formState,
    required this.title,
    required this.onBack,
    this.onEditMode,
  });

  final WranglerFormState formState;
  final String title;
  final VoidCallback onBack;
  final VoidCallback? onEditMode;

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);

  @override
  Widget build(BuildContext context) {
    final receiptModel = Provider.of<ReceiptModel>(context, listen: false);
    
    return AppBar(
      title: Text(title),
      leading: IconButton(
        icon: const Icon(Icons.arrow_back),
        onPressed: onBack,
      ),
      actions: [_buildCommentActions(context, receiptModel)],
    );
  }

  Widget _buildCommentActions(BuildContext context, ReceiptModel receiptModel) {
    if (formState == WranglerFormState.view) {
      return ReceiptEditPopupMenu(
          groupId: receiptModel.receipt.groupId,
          popupMenuChildren: [
            PopupMenuItem(
              child: Text("Edit"),
              onTap: () {
                if (onEditMode != null) {
                  onEditMode!();
                }
              },
            )
          ],
          formState: formState);
    }

    return SizedBox.shrink();
  }
}