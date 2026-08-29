import 'package:flutter/material.dart';

import '../functions/receipt_entry.dart';

/// The receipts screen's overflow menu.
///
/// The Scan slot's extra actions live behind a long-press, which is not
/// discoverable and not reachable by every input method. This exposes the same
/// permission-gated items — they come from [buildReceiptEntryMenuItems], so the
/// two menus cannot drift apart — through an ordinary button.
///
/// Renders nothing when the user can neither quick-scan nor create receipts:
/// an overflow that only ever says "no" is worse than no overflow.
class ReceiptEntryOverflowMenu extends StatelessWidget {
  const ReceiptEntryOverflowMenu({super.key});

  @override
  Widget build(BuildContext context) {
    final items = buildReceiptEntryMenuItems(context);
    if (items.isEmpty) {
      return const SizedBox.shrink();
    }

    return PopupMenuButton(
      key: const ValueKey("receipt-entry-overflow-menu"),
      icon: const Icon(Icons.more_vert),
      tooltip: "Add a receipt",
      itemBuilder: (context) => buildReceiptEntryMenuItems(context),
    );
  }
}
