import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/shared/functions/quick_scan.dart';

import '../../models/auth_model.dart';

List<PopupMenuItem> buildPopupMenuItems(BuildContext context) {
  var authModel = Provider.of<AuthModel>(context, listen: false);
  List<PopupMenuItem> items = [
    PopupMenuItem(
      value: 0,
      child: Text("Add Manual Receipt"),
    ),
  ];

  if (authModel.featureConfig.aiPoweredReceipts) {
    items.add(
      PopupMenuItem(
        value: 1,
        child: Text("Quick Scan"),
        onTap: () => showQuickScanBottomSheet(context),
      ),
    );
  }

  return items;
}

void showAddMenu(BuildContext context, GlobalKey addButtonKey) {
  final RenderBox renderBox =
      addButtonKey.currentContext?.findRenderObject() as RenderBox;
  renderBox.globalToLocal(Offset.zero);
  final Offset offset = renderBox.localToGlobal(Offset.zero);
  final Size size = renderBox.size;

  final RelativeRect position = RelativeRect.fromLTRB(
    offset.dx,
    offset.dy,
    offset.dx + size.width,
    offset.dy + size.height,
  );

  showMenu(context: context, position: position, items: [
    PopupMenuItem(
      value: 0,
      child: Text("Add Manual Receipt"),
      onTap: () => context.go("/receipts/add"),
    ),
    PopupMenuItem(
      value: 1,
      child: Text("Quick Scan"),
      onTap: () => showQuickScanBottomSheet(context),
    ),
  ]);
}
