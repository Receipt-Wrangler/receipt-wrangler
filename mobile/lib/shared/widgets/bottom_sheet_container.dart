import 'package:flutter/material.dart';

class BottomSheetContainer extends StatefulWidget {
  const BottomSheetContainer(
      {super.key, required this.child, required this.header, this.actions});

  final Widget child;
  final String header;
  final List<Widget>? actions;

  @override
  State<BottomSheetContainer> createState() => _BottomSheetContainerState();
}

class _BottomSheetContainerState extends State<BottomSheetContainer> {
  @override
  Widget build(BuildContext context) {
    return CustomScrollView(
      slivers: [
        SliverAppBar(
          pinned: true,
          centerTitle: false,
          automaticallyImplyLeading: false,
          surfaceTintColor: Colors.white,
          floating: false,
          actions: widget.actions,
          title: Text(widget.header),
        ),
        SliverToBoxAdapter(
          child: widget.child,
        ),
      ],
    );
  }
}
