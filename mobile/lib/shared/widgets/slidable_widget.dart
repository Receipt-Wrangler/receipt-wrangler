import 'package:flutter/material.dart';
import 'package:flutter_slidable/flutter_slidable.dart';

class SlidableWidget extends StatefulWidget {
  const SlidableWidget({
    super.key,
    required this.slideEnabled,
    required this.endActionPaneChildren,
    required this.slidableChild,
  });

  final bool slideEnabled;

  final Widget slidableChild;

  final List<Widget> endActionPaneChildren;

  @override
  State<SlidableWidget> createState() => _SlidableWidget();
}

class _SlidableWidget extends State<SlidableWidget> {
  @override
  Widget build(BuildContext context) {
    return Slidable(
        enabled: widget.slideEnabled,
        endActionPane: ActionPane(
          motion: const DrawerMotion(),
          children: widget.endActionPaneChildren,
        ),
        child: widget.slidableChild);
  }
}
