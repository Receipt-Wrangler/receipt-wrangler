import { Pipe, PipeTransform } from "@angular/core";
import { Role } from "../open-api";

@Pipe({
    name: "roleName",
    standalone: false
})
export class RoleNamePipe implements PipeTransform {
  public transform(id: number | undefined | null, roles: Role[]): string {
    if (id == null) {
      return "";
    }
    return roles.find((role) => role.id === id)?.name ?? "";
  }
}
