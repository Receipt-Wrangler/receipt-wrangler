import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/top_app_bar.dart';

class GroupSelectAppBar extends StatefulWidget implements PreferredSizeWidget {
  const GroupSelectAppBar({super.key});

  @override
  State<GroupSelectAppBar> createState() => _GroupSelectAppBar();

  @override
  Size get preferredSize => AppBar().preferredSize;
}

class _GroupSelectAppBar extends State<GroupSelectAppBar> {
  @override
  Widget build(BuildContext context) {
    return const TopAppBar(
      titleText: "Groups",
    );
  }
}
