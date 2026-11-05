import 'package:flutter/material.dart';

import '../../constants/colors.dart';

/// "Receipt History" with the summary line under it once loaded. A plain
/// `AppBar`, like `ReceiptCommentAppBar`: `TopAppBar` has no subtitle.
class ReceiptHistoryAppBar extends StatelessWidget implements PreferredSizeWidget {
  const ReceiptHistoryAppBar({super.key, required this.onBack, this.summary});

  final VoidCallback onBack;

  final String? summary;

  static const summaryKey = ValueKey('receipt-history-summary');

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);

  @override
  Widget build(BuildContext context) {
    return AppBar(
      leading: IconButton(
        icon: const Icon(Icons.arrow_back),
        tooltip: 'Back',
        onPressed: onBack,
      ),
      title: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        mainAxisSize: MainAxisSize.min,
        children: [
          const Text('Receipt History'),
          if (summary != null)
            Text(
              summary!,
              key: summaryKey,
              style: const TextStyle(fontSize: 12, color: slate500),
            ),
        ],
      ),
    );
  }
}
