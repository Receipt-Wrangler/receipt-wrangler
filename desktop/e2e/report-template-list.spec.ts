import { expect, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiCreateReportTemplate,
  apiDeleteReportTemplateById,
  apiFirstReportCategory,
  apiFirstReportGroupId,
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

  test('hydrates a saved category filter into the builder on open', async ({ page }) => {
    // Regression guard: opening a template that filters on a category used to render
    // no filter row at all (the visible-row set wasn't seeded from the hydrated form).
    // Seed a template filtering on a real category, open it, and assert the row and the
    // selected category's chip are present.
    const filterName = uniqueName('report-template-catfilter');
    const { groupId, categoryId, categoryName } = await withAdminApi((api) =>
      apiFirstReportCategory(api),
    );
    const filtered = await withAdminApi((api) =>
      apiCreateReportTemplate(api, {
        name: filterName,
        groupIds: [groupId],
        filter: { categories: { operation: 'CONTAINS', value: [categoryId] } },
      }),
    );
    extraIds.push(filtered.id);

    await gotoReports(page);
    const row = page.getByRole('row').filter({ hasText: filterName });
    await expect(row).toBeVisible();

    await row.getByTestId('report-template-name').click();
    await expect(page).toHaveURL(/\/reports\/\d+\/edit$/);

    // The category filter row is rendered (its remove control) and the stored category
    // resolves to its name chip — proving both the row-seeding fix and that the value
    // carried over.
    await expect(page.getByTestId('report-filter-remove')).toBeVisible();
    await expect(page.getByText(categoryName, { exact: true })).toBeVisible();
  });

  test('updates an opened template in place (same id, no new row)', async ({ page }) => {
    await gotoReports(page);
    await page
      .getByRole('row')
      .filter({ hasText: seededName })
      .getByTestId('report-template-name')
      .click();
    await expect(page).toHaveURL(/\/reports\/\d+\/edit$/);

    // Hydrated: the name field carries the stored name; the edit-route Save is "Update".
    const nameInput = page.getByLabel('Report name');
    await expect(nameInput).toHaveValue(seededName);
    await expect(page.getByTestId('report-save-template')).toContainText('Update Template');

    // Rename and save — the edit route updates in place (PUT the seeded id), never creates.
    const renamed = uniqueName('report-template-renamed');
    await nameInput.fill(renamed);
    const [response] = await Promise.all([
      page.waitForResponse(
        (r) => /\/api\/report\/template\/\d+$/.test(r.url()) && r.request().method() === 'PUT',
      ),
      page.getByTestId('report-save-template').click(),
    ]);
    expect(response.status()).toBe(200);
    expect(new URL(response.url()).pathname.endsWith(`/template/${seededId}`)).toBe(true);

    // Back on the list: the same row is renamed and the old name is gone — no duplicate.
    await gotoReports(page);
    await expect(page.getByRole('row').filter({ hasText: renamed })).toHaveCount(1);
    await expect(page.getByRole('row').filter({ hasText: seededName })).toHaveCount(0);
  });

  test('hydrates a saved "report generator" paid-by filter into the builder', async ({ page }) => {
    // A template can filter paid-by on the dynamic report-generator sentinel (-1) rather
    // than a static user. Opening it must render the paid-by row and resolve -1 to the
    // pinned "Whoever generates the report" option chip.
    const filterName = uniqueName('report-template-mefilter');
    const groupId = await withAdminApi((api) => apiFirstReportGroupId(api));
    const filtered = await withAdminApi((api) =>
      apiCreateReportTemplate(api, {
        name: filterName,
        groupIds: [groupId],
        filter: { paidBy: { operation: 'CONTAINS', value: [-1] } },
      }),
    );
    extraIds.push(filtered.id);

    await gotoReports(page);
    await page
      .getByRole('row')
      .filter({ hasText: filterName })
      .getByTestId('report-template-name')
      .click();
    await expect(page).toHaveURL(/\/reports\/\d+\/edit$/);

    await expect(page.getByTestId('report-filter-remove')).toBeVisible();
    await expect(page.getByText('Whoever generates the report', { exact: true })).toBeVisible();
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
