import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';

// Permissions are ingested as the raw wire strings AppData delivers, and
// queried with the generated Permission enum. Seeding with literal strings
// below is deliberate: it is what the server actually sends, and it keeps the
// ingest side honest about being an open set rather than a closed enum.
//
// Permission enum constants used here map to these wire strings:
//   appPeriodUsersPeriodRead       -> app.users.read
//   appPeriodUsersPeriodCreate     -> app.users.create
//   appPeriodGroupsPeriodRead      -> app.groups.read
//   groupPeriodReceiptsPeriodRead  -> group.receipts.read
//   groupPeriodReceiptsPeriodCreate-> group.receipts.create
//   groupPeriodReceiptsPeriodUpdate-> group.receipts.update
//   groupPeriodReceiptsPeriodDelete-> group.receipts.delete
void main() {
  group('PermissionsModel', () {
    late PermissionsModel model;

    setUp(() {
      model = PermissionsModel();
    });

    test('denies everything before hydration', () {
      expect(model.hasAppPermission(Permission.appPeriodUsersPeriodRead), false);
      expect(
        model.hasGroupPermission(1, Permission.groupPeriodReceiptsPeriodUpdate),
        false,
      );
      expect(
        model.hasGroupPermissionInAnyGroup(
            Permission.groupPeriodReceiptsPeriodCreate),
        false,
      );
    });

    test('hydrates app and group permissions from AppData', () {
      model.setPermissions(
        ['app.users.read'],
        {
          '7': ['group.receipts.update', 'group.receipts.create'],
          '9': ['group.receipts.read'],
        },
      );

      // App scope.
      expect(model.hasAppPermission(Permission.appPeriodUsersPeriodRead), true);
      expect(
          model.hasAppPermission(Permission.appPeriodUsersPeriodCreate), false);
      expect(
        model.hasAnyAppPermission([
          Permission.appPeriodUsersPeriodCreate,
          Permission.appPeriodUsersPeriodRead,
        ]),
        true,
      );

      // Group scope (keyed by int parsed from the string map keys).
      expect(
        model.hasGroupPermission(7, Permission.groupPeriodReceiptsPeriodUpdate),
        true,
      );
      expect(
        model.hasGroupPermission(7, Permission.groupPeriodReceiptsPeriodDelete),
        false,
      );
      expect(
        model.hasGroupPermission(9, Permission.groupPeriodReceiptsPeriodUpdate),
        false,
      );
      expect(
        model.hasGroupPermission(9, Permission.groupPeriodReceiptsPeriodRead),
        true,
      );
      // An unknown group denies.
      expect(
        model.hasGroupPermission(
            999, Permission.groupPeriodReceiptsPeriodRead),
        false,
      );
    });

    test('orApp override bypasses the group check', () {
      model.setPermissions(
        ['app.groups.read'],
        {
          '7': ['group.receipts.read'],
        },
      );

      // group.receipts.update is absent in group 7, but the app override is held.
      expect(
        model.hasGroupPermission(
          7,
          Permission.groupPeriodReceiptsPeriodUpdate,
          orApp: [Permission.appPeriodGroupsPeriodRead],
        ),
        true,
      );
      // Without the override, denied.
      expect(
        model.hasGroupPermission(7, Permission.groupPeriodReceiptsPeriodUpdate),
        false,
      );
      // Override not held -> falls through to the (absent) group permission.
      expect(
        model.hasGroupPermission(
          7,
          Permission.groupPeriodReceiptsPeriodUpdate,
          orApp: [Permission.appPeriodUsersPeriodCreate],
        ),
        false,
      );
    });

    test('hasGroupPermissionInAnyGroup checks across all groups', () {
      model.setPermissions(
        const [],
        {
          '7': ['group.receipts.read'],
          '9': ['group.receipts.create'],
        },
      );

      expect(
        model.hasGroupPermissionInAnyGroup(
            Permission.groupPeriodReceiptsPeriodCreate),
        true,
      );
      expect(
        model.hasGroupPermissionInAnyGroup(
            Permission.groupPeriodReceiptsPeriodDelete),
        false,
      );
    });

    test('notifies listeners on setPermissions', () {
      var notified = 0;
      model.addListener(() => notified++);
      model.setPermissions(const [], const {});
      expect(notified, 1);
    });

    // Forward-compatibility guard. A permission key added by a newer backend
    // must be inert on an older build, never fatal. When these rode the closed
    // `Permission` enum, an unrecognized key threw out of the generated
    // deserializer and failed the entire AppData parse, which hard-failed
    // login -- twice in production (2026-07-24 group.members.create;
    // 2026-08-06 group.members.grants.update).
    test('ignores unknown permission keys without disturbing known ones', () {
      model.setPermissions(
        ['app.users.read', 'app.not.yet.invented'],
        {
          '7': ['group.receipts.read', 'group.members.grants.update'],
        },
      );

      expect(model.hasAppPermission(Permission.appPeriodUsersPeriodRead), true);
      expect(
        model.hasGroupPermission(7, Permission.groupPeriodReceiptsPeriodRead),
        true,
      );
      // The unknown key grants nothing it wasn't asked to grant.
      expect(
        model.hasGroupPermission(7, Permission.groupPeriodReceiptsPeriodDelete),
        false,
      );
    });

    // Wildcard grants are supported by the matcher on both clients and by the
    // backend, and could never have been represented by the Permission enum.
    test('honors a wildcard grant', () {
      model.setPermissions(
        const [],
        {
          '7': ['group.receipts.*'],
        },
      );

      expect(
        model.hasGroupPermission(7, Permission.groupPeriodReceiptsPeriodUpdate),
        true,
      );
      expect(
        model.hasGroupPermission(7, Permission.groupPeriodCommentsPeriodCreate),
        false,
      );
    });
  });
}
