import 'package:flutter/material.dart';
import 'package:flutter_slidable/flutter_slidable.dart';

class SlidableDeleteButton extends StatefulWidget {
  const SlidableDeleteButton({super.key, required this.onPressed});

  final void Function() onPressed;

  @override
  State<SlidableDeleteButton> createState() => _SlidableDeleteButton();
}

class _SlidableDeleteButton extends State<SlidableDeleteButton> {
  @override
  Widget build(BuildContext context) {
    return SlidableAction(
        onPressed: (BuildContext context) {
          widget.onPressed();
        },
        icon: Icons.delete,
        foregroundColor: Colors.red,
        label: "Delete");
  }
}
