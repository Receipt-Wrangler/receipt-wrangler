import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/image_viewer.dart';

class FullScreenImageViewer extends StatefulWidget {
  const FullScreenImageViewer({super.key, required this.image});

  final Image image;

  @override
  State<FullScreenImageViewer> createState() => _FullScreenImageViewer();
}

class _FullScreenImageViewer extends State<FullScreenImageViewer> {
  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [ImageViewer(image: widget.image)],
    );
  }
}
