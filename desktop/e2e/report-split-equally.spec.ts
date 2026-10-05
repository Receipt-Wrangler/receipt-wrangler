import { expect, Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateCategory,
  apiCreateCustomField,
  apiCreateGroup,
  apiCreateReceipt,
  apiCreateTag,
  apiDeleteCategoryById,
  apiDeleteCustomFieldById,
  apiDeleteGroupById,
  apiDeleteReportTemplateById,
  apiDeleteTagById,
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

// The builder is gated by app.reports.read, and seeding categories, tags and a
// custom field needs their create permissions — all Legacy Admin.
test.use({ storageState: 'e2e/.auth/admin.json' });

// Serial: every test shares one seeded group, and categories, tags and custom
// fields are global pools, so concurrent runs would multiply the seeding.
test.describe.configure({ mode: 'serial' });

/** The preview is server-rendered HTML in an iframe; assert against its srcdoc. */
function preview(page: Page) {
  return page.getByTitle('Report preview');
}

/**
 * Money, tolerant of the decimal separator: the currency configuration is a
 * global System Setting on the shared CI backend that this spec must not mutate.
 */
function money(whole: string, cents: string): RegExp {
  return new RegExp(`${whole}[.,]${cents}`);
}

/**
 * The split checkboxes on the Report Builder. A receipt with three categories
 * normally counts in full in every category; ticked, each category gets a third,
 * with the leftover cent going to the first, so the buckets add back up to what
 * was spent.
 *
 * The spec seeds its own group so the preview holds exactly these two receipts:
 *   - 100.00 in Food, Fuel and Toys, tagged Alex and Sam
 *   - 20.00 in Food only
 * Unsplit, the category grand total is 320.00 (100 three times, plus 20). Split,
 * it is the true 120.00: Food 53.34, Fuel 33.33, Toys 33.33.
 */
test.describe('Report Builder — split amounts equally', () => {
  const groupName = uniqueName('split-grp');
  const foodName = uniqueName('Food');
  const fuelName = uniqueName('Fuel');
  const toysName = uniqueName('Toys');
  const alexName = uniqueName('Alex');
  const samName = uniqueName('Sam');
  const tipName = uniqueName('Tip');

  let groupId: number;
  let categoryIds: number[] = [];
  let tagIds: number[] = [];
  let tipId: number;
  let templateId: number | undefined;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      const adminId = await apiGetUserId(api, creds('admin').username);
      groupId = (await apiCreateGroup(api, groupName)).id;

      const food = await apiCreateCategory(api, foodName);
      const fuel = await apiCreateCategory(api, fuelName);
      const toys = await apiCreateCategory(api, toysName);
      categoryIds = [food.id, fuel.id, toys.id];

      const alex = await apiCreateTag(api, alexName);
      const sam = await apiCreateTag(api, samName);
      tagIds = [alex.id, sam.id];

      tipId = (await apiCreateCustomField(api, { name: tipName, type: 'CURRENCY' })).id;

      await apiCreateReceipt(api, {
        groupId,
        paidByUserId: adminId,
        name: uniqueName('split-shared'),
        amount: '100.00',
        categories: [food, fuel, toys],
        tags: [alex, sam],
        customFields: [{ customFieldId: tipId, currencyValue: '10.00' }],
      });
      await apiCreateReceipt(api, {
        groupId,
        paidByUserId: adminId,
        name: uniqueName('split-food'),
        amount: '20.00',
        categories: [food],
      });
    });
  });

  // Each id is checked because a seeding failure leaves the later ones unassigned.
  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        if (templateId) {
          await apiDeleteReportTemplateById(api, templateId);
        }
        // The group takes its receipts with it, so the pools are safe to drop after.
        if (groupId) {
          await apiDeleteGroupById(api, String(groupId));
        }
        for (const id of categoryIds) {
          await apiDeleteCategoryById(api, id);
        }
        for (const id of tagIds) {
          await apiDeleteTagById(api, id);
        }
        if (tipId) {
          await apiDeleteCustomFieldById(api, tipId);
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
   * Open the builder on the seeded group over a window containing the receipts
   * (they are dated 2024-01-01), so nothing depends on the wall clock. The
   * builder's default report aggregates by category with Count and Total columns
   * and a grand total.
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

    await expect(page.getByTestId('report-receipt-count')).toContainText('2 receipts', {
      timeout: 20_000,
    });
  }

  function splitCategories(page: Page) {
    return page.getByRole('checkbox', { name: 'Split amounts equally across categories' });
  }

  function splitTags(page: Page) {
    return page.getByRole('checkbox', { name: 'Split amounts equally across tags' });
  }

  test('splits each receipt equally across its categories', async ({ page }) => {
    await openBuilderOnSeededData(page);

    // Unsplit, the three-category receipt counts in full three times.
    await expect(preview(page)).toHaveAttribute('srcdoc', money('320', '00'), { timeout: 20_000 });

    await Promise.all([waitForPreview(page), splitCategories(page).check()]);

    await expect(preview(page)).toHaveAttribute('srcdoc', money('53', '34'), { timeout: 20_000 });
    const html = (await preview(page).getAttribute('srcdoc')) ?? '';
    expect(html).toMatch(money('33', '33'));
    expect(html).toMatch(money('120', '00'));
    expect(html).not.toMatch(money('320', '00'));
    // The count is still receipts, not the shares they were split into.
    await expect(page.getByTestId('report-receipt-count')).toContainText('2 receipts');
    // Aggregating by category, the split applies, so no note says otherwise.
    await expect(page.getByTestId('report-split-categories-inactive')).toHaveCount(0);
  });

  test('flags a tag split until the report groups by tag', async ({ page }) => {
    await openBuilderOnSeededData(page);

    await Promise.all([waitForPreview(page), splitTags(page).check()]);
    await expect(page.getByTestId('report-split-tags-inactive')).toBeVisible();

    await addGroupingLevel(page, 'Tag');
    await expect(page.getByTestId('report-split-tags-inactive')).toHaveCount(0);

    // Grouped by tag, the 100.00 receipt halves across Alex and Sam. Categories are
    // not split here, so each half still counts once per category: Alex is
    // 3 x 50.00 = 150.00, where an unsplit tag would have given 300.00.
    await expect(preview(page)).toHaveAttribute('srcdoc', money('150', '00'), { timeout: 20_000 });
    expect(await preview(page).getAttribute('srcdoc')).not.toMatch(money('300', '00'));
  });

  test('offers the currency fields to leave whole, and sends the pick', async ({ page }) => {
    await openBuilderOnSeededData(page);

    // Nothing to exclude until a split is ticked.
    await expect(page.getByTestId('report-split-excluded')).toHaveCount(0);
    await Promise.all([waitForPreview(page), splitCategories(page).check()]);

    const excluded = page.getByTestId('report-split-excluded');
    await expect(excluded).toBeVisible();

    const [request] = await Promise.all([
      waitForPreview(page).then((response) => response.request()),
      (async () => {
        // fill, not click: the empty field's floating label sits over the input
        // and intercepts a click. Typing opens the panel just the same.
        await excluded.getByRole('combobox').fill(tipName);
        await page.getByRole('option', { name: tipName, exact: true }).click();
      })(),
    ]);
    expect(request.postDataJSON()).toMatchObject({
      splitCategoriesEqually: true,
      splitExcludedFields: [`custom_${tipId}`],
    });
  });

  test('saves the split options into a template and reopens them', async ({ page }) => {
    await openBuilderOnSeededData(page);
    await Promise.all([waitForPreview(page), splitCategories(page).check()]);

    const templateName = uniqueName('split-template');
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
    const saved = (await response.json()) as { id: number; configuration: Record<string, unknown> };
    templateId = saved.id;
    expect(saved.configuration).toMatchObject({ splitCategoriesEqually: true });
    expect(saved.configuration).not.toHaveProperty('splitTagsEqually');

    await gotoReports(page);
    await page
      .getByRole('row')
      .filter({ hasText: templateName })
      .getByTestId('report-template-name')
      .click();
    await expect(page).toHaveURL(/\/reports\/\d+\/edit$/);

    await expect(splitCategories(page)).toBeChecked();
    await expect(splitTags(page)).not.toBeChecked();
    await expect(preview(page)).toHaveAttribute('srcdoc', money('53', '34'), { timeout: 20_000 });
  });
});
