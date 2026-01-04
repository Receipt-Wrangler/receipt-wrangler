import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;

class GroupListCard extends StatefulWidget {
  const GroupListCard(
      {super.key, required this.group, required this.onGroupTap});

  final api.Group group;

  final void Function(api.Group group) onGroupTap;

  @override
  State<GroupListCard> createState() => _GroupListCard();
}

class _GroupListCard extends State<GroupListCard> {
  @override
  Widget build(BuildContext context) {
    return ListTile(
      title: Text(widget.group.name),
      onTap: () => widget.onGroupTap(widget.group),
    );
  }
}
