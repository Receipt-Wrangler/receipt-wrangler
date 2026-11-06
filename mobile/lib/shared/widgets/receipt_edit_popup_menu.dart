import 'package:flutter/material.dart';
import 'package:provider/provider.dart';

import '../../enums/form_state.dart';
import '../../models/permissions_model.dart';
import '../../utils/forms.dart';
import '../functions/permissions.dart';

/// The ⋮ menu on a receipt's view screens.
///
/// [popupMenuChildren] are the edit actions, shown only to a user who may edit
/// the receipt (or on the add form). [alwaysShownChildren] are entries the
/// caller has gated itself -- the receipt view's "History" -- shown whatever the
/// edit permission, so a read-only member who may see the history still gets a
/// menu holding only that. With nothing to show, the widget renders no button
/// at all.
///
/// Integration specs wait for this widget to mean "the /view shell has
/// mounted", so it stays on the view screens only; the edit screen's menu is a
/// plain `PopupMenuButton` (see `ReceiptAppBarActionBuilder`).
class ReceiptEditPopupMenu extends StatelessWidget {
  const ReceiptEditPopupMenu({
    super.key,
    required this.groupId,
    required this.popupMenuChildren,
    this.alwaysShownChildren = const [],
    this.formState,
  });

  final int groupId;

  final List<PopupMenuEntry> popupMenuChildren;

  final List<PopupMenuEntry> alwaysShownChildren;

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

    final children = [
      if (canEdit) ...popupMenuChildren,
      ...alwaysShownChildren,
    ];

    if (children.isEmpty) {
      return const SizedBox.shrink();
    }

    return PopupMenuButton(
      itemBuilder: (BuildContext context) {
        return children;
      },
    );
  }
}
