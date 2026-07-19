import 'package:built_collection/built_collection.dart';
import 'package:built_value/json_object.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/dashboard_widgets/report_widget.dart';

// The WebView render path is exercised on-device, not in a widget test — the
// webview_flutter platform channel is unavailable under flutter_test (which is
// why report_preview_screen.dart likewise has no widget test). These tests pin
// the two WebView-free pieces of logic the widget depends on: reading the
// pinned template id out of the untyped config blob, and the server-driven
// download gate.
void main() {
  BuiltMap<String, JsonObject?> config(Map<String, JsonObject?> entries) =>
      BuiltMap<String, JsonObject?>(entries);

  group('reportTemplateIdFromConfig', () {
    test('reads an integer reportTemplateId', () {
      expect(
        reportTemplateIdFromConfig(config({'reportTemplateId': JsonObject(7)})),
        7,
      );
    });

    test('accepts a whole-number double', () {
      // JSON numbers arrive as num; a whole id in double form is still valid.
      expect(
        reportTemplateIdFromConfig(
            config({'reportTemplateId': JsonObject(7.0)})),
        7,
      );
    });

    test('returns null for a fractional reportTemplateId', () {
      // Must not truncate 7.5 -> 7 and load the wrong template.
      expect(
        reportTemplateIdFromConfig(
            config({'reportTemplateId': JsonObject(7.5)})),
        isNull,
      );
    });

    test('returns null when the key is absent', () {
      expect(
        reportTemplateIdFromConfig(config({'chartGrouping': JsonObject('X')})),
        isNull,
      );
    });

    test('returns null for a null config', () {
      expect(reportTemplateIdFromConfig(null), isNull);
    });

    test('returns null when the value is not a number', () {
      expect(
        reportTemplateIdFromConfig(
            config({'reportTemplateId': JsonObject('nope')})),
        isNull,
      );
    });
  });

  group('reportWidgetCanDownload', () {
    test('true when allowedActions contains generate', () {
      expect(
        reportWidgetCanDownload(BuiltList<String>(['read', 'generate'])),
        isTrue,
      );
    });

    test('false when allowedActions omits generate', () {
      expect(
        reportWidgetCanDownload(BuiltList<String>(['read'])),
        isFalse,
      );
    });

    test('false for empty allowedActions', () {
      expect(reportWidgetCanDownload(BuiltList<String>()), isFalse);
    });

    test('false for null allowedActions', () {
      expect(reportWidgetCanDownload(null), isFalse);
    });
  });
}
