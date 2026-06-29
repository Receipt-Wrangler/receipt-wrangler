import { expect, test, type Page } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';

// Role management is admin-only (the /roles route is guarded by UserRole.Admin),
// so override the default regular-user storage state with the admin one.
test.use({ storageState: 'e2e/.auth/admin.json' });

function uniqueName(tag: string) {
  return `e2e-role-${tag}-${Date.now()}`;
}

type RoleType = 'app' | 'group';

// Roles created during a test are tracked here and removed in afterEach so the
// shared DB does not accumulate test data. The array is reset per test; under
// Playwright's worker model each worker runs its tests sequentially, so this
// module-level state is safe. Tests still use unique names and assert on the
// *specific* role they created (never absolute totals) to stay parallel-safe.
let createdRoles: string[] = [];

async function createRole(
  page: Page,
  opts: { name: string; type: RoleType; template: string },
) {
  createdRoles.push(opts.name);
  await page.goto('/roles');
  // Header "Add Role" button (the empty-state one, if present, is identical).
  await page.getByRole('button', { name: 'Add Role' }).first().click();
  await expect(page).toHaveURL(/\/roles\/new$/);
  await expect(page.getByLabel('Role Name')).toBeVisible();

  // Role type is the up-front choice (Application is the default). Clicking the
  // already-selected type is a no-op; switching resets the template/permissions.
  const typeName = opts.type === 'app' ? /Application role/ : /Group role/;
  await page.getByRole('button', { name: typeName }).click();

  await page.getByLabel('Role Name').fill(opts.name);

  // Seed permissions from a per-scope template, then save.
  await page.getByRole('button', { name: opts.template }).click();
  await page.getByRole('button', { name: 'Save Role' }).click();

  // Saving navigates back to the list.
  await expect(page).toHaveURL(/\/roles$/);
}

// Opens a role's editor from the list. The edit icon button has no accessible
// name (Material's mat-icon is aria-hidden and matTooltip sets no label), so it
// carries a data-testid.
async function openRoleEditor(page: Page, name: string) {
  await page.goto('/roles');
  const row = page.getByRole('row').filter({ hasText: name }).first();
  await expect(row).toBeVisible();
  await row.getByTestId('role-edit').click();
  await expect(page).toHaveURL(/\/roles\/\d+\/edit/);
  // The form populates from getRoles(); wait for the name before asserting more.
  await expect(page.getByLabel('Role Name')).toHaveValue(name);
}

// Deletes a role from the list by name and asserts the row is gone. The delete
// icon button has no accessible name, so it carries a data-testid. The
// confirmation dialog's confirm button is the icon-only "done" submit button.
async function deleteRole(page: Page, name: string) {
  await page.goto('/roles');
  const row = page.getByRole('row').filter({ hasText: name }).first();
  await expect(row).toBeVisible();
  await row.getByTestId('role-delete').click();

  const dialog = page.getByRole('dialog');
  await expect(dialog).toBeVisible();
  await dialog.getByTestId('dialog-submit-button').click();

  await expect(page.getByRole('row').filter({ hasText: name })).toHaveCount(0);
}

// Best-effort cleanup: removes every list row whose name contains `name`
// (catches renamed variants, e.g. "<name>-renamed"). Skips rows whose delete is
// disabled (system/assigned) and stops once none remain.
async function cleanupRole(page: Page, name: string) {
  for (let i = 0; i < 20; i++) {
    await page.goto('/roles');
    // The table renders only after getRoles() resolves; wait for it before
    // counting so a not-yet-loaded list isn't mistaken for "already gone".
    await expect(page.locator('app-table')).toBeVisible();
    const rows = page.getByRole('row').filter({ hasText: name });
    if ((await rows.count()) === 0) {
      return;
    }
    const del = rows.first().getByTestId('role-delete');
    if (await del.isDisabled()) {
      return;
    }
    await del.click();
    const dialog = page.getByRole('dialog');
    await dialog.getByTestId('dialog-submit-button').click();
    await expect(dialog).toHaveCount(0);
  }
}

// The summary panel's granted-permission count (driven by granted().size).
async function grantedCount(page: Page): Promise<number> {
  return Number((await page.getByTestId('granted-permission-count').innerText()).trim());
}

// The submit button lives in a fixed bottom bar and carries a matTooltip
// overlay, which makes a direct click flaky on tall forms. Submit through the
// form's implicit Enter handler from the Role Name field instead — a real user
// action that sidesteps the overlay and viewport edge cases.
async function saveRole(page: Page) {
  const name = page.getByLabel('Role Name');
  await expect(name).toBeEnabled();
  await name.press('Enter');
}

