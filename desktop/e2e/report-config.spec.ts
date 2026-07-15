import { expect, Page, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import { addFirstGroupToScope, addGroupingLevel, gotoReports, openComboboxAndPick } from './helpers/reports';

// The Report Builder is gated by app.reports.read, carried only by Legacy Admin.
test.use({ storageState: 'e2e/.auth/admin.json' });

async function downloadOnGenerate(page: Page) {
  const [download] = await Promise.all([
    page.waitForEvent('download'),
    page.getByTestId('report-generate').click(),
  ]);
  return download;
}

// Config-only coverage of the builder's sections. A scope group is added in
// beforeEach so the live preview is active; none of these assertions depend on
// seeded receipts (the preview renders at 0 rows).
test.describe('Report Builder — configuration', () => {
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
    await gotoReports(page);
    await addFirstGroupToScope(page);
    // Settle the initial preview so per-test picks start from a stable panel.
    await expect(page.getByTestId('report-receipt-count')).toBeVisible({ timeout: 20_000 });
  });

  test('adds and removes grouping levels', async ({ page }) => {
    const removes = page.getByTestId('report-grouping-remove');

    await addGroupingLevel(page, 'Paid By');
    await expect(removes).toHaveCount(1);
    await addGroupingLevel(page, 'Tag');
    await expect(removes).toHaveCount(2);

    await removes.first().click();
    await expect(removes).toHaveCount(1);
  });

  test('reorders grouping levels', async ({ page }) => {
    const labels = page.getByTestId('report-grouping-label');

    await addGroupingLevel(page, 'Paid By');
    await addGroupingLevel(page, 'Tag');
    await expect(labels).toHaveText(['Paid By', 'Tag']);

    // Moving the first level down swaps the order.
    await page.getByTestId('report-grouping-down').first().click();
    await expect(labels).toHaveText(['Tag', 'Paid By']);
  });

  test('switching detail mode shows or hides the columns section', async ({ page }) => {
    // Default is Aggregate: the custom-column UI is present.
    await expect(page.getByTestId('report-add-column')).toBeVisible();

    // Records mode: columns are the receipt fields, so the custom-column UI is gone.
    await page.getByTestId('report-detail-records').click();
    await expect(page.getByTestId('report-add-column')).toHaveCount(0);
    await expect(page.getByText('Switch to Aggregate to define custom columns')).toBeVisible();

    // Back to Aggregate: the Add column control returns.
    await page.getByTestId('report-detail-aggregate').click();
    await expect(page.getByTestId('report-add-column')).toBeVisible();
  });

  test('the custom period range reveals date pickers', async ({ page }) => {
    await openComboboxAndPick(
      page,
      page.getByRole('combobox', { name: 'Period covering' }),
      page.getByRole('option', { name: /Custom range/ }),
    );

    const start = page.getByLabel('Start', { exact: true });
    const end = page.getByLabel('End', { exact: true });
    await expect(start).toBeVisible();
    await expect(end).toBeVisible();

    await start.fill('01/15/2026');
    await end.fill('06/15/2026');
    await start.blur();

    // The preview still renders after choosing a custom range.
    await expect(page.getByTestId('report-receipt-count')).toBeVisible({ timeout: 20_000 });
  });

  test('the document title becomes the preview heading', async ({ page }) => {
    await expect(page.getByTestId('report-receipt-count')).toBeVisible({ timeout: 20_000 });

    await page.getByLabel('Title', { exact: true }).fill('Quarterly Summary');

    // The preview is server-rendered HTML in an iframe; assert the heading via its
    // srcdoc rather than reaching across the frame.
    await expect(page.getByTitle('Report preview')).toHaveAttribute(
      'srcdoc',
      /Quarterly Summary/,
      { timeout: 20_000 },
    );
  });

  test('format selection determines the generated file', async ({ page }) => {
    await expect(page.getByTestId('report-receipt-count')).toBeVisible({ timeout: 20_000 });
    // report-generate is an app-button; its disabled state is on the inner <button>.
    const generate = page.getByTestId('report-generate').locator('button');

    // Removing the only default format (PDF) disables generation.
    await page.getByTestId('report-format-pdf').click();
    await expect(generate).toBeDisabled();

    // XLSX only: a single .xlsx file.
    await page.getByTestId('report-format-xlsx').click();
    await expect(generate).toBeEnabled();
    expect((await downloadOnGenerate(page)).suggestedFilename()).toContain('.xlsx');

    // CSV + XLSX: a zipped bundle.
    await page.getByTestId('report-format-csv').click();
    expect((await downloadOnGenerate(page)).suggestedFilename()).toContain('.zip');
  });

  test('removing the scope group returns the preview to its placeholder', async ({ page }) => {
    await expect(page.getByTestId('report-receipt-count')).toBeVisible({ timeout: 20_000 });

    await page.getByTestId('report-scope-remove').first().click();

    await expect(page.getByText('Add a group and at least one column to see a preview.')).toBeVisible();
    await expect(page.getByTestId('report-receipt-count')).toHaveCount(0);
  });
});
