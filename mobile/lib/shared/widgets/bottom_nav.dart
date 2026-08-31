import 'dart:async';

import 'package:flutter/material.dart';

/// A bottom-nav destination together with the identity its owner switches on.
///
/// Destinations are permission-gated, so a slot can be missing and every index
/// after it shifts. Carrying an [id] lets the owning nav resolve "which
/// destination was tapped" and "which one matches the current route" by identity
/// instead of by position, which is what makes hiding a middle slot safe.
@immutable
class NavDestinationItem {
  const NavDestinationItem({
    required this.id,
    required this.destination,
    this.onLongPress,
  });

  /// Stable identifier, unique within one nav.
  final String id;

  final NavigationDestination destination;

  /// Optional long-press action for this slot. Only the Scan slot uses it.
  final VoidCallback? onLongPress;
}

class BottomNav extends StatefulWidget {
  const BottomNav({
    super.key,
    required this.items,
    required this.getInitialSelectedIndex,
    required this.onDestinationSelected,
    required this.indexSelectedController,
  });

  final List<NavDestinationItem> items;

  final void Function(int) onDestinationSelected;

  final int Function() getInitialSelectedIndex;

  final StreamController<int> indexSelectedController;

  @override
  State<BottomNav> createState() => _BottomNav();
}

class _BottomNav extends State<BottomNav> {
  var indexSelected = 0;

  /// Held so it can be cancelled: the controller is owned by the caller and
  /// outlives this State, so an uncancelled subscription both leaks and calls
  /// setState after dispose.
  StreamSubscription<int>? _indexSelectedSubscription;

  @override
  void initState() {
    super.initState();

    _indexSelectedSubscription =
        widget.indexSelectedController.stream.listen((index) {
      setState(() {
        indexSelected = index;
      });
    });
  }

  @override
  void dispose() {
    _indexSelectedSubscription?.cancel();
    super.dispose();
  }

  Widget buildNavigationBar() {
    return NavigationBar(
      destinations: widget.items.map((item) => item.destination).toList(),
      onDestinationSelected: widget.onDestinationSelected,
      selectedIndex: widget.getInitialSelectedIndex(),
    );
  }

  /// Transparent long-press targets laid over the bar, one slice per
  /// destination, so a slot can respond to a hold as well as a tap.
  ///
  /// `NavigationBar` has no long-press API and lays its destinations out as
  /// equal-flex `Expanded` children, so an equal-flex `Row` on top lines up with
  /// them. `HitTestBehavior.translucent` keeps the bar underneath hit-testable:
  /// a short tap loses the gesture arena to the bar's own `InkWell` and selects
  /// the destination as usual, while a hold is claimed by the long-press
  /// recognizer at the 500ms mark before the tap can complete. Slices without a
  /// long-press action stay inert so they never interfere.
  ///
  /// The action runs on **`onLongPressUp`**, not `onLongPress`. `onLongPress`
  /// fires at the 500ms mark with the finger still down, and these actions open
  /// a route — pushing one mid-gesture leaves the recognizer holding the pointer
  /// it never saw released, after which that slice ignores every later tap and
  /// hold. (Observed as a slot that worked once and then went dead;
  /// bottom_nav_test covers it.) Waiting for the release costs nothing: the
  /// recognizer has already won the arena at 500ms, so the tap underneath stays
  /// suppressed either way.
  Widget buildLongPressOverlay() {
    return Row(
      children: widget.items.map((item) {
        final onLongPress = item.onLongPress;
        return Expanded(
          child: onLongPress == null
              ? const SizedBox.expand()
              : GestureDetector(
                  behavior: HitTestBehavior.translucent,
                  onLongPressUp: onLongPress,
                  child: const SizedBox.expand(),
                ),
        );
      }).toList(),
    );
  }

  @override
  Widget build(BuildContext context) {
    // Material's NavigationBar asserts on fewer than two destinations, and
    // permission gating can get us there: on group-select, a user who can
    // neither add receipts nor search is left with "Groups" alone. A one-item
    // nav offers nothing anyway -- the user is already on that screen -- so
    // render no bar rather than crash.
    if (widget.items.length < 2) {
      return const SizedBox.shrink();
    }

    final hasLongPress =
        widget.items.any((item) => item.onLongPress != null);
    if (!hasLongPress) {
      return buildNavigationBar();
    }

    return Stack(
      children: [
        buildNavigationBar(),
        Positioned.fill(child: buildLongPressOverlay()),
      ],
    );
  }
}
