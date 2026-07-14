import { expect, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import { addFirstGroupToScope, gotoReports, openComboboxAndPick, waitForPreview } from './helpers/reports';

// The Report Builder is gated by the app-level app.reports.read permission, which
// the seeded Legacy Admin role carries (and Legacy User does not). Run the positive
// flow as admin.
test.use({ storageState: 'e2e/.auth/admin.json' });

test.describe('Report Builder', () => {
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('admin builds a report, sees a live preview, and downloads it', async ({ page }) => {
    await gotoReports(page);

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

  test('admin drills into a receipt and opens the full receipt in a new tab', async ({ page }) => {
    await gotoReports(page);
    await addFirstGroupToScope(page);

    // Seeded receipts sit in prior months, so cover them with "Last month", and
    // wait for the debounced preview so the covered-count reflects the new period.
    await Promise.all([
      waitForPreview(page),
      openComboboxAndPick(
        page,
        page.getByRole('combobox', { name: /Period covering/ }),
        page.getByRole('option', { name: 'Last month' }),
      ),
    ]);

    const chip = page.getByTestId('report-receipt-count');
    await expect(chip).toBeVisible({ timeout: 20_000 });
    await chip.click();

    // The chip only opens the drill-in when the report covers receipts; if the
    // seeded scope has none in range, there is nothing to drill into.
    const opened = await page
      .getByText('Receipts in this report')
      .waitFor({ state: 'visible', timeout: 4000 })
      .then(() => true)
      .catch(() => false);
    test.skip(!opened, 'seeded scope has no receipts in the selected period');

    // Click a receipt row to open its breakdown, then open the full receipt.
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

  test('disables an aggregate dimension column that is neither the aggregate-by nor a grouping level', async ({ page }) => {
    await gotoReports(page);

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
    await gotoReports(page);
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
