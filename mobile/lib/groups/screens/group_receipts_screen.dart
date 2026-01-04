import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/group_receipts_list.dart';

class GroupReceiptsScreen extends StatefulWidget {
  const GroupReceiptsScreen({super.key});

  @override
  State<GroupReceiptsScreen> createState() => _GroupReceiptsScreen();
}

class _GroupReceiptsScreen extends State<GroupReceiptsScreen> {
  @override
  Widget build(BuildContext context) {
    return const GroupReceiptsList();
  }
}
