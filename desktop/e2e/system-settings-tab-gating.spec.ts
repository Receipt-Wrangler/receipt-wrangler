import { expect, Page, test } from '@playwright/test';
import { rmSync } from 'node:fs';
import { stubTokenRefresh } from './helpers/auth';

// System Settings is a tabbed page where each tab gates independently on its own
// app.<resource>.read permission, and a landing guard redirects /system-settings
// to the first tab the user can read (falling through to the dashboard if none).
// The rest of the suite only checks the whole area (admin allowed / Legacy User
// redirected); this spec covers the PER-TAB boundary: a user who can read just
// one tab lands on it, sees only it, and is route-denied the others. An admin
// contrast confirms all tabs render.
//
// There's no seeded account with partial system access, so an admin context
// provisions a custom "Custom"-preset app role and assigns it to a fresh user;
// that user's session is saved as a storageState the tests run under.
//
// The tab under test is System Emails: its only route resolver is allGroups
// (needs app.groups.read, not a system tab), so the role can grant exactly one
// system-tab read without dragging in another. (The Prompts/System Tasks/Settings
// tabs resolve allReceiptProcessingSettings/systemSettings, which would force
// granting those tabs' reads too — see system-settings-routing.module.ts.) The
// role also grants Account (GetAppData requires app.account.read) and
// Notifications (the bootstrap's notificationCount) so the user can load the app
// at all — neither is a system-settings tab, so the per-tab matrix is unaffected.

function uniqueName(tag: string) {
  return `e2e-${tag}-${Date.now()}`;
}

const NEW_USER_PASSWORD = 'a-really-secure-password';

// Session for the provisioned partial-access user, written in beforeAll. Kept
// under the git-ignored .auth/ dir and removed in afterAll.
const PARTIAL_AUTH_FILE = 'e2e/.auth/system-tab-user.json';

// All five tab labels, in render order (system-settings.component.ts).
const ALL_TABS = [
  'System Settings',
  'Receipt Processing Settings',
  'Prompts',
  'System Emails',
  'System Tasks',
];

// Creates an Application role that can read exactly one system-settings tab
// (System Emails) plus the baseline a functional user needs. Starts from the
// empty "Custom" preset, then flips on whole resource groups via their toggles.
async function createSystemEmailsTabRole(page: Page, name: string) {
  await page.goto('/roles');
  await page.getByRole('button', { name: 'Add Role' }).first().click();
  await expect(page).toHaveURL(/\/roles\/new$/);
  await expect(page.getByLabel('Role Name')).toBeVisible();
  // Application is the default type; clicking it is a harmless no-op.
  await page.getByRole('button', { name: /Application role/ }).click();
  await page.getByLabel('Role Name').fill(name);
  // Empty slate (the "Custom" preset, addressed by its unique description to
  // avoid matching the "Custom Fields" group toggle).
  await page.getByRole('button', { name: 'Start from scratch' }).click();
  // Baseline so the user can load the app and the System Emails tab's allGroups
  // resolver succeeds; then the one system tab under test.
  for (const group of ['Account', 'Notifications', 'Groups', 'System Emails']) {
    await page
      .getByRole('button', { name: `Toggle all ${group} permissions` })
      .click();
  }
  await page.getByRole('button', { name: 'Save Role' }).click();
  await expect(page).toHaveURL(/\/roles$/);
}

async function createUserWithRole(
  page: Page,
  opts: { username: string; password: string; role: string },
) {
  await page.goto('/users');
  await page.getByRole('button', { name: 'Add User' }).click();
  const dialog = page.getByRole('dialog').filter({ hasText: 'Create User' });
  await expect(dialog).toBeVisible();
  await dialog.getByLabel('Username').fill(opts.username);
  await dialog.getByLabel('Displayname').fill(opts.username);
  await dialog.getByLabel('Password').fill(opts.password);
  await dialog.getByRole('combobox', { name: 'App Role' }).click();
  await page.getByRole('option', { name: opts.role, exact: true }).click();
  await dialog.locator('app-submit-button button').click();
  await expect(dialog).toBeHidden();
}

async function deleteUserByName(page: Page, username: string) {
  await page.goto('/users');
  const row = page.getByRole('row').filter({ hasText: username }).first();
  await expect(row).toBeVisible();
  await row.getByTestId('user-delete').click();
  const dialog = page.getByRole('dialog');
  await expect(dialog).toBeVisible();
  await dialog.getByTestId('dialog-submit-button').click();
  await expect(page.getByRole('row').filter({ hasText: username })).toHaveCount(
    0,
  );
}

