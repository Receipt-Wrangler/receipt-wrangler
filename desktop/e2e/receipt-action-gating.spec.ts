import { BrowserContext, expect, Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateReceipt,
  apiDeleteGroupById,
  apiGetUserId,
  createGroupWithMember,
  uniqueName,
  withAdminApi,
  withApiAs,
} from './helpers/provisioning';

// Receipt write actions are gated by group permissions a Legacy Viewer lacks:
// the receipts-table row's Duplicate (group.receipts.duplicate) and Delete
// (group.receipts.delete) actions, the receipt edit route (group.receipts.update,
// via receiptGuardGuard), and creating a receipt (group.receipts.create, server-
// enforced). The receipts-table *edit* action itself is not template-gated, so
// the edit denial is asserted at the route level, plus an API-403 proves the
// backend enforces create independently of the hidden UI.
//
// A focused, self-contained spec: it provisions its own Legacy Viewer group +
// receipt and asserts as e2e-user. Kept separate from
// group-viewer-visibility.spec.ts (which covers GROUP-level gating) to keep each
// spec's concern narrow.
test.describe('Group Viewer is denied receipt write actions', () => {
  test.describe.configure({ mode: 'serial' });

  let adminContext: BrowserContext;
  let adminPage: Page;
  const groupName = uniqueName('rcpt-action-grp');
  const receiptName = uniqueName('rcpt-action');
  let groupId: string;
  let receiptId: number;
  let userId: number; // e2e-user's id (admin-only user list)

  test.beforeAll(async ({ browser }) => {
    adminContext = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    adminPage = await adminContext.newPage();
    await stubTokenRefresh(adminPage);

    // Group with e2e-user ("E2E User") as a Legacy Viewer (holds
    // group.receipts.read, none of update/delete/duplicate/create).
    groupId = await createGroupWithMember(adminPage, {
      groupName,
      memberDisplayName: 'E2E User',
      roleName: 'Legacy Viewer',
    });

    // Seed a receipt paid by admin so it appears in the Viewer's list (the
    // Viewer holds read, no paid-by restriction).
    await withAdminApi(async (api) => {
      const adminId = await apiGetUserId(api, creds('admin').username);
      userId = await apiGetUserId(api, creds('user').username);
      receiptId = await apiCreateReceipt(api, {
        groupId,
        paidByUserId: adminId,
        name: receiptName,
      });
    });
  });

  test.afterAll(async () => {
    try {
      // Deleting the group cascades its receipt and frees the member's role.
      await withAdminApi(async (api) => apiDeleteGroupById(api, groupId));
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
    await adminContext?.close();
  });

  // Assertions run as the default project user (e2e-user = the Viewer).
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('receipts row exposes no duplicate or delete action', async ({
    page,
  }) => {
    await page.goto(`/receipts/group/${groupId}`);
    // The Viewer holds group.receipts.read, so the seeded receipt is listed.
    const row = page.getByRole('row').filter({ hasText: receiptName }).first();
    await expect(row).toBeVisible();
    // Duplicate / delete row actions are gated on permissions not held.
    await expect(row.getByTestId('receipt-duplicate')).toHaveCount(0);
    await expect(row.getByTestId('receipt-delete')).toHaveCount(0);
  });

  test('navigating to the receipt edit route redirects away', async ({
    page,
  }) => {
    await page.goto(`/receipts/${receiptId}/edit`);
    // receiptGuardGuard checks group.receipts.update (not held) and redirects
    // rather than admitting a 403'd edit.
    await page.waitForURL(
      (url) => !url.href.includes(`/receipts/${receiptId}/edit`),
      { timeout: 10_000 },
    );
  });

  test('creating a receipt in the group is denied (API 403)', async () => {
    await withApiAs('user', async (api) => {
      const res = await api.post('/api/receipt', {
        data: {
          name: uniqueName('viewer-denied-rcpt'),
          amount: '10.00',
          date: '2024-01-01T00:00:00Z',
          groupId: Number(groupId),
          paidByUserId: userId,
          status: 'OPEN',
        },
      });
      // The Viewer lacks group.receipts.create — the backend 403s.
      expect(res.status()).toBe(403);
    });
  });
});
