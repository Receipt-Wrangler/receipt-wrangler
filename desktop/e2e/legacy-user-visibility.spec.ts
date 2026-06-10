import { expect, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';

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
