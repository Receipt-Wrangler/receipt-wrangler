import { expect, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiCreateReportTemplate,
  apiDeleteReportTemplateById,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import { gotoReports } from './helpers/reports';

// The Reports list + its row actions (generate / open / duplicate / delete) run as
// admin, who holds every app.reports.* permission (Legacy Admin, add-only
// reconciliation). Each test seeds its own uniquely-named template through the API
// and tears down what it created (plus any duplicate) so runs don't accumulate rows.
test.use({ storageState: 'e2e/.auth/admin.json' });

test.describe('Reports — templates list', () => {
  let seededId: number | undefined;
  let seededName: string;
  const extraIds: number[] = [];

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
    seededName = uniqueName('report-template');
    const template = await withAdminApi((api) => apiCreateReportTemplate(api, { name: seededName }));
    seededId = template.id;
  });

  test.afterEach(async () => {
    const ids = [seededId, ...extraIds].filter((id): id is number => id !== undefined);
    try {
      await withAdminApi(async (api) => {
        for (const id of ids) {
          await apiDeleteReportTemplateById(api, id);
        }
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a teardown error.
    }
    seededId = undefined;
    extraIds.length = 0;
  });

  test('lists a saved template and opens it in the builder', async ({ page }) => {
    await gotoReports(page);

    const row = page.getByRole('row').filter({ hasText: seededName });
    await expect(row).toBeVisible();

    // Open the template in the builder via its name link.
    await row.getByTestId('report-template-name').click();
    await expect(page).toHaveURL(/\/reports\/\d+\/edit$/);
    await expect(page.getByText('Report Builder')).toBeVisible();

    // The breadcrumb shows the loaded template name and a live preview fires — both
    // only hold if the stored (scoped, runnable) configuration hydrated the form. A
    // blank "New report" builder has no scope, so it never reaches a preview.
    await expect(page.getByText(seededName)).toBeVisible();
    await expect(page.getByTestId('report-receipt-count')).toBeVisible({ timeout: 20_000 });
  });

  test('generates a report from a row action', async ({ page }) => {
    await gotoReports(page);
    const row = page.getByRole('row').filter({ hasText: seededName });

    const [response] = await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/report/generate') && r.request().method() === 'POST',
      ),
      row.getByTestId('report-template-generate').click(),
    ]);
    expect(response.status()).toBe(200);
  });

  test('duplicates a template from a row action', async ({ page }) => {
    await gotoReports(page);
    const row = page.getByRole('row').filter({ hasText: seededName });

    const [response] = await Promise.all([
      page.waitForResponse(
        (r) =>
          /\/api\/report\/template\/\d+\/duplicate$/.test(r.url()) && r.request().method() === 'POST',
      ),
      row.getByTestId('report-template-duplicate').click(),
    ]);
    expect(response.status()).toBe(200);
    const duplicate = (await response.json()) as { id: number };
    extraIds.push(duplicate.id);

    // The list reloads and the copy (name + " duplicate") appears.
    await expect(page.getByRole('row').filter({ hasText: `${seededName} duplicate` })).toBeVisible();
  });

  test('deletes a template from a row action', async ({ page }) => {
    await gotoReports(page);
    const row = page.getByRole('row').filter({ hasText: seededName });

    await row.getByTestId('report-template-delete').click();
    await page.getByTestId('dialog-submit-button').click();

    await expect(page.getByRole('row').filter({ hasText: seededName })).toHaveCount(0);
    // The test already deleted it; afterEach's delete is a harmless idempotent no-op.
  });
});
