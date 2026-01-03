import 'package:openapi/openapi.dart';
import '../../models/auth_model.dart';
import '../../models/group_model.dart';

bool canEditReceipt(AuthModel authModel, GroupModel groupModel, int groupId) {
  var userId = authModel.claims?.userId;
  if (userId == null) {
    return false;
  }

  var groupRole = getUserRoleInGroup(userId, groupId, groupModel);
  if (groupRole == null) {
    return false;
  }

  return groupRole == GroupRole.EDITOR || groupRole == GroupRole.OWNER;
}

GroupRole? getUserRoleInGroup(int userId, int groupId, GroupModel groupModel) {
  var group = groupModel.getGroupById(groupId.toString());
  if (group == null) {
    return null;
  }

  var groupMember = group.groupMembers.firstWhere(
    (groupMember) => groupMember.userId == userId,
  );
  if (groupMember == null) {
    return null;
  }

  return groupMember.groupRole;
}
