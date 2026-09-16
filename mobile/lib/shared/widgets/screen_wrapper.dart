import 'package:flutter/material.dart';

class ScreenWrapper extends StatefulWidget {
  const ScreenWrapper({
    super.key,
    required this.child,
    this.bottomNavigationBarWidget,
    this.appBarWidget,
    this.bodyPadding,
    this.bottomSheetWidget,
    this.backgroundColor,
  });

  final Widget child;
  final Widget? bottomNavigationBarWidget;
  final PreferredSizeWidget? appBarWidget;
  final EdgeInsets? bodyPadding;
  final Widget? bottomSheetWidget;

  /// Overrides the scaffold's own background. Null keeps `Scaffold`'s default,
  /// so every existing caller is unchanged; a screen passes a value when its
  /// content is meant to read as raised cards on a canvas.
  final Color? backgroundColor;

  @override
  State<ScreenWrapper> createState() => _ScreenWrapper();
}

class _ScreenWrapper extends State<ScreenWrapper> {
  @override
  void initState() {
    super.initState();
  }

  @override
  Widget build(BuildContext context) {
    return SafeArea(
      bottom: true,
      top: false,
      child: Scaffold(
        backgroundColor: widget.backgroundColor,
        appBar: widget.appBarWidget,
        bottomSheet: widget.bottomSheetWidget,
        bottomNavigationBar: widget.bottomNavigationBarWidget,
        body: Container(
          padding:
              widget.bodyPadding ?? const EdgeInsets.only(left: 16, right: 16),
          width: MediaQuery.of(context).size.width,
          child: widget.child,
        ),
      ),
    );
  }
}
