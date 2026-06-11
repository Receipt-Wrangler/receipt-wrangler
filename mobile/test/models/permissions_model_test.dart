import 'package:built_collection/built_collection.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';

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

    test('hydrates app and group permissions from AppData enums', () {
      model.setPermissions(
        BuiltList<Permission>([Permission.appPeriodUsersPeriodRead]),
        BuiltMap<String, BuiltList<Permission>>({
          '7': BuiltList<Permission>([
            Permission.groupPeriodReceiptsPeriodUpdate,
            Permission.groupPeriodReceiptsPeriodCreate,
          ]),
          '9': BuiltList<Permission>([
            Permission.groupPeriodReceiptsPeriodRead,
          ]),
        }),
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
        BuiltList<Permission>([Permission.appPeriodGroupsPeriodRead]),
        BuiltMap<String, BuiltList<Permission>>({
          '7': BuiltList<Permission>(
              [Permission.groupPeriodReceiptsPeriodRead]),
        }),
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
        BuiltList<Permission>(),
        BuiltMap<String, BuiltList<Permission>>({
          '7': BuiltList<Permission>(
              [Permission.groupPeriodReceiptsPeriodRead]),
          '9': BuiltList<Permission>(
              [Permission.groupPeriodReceiptsPeriodCreate]),
        }),
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
      model.setPermissions(
        BuiltList<Permission>(),
        BuiltMap<String, BuiltList<Permission>>(),
      );
      expect(notified, 1);
    });
  });
}
