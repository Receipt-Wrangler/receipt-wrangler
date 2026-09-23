import { expect, Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateGroup,
  apiCreateReceipt,
  apiDeleteGroupById,
  apiDeleteReportTemplateById,
  apiGetUserId,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import {
  addGroupToScopeByName,
  gotoReportBuilder,
  gotoReports,
  openComboboxAndPick,
  waitForPreview,
} from './helpers/reports';

// The Report Builder is gated by app.reports.read, carried only by Legacy Admin.
test.use({ storageState: 'e2e/.auth/admin.json' });

// The same fields, in the same order, as the receipts table's quick date filter.
const DATE_FIELDS = ['Receipt Date', 'Resolved Date', 'Added At'];

function dateFieldSelect(page: Page) {
  return page.getByRole('combobox', { name: 'Date field', exact: true });
}

/** Points the period at a date field and waits for the preview to recompute. */
async function pickDateField(page: Page, label: string): Promise<void> {
  await Promise.all([
    waitForPreview(page),
    openComboboxAndPick(page, dateFieldSelect(page), page.getByRole('option', { name: label, exact: true })),
  ]);
}

/** Sets a custom period. Dates are typed in the datepicker's MM/DD/YYYY. */
async function setCustomPeriod(page: Page, start: string, end: string): Promise<void> {
  await openComboboxAndPick(
    page,
    page.getByRole('combobox', { name: /Period covering/ }),
    page.getByRole('option', { name: /Custom range/ }),
  );
  await page.getByLabel('Start', { exact: true }).fill(start);
  await page.getByLabel('End', { exact: true }).fill(end);
  await page.getByLabel('End', { exact: true }).blur();
}

/**
 * The count chip reads "receipt_long<N> receipts": its icon's ligature text sits
 * flush against the number, so a plain "2 receipts" would also match "12 receipts".
 */
async function expectReceiptCount(page: Page, count: number): Promise<void> {
  await expect(page.getByTestId('report-receipt-count')).toContainText(new RegExp(`(?<!\\d)${count} receipts`), {
    timeout: 20_000,
  });
}

// The period can cover any of the three receipt dates. Two receipts dated
// 2024-01-01 are seeded, one of them resolved: the server stamps a receipt's
// resolved date and its added-at date with "now", so a 2023-2024 window separates
// the receipt date from the other two, and a window from 2025 on catches them.
test.describe.serial('Report Builder — period date field', () => {
  const groupName = uniqueName('report-period-field-grp');
  let groupId: number;
  let createdTemplateId: number | undefined;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      const adminId = await apiGetUserId(api, creds('admin').username);
      groupId = (await apiCreateGroup(api, groupName)).id;
      await apiCreateReceipt(api, { groupId, paidByUserId: adminId, name: uniqueName('open') });
      await apiCreateReceipt(api, {
        groupId,
        paidByUserId: adminId,
        name: uniqueName('resolved'),
        status: 'RESOLVED',
      });
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        if (createdTemplateId !== undefined) {
          await apiDeleteReportTemplateById(api, createdTemplateId);
        }
        await apiDeleteGroupById(api, String(groupId));
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
    await gotoReportBuilder(page);
    await addGroupToScopeByName(page, groupName);
    await expect(page.getByTestId('report-receipt-count')).toBeVisible({ timeout: 20_000 });
  });

  test('offers the quick date filter fields in order, defaulting to Receipt Date', async ({ page }) => {
    await expect(dateFieldSelect(page)).toContainText('Receipt Date');
    await expect(page.getByText(/Resolves to .* on Receipt Date/)).toBeVisible();

    await expect(async () => {
      await dateFieldSelect(page).click();
      await expect(page.getByRole('option')).toHaveText(DATE_FIELDS, { timeout: 2000 });
    }).toPass({ timeout: 15_000 });
  });

  test('a past window matches only on the receipt date', async ({ page }) => {
    await setCustomPeriod(page, '01/01/2023', '12/31/2024');
    await expectReceiptCount(page, 2);

    await pickDateField(page, 'Resolved Date');
    await expectReceiptCount(page, 0);

    await pickDateField(page, 'Added At');
    await expectReceiptCount(page, 0);
  });

  test('a window from 2025 on matches on the resolved and added-at dates', async ({ page }) => {
    await setCustomPeriod(page, '01/01/2025', '12/31/2099');
    await expectReceiptCount(page, 0);

    await pickDateField(page, 'Resolved Date');
    await expectReceiptCount(page, 1);
    await expect(page.getByText(/Resolves to .* on Resolved Date/)).toBeVisible();

    await pickDateField(page, 'Added At');
    await expectReceiptCount(page, 2);

    // The drill-in lists what the report covers, so it narrows on the same field.
    await page.getByTestId('report-receipt-count').click();
    const drillIn = page.getByRole('dialog');
    await expect(drillIn.getByText(/2025-01-01 to 2099-12-31 on Added At/)).toBeVisible();
    await expect(drillIn.getByTestId('report-receipt-row')).toHaveCount(2);
  });

  test('a saved template reopens on its date field', async ({ page }) => {
    const templateName = uniqueName('period-date-field');
    await pickDateField(page, 'Resolved Date');
    await page.getByLabel('Report name').fill(templateName);

    const save = page.getByTestId('report-save-template');
    await expect(save.locator('button')).toBeEnabled();
    const [response] = await Promise.all([
      page.waitForResponse((r) => r.url().includes('/api/report/template') && r.request().method() === 'POST'),
      save.click(),
    ]);
    expect(response.status()).toBe(200);
    const saved = (await response.json()) as { id: number; configuration: { period: { dateField?: string } } };
    createdTemplateId = saved.id;
    expect(saved.configuration.period.dateField).toBe('resolvedDate');

    await gotoReports(page);
    await page.getByRole('row').filter({ hasText: templateName }).getByTestId('report-template-name').click();
    await expect(page).toHaveURL(/\/reports\/\d+\/edit$/);
    await expect(dateFieldSelect(page)).toContainText('Resolved Date');
  });
});
