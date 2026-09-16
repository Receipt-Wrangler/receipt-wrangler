import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/constants/colors.dart';

/// The app's single [ThemeData].
///
/// Extracted from `main.dart` so the color roles below can be asserted on
/// directly (`test/theme/app_theme_test.dart`) and so a harness can render a
/// screen exactly as the app does.
///
/// **Every neutral role is spelled out on purpose.** `ColorScheme`'s getters
/// fall back when a role is left null -- `outline` and `outlineVariant` to
/// `onBackground`, `onSurfaceVariant` to `onSurface`, the whole
/// `surfaceContainer*` family to `surface`, `secondaryContainer` to `secondary`.
/// With this scheme's black `onSurface` and white `surface` that gave every
/// Material 3 component a hard black border, black "muted" text with no
/// hierarchy, no surface tone to layer against, and a selected chip in the
/// `#8EA1AC` slate that reads as disabled. Do not drop these.
const String appFontFamily = "Raleway";

ThemeData buildAppTheme() {
  const colorScheme = ColorScheme(
    brightness: Brightness.light,
    primary: Color(0xFF27B1FF),
    onPrimary: Color(0xFFFFFFFF),
    primaryContainer: accentTint,
    onPrimaryContainer: accentBlueDark,
    secondary: Color(0xFF8EA1AC),
    onSecondary: Color(0xFF000000),
    secondaryContainer: slate100,
    onSecondaryContainer: slate700,
    error: Color(0xFFd63333),
    onError: Color(0xFFFFFFFF),
    surface: Color(0xFFFFFFFF),
    onSurface: Color(0xFF000000),
    onSurfaceVariant: slate500,
    outline: slate300,
    outlineVariant: slate200,
    surfaceDim: slate50,
    surfaceBright: Color(0xFFFFFFFF),
    surfaceContainerLowest: Color(0xFFFFFFFF),
    // Cards stay white so they read as raised against a slate canvas.
    surfaceContainerLow: Color(0xFFFFFFFF),
    surfaceContainer: slate100,
    surfaceContainerHigh: slate100,
    surfaceContainerHighest: slate200,
    background: Color(0xFFFFFFFF),
    onBackground: Color(0xFF000000),
  );

  return ThemeData(
    fontFamily: appFontFamily,
    inputDecorationTheme: InputDecorationTheme(
      border: OutlineInputBorder(
        borderRadius: BorderRadius.circular(10),
      ),
    ),
    chipTheme: ChipThemeData(
      shape: RoundedRectangleBorder(
        borderRadius: BorderRadius.circular(50),
      ),
      // A selected chip is a filled accent chip everywhere in the app.
      //
      // The label colour has to ride on `labelStyle` as a WidgetStateColor:
      // `RawChip` resolves `labelStyle.color` through
      // `WidgetStateProperty.resolveAs` and, under Material 3, never consults
      // `secondaryLabelStyle` at all. Left to the defaults the selected label
      // is `onSecondaryContainer` -- dark slate on `#27B1FF`, which is the
      // contrast bug the category and status pickers shipped with.
      // `labelStyle` REPLACES the defaults rather than merging, so the family
      // and size have to be restated here.
      selectedColor: colorScheme.primary,
      checkmarkColor: colorScheme.onPrimary,
      labelStyle: TextStyle(
        fontFamily: appFontFamily,
        fontSize: 14,
        fontWeight: FontWeight.w500,
        color: WidgetStateColor.resolveWith((states) =>
            states.contains(WidgetState.selected)
                ? colorScheme.onPrimary
                : slate700),
      ),
    ),
    // Both halves of a toolbar's icons, together. M3 takes `leading` from
    // `onSurface` (black here) and `actions` from `onSurfaceVariant` (slate),
    // so an app bar with both renders two different greys.
    appBarTheme: const AppBarTheme(
      iconTheme: IconThemeData(color: slate700),
      actionsIconTheme: IconThemeData(color: slate700),
    ),
    bottomSheetTheme: const BottomSheetThemeData(
      backgroundColor: Colors.white,
      modalBackgroundColor: Colors.white,
      surfaceTintColor: Colors.white,
    ),
    colorScheme: colorScheme,
    useMaterial3: true,
  );
}
