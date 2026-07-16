import { expect, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import { apiDeleteReportTemplateById, withAdminApi } from './helpers/provisioning';
import { addFirstGroupToScope, gotoReportBuilder } from './helpers/reports';

// Saving a report template is gated by the app-level app.reports.create permission,
// which the seeded Legacy Admin role carries (add-only reconciliation). Run as admin.
test.use({ storageState: 'e2e/.auth/admin.json' });

test.describe('Report Builder — save template', () => {
  // The id of the template created by the running test, torn down afterwards so
  // runs don't accumulate saved templates in the DB.
  let createdTemplateId: number | undefined;

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test.afterEach(async () => {
    if (createdTemplateId === undefined) {
      return;
    }
    try {
      await withAdminApi((api) => apiDeleteReportTemplateById(api, createdTemplateId!));
    } catch {
      // Best-effort cleanup — don't mask a test failure with a teardown error.
    }
    createdTemplateId = undefined;
  });

  test('saves the current configuration as a template and confirms with a snackbar', async ({ page }) => {
    await gotoReportBuilder(page);
    await addFirstGroupToScope(page);

    // With a group and the default columns/format the config is buildable — the
    // preview activates and the Save Template button (Generate's validity plus a
    // non-empty name) enables.
    await expect(page.getByTestId('report-receipt-count')).toBeVisible({ timeout: 20_000 });

    const save = page.getByTestId('report-save-template');
    await expect(save.locator('button')).toBeEnabled();

    // Clicking it POSTs the current configuration to /api/report/template (200).
    const [response] = await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/report/template') && r.request().method() === 'POST',
      ),
      save.click(),
    ]);
    expect(response.status()).toBe(200);

    // Remember the new template so afterEach can delete it, and confirm the backend
    // stamped it with the current config schema version.
    const body = (await response.json()) as { id: number; configurationVersion: number };
    createdTemplateId = body.id;
    expect(body.configurationVersion).toBe(1);

    // The success snackbar confirms the save.
    await expect(page.getByText('Template saved')).toBeVisible();
  });
});
