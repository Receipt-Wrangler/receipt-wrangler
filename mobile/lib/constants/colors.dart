import 'dart:ui';

// ---------------------------------------------------------------------------
// Where these values come from
//
// The canonical Receipt Wrangler palette is the set of Angular Material M2 maps
// in `desktop/src/variables.scss`, which is what the web app's theme is actually
// built from (`mat.m2-define-palette` in `desktop/src/styles.scss`). Almost every
// constant below is a stop out of one of those maps, and each says which:
//
//   $primary-palette   50 #ccecff · 100 #bbe6ff · 200 #a4deff · 300 #85d2ff
//                      400 #5dc4ff · 500 #27b1ff · 600 #009efa · 700 #0086d4
//                      800 #0072b4 · 900 #006199
//   $accent-palette    Tailwind slate, 50 #f8fafc through 900 #0f172a
//   $warn-palette      500 #d63333
//
// **Source a new value from those maps before inventing one.** Exactly two
// constants here are not stops -- [accentTint] and [borderSlate] -- and both say
// why in their own doc comment. Nothing mechanically ties this file to the SCSS,
// so the two clients stay in step by hand.
// ---------------------------------------------------------------------------

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

/// Page canvas behind raised white content. `$accent-palette` 50.
const slate50 = Color(0xFFF8FAFC);

/// A filled, low-emphasis surface -- the operation pill on a filter condition.
/// `$accent-palette` 100.
const slate100 = Color(0xFFF1F5F9);

/// Dividers, and the hairline above a bottom action bar. `$accent-palette` 200.
const slate200 = Color(0xFFE2E8F0);

/// Disabled fills. `$accent-palette` 300, and deliberately NOT a border colour
/// -- see [borderSlate].
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
///
/// **One of the two values here that is not a palette stop**, and knowingly so:
/// `$accent-palette` jumps 400 `#94A3B8` (2.56:1, under the floor) straight to
/// 500 `#64748B` (4.76:1, as dark as the label text and visibly heavy on every
/// field in the app). There is nothing between them, so this is a half-step.
const borderSlate = Color(0xFF7E8DA1);

/// Purely decorative greys; `$accent-palette` 400. A glyph that carries meaning
/// -- a chevron, a dismiss X -- belongs at `onSurfaceVariant` instead: this is
/// 2.56:1, under SC 1.4.11.
const slate400 = Color(0xFF94A3B8);

/// Muted text: field labels, hints, list subtitles, section micro-labels.
/// `$accent-palette` 500, which desktop also aliases as `$basic-gray`.
const slate500 = Color(0xFF64748B);

/// Secondary text that still has to be read comfortably. `$accent-palette` 700.
const slate700 = Color(0xFF334155);

/// Selected-state tint for a chip that picks a *mode* rather than a value.
///
/// **The other non-stop value**: the design's `rgba(39,177,255,.08)` flattened
/// against white. It is lighter than `$primary-palette` 50 ([accentContainer])
/// on purpose -- [accentBlueDark] reads at 3.59:1 on this and 3.17:1 on that, and
/// this one sits under 13px chip text where that one sits behind a 24px icon.
const accentTint = Color(0xFFEAF7FF);

/// A filled accent container -- the bottom nav's selected pill.
/// `$primary-palette` 50; desktop pairs the same stop with [accentBlueDark] in
/// `roles/role-presets.ts` (`PRIMARY_TINT` / `PRIMARY_COLOR`).
///
/// Distinct from [accentTint], and both are needed. This one is strong enough
/// that a 24px icon on it reads as selected across the room; the lighter tint is
/// not, which is exactly the bug it was introduced to fix. Conversely
/// [accentBlueDark] clears AA on [accentTint] at chip-label size but only clears
/// the large-text / UI threshold on this one.
const accentContainer = Color(0xFFCCECFF);

/// Foreground for [accentTint] and [accentContainer].
///
/// `$primary-palette` **700** -- a brand colour, not an off-brand darkening.
/// `mat.m2-define-palette` makes 700 the web theme's "darker" variant, and
/// desktop uses it for exactly this job in about sixteen places (`filter-bar`,
/// `month-stepper`, `role-presets.ts`, the reports panels).
///
/// **Not** a general "accent on white": accent icons and button labels on white
/// are the theme's `primary` (stop 500) -- what `receipt_form.dart` already does
/// and what Material 3 gives every text and outlined button. This darker step
/// exists because 500 over the accent tints is barely 2:1.
const accentBlueDark = Color(0xFF0086D4);

/// Border for a chip in the tinted "selected mode" treatment (the filter
/// editor's operation chips). `$primary-palette` 100.
const selectedChipBorder = Color(0xFFBBE6FF);
