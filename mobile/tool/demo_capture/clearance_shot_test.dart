import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:image/image.dart' as img;

import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';

import '../../test/helpers/receipt_form_test_helpers.dart';
import 'capture.dart';

/// Captures **one** panel of the receipt-form clearance still.
///
/// Not part of the test suite — `flutter test` with no arguments scans only
/// `test/`, so CI never runs this. Driven by `tool/record_clearance_shot.sh`,
/// which invokes it twice (the second time around a temporary revert of
/// `submitButtonSpacing`) and stitches the two panels:
///
/// ```bash
/// cd mobile && ./tool/record_clearance_shot.sh
/// ```
///
/// Unlike the keyboard demo there is no `debugDisable…` seam to flip: no
/// production flag guards the spacer, and adding one for a screenshot would not
/// be worth it. The runner reverts the source instead, under a `trap`.
///
/// `$CLEARANCE_SHOT_OUT` is the panel's output path; `$CLEARANCE_SHOT_LABEL`
/// picks the caption.
void main() {
  /// `BottomSubmitButton` is a fixed 50px and this Scaffold has no bottom bar,
  /// so the floating button covers exactly the last 50px of the viewport.
  const buttonHeight = 50.0;

  const beforeColor = Color(0xFFB4232B);
  const afterColor = Color(0xFF15803D);

  testWidgets('receipt form at max scroll', (tester) async {
    final out =
        Platform.environment['CLEARANCE_SHOT_OUT'] ?? '/tmp/clearance.png';
    final isBefore =
        (Platform.environment['CLEARANCE_SHOT_LABEL'] ?? 'BEFORE') == 'BEFORE';

    addTearDown(tester.view.reset);
    await tester.runAsync(loadDemoFonts);

    tester.view.physicalSize = const Size(demoPhoneWidth, demoPanelHeight);
    tester.view.devicePixelRatio = 1.0;
    tester.view.viewInsets = FakeViewPadding.zero;

    await pumpReceiptForm(
      tester,
      groups: [buildGroup(id: 1, name: 'Household')],
      // The shape the bug lives in: the submit button in the floating
      // `Scaffold.bottomSheet` slot. A bare Scaffold has nothing to bury under.
      pinnedSubmitButton: true,
      showDebugBanner: false,
      wrap: (app) => buildDemoSurface(
        label: isBefore
            ? 'BEFORE — tail under the button'
            : 'AFTER — tail scrolls clear',
        labelColor: isBefore ? beforeColor : afterColor,
        app: Stack(
          children: [
            Positioned.fill(child: app),
            // Without this the burial is nearly invisible: what gets covered is
            // the bottom border and content padding of the "Shared With"
            // decorator, not anything carrying text.
            Positioned(
              left: 0,
              right: 0,
              bottom: buttonHeight,
              child: IgnorePointer(
                child: Row(
                  crossAxisAlignment: CrossAxisAlignment.end,
                  children: [
                    const Expanded(
                      child: SizedBox(
                        height: 3,
                        child: ColoredBox(color: Color(0xFFE5484D)),
                      ),
                    ),
                    Container(
                      color: const Color(0xFFE5484D),
                      padding: const EdgeInsets.symmetric(
                          horizontal: 6, vertical: 2),
                      child: const Text(
                        'top of Submit button',
                        style: TextStyle(
                          fontFamily: 'Raleway',
                          color: Colors.white,
                          fontSize: 11,
                          fontWeight: FontWeight.w600,
                        ),
                      ),
                    ),
                  ],
                ),
              ),
            ),
          ],
        ),
      ),
    );

    final scrollable =
        tester.state<ScrollableState>(find.byType(Scrollable).first);
    scrollable.position.jumpTo(scrollable.position.maxScrollExtent);
    await tester.pumpAndSettle();

    // Report the geometry the still is meant to show, so a panel that failed to
    // reproduce the bug is obvious in the run output rather than only in the
    // image.
    final button = tester.getRect(find.byType(BottomSubmitButton));
    final field = tester.getRect(find
        .ancestor(of: find.text('Shared With'), matching: find.byType(InputDecorator))
        .first);
    // ignore: avoid_print
    print('panel=${isBefore ? "BEFORE" : "AFTER"} '
        'fieldBottom=${field.bottom} buttonTop=${button.top} '
        'buried=${(field.bottom - button.top).clamp(0, double.infinity)}px');

    final frame = await grabFrame(tester);
    File(out).writeAsBytesSync(img.encodePng(frame));
    // ignore: avoid_print
    print('wrote $out (${frame.width}x${frame.height})');
  }, timeout: const Timeout(Duration(minutes: 3)));
}
