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
