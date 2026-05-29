import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/shared/functions/quick_scan.dart';

import '../../models/auth_model.dart';

void showAddMenu(BuildContext context, GlobalKey addButtonKey) {
  final RenderBox renderBox =
      addButtonKey.currentContext?.findRenderObject() as RenderBox;
  final Offset offset = renderBox.localToGlobal(Offset.zero);
  final Size size = renderBox.size;

  final RelativeRect position = RelativeRect.fromLTRB(
    offset.dx,
    offset.dy,
    offset.dx + size.width,
    offset.dy + size.height,
  );

  final authModel = Provider.of<AuthModel>(context, listen: false);
  final items = <PopupMenuItem>[
    PopupMenuItem(
      value: 0,
      child: const Text("Add Manual Receipt"),
      onTap: () => context.go("/receipts/add"),
    ),
    if (authModel.featureConfig.aiPoweredReceipts)
      PopupMenuItem(
        value: 1,
        child: const Text("Quick Scan"),
        onTap: () => showQuickScanBottomSheet(context),
      ),
  ];

  showMenu(context: context, position: position, items: items);
}
