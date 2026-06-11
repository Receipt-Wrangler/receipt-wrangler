import 'package:built_collection/built_collection.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/shared/functions/permissions.dart';

PermissionsModel _modelWithGroup(int groupId, List<Permission> permissions) {
  final model = PermissionsModel();
  model.setPermissions(
    BuiltList<Permission>(),
    BuiltMap<String, BuiltList<Permission>>({
      '$groupId': BuiltList<Permission>(permissions),
    }),
  );
  return model;
}

void main() {
  group('canEditReceipt', () {
    test('true when the user holds group.receipts.update in the group', () {
      final model =
          _modelWithGroup(3, [Permission.groupPeriodReceiptsPeriodUpdate]);
      expect(canEditReceipt(model, 3), true);
    });

    test('false when the group only grants read', () {
      final model =
          _modelWithGroup(3, [Permission.groupPeriodReceiptsPeriodRead]);
      expect(canEditReceipt(model, 3), false);
    });

    test('false for a group the user is not a member of', () {
      final model =
          _modelWithGroup(3, [Permission.groupPeriodReceiptsPeriodUpdate]);
      expect(canEditReceipt(model, 99), false);
    });
  });
}
