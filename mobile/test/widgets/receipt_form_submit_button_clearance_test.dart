import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';

import '../helpers/receipt_form_test_helpers.dart';

/// The receipt form's tail must stay clear of its own submit button.
///
/// The button is `ScreenWrapper.bottomSheetWidget` → `Scaffold.bottomSheet`,
/// which buys keyboard avoidance (it lands at `contentBottom`, already net of
/// `viewInsets.bottom`) at the cost of **floating over** the body rather than
/// reserving space. `SingleChildScrollView` at `maxScrollExtent` only brings
/// content flush with the viewport's bottom edge, so anything in the last 50px
/// is unreachable at *every* scroll offset — not merely awkward to get to.
/// `submitButtonSpacing` at the end of the column is what reserves it.
///
/// Assert **geometry, not finders**: `findsOneWidget` is true for a field
/// parked underneath a floating button, and `ensureVisible` would drag it into
/// the viewport even where a user could not. Same lesson as the Quick Scan
/// sheet and the multi-select sheet — see `mobile/CLAUDE.md`.
///
/// This is only testable because the `kDebugMode` "Check form value" button
/// that used to end this column is gone. It contributed 48px in debug and
/// nothing in release, so the debug tree cleared the button while shipped
/// builds did not — the suite could not see the bug at all. (That button was
/// also drawn entirely underneath the submit button, so it was unclickable.)
void main() {
  // Deliberately wider than a phone. The bug under test is purely vertical,
  // and at 390px the details header's Row overflows horizontally under
  // `--use-test-fonts`, whose stub glyphs are wider than Raleway -- a harness
  // artifact that would otherwise mask the assertions behind a layout error.
  const screen = Size(800, 844);

  Future<void> pumpAtBottom(WidgetTester tester) async {
    tester.view.physicalSize = screen;
    tester.view.devicePixelRatio = 1.0;
    addTearDown(tester.view.reset);

    await pumpReceiptForm(
      tester,
      groups: [buildGroup(id: 1, name: 'Household')],
      // The only shape that reproduces this: a bare Scaffold has no floating
      // button to be buried under.
      pinnedSubmitButton: true,
    );

    final scrollable = tester.state<ScrollableState>(find.byType(Scrollable).first);
    scrollable.position.jumpTo(scrollable.position.maxScrollExtent);
    await tester.pumpAndSettle();
  }

  /// The last real content in the column — the "Shared With" decorator, whose
  /// bottom border and content padding are what got buried.
  Finder lastField() => find
      .ancestor(
        of: find.text('Shared With'),
        matching: find.byType(InputDecorator),
      )
      .first;

  testWidgets('the last field clears the submit button at max scroll',
      (tester) async {
    await pumpAtBottom(tester);

    final button = tester.getRect(find.byType(BottomSubmitButton));
    final field = tester.getRect(lastField());

    // The assertion that fails without `submitButtonSpacing`: the field's
    // bottom edge lands 30px inside the button's band.
    expect(field.bottom, lessThanOrEqualTo(button.top),
        reason: 'the tail of the form must be scrollable clear of the button '
            'that floats over it');
    expect(field.overlaps(button), isFalse);
  });

  testWidgets('the whole field is inside the window, not just above the button',
      (tester) async {
    await pumpAtBottom(tester);

    final field = tester.getRect(lastField());

    // Guards the other direction: a fix that over-reserved, or a body that
    // stopped scrolling, would push the field off the top instead.
    expect(field.top, greaterThanOrEqualTo(0));
    expect(field.bottom, lessThanOrEqualTo(screen.height));
  });

  testWidgets('the button itself still sits on the bottom of the screen',
      (tester) async {
    await pumpAtBottom(tester);

    // The control. Without it the clearance assertion above would also pass on
    // a form whose button had simply drifted up the screen for some unrelated
    // reason, which would prove nothing about reachability.
    expect(tester.getRect(find.byType(BottomSubmitButton)).bottom, screen.height);
  });
}
