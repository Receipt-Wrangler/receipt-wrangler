import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/persistence/global_shared_preferences.dart';
import 'package:receipt_wrangler_mobile/service/crash_reporting.dart';
import 'package:sentry_flutter/sentry_flutter.dart';
import 'package:shared_preferences/shared_preferences.dart';

void main() {
  setUpAll(() async {
    TestWidgetsFlutterBinding.ensureInitialized();
    // Must run before the first getInstance(); GlobalSharedPreferences caches a
    // single instance and has no reset, so we mutate the live instance per test.
    SharedPreferences.setMockInitialValues({});
    await GlobalSharedPreferences.initialize();
  });

  setUp(() async {
    await GlobalSharedPreferences.instance.remove(kCrashReportingKey);
  });

  group('isCrashReportingEnabled', () {
    test('defaults to true (opted in) when the preference is unset', () {
      expect(isCrashReportingEnabled(), isTrue);
    });

    test('returns false when the preference is stored false', () async {
      await GlobalSharedPreferences.instance.setBool(kCrashReportingKey, false);
      expect(isCrashReportingEnabled(), isFalse);
    });

    test('returns true when the preference is stored true', () async {
      await GlobalSharedPreferences.instance.setBool(kCrashReportingKey, true);
      expect(isCrashReportingEnabled(), isTrue);
    });
  });

  group('configureSentry', () {
    test('applies the privacy-hardened configuration', () {
      final options = SentryFlutterOptions();

      configureSentry(options);

      expect(options.dsn, isNotEmpty);
      expect(options.sendDefaultPii, isFalse); // no IP / no PII
      expect(options.enableAutoSessionTracking, isFalse); // GlitchTip: no sessions
      expect(options.enablePrintBreadcrumbs, isFalse); // logs not turned into breadcrumbs
      expect(options.attachScreenshot, isFalse); // never capture screen contents
      expect(options.attachViewHierarchy, isFalse);
      expect(options.tracesSampleRate, 0.0); // crash/error only, no perf tracing
    });
  });
}
