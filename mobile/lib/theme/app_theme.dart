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
    outline: borderSlate,
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
      // A WidgetStateInputBorder, not a plain OutlineInputBorder, and that is
      // load-bearing rather than fussy. For a non-filled field
      // `InputDecorator._getDefaultBorder` THROWS AWAY whatever `borderSide` the
      // theme's border carries and substitutes `_InputDecoratorDefaultsM3`'s --
      // which is built fresh from the context and never merged with this theme,
      // so `InputDecorationTheme.outlineBorder` is ignored too. Only the shape
      // survives. Its one escape hatch is an early return when the border IS a
      // `WidgetStateProperty<InputBorder>`, which bypasses the whole resolution.
      //
      // The cost of taking that hatch: this theme now owns EVERY state, so all
      // four are written out. Dropping the error branch would leave every failed
      // validator in the app without its red border, and nothing would fail to
      // compile -- hence the state test in `app_theme_test.dart`.
      // The focused floating label keeps `onSurfaceVariant` (4.76:1) rather than
      // turning accent. Material's default here is `primary` -- 2.38:1, unusable
      // -- and even `accentBlueDark` is only 3.92:1, which clears the 3:1 floor
      // for a non-text boundary but not the 4.5:1 one this *is* text. Focus is
      // carried by the border instead, which goes 1.5px slate to 2px accent:
      // unmistakable, and it costs no third blue.
      floatingLabelStyle: WidgetStateTextStyle.resolveWith((states) {
        if (states.contains(WidgetState.focused) &&
            !states.contains(WidgetState.error)) {
          return TextStyle(color: colorScheme.onSurfaceVariant);
        }
        // Empty merges as a no-op, so error red and the disabled 38% stay
        // Material's own.
        return const TextStyle();
      }),
      border: WidgetStateInputBorder.resolveWith((states) {
        final radius = BorderRadius.circular(10);

        if (states.contains(WidgetState.disabled)) {
          // Matches the M3 default. Disabled controls are exempt from the
          // contrast rule below, so this stays faint on purpose.
          return OutlineInputBorder(
            borderRadius: radius,
            borderSide:
                BorderSide(color: colorScheme.onSurface.withValues(alpha: 0.12)),
          );
        }
        if (states.contains(WidgetState.error)) {
          return OutlineInputBorder(
            borderRadius: radius,
            borderSide: BorderSide(
              color: colorScheme.error,
              width: states.contains(WidgetState.focused) ? 2 : 1.5,
            ),
          );
        }
        if (states.contains(WidgetState.focused)) {
          // `primary` is 2.38:1 on white -- too weak for the one state that has
          // to be unmistakable. accentBlueDark is 3.92:1.
          return OutlineInputBorder(
            borderRadius: radius,
            borderSide: const BorderSide(color: accentBlueDark, width: 2),
          );
        }
        return OutlineInputBorder(
          borderRadius: radius,
          borderSide: BorderSide(color: colorScheme.outline, width: 1.5),
        );
      }),
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
    // NavigationBar takes its indicator from `secondaryContainer`, which used
    // to fall through to the `#8EA1AC` secondary -- a solid slate-blue pill.
    // Naming that role slate-100 turned the selected destination into a
    // near-white pill on a near-white bar, i.e. invisible, so the bar states its
    // own selected treatment rather than riding a role it does not want. Values
    // are the design's own nav.
    navigationBarTheme: NavigationBarThemeData(
      indicatorColor: accentContainer,
      iconTheme: WidgetStateProperty.resolveWith(
        (states) => IconThemeData(
          color: states.contains(WidgetState.selected)
              ? accentBlueDark
              : slate700,
        ),
      ),
      labelTextStyle: WidgetStateProperty.resolveWith(
        (states) => TextStyle(
          fontFamily: appFontFamily,
          fontSize: 12,
          fontWeight: FontWeight.w500,
          color: states.contains(WidgetState.selected)
              ? accentBlueDark
              : slate700,
        ),
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
