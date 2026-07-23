import { expect, test, type Page } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';

function uniqueName(tag: string) {
  return `e2e-${tag}-${Date.now()}`;
}

// A mat-select trigger addressed by its field label via the combobox role. Use
// the combobox role rather than getByLabel: while the option panel is open the
// listbox shares the field's aria label, so getByLabel would match two elements.
function selectByLabel(page: Page, label: string) {
  return page.getByRole('combobox', { name: label });
}

// User creation and group-member management now assign modern roles (an app role
// for users, a group role for members) instead of the legacy enums, and each
// selector has a Preview button that opens the shared role-preview dialog.
test.describe('assigning modern roles (admin)', () => {
  test.use({ storageState: 'e2e/.auth/admin.json' });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('creates a user with a modern app role and previews it', async ({ page }) => {
    const username = uniqueName('approle');
    await page.goto('/users');

    await page.getByRole('button', { name: 'Add User' }).click();
    const dialog = page.getByRole('dialog').filter({ hasText: 'Create User' });
    await expect(dialog).toBeVisible();

    await dialog.getByLabel('Username').fill(username);
    await dialog.getByLabel('Displayname').fill('E2E App Role');
    await dialog.getByLabel('Password').fill('a-really-secure-password');

    // The app-role selector defaults to the configured default app role.
    const roleSelect = selectByLabel(page, 'App Role');
    await expect(roleSelect).toContainText('Legacy User');

    // Previewing opens a separate dialog describing the selected role.
    await dialog.getByTestId('role-preview').click();
    const preview = page.getByRole('dialog').filter({ hasText: 'Application role' });
    await expect(preview).toBeVisible();
    await expect(preview).toContainText('Legacy User');
    await preview.getByRole('button', { name: 'Close' }).click();
    await expect(preview).toBeHidden();

    // Assign a different app role, then save.
    await roleSelect.click();
    await page.getByRole('option', { name: 'Legacy Admin', exact: true }).click();
    await dialog.locator('app-submit-button button').click();

    // The dialog closes on success and the new user appears in the list.
    await expect(dialog).toBeHidden();
    const row = page.getByRole('row').filter({ hasText: username });
    await expect(row).toBeVisible();

    // Clean up the user created by this test.
    await row.getByTestId('user-delete').click();
    const confirm = page.getByRole('dialog');
    await expect(confirm).toBeVisible();
    await confirm.getByTestId('dialog-submit-button').click();
    await expect(page.getByRole('row').filter({ hasText: username })).toHaveCount(0);
  });

  test('previews a group role when adding a group member', async ({ page }) => {
    await page.goto('/groups/create');
    await expect(page.getByLabel('Group Name')).toBeVisible();

    await page.getByTestId('add-group-member').click();
    const dialog = page.getByRole('dialog').filter({ hasText: 'Add Group Member' });
    await expect(dialog).toBeVisible();

    // The role selector lists modern group roles, defaulting to the default group role.
    const roleSelect = selectByLabel(page, 'Role');
    await expect(roleSelect).toContainText('Legacy Owner');

    // Previewing opens a separate dialog describing the selected group role.
    await dialog.getByTestId('role-preview').click();
    const preview = page.getByRole('dialog').filter({ hasText: 'Group role' });
    await expect(preview).toBeVisible();
    await expect(preview).toContainText('Legacy Owner');
    await preview.getByRole('button', { name: 'Close' }).click();
    await expect(preview).toBeHidden();

    // Close the member dialog without saving — no group state is mutated.
    await page.keyboard.press('Escape');
    await expect(dialog).toBeHidden();
  });
});

// Deny-path contrast: a non-admin who lacks app.roles.read still reaches the
// group-create form (Legacy User holds app.groups.create), but the role selector
// degrades to EMPTY instead of erroring. The group-member-form loads roles with
// RoleService.getRoles() wrapped in catchError(() => of([])), so the background
// GET /role 403 is swallowed and groupRoleOptions() is empty — the dialog still
// opens and the page doesn't redirect/crash.
test.describe('role selector degrades for a non-admin (no app.roles.read)', () => {
  // Inherits the chromium project's e2e/.auth/user.json (e2e-user = Legacy User,
  // which lacks app.roles.read).
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('the Add Group Member role selector is empty rather than erroring', async ({
    page,
  }) => {
    await page.goto('/groups/create');
    // Legacy User holds app.groups.create, so the form loads (no redirect).
    await expect(page).toHaveURL(/\/groups\/create/);
    await expect(page.getByLabel('Group Name')).toBeVisible();

    await page.getByTestId('add-group-member').click();
    const dialog = page
      .getByRole('dialog')
      .filter({ hasText: 'Add Group Member' });
    await expect(dialog).toBeVisible();

    // Open the Role select: getRoles() 403'd and was caught, so it has no
    // options (the empty list is the graceful-degrade path, not an error).
    await dialog.getByRole('combobox', { name: 'Role' }).click();
    await expect(page.getByRole('option')).toHaveCount(0);

    // The dialog is still usable (no crash/redirect) — the Preview button is
    // simply disabled because no role can be selected (app-button binds the
    // disabled state onto its inner native button).
    await page.keyboard.press('Escape');
    await expect(
      dialog.getByTestId('role-preview').locator('button'),
    ).toBeDisabled();
  });
});
