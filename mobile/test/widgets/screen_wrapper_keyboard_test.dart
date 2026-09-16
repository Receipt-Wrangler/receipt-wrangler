import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/screen_wrapper.dart';

/// Keyboard contract for [ScreenWrapper]'s two bottom slots.
///
/// `Scaffold` lifts its `bottomSheet` above the keyboard for free — that slot
/// lands at `contentBottom`, which already subtracts `viewInsets.bottom`. It
/// does **not** lift `bottomNavigationBar`, which it pins to the physical
/// bottom of the scaffold. So the slot `showFullscreenBottomSheet` picks with
/// `bodyFillsSheet: true` — chosen precisely *because* it reserves its space
/// rather than floating over the form's last field — was the one an open
/// keyboard buried. `ScreenWrapper._liftAboveKeyboard` closes that gap.
///
/// These assert **geometry, not finders**, for the reason `mobile/CLAUDE.md`
/// gives throughout: `findsOneWidget` is true for a widget parked off-screen
/// under the keyboard, so only a rect can tell the fix from the bug.
///
/// Every "keyboard" case here **fails on the pre-fix tree** (the bar's bottom
/// edge reads the full screen height, 844, instead of 544). The keyboard-down
/// cases are the controls: without them a bar that had simply been shoved
/// off-screen for some unrelated reason would satisfy the guards.
void main() {
  // A phone, with dpr 1.0 so FakeViewPadding's physical pixels are also the
  // logical ones the rects come back in. Mirrors filter_multiselect_test.
  const screen = Size(390, 844);
  const keyboard = 300.0;
  const barHeight = 80.0;
  const sheetHeight = 56.0;

  const barKey = ValueKey('bottom-bar');
  const sheetKey = ValueKey('bottom-sheet');

  Future<void> pumpWrapper(
    WidgetTester tester, {
    bool withBar = true,
    bool withSheet = false,
    double keyboardHeight = 0,
  }) async {
    tester.view.physicalSize = screen;
    tester.view.devicePixelRatio = 1.0;
    tester.view.viewInsets = FakeViewPadding(bottom: keyboardHeight);
    addTearDown(tester.view.reset);

    await tester.pumpWidget(
      MaterialApp(
        home: ScreenWrapper(
          bottomNavigationBarWidget: withBar
              ? const SizedBox(key: barKey, height: barHeight, width: double.infinity)
              : null,
          bottomSheetWidget: withSheet
              ? const SizedBox(key: sheetKey, height: sheetHeight, width: double.infinity)
              : null,
          child: const SizedBox.expand(),
        ),
      ),
    );
    await tester.pump();
  }

  testWidgets('with no keyboard the bar sits on the bottom of the screen',
      (tester) async {
    await pumpWrapper(tester);

    expect(tester.getRect(find.byKey(barKey)).bottom, screen.height,
        reason: 'the control: nothing to avoid, so nothing is padded');
  });

  testWidgets('an open keyboard pushes the bar clear of it', (tester) async {
    await pumpWrapper(tester, keyboardHeight: keyboard);

    final bar = tester.getRect(find.byKey(barKey));

    expect(bar.bottom, lessThanOrEqualTo(screen.height - keyboard),
        reason: 'the bar must end at or above the top of the keyboard; '
            'pre-fix it is pinned at ${screen.height}');
    // Clear of the keyboard but not floated off into the middle of the screen.
    expect(bar.bottom, screen.height - keyboard);
    expect(bar.height, barHeight, reason: 'lifted, not squashed');
  });

  testWidgets('the bottomSheet slot is left alone — Scaffold already lifts it',
      (tester) async {
    await pumpWrapper(tester,
        withBar: false, withSheet: true, keyboardHeight: keyboard);

    // Exactly on the keyboard, not a bar's height above it: padding this slot
    // as well would double-count the inset.
    expect(tester.getRect(find.byKey(sheetKey)).bottom, screen.height - keyboard);
  });

  testWidgets('a screen using both slots stacks sheet, bar, keyboard',
      (tester) async {
    // This is the /search shell's shape: WranglerSearchBar in the bottomSheet
    // slot, the group nav bar in the bottomNavigationBar slot.
    await pumpWrapper(tester,
        withBar: true, withSheet: true, keyboardHeight: keyboard);

    final bar = tester.getRect(find.byKey(barKey));
    final sheet = tester.getRect(find.byKey(sheetKey));

    expect(bar.bottom, screen.height - keyboard);
    expect(sheet.bottom, lessThanOrEqualTo(bar.top),
        reason: 'the search field now rides above the nav bar rather than '
            'sitting directly on the keyboard with the bar hidden behind it');
    expect(sheet.height, sheetHeight);
  });

  testWidgets('the bar returns to the bottom when the keyboard closes',
      (tester) async {
    await pumpWrapper(tester, keyboardHeight: keyboard);
    expect(tester.getRect(find.byKey(barKey)).bottom, screen.height - keyboard);

    // The engine reports the inset shrinking back as the keyboard dismisses;
    // the padding has to follow it back down rather than latching.
    tester.view.viewInsets = FakeViewPadding.zero;
    await tester.pump();

    expect(tester.getRect(find.byKey(barKey)).bottom, screen.height);
  });
}