test.describe('roles', () => {
  test.beforeEach(async ({ page }) => {
    createdRoles = [];
    await stubTokenRefresh(page);
  });

  test.afterEach(async ({ page }) => {
    for (const name of createdRoles) {
      await cleanupRole(page, name);
    }
  });

  test('create an application role and see it in the table', async ({ page }) => {
    const name = uniqueName('app');
    await createRole(page, { name, type: 'app', template: 'Administrator' });

    const row = page.getByRole('row').filter({ hasText: name }).first();
    await expect(row).toBeVisible();
    // Single Type chip reflects the chosen scope.
    await expect(row).toContainText('Application');
  });

  test('create a group role and see it in the table', async ({ page }) => {
    const name = uniqueName('group');
    await createRole(page, { name, type: 'group', template: 'Group Manager' });

    const row = page.getByRole('row').filter({ hasText: name }).first();
    await expect(row).toBeVisible();
    await expect(row).toContainText('Group');
  });

  test('the type filter narrows the table to the selected scope', async ({ page }) => {
    const appName = uniqueName('filter-app');
    const groupName = uniqueName('filter-group');
    await createRole(page, { name: appName, type: 'app', template: 'Read Only' });
    await createRole(page, { name: groupName, type: 'group', template: 'Viewer' });

    await page.goto('/roles');
    const appRow = () => page.getByRole('row').filter({ hasText: appName });
    const groupRow = () => page.getByRole('row').filter({ hasText: groupName });

    // Both visible under "All roles" (the default).
    await expect(appRow().first()).toBeVisible();
    await expect(groupRow().first()).toBeVisible();

    // Filtering to Application hides the group role and vice versa. The filter
    // bar renders its options as ARIA tabs.
    await page.getByRole('tab', { name: 'Application' }).click();
    await expect(appRow().first()).toBeVisible();
    await expect(groupRow()).toHaveCount(0);

    await page.getByRole('tab', { name: 'Group' }).click();
    await expect(groupRow().first()).toBeVisible();
    await expect(appRow()).toHaveCount(0);
  });

  test('Add Role opens the create form, Cancel returns without creating', async ({ page }) => {
    await page.goto('/roles');
    await page.getByRole('button', { name: 'Add Role' }).first().click();
    await expect(page).toHaveURL(/\/roles\/new$/);

    await page.getByRole('button', { name: 'Cancel' }).click();
    await expect(page).toHaveURL(/\/roles$/);
  });

  test('editing an application role persists name, description, and permissions', async ({ page }) => {
    const original = uniqueName('edit-app');
    await createRole(page, { name: original, type: 'app', template: 'Read Only' });

    await openRoleEditor(page, original);
    const readOnlyCount = await grantedCount(page);

    // Granting everything via the Administrator preset is a deterministic
    // increase over Read Only, without depending on the absolute permission total.
    await page.getByRole('button', { name: 'Administrator' }).click();
    await expect.poll(() => grantedCount(page)).toBeGreaterThan(readOnlyCount);
    const adminCount = await grantedCount(page);

    const renamed = `${original}-renamed`;
    const description = `desc-${Date.now()}`;
    await page.getByLabel('Role Name').fill(renamed);
    await page.getByLabel('Description').fill(description);
    await saveRole(page);
    await expect(page).toHaveURL(/\/roles$/);

    const row = page.getByRole('row').filter({ hasText: renamed }).first();
    await expect(row).toBeVisible();
    await expect(row).toContainText('Application');

    // Re-open and confirm everything round-tripped through the PUT.
    await openRoleEditor(page, renamed);
    await expect(page.getByLabel('Description')).toHaveValue(description);
    await expect.poll(() => grantedCount(page)).toBe(adminCount);
  });

  test('the role type is locked when editing', async ({ page }) => {
    const name = uniqueName('lock');
    await createRole(page, { name, type: 'app', template: 'Read Only' });

    await openRoleEditor(page, name);
    // The edit URL carries the scope that disambiguates app/group ids.
    await expect(page).toHaveURL(/\/roles\/\d+\/edit\?scope=app/);

    // Neither type card can be chosen — an app role stays an app role.
    await expect(page.getByRole('button', { name: /Application role/ })).toBeDisabled();
    await expect(page.getByRole('button', { name: /Group role/ })).toBeDisabled();
  });

  test('editing a group role updates its name and keeps it group-scoped', async ({ page }) => {
    const original = uniqueName('edit-group');
    await createRole(page, { name: original, type: 'group', template: 'Group Manager' });

    await openRoleEditor(page, original);
    await expect(page).toHaveURL(/\/roles\/\d+\/edit\?scope=group/);
    await expect(page.getByRole('button', { name: /Group role/ })).toBeDisabled();

    const renamed = `${original}-renamed`;
    await page.getByLabel('Role Name').fill(renamed);
    await saveRole(page);
    await expect(page).toHaveURL(/\/roles$/);

    const row = page.getByRole('row').filter({ hasText: renamed }).first();
    await expect(row).toBeVisible();
    await expect(row).toContainText('Group');
  });

  test('create a role then delete it removes it from the table', async ({ page }) => {
    const name = uniqueName('delete');
    await createRole(page, { name, type: 'app', template: 'Read Only' });

    // A freshly created, unassigned role is deletable.
    await page.goto('/roles');
    const row = page.getByRole('row').filter({ hasText: name }).first();
    await expect(row.getByTestId('role-delete')).toBeEnabled();

    await deleteRole(page, name);
  });

  test('cancelling an edit discards changes', async ({ page }) => {
    const original = uniqueName('cancel');
    await createRole(page, { name: original, type: 'app', template: 'Read Only' });

    await openRoleEditor(page, original);
    await page.getByLabel('Role Name').fill(`${original}-DISCARDED`);
    await page.getByRole('button', { name: 'Cancel' }).click();
    await expect(page).toHaveURL(/\/roles$/);

    // The original name survives; the discarded edit never persisted.
    await expect(
      page.getByRole('row').filter({ hasText: original }).first(),
    ).toBeVisible();
    await expect(
      page.getByRole('row').filter({ hasText: `${original}-DISCARDED` }),
    ).toHaveCount(0);
  });
});
