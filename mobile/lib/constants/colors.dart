import 'dart:ui';

const successGreen = Color.fromRGBO(144, 238, 144, 1);
const errorRed = const Color.fromRGBO(242, 191, 191, 1);

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
