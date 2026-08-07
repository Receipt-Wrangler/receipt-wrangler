import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/shared/functions/permissions.dart';

import '../../helpers/permission_test_helpers.dart';

PermissionsModel _modelWithGroup(int groupId, List<Permission> permissions) =>
    seededPermissions(group: {groupId: permissions});

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

  group('canCommentCreate', () {
    test('true when the user holds group.comments.create in the group', () {
      final model =
          _modelWithGroup(4, [Permission.groupPeriodCommentsPeriodCreate]);
      expect(canCommentCreate(model, 4), true);
    });

    test('false when the group does not grant comment create', () {
      final model =
          _modelWithGroup(4, [Permission.groupPeriodReceiptsPeriodRead]);
      expect(canCommentCreate(model, 4), false);
    });

    test('false for a group the user is not a member of', () {
      final model =
          _modelWithGroup(4, [Permission.groupPeriodCommentsPeriodCreate]);
      expect(canCommentCreate(model, 99), false);
    });
  });

  group('canCommentDelete', () {
    test('true when the user holds group.comments.delete in the group', () {
      final model =
          _modelWithGroup(4, [Permission.groupPeriodCommentsPeriodDelete]);
      expect(canCommentDelete(model, 4), true);
    });

    test('false when the group only grants comment create', () {
      final model =
          _modelWithGroup(4, [Permission.groupPeriodCommentsPeriodCreate]);
      expect(canCommentDelete(model, 4), false);
    });
  });
}
