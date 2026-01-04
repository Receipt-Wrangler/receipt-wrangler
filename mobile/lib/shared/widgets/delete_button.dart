import 'package:flutter/material.dart';

class DeleteButton extends StatefulWidget {
  const DeleteButton({
    super.key,
    this.onPressed,
  });

  final void Function()? onPressed;

  @override
  _DeleteButtonState createState() {
    return _DeleteButtonState();
  }
}

class _DeleteButtonState extends State<DeleteButton> {
  @override
  void initState() {
    super.initState();
  }

  @override
  void dispose() {
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return IconButton(
        icon: const Icon(Icons.delete),
        color: Colors.red,
        onPressed: widget.onPressed);
  }
}
