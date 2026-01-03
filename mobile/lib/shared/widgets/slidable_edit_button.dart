import 'package:flutter/material.dart';
import 'package:flutter_slidable/flutter_slidable.dart';

class SlidableEditButton extends StatefulWidget {
  const SlidableEditButton({super.key, required this.onPressed});

  final void Function() onPressed;

  @override
  State<SlidableEditButton> createState() => _SlidableEditButton();
}

class _SlidableEditButton extends State<SlidableEditButton> {
  @override
  Widget build(BuildContext context) {
    return SlidableAction(
      icon: Icons.edit,
      label: "Edit",
      foregroundColor: Theme.of(context).colorScheme.primary,
      onPressed: (BuildContext context) {
        widget.onPressed();
      },
    );
  }
}
