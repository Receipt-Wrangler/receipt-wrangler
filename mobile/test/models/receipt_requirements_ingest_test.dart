import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';

// `AppData.groupReceiptRequirements` carries the role-required receipt fields,
// already resolved server-side (waivers applied) and only for groups where
// something is required. It is deliberately NOT in AppData's swagger `required`
// list, so a payload from an older server — which omits it — must still parse
// and must read as "nothing required anywhere". These go through the real
// generated deserializer, where a login-breaking regression would surface.
Map<String, Object?> _appDataJson({Object? groupReceiptRequirements}) => {
  'about': {'buildDate': '2026-01-01', 'version': '0.0.0'},
  'claims': {
    'userId': 1,
    'displayName': 'tester',
    'defaultAvatarColor': '#ffffff',
    'username': 'tester',
    'iss': 'receiptWrangler',
    'exp': 4102444800,
  },
  'groups': <Object?>[],
  'users': <Object?>[],
  'userPreferences': {'id': 1, 'createdAt': '2026-01-01', 'userId': 1},
  'featureConfig': {'aiPoweredReceipts': false, 'enableLocalSignUp': false},
  'categories': <Object?>[],
  'tags': <Object?>[],
  'currencyDisplay': r'$',
  'icons': <Object?>[],
  'appPermissions': <String>[],
  'groupPermissions': <String, Object?>{},
  if (groupReceiptRequirements != null)
    'groupReceiptRequirements': groupReceiptRequirements,
};

PermissionsModel _ingest(Map<String, Object?> json) {
  final appData =
      standardSerializers.deserializeWith(AppData.serializer, json)!;
  // The same call storeAppData makes.
  return PermissionsModel()
    ..setReceiptRequirements(appData.groupReceiptRequirements?.toMap());
}

void main() {
  group('AppData.groupReceiptRequirements ingest', () {
    test('stores each group\'s requirements keyed by int group id', () {
      final model = _ingest(
        _appDataJson(
          groupReceiptRequirements: {
            '7': {'commentRequired': true, 'imageRequired': false},
            '9': {'commentRequired': false, 'imageRequired': true},
          },
        ),
      );

      expect(model.receiptRequirements(7).commentRequired, isTrue);
      expect(model.receiptRequirements(7).imageRequired, isFalse);
      expect(model.receiptRequirements(9).commentRequired, isFalse);
      expect(model.receiptRequirements(9).imageRequired, isTrue);
    });

    test('an absent group requires nothing', () {
      final model = _ingest(
        _appDataJson(
          groupReceiptRequirements: {
            '7': {'commentRequired': true, 'imageRequired': true},
          },
        ),
      );

      expect(model.receiptRequirements(8).commentRequired, isFalse);
      expect(model.receiptRequirements(8).imageRequired, isFalse);
    });

    test(
      'an older server that omits the field still parses, requiring nothing',
      () {
        final model = _ingest(_appDataJson());

        expect(model.receiptRequirements(7).commentRequired, isFalse);
        expect(model.receiptRequirements(7).imageRequired, isFalse);
      },
    );

    test('an empty map requires nothing', () {
      final model = _ingest(_appDataJson(groupReceiptRequirements: {}));

      expect(model.receiptRequirements(7).commentRequired, isFalse);
      expect(model.receiptRequirements(7).imageRequired, isFalse);
    });

    test('a later payload replaces the earlier one', () {
      // AppData is re-fetched during a session; a role that stopped requiring
      // something must stop being enforced client-side too.
      final model = _ingest(
        _appDataJson(
          groupReceiptRequirements: {
            '7': {'commentRequired': true, 'imageRequired': true},
          },
        ),
      );
      model.setReceiptRequirements(null);

      expect(model.receiptRequirements(7).commentRequired, isFalse);
      expect(model.receiptRequirements(7).imageRequired, isFalse);
    });

    test('an unparseable group key is skipped', () {
      final model =
          PermissionsModel()..setReceiptRequirements({
            'not-a-group': ReceiptRequirements(
              (b) =>
                  b
                    ..commentRequired = true
                    ..imageRequired = true,
            ),
          });

      expect(model.receiptRequirements(0).commentRequired, isFalse);
    });
  });
}
