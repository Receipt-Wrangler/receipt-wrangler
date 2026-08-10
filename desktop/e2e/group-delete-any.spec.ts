import { expect, test, type Browser, type Page } from '@playwright/test';
import { rmSync } from 'node:fs';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiCreateGroup,
  apiDeleteGroupById,
  apiDeleteRoleByName,
  apiDeleteUserByName,
  createRole,
  createUserWithRole,
  uniqueName,
  withAdminApi,
  withApiAsCreds,
} from './helpers/provisioning';

// The group-delete row action is gated on the group-scoped `group.delete` OR the
// app-scoped `app.groups.delete` (`*hasGroupPermission` with `orApp`), mirroring
// the backend's `OrAppPermissions` on DeleteGroup. This is what lets an
// administrator who can already SEE every group (app.groups.read) clean up
// abandoned ones they are not a member of.
//
// No seeded account fits either side, so an admin context provisions two app
// roles through the real UI: a "deleter" holding the whole Groups resource
// (read + delete), and a "reader" identical but with Delete Any Group switched
// off. Each gets its own user + captured session; teardown is API-based.

// Generated at runtime (no static secret in the repo). The -Aa1! suffix
// guarantees upper/lower/digit/symbol in case a password policy is ever added.
const PASSWORD = `${uniqueName('pw')}-Aa1!`;

const DELETER_AUTH_FILE = 'e2e/.auth/group-deleter.json';
const READER_AUTH_FILE = 'e2e/.auth/group-reader.json';

/** Switches the groups list to the all-groups view via the filter dialog. */
async function showAllGroups(page: Page): Promise<void> {
  await page.goto('/groups');
  // Icon-only control: the tooltip is aria-describedby, not the a11y name.
  await page.getByTestId('group-filter').click();
  const dialog = page.getByRole('dialog').filter({ hasText: 'Filter Groups' });
  await expect(dialog).toBeVisible();
  await dialog.getByRole('combobox', { name: 'Groups to View' }).click();
  await page.getByRole('option', { name: 'All Groups', exact: true }).click();
  await dialog.getByTestId('dialog-submit-button').click();
  await expect(dialog).toBeHidden();
  await expect(page.getByText('All Groups', { exact: true })).toBeVisible();
}

/**
 * Logs in as [username] once and persists the session, waiting until the held
 * permission has landed in the "auth" localStorage slice so the saved state
 * hydrates permissions synchronously on the next load.
 */
async function captureSession(
  browser: Browser,
  username: string,
  heldPermission: string,
  path: string,
): Promise<void> {
  const context = await browser.newContext({ storageState: undefined });
  const page = await context.newPage();
  await page.goto('/auth/login');
  await page.getByLabel('Username').fill(username);
  await page.getByLabel('Password').fill(PASSWORD);
  await page.getByRole('button', { name: 'Login' }).click();
  await expect(page).toHaveURL(/\/dashboard\/group\/\d+/, { timeout: 15_000 });
  await page.waitForFunction(
    (permission) => (localStorage.getItem('auth') ?? '').includes(permission),
    heldPermission,
  );
  await context.storageState({ path });
  await context.close();
}

