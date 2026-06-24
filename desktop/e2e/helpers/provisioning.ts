import {
  type APIRequestContext,
  expect,
  type Page,
  request as pwRequest,
} from '@playwright/test';
import { creds } from './auth';

// Shared provisioning for permission e2e specs. Roles/users/groups are CREATED
// through the real app UI (the realistic flow under test); TEARDOWN goes through
// the admin API. UI teardown can't reliably remove a custom role: the role-list
// delete button is disabled while the role is assigned, and the UI's *bulk*
// user-delete dummy-converts a group-owning user (keeping the role assigned), so
// the role is never freed. `DELETE /api/user/{id}` hard-deletes, and deleting a
// group frees its members' group-role assignment — so the API tears down cleanly.

export function uniqueName(tag: string): string {
  return `e2e-${tag}-${Date.now()}`;
}

export interface CreateRoleOptions {
  name: string;
  /** Role-type card to pick, e.g. 'Application role' | 'Group role'. */
  type: 'Application role' | 'Group role';
  /** Starting template button, matched by (sub)string, e.g. 'Start from scratch' | 'Viewer'. */
  preset: string;
  /** Resource categories to flip fully on, e.g. ['Account', 'Notifications', 'Groups']. */
  enableCategories?: string[];
  /**
   * Individual permissions to switch OFF after applying the preset. Identify the
   * accordion by its resource key (key minus last segment, e.g. 'group.dashboards')
   * and the row by its label (e.g. 'Read Dashboards').
   */
  disablePermissions?: { panelKey: string; label: string }[];
}

/**
 * Creates a role via the role form: pick type → name it → apply a preset →
 * optionally flip whole resource groups on → optionally switch individual
 * permissions off. Leaves the browser on `/roles`.
 */
export async function createRole(page: Page, opts: CreateRoleOptions): Promise<void> {
  await page.goto('/roles');
  await page.getByRole('button', { name: 'Add Role' }).first().click();
  await expect(page).toHaveURL(/\/roles\/new$/);
  await expect(page.getByLabel('Role Name')).toBeVisible();

  // Pick the scope first — it swaps the available presets.
  await page.getByRole('button', { name: new RegExp(opts.type) }).click();
  await page.getByLabel('Role Name').fill(opts.name);
  await page.getByRole('button', { name: opts.preset }).click();

  for (const category of opts.enableCategories ?? []) {
    await page
      .getByRole('button', { name: `Toggle all ${category} permissions` })
      .click();
  }

  for (const { panelKey, label } of opts.disablePermissions ?? []) {
    // Expand the accordion by clicking its sub-text (the resource key is unique
    // on the page, so this won't collide with e.g. a sidebar "Dashboards" link),
    // then switch the individual permission off.
    const subText = new RegExp(`${panelKey.replace(/\./g, '\\.')} \\u00b7`);
    await page.getByText(subText).click();
    await page.getByRole('button', { name: `Toggle ${label}` }).click();
  }

  await page.getByRole('button', { name: 'Save Role' }).click();
  await expect(page).toHaveURL(/\/roles$/);
}

/** Creates a user assigned [opts.role] via the Create User dialog. */
export async function createUserWithRole(
  page: Page,
  opts: { username: string; password: string; role: string },
): Promise<void> {
  await page.goto('/users');
  await page.getByRole('button', { name: 'Create User' }).click();
  const dialog = page.getByRole('dialog').filter({ hasText: 'Create User' });
  await expect(dialog).toBeVisible();
  await dialog.getByLabel('Username').fill(opts.username);
  await dialog.getByLabel('Displayname').fill(opts.username);
  await dialog.getByLabel('Password').fill(opts.password);
  await dialog.getByRole('combobox', { name: 'App Role' }).click();
  // The option panel is a floating overlay rendered on the page, not the dialog.
  await page.getByRole('option', { name: opts.role, exact: true }).click();
  await dialog.locator('app-submit-button button').click();
  await expect(dialog).toBeHidden();
}

