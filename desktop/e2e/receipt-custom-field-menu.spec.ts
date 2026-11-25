import { expect, Page, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiCreateCustomField,
  apiDeleteCustomFieldById,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';

// The menu is gated on app.custom-fields.read, and seeding/removing the field
// needs .create/.delete - all Legacy Admin.
test.use({ storageState: 'e2e/.auth/admin.json' });

const EMPTY_ROW = 'No items found';

/**
 * The receipt form's "Manage custom fields" menu (`app-filtered-stateful-menu`).
 *
 * #748: its empty-state row was guarded by `filteredItems.length`, the computed
 * signal's parameter count (always 0), so "No items found" rendered above the
 * real items every time. This pins that the row appears only when the filter
 * matches nothing.
 *
 * It seeds its own custom field so the menu is never empty on a fresh database,
 * and only reads it, so the tests can run in parallel.
 */
test.describe('Receipt form — manage custom fields menu', () => {
  const fieldName = uniqueName('menu-cf');
  let fieldId: number;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      fieldId = (await apiCreateCustomField(api, { name: fieldName, type: 'TEXT' })).id;
    });
  });

  test.afterAll(async () => {
    await withAdminApi(async (api) => {
      if (fieldId) {
        await apiDeleteCustomFieldById(api, fieldId);
      }
    });
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  async function openMenu(page: Page): Promise<void> {
    await page.goto('/receipts/add');
    // .locator('button'): the cdkMenuTrigger also puts role="button" on the
    // <app-button> host, so getByRole('button') matches two elements.
    await page.getByTestId('receipt-manage-custom-fields').locator('button').click();
  }

  const fieldItem = (page: Page) => page.getByRole('menuitem', { name: fieldName });

  test('lists the custom fields without a "No items found" row', async ({ page }) => {
    await openMenu(page);

    await expect(fieldItem(page)).toBeVisible();
    await expect(page.getByText(EMPTY_ROW)).toHaveCount(0);
  });

  test('shows "No items found" only when the filter matches nothing', async ({ page }) => {
    await openMenu(page);
    const filter = page.getByLabel('Filter custom fields');

    await filter.fill(uniqueName('no-such-field'));
    await expect(page.getByText(EMPTY_ROW)).toBeVisible();
    await expect(fieldItem(page)).toHaveCount(0);

    await filter.fill(fieldName);
    await expect(fieldItem(page)).toBeVisible();
    await expect(page.getByText(EMPTY_ROW)).toHaveCount(0);
  });
});
