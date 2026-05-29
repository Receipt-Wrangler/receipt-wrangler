import 'package:flutter_test/flutter_test.dart';

/// Pumps frames until [finder] matches at least one widget, or [timeout]
/// elapses. Use instead of [WidgetTester.pumpAndSettle] on screens that host
/// infinite animations (bootstrap loaders, shimmers, carousels) — those never
/// settle and `pumpAndSettle` will hang until its own timeout fires.
Future<void> pumpUntilFound(
  WidgetTester tester,
  Finder finder, {
  Duration timeout = const Duration(seconds: 10),
  Duration step = const Duration(milliseconds: 100),
}) async {
  final deadline = DateTime.now().add(timeout);
  while (DateTime.now().isBefore(deadline)) {
    await tester.pump(step);
    // Some finders (e.g. `.last`, `.first`) throw StateError when the
    // parent matches nothing -- they call `iterable.last` directly. Treat
    // any exception during evaluation as "not found yet" so polling
    // continues until the parent populates or the timeout fires.
    try {
      if (finder.evaluate().isNotEmpty) return;
    } catch (_) {
      // Keep polling.
    }
  }
  throw StateError(
    'Timed out after ${timeout.inSeconds}s waiting for $finder',
  );
}

/// Inverse of [pumpUntilFound]: pumps frames until [finder] no longer
/// matches anything, or [timeout] elapses. Useful for asserting that a
/// widget disappears after a user action (e.g. swipe-delete on a list
/// row) where the disappearance is driven by an API round-trip rather
/// than an immediate setState.
Future<void> pumpUntilGone(
  WidgetTester tester,
  Finder finder, {
  Duration timeout = const Duration(seconds: 10),
  Duration step = const Duration(milliseconds: 100),
}) async {
  final deadline = DateTime.now().add(timeout);
  while (DateTime.now().isBefore(deadline)) {
    await tester.pump(step);
    try {
      if (finder.evaluate().isEmpty) return;
    } catch (_) {
      return;
    }
  }
  throw StateError(
    'Timed out after ${timeout.inSeconds}s waiting for $finder to disappear',
  );
}
