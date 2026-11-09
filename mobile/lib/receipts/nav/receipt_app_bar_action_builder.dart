import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';

import '../../models/permissions_model.dart';
import '../../models/receipt_model.dart';
import '../../shared/functions/permissions.dart';
import '../../shared/widgets/receipt_edit_popup_menu.dart';
import '../screens/receipt_history_screen.dart';

/// The receipt form screen's app-bar menu.
///
/// This used to also carry image and comment menus, selected by testing
/// `state.fullPath` for "images" / "comments". No such route exists -- the full
/// set is `/receipts/add`, `/receipts/:id/view` and `/receipts/:id/edit` (plus
/// the group, profile, reports and search routes) -- and the image screen is
/// pushed as a plain `MaterialPageRoute` from `receipt_form.dart`, so the path
/// never changes when it opens. Both branches were unreachable, and the image
/// one was a stale copy of `ReceiptImageAppBar` that navigated to a URL
/// (`/receipts/:id/images/edit`) the router does not define. They were deleted
/// rather than merged; `ReceiptImageAppBar` is the live implementation.
class ReceiptAppBarActionBuilder {
  late final ReceiptModel receiptModel;

  late final BuildContext context;

  ReceiptAppBarActionBuilder(BuildContext context, ReceiptModel receiptModel) {
    this.context = context;
    this.receiptModel = receiptModel;
  }

  static const historyMenuItemKey = ValueKey('receipt-history-menu-item');

  static const editOverflowMenuKey = ValueKey('receipt-edit-overflow-menu');

  List<Widget> buildAppBarMenu(GoRouterState state) {
    return [
      _buildReceiptAppBarMenu(
          state.fullPath!, state.pathParameters['receiptId']),
    ];
  }

  Widget _buildReceiptAppBarMenu(String fullPath, String? routeReceiptId) {
    // The app bar renders before the receipt loads, while the model still
    // holds the default receipt (or the previous screen's). History is only
    // offered once it holds the receipt this route names: earlier it would
    // open the wrong receipt's history, and it would make the menu appear
    // before the receipt loaded, which specs read as "loaded".
    // Judged in the group the receipt is saved in (`receipt`, not the edit
    // form's working copy), as the server does.
    final receiptLoaded =
        receiptModel.receipt.id.toString() == routeReceiptId;
    final canViewHistory = receiptLoaded &&
        canViewReceiptHistory(
            Provider.of<PermissionsModel>(context, listen: false),
            receiptModel.receipt.groupId);

    if (fullPath.contains("view")) {
      return ReceiptEditPopupMenu(
          groupId: receiptModel.receipt.groupId,
          popupMenuChildren: [
            PopupMenuItem(
              child: Text("Edit"),
              onTap: () => _goAfterMenu(
                  "/receipts/${receiptModel.receipt.id}/edit"),
            )
          ],
          alwaysShownChildren: [
            if (canViewHistory) _buildHistoryMenuItem(),
          ]);
    }

    // Deliberately a plain PopupMenuButton, not ReceiptEditPopupMenu:
    // integration specs wait for that widget to mean "the /view shell has
    // mounted", which this screen must not satisfy.
    if (fullPath.contains("edit") && canViewHistory) {
      return PopupMenuButton(
        key: editOverflowMenuKey,
        itemBuilder: (context) => [_buildHistoryMenuItem()],
      );
    }

    return SizedBox.shrink();
  }

  PopupMenuItem _buildHistoryMenuItem() {
    return PopupMenuItem(
      key: historyMenuItemKey,
      child: Text("History"),
      onTap: _pushHistoryAfterMenu,
    );
  }

  // A PopupMenuItem invokes onTap *after* it pops the menu route. Navigating
  // straight from the (now-departing) menu/page context is lost under the
  // receipt routes' NoTransitionPage swap, which removes the old page in the
  // same frame. Capture the stable GoRouter up front and navigate once the
  // menu has finished dismissing so the go isn't clobbered by the pop.
  void _goAfterMenu(String location) {
    final router = GoRouter.of(context);
    Future.delayed(
        const Duration(milliseconds: 50), () => router.go(location));
  }

  // History is pushed onto the navigator rather than routed: the form derives
  // view/edit from the live URL (getFormState), so a /history route would flip
  // an open edit form underneath to view. Pushing leaves the URL alone, as the
  // Comments and Images sub-screens do, and the edit form keeps its state.
  void _pushHistoryAfterMenu() {
    final navigator = Navigator.of(context);
    final receiptId = receiptModel.receipt.id;
    Future.delayed(
      const Duration(milliseconds: 50),
      () => navigator.push(MaterialPageRoute(
        builder: (_) => ReceiptHistoryScreen(receiptId: receiptId),
      )),
    );
  }
}