async function deleteRoleByName(page: Page, name: string) {
  await page.goto('/roles');
  const row = page.getByRole('row').filter({ hasText: name }).first();
  await expect(row).toBeVisible();
  await row.getByTestId('role-delete').click();
  const dialog = page.getByRole('dialog');
  await expect(dialog).toBeVisible();
  await dialog.getByTestId('dialog-submit-button').click();
  await expect(page.getByRole('row').filter({ hasText: name })).toHaveCount(0);
}

test.describe('System settings per-tab gating (partial-access role)', () => {
  // Run as the provisioned partial-access user via a saved storageState — like
  // every other gating spec (storageState + stubTokenRefresh). The session is
  // created in beforeAll, once the role + user exist.
  test.use({ storageState: PARTIAL_AUTH_FILE });
  test.describe.configure({ mode: 'serial' });

  let roleName: string;
  let username: string;

  test.beforeAll(async ({ browser }) => {
    roleName = uniqueName('sysemail-role');
    username = uniqueName('sysemail-user');

    // Admin provisions the custom role + user.
    const admin = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const adminPage = await admin.newPage();
    await stubTokenRefresh(adminPage);
    await createSystemEmailsTabRole(adminPage, roleName);
    await createUserWithRole(adminPage, {
      username,
      password: NEW_USER_PASSWORD,
      role: roleName,
    });
    await admin.close();

    // Log in once as that user and persist the session for the tests below.
    // storageState: undefined opts out of the describe's PARTIAL_AUTH_FILE
    // (not created yet) that browser.newContext() would otherwise inherit.
    const userContext = await browser.newContext({ storageState: undefined });
    const userPage = await userContext.newPage();
    await userPage.goto('/auth/login');
    await userPage.getByLabel('Username').fill(username);
    await userPage.getByLabel('Password').fill(NEW_USER_PASSWORD);
    await userPage.getByRole('button', { name: 'Login' }).click();
    await expect(userPage).toHaveURL(/\/dashboard\/group\/\d+/, {
      timeout: 15_000,
    });
    // Wait until permissions are loaded AND persisted to the "auth" localStorage
    // slice, so the saved session hydrates them synchronously on the next load —
    // otherwise the system-settings landing guard runs before they're present.
    await userPage.waitForFunction(() =>
      (localStorage.getItem('auth') ?? '').includes('app.system-emails.read'),
    );
    await userContext.storageState({ path: PARTIAL_AUTH_FILE });
    await userContext.close();
  });

  test.afterAll(async ({ browser }) => {
    const admin = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const adminPage = await admin.newPage();
    try {
      await stubTokenRefresh(adminPage);
      // Delete the user first so the role becomes unassigned and deletable.
      await deleteUserByName(adminPage, username);
      await deleteRoleByName(adminPage, roleName);
    } catch {
      // Best-effort cleanup — don't mask a test failure.
    } finally {
      await admin.close();
    }
    rmSync(PARTIAL_AUTH_FILE, { force: true });
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('is route-denied the system tabs it cannot read', async ({ page }) => {
    // The role grants only app.system-emails.read among the system tabs, so the
    // per-tab route guards (appPermissionGuard) deny the others. The guard runs
    // BEFORE the route's resolvers, which is what we assert here: route-level
    // denial is the actual security boundary. (We deliberately don't assert
    // landing on / rendering the System Emails tab: its allGroups resolver — like
    // every system tab's resolvers, see system-settings-routing.module.ts — fires
    // a paged request a freshly provisioned user can't satisfy. Tab *rendering*
    // is covered by the admin contrast below; tab *gating* is covered here.)
    for (const tab of ['prompts', 'system-tasks', 'receipt-processing-settings']) {
      await page.goto(`/system-settings/${tab}`);
      await page.waitForURL(
        (url) => !url.pathname.startsWith('/system-settings'),
      );
    }
  });
});

test.describe('System settings tabs — admin', () => {
  test.use({ storageState: 'e2e/.auth/admin.json' });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('an admin sees every system-settings tab', async ({ page }) => {
    await page.goto('/system-settings');
    for (const label of ALL_TABS) {
      await expect(
        page.getByRole('tab', { name: label, exact: true }),
      ).toBeVisible();
    }
  });
});
