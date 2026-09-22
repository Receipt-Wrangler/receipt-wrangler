import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/utils/users.dart';

/// Option sources for the filter's list and users conditions.
///
/// Categories and tags need nothing here: the synthetic "All" group is a real
/// row the user is a member of (`GroupRepository.CreateAllGroup`), so
/// `GetAppData` builds it a `groupCategories` / `groupTags` entry like any other
/// group. `CategorySelectField` / `TagSelectField` therefore work unchanged with
/// the route's group id -- which is also what desktop does.
///
/// Paid-by is the exception: the All group's roster is just the caller, so its
/// members are not the people whose receipts the view actually spans.

/// Whether [groupId] is the synthetic "All" group -- a real row that spans the
/// user's groups rather than holding receipts or members of its own.
bool isAllGroupId(GroupModel groupModel, String groupId) {
  return groupModel.getGroupById(groupId)?.isAllGroup ?? false;
}

/// The people a receipt in this view could have been paid by.
///
/// A real group's own members, or -- on the All group, whose roster is just the
/// caller -- the union across every real group, de-duplicated by id and keeping
/// first-seen order so the picker is stable between opens.
///
/// Built from the per-group rosters rather than the flat `UserModel.users`,
/// which is the install-wide name table and would offer people the caller has
/// never shared a group with.
List<api.UserView> filterPaidByOptions(
    UserModel userModel, GroupModel groupModel, String groupId) {
  if (!isAllGroupId(groupModel, groupId)) {
    return getUsersInGroup(userModel, groupModel, groupId);
  }

  final seen = <int>{};
  final union = <api.UserView>[];

  for (final group in groupModel.groupsWithoutAllGroup) {
    for (final user
        in getUsersInGroup(userModel, groupModel, group.id.toString())) {
      if (seen.add(user.id)) {
        union.add(user);
      }
    }
  }

  return union;
}

/// The groups the "Group" condition offers.
///
/// Always `groupsWithoutAllGroup`: the All group is the *view* being narrowed,
/// never one of the values to narrow it to.
List<api.Group> filterGroupOptions(GroupModel groupModel) {
  return groupModel.groupsWithoutAllGroup;
}

/// Every receipt status a user can filter on.
///
/// `ReceiptStatus.empty` is excluded: it is the deserialization fallback for a
/// status this build predates, not a value anyone means to select.
List<api.ReceiptStatus> filterStatusOptions() {
  return api.ReceiptStatus.values
      .where((status) => status != api.ReceiptStatus.empty)
      .toList();
}
