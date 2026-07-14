import { expect, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import { addFirstGroupToScope, gotoReports, openComboboxAndPick } from './helpers/reports';

// The Report Builder is gated by app.reports.read, carried only by Legacy Admin.
test.use({ storageState: 'e2e/.auth/admin.json' });

// These specs exercise the column picker (add/edit/reorder/remove and the three
// column kinds). They are config-only: the default report (Category dimension,
// Count + Total aggregates, Aggregate detail mode) exists without any receipts,
// so nothing here depends on seeded data.
test.describe('Report Builder — columns', () => {
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
    await gotoReports(page);
    await addFirstGroupToScope(page);
  });

  test('adds aggregate SUM and COUNT columns', async ({ page }) => {
    const rows = page.getByTestId('report-column-row');
    const before = await rows.count();

    // SUM(amount): the aggregate step presets SUM with the first measure (Amount).
    await page.getByTestId('report-add-column').click();
    await page.getByTestId('picker-kind-aggregate').click();
    await page.getByTestId('picker-save').click();

    // COUNT: choosing COUNT drops the Measure select.
    await page.getByTestId('report-add-column').click();
    await page.getByTestId('picker-kind-aggregate').click();
    await page.getByRole('button', { name: 'COUNT' }).click();
    await expect(page.getByRole('combobox', { name: 'Measure' })).toHaveCount(0);
    await page.getByTestId('picker-save').click();

    await expect(rows).toHaveCount(before + 2);
  });

  test('builds a formula column from clicked chips', async ({ page }) => {
    await page.getByTestId('report-add-column').click();
    await page.getByTestId('picker-kind-formula').click();

    await page.getByLabel('Column label').fill('Average');
    // The referenceable columns are the aggregates Total and Count (a dimension
    // cannot be referenced), inserted by their engine name via the chips.
    await page.getByTestId('picker-col-Total').click();
    await page.getByTestId('picker-op-/').click();
    await page.getByTestId('picker-col-Count').click();
    await expect(page.locator('.picker__status--ok')).toBeVisible();

    await page.getByTestId('picker-save').click();
    await expect(page.getByTestId('report-column-row').filter({ hasText: 'Average' })).toBeVisible();
  });

  test('the formula picker builds by chip and supports backspace/clear', async ({ page }) => {
    await page.getByTestId('report-add-column').click();
    await page.getByTestId('picker-kind-formula').click();

    // picker-save is an app-button; its disabled state lives on the inner <button>.
    const save = page.getByTestId('picker-save').locator('button');
    const exprText = page.locator('.picker__expr-text');

    // A label is required, so an unlabelled formula cannot be saved.
    await expect(save).toBeDisabled();

    // The expression is built only by clicking chips (it is read-only).
    await page.getByLabel('Column label').fill('Average');
    await page.getByTestId('picker-col-Total').click();
    await page.getByTestId('picker-op-/').click();
    await page.getByTestId('picker-col-Count').click();
    await expect(exprText).toHaveText('Total / Count');
    await expect(save).toBeEnabled();

    // Backspace drops the last token; clear empties the expression.
    await page.getByTestId('picker-formula-backspace').click();
    await expect(exprText).toHaveText('Total /');
    await page.getByTestId('picker-formula-clear').click();
    await expect(exprText).toHaveText('Click a column, then operators, to build the formula');
  });

  test('adds an enabled dimension column when grouped by that field', async ({ page }) => {
    // Group by Paid By so a Paid By dimension column is valid in aggregate mode.
    await openComboboxAndPick(
      page,
      page.getByRole('combobox', { name: /Add grouping level/ }),
      page.getByRole('option', { name: 'Paid By', exact: true }),
    );
    await expect(page.getByTestId('report-grouping-remove')).toBeVisible();

    // The dimension step presets the Field to the first dimension (Paid By).
    await page.getByTestId('report-add-column').click();
    await page.getByTestId('picker-kind-dimension').click();
    await page.getByTestId('picker-save').click();

    // The new column is enabled (a valid grouped dimension), not disabled.
    await expect(page.getByTestId('report-column-row').filter({ hasText: 'Paid By' })).toBeVisible();
    await expect(page.getByTestId('report-column-disabled').filter({ hasText: 'Paid By' })).toHaveCount(0);
  });

  test('edits, reorders, and removes columns', async ({ page }) => {
    const rows = page.getByTestId('report-column-row');
    const initial = await rows.count();

    // Edit: rename the Total aggregate.
    await rows.filter({ hasText: 'Total' }).getByTestId('report-column-edit').click();
    await page.getByLabel('Column label').fill('Grand Total');
    await page.getByTestId('picker-save').click();
    await expect(rows.filter({ hasText: 'Grand Total' })).toBeVisible();

    // Reorder: moving the first column down changes which column leads the list.
    const firstBefore = await rows.first().innerText();
    await rows.first().getByTestId('report-column-down').click();
    await expect(rows.first()).not.toHaveText(firstBefore);

    // Remove: one fewer column row.
    await rows.first().getByTestId('report-column-remove').click();
    await expect(rows).toHaveCount(initial - 1);
  });
});
