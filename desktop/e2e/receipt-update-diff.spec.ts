import { expect, test, type Page } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateGroup,
  apiCreateReceipt,
  apiDeleteGroupById,
  apiGetUserId,
  apiPagedSystemTasks,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';

// The System Tasks "Updated Receipt" row, end to end.
//
// The API stores a RECEIPT_UPDATED description double-encoded (each of
// before/after is itself a JSON string), which the generic pretty-json pipe
// could not parse, so the row showed the raw escaped text. The Jest specs build
// that format by hand; this proves the desktop reads what the real server
// writes, and that the before/after diff puts each side where it belongs.
//
// The table is shared with every other spec's tasks, so the paged response is
// narrowed to this spec's own task. The row itself is still the real server's.
//
// The receipt carries an item on both sides of the update under test. The
// server deletes and recreates items on every save, so the item comes back with
// a new id and timestamps; the summary must still list only what was edited.

test.use({ storageState: 'e2e/.auth/admin.json' });

test.describe('Receipt update diff', () => {
  test.describe.configure({ mode: 'serial' });

  const oldName = uniqueName('diff-before');
  const newName = uniqueName('diff-after');
  let groupId: number;
  let taskId: number;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      const paidByUserId = await apiGetUserId(api, creds('admin').username);
      groupId = (await apiCreateGroup(api, uniqueName('diff-group'))).id;
      const receiptId = await apiCreateReceipt(api, {
        groupId,
        paidByUserId,
        name: oldName,
        amount: '10.00',
      });

      const update = async (name: string, amount: string) => {
        const res = await api.put(`/api/receipt/${receiptId}`, {
          data: {
            name,
            amount,
            date: '2024-01-01T00:00:00Z',
            groupId,
            paidByUserId,
            status: 'OPEN',
            receiptItems: [
              {
                receiptId,
                name: 'Pizza',
                amount: '10.00',
                chargedToUserId: paidByUserId,
                status: 'OPEN',
              },
            ],
          },
        });
        expect(res.ok()).toBe(true);
      };

      // Adds the item; this spec does not look at that update's row.
      await update(oldName, '10.00');
      // The update under test: same item, new name and amount.
      await update(newName, '12.50');

      const tasks = await apiPagedSystemTasks(api, {
        type: { operation: 'CONTAINS', value: ['RECEIPT_UPDATED'] },
      });
      const ownTasks = tasks.data.filter((row) => row.associatedEntityId === receiptId);
      expect(ownTasks).toHaveLength(2);
      const task = ownTasks.reduce((latest, row) => (row.id > latest.id ? row : latest));
      taskId = task.id;

      // The current API marks what it writes, so the desktop trusts its own
      // "before" rather than rebuilding it.
      expect(JSON.parse(task.resultDescription).version).toBe(2);
    });
  });

  test.afterAll(async () => {
    await withAdminApi(async (api) => {
      try {
        await apiDeleteGroupById(api, String(groupId));
      } catch {
        // Best effort; a leaked group does not affect other specs.
      }
    });
  });

  const gotoOwnTask = async (page: Page): Promise<void> => {
    await stubTokenRefresh(page);
    await page.route('**/api/systemTask/getPagedSystemTasks', async (route) => {
      const response = await route.fetch();
      const body = await response.json();
      const data = (body.data ?? []).filter((row: { id: number }) => row.id === taskId);
      await route.fulfill({ response, json: { ...body, data, totalCount: data.length } });
    });
    await page.goto('/system-settings/system-tasks');
  };

  test('the row summarizes only what was edited, not the recreated item', async ({ page }) => {
    await gotoOwnTask(page);

    await expect(page.getByTestId('receipt-update-summary')).toHaveText('Changed: name, amount');
    await expect(page.getByText('\\"before\\"')).toHaveCount(0);
  });

  test('the dialog shows the old values on the left and the new ones on the right', async ({ page }) => {
    await gotoOwnTask(page);
    await page.getByTestId('receipt-update-diff-open').click();

    const dialog = page.getByTestId('receipt-diff-dialog');
    await expect(dialog.getByRole('heading', { name: `Receipt update: ${newName}` })).toBeVisible();

    const changed = dialog.locator('[data-testid="receipt-diff-row"][data-kind="changed"]');
    const nameRow = changed.filter({ hasText: '"name"' });
    await expect(nameRow.locator('[data-side="left"]')).toHaveText(`  "name": "${oldName}",`);
    await expect(nameRow.locator('[data-side="right"]')).toHaveText(`  "name": "${newName}",`);
    await expect(changed.filter({ hasText: '"amount"' })).toHaveCount(1);

    // A current row needs no explanation of where "before" came from.
    await expect(dialog.getByTestId('receipt-diff-version-notice')).toHaveCount(0);

    // Every line shows by default.
    await expect(dialog.locator('[data-kind="equal"]').first()).toBeVisible();
    await expect(dialog.getByTestId('receipt-diff-collapsed')).toHaveCount(0);
  });

  test('changes only collapses the unchanged lines and keeps every change', async ({ page }) => {
    await gotoOwnTask(page);
    await page.getByTestId('receipt-update-diff-open').click();

    const dialog = page.getByTestId('receipt-diff-dialog');
    const changedRows = dialog.locator('[data-testid="receipt-diff-row"]:not([data-kind="equal"])');
    const changedCount = await changedRows.count();

    await dialog.getByRole('tab', { name: /Changes only/ }).click();

    await expect(dialog.getByTestId('receipt-diff-collapsed').first()).toContainText('unchanged lines');
    await expect(changedRows).toHaveCount(changedCount);
  });
});
