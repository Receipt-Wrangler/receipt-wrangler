import { expect, Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateComment,
  apiCreateGroup,
  apiCreateReceipt,
  apiDeleteGroupById,
  apiGetUserId,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import { openReceiptsOverflowMenu } from './helpers/receipts-table';

test.use({ storageState: 'e2e/.auth/admin.json' });

// Serial: the tests share one seeded group. The column configuration does NOT
// carry between them - it lives in localStorage and every test gets a fresh
// context - so each test switches the column on itself.
test.describe.configure({ mode: 'serial' });

/**
 * The receipts table's Comment column: each receipt's FIRST comment, sortable.
 *
 * The Jest specs build columns against a mocked store, so they prove nothing
 * about the wire. What only an e2e can show: the paged list actually carrying
 * `firstComment`, the earliest comment winning over a later one, and the
 * `first_comment` sort key round-tripping to the API instead of being rejected.
 *
 * It seeds its own group so the rows are exactly these, in a known order.
 */
test.describe('Receipts table — Comment column', () => {
  const groupName = uniqueName('comment-col-grp');
  let groupId: number;

  // `later` is written second and sorts FIRST alphabetically, so a column that
  // showed - or sorted by - the lowest or the latest comment would be visibly
  // wrong on this receipt.
  const receipts = [
    { name: 'comment-col-bravo', comments: ['bravo note'] },
    { name: 'comment-col-zulu', comments: ['zulu note', 'aaa later note'] },
    { name: 'comment-col-alpha', comments: ['alpha note'] },
  ];

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      const adminId = await apiGetUserId(api, creds('admin').username);
      groupId = (await apiCreateGroup(api, groupName)).id;

      for (const receipt of receipts) {
        const receiptId = await apiCreateReceipt(api, {
          groupId,
          paidByUserId: adminId,
          name: receipt.name,
        });
        // Sequential on purpose: "first" is the earliest, and the id breaks a
        // created_at tie in insertion order.
        for (const comment of receipt.comments) {
          await apiCreateComment(api, { receiptId, comment });
        }
      }
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        if (groupId) {
          await apiDeleteGroupById(api, String(groupId));
        }
      });
    } catch {
      // Best-effort cleanup: never mask a real failure with a teardown error.
    }
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  async function gotoTable(page: Page): Promise<void> {
    await page.goto(`/receipts/group/${groupId}`);
    await expect(page.getByTestId('receipts-overflow-menu')).toBeVisible();
  }

  async function openConfigureColumns(page: Page) {
    await openReceiptsOverflowMenu(page);
    await page.getByTestId('configure-columns').click();
    const dialog = page.getByRole('dialog');
    await expect(dialog).toBeVisible();
    return dialog;
  }

  async function showCommentColumn(page: Page): Promise<void> {
    const dialog = await openConfigureColumns(page);
    await dialog.getByRole('checkbox', { name: 'Comment', exact: true }).check();
    await dialog.getByTestId('dialog-submit-button').click();
    await expect(dialog).toBeHidden();
    await expect(page.getByRole('columnheader', { name: 'Comment' })).toBeVisible();
  }

  /** The rendered row order, which is what every sort assertion compares. */
  async function receiptNames(page: Page): Promise<string[]> {
    const names = await page.locator('td a[href^="/receipts/"]').allTextContents();
    return names.map((name) => name.trim());
  }

  function rowFor(page: Page, receiptName: string) {
    return page.getByRole('row').filter({ hasText: receiptName });
  }

  test('is offered in Configure Columns, unchecked by default', async ({ page }) => {
    await gotoTable(page);

    const dialog = await openConfigureColumns(page);
    const checkbox = dialog.getByRole('checkbox', { name: 'Comment', exact: true });
    await expect(checkbox).toBeVisible();
    await expect(checkbox).not.toBeChecked();

    // It is a built-in column, not a custom field.
    await expect(
      dialog.locator('.column-item').filter({ hasText: 'Comment' }).getByTestId('column-config-custom'),
    ).toHaveCount(0);

    await expect(page.getByRole('columnheader', { name: 'Comment' })).toHaveCount(0);
  });

  test("shows each receipt's first comment", async ({ page }) => {
    await gotoTable(page);
    await showCommentColumn(page);

    await expect(rowFor(page, 'comment-col-zulu').getByTestId('receipt-first-comment')).toHaveText(
      'zulu note',
    );
    await expect(rowFor(page, 'comment-col-alpha').getByTestId('receipt-first-comment')).toHaveText(
      'alpha note',
    );
    await expect(page.getByText('aaa later note')).toHaveCount(0);
  });

  test('sorts by the first comment in both directions', async ({ page }) => {
    await gotoTable(page);
    await showCommentColumn(page);

    const header = page.getByRole('columnheader', { name: 'Comment' });

    // A 500 here would mean the API rejected first_comment as an untrusted
    // orderBy. The zulu receipt sorts by "zulu note", never by its later "aaa".
    await header.click();
    await expect
      .poll(() => receiptNames(page))
      .toEqual(['comment-col-alpha', 'comment-col-bravo', 'comment-col-zulu']);

    await header.click();
    await expect
      .poll(() => receiptNames(page))
      .toEqual(['comment-col-zulu', 'comment-col-bravo', 'comment-col-alpha']);
  });
});
