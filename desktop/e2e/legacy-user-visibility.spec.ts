import { expect, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import { uniqueName, withAdminApi, withApiAs } from './helpers/provisioning';

// A "Legacy User" is the app role assigned to converted/normal users; its
// permission set is canonical and stable. These tests lock in that, post
// conversion, such a user cannot reach admin-only areas and, on shared resource
// pages, sees only the actions their permissions allow (create yes; update /
// delete no). Access is asserted via route guards (bulletproof, no UI fiddling)
// and via the resource-table action buttons (gated by *hasAppPermission, so they
// simply do not render for a Legacy User). The sidebar nav items use the same
// permission gates and route guards, so they are covered transitively here plus
// by the directive unit tests.

// Admin-only areas a Legacy User must be redirected away from. (`/users` is the
// case corrected in this change: the Legacy User role no longer holds
// app.users.read, so the route guard now denies it.)
const ADMIN_ONLY_ROUTES = ['/users', '/roles', '/system-settings'];

test.describe('Legacy User visibility', () => {
  // The default chromium project uses e2e/.auth/user.json (e2e-user = Legacy User).
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  for (const path of ADMIN_ONLY_ROUTES) {
    test(`is redirected to the dashboard from ${path}`, async ({ page }) => {
      await page.goto(path);
      await expect(page).toHaveURL(/\/dashboard\/group\/\d+/);
    });
  }

  test('categories: can add, cannot edit or delete', async ({ page }) => {
    await page.goto('/categories');
    // Legacy User holds app.categories.create.
    await expect(page.getByTestId('category-add')).toBeVisible();
    // ...but not app.categories.update / .delete — those row actions never render.
    await expect(page.getByTestId('category-edit')).toHaveCount(0);
    await expect(page.getByTestId('category-delete')).toHaveCount(0);
  });

  test('tags: can add, cannot edit or delete', async ({ page }) => {
    await page.goto('/tags');
    await expect(page.getByTestId('tag-add')).toBeVisible();
    await expect(page.getByTestId('tag-edit')).toHaveCount(0);
    await expect(page.getByTestId('tag-delete')).toHaveCount(0);
  });

  test('custom fields: can add, cannot delete', async ({ page }) => {
    await page.goto('/custom-fields');
    // Custom fields have no update permission — only create (held) and delete (not).
    await expect(page.getByTestId('custom-field-add')).toBeVisible();
    await expect(page.getByTestId('custom-field-delete')).toHaveCount(0);
  });

  test('groups: can reach the create form (holds app.groups.create)', async ({
    page,
  }) => {
    await page.goto('/groups/create');
    // Legacy User holds app.groups.create (but not app.groups.read, so the
    // /groups list is off-limits). The newly-added appPermissionGuard on the
    // create route must therefore admit them rather than redirect.
    await expect(page).toHaveURL(/\/groups\/create/);
    await expect(page.getByLabel('Group Name')).toBeVisible();
  });

  test('header search bar renders (holds app.receipts.search)', async ({
    page,
  }) => {
    await page.goto('/groups/create');
    // The search bar is now gated on app.receipts.search, which the Legacy User
    // holds — so the always-visible header search input still renders for them.
    await expect(page.getByPlaceholder('Search receipts')).toBeVisible();
  });
});

// Contrast: an admin (Legacy Admin) holds the corresponding permissions and can
// reach the same areas — proving the gating is selective, not globally broken.
test.describe('Admin can reach admin-only areas (contrast)', () => {
  test.use({ storageState: 'e2e/.auth/admin.json' });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('can open Manage Users', async ({ page }) => {
    await page.goto('/users');
    await expect(page).toHaveURL(/\/users/);
  });

  test('can open Manage Roles', async ({ page }) => {
    await page.goto('/roles');
    await expect(page).toHaveURL(/\/roles/);
  });

  test('can open System Settings', async ({ page }) => {
    await page.goto('/system-settings');
    // The landing guard redirects to the first tab the admin can read.
    await expect(page).toHaveURL(/\/system-settings\//);
  });
});

// The UI hides the category/tag edit/delete row actions for a Legacy User, but
// the server is the real enforcer: it must 403 a direct DELETE even if the UI is
// bypassed. An admin context seeds a category and a tag via the API; the Legacy
// User's DELETE on each must 403; admin teardown removes the survivors.
test.describe('Legacy User category/tag delete is denied (API 403)', () => {
  // Runs against the default project user (e2e-user = Legacy User) for the API
  // assertions; provisioning/teardown go through the admin API.
  test.describe.configure({ mode: 'serial' });

  let categoryId: number;
  let tagId: number;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      const catRes = await api.post('/api/category/', {
        data: { name: uniqueName('legacy-del-cat') },
      });
      expect(catRes.ok()).toBeTruthy();
      categoryId = ((await catRes.json()) as { id: number }).id;

      const tagRes = await api.post('/api/tag/', {
        data: { name: uniqueName('legacy-del-tag') },
      });
      expect(tagRes.ok()).toBeTruthy();
      tagId = ((await tagRes.json()) as { id: number }).id;
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        // The deletes above are denied, so both resources survive — remove them.
        await api.delete(`/api/category/${categoryId}`);
        await api.delete(`/api/tag/${tagId}`);
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
  });

  test('DELETE /api/category/:id 403s for a Legacy User', async () => {
    await withApiAs('user', async (api) => {
      const res = await api.delete(`/api/category/${categoryId}`);
      // Legacy User lacks app.categories.delete — the backend 403s.
      expect(res.status()).toBe(403);
    });
  });

  test('DELETE /api/tag/:id 403s for a Legacy User', async () => {
    await withApiAs('user', async (api) => {
      const res = await api.delete(`/api/tag/${tagId}`);
      // Legacy User lacks app.tags.delete — the backend 403s.
      expect(res.status()).toBe(403);
    });
  });
});
