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
import { openReceiptsOverflowMenu } from './helpers/receipts-table';

// Bulk delete from the receipts table, where the selection can span groups (the
// All view) and the caller may delete in only some of them. The server judges
// each receipt in its own group and refuses the rest without refusing the
// others; the desktop names each refusal in a results dialog.
//
// Only an e2e proves this: the Jest specs mock the response, so they cannot
// show the server deleting one receipt while keeping its neighbour.
//
// e2e-user is a Legacy Viewer (no group.receipts.delete) in one provisioned
// group and a Legacy Owner in the other. Assertions run as e2e-user.
test.describe('Bulk receipt delete', () => {
  test.describe.configure({ mode: 'serial' });

  let adminContext: BrowserContext;
  let adminPage: Page;
  const viewerGroupName = uniqueName('bulk-del-viewer');
  const ownerGroupName = uniqueName('bulk-del-owner');
  const viewerReceiptName = uniqueName('bulk-del-kept');
  const ownerReceiptName = uniqueName('bulk-del-gone');
  let viewerGroupId: string;
  let ownerGroupId: string;
  let allGroupId: number;

  const row = (page: Page, name: string) =>
    page.getByRole('row').filter({ hasText: name });

  const selectRow = async (page: Page, name: string) => {
    await row(page, name).locator('mat-checkbox').click();
  };

  test.beforeAll(async ({ browser }) => {
    adminContext = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    adminPage = await adminContext.newPage();
    await stubTokenRefresh(adminPage);

    viewerGroupId = await createGroupWithMember(adminPage, {
      groupName: viewerGroupName,
      memberDisplayName: 'E2E User',
      roleName: 'Legacy Viewer',
    });
    ownerGroupId = await createGroupWithMember(adminPage, {
      groupName: ownerGroupName,
      memberDisplayName: 'E2E User',
      roleName: 'Legacy Owner',
    });

    await withAdminApi(async (api) => {
      const adminId = await apiGetUserId(api, creds('admin').username);
      await apiCreateReceipt(api, {
        groupId: viewerGroupId,
        paidByUserId: adminId,
        name: viewerReceiptName,
      });
      await apiCreateReceipt(api, {
        groupId: ownerGroupId,
        paidByUserId: adminId,
        name: ownerReceiptName,
      });
    });

    // The All group is synthetic and its id varies per install, so it is read
    // from the API's own marker, as e2e-user sees it.
    allGroupId = await withApiAs('user', async (api) => {
      const groups = (await (await api.get('/api/group/')).json()) as {
        id: number;
        isAllGroup?: boolean;
      }[];
      const all = groups.find((group) => group.isAllGroup);
      if (!all) {
        throw new Error('no All group visible to e2e-user');
      }
      return all.id;
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        await apiDeleteGroupById(api, viewerGroupId);
        await apiDeleteGroupById(api, ownerGroupId);
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
    await adminContext?.close();
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('is not offered for a selection the user may not delete', async ({
    page,
  }) => {
    await page.goto(`/receipts/group/${viewerGroupId}`);
    await expect(row(page, viewerReceiptName)).toBeVisible();
    await selectRow(page, viewerReceiptName);

    // Asserted by testid with the menu open: a menuitem is not a button, so a
    // role-based negative would pass whether or not the entry rendered.
    await openReceiptsOverflowMenu(page);
    await expect(page.getByTestId('receipts-export-selected')).toBeVisible();
    await expect(page.getByTestId('receipts-bulk-delete')).toHaveCount(0);
  });

  test('deletes what it may across groups and names each receipt it refused', async ({
    page,
  }) => {
    await page.goto(`/receipts/group/${allGroupId}`);
    await expect(row(page, viewerReceiptName)).toBeVisible();
    await expect(row(page, ownerReceiptName)).toBeVisible();

    await selectRow(page, viewerReceiptName);
    await selectRow(page, ownerReceiptName);

    await openReceiptsOverflowMenu(page);
    await page.getByTestId('receipts-bulk-delete').click();

    const confirmation = page
      .getByRole('dialog')
      .filter({ hasText: 'Delete Receipts' });
    await expect(confirmation).toContainText('delete 2 selected receipts');
    await confirmation.getByTestId('dialog-submit-button').click();

    const result = page
      .getByRole('dialog')
      .filter({ hasText: 'Some receipts could not be deleted' });
    await expect(result.getByTestId('bulk-delete-result-summary')).toHaveText(
      'Deleted 1 of 2 receipts.',
    );
    const failures = result.getByTestId('bulk-delete-result-failures');
    await expect(failures).toContainText(viewerReceiptName);
    await expect(failures).toContainText(
      'You do not have permission to delete this receipt',
    );
    await expect(failures).not.toContainText(ownerReceiptName);

    // The shared footer in dismiss-only mode: its cancel button is the only action.
    await expect(result.getByTestId('dialog-submit-button')).toHaveCount(0);
    // Found by role and name: the icon-only button must carry an accessible name.
    await result.getByRole('button', { name: 'Close' }).click();
    await expect(result).toBeHidden();

    // The table refetched: the permitted receipt is gone, its neighbour stays.
    await expect(row(page, ownerReceiptName)).toHaveCount(0);
    await expect(row(page, viewerReceiptName)).toBeVisible();
    // The selection was cleared, so nothing still counts the deleted receipt.
    await expect(page.locator('app-summary-card')).toHaveCount(0);

    // And the server agrees, independently of what the table shows.
    await withAdminApi(async (api) => {
      const res = await api.post(`/api/receipt/group/${ownerGroupId}`, {
        data: { page: 1, pageSize: 50 },
      });
      const body = (await res.json()) as { totalCount: number };
      expect(body.totalCount).toBe(0);
    });
  });
});
