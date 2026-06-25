import { BrowserContext, expect, Page, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiDeleteGroupById,
  apiDeleteRoleByName,
  createGroupWithMember,
  createRole,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';

// The group-dashboards header Add/Edit/Delete buttons are gated with
// *hasGroupPermission on group.dashboards.create/update/delete
// (group-dashboards.component.html). A member holding group.dashboards.read (so
// they can VIEW dashboards) but lacking create/update/delete must not see the
// Add Dashboard control, while an owner of their own group does.
//
// The "Viewer" group preset grants group.dashboards.read but NOT
// create/update/delete (see role-presets.ts VIEWER_KEYS), so it is the exact
// fixture: an admin context provisions that role + a group with e2e-user as the
// member, then asserts as e2e-user. The admin who creates the group is its Owner
// (full perms incl. group.dashboards.create), giving the allow contrast.

test.describe('Group dashboard CRUD gating', () => {
  test.describe.configure({ mode: 'serial' });

  let adminContext: BrowserContext;
  let adminPage: Page;
  let roleName: string;
  let groupName: string;
  let groupId: string;

  test.beforeAll(async ({ browser }) => {
    roleName = uniqueName('dash-crud-role');
    groupName = uniqueName('dash-crud-grp');

    adminContext = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    adminPage = await adminContext.newPage();
    await stubTokenRefresh(adminPage);

    // Viewer holds group.dashboards.read but not create/update/delete.
    await createRole(adminPage, {
      name: roleName,
      type: 'Group role',
      preset: 'Viewer',
    });

    groupId = await createGroupWithMember(adminPage, {
      groupName,
      memberDisplayName: 'E2E User',
      roleName,
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        // Delete the group first (frees the member's group-role assignment),
        // then the now-unassigned role.
        await apiDeleteGroupById(api, groupId);
        await apiDeleteRoleByName(api, roleName, 'GROUP');
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
    await adminContext?.close();
  });

  // Runs as the default project user (e2e-user = the read-only viewer).
  test('a viewer holding only dashboards.read sees no Add Dashboard button', async ({
    page,
  }) => {
    await stubTokenRefresh(page);
    await page.goto(`/dashboard/group/${groupId}`);
    // The dashboards header heading renders for a member who holds .read.
    await expect(
      page.getByRole('heading', { name: /Dashboards$/ }),
    ).toBeVisible();
    // Add/Edit/Delete are gated on create/update/delete which the Viewer lacks.
    await expect(page.locator('app-add-button')).toHaveCount(0);
    await expect(page.locator('app-edit-button')).toHaveCount(0);
    await expect(page.locator('app-delete-button')).toHaveCount(0);
  });

  test.describe('owner with dashboards.create', () => {
    // The admin created the group, so they are its Owner and hold the permission.
    test.use({ storageState: 'e2e/.auth/admin.json' });

    test('sees the Add Dashboard button (contrast)', async ({ page }) => {
      await stubTokenRefresh(page);
      await page.goto(`/dashboard/group/${groupId}`);
      await expect(
        page.getByRole('heading', { name: /Dashboards$/ }),
      ).toBeVisible();
      await expect(page.locator('app-add-button')).toBeVisible();
    });
  });
});
