import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/group_list.dart';

class GroupSelect extends StatefulWidget {
  const GroupSelect({super.key});

  @override
  State<GroupSelect> createState() => _GroupSelect();
}

class _GroupSelect extends State<GroupSelect> {
  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return const GroupList();
  }
}
