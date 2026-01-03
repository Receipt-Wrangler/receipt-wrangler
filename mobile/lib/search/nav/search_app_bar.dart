import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/top_app_bar.dart';

class SearchAppBar extends StatefulWidget implements PreferredSizeWidget {
  const SearchAppBar({super.key});

  @override
  State<SearchAppBar> createState() => _SearchAppBar();

  @override
  Size get preferredSize => AppBar().preferredSize;
}

class _SearchAppBar extends State<SearchAppBar> {
  @override
  Widget build(BuildContext context) {
    return const TopAppBar(
      titleText: "Search",
    );
  }
}
