import 'package:flutter/material.dart';

class ScreenWrapper extends StatefulWidget {
  const ScreenWrapper({
    super.key,
    required this.child,
    this.bottomNavigationBarWidget,
    this.appBarWidget,
    this.bodyPadding,
    this.bottomSheetWidget,
  });

  final Widget child;
  final Widget? bottomNavigationBarWidget;
  final PreferredSizeWidget? appBarWidget;
  final EdgeInsets? bodyPadding;
  final Widget? bottomSheetWidget;

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
