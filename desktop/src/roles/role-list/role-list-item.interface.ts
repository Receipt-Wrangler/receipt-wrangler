/**
 * View-model for a single row in the roles list.
 *
 * Mapped from the generated `Role` model. A role belongs to exactly one scope
 * (application OR group). `userCount` is how many users/group members are
 * currently assigned the role; it drives both the Members column and whether
 * the role can be deleted (assigned roles cannot be deleted).
 */
export interface RoleListItem {
  id: string;
  name: string;
  description: string;
  scope: RoleScope;
  permissionCount: number;
  userCount: number;
  isDefault: boolean;
  isSystem: boolean;
  icon: string;
  iconColor: string;
  iconTint: string;
}

export type RoleScope = "app" | "group";

/** The type filter above the table, plus the implicit "all". */
export type RoleListFilter = "all" | RoleScope;
