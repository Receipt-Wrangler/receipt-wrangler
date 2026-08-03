import { BrowserContext, expect, Page, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiDeleteRoleByName,
  apiDeleteUserByName,
  apiGroupNames,
  createRole,
  createUserWithRole,
  uniqueName,
  withAdminApi,
  withApiAsCreds,
} from './helpers/provisioning';

// An app role can opt its new users out of the automatic personal "My Receipts"
// group, for accounts that only ever belong to a few specific shared groups. The
// virtual "All" group is always created, so the account keeps a working
// dashboard. See api/CLAUDE.md → "Skipping the personal group per app role".
//
// Three layers are covered: the app-scope-only card in the role editor (present
// on APP, absent on GROUP, round-trips both ways through save/reload), the
// end-to-end effect on a real provisioned account, and the server-side scope
// guard that rejects the flag on a group role.

// Role management is admin-only, so override the default regular-user state.
test.use({ storageState: 'e2e/.auth/admin.json' });

const PERSONAL_GROUP = 'My Receipts';
const ALL_GROUP = 'All';
const PASSWORD = 'a really secure password';

// Each scope-specific card is located by its own <h2>, not by the shared
// .rw-card styling class — a raw CSS selector would break silently on a styling
// refactor even though the card is unchanged (see CLAUDE.md → "Locators").
const groupCreationCard = (page: Page) =>
  page.getByRole('heading', { name: 'Group creation' });

const memberVisibilityCard = (page: Page) =>
  page.getByRole('heading', { name: 'Member visibility' });

const skipCheckbox = (page: Page) =>
  page.getByTestId('skip-default-group').getByRole('checkbox');

// Opens a role's editor from the list. The edit icon button has no accessible
// name (mat-icon is aria-hidden and matTooltip sets none), so it has a testid.
async function openRoleEditor(page: Page, name: string) {
  await page.goto('/roles');
  const row = page.getByRole('row').filter({ hasText: name }).first();
  await expect(row).toBeVisible();
  await row.getByTestId('role-edit').click();
  await expect(page).toHaveURL(/\/roles\/\d+\/edit/);
  // The form populates from getRoles(); wait for the name before asserting more.
  await expect(page.getByLabel('Role Name')).toHaveValue(name);
}

// The submit button sits in a fixed bottom bar under a matTooltip overlay, which
// makes a direct click flaky on tall forms. Submit through the form's implicit
// Enter handler instead — a real user action that sidesteps the overlay.
async function saveRole(page: Page) {
  const name = page.getByLabel('Role Name');
  await expect(name).toBeEnabled();
  await name.press('Enter');
}

test.describe('role editor: the Group creation card is app-scope only', () => {
  const createdRoles: string[] = [];

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        for (const name of createdRoles) {
          await apiDeleteRoleByName(api, name, 'APP');
        }
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
  });

  test('renders for an application role and is unchecked by default', async ({
    page,
  }) => {
    await page.goto('/roles/new');
    await expect(page.getByLabel('Role Name')).toBeVisible();
    await page.getByRole('button', { name: /Application role/ }).click();

    await expect(groupCreationCard(page)).toBeVisible();
    await expect(
      page.getByText("Don't create a personal group for new users with this role"),
    ).toBeVisible();
    await expect(skipCheckbox(page)).not.toBeChecked();
  });

  test('is hidden for a group role, which shows the group-only card instead', async ({
    page,
  }) => {
    await page.goto('/roles/new');
    await expect(page.getByLabel('Role Name')).toBeVisible();
    await page.getByRole('button', { name: /Group role/ }).click();

    // Personal-group creation is meaningless on a group role (it is assigned to
    // a membership long after the account was created).
    await expect(groupCreationCard(page)).toHaveCount(0);
    // The group-scoped card still renders, proving the form isn't just blank.
    await expect(memberVisibilityCard(page)).toBeVisible();
  });

  test('resets when switching role type', async ({ page }) => {
    await page.goto('/roles/new');
    await expect(page.getByLabel('Role Name')).toBeVisible();
    await page.getByRole('button', { name: /Application role/ }).click();
    await skipCheckbox(page).check();
    await expect(skipCheckbox(page)).toBeChecked();

    // Switching scope clears every scope-specific setting; switching back must
    // not resurrect the old value.
    await page.getByRole('button', { name: /Group role/ }).click();
    await expect(groupCreationCard(page)).toHaveCount(0);
    await page.getByRole('button', { name: /Application role/ }).click();
    await expect(skipCheckbox(page)).not.toBeChecked();
  });

  test('round-trips through save, and can be turned back off', async ({ page }) => {
    const name = uniqueName('skip-roundtrip');
    createdRoles.push(name);

    await createRole(page, {
      name,
      type: 'Application role',
      preset: 'Read Only',
      skipDefaultGroup: true,
    });

    // Reopening reads the persisted value back off the API.
    await openRoleEditor(page, name);
    await expect(skipCheckbox(page)).toBeChecked();

    // Turning it back off must stick. This is the regression guard for the
    // update path: GORM's struct-form Updates skips zero-value bools, so a
    // toggled-off flag would silently stay on.
    await skipCheckbox(page).uncheck();
    await saveRole(page);
    await expect(page).toHaveURL(/\/roles$/);

    await openRoleEditor(page, name);
    await expect(skipCheckbox(page)).not.toBeChecked();
  });
});

