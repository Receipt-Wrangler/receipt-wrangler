import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../enums/form_state.dart';
import '../../models/permissions_model.dart';
import '../../utils/forms.dart';
import '../functions/permissions.dart';

class ReceiptEditPopupMenu extends StatelessWidget {
  const ReceiptEditPopupMenu(
      {super.key, required this.groupId, required this.popupMenuChildren, this.formState});

  final int groupId;

  final List<PopupMenuEntry> popupMenuChildren;

  final WranglerFormState? formState;

  @override
  Widget build(BuildContext context) {
    late final permissionsModel =
        Provider.of<PermissionsModel>(context, listen: false);
    late final formStateToUse = formState ?? getFormStateFromContext(context);
    var canEdit = canEditReceipt(permissionsModel, groupId);

    if (formStateToUse == WranglerFormState.add) {
      canEdit = true;
    }

    if (canEdit) {
      return PopupMenuButton(
        itemBuilder: (BuildContext context) {
          return popupMenuChildren;
        },
      );
    } else {
      return const SizedBox.shrink();
    }
  }
}
