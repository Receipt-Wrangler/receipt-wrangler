import { BrowserContext, expect, Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateReceipt,
  apiDeleteGroupById,
  apiDeleteRoleByName,
  apiGetUserId,
  createGroupWithMember,
  createRole,
  uniqueName,
  withAdminApi,
  withApiAs,
} from './helpers/provisioning';

// A group member whose group role restricts paid-by visibility to "their own
// receipts" must not see — or be able to fetch — receipts paid by anyone else,
// while still seeing their own. An admin context provisions the role (Viewer
// preset, so the member holds group.receipts.read — the denial is paid-by, not a
// missing permission) + a group with e2e-user as that member + two receipts (one
// paid by admin, one by e2e-user). Assertions then run as e2e-user.
//
// The API assertions are the core guarantee (a hidden receipt GET must 403); the
// UI assertions prove the list hides it and the receipt-view guard redirects
// cleanly instead of admitting the user to a 403'd fetch.

test.describe('Group role paid-by visibility', () => {
  test.describe.configure({ mode: 'serial' });

  let adminContext: BrowserContext;
  let adminPage: Page;

  const roleName = uniqueName('paidby-role');
  const groupName = uniqueName('paidby-grp');
  const hiddenReceiptName = uniqueName('paidby-hidden');
  const ownReceiptName = uniqueName('paidby-own');

  let groupId: string;
  let hiddenReceiptId: number; // paid by admin — outside the member's grant
  let ownReceiptId: number; // paid by e2e-user — their own

  test.beforeAll(async ({ browser }) => {
    adminContext = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    adminPage = await adminContext.newPage();
    await stubTokenRefresh(adminPage);

    // Group role restricted to the member's own receipts on the paid-by axis.
    await createRole(adminPage, {
      name: roleName,
      type: 'Group role',
      preset: 'Viewer',
      paidByOwn: true,
    });

    // Group with e2e-user (display name "E2E User") holding that role.
    groupId = await createGroupWithMember(adminPage, {
      groupName,
      memberDisplayName: 'E2E User',
      roleName,
    });

    // Seed a receipt paid by admin (hidden) and one paid by e2e-user (visible).
    await withAdminApi(async (api) => {
      const adminId = await apiGetUserId(api, creds('admin').username);
      const userId = await apiGetUserId(api, creds('user').username);
      hiddenReceiptId = await apiCreateReceipt(api, {
        groupId,
        paidByUserId: adminId,
        name: hiddenReceiptName,
      });
      ownReceiptId = await apiCreateReceipt(api, {
        groupId,
        paidByUserId: userId,
        name: ownReceiptName,
      });
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        // Group delete frees the member's role assignment; then the role deletes.
        await apiDeleteGroupById(api, groupId);
        await apiDeleteRoleByName(api, roleName, 'GROUP');
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
    await adminContext?.close();
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('denies fetching a hidden-payer receipt but allows their own (API)', async () => {
    await withApiAs('user', async (api) => {
      const hidden = await api.get(`/api/receipt/${hiddenReceiptId}`);
      expect(hidden.status()).toBe(403);

      const own = await api.get(`/api/receipt/${ownReceiptId}`);
      expect(own.status()).toBe(200);
    });
  });

  test('receipts list shows their own receipt and hides the other-payer one', async ({
    page,
  }) => {
    await page.goto(`/receipts/group/${groupId}`);
    // Assert the own receipt first so the list has loaded before the negative check.
    await expect(page.getByText(ownReceiptName)).toBeVisible();
    await expect(page.getByText(hiddenReceiptName)).toHaveCount(0);
  });

  test('navigating directly to a hidden receipt redirects away', async ({
    page,
  }) => {
    await page.goto(`/receipts/${hiddenReceiptId}/view`);
    // The receipt-route guard's hasAccess probe is now paid-by-aware, so it
    // redirects instead of admitting the member to a receipt that 403s on fetch.
    await page.waitForURL(
      (url) => !url.href.includes(`/receipts/${hiddenReceiptId}/view`),
      { timeout: 10_000 },
    );
  });
});
