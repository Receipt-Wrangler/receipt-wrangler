import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter_options.dart';

import '../helpers/receipt_form_test_helpers.dart';

/// The All group is a real row the user belongs to, so its category and tag
/// catalogs arrive on AppData like any other group's -- those pickers need no
/// special casing at all. Its *roster*, though, is just the caller, which is why
/// paid-by is the one option source with a rule of its own.
void main() {
  const allGroupId = 1;
  const householdId = 2;
  const officeId = 3;

  late GroupModel groupModel;
  late UserModel userModel;

  setUp(() {
    groupModel = GroupModel();
    userModel = UserModel();

    userModel.setUsers([
      buildUserView(id: 10, displayName: "Noah Hall"),
      buildUserView(id: 11, displayName: "Dana Kim"),
      buildUserView(id: 12, displayName: "Priya Shah"),
    ]);

    groupModel.setGroups([
      buildGroup(
        id: allGroupId,
        name: "All",
        isAllGroup: true,
        // The All group is created for its owner alone.
        members: [buildGroupMember(userId: 10, groupId: allGroupId)],
      ),
      buildGroup(id: householdId, name: "Household", members: [
        buildGroupMember(userId: 10, groupId: householdId),
        buildGroupMember(userId: 11, groupId: householdId),
      ]),
      buildGroup(id: officeId, name: "Office", members: [
        buildGroupMember(userId: 11, groupId: officeId),
        buildGroupMember(userId: 12, groupId: officeId),
      ]),
    ]);
  });

  group("isAllGroupId", () {
    test("is true only for the synthetic group", () {
      expect(isAllGroupId(groupModel, "$allGroupId"), isTrue);
      expect(isAllGroupId(groupModel, "$householdId"), isFalse);
    });

    test("is false for a group the model does not know", () {
      expect(isAllGroupId(groupModel, "999"), isFalse);
      expect(isAllGroupId(groupModel, "not-a-number"), isFalse);
    });
  });

  group("filterPaidByOptions", () {
    test("a real group offers exactly its own roster", () {
      final names = filterPaidByOptions(userModel, groupModel, "$householdId")
          .map((user) => user.displayName)
          .toList();

      expect(names, ["Noah Hall", "Dana Kim"]);
    });

    test("the All group unions every real group's roster", () {
      // Its own roster is just the caller, which would make a paid-by filter
      // there useless -- the view spans receipts paid by anyone in any group.
      final names = filterPaidByOptions(userModel, groupModel, "$allGroupId")
          .map((user) => user.displayName)
          .toList();

      expect(names, ["Noah Hall", "Dana Kim", "Priya Shah"]);
    });

    test("someone in two groups appears once", () {
      final ids = filterPaidByOptions(userModel, groupModel, "$allGroupId")
          .map((user) => user.id)
          .toList();

      expect(ids.toSet().length, ids.length);
    });

    test("an unknown group offers nobody rather than throwing", () {
      expect(filterPaidByOptions(userModel, groupModel, "999"), isEmpty);
    });
  });

  group("filterGroupOptions", () {
    test("offers the real groups and never the All group itself", () {
      // The All group is the view being narrowed, not a value to narrow it to.
      expect(filterGroupOptions(groupModel).map((group) => group.name).toList(),
          ["Household", "Office"]);
    });
  });

  group("filterStatusOptions", () {
    test("offers every real status", () {
      expect(filterStatusOptions(), [
        api.ReceiptStatus.OPEN,
        api.ReceiptStatus.NEEDS_ATTENTION,
        api.ReceiptStatus.RESOLVED,
        api.ReceiptStatus.DRAFT,
        api.ReceiptStatus.DECLINED,
      ]);
    });

    test("never offers the empty fallback status", () {
      // ReceiptStatus.empty exists so an unknown wire value degrades instead of
      // failing the whole payload; it is not something a user selects.
      expect(filterStatusOptions(), isNot(contains(api.ReceiptStatus.empty)));
    });
  });
}
