import { expect, test, type Page } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';

// Role management is admin-only (the /roles route is guarded by UserRole.Admin),
// so override the default regular-user storage state with the admin one.
test.use({ storageState: 'e2e/.auth/admin.json' });

function uniqueName(tag: string) {
  return `e2e-role-${tag}-${Date.now()}`;
}

type RoleType = 'app' | 'group';

// NOTE: there is no delete-role flow yet (no endpoint/UI), so created roles can't
// be cleaned up. Tests therefore use unique names and assert that the *specific*
// role they created appears — never absolute row/count totals — so they stay
// correct as roles accumulate on the shared DB. Add an afterEach cleanup here
// once a delete flow exists.
async function createRole(
  page: Page,
  opts: { name: string; type: RoleType; template: string },
) {
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

test.describe('roles', () => {
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
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

  test('edit a role updates its name, with the type locked', async ({ page }) => {
    const original = uniqueName('edit');
    await createRole(page, { name: original, type: 'app', template: 'Read Only' });

    await page.goto('/roles');
    const row = page.getByRole('row').filter({ hasText: original }).first();
    await expect(row).toBeVisible();
    // The row's edit action is the mat-icon "edit" inside the action cell.
    await row.locator('button:has(mat-icon:has-text("edit"))').click();

    await expect(page).toHaveURL(/\/roles\/\d+\/edit/);
    await expect(page.getByLabel('Role Name')).toHaveValue(original);
    // The role type can't be switched once the role exists.
    await expect(page.getByRole('button', { name: /Application role/ })).toBeDisabled();

    const renamed = `${original}-renamed`;
    await page.getByLabel('Role Name').fill(renamed);
    await page.getByRole('button', { name: 'Save Role' }).click();

    await expect(page).toHaveURL(/\/roles$/);
    await expect(
      page.getByRole('row').filter({ hasText: renamed }).first(),
    ).toBeVisible();
  });
});
