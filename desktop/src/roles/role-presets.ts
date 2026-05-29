import { PermissionScope } from "../open-api";

/**
 * A role applies to exactly one scope. The "Role type" choice drives the
 * available templates, the permission list, and the summary.
 */
export type RoleType = "app" | "group";

export interface RoleTypeMeta {
  id: RoleType;
  name: string;
  /** Material icon name. */
  icon: string;
  /** Accent foreground color for the icon tile. */
  color: string;
  /** Accent background color for the icon tile. */
  tint: string;
  description: string;
  scope: PermissionScope;
}

export interface RolePreset {
  id: string;
  name: string;
  description: string;
  /** Material icon name. */
  icon: string;
  color: string;
  tint: string;
  /**
   * Resolves the preset to a set of permission keys. The full registry keys
   * (already filtered to the active scope) are passed in so presets can be
   * computed dynamically and stay resilient to unknown keys.
   */
  resolve: (scopeKeys: readonly string[]) => Set<string>;
}

// Violet "Group" accent literals, matching the role-list page.
const GROUP_TINT = "#ede9fe";
const GROUP_COLOR = "#6d28d9";

// Primary tint/color literals (sky blue) sourced from variables.scss.
const PRIMARY_TINT = "#ccecff";
const PRIMARY_COLOR = "#0086d4";

export const ROLE_TYPES: RoleTypeMeta[] = [
  {
    id: "app",
    name: "Application role",
    icon: "apps",
    color: PRIMARY_COLOR,
    tint: PRIMARY_TINT,
    description:
      "Controls app-wide capabilities — users, settings, prompts, API keys and more.",
    scope: PermissionScope.App,
  },
  {
    id: "group",
    name: "Group role",
    icon: "workspaces",
    color: GROUP_COLOR,
    tint: GROUP_TINT,
    description:
      "Controls what a member can do inside a group — receipts, dashboards, activity.",
    scope: PermissionScope.Group,
  },
];

/**
 * Derives the resource key for a permission: the full key minus its last
 * `.segment`. e.g. `app.users.create` -> `app.users`, `group.view` -> `group`.
 */
export function resourceKeyOf(permissionKey: string): string {
  const lastDot = permissionKey.lastIndexOf(".");
  return lastDot === -1 ? permissionKey : permissionKey.slice(0, lastDot);
}

/** Filters the registry keys down to those that exist, preserving resilience. */
function intersect(wanted: readonly string[], available: readonly string[]): Set<string> {
  const availableSet = new Set(available);
  return new Set(wanted.filter((k) => availableSet.has(k)));
}

const USER_MANAGER_RESOURCES = new Set(["app.users", "app.api-keys", "app.system-emails"]);

const GROUP_MANAGER_KEYS = [
  "group.receipts.create",
  "group.receipts.read",
  "group.receipts.update",
  "group.receipts.delete",
  "group.receipts.duplicate",
  "group.receipts.magic-fill",
  "group.comments.create",
  "group.comments.delete",
  "group.dashboards.create",
  "group.dashboards.read",
  "group.dashboards.update",
  "group.dashboards.delete",
  "group.widgets.read",
  "group.activities.read",
  "group.activities.rerun",
  "group.view",
  "group.update",
  "group.export-ocr",
];

const RECEIPT_EDITOR_KEYS = [
  "group.receipts.create",
  "group.receipts.read",
  "group.receipts.update",
  "group.receipts.duplicate",
  "group.receipts.magic-fill",
  "group.comments.create",
  "group.dashboards.read",
  "group.widgets.read",
  "group.view",
];

const VIEWER_KEYS = [
  "group.receipts.read",
  "group.dashboards.read",
  "group.widgets.read",
  "group.activities.read",
  "group.view",
];

