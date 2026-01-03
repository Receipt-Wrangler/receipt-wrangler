import { FormControl, FormGroup, Validators } from "@angular/forms";
import { GroupMember } from "../../open-api";

export function buildGroupMemberForm(groupMember?: GroupMember): FormGroup {
  return new FormGroup({
    userId: new FormControl(groupMember?.userId ?? "", Validators.required),
    groupRole: new FormControl(
      groupMember?.groupRole ?? "",
      Validators.required
    ),
    groupId: new FormControl(groupMember?.groupId ?? undefined),
  });
}
