import 'dart:math' as math;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/constants/colors.dart';
import 'package:receipt_wrangler_mobile/theme/app_theme.dart';

/// `ColorScheme` silently falls back for any role it is not given, and in this
/// scheme every fallback lands on pure black or pure white. So a role dropped
/// from `buildAppTheme` does not fail to compile and does not fail any widget
/// test -- it just puts black borders and black "muted" text back across the
/// whole app. These assertions are the only thing that notices.
void main() {
  final theme = buildAppTheme();
  final scheme = theme.colorScheme;

  group("the neutral roles are spelled out, not left to fall back", () {
    test("outline and outlineVariant are slate, not onBackground", () {
      // outlineVariant is $accent-palette 200; outline is the one deliberate
      // half-step off the palette (see borderSlate's doc comment).
      expect(scheme.outline, borderSlate);
      expect(scheme.outlineVariant, slate200);
      expect(scheme.outline, isNot(scheme.onSurface));
      expect(scheme.outline, isNot(slate300));
      expect(scheme.outlineVariant, isNot(scheme.onSurface));
    });

    test("onSurfaceVariant is muted, not a second copy of onSurface", () {
      expect(scheme.onSurfaceVariant, slate500);
      expect(scheme.onSurfaceVariant, isNot(scheme.onSurface));
    });

    test("the surfaceContainer family is tonal, not all one white", () {
      expect(scheme.surfaceContainerLowest, Colors.white);
      expect(scheme.surfaceContainerLow, Colors.white);
      expect(scheme.surfaceContainer, slate100);
      expect(scheme.surfaceContainerHigh, slate100);
      expect(scheme.surfaceContainerHighest, slate200);
      expect(scheme.surfaceDim, slate50);
      expect(scheme.surfaceContainer, isNot(scheme.surface));
      expect(scheme.surfaceContainerHighest, isNot(scheme.surface));
    });

    test("secondaryContainer does not fall through to the #8EA1AC slate", () {
      expect(scheme.secondaryContainer, slate100);
      expect(scheme.secondaryContainer, isNot(scheme.secondary));
      expect(scheme.onSecondaryContainer, slate700);
    });

    test("primaryContainer carries the accent tint and its readable text", () {
      // accentTint is the design's 8% wash; onPrimaryContainer is
      // $primary-palette 700, the web theme's own "darker" stop.
      expect(scheme.primaryContainer, accentTint);
      expect(scheme.onPrimaryContainer, accentBlueDark);
      expect(scheme.primaryContainer, isNot(scheme.primary));
    });
  });

  group("the roles the app already had are unchanged", () {
    test("primary, error and the surface pair", () {
      // Brand palette stops, from desktop/src/variables.scss: primary 500,
      // accent 500 (the old secondary), warn 500.
      expect(scheme.primary, const Color(0xFF27B1FF));
      expect(scheme.onPrimary, const Color(0xFFFFFFFF));
      expect(scheme.secondary, const Color(0xFF8EA1AC));
      expect(scheme.error, const Color(0xFFd63333));
      expect(scheme.surface, const Color(0xFFFFFFFF));
      expect(scheme.onSurface, const Color(0xFF000000));
      expect(scheme.brightness, Brightness.light);
    });

    test("Material 3 and Raleway", () {
      expect(theme.useMaterial3, isTrue);
      expect(theme.textTheme.bodyMedium?.fontFamily, appFontFamily);
    });
  });

  group("a selected chip is readable", () {
    // Material 3 resolves a chip's label colour out of `labelStyle`, never
    // `secondaryLabelStyle`, so the selected colour has to be a WidgetStateColor
    // on that style. Left to the M3 defaults it is `onSecondaryContainer` --
    // dark slate on the #27B1FF fill, which is what the pickers shipped with.
    Color labelColorFor(Set<WidgetState> states) {
      final color = theme.chipTheme.labelStyle!.color!;
      return WidgetStateProperty.resolveAs<Color?>(color, states)!;
    }

    test("selected is onPrimary, over the primary fill", () {
      expect(theme.chipTheme.selectedColor, scheme.primary);
      expect(labelColorFor({WidgetState.selected}), scheme.onPrimary);
    });

    test("unselected is the readable slate", () {
      expect(labelColorFor(const <WidgetState>{}), slate700);
    });

    test("the label style keeps the app font, since it replaces the default",
        () {
      expect(theme.chipTheme.labelStyle?.fontFamily, appFontFamily);
      expect(theme.chipTheme.labelStyle?.fontSize, isNotNull);
    });
  });

  test("the bottom nav's selected pill does not ride secondaryContainer", () {
    // The regression this guards: NavigationBar takes its indicator from
    // `secondaryContainer`. Before that role was named it fell through to the
    // #8EA1AC secondary -- a solid pill. Naming it slate-100 (correct for its
    // other consumers) made the selected destination a near-white pill on a
    // near-white bar, on every main screen. Nothing else notices.
    final nav = theme.navigationBarTheme;

    // accentContainer is $primary-palette 50 -- the same stop desktop pairs
    // with accentBlueDark in roles/role-presets.ts.
    expect(nav.indicatorColor, accentContainer);
    expect(nav.indicatorColor, isNot(scheme.secondaryContainer));
    expect(nav.indicatorColor, isNot(scheme.surface));

    Color iconColorFor(Set<WidgetState> states) =>
        nav.iconTheme!.resolve(states)!.color!;
    Color labelColorFor(Set<WidgetState> states) =>
        nav.labelTextStyle!.resolve(states)!.color!;

    expect(iconColorFor({WidgetState.selected}), accentBlueDark);
    expect(labelColorFor({WidgetState.selected}), accentBlueDark);
    expect(iconColorFor(const <WidgetState>{}), slate700);
    expect(labelColorFor(const <WidgetState>{}), slate700);
  });

  test("an app bar's leading and action icons are the same grey", () {
    // M3 takes leading from onSurface and actions from onSurfaceVariant, so a
    // bar carrying both renders two different colours unless both are pinned.
    expect(theme.appBarTheme.iconTheme?.color, slate700);
    expect(theme.appBarTheme.actionsIconTheme?.color, slate700);
  });

  group("non-text contrast clears WCAG 2.1 SC 1.4.11", () {
    // Measured rather than pinned to hex values, so the assertions keep holding
    // through a palette tweak and fail only when one actually breaks
    // accessibility. The borders started at pure black (21:1), went to the
    // design's #CBD5E1 (1.48:1, a fail) and now sit here.
    double relativeLuminance(Color c) {
      double channel(double v) =>
          v <= 0.03928 ? v / 12.92 : math.pow((v + 0.055) / 1.055, 2.4) as double;
      return 0.2126 * channel(c.r) +
          0.7152 * channel(c.g) +
          0.0722 * channel(c.b);
    }

    double contrastRatio(Color a, Color b) {
      final la = relativeLuminance(a);
      final lb = relativeLuminance(b);
      return (math.max(la, lb) + 0.05) / (math.min(la, lb) + 0.05);
    }

    void expectAtLeast(String what, Color fg, Color bg, double floor) {
      final measured = contrastRatio(fg, bg);
      expect(measured, greaterThanOrEqualTo(floor),
          reason: "$what measures ${measured.toStringAsFixed(2)}:1 against its "
              "background, under the ${floor}:1 floor");
    }

    BorderSide sideFor(Set<WidgetState> states) {
      final border = theme.inputDecorationTheme.border!;
      return (border as WidgetStateProperty<InputBorder>).resolve(states).borderSide;
    }

    test("a field's resting border is a visible boundary", () {
      expectAtLeast("outline", scheme.outline, scheme.surface, 3.0);
      expectAtLeast("the resting field border",
          sideFor(const <WidgetState>{}).color, scheme.surface, 3.0);
    });

    test("the focused border is visible -- primary alone is not", () {
      final focused = sideFor({WidgetState.focused}).color;

      expectAtLeast("the focused field border", focused, scheme.surface, 3.0);
      // The specific trap: M3 uses colorScheme.primary here, and this app's
      // primary is #27B1FF at 2.38:1.
      expect(focused, isNot(scheme.primary));
    });

    test("muted text clears the 4.5:1 text floor, not just 3:1", () {
      expectAtLeast(
          "onSurfaceVariant", scheme.onSurfaceVariant, scheme.surface, 4.5);
    });

    test("the accent used as a foreground is readable on its tints", () {
      // 3:1, the non-text floor. These are the two places accentBlueDark is
      // also used for small TEXT (the nav's selected label, the editor chip's
      // label), where the floor is 4.5:1 and they measure ~3.6. Closeable
      // on-palette by moving the accent from $primary-palette 700 to 800
      // (#0072b4, which clears every floor); left at 700 deliberately because
      // 800 reads navy. Recorded in mobile/CLAUDE.md.
      expectAtLeast("accentBlueDark on primaryContainer", scheme.onPrimaryContainer,
          scheme.primaryContainer, 3.0);
      expectAtLeast("the nav's selected icon on its pill", accentBlueDark,
          accentContainer, 3.0);
    });

    test("the focused floating label is text, so it clears the text floor", () {
      final style = (theme.inputDecorationTheme.floatingLabelStyle!
              as WidgetStateProperty<TextStyle>)
          .resolve({WidgetState.focused});

      expectAtLeast("the focused floating label", style.color!, scheme.surface, 4.5);
      // Material's default is colorScheme.primary at 2.38:1.
      expect(style.color, isNot(scheme.primary));
    });
  });

  group("the input border resolves every state", () {
    // Setting `border` to a WidgetStateInputBorder makes InputDecorator skip
    // `_getDefaultBorder` entirely, so this theme owns all four states rather
    // than only the ones it meant to change. A dropped branch is invisible:
    // nothing fails to compile and no other test looks.
    BorderSide sideFor(Set<WidgetState> states) {
      final border = theme.inputDecorationTheme.border!;
      return (border as WidgetStateProperty<InputBorder>).resolve(states).borderSide;
    }

    test("enabled is the outline, at the design's 1.5px", () {
      final side = sideFor(const <WidgetState>{});

      expect(side.color, scheme.outline);
      expect(side.width, 1.5);
    });

    test("focused is the darker accent at 2px", () {
      final side = sideFor({WidgetState.focused});

      expect(side.color, accentBlueDark);
      expect(side.width, 2);
    });

    test("error is the error colour, focused or not", () {
      expect(sideFor({WidgetState.error}).color, scheme.error);
      expect(sideFor({WidgetState.error, WidgetState.focused}).color, scheme.error);
      expect(sideFor({WidgetState.error, WidgetState.focused}).width, 2);
    });

    test("error beats focus, so a bad value never looks merely focused", () {
      expect(sideFor({WidgetState.error, WidgetState.focused}).color,
          isNot(accentBlueDark));
    });

    test("disabled is faint and is neither of the above", () {
      final side = sideFor({WidgetState.disabled});

      expect(side.color, isNot(scheme.outline));
      expect(side.color, isNot(scheme.error));
      expect(side.color.a, lessThan(1.0));
    });

    test("every state keeps the 10px radius", () {
      final border = theme.inputDecorationTheme.border!
          as WidgetStateProperty<InputBorder>;

      for (final states in const [
        <WidgetState>{},
        {WidgetState.focused},
        {WidgetState.error},
        {WidgetState.disabled},
      ]) {
        final resolved = border.resolve(states) as OutlineInputBorder;
        expect(resolved.borderRadius, BorderRadius.circular(10));
      }
    });
  });
}
