import 'package:flutter/material.dart';

class ListItemColorBlock extends StatefulWidget {
  const ListItemColorBlock({
    super.key,
    required this.color,
  });

  final Color color;

  @override
  _ListItemColorBlockState createState() {
    return _ListItemColorBlockState();
  }
}

class _ListItemColorBlockState extends State<ListItemColorBlock> {
  @override
  Widget build(BuildContext context) {
    return Container(
      width: 10,
      height: 50,
      color: widget.color,
    );
  }
}
