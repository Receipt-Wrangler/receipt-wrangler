import { expect, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateGroup,
  apiCreateReceipt,
  apiDeleteGroupById,
  apiGetUserId,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import { addFirstGroupToScope, addGroupToScopeByName, gotoReportBuilder, openComboboxAndPick } from './helpers/reports';

// The Report Builder is gated by the app-level app.reports.read permission, which
// the seeded Legacy Admin role carries (and Legacy User does not). Run the positive
// flow as admin.
test.use({ storageState: 'e2e/.auth/admin.json' });

test.describe('Report Builder', () => {
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('admin builds a report, sees a live preview, and downloads it', async ({ page }) => {
    await gotoReportBuilder(page);

    // The reports route opts into the shell's full-height frame: the content area
    // becomes a bounded flex column and the outlet drops the default p-4 padding so
    // the two panes scroll independently under the sticky page-bar/generate-bar.
    await expect(page.locator('mat-drawer-content.drawer-content--full-height')).toBeVisible();
    const outletPadding = await page.locator('.drawer-outlet').evaluate(
      (el) => getComputedStyle(el).padding,
    );
    expect(outletPadding).toBe('0px');

    await addFirstGroupToScope(page);

    // With a group and the default columns, the preview activates and its
    // receipt-count chip appears (it renders the engine's HTML server-side).
    await expect(page.getByTestId('report-receipt-count')).toBeVisible({ timeout: 20_000 });
    await expect(page.getByTitle('Report preview')).toBeVisible();

    // Choose a single CSV output (default is PDF) and generate; expect a download.
    await page.getByTestId('report-format-pdf').click();
    await page.getByTestId('report-format-csv').click();

    const downloadPromise = page.waitForEvent('download');
    await page.getByTestId('report-generate').click();
    const download = await downloadPromise;
    expect(download.suggestedFilename()).toContain('.csv');
  });

  test('disables an aggregate dimension column that is neither the aggregate-by nor a grouping level', async ({ page }) => {
    await gotoReportBuilder(page);

    // Scope to a group (the disable behavior is independent of the data).
    await addFirstGroupToScope(page);

    // Default is aggregate-by = Category with a valid Category dimension column.
    // Switch aggregate-by to Tag — now the Category column reads neither the
    // aggregate-by (Tag) nor a grouping level, so it becomes invalid.
    await openComboboxAndPick(
      page,
      page.locator('mat-select').filter({ hasText: 'Category' }),
      page.locator('mat-option').filter({ hasText: 'Tag' }),
    );

    // The Category column is kept but shown disabled (not removed), and the
    // preview still renders from the remaining columns — no raw engine error.
    await expect(page.getByTestId('report-column-disabled')).toBeVisible();
    await expect(page.getByTestId('report-receipt-count')).toBeVisible({ timeout: 20_000 });
    await expect(page.getByTitle('Report preview')).toBeVisible();
  });

  test('adds a grouping level and a filter through the shared app-select pickers', async ({ page }) => {
    await gotoReportBuilder(page);
    await addFirstGroupToScope(page);

    // Grouping: the "Add grouping level…" picker is an app-select (a combobox),
    // not a native <select>. Picking a dimension appends a grouping row.
    await openComboboxAndPick(
      page,
      page.getByRole('combobox', { name: /Add grouping level/ }),
      page.getByRole('option', { name: 'Paid By', exact: true }),
    );
    await expect(page.getByTestId('report-grouping-remove')).toBeVisible();

    // Filters: the "Add filter…" picker is likewise an app-select. Picking a field
    // adds its filter row.
    await openComboboxAndPick(
      page,
      page.getByRole('combobox', { name: /Add filter/ }),
      page.getByRole('option', { name: 'Name', exact: true }),
    );
    await expect(page.getByTestId('report-filter-remove')).toBeVisible();
  });

  test('the /reports route is gated by app.reports.read', async ({ browser }) => {
    // A regular user (Legacy User) lacks app.reports.read, so the route guard
    // redirects them away from /reports.
    const context = await browser.newContext({ storageState: 'e2e/.auth/user.json' });
    const page = await context.newPage();
    await stubTokenRefresh(page);

    await page.goto('/reports');
    await expect(page).not.toHaveURL(/\/reports/, { timeout: 10_000 });

    await context.close();
  });
});

// The drill-in flow needs a receipt to exist, so it provisions its own group +
// receipt (dated 2024-01-01) via the API and tears them down — deterministic and
// parallel-safe (unique names), rather than depending on seeded data.
test.describe('Report Builder — receipt drill-in', () => {
  const groupName = uniqueName('report-drill-grp');
  const receiptName = uniqueName('report-drill-rcpt');
  let groupId: number;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      const adminId = await apiGetUserId(api, creds('admin').username);
      groupId = (await apiCreateGroup(api, groupName)).id;
      await apiCreateReceipt(api, { groupId, paidByUserId: adminId, name: receiptName });
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi((api) => apiDeleteGroupById(api, String(groupId)));
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('drills into a receipt and opens the full receipt in a new tab', async ({ page }) => {
    await gotoReportBuilder(page);
    await addGroupToScopeByName(page, groupName);

    // Cover the provisioned 2024-01-01 receipt with a generous custom range (a
    // full year of slack so time zones can't push it out of the window).
    await openComboboxAndPick(
      page,
      page.getByRole('combobox', { name: /Period covering/ }),
      page.getByRole('option', { name: /Custom range/ }),
    );
    await page.getByLabel('Start', { exact: true }).fill('01/01/2023');
    await page.getByLabel('End', { exact: true }).fill('12/31/2024');
    await page.getByLabel('End', { exact: true }).blur();

    // The preview now covers exactly the one provisioned receipt.
    const chip = page.getByTestId('report-receipt-count');
    await expect(chip).toContainText('1 receipts', { timeout: 20_000 });
    await chip.click();

    await expect(page.getByText('Receipts in this report')).toBeVisible();

    // Click the receipt row to open its breakdown, then open the full receipt.
    const row = page.getByTestId('report-receipt-row').first();
    await expect(row).toBeVisible();
    await row.click();

    const openFull = page.getByTestId('report-receipt-open-full');
    await expect(openFull).toBeVisible();

    const [popup] = await Promise.all([page.waitForEvent('popup'), openFull.click()]);
    await expect(popup).toHaveURL(/\/receipts\/\d+\/view/);

    // Back returns to the list.
    await page.getByTestId('report-receipt-back').click();
    await expect(row).toBeVisible();
  });
});
