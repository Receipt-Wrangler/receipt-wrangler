import { expect, Page, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import { apiDeleteReportTemplateById, uniqueName, withAdminApi } from './helpers/provisioning';
import {
  addFirstGroupToScope,
  addGroupingLevel,
  gotoReportBuilder,
  gotoReports,
  waitForPreview,
} from './helpers/reports';

// The Report Builder is gated by app.reports.read, carried only by Legacy Admin.
test.use({ storageState: 'e2e/.auth/admin.json' });

/** The preview is server-rendered HTML in an iframe; assert against its srcdoc. */
function preview(page: Page) {
  return page.getByTitle('Report preview');
}

/** Matches a heading rendered as its own cell, not a substring of some other text. */
function asCell(text: string): RegExp {
  return new RegExp(`>\\s*${text}\\s*<`);
}

/** Opens the rename dialog for the grouping level at `index` and returns its label input. */
async function openGroupingRename(page: Page, index = 0) {
  await page.getByTestId('report-grouping-edit').nth(index).click();
  const label = page.getByLabel('Column label');
  await expect(label).toBeVisible();
  return label;
}

/** Renames the open dialog's column and waits for the preview to catch up. */
async function saveRename(page: Page, heading: string) {
  await page.getByLabel('Column label').fill(heading);
  await Promise.all([waitForPreview(page), page.getByTestId('picker-save').click()]);
}

// A grouping level adds a leading column to the rendered report whose heading the
// backend derives from the field catalog. These cover the pencil that overrides it:
// the heading reaches the real rendered report (the preview is the engine's own
// output), resets when the default is retyped, and survives a save/reopen.
test.describe('Report Builder — grouping column headings', () => {
  // The template created by the running test, torn down afterwards so runs don't
  // accumulate saved templates.
  let createdTemplateId: number | undefined;

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
    await gotoReportBuilder(page);
    await addFirstGroupToScope(page);
    // Settle the initial preview so the grouping picks start from a stable panel.
    await expect(page.getByTestId('report-receipt-count')).toBeVisible({ timeout: 20_000 });
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

  test('renames the column a grouping level adds to the report', async ({ page }) => {
    const heading = uniqueName('Payer');

    await addGroupingLevel(page, 'Paid By');
    // The grouping level's column is rendered by the backend, headed by the field
    // catalog's own label until it is overridden.
    await expect(preview(page)).toHaveAttribute('srcdoc', asCell('Paid By'), { timeout: 20_000 });

    await openGroupingRename(page);

    // The dialog names the grouping column rather than claiming to add one, and the
    // field is locked: it is chosen by the grouping level, not here.
    await expect(page.getByText('Grouping column')).toBeVisible();
    await expect(page.getByTestId('picker-back')).toHaveCount(0);
    const field = page.getByLabel('Field');
    await expect(field).toHaveJSProperty('readOnly', true);
    await expect(field).toHaveValue('Paid By');
    // The label is seeded with the heading the column currently carries.
    await expect(page.getByLabel('Column label')).toHaveValue('Paid By');

    await saveRename(page, heading);

    // The grouping row shows the heading its column will carry...
    await expect(page.getByTestId('report-grouping-label')).toHaveText(heading);
    // ...and so does the rendered report, in place of the catalog label.
    await expect(preview(page)).toHaveAttribute('srcdoc', asCell(heading), { timeout: 20_000 });
    expect(await preview(page).getAttribute('srcdoc')).not.toMatch(asCell('Paid By'));
  });

  test('resets to the field label when the default is retyped', async ({ page }) => {
    await addGroupingLevel(page, 'Paid By');
    await expect(preview(page)).toHaveAttribute('srcdoc', asCell('Paid By'), { timeout: 20_000 });

    const heading = uniqueName('Payer');
    await openGroupingRename(page);
    await saveRename(page, heading);
    await expect(page.getByTestId('report-grouping-label')).toHaveText(heading);
    await expect(preview(page)).toHaveAttribute('srcdoc', asCell(heading), { timeout: 20_000 });

    // Retyping the field's own label is how a user resets: no override is stored, so
    // the report goes back to the catalog heading.
    await openGroupingRename(page);
    await saveRename(page, 'Paid By');

    await expect(page.getByTestId('report-grouping-label')).toHaveText('Paid By');
    await expect(preview(page)).toHaveAttribute('srcdoc', asCell('Paid By'), { timeout: 20_000 });
    expect(await preview(page).getAttribute('srcdoc')).not.toMatch(asCell(heading));
  });

  test('persists the heading into a saved template and its list row', async ({ page }) => {
    const heading = uniqueName('Payer');
    const templateName = uniqueName('grouping-heading');

    await addGroupingLevel(page, 'Paid By');
    await openGroupingRename(page);
    await saveRename(page, heading);

    await page.getByLabel('Report name').fill(templateName);

    const save = page.getByTestId('report-save-template');
    await expect(save.locator('button')).toBeEnabled();
    const [response] = await Promise.all([
      page.waitForResponse(
        (r) => r.url().includes('/api/report/template') && r.request().method() === 'POST',
      ),
      save.click(),
    ]);
    expect(response.status()).toBe(200);
    createdTemplateId = ((await response.json()) as { id: number }).id;

    // The list summarizes the grouping with the heading the report renders, not the
    // underlying field's name.
    await gotoReports(page);
    const row = page.getByRole('row').filter({ hasText: templateName });
    await expect(row).toContainText(heading);
    await expect(row).not.toContainText('Paid By');

    // Reopening rehydrates the override, so it is not lost on the next edit.
    await row.getByTestId('report-template-name').click();
    await expect(page).toHaveURL(/\/reports\/\d+\/edit$/);
    await expect(page.getByTestId('report-grouping-label')).toHaveText(heading);
  });
});
