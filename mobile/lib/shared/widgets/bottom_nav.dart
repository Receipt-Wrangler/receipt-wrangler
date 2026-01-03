import 'dart:async';

import 'package:flutter/material.dart';

class BottomNav extends StatefulWidget {
  const BottomNav({
    super.key,
    required this.destinations,
    required this.getInitialSelectedIndex,
    required this.onDestinationSelected,
    required this.indexSelectedController,
  });

  final List<NavigationDestination> destinations;

  final void Function(int) onDestinationSelected;

  final int Function() getInitialSelectedIndex;

  final StreamController<int> indexSelectedController;

  @override
  State<BottomNav> createState() => _BottomNav();
}

class _BottomNav extends State<BottomNav> {
  var indexSelected = 0;

  @override
  void initState() {
    super.initState();

    widget.indexSelectedController.stream.listen((index) {
      setState(() {
        indexSelected = index;
      });
    });
  }

  Widget buildNavigationBar() {
    return NavigationBar(
      destinations: widget.destinations,
      onDestinationSelected: widget.onDestinationSelected,
      selectedIndex: widget.getInitialSelectedIndex(),
    );
  }

  @override
  Widget build(BuildContext context) {
    return buildNavigationBar();
  }
}
