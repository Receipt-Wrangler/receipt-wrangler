import { FormControl, FormGroup, Validators } from "@angular/forms";
import { GroupMember } from "../../open-api";

export function buildGroupMemberForm(groupMember?: GroupMember): FormGroup {
  return new FormGroup({
    userId: new FormControl(groupMember?.userId ?? "", Validators.required),
    // The member's modern group role is what the user selects; it is the
    // required choice and is what the backend persists.
    groupRoleId: new FormControl(
      groupMember?.groupRoleId ?? null,
      Validators.required
    ),
    groupId: new FormControl(groupMember?.groupId ?? undefined),
  });
}
