import { Pipe, PipeTransform } from "@angular/core";
import { Group, GroupRole } from "../../open-api";
import { GroupUtil } from "../../utils/index";

@Pipe({
    name: "groupTableEditButton",
    standalone: false
})
export class GroupTableEditButtonPipe implements PipeTransform {

  constructor(private groupUtil: GroupUtil) {}

  public transform(group: Group, isAdmin: boolean): {
    routerLink: string[]
    queryParams: any
  } {
    const isGroupOwner = this.groupUtil.hasGroupAccess(
      group.id,
      GroupRole.Owner,
      false,
      false
    );

    if (isGroupOwner) {
      return {
        routerLink: [`/groups/${group.id}/details/edit`],
        queryParams: { tab: "details" }
      };
    } else if (isAdmin) {
      return {
        routerLink: [`/groups/${group.id}/settings/edit`],
        queryParams: { tab: "settings" }
      };
    } else {
      return {
        routerLink: [`/groups/${group.id}/details/view`],
        queryParams: { tab: "details" }
      };
    }
  }
}
