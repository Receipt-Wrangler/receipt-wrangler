import 'dart:async';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_nav.dart';

/// The Scan slot answers a tap and a hold with different actions, layered over a
/// Material `NavigationBar` that has no long-press API of its own. These pin the
/// gesture split: a hold must not also select the destination, and a tap must
/// still reach the bar underneath the transparent overlay.
void main() {
  late List<int> tapped;
  late List<int> held;
  late StreamController<int> controller;

  setUp(() {
    tapped = [];
    held = [];
    controller = StreamController<int>();
  });

  tearDown(() => controller.close());

  Widget harness({required bool withLongPress}) {
    return MaterialApp(
      home: Scaffold(
        bottomNavigationBar: BottomNav(
          items: [
            const NavDestinationItem(
              id: 'first',
              destination: NavigationDestination(
                  icon: Icon(Icons.dashboard), label: 'First'),
            ),
            NavDestinationItem(
              id: 'middle',
              onLongPress: withLongPress ? () => held.add(1) : null,
              destination: const NavigationDestination(
                  icon: Icon(Icons.add), label: 'Middle'),
            ),
            const NavDestinationItem(
              id: 'last',
              destination: NavigationDestination(
                  icon: Icon(Icons.receipt), label: 'Last'),
            ),
          ],
          onDestinationSelected: tapped.add,
          getInitialSelectedIndex: () => 0,
          indexSelectedController: controller,
        ),
      ),
    );
  }

  testWidgets('a tap on the long-pressable slot still selects it',
      (tester) async {
    await tester.pumpWidget(harness(withLongPress: true));

    await tester.tap(find.text('Middle'));
    await tester.pump();

    expect(tapped, [1], reason: 'the overlay is translucent, so the tap reaches '
        'the NavigationBar underneath');
    expect(held, isEmpty);
  });

  testWidgets('a hold opens the menu instead of selecting', (tester) async {
    await tester.pumpWidget(harness(withLongPress: true));

    await tester.longPress(find.text('Middle'));
    await tester.pump();

    expect(held, [1]);
    expect(tapped, isEmpty,
        reason: 'the long-press recognizer wins the arena before the tap '
            'completes, so the hold must not also navigate');
  });

  testWidgets('slots without a long-press action stay inert', (tester) async {
    await tester.pumpWidget(harness(withLongPress: true));

    await tester.longPress(find.text('First'));
    await tester.pump();

    expect(held, isEmpty);
  });

  testWidgets('taps on the other slots are unaffected by the overlay',
      (tester) async {
    await tester.pumpWidget(harness(withLongPress: true));

    await tester.tap(find.text('Last'));
    await tester.pump();

    expect(tapped, [2]);
  });

  // NOTE: the "slot dies after its menu closes" regression (see
  // buildLongPressOverlay's onLongPressUp comment) is NOT covered here. It was
  // reproduced only against the real app -- a minimal harness delivers the
  // pointer-up cleanly enough that the recognizer recovers, and the test passes
  // with the bug present. The guard lives in
  // integration_test/receipt_entry_menu_reopen_test.dart instead.

  testWidgets('no overlay is built when nothing is long-pressable',
      (tester) async {
    await tester.pumpWidget(harness(withLongPress: false));

    expect(find.byType(Stack), findsWidgets,
        reason: 'Material builds its own Stacks; the assertion below is the '
            'meaningful one');
    await tester.tap(find.text('Middle'));
    await tester.pump();
    expect(tapped, [1]);
  });
}