export const APP_PRESETS: RolePreset[] = [
  {
    id: "admin",
    name: "Administrator",
    description: "Everything in the app",
    icon: "shield_person",
    color: PRIMARY_COLOR,
    tint: PRIMARY_TINT,
    resolve: (keys) => new Set(keys),
  },
  {
    id: "user-mgr",
    name: "User Manager",
    description: "Manage users & keys",
    icon: "manage_accounts",
    color: PRIMARY_COLOR,
    tint: PRIMARY_TINT,
    resolve: (keys) =>
      new Set(keys.filter((k) => USER_MANAGER_RESOURCES.has(resourceKeyOf(k)))),
  },
  {
    id: "app-viewer",
    name: "Read Only",
    description: "View settings only",
    icon: "visibility",
    color: "#10b981",
    tint: "#d1fae5",
    resolve: (keys) => new Set(keys.filter((k) => k.endsWith(".read"))),
  },
  {
    id: "custom",
    name: "Custom",
    description: "Start from scratch",
    icon: "tune",
    color: "#64748b",
    tint: "#f1f5f9",
    resolve: () => new Set<string>(),
  },
];

export const GROUP_PRESETS: RolePreset[] = [
  {
    id: "manager",
    name: "Group Manager",
    description: "Run groups & receipts",
    icon: "workspaces",
    color: GROUP_COLOR,
    tint: GROUP_TINT,
    resolve: (keys) => intersect(GROUP_MANAGER_KEYS, keys),
  },
  {
    id: "editor",
    name: "Receipt Editor",
    description: "Add & edit receipts",
    icon: "edit_note",
    color: "#10b981",
    tint: "#d1fae5",
    resolve: (keys) => intersect(RECEIPT_EDITOR_KEYS, keys),
  },
  {
    id: "viewer",
    name: "Viewer",
    description: "Read-only access",
    icon: "visibility",
    color: "#64748b",
    tint: "#f1f5f9",
    resolve: (keys) => intersect(VIEWER_KEYS, keys),
  },
  {
    id: "custom",
    name: "Custom",
    description: "Start from scratch",
    icon: "tune",
    color: "#64748b",
    tint: "#f1f5f9",
    resolve: () => new Set<string>(),
  },
];

export const CUSTOM_PRESET_ID = "custom";

export function presetsForType(type: RoleType): RolePreset[] {
  return type === "app" ? APP_PRESETS : GROUP_PRESETS;
}

/**
 * Maps a resource key to a Material icon name. Derived from the design's
 * resource icon set; falls back to `lock` for unknown resources.
 */
const RESOURCE_ICONS: Record<string, string> = {
  "app.users": "person",
  "app.groups": "group",
  "app.categories": "sell",
  "app.tags": "label",
  "app.custom-fields": "data_object",
  "app.prompts": "auto_awesome",
  "app.receipt-processing-settings": "document_scanner",
  "app.system-settings": "settings",
  "app.system-emails": "mail",
  "app.system-tasks": "task_alt",
  "app.imports": "file_upload",
  "app.api-keys": "key",
  group: "workspaces",
  "group.receipts": "receipt_long",
  "group.comments": "chat_bubble",
  "group.dashboards": "dashboard",
  "group.widgets": "widgets",
  "group.activities": "history",
  "group.email": "alternate_email",
};

export function iconForResource(resourceKey: string): string {
  return RESOURCE_ICONS[resourceKey] ?? "lock";
}

const RESOURCE_NAME_OVERRIDES: Record<string, string> = {
  "custom-fields": "Custom Fields",
  "receipt-processing-settings": "Receipt Processing",
  "api-keys": "API Keys",
  "system-settings": "System Settings",
  "system-emails": "System Emails",
  "system-tasks": "System Tasks",
};

/** Prettifies the last segment of a resource key into a friendly title. */
export function friendlyResourceName(resourceKey: string): string {
  const lastDot = resourceKey.lastIndexOf(".");
  const segment = lastDot === -1 ? resourceKey : resourceKey.slice(lastDot + 1);

  if (RESOURCE_NAME_OVERRIDES[segment]) {
    return RESOURCE_NAME_OVERRIDES[segment];
  }

  return segment
    .split("-")
    .map((part) => part.charAt(0).toUpperCase() + part.slice(1))
    .join(" ");
}

/** Prettifies a permission key's action (last segment) into a label. */
export function friendlyActionLabel(permissionKey: string): string {
  const lastDot = permissionKey.lastIndexOf(".");
  const action = lastDot === -1 ? permissionKey : permissionKey.slice(lastDot + 1);
  const spaced = action.replace(/-/g, " ");
  return spaced.charAt(0).toUpperCase() + spaced.slice(1);
}
