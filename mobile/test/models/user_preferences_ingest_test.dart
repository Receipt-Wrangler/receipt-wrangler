import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart';

/// Guards the API boundary for `UserPreferences.closeChipSelectOnSelect`.
///
/// The generated client is regenerated wholesale from `swagger.yml`, and the
/// stub the generator drops in `mobile/api/test/` is never a guard: the
/// generator does not overwrite an existing test file, so that stub still
/// describes the fields the model had the day it was first written and gains
/// nothing when a field is added. It also lives under `mobile/api/`, which is
/// generated output nobody may hand-edit. So the real coverage lives here,
/// beside the two other API-boundary guards
/// (`app_data_permission_ingest_test.dart`, `receipt_status_ingest_test.dart`).
///
/// `closeChipSelectOnSelect` is a desktop behavior flag; mobile carries the
/// field and never reads it. That is exactly why it needs a test: nothing in
/// the app would notice if a regen dropped it, if `build_runner` failed to
/// rebuild the serializer, or if the field started arriving as null.
Map<String, Object?> _userPreferencesJson({Object? closeChipSelectOnSelect}) => {
      // UserPreferences implements BaseModel, so id/createdAt are required.
      'id': 1,
      'createdAt': '2026-01-01',
      'userId': 1,
      if (closeChipSelectOnSelect != null)
        'closeChipSelectOnSelect': closeChipSelectOnSelect,
    };

UserPreferences _deserialize(Map<String, Object?> json) =>
    standardSerializers.deserializeWith(UserPreferences.serializer, json)!;

void main() {
  group('closeChipSelectOnSelect ingest', () {
    // The headline. A server that predates the field - or any payload that
    // simply omits the key - must yield the documented default rather than
    // null, because `false` IS the shipped behavior (the option list stays
    // open). A null here would be a third state nothing handles.
    test('an absent key deserializes to false, not null', () {
      final preferences = _deserialize(_userPreferencesJson());

      expect(preferences.closeChipSelectOnSelect, isNotNull);
      expect(preferences.closeChipSelectOnSelect, isFalse);
    });

    test('true on the wire deserializes to true', () {
      expect(
        _deserialize(_userPreferencesJson(closeChipSelectOnSelect: true))
            .closeChipSelectOnSelect,
        isTrue,
      );
    });

    // Not redundant with the absent-key case above: this is the one that fails
    // if the serializer stops reading the key at all, since a dropped `false`
    // and a defaulted `false` are indistinguishable without it.
    test('false on the wire deserializes to false', () {
      expect(
        _deserialize(_userPreferencesJson(closeChipSelectOnSelect: false))
            .closeChipSelectOnSelect,
        isFalse,
      );
    });

    // The write direction. The app never PUTs preferences today, but the field
    // has to survive a round trip or a future writer would silently reset it.
    test('serializes both values back onto the wire', () {
      for (final value in [true, false]) {
        final preferences =
            _deserialize(_userPreferencesJson(closeChipSelectOnSelect: value));
        final json = standardSerializers.serializeWith(
          UserPreferences.serializer,
          preferences,
        ) as Map<Object?, Object?>;

        expect(json['closeChipSelectOnSelect'], value);
      }
    });

    // The blast radius that actually matters. UserPreferences rides on AppData,
    // which hydrates at login, so a field this payload cannot parse does not
    // degrade a screen - it fails login outright. That is the mechanism behind
    // the two outages `app_data_permission_ingest_test.dart` documents.
    test('does not fail the enclosing AppData payload', () {
      final appData = standardSerializers.deserializeWith(AppData.serializer, {
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
        'userPreferences':
            _userPreferencesJson(closeChipSelectOnSelect: true),
        'featureConfig': {
          'aiPoweredReceipts': false,
          'enableLocalSignUp': false,
        },
        'categories': <Object?>[],
        'tags': <Object?>[],
        'currencyDisplay': r'$',
        'icons': <Object?>[],
        'appPermissions': <String>[],
        'groupPermissions': <String, List<String>>{},
      });

      expect(appData, isNotNull);
      expect(appData!.userPreferences.closeChipSelectOnSelect, isTrue);
    });
  });
}
