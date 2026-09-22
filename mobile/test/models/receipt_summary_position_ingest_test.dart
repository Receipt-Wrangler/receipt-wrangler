import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_summary.dart';

/// Guards the API boundary for `ReceiptSummaryPosition`, the newest member of the family
/// that caused the two production login outages.
///
/// built_value renders an OpenAPI enum as a CLOSED set: without the `fallback: true`
/// patch on `BOTTOM`, `_$valueOf` ends in `default: throw ArgumentError(name)` and a
/// single unrecognized wire value fails the *entire* enclosing payload. This one rides on
/// `GroupReceiptSettings`, which rides on `AppData` -- so the enclosing payload is
/// **login**. Adding a third position server-side would brick every already-released
/// build at the sign-in screen.
///
/// The fallback lands on `BOTTOM` rather than a sentinel member because that is a real
/// placement -- and the same one the server's `ReceiptSummaryPosition.OrDefault()` would
/// have sent -- so an unknown value degrades to where the block always rendered instead of
/// to a value the rest of the app has to special-case.
///
/// `api/generate-client.sh` re-applies the patch on every regen and fails loudly if it
/// cannot, but the annotation compiles perfectly well when absent, so this suite is what
/// proves the behaviour is actually there.
Map<String, Object?> _settingsJson(String position) => {
      'id': 1,
      'createdAt': '2026-01-01T00:00:00Z',
      'groupId': 1,
      'receiptSummaryEnabled': true,
      'receiptSummaryPosition': position,
    };

void main() {
  group('unknown wire position', () {
    test('deserializes to BOTTOM instead of throwing', () {
      expect(
        standardSerializers.deserializeWith(
            ReceiptSummaryPosition.serializer, 'FLOATING'),
        ReceiptSummaryPosition.BOTTOM,
      );
    });

    // The server normalizes an empty position away at every emit point, so this is belt
    // and braces -- but an empty string is the one unknown value a partially migrated
    // install could realistically produce.
    test('an empty wire value degrades the same way', () {
      expect(
        standardSerializers.deserializeWith(ReceiptSummaryPosition.serializer, ''),
        ReceiptSummaryPosition.BOTTOM,
      );
    });

    // The blast radius that actually matters: one bad value must not take the enclosing
    // settings -- and therefore the whole AppData payload, and therefore login -- down.
    test('does not fail the enclosing GroupReceiptSettings payload', () {
      final settings = standardSerializers.deserializeWith(
          GroupReceiptSettings.serializer, _settingsJson('FLOATING'));

      expect(settings, isNotNull);
      expect(settings!.groupId, 1);
      expect(settings.receiptSummaryEnabled, isTrue);
    });

    // And the degradation is the historical placement, not a blank screen.
    test('renders at the bottom, which is where the block always was', () {
      final settings = standardSerializers.deserializeWith(
          GroupReceiptSettings.serializer, _settingsJson('FLOATING'))!;

      expect(isReceiptSummaryAtTop(settings.receiptSummaryPosition), isFalse);
    });
  });

  group('known wire positions', () {
    test('round-trip', () {
      for (final position in [
        ReceiptSummaryPosition.TOP,
        ReceiptSummaryPosition.BOTTOM,
      ]) {
        final wire = standardSerializers.serializeWith(
            ReceiptSummaryPosition.serializer, position);
        expect(
          standardSerializers.deserializeWith(
              ReceiptSummaryPosition.serializer, wire),
          position,
        );
      }
    });
  });
}
