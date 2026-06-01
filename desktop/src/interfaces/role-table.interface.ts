import { PagedTableInterface } from "src/interfaces/paged-table.interface";
import { RoleListFilter } from "src/roles/role-list/role-list-item.interface";

/**
 * Paged-table state for the roles list. Extends the shared paging state with the
 * scope filter ("all" | "app" | "group") that backs the filter bar above the
 * table; it is mapped to the API's PermissionScope when building the request.
 */
export interface RoleTableInterface extends PagedTableInterface {
  scope: RoleListFilter;
}
