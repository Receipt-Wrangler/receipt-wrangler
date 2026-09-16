import 'dart:ui';

const successGreen = Color.fromRGBO(144, 238, 144, 1);
const errorRed = const Color.fromRGBO(242, 191, 191, 1);

/// Receipt status tint for NEEDS_ATTENTION. It gave [errorRed] up to DECLINED, so
/// "needs attention" now reads as a warning rather than a rejection. Matches the
/// desktop chip's `$warning-amber`; pale like every other status tint, because
/// `ListItemTrailingStatus` paints its label in the default dark `onBackground`.
const warningAmber = Color.fromRGBO(255, 224, 178, 1);

/// Receipt status tint for DRAFT, and the fallback for a status this build does
/// not recognize.
const neutralStatusGrey = Color.fromRGBO(224, 224, 224, 1);

/// Background for the in-page notices (the "Quick Scan unavailable" banner on
/// the receipt form, and the queued confirmation in the Quick Scan sheet).
///
/// A darkened form of the theme's `secondary` slate, chosen so [onNoticeSurface]
/// white text sits at ~8.9:1 contrast. The theme's `ColorScheme` never defines
/// `secondaryContainer`, so it falls back to `secondary` (#8EA1AC) with black
/// `onSecondary` text — legible, but muddy at the small type these notices use.
const noticeSurface = Color.fromRGBO(62, 76, 85, 1);

/// Foreground for [noticeSurface].
const onNoticeSurface = Color.fromRGBO(255, 255, 255, 1);

// ---------------------------------------------------------------------------
// Neutral scale
//
// The app's `ColorScheme` (lib/theme/app_theme.dart) spells out every neutral
// role, because Material 3 falls back to `onBackground` / `surface` for the ones
// it is not given -- which in this scheme means pure black borders, pure black
// "muted" text and no surface tone at all. These are the values it is given.
//
// The scale is Tailwind's slate, which is what the receipt-filter design
// (`Quick Date Filtering.dc.html`) is drawn in.
// ---------------------------------------------------------------------------

/// Page canvas behind raised white content.
const slate50 = Color(0xFFF8FAFC);

/// A filled, low-emphasis surface -- the operation pill on a filter condition.
const slate100 = Color(0xFFF1F5F9);

/// Dividers, and the hairline above a bottom action bar.
const slate200 = Color(0xFFE2E8F0);

/// Input, chip and outlined-button borders.
const slate300 = Color(0xFFCBD5E1);

/// De-emphasized glyphs: a chevron, a dismiss X, a dashed placeholder border.
const slate400 = Color(0xFF94A3B8);

/// Muted text: field labels, hints, list subtitles, section micro-labels.
const slate500 = Color(0xFF64748B);

/// Secondary text that still has to be read comfortably.
const slate700 = Color(0xFF334155);

/// Selected-state tint for a chip that picks a *mode* rather than a value.
/// The design's `rgba(39,177,255,.08)`, flattened against white.
const accentTint = Color(0xFFEAF7FF);

/// The accent for text and icons **on white**.
///
/// The theme's `primary` (`#27B1FF`) is 2.2:1 against white and fails WCAG for
/// anything but a large solid fill, so it must never be used as a foreground
/// color. This darker step is the design's own `#0086D4` and clears AA.
const accentBlueDark = Color(0xFF0086D4);

/// Border for a chip in the tinted "selected mode" treatment (the filter
/// editor's operation chips) -- the design's `#bbe6ff`.
const selectedChipBorder = Color(0xFFBBE6FF);
