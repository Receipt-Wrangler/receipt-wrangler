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

/// Disabled fills. Deliberately NOT a border colour -- see [borderSlate].
const slate300 = Color(0xFFCBD5E1);

/// Input, chip and outlined-button borders.
///
/// The design draws these in slate-300 (`#CBD5E1`), but at 1.5-2px; Flutter
/// renders 1px, and slate-300 on white is **1.48:1** -- far under the 3:1 WCAG
/// 2.1 SC 1.4.11 asks for the visual boundary that identifies a component, which
/// is exactly what an outlined field's border is. This is that role pushed to
/// **3.38:1**, with margin over the bar rather than the 3.05:1 that just clears
/// it, because a 1px antialiased hairline reads lighter than its nominal colour.
/// For scale, Material 3's own baseline `outline` is `#79747E`, 4.55:1.
const borderSlate = Color(0xFF7E8DA1);

/// Purely decorative greys. A glyph that carries meaning -- a chevron, a dismiss
/// X -- belongs at `onSurfaceVariant` instead: this is 2.56:1, under SC 1.4.11.
const slate400 = Color(0xFF94A3B8);

/// Muted text: field labels, hints, list subtitles, section micro-labels.
const slate500 = Color(0xFF64748B);

/// Secondary text that still has to be read comfortably.
const slate700 = Color(0xFF334155);

/// Selected-state tint for a chip that picks a *mode* rather than a value.
/// The design's `rgba(39,177,255,.08)`, flattened against white.
const accentTint = Color(0xFFEAF7FF);

/// A filled accent container -- the bottom nav's selected pill.
///
/// Distinct from [accentTint], and both are needed. This one is strong enough
/// that a 24px icon on it reads as selected across the room; the lighter tint is
/// not, which is exactly the bug it was introduced to fix. Conversely
/// [accentBlueDark] clears AA on [accentTint] at chip-label size but only clears
/// the large-text / UI threshold on this one.
const accentContainer = Color(0xFFCCECFF);

/// Foreground for [accentTint] and [accentContainer].
///
/// **Not** a general "accent on white". Accent icons and button labels on white
/// are the theme's `primary` -- that is what `receipt_form.dart` already does
/// and what Material 3 gives every text and outlined button. This darker step
/// exists because `primary` over the accent tints is barely 2:1.
const accentBlueDark = Color(0xFF0086D4);

/// Border for a chip in the tinted "selected mode" treatment (the filter
/// editor's operation chips) -- the design's `#bbe6ff`.
const selectedChipBorder = Color(0xFFBBE6FF);
