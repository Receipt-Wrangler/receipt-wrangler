import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/list_item_color_block.dart';

import 'date_block.dart';

class ListItemLead extends StatefulWidget {
  const ListItemLead({
    super.key,
    required this.date,
    required this.color,
  });
  final Color color;

  final String date;

  @override
  State<ListItemLead> createState() => _ListItemLead();
}

class _ListItemLead extends State<ListItemLead> {
  @override
  Widget build(BuildContext context) {
    return SizedBox(
      width: 50,
      child: Row(
        mainAxisAlignment: MainAxisAlignment.start,
        children: [
          ListItemColorBlock(color: widget.color),
          const SizedBox(width: 10),
          DateBlock(date: widget.date),
        ],
      ),
    );
  }
}
