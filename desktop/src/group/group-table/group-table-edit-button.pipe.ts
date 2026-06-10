import { Pipe, PipeTransform } from "@angular/core";
import { Group, Permission } from "../../open-api";
import { hasAll, hasAny } from "../../utils/permission.utils";

@Pipe({
    name: "groupTableEditButton",
    standalone: false
})
export class GroupTableEditButtonPipe implements PipeTransform {

  public transform(
    group: Group,
    appPermissions: string[],
    groupPermissions: { [id: number]: string[] }
  ): string {
    if (hasAll(groupPermissions?.[group.id] ?? [], Permission.GroupUpdate)) {
      return `/groups/${group.id}/details/edit`;
    } else if (hasAny(appPermissions, Permission.AppGroupsUpdateSettings)) {
      return `/groups/${group.id}/settings/edit`;
    } else {
      return `/groups/${group.id}/details/view`;
    }
  }
}
