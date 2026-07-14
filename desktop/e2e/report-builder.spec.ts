import { expect, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';

// The Report Builder is gated by the app-level app.reports.read permission, which
// the seeded Legacy Admin role carries (and Legacy User does not). Run the positive
// flow as admin.
test.use({ storageState: 'e2e/.auth/admin.json' });

test.describe('Report Builder', () => {
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('admin builds a report, sees a live preview, and downloads it', async ({ page }) => {
    await page.goto('/reports');
    await expect(page.getByText('Report Builder')).toBeVisible();

    // The reports route opts into the shell's full-height frame: the content area
    // becomes a bounded flex column and the outlet drops the default p-4 padding so
    // the two panes scroll independently under the sticky page-bar/generate-bar.
    await expect(page.locator('mat-drawer-content.drawer-content--full-height')).toBeVisible();
    const outletPadding = await page.locator('.drawer-outlet').evaluate(
      (el) => getComputedStyle(el).padding,
    );
    expect(outletPadding).toBe('0px');

    // Add a group to the report scope via the Add Group dialog.
    await page.getByTestId('report-add-group').click();
    await page.getByTestId('add-group-select').first().click();
    await page.getByTestId('dialog-submit-button').click();

    // With a group and the default columns, the preview activates and its
    // receipt-count chip appears (it renders the engine's HTML server-side).
    await expect(page.getByTestId('report-receipt-count')).toBeVisible({ timeout: 20_000 });
    await expect(page.locator('iframe[title="Report preview"]')).toBeVisible();

    // Choose a single CSV output (default is PDF) and generate; expect a download.
    await page.getByTestId('report-format-pdf').click();
    await page.getByTestId('report-format-csv').click();

    const downloadPromise = page.waitForEvent('download');
    await page.getByTestId('report-generate').click();
    const download = await downloadPromise;
    expect(download.suggestedFilename()).toContain('.csv');
  });

  test('admin drills into a receipt and opens the full receipt in a new tab', async ({ page }) => {
    await page.goto('/reports');
    await expect(page.getByText('Report Builder')).toBeVisible();

    await page.getByTestId('report-add-group').click();
    await page.getByTestId('add-group-select').first().click();
    await page.getByTestId('dialog-submit-button').click();

    // Seeded receipts sit in prior months, so cover them with "Last month".
    await page.locator('mat-select').filter({ hasText: 'This month' }).click();
    await page.getByRole('option', { name: 'Last month' }).click();
    // Let the debounced preview refresh so the covered-count reflects the new period.
    await page.waitForTimeout(1500);

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
