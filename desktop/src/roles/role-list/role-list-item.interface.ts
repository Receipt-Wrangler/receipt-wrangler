/**
 * View-model for a single row in the roles list.
 *
 * Mapped from the generated `Role` model. A role belongs to exactly one scope
 * (application OR group). Member data is not populated yet — assigning roles to
 * users is a later slice — so `members`/`userCount` render as "No users".
 */
export interface RoleListItem {
  id: string;
  name: string;
  description: string;
  scope: RoleScope;
  permissionCount: number;
  members: RoleMember[];
  userCount: number;
  isSystem: boolean;
  icon: string;
  iconColor: string;
  iconTint: string;
}

export type RoleScope = "app" | "group";

/** The type filter above the table, plus the implicit "all". */
export type RoleListFilter = "all" | RoleScope;

export interface RoleMember {
  name: string;
  color: string;
}