test.describe('Delete any group (app.groups.delete)', () => {
  test.describe.configure({ mode: 'serial' });

  let deleterRole: string;
  let readerRole: string;
  let deleterUser: string;
  let readerUser: string;
  // Two groups so the destructive test can't strand the read-only assertions.
  let deletableGroup: { id: number; name: string };
  let untouchedGroup: { id: number; name: string };

  test.beforeAll(async ({ browser }) => {
    deleterRole = uniqueName('group-deleter-role');
    readerRole = uniqueName('group-reader-role');
    deleterUser = uniqueName('group-deleter');
    readerUser = uniqueName('group-reader');

    const admin = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const adminPage = await admin.newPage();
    await stubTokenRefresh(adminPage);

    // The whole Groups resource: app.groups.create/read/update-settings/delete.
    await createRole(adminPage, {
      name: deleterRole,
      type: 'Application role',
      preset: 'Start from scratch',
      enableCategories: ['Account', 'Notifications', 'Groups'],
    });
    // Identical, minus the one permission under test.
    await createRole(adminPage, {
      name: readerRole,
      type: 'Application role',
      preset: 'Start from scratch',
      enableCategories: ['Account', 'Notifications', 'Groups'],
      disablePermissions: [{ panelKey: 'app.groups', label: 'Delete Any Group' }],
    });

    await createUserWithRole(adminPage, {
      username: deleterUser,
      password: PASSWORD,
      role: deleterRole,
    });
    await createUserWithRole(adminPage, {
      username: readerUser,
      password: PASSWORD,
      role: readerRole,
    });
    await admin.close();

    // Groups owned by the admin — neither provisioned user is a member, so the
    // group-scoped group.delete can never apply to them.
    await withAdminApi(async (api) => {
      deletableGroup = await apiCreateGroup(api, uniqueName('abandoned-group'));
      untouchedGroup = await apiCreateGroup(api, uniqueName('leftover-group'));
    });

    await captureSession(
      browser,
      deleterUser,
      'app.groups.delete',
      DELETER_AUTH_FILE,
    );
    await captureSession(
      browser,
      readerUser,
      'app.groups.read',
      READER_AUTH_FILE,
    );
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        // Groups/users first — deleting them frees the role assignments, without
        // which the roles can't be removed.
        for (const group of [deletableGroup, untouchedGroup]) {
          if (group) {
            await apiDeleteGroupById(api, String(group.id));
          }
        }
        await apiDeleteUserByName(api, deleterUser);
        await apiDeleteUserByName(api, readerUser);
        await apiDeleteRoleByName(api, deleterRole, 'APP');
        await apiDeleteRoleByName(api, readerRole, 'APP');
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
    rmSync(DELETER_AUTH_FILE, { force: true });
    rmSync(READER_AUTH_FILE, { force: true });
  });

  // The sessions are written in beforeAll, which runs after Playwright would
  // resolve a `test.use({ storageState })` option — so each test opens its own
  // context from the file instead of declaring it as a fixture.
  test('an app.groups.delete holder deletes a group it is not a member of', async ({
    browser,
  }) => {
    const context = await browser.newContext({
      storageState: DELETER_AUTH_FILE,
    });
    const page = await context.newPage();
    try {
      await stubTokenRefresh(page);
      await showAllGroups(page);

      const row = page
        .getByRole('row')
        .filter({ hasText: deletableGroup.name })
        .first();
      await expect(row).toBeVisible();

      const deleteButton = row.getByTestId('group-delete');
      await expect(deleteButton).toBeVisible();
      // The disabled binding must not fall back to the caller's own group count.
      await expect(deleteButton.locator('button')).toBeEnabled();

      await deleteButton.click();
      await page.getByTestId('dialog-submit-button').click();

      // The table refetches, so the row goes away while the view stays on the
      // all-groups filter (a collapse to "My Groups" would also drop the row).
      await expect(
        page.getByRole('row').filter({ hasText: deletableGroup.name }),
      ).toHaveCount(0);
      await expect(page.getByText('All Groups', { exact: true })).toBeVisible();
      await expect(
        page.getByRole('row').filter({ hasText: untouchedGroup.name }),
      ).toHaveCount(1);
    } finally {
      await context.close();
    }
  });

  test('a reader without app.groups.delete sees every group but no delete action', async ({
    browser,
  }) => {
    const context = await browser.newContext({ storageState: READER_AUTH_FILE });
    const page = await context.newPage();
    try {
      await stubTokenRefresh(page);
      await showAllGroups(page);

      const row = page
        .getByRole('row')
        .filter({ hasText: untouchedGroup.name })
        .first();
      // app.groups.read alone still lists it...
      await expect(row).toBeVisible();
      // ...but neither gate (group.delete / app.groups.delete) is held.
      await expect(row.getByTestId('group-delete')).toHaveCount(0);
    } finally {
      await context.close();
    }
  });

  test('a reader without app.groups.delete is denied by the server', async () => {
    // The hidden button is a UI hint; the endpoint is the real gate.
    await withApiAsCreds(readerUser, PASSWORD, async (api) => {
      const res = await api.delete(`/api/group/${untouchedGroup.id}`);
      expect(res.status()).toBe(403);
    });
  });
});
