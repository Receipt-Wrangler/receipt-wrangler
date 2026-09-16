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
      expect(scheme.outline, slate300);
      expect(scheme.outlineVariant, slate200);
      expect(scheme.outline, isNot(scheme.onSurface));
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
      expect(scheme.primaryContainer, accentTint);
      expect(scheme.onPrimaryContainer, accentBlueDark);
      expect(scheme.primaryContainer, isNot(scheme.primary));
    });
  });

  group("the roles the app already had are unchanged", () {
    test("primary, error and the surface pair", () {
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

  test("an app bar's leading and action icons are the same grey", () {
    // M3 takes leading from onSurface and actions from onSurfaceVariant, so a
    // bar carrying both renders two different colours unless both are pinned.
    expect(theme.appBarTheme.iconTheme?.color, slate700);
    expect(theme.appBarTheme.actionsIconTheme?.color, slate700);
  });
}