/**
 * Creates a group with [opts.memberDisplayName] added as [opts.roleName], and
 * returns the new group's id. Mirrors the group-member form flow.
 */
export async function createGroupWithMember(
  page: Page,
  opts: { groupName: string; memberDisplayName: string; roleName: string },
): Promise<string> {
  await page.goto('/groups/create');
  await expect(page.getByLabel('Group Name')).toBeVisible();
  await page.getByLabel('Group Name').fill(opts.groupName);

  await page.getByTestId('add-group-member').click();
  const memberDialog = page
    .getByRole('dialog')
    .filter({ hasText: 'Add Group Member' });
  await expect(memberDialog).toBeVisible();

  // The user autocomplete filters/displays by display name; while the option
  // panel is open the listbox shares the field's aria label, so target the combobox.
  const userField = memberDialog.getByRole('combobox', { name: 'User' });
  await userField.fill(opts.memberDisplayName);
  await page
    .getByRole('option', { name: opts.memberDisplayName, exact: true })
    .click();

  // Role defaults to "Legacy Owner" — switch it to the chosen role.
  const roleSelect = memberDialog.getByRole('combobox', { name: 'Role' });
  await roleSelect.click();
  await page.getByRole('option', { name: opts.roleName, exact: true }).click();

  await memberDialog.getByTestId('dialog-submit-button').click();
  await expect(memberDialog).toBeHidden();

  // The page-level save is icon-only; the app-form submits on Enter from a field.
  await page.getByLabel('Group Name').press('Enter');
  await page.waitForURL(/\/groups\/\d+\/details\/view/);
  return page.url().match(/\/groups\/(\d+)\//)![1];
}

// --- API teardown -----------------------------------------------------------
//
// Cleanup runs over the admin API (see the file header for why the UI can't do
// it). `request.newContext` against E2E_BASE_URL reuses the dev-server proxy to
// `/api`; the admin login cookie is stored on the context for later DELETEs.

const apiBaseUrl = (): string =>
  process.env.E2E_BASE_URL ?? 'http://localhost:4200';

/**
 * Opens an admin-authenticated `APIRequestContext`, runs [fn], and disposes it.
 * Use in `afterAll` for best-effort teardown of provisioned roles/users/groups.
 */
export async function withAdminApi<T>(
  fn: (api: APIRequestContext) => Promise<T>,
): Promise<T> {
  const api = await pwRequest.newContext({ baseURL: apiBaseUrl() });
  try {
    const { username, password } = creds('admin');
    const res = await api.post('/api/login', { data: { username, password } });
    if (!res.ok()) {
      throw new Error(`admin API login failed: HTTP ${res.status()}`);
    }
    return await fn(api);
  } finally {
    await api.dispose();
  }
}

/** Hard-deletes the user with [username] (frees any app-role assignment). */
export async function apiDeleteUserByName(
  api: APIRequestContext,
  username: string,
): Promise<void> {
  const users = (await (await api.get('/api/user/')).json()) as {
    id: number;
    username: string;
  }[];
  const user = users.find((u) => u.username === username);
  if (user) {
    await api.delete(`/api/user/${user.id}`);
  }
}

/** Deletes the group [groupId] (frees its members' group-role assignment). */
export async function apiDeleteGroupById(
  api: APIRequestContext,
  groupId: string,
): Promise<void> {
  await api.delete(`/api/group/${groupId}`);
}

/**
 * Deletes the [scope] role named [name]. Only succeeds once it's unassigned, so
 * call after the user/group that referenced it is gone.
 */
export async function apiDeleteRoleByName(
  api: APIRequestContext,
  name: string,
  scope: 'APP' | 'GROUP',
): Promise<void> {
  const roles = (await (await api.get('/api/role')).json()) as {
    id: number;
    name: string;
    scope: string;
  }[];
  const role = roles.find((r) => r.name === name && r.scope === scope);
  if (role) {
    await api.delete(`/api/role/${role.id}?scope=${scope}`);
  }
}
