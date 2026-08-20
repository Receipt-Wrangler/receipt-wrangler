import { expect, Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateCustomField,
  apiCreateGroup,
  apiCreateReceipt,
  apiCreateReportTemplate,
  apiDeleteCustomFieldById,
  apiDeleteGroupById,
  apiDeleteReportTemplateById,
  apiGetUserId,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import {
  addGroupingLevel,
  addGroupToScopeByName,
  gotoReportBuilder,
  gotoReports,
  openComboboxAndPick,
  waitForPreview,
} from './helpers/reports';

// The builder is gated by app.reports.read, and creating/deleting custom fields
// needs app.custom-fields.create/.delete — all Legacy Admin.
test.use({ storageState: 'e2e/.auth/admin.json' });

// Serial: every test shares one seeded group + custom-field pool, and the pool is
// global, so running them concurrently would multiply the seeding.
test.describe.configure({ mode: 'serial' });

/**
 * Custom fields in the Report Builder: a field of ANY type can be grouped by, a
 * currency one can also be aggregated, a date one offers derived calendar-period
 * levels, and every one is badged "Custom".
 *
 * The spec seeds its own group so the live preview contains exactly these four
 * receipts — the assertions are about rendered bucket text, which shared data
 * would make non-deterministic. Deleting the group cascades the receipts, which
 * matters because deleting a custom field destroys every value stored against it.
 */
test.describe('Report Builder — custom fields', () => {
  const groupName = uniqueName('cf-report-grp');
  const tipName = uniqueName('Tip');
  const dueName = uniqueName('Due');
  const reimbursedName = uniqueName('Reimbursed');

  let groupId: number;
  let tipId: number;
  let dueId: number;
  let reimbursedId: number;
  let templateId: number;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      const adminId = await apiGetUserId(api, creds('admin').username);
      groupId = (await apiCreateGroup(api, groupName)).id;

      tipId = (await apiCreateCustomField(api, { name: tipName, type: 'CURRENCY' })).id;
      dueId = (await apiCreateCustomField(api, { name: dueName, type: 'DATE' })).id;
      reimbursedId = (await apiCreateCustomField(api, { name: reimbursedName, type: 'BOOLEAN' })).id;

      // 1500.50 and 1500.5 are the same amount written two ways: they must land in
      // ONE bucket (bucket keys are canonical decimals). The fourth receipt carries
      // no custom field values at all, which is the (None) bucket.
      const receipts: { name: string; tip?: string; due?: string; reimbursed?: boolean }[] = [
        { name: 'cf-a', tip: '1500.50', due: '2024-03-15T00:00:00Z', reimbursed: true },
        { name: 'cf-b', tip: '1500.5', due: '2024-03-20T00:00:00Z', reimbursed: false },
        { name: 'cf-c', tip: '2.00', due: '2024-04-05T00:00:00Z', reimbursed: true },
        { name: 'cf-none' },
      ];

      for (const receipt of receipts) {
        const customFields = [];
        if (receipt.tip !== undefined) {
          customFields.push({ customFieldId: tipId, currencyValue: receipt.tip });
        }
        if (receipt.due !== undefined) {
          customFields.push({ customFieldId: dueId, dateValue: receipt.due });
        }
        if (receipt.reimbursed !== undefined) {
          customFields.push({ customFieldId: reimbursedId, booleanValue: receipt.reimbursed });
        }
        await apiCreateReceipt(api, {
          groupId,
          paidByUserId: adminId,
          name: uniqueName(receipt.name),
          customFields,
        });
      }
    });
  });

  // Each id is checked because a seeding failure leaves the later ones unassigned,
  // and deleting `undefined` turns one real error into a wall of cleanup noise.
  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        if (templateId) {
          await apiDeleteReportTemplateById(api, templateId);
        }
        // The group takes its receipts (and their custom field values) with it, so
        // the fields are safe to drop afterwards.
        if (groupId) {
          await apiDeleteGroupById(api, String(groupId));
        }
        for (const id of [tipId, dueId, reimbursedId]) {
          if (id) {
            await apiDeleteCustomFieldById(api, id);
          }
        }
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  /**
   * Open the builder on the seeded group over a window that contains the receipts
   * (they are dated 2024-01-01), so nothing depends on the wall clock.
   */
  async function openBuilderOnSeededData(page: Page): Promise<void> {
    await gotoReportBuilder(page);
    await addGroupToScopeByName(page, groupName);

    await openComboboxAndPick(
      page,
      page.getByRole('combobox', { name: /Period covering/ }),
      page.getByRole('option', { name: /Custom range/ }),
    );
    await page.getByLabel('Start', { exact: true }).fill('01/01/2023');
    await page.getByLabel('End', { exact: true }).fill('12/31/2024');
    await page.getByLabel('End', { exact: true }).blur();

    await expect(page.getByTestId('report-receipt-count')).toContainText('4 receipts', {
      timeout: 20_000,
    });
  }

  /** The preview is server-rendered HTML in an iframe; assert against its srcdoc. */
  function preview(page: Page) {
    return page.getByTitle('Report preview');
  }

  test('offers every custom field in the grouping picker, badged', async ({ page }) => {
    await openBuilderOnSeededData(page);

    await page.getByRole('combobox', { name: /Add grouping level/ }).click();

    // Every type is groupable — a CURRENCY field is a measure, but measuring is the
    // only thing its type restricts.
    for (const name of [tipName, dueName, reimbursedName]) {
      await expect(page.getByRole('option', { name: new RegExp(`^${name} Custom$`) })).toBeVisible();
    }

    // A DATE field also contributes the derived calendar-period levels.
    for (const period of ['Day', 'Month', 'Year']) {
      await expect(
        page.getByRole('option', { name: new RegExp(`^${dueName} \\(${period}\\) Custom$`) }),
      ).toBeVisible();
    }

    // The badge marks custom fields only — a built-in carries none.
    await expect(page.getByRole('option', { name: 'Paid By', exact: true })).toBeVisible();
  });

  test('groups by a currency custom field and formats its buckets as money', async ({ page }) => {
    await openBuilderOnSeededData(page);
    await addGroupingLevel(page, new RegExp(`^${tipName} Custom$`));

    // The picked level keeps the badge, so a custom field stays identifiable after
    // it has been chosen.
    await expect(page.getByTestId('report-grouping-label')).toHaveText([tipName]);
    await expect(page.getByTestId('report-grouping-custom')).toHaveCount(1);

    // Money is formatted per the app's currency configuration, whose symbol and
    // separators are global settings this spec must not mutate — so assert the
    // thousands separator and the two decimal places rather than a symbol. The
    // engine's raw rendering of 1500.5 has neither, so this fails if the
    // type-aware formatting regresses.
    await expect(preview(page)).toHaveAttribute('srcdoc', /1[,. ]500[.,]50/, { timeout: 20_000 });
    await expect(preview(page)).toHaveAttribute('srcdoc', /2[.,]00/);
    // A receipt carrying no value for the field is its own bucket, never dropped.
    await expect(preview(page)).toHaveAttribute('srcdoc', /\(None\)/);
  });

  test('renders a boolean custom field as Yes/No', async ({ page }) => {
    await openBuilderOnSeededData(page);
    await addGroupingLevel(page, new RegExp(`^${reimbursedName} Custom$`));

    // Assert through the retrying matcher first: the preview is debounced, so a
    // plain getAttribute can still be holding the pre-grouping render.
    await expect(preview(page)).toHaveAttribute('srcdoc', />\s*Yes\s*</, { timeout: 20_000 });
    await expect(preview(page)).toHaveAttribute('srcdoc', />\s*No\s*</);
    await expect(preview(page)).toHaveAttribute('srcdoc', /\(None\)/);

    // The grouped preview is on screen now, so one read is safe for the negatives —
    // and they must be one-shot: a `not.toHaveAttribute` would pass simply by
    // catching a stale render that hasn't got there yet.
    const srcdoc = await preview(page).getAttribute('srcdoc');
    expect(srcdoc).not.toMatch(/>\s*true\s*</);
    expect(srcdoc).not.toMatch(/>\s*false\s*</);
  });

  test('buckets a date custom field by calendar month', async ({ page }) => {
    await openBuilderOnSeededData(page);
    await addGroupingLevel(page, new RegExp(`^${dueName} \\(Month\\) Custom$`));

    // Two receipts fall in March and one in April, so grouping by the derived month
    // yields three buckets — not the one-per-receipt the raw instant would give.
    await expect(preview(page)).toHaveAttribute('srcdoc', /2024-03/, { timeout: 20_000 });
    await expect(preview(page)).toHaveAttribute('srcdoc', /2024-04/);
    await expect(preview(page)).toHaveAttribute('srcdoc', /\(None\)/);

    // The raw timestamp never reaches the reader (see above on one-shot negatives).
    const srcdoc = await preview(page).getAttribute('srcdoc');
    expect(srcdoc).not.toContain('2024-03-15T');
  });

  test('aggregates a custom currency field as a measure', async ({ page }) => {
    await openBuilderOnSeededData(page);

    await page.getByTestId('report-add-column').click();
    await page.getByTestId('picker-kind-aggregate').click();

    // A currency custom field is offered as a measure, badged there too.
    const measure = page.getByRole('combobox', { name: 'Measure' });
    const tipOption = page.getByRole('option', { name: new RegExp(`^${tipName} Custom$`) });
    await openComboboxAndPick(page, measure, tipOption);

    await Promise.all([waitForPreview(page), page.getByTestId('picker-save').click()]);

    // The column row carries the Custom badge because the measure it reads is custom.
    const column = page.getByTestId('report-column-row').filter({ hasText: tipName });
    await expect(column).toBeVisible();
    await expect(column.getByTestId('report-column-custom')).toBeVisible();

    // 1500.50 + 1500.50 + 2.00, with the fourth receipt contributing nothing.
    await expect(preview(page)).toHaveAttribute('srcdoc', /3[,. ]003[.,]00/, { timeout: 20_000 });
  });

  test('a saved template names its custom fields and reopens badged', async ({ page }) => {
    const templateName = uniqueName('cf-template');
    await withAdminApi(async (api) => {
      templateId = (
        await apiCreateReportTemplate(api, {
          name: templateName,
          groupIds: [String(groupId)],
          groupBy: [`custom_${dueId}_month`],
          detail: { mode: 'aggregate', by: `custom_${reimbursedId}` },
          columns: [
            {
              kind: 'dimension',
              name: 'Reimbursed',
              label: 'Reimbursed',
              field: `custom_${reimbursedId}`,
            },
          ],
        })
      ).id;
    });

    await gotoReports(page);
    const row = page.getByRole('row').filter({ hasText: templateName });

    // The list resolves custom_<id> to the field's name; the raw key never shows.
    await expect(row).toContainText(`${dueName} (Month)`);
    await expect(row).toContainText(`Aggregate by ${reimbursedName}`);
    await expect(row).not.toContainText(`custom_${dueId}`);

    // Reopening rehydrates the grouping level, still marked custom.
    await row.getByTestId('report-template-edit').click();
    await page.waitForURL(/\/reports\/\d+\/edit/);
    await expect(page.getByTestId('report-grouping-label')).toHaveText([`${dueName} (Month)`]);
    await expect(page.getByTestId('report-grouping-custom')).toHaveCount(1);
  });
});
