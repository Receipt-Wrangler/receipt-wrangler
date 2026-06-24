import { expect, test } from '@playwright/test';
import { rmSync } from 'node:fs';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiDeleteRoleByName,
  apiDeleteUserByName,
  createRole,
  createUserWithRole,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';

// The header search bar is gated on app.receipts.search
// (`*hasAppPermission="Permission.AppReceiptsSearch"` in header.component.html).
// legacy-user-visibility.spec.ts covers the POSITIVE side (Legacy User holds the
// permission -> the bar renders). This is the negative side: a user WITHOUT it
// never sees the bar.
//
// No seeded account lacks app.receipts.search, so an admin context provisions a
// custom app role omitting the Receipts resource entirely (its only app
// permission is app.receipts.search) and assigns it to a fresh user. The role
// still grants Account/Notifications/Groups so the user can load the app.

const NEW_USER_PASSWORD = 'a-really-secure-password';

// Session for the provisioned no-search user, written in beforeAll under the
// git-ignored .auth/ dir and removed in afterAll.
const PARTIAL_AUTH_FILE = 'e2e/.auth/no-search-user.json';

test.describe('Header search bar gating (no app.receipts.search)', () => {
  test.use({ storageState: PARTIAL_AUTH_FILE });
  test.describe.configure({ mode: 'serial' });

  let roleName: string;
  let username: string;

  test.beforeAll(async ({ browser }) => {
    roleName = uniqueName('no-search-role');
    username = uniqueName('no-search-user');

    const admin = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const adminPage = await admin.newPage();
    await stubTokenRefresh(adminPage);
    await createRole(adminPage, {
      name: roleName,
      type: 'Application role',
      preset: 'Start from scratch',
      // Baseline to load the app; deliberately no Receipts -> no app.receipts.search.
      enableCategories: ['Account', 'Notifications', 'Groups'],
    });
    await createUserWithRole(adminPage, {
      username,
      password: NEW_USER_PASSWORD,
      role: roleName,
    });
    await admin.close();

    // Log in once as that user and persist the session for the test below.
    // storageState: undefined opts out of the describe's (not-yet-created) file.
    const userContext = await browser.newContext({ storageState: undefined });
    const userPage = await userContext.newPage();
    await userPage.goto('/auth/login');
    await userPage.getByLabel('Username').fill(username);
    await userPage.getByLabel('Password').fill(NEW_USER_PASSWORD);
    await userPage.getByRole('button', { name: 'Login' }).click();
    await expect(userPage).toHaveURL(/\/dashboard\/group\/\d+/, {
      timeout: 15_000,
    });
    // Wait until permissions are persisted to the "auth" localStorage slice so
    // the saved session hydrates them synchronously on the next load.
    await userPage.waitForFunction(() =>
      (localStorage.getItem('auth') ?? '').includes('app.account.read'),
    );
    await userContext.storageState({ path: PARTIAL_AUTH_FILE });
    await userContext.close();
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        // Delete the user first (hard delete frees the role assignment), then
        // the now-unassigned role.
        await apiDeleteUserByName(api, username);
        await apiDeleteRoleByName(api, roleName, 'APP');
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure.
    }
    rmSync(PARTIAL_AUTH_FILE, { force: true });
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('the search bar never renders without app.receipts.search', async ({
    page,
  }) => {
    // Same page the positive test uses (the user holds app.groups.create).
    await page.goto('/groups/create');
    await expect(page.getByLabel('Group Name')).toBeVisible();
    await expect(page.getByPlaceholder('Search receipts')).toHaveCount(0);
  });
});
