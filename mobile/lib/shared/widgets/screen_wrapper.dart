import 'package:flutter/material.dart';

/// Keeps a pinned bottom bar clear of the software keyboard.
///
/// `Scaffold` lifts its `bottomSheet` slot for free -- that one lands at
/// `contentBottom`, which already subtracts `viewInsets.bottom` -- but it pins
/// `bottomNavigationBar` to the physical bottom of the scaffold
/// (`max(0, height - bottomWidgetsHeight)`, with no inset term). So the slot
/// chosen precisely *because* it reserves its space is the one an open keyboard
/// buries, which is how the Quick Scan submit button ended up under it.
///
/// The value is not computed or guessed: `viewInsets.bottom` *is* the top of the
/// keyboard, reported by the platform and updated every frame while it slides,
/// so the bar follows it exactly -- including split, floating and hardware
/// keyboards.
///
/// The [Builder] is load-bearing. Widgets are inflated where they are mounted,
/// not where they are constructed, so its context lands inside the bar slot --
/// which is the only place the value survives. `Scaffold` strips the bottom
/// inset from the **body** slot's `MediaQuery` (`removeBottomInset`) but not
/// from this one, so reading it any higher, or from inside the body, yields 0
/// and this silently does nothing. It also scopes the rebuild to the bar.
///
/// Nothing double-counts against the [SafeArea] below: that consumes
/// `padding.bottom`, which is `viewPadding - viewInsets` floored at 0 and so is
/// reported as 0 for exactly as long as the keyboard is up. The two are
/// complementary, never additive.
///
/// Plain [Padding], deliberately not `AnimatedPadding`: the scaffold body
/// resizes un-animated off the same value, so animating here would put the bar
/// out of step with the body it has to sit against.
Widget _liftAboveKeyboard(Widget child) {
  return Builder(
    builder: (BuildContext context) => Padding(
      padding: EdgeInsets.only(
        bottom: debugDisableKeyboardLift
            ? 0
            : MediaQuery.viewInsetsOf(context).bottom,
      ),
      child: child,
    ),
  );
}

/// Turns [_liftAboveKeyboard] back off, reproducing the pre-fix layout.
///
/// Exists for the demo harness in `tool/keyboard_demo/`, which records the
/// before/after pair by running the **real** screens twice with only this
/// flipped -- so the "before" side is the production tree rather than a copy of
/// it that could drift. Mirrors the existing debug seams (`debugCameraAccessOverride`,
/// `debugForcePermissionDenied`); the regression tests deliberately do **not**
/// use it, since a test driven by this flag would only be testing the flag.
@visibleForTesting
bool debugDisableKeyboardLift = false;

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
        // Not lifted: Scaffold already places its bottomSheet above the
        // keyboard, so padding it here would double-count.
        bottomSheet: widget.bottomSheetWidget,
        bottomNavigationBar: widget.bottomNavigationBarWidget == null
            ? null
            : _liftAboveKeyboard(widget.bottomNavigationBarWidget!),
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
