import { expect, test, type Page } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';

function uniqueName(tag: string) {
  return `e2e-${tag}-${Date.now()}`;
}

// The role-list table row whose text contains `name`. Role names are unique, so
// the first match is the role's row.
function roleRow(page: Page, name: string) {
  return page.getByRole('row').filter({ hasText: name }).first();
}

// The `app-select` trigger (a mat-select) addressed by its label. Use the
// combobox role rather than getByLabel: while the option panel is open (and
// during its close animation) the listbox panel shares the same aria label, so
// getByLabel would match two elements.
function defaultSelect(page: Page, label: string) {
  return page.getByRole('combobox', { name: label });
}

// Opens a default-role selector and picks an option by its exact text.
async function selectDefault(page: Page, label: string, optionName: string) {
  await defaultSelect(page, label).click();
  await page.getByRole('option', { name: optionName, exact: true }).click();
}

test.describe('default roles (admin)', () => {
  test.use({ storageState: 'e2e/.auth/admin.json' });
  // The default per scope is global, mutable server state. Run serially and
  // restore after each test so the tests don't race each other.
  test.describe.configure({ mode: 'serial' });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test.afterEach(async ({ page }) => {
    // Restore the seeded default app role even if a test failed mid-way.
    // Re-selecting the already-active option is a no-op (mat-select emits no
    // change), so this is safe whether or not the test already restored it.
    try {
      await page.goto('/roles');
      await selectDefault(page, 'Default application role', 'Legacy User');
    } catch {
      // Best-effort — don't mask the test's own failure.
    }
  });

  test('seeded defaults are preselected and badged', async ({ page }) => {
    await page.goto('/roles');

    // The two selectors reflect the legacy-equivalent defaults seeded on boot.
    await expect(defaultSelect(page, 'Default application role')).toContainText('Legacy User');
    await expect(defaultSelect(page, 'Default group role')).toContainText('Legacy Owner');

    // The same roles carry a "Default" badge in the table.
    await expect(roleRow(page, 'Legacy User')).toContainText('Default');
    await expect(roleRow(page, 'Legacy Owner')).toContainText('Default');
  });

  test('changing the default application role persists and moves the badge', async ({ page }) => {
    await page.goto('/roles');

    await selectDefault(page, 'Default application role', 'Legacy Admin');

    await expect(page.getByText('Default role updated')).toBeVisible();
    await expect(defaultSelect(page, 'Default application role')).toContainText('Legacy Admin');
    await expect(roleRow(page, 'Legacy Admin')).toContainText('Default');
    await expect(roleRow(page, 'Legacy User')).not.toContainText('Default');

    // The change persists across a reload (it was written server-side).
    await page.goto('/roles');
    await expect(defaultSelect(page, 'Default application role')).toContainText('Legacy Admin');
    await expect(roleRow(page, 'Legacy Admin')).toContainText('Default');
    await expect(roleRow(page, 'Legacy User')).not.toContainText('Default');

    // Restore so the group-creation test (and reruns) see the seeded default.
    await selectDefault(page, 'Default application role', 'Legacy User');
    await expect(defaultSelect(page, 'Default application role')).toContainText('Legacy User');
    await expect(roleRow(page, 'Legacy User')).toContainText('Default');
  });
});

test.describe('group creation assigns the default group role (user)', () => {
  test.use({ storageState: 'e2e/.auth/user.json' });

  // Set after a group is created so afterEach can reap it if the test fails
  // before its own cleanup runs.
  let createdGroupName: string | null = null;

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test.afterEach(async ({ page }) => {
    if (!createdGroupName) return;
    const leftover = createdGroupName;
    createdGroupName = null;
    try {
      await deleteGroupByName(page, leftover);
    } catch {
      // Best-effort — the test already failed and we don't want to mask it.
    }
  });

  test('a regular user can create a group and open it', async ({ page }) => {
    const name = uniqueName('grp');

    await page.goto('/groups/create');
    await expect(page.getByLabel('Group Name')).toBeVisible();
    await page.getByLabel('Group Name').fill(name);
    // The app-form submits on Enter from a field (the submit button is an
    // icon-only control with no accessible name).
    await page.getByLabel('Group Name').press('Enter');

    // Landing on the group's detail view means the modern, permission-gated
    // group resolver succeeded — i.e. the creator was assigned the default
    // group role and is not locked out of the group they just created.
    await expect(page).toHaveURL(/\/groups\/\d+\/details\/view/);
    createdGroupName = name;
    await expect(page.getByLabel('Group Name')).toHaveValue(name);

    await deleteGroupByName(page, name);
    createdGroupName = null;
  });
});

// Deletes a group from the list by name. The delete button is an icon-only
// control (no accessible name), so it carries a data-testid; the confirmation
// dialog's confirm button likewise carries a data-testid.
async function deleteGroupByName(page: Page, name: string) {
  await page.goto('/groups');
  const row = page.getByRole('row').filter({ hasText: name }).first();
  await expect(row).toBeVisible();
  await row.getByTestId('group-delete').click();

  const dialog = page.getByRole('dialog');
  await expect(dialog).toBeVisible();
  await dialog.getByTestId('dialog-submit-button').click();

  await expect(page.getByRole('row').filter({ hasText: name })).toHaveCount(0);
}
