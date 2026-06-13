import { Store } from "@ngxs/store";
import { catchError, Observable, of } from "rxjs";
import { Permission, Role, RoleService } from "../open-api";
import { AuthState } from "../store";

/**
 * Loads the roles a UI needs to resolve role ids to display names — but only
 * when the caller holds `app.roles.read`. Non-holders get an empty list WITHOUT
 * issuing `GET /role`: that request would 403, and the global HTTP interceptor
 * turns a 403 into a forced logout. This matters because some role-name readers
 * (the group-member table in `group-form`) are rendered for ordinary group
 * members who are not admins. Request errors still degrade to an empty list.
 */
export function loadAssignableRoles(
  store: Store,
  roleService: RoleService
): Observable<Role[]> {
  if (!store.selectSnapshot(AuthState.hasAppPermission(Permission.AppRolesRead))) {
    return of([]);
  }
  return roleService.getRoles().pipe(catchError(() => of([])));
}
