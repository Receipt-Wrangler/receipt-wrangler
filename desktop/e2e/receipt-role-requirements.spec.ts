import { readFileSync } from 'node:fs';
import { BrowserContext, expect, Page, type Route, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiDeleteGroupById,
  apiDeleteRoleByName,
  createGroupWithMember,
  createRole,
  uniqueName,
  withAdminApi,
  withApiAs,
} from './helpers/provisioning';
import {
  RECEIPT_PNG,
  injectQuickScanAppData,
  openQuickScanDialog,
  selectImageGroup,
  uploadQuickScanImages,
} from './helpers/quick-scan';

// A group role can require its members to keep at least one comment and one
// image on the group's receipts. The desktop mirrors the server's check before
// submitting, and create is ONE multipart call (POST /receipt/withFiles)
// carrying the receipt, its comments and its images — the only way the server
// can see the images at create time.
//
// An admin provisions a role with both flags through the real role form, and a
// group with e2e-user as its member; every assertion then runs as e2e-user (the
// default project storageState — AppData is re-fetched on navigation, so the new
// membership and its requirements reach the store without a re-login).

test.describe('Role-required receipt comment and image', () => {
  test.describe.configure({ mode: 'serial' });

  let adminContext: BrowserContext;
  let adminPage: Page;

  const roleName = uniqueName('receipt-req-role');
  const groupName = uniqueName('receipt-req-grp');
  let groupId: string;

  // Created by the first create test, reused by the edit test.
  let receiptId: number;
  const commentText = `required-comment-${Date.now()}`;

  test.beforeAll(async ({ browser }) => {
    adminContext = await browser.newContext({ storageState: 'e2e/.auth/admin.json' });
    adminPage = await adminContext.newPage();
    await stubTokenRefresh(adminPage);

    // Group Manager covers create/update and comment create/delete; flipping the
    // whole Receipts resource on adds group.receipts.quick-scan for the dialog.
    await createRole(adminPage, {
      name: roleName,
      type: 'Group role',
      preset: 'Group Manager',
      enableCategories: ['Receipts'],
      requireReceiptComment: true,
      requireReceiptImage: true,
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
        // The group delete cascades its receipts and frees the role.
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

  test('the role form persists both requirement flags', async () => {
    const role = await withAdminApi(async (api) => {
      const roles = (await (await api.get('/api/role')).json()) as {
        name: string;
        requireReceiptComment?: boolean;
        requireReceiptImage?: boolean;
      }[];
      return roles.find((candidate) => candidate.name === roleName);
    });

    expect(role?.requireReceiptComment).toBe(true);
    expect(role?.requireReceiptImage).toBe(true);

    // And the edit form rehydrates them.
    await adminPage.goto('/roles');
    await adminPage.getByRole('row').filter({ hasText: roleName }).first().getByTestId('role-edit').click();
    await expect(adminPage).toHaveURL(/\/roles\/\d+\/edit\?scope=group/);
    await expect(
      adminPage.getByTestId('require-receipt-comment').getByRole('checkbox'),
    ).toBeChecked();
    await expect(adminPage.getByTestId('require-receipt-image').getByRole('checkbox')).toBeChecked();
  });

  test('create is blocked without a comment and an image, then succeeds in ONE request', async ({
    page,
  }) => {
    const name = uniqueName('req-receipt');
    await openAddFormInGroup(page, groupId, groupName);

    await page.getByLabel('Name').fill(name);
    await page.getByLabel('Amount').fill('18.00');
    await selectFirstOption(page, 'Paid By');

    // Both inline hints show while nothing is attached.
    await expect(page.getByTestId('receipt-image-required-hint')).toBeVisible();
    await expect(page.getByTestId('receipt-comment-required-hint')).toBeVisible();

    const writes = trackReceiptWrites(page);

    await clickSave(page);
    await expect(
      page.getByText("Your role requires an image and a comment on this group's receipts."),
    ).toBeVisible();
    await expect(page).toHaveURL(/\/receipts\/add/);

    // A comment alone is not enough.
    await page.getByLabel('Comment').fill(commentText);
    await page.getByRole('button', { name: 'Comment', exact: true }).click();
    await expect(page.getByTestId('receipt-comment-required-hint')).toHaveCount(0);

    await clickSave(page);
    await expect(
      page.getByText("Your role requires an image on this group's receipts."),
    ).toBeVisible();
    await expect(page).toHaveURL(/\/receipts\/add/);
    expect(writes).toEqual([]);

    // Queue an image; the save now goes through as a single multipart create.
    await page.locator('app-upload-image input[type="file"]').first().setInputFiles(RECEIPT_PNG);
    await expect(page.getByTestId('receipt-image-required-hint')).toHaveCount(0);

    await clickSave(page);
    await expect(page).toHaveURL(/\/receipts\/\d+\/view/);
    receiptId = Number(page.url().match(/\/receipts\/(\d+)\/view/)![1]);

    expect(writes).toEqual(['/api/receipt/withFiles']);

    // The receipt landed with both the comment and the image.
    const saved = await withApiAs('user', async (api) =>
      (await api.get(`/api/receipt/${receiptId}`)).json(),
    );
    expect(saved.imageFiles).toHaveLength(1);
    expect(saved.comments.map((comment: { comment: string }) => comment.comment)).toEqual([
      commentText,
    ]);
  });

  test('edit mode refuses deleting the last image and the last comment', async ({ page }) => {
    test.skip(!receiptId, 'depends on the create test');

    await page.goto(`/receipts/${receiptId}/edit`);
    await expect(page.getByLabel('Name')).toBeVisible();
    await expect(page.getByText(commentText)).toBeVisible();

    // The only image cannot be removed, and the only comment offers no delete.
    await expect(page.getByTestId('receipt-image-remove').locator('button')).toBeDisabled();
    await expect(page.getByTestId('comment-delete')).toHaveCount(0);

    // The server refuses both deletes too, so an older client cannot bypass it.
    await withApiAs('user', async (api) => {
      const saved = await (await api.get(`/api/receipt/${receiptId}`)).json();
      const imageDelete = await api.delete(`/api/receiptImage/${saved.imageFiles[0].id}`);
      expect(imageDelete.status()).toBe(400);
      const commentDelete = await api.delete(`/api/comment/${saved.comments[0].id}`);
      expect(commentDelete.status()).toBe(400);
    });

    // A second comment makes either one deletable again.
    await page.getByLabel('Comment').fill('a second comment');
    await page.getByRole('button', { name: 'Comment', exact: true }).click();
    await expect(page.getByText('a second comment')).toBeVisible();
    await expect(page.getByTestId('comment-delete')).toHaveCount(2);
  });

  test('quick scan requires a comment when the role does', async ({ page }) => {
    // The group's own quick-scan config leaves the comment OFF; the role's
    // requirement (from the real AppData) must switch it on and make it required.
    await injectQuickScanAppData(page, {
      groupConfigs: [
        {
          groupId: Number(groupId),
          config: {
            quickScanPaidByEnabled: false,
            quickScanStatusEnabled: false,
            quickScanDefaultStatus: 'OPEN',
            quickScanCategoriesEnabled: false,
            quickScanTagsEnabled: false,
            quickScanCommentEnabled: false,
            quickScanCommentRequired: false,
          },
        },
      ],
    });

    const requests: string[] = [];
    await page.route('**/api/receipt/quickScan', async (route: Route) => {
      requests.push(route.request().url());
      await route.fulfill({ status: 200, contentType: 'application/json', body: '{}' });
    });

    const dialog = await openQuickScanDialog(page, Number(groupId));
    await uploadQuickScanImages(dialog, 1);
    await selectImageGroup(page, dialog, groupName);

    const comment = dialog.getByTestId('quick-scan-comment').getByRole('textbox');
    await expect(comment).toBeVisible();

    await dialog.getByTestId('dialog-submit-button').click();
    await expect(page.getByText('Please fill in all required fields', { exact: false })).toBeVisible();
    expect(requests).toHaveLength(0);

    await comment.fill('Scanned with a note');
    await dialog.getByTestId('dialog-submit-button').click();
    await expect(page.getByText('Successfully queued', { exact: false })).toBeVisible();
    expect(requests).toHaveLength(1);

    // The server enforces it on its own: the same scan without a comment is a 400.
    const userId = await withApiAs('user', async (api) => {
      const appData = await (await api.get('/api/user/appData')).json();
      return appData.claims.userId as number;
    });
    const status = await withApiAs('user', async (api) => {
      const res = await api.post('/api/receipt/quickScan', {
        multipart: {
          files: {
            name: 'receipt.png',
            mimeType: 'image/png',
            buffer: readFileSync(RECEIPT_PNG),
          },
          groupIds: groupId,
          paidByUserIds: String(userId),
          statuses: 'OPEN',
          categoryIds: '',
          tagIds: '',
          comments: '',
        },
      });
      return res.status();
    });
    expect(status).toBe(400);
  });
});

/**
 * Opens the add form with [groupName] as its group. The form seeds the group
 * being browsed, but the member belongs to several groups, so the field is
 * re-picked unless it already holds this one.
 */
async function openAddFormInGroup(page: Page, groupId: string, groupName: string) {
  await page.goto(`/receipts/group/${groupId}`);
  await page.goto('/receipts/add');
  const groupField = page.getByRole('combobox', { name: 'Group' });
  await expect(page.getByLabel('Name')).toBeVisible();
  if ((await groupField.inputValue()) !== groupName) {
    const clear = page.getByTestId('receipt-group').getByTestId('autocomplete-clear');
    if (await clear.count()) {
      await clear.first().click();
    }
    await groupField.click();
    await groupField.fill(groupName);
    await page.getByRole('option', { name: groupName, exact: true }).click();
  }
  await expect(groupField).toHaveValue(groupName);
}

async function selectFirstOption(page: Page, label: string) {
  const field = page.getByRole('combobox', { name: label });
  if (await field.inputValue()) {
    return;
  }
  await field.click();
  await page.getByRole('option').first().click();
}

async function clickSave(page: Page) {
  await page.getByRole('button', { name: 'Save', exact: true }).first().click({ force: true });
}

/**
 * Records the path of every receipt/image WRITE the page sends — the old flow
 * was POST /receipt then one POST /receiptImage per image, the new one is a
 * single POST /receipt/withFiles.
 */
function trackReceiptWrites(page: Page): string[] {
  const writes: string[] = [];
  page.on('request', (request) => {
    if (request.method() !== 'POST') {
      return;
    }
    const path = new URL(request.url()).pathname.replace(/\/$/, '');
    if (/^\/api\/(receipt|receipt\/withFiles|receiptImage)$/.test(path)) {
      writes.push(path);
    }
  });
  return writes;
}