test.describe('a flagged role skips the new user\'s personal group', () => {
  // One admin context provisions both roles and both users, then the assertions
  // run against them in order.
  test.describe.configure({ mode: 'serial' });

  let adminContext: BrowserContext;
  let adminPage: Page;
  let skipRole: string;
  let normalRole: string;
  let skipUser: string;
  let normalUser: string;
  // A third pair, used only by the creation-time-only case, which flips its
  // role's flag mid-test — so it must not share a role with the cases above.
  let laterRole: string;
  let laterUser: string;

  test.beforeAll(async ({ browser }) => {
    skipRole = uniqueName('skip-role');
    normalRole = uniqueName('normal-role');
    skipUser = uniqueName('skip-user');
    normalUser = uniqueName('normal-user');
    laterRole = uniqueName('later-role');
    laterUser = uniqueName('later-user');

    adminContext = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    adminPage = await adminContext.newPage();
    await stubTokenRefresh(adminPage);

    // "Read Only" grants every *.read permission, which includes
    // app.account.read — what GET /api/group/ is gated on, so each provisioned
    // user can list their own groups below.
    await createRole(adminPage, {
      name: skipRole,
      type: 'Application role',
      preset: 'Read Only',
      skipDefaultGroup: true,
    });
    await createRole(adminPage, {
      name: normalRole,
      type: 'Application role',
      preset: 'Read Only',
    });
    // Starts unflagged; the creation-time-only case turns it on afterwards.
    await createRole(adminPage, {
      name: laterRole,
      type: 'Application role',
      preset: 'Read Only',
    });

    await createUserWithRole(adminPage, {
      username: skipUser,
      password: PASSWORD,
      role: skipRole,
    });
    await createUserWithRole(adminPage, {
      username: normalUser,
      password: PASSWORD,
      role: normalRole,
    });
    await createUserWithRole(adminPage, {
      username: laterUser,
      password: PASSWORD,
      role: laterRole,
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        // Users first — a role cannot be deleted while it is still assigned.
        await apiDeleteUserByName(api, skipUser);
        await apiDeleteUserByName(api, normalUser);
        await apiDeleteUserByName(api, laterUser);
        await apiDeleteRoleByName(api, skipRole, 'APP');
        await apiDeleteRoleByName(api, normalRole, 'APP');
        await apiDeleteRoleByName(api, laterRole, 'APP');
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
    await adminContext?.close();
  });

  test('the flagged user gets the All group but no personal group', async () => {
    const groups = await withApiAsCreds(skipUser, PASSWORD, apiGroupNames);

    expect(groups).not.toContain(PERSONAL_GROUP);
    // The virtual "All" group is always created, so the account still has a
    // dashboard to land on.
    expect(groups).toContain(ALL_GROUP);
  });

  test('a user on an unflagged role still gets the personal group', async () => {
    const groups = await withApiAsCreds(normalUser, PASSWORD, apiGroupNames);

    // Contrast case: proves the skip is driven by the role's flag, not by
    // something else about how these accounts are provisioned.
    expect(groups).toContain(PERSONAL_GROUP);
    expect(groups).toContain(ALL_GROUP);
  });

  test('turning the flag on later leaves an existing user\'s groups alone', async () => {
    // laterUser was created while its role was unflagged, so it has the group.
    let groups = await withApiAsCreds(laterUser, PASSWORD, apiGroupNames);
    expect(groups).toContain(PERSONAL_GROUP);

    // Flip the flag on after the fact.
    await openRoleEditor(adminPage, laterRole);
    await skipCheckbox(adminPage).check();
    await saveRole(adminPage);
    await expect(adminPage).toHaveURL(/\/roles$/);
    await openRoleEditor(adminPage, laterRole);
    await expect(skipCheckbox(adminPage)).toBeChecked();

    // The flag is evaluated once, at user-creation time — it is not a live
    // property of the account, so the existing user keeps its personal group.
    groups = await withApiAsCreds(laterUser, PASSWORD, apiGroupNames);
    expect(groups).toContain(PERSONAL_GROUP);
    expect(groups).toContain(ALL_GROUP);
  });
});

test.describe('the flag is app-scope only on the server', () => {
  test('POST /api/role rejects it on a group role', async () => {
    await withAdminApi(async (api) => {
      const res = await api.post('/api/role', {
        data: {
          name: uniqueName('bad-group-role'),
          description: '',
          scope: 'GROUP',
          permissions: ['group.receipts.read'],
          skipDefaultGroupCreation: true,
        },
      });

      // The UI never sends it on a group role, but the server is the enforcer:
      // a crafted request must be rejected rather than silently ignored.
      expect(res.status()).toBe(400);
      expect(await res.text()).toContain('skipDefaultGroupCreation');
    });
  });
});
