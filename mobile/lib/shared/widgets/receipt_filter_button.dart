import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/receipt-list-model.dart';
import 'package:receipt_wrangler_mobile/receipts/screens/receipt_filter_screen.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';

/// The receipts app bar's filter action, badged with how many conditions are
/// narrowing the list.
///
/// It awaits nothing: the model is read and the route pushed synchronously, so
/// this never uses its own BuildContext after a gap. That matters here
/// specifically -- an app-bar action that awaits can find itself unmounted (see
/// mobile/CLAUDE.md, "`TopAppBar`'s loading bar is always mounted"). The applied
/// filter travels back through ReceiptListModel rather than as a route result,
/// so there is nothing to await for.
class ReceiptFilterButton extends StatelessWidget {
  const ReceiptFilterButton({super.key});

  @override
  Widget build(BuildContext context) {
    // listen: true -- the badge has to follow setFilter.
    final count = context.watch<ReceiptListModel>().activeFilterCount;
    final groupId = getGroupId(context);

    return Badge(
      isLabelVisible: count > 0,
      label: Text("$count"),
      child: IconButton(
        key: const ValueKey("receipt-filter-button"),
        icon: const Icon(Icons.filter_alt),
        tooltip: "Filter receipts",
        onPressed: () => showReceiptFilterScreen(context, groupId),
      ),
    );
  }
}
