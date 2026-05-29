/**
 * View-model for a single row in the roles list.
 *
 * This is a temporary, frontend-only shape used to render the list page until
 * the backend role-list endpoint and its generated client exist. Once that
 * lands, this will be replaced by (or mapped from) the generated role model.
 */
export interface RoleListItem {
  id: string;
  name: string;
  description: string;
  scopes: RoleScope[];
  appCount: number;
  groupCount: number;
  members: RoleMember[];
  userCount: number;
  isSystem: boolean;
  icon: string;
  iconColor: string;
  iconTint: string;
}

export type RoleScope = "app" | "group";

export interface RoleMember {
  name: string;
  color: string;
}
