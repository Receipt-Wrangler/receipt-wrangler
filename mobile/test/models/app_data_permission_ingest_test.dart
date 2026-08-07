import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';

/// Regression guard for two production login outages.
///
/// `AppData.appPermissions` / `.groupPermissions` used to be typed as the
/// generated `Permission` enum. built_value renders that as a CLOSED enum whose
/// `_$valueOf` ends in `default: throw ArgumentError(name)`, so a single
/// unrecognized wire value failed the *entire* `AppData` deserialization. Since
/// permissions hydrate on login, every already-released mobile build hard-failed
/// login the moment the backend added a permission -- the request returned HTTP
/// 200 and the user saw a generic red "An error occurred" snackbar.
///
/// It happened twice: 2026-07-24 (`group.members.create`, GHSA-89mm-9qfv-cjg3)
/// and 2026-08-06 (`group.members.grants.update`, PR #661). The fix was to type
/// the effective-permission payload as plain strings in `swagger.yml` -- the
/// enum remains the contract for the *catalog* (`Role.permissions`,
/// `UpsertRoleCommand.permissions`, `PermissionDescriptor.key`), but "which
/// permissions does this user hold" is server-authoritative data.
///
/// This test exercises the real generated deserializer, which is where the
/// throw happened -- upstream of `PermissionsModel`. It fails against a client
/// generated from the pre-fix swagger.
Map<String, Object?> _appDataJson({
  required List<String> appPermissions,
  Map<String, List<String>> groupPermissions = const {},
}) =>
    {
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
      // UserPreferences implements BaseModel, so id/createdAt are required too.
      'userPreferences': {'id': 1, 'createdAt': '2026-01-01', 'userId': 1},
      'featureConfig': {
        'aiPoweredReceipts': false,
        'enableLocalSignUp': false,
      },
      'categories': <Object?>[],
      'tags': <Object?>[],
      'currencyDisplay': r'$',
      'icons': <Object?>[],
      'appPermissions': appPermissions,
      'groupPermissions': groupPermissions,
    };

void main() {
  group('AppData permission ingest', () {
    // Pins the mechanism the tests below exist to route around, and keeps them
    // from silently going vacuous: `Permission` is a CLOSED enum, so anything
    // typed as it rejects a value a newer server might send. That is correct
    // for the catalog and fatal for the effective-permission payload.
    test('the Permission enum itself still rejects an unknown value', () {
      expect(
        () => standardSerializers.deserializeWith(
            Permission.serializer, 'group.not.a.real.permission'),
        throwsA(
          // built_value wraps the underlying ArgumentError. This is verbatim
          // the logcat line both outages produced.
          predicate((e) => e.toString().contains(
              "Deserializing to 'Permission' failed due to: "
              'Invalid argument(s): group.not.a.real.permission')),
        ),
      );
    });

    test('deserializes a payload carrying an unknown permission key', () {
      final json = _appDataJson(
        appPermissions: ['app.users.read', 'app.not.yet.invented'],
        groupPermissions: {
          '7': [
            'group.receipts.read',
            // Shipped in #661 and broke every released build.
            'group.members.grants.update',
            'group.also.not.yet.invented',
          ],
        },
      );

      final appData =
          standardSerializers.deserializeWith(AppData.serializer, json)!;

      expect(appData.appPermissions, contains('app.not.yet.invented'));
      expect(
        appData.groupPermissions['7'],
        contains('group.members.grants.update'),
      );
    });

    test('flows through to PermissionsModel with known keys intact', () {
      final json = _appDataJson(
        appPermissions: ['app.users.read', 'app.not.yet.invented'],
        groupPermissions: {
          '7': ['group.receipts.read', 'group.also.not.yet.invented'],
        },
      );

      final appData =
          standardSerializers.deserializeWith(AppData.serializer, json)!;
      final model = PermissionsModel()
        ..setPermissions(
          appData.appPermissions,
          appData.groupPermissions.toMap(),
        );

      expect(model.hasAppPermission(Permission.appPeriodUsersPeriodRead), true);
      expect(
        model.hasGroupPermission(7, Permission.groupPeriodReceiptsPeriodRead),
        true,
      );
      expect(
        model.hasGroupPermission(7, Permission.groupPeriodReceiptsPeriodDelete),
        false,
      );
    });
  });
}
