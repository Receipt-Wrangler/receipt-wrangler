import {
  type APIRequestContext,
  expect,
  type Page,
  request as pwRequest,
} from '@playwright/test';
import { creds, type Role } from './auth';

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
  /**
   * Group-role paid-by visibility (the "Visible paid-by users" picker). Pin
   * "Their own receipts" and/or specific users by display name. Leaving both
   * unset keeps the role unrestricted (members see every payer's receipts).
   */
  paidByOwn?: boolean;
  paidByUsers?: string[];
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

  // Expand each accordion at most once (clicking its sub-text toggles it), then
  // switch off every listed permission under it. Opening a panel per entry would
  // re-close it when two permissions share a panel (e.g. two app.reports perms).
  const openedPanels = new Set<string>();
  for (const { panelKey, label } of opts.disablePermissions ?? []) {
    if (!openedPanels.has(panelKey)) {
      // The resource key is unique on the page, so this won't collide with e.g. a
      // sidebar "Dashboards" link.
      const subText = new RegExp(`${panelKey.replace(/\./g, '\\.')} \\u00b7`);
      await page.getByText(subText).click();
      openedPanels.add(panelKey);
    }
    await page.getByRole('button', { name: `Toggle ${label}` }).click();
  }

  if (opts.paidByOwn || (opts.paidByUsers?.length ?? 0) > 0) {
    // The picker is a multi-select autocomplete; while its panel is open the
    // listbox shares the field's label, so target the combobox role.
    const paidByField = page.getByRole('combobox', {
      name: 'Visible paid-by users',
    });
    if (opts.paidByOwn) {
      await paidByField.click();
      await page
        .getByRole('option', { name: 'Their own receipts', exact: true })
        .click();
    }
    for (const displayName of opts.paidByUsers ?? []) {
      await paidByField.click();
      await paidByField.fill(displayName);
      await page.getByRole('option', { name: displayName, exact: true }).click();
    }
    // Selecting an option re-opens the panel; close it so Save is clickable.
    await page.keyboard.press('Escape');
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
 * Opens an `APIRequestContext` authenticated as the given e2e [role], runs [fn],
 * and disposes it. Use to drive API calls AS a specific user — e.g. asserting a
 * restricted user's request is denied, or admin teardown.
 */
export async function withApiAs<T>(
  role: Role,
  fn: (api: APIRequestContext) => Promise<T>,
): Promise<T> {
  const api = await pwRequest.newContext({ baseURL: apiBaseUrl() });
  try {
    const { username, password } = creds(role);
    const res = await api.post('/api/login', { data: { username, password } });
    if (!res.ok()) {
      throw new Error(`${role} API login failed: HTTP ${res.status()}`);
    }
    return await fn(api);
  } finally {
    await api.dispose();
  }
}

/**
 * Opens an admin-authenticated `APIRequestContext`, runs [fn], and disposes it.
 * Use in `afterAll` for best-effort teardown of provisioned roles/users/groups.
 */
export async function withAdminApi<T>(
  fn: (api: APIRequestContext) => Promise<T>,
): Promise<T> {
  return withApiAs('admin', fn);
}

/** Returns the id of the user with [username] (via the admin user list). */
export async function apiGetUserId(
  api: APIRequestContext,
  username: string,
): Promise<number> {
  const users = (await (await api.get('/api/user/')).json()) as {
    id: number;
    username: string;
  }[];
  const user = users.find((u) => u.username === username);
  if (!user) {
    throw new Error(`user ${username} not found`);
  }
  return user.id;
}

/** Creates a minimal OPEN receipt and returns its id. */
export async function apiCreateReceipt(
  api: APIRequestContext,
  opts: { groupId: number | string; paidByUserId: number; name: string },
): Promise<number> {
  const res = await api.post('/api/receipt', {
    data: {
      name: opts.name,
      amount: '10.00',
      date: '2024-01-01T00:00:00Z',
      groupId: Number(opts.groupId),
      paidByUserId: opts.paidByUserId,
      status: 'OPEN',
    },
  });
  if (!res.ok()) {
    throw new Error(
      `create receipt failed: HTTP ${res.status()} ${await res.text()}`,
    );
  }
  const receipt = (await res.json()) as { id: number };
  return receipt.id;
}

/**
 * Creates a comment on [receiptId] and returns its id. The comment is owned by
 * whoever the [api] context is authenticated as (the per-comment delete control
 * only renders for the logged-in user's own comments).
 */
export async function apiCreateComment(
  api: APIRequestContext,
  opts: { receiptId: number; comment: string },
): Promise<number> {
  const res = await api.post('/api/comment/', {
    data: { comment: opts.comment, receiptId: opts.receiptId },
  });
  if (!res.ok()) {
    throw new Error(
      `create comment failed: HTTP ${res.status()} ${await res.text()}`,
    );
  }
  const comment = (await res.json()) as { id: number };
  return comment.id;
}

/**
 * Creates a group via the admin API (the caller is auto-added as owner) and
 * returns its id + name. Create only needs name+status; members are required
 * on update, not create.
 */
export async function apiCreateGroup(
  api: APIRequestContext,
  name: string,
): Promise<{ id: number; name: string }> {
  const res = await api.post('/api/group', {
    data: { name, status: 'ACTIVE' },
  });
  if (!res.ok()) {
    throw new Error(
      `create group failed: HTTP ${res.status()} ${await res.text()}`,
    );
  }
  const group = (await res.json()) as { id: number; name: string };
  return { id: group.id, name: group.name };
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

/** Returns a group id (as a string) the caller can generate reports over. */
export async function apiFirstReportGroupId(api: APIRequestContext): Promise<string> {
  const appData = (await (await api.get('/api/user/appData')).json()) as {
    groupPermissions?: Record<string, string[]>;
  };
  const groupPermissions = appData.groupPermissions ?? {};
  const groupId = Object.keys(groupPermissions).find((id) =>
    (groupPermissions[id] ?? []).includes('group.reports.read'),
  );
  if (!groupId) {
    throw new Error('no group with group.reports.read for the caller');
  }
  return groupId;
}

/**
 * Creates a report template via the API and returns its id + name (requires
 * app.reports.create on the caller). The body is a complete, buildable
 * ReportRequestCommand — records mode over a group the caller can report on.
 */
export async function apiCreateReportTemplate(
  api: APIRequestContext,
  opts?: { name?: string; groupIds?: string[]; formats?: string[] },
): Promise<{ id: number; name: string }> {
  const name = opts?.name ?? uniqueName('report-template');
  // Scope a real group the caller can report on, so generating over it isn't
  // group-gated out (the seeded admin's group ids are not necessarily [1]).
  const groupIds = opts?.groupIds ?? [await apiFirstReportGroupId(api)];
  const res = await api.post('/api/report/template', {
    data: {
      name,
      groupIds,
      period: { preset: 'this_month' },
      detail: { mode: 'records' },
      columns: [{ kind: 'dimension', name: 'Name', label: 'Name', field: 'name' }],
      formats: opts?.formats ?? ['csv'],
    },
  });
  if (!res.ok()) {
    throw new Error(
      `create report template failed: HTTP ${res.status()} ${await res.text()}`,
    );
  }
  const template = (await res.json()) as { id: number; name: string };
  return { id: template.id, name: template.name };
}

/** Deletes the report template [id] (requires app.reports.delete on the caller). */
export async function apiDeleteReportTemplateById(
  api: APIRequestContext,
  id: number | string,
): Promise<void> {
  await api.delete(`/api/report/template/${id}`);
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
