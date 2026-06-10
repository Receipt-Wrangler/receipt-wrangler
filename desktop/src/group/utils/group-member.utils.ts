import { FormControl, FormGroup, Validators } from "@angular/forms";
import { GroupMember, GroupRole, Role } from "../../open-api";

// Maps a modern group role onto the legacy GroupRole enum by its (immutable)
// system-role name. Transitional bridge so the group form's member table and its
// "keep an owner" check — which still read the legacy enum — keep working for
// newly assigned members; the backend derives the same value on save. Custom
// roles fall back to the least-privilege VIEWER.
export function legacyGroupRoleFromRole(role: Role): GroupRole {
  switch (role.name) {
    case "Legacy Owner":
      return GroupRole.Owner;
    case "Legacy Editor":
      return GroupRole.Editor;
    default:
      return GroupRole.Viewer;
  }
}

export function buildGroupMemberForm(groupMember?: GroupMember): FormGroup {
  return new FormGroup({
    userId: new FormControl(groupMember?.userId ?? "", Validators.required),
    // The member's modern group role is what the user now selects; it is the
    // required choice and is what the backend persists.
    groupRoleId: new FormControl(
      groupMember?.groupRoleId ?? null,
      Validators.required
    ),
    // The legacy enum is retained transitionally — it is carried (not edited)
    // for existing members so the parent group form's "must keep an owner" check
    // and any other legacy readers keep working until those are migrated. The
    // backend derives it from groupRoleId on save.
    groupRole: new FormControl(groupMember?.groupRole ?? ""),
    groupId: new FormControl(groupMember?.groupId ?? undefined),
  });
}
