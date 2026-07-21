import 'package:receipt_wrangler_mobile/persistence/global_shared_preferences.dart';
import 'package:sentry_flutter/sentry_flutter.dart';

/// Privacy-respecting crash & error reporting via GlitchTip (Sentry-compatible).
///
/// - Reports crashes / uncaught exceptions / native hangs only — no behavioral
///   analytics, no ads, no advertising IDs.
/// - No PII: `sendDefaultPii = false` (no IP), never calls `setUser`, no
///   screenshots/view-hierarchy, and app logs are not turned into breadcrumbs.
/// - Opt-out: on by default, but the user can disable it live from the profile
///   screen ([setCrashReportingEnabled]) — `Sentry.close()` stops the whole SDK
///   (Dart + native) immediately; re-enabling re-inits it, no app restart.
const String kCrashReportingKey = 'crashReportingEnabled';

const String _dsn =
    'https://167ea649143c46418671b1597dfd8a9b@app.glitchtip.com/25926';

/// Whether the user has crash reporting enabled (default: true / opted-in).
bool isCrashReportingEnabled() =>
    GlobalSharedPreferences.instance.getBool(kCrashReportingKey) ?? true;

/// Shared, privacy-hardened Sentry configuration.
void configureSentry(SentryFlutterOptions options) {
  options.dsn = _dsn;
  options.sendDefaultPii = false; // no IP address / no PII
  options.enableAutoSessionTracking = false; // GlitchTip has no sessions
  options.enablePrintBreadcrumbs = false; // don't turn app logs into breadcrumbs
  options.attachScreenshot = false; // never capture screen contents
  options.attachViewHierarchy = false;
  options.tracesSampleRate = 0.0; // crash/error reporting only, no perf tracing
}

/// Toggle crash reporting at runtime (opt-out; default on). Effective
/// immediately — no restart.
Future<void> setCrashReportingEnabled(bool enabled) async {
  // Apply the SDK change first, then persist — a failed init/close must not
  // leave the stored preference out of sync with what the SDK is actually doing.
  if (enabled) {
    await SentryFlutter.init(configureSentry);
  } else {
    await Sentry.close();
  }
  await GlobalSharedPreferences.instance.setBool(kCrashReportingKey, enabled);
}
