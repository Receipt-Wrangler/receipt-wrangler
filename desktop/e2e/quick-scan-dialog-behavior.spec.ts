import { expect, type Route, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateGroup,
  apiDeleteGroupById,
  apiGetUserId,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import {
  injectQuickScanAppData,
  openQuickScanDialog,
  parseMultipartFields,
  selectImageGroup,
  uploadQuickScanImages,
} from './helpers/quick-scan';

// Deeper Quick Scan DIALOG behavior, complementing quick-scan-dialog.spec.ts
// (which asserts a single static config snapshot). Everything is driven by
// client-side AppData injection (see helpers/quick-scan.ts) so no server config
// is mutated. The two SUBMIT specs mock POST /receipt/quickScan: the backend
// validates each group's PERSISTED config (which we intentionally don't touch),
// so a real submit would 400 — capturing the request instead lets us assert the
// exact multipart the client builds (the "falls off the submission" half the
// mobile suite can't observe, since its queued receipt has no id).

test.use({ storageState: 'e2e/.auth/admin.json' });

test.describe('Quick scan dialog behavior', () => {
  let groupA: { id: number; name: string };
  let groupB: { id: number; name: string };
  let adminId: number;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      groupA = await apiCreateGroup(api, uniqueName('qs-a'));
      groupB = await apiCreateGroup(api, uniqueName('qs-b'));
      adminId = await apiGetUserId(api, creds('admin').username);
    });
  });

  test.afterAll(async () => {
    await withAdminApi(async (api) => {
      await apiDeleteGroupById(api, String(groupA.id));
      await apiDeleteGroupById(api, String(groupB.id));
    });
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  // Gap #1 — a user-preference paid-by preset "falls off the FORM" when the
  // selected group hides paid-by, while a preset for a shown field is honored.
  test('a preset paid-by falls off the form when the group hides it; a preset status is kept', async ({
    page,
  }) => {
    await injectQuickScanAppData(page, {
      groupConfigs: [
        {
          groupId: groupA.id,
          config: {
            quickScanPaidByEnabled: false, // hidden -> preset must fall off
            quickScanDefaultPaidByType: 'UPLOADER',
            quickScanStatusEnabled: true, // shown -> preset honored
            quickScanStatusRequired: false,
            quickScanDefaultStatus: 'OPEN',
            quickScanCategoriesEnabled: false,
            quickScanTagsEnabled: false,
          },
        },
      ],
      // The caller has a preset group (A), paid-by (admin) and status (OPEN).
      userPreferences: {
        quickScanDefaultGroupId: groupA.id,
        quickScanDefaultPaidById: adminId,
        quickScanDefaultStatus: 'OPEN',
      },
    });

    const dialog = await openQuickScanDialog(page, groupA.id);
    await uploadQuickScanImages(dialog, 1);
    // Group A is preselected from the prefill, so its config drives the fields.
    await expect(dialog.getByRole('combobox', { name: 'Group' })).toBeVisible();

    // Paid-by is hidden by group A -> the preset admin fell off (field absent).
    await expect(dialog.getByRole('combobox', { name: 'Paid By' })).toHaveCount(0);

    // Status is shown -> the preset OPEN is prefilled (trigger reads "Open").
    const status = dialog.getByRole('combobox', { name: 'Status' });
    await expect(status).toBeVisible();
    await expect(status).toContainText('Open');
  });

  // Gap #4 — the same preset paid-by "falls off the SUBMISSION": it is sent as
  // the empty sentinel (not the stale admin id), while the shown status is sent.
  test('a preset paid-by falls off the submission (sent as the empty sentinel)', async ({
    page,
  }) => {
    await injectQuickScanAppData(page, {
      groupConfigs: [
        {
          groupId: groupA.id,
          config: {
            quickScanPaidByEnabled: false,
            quickScanDefaultPaidByType: 'UPLOADER',
            quickScanStatusEnabled: true,
            quickScanStatusRequired: true,
            quickScanDefaultStatus: 'OPEN',
            quickScanCategoriesEnabled: false,
            quickScanTagsEnabled: false,
          },
        },
      ],
      userPreferences: {
        quickScanDefaultGroupId: groupA.id,
        quickScanDefaultPaidById: adminId,
        quickScanDefaultStatus: 'OPEN',
      },
    });

    const requests: { body: Buffer | null; contentType: string }[] = [];
    await page.route('**/api/receipt/quickScan', async (route: Route) => {
      const req = route.request();
      requests.push({
        body: req.postDataBuffer(),
        contentType: req.headers()['content-type'] ?? '',
      });
      await route.fulfill({ status: 200, contentType: 'application/json', body: '{}' });
    });

    const dialog = await openQuickScanDialog(page, groupA.id);
    await uploadQuickScanImages(dialog, 1);
    await expect(dialog.getByRole('combobox', { name: 'Paid By' })).toHaveCount(0);

    await dialog.getByTestId('dialog-submit-button').click();
    await expect(page.getByText('Successfully queued', { exact: false })).toBeVisible();
    await expect(dialog).toBeHidden();

    expect(requests).toHaveLength(1);
    expect(requests[0].body).not.toBeNull();
    const fields = parseMultipartFields(requests[0].body!, requests[0].contentType);
    expect(fields.get('groupIds')).toEqual([String(groupA.id)]);
    expect(fields.get('paidByUserIds')).toEqual(['']); // preset admin discarded
    expect(fields.get('statuses')).toEqual(['OPEN']); // preset status kept
  });

  // Gap #3 — changing an image's group re-runs configureImages and flips which
  // fields are shown/required for that image.
  test("switching an image's group re-flips which fields are shown", async ({ page }) => {
    await injectQuickScanAppData(page, {
      groupConfigs: [
        {
          groupId: groupA.id,
          config: {
            quickScanPaidByEnabled: false, // A: paid-by hidden, status shown
            quickScanStatusEnabled: true,
            quickScanStatusRequired: false,
            quickScanDefaultStatus: 'OPEN',
            quickScanCategoriesEnabled: false,
            quickScanTagsEnabled: false,
          },
        },
        {
          groupId: groupB.id,
          config: {
            quickScanPaidByEnabled: true, // B: paid-by shown, status hidden
            quickScanStatusEnabled: false,
            quickScanDefaultStatus: 'OPEN',
            quickScanCategoriesEnabled: false,
            quickScanTagsEnabled: false,
          },
        },
      ],
    });

    const dialog = await openQuickScanDialog(page, groupA.id);
    await uploadQuickScanImages(dialog, 1);

    // Group A: paid-by hidden, status shown.
    await selectImageGroup(page, dialog, groupA.name);
    await expect(dialog.getByRole('combobox', { name: 'Status' })).toBeVisible();
    await expect(dialog.getByRole('combobox', { name: 'Paid By' })).toHaveCount(0);

    // Switch to group B: the field set flips — paid-by shown, status hidden.
    await selectImageGroup(page, dialog, groupB.name);
    await expect(dialog.getByRole('combobox', { name: 'Paid By' })).toBeVisible();
    await expect(dialog.getByRole('combobox', { name: 'Status' })).toHaveCount(0);
  });

  // Gap #5 — the positive category path: a required category picked from the
  // per-group catalog lets the submit through and rides the multipart.
  test('a required category selected via the picker lets the submit through', async ({
    page,
  }) => {
    const category = { id: 987654, name: 'E2E QS Category' };
    await injectQuickScanAppData(page, {
      groupConfigs: [
        {
          groupId: groupA.id,
          config: {
            quickScanPaidByEnabled: false,
            quickScanStatusEnabled: false,
            quickScanDefaultStatus: 'OPEN',
            quickScanCategoriesEnabled: true, // shown + required
            quickScanCategoriesRequired: true,
            quickScanTagsEnabled: false,
          },
        },
      ],
      groupCategories: { [groupA.id]: [category] },
    });

    const requests: { body: Buffer | null; contentType: string }[] = [];
    await page.route('**/api/receipt/quickScan', async (route: Route) => {
      const req = route.request();
      requests.push({
        body: req.postDataBuffer(),
        contentType: req.headers()['content-type'] ?? '',
      });
      await route.fulfill({ status: 200, contentType: 'application/json', body: '{}' });
    });

    const dialog = await openQuickScanDialog(page, groupA.id);
    await uploadQuickScanImages(dialog, 1);
    await selectImageGroup(page, dialog, groupA.name);

    const categories = dialog.getByRole('combobox', { name: 'Categories' });
    await expect(categories).toBeVisible();
    await categories.click();
    await categories.fill(category.name);
    await page.getByRole('option', { name: category.name, exact: true }).click();

    await dialog.getByTestId('dialog-submit-button').click();
    await expect(page.getByText('Successfully queued', { exact: false })).toBeVisible();
    await expect(dialog).toBeHidden();

    expect(requests).toHaveLength(1);
    expect(requests[0].body).not.toBeNull();
    const fields = parseMultipartFields(requests[0].body!, requests[0].contentType);
    expect(fields.get('groupIds')).toEqual([String(groupA.id)]);
    expect(fields.get('categoryIds')).toEqual([String(category.id)]);
  });

  // Gap #6 — two images on two groups get independent field sets, and one
  // image's unmet required field blocks the WHOLE submit (nothing is sent).
  test("two images on different groups get independent fields; one missing required field blocks submit", async ({
    page,
  }) => {
    await injectQuickScanAppData(page, {
      groupConfigs: [
        {
          groupId: groupA.id,
          config: {
            quickScanPaidByEnabled: false, // image 0 (A): status required
            quickScanStatusEnabled: true,
            quickScanStatusRequired: true,
            quickScanDefaultStatus: 'OPEN',
            quickScanCategoriesEnabled: false,
            quickScanTagsEnabled: false,
          },
        },
        {
          groupId: groupB.id,
          config: {
            quickScanPaidByEnabled: true, // image 1 (B): paid-by required
            quickScanPaidByRequired: true,
            quickScanStatusEnabled: false,
            quickScanDefaultStatus: 'OPEN',
            quickScanCategoriesEnabled: false,
            quickScanTagsEnabled: false,
          },
        },
      ],
    });

    // The form must block this client-side — fail loudly if a POST escapes.
    let posted = false;
    await page.route('**/api/receipt/quickScan', async (route: Route) => {
      posted = true;
      await route.fulfill({ status: 200, contentType: 'application/json', body: '{}' });
    });

    const dialog = await openQuickScanDialog(page, groupA.id);
    await uploadQuickScanImages(dialog, 2);

    // Image 0 -> group A: status shown+required, paid-by hidden. Fill its status
    // so only image 1 is left missing a required field.
    const slide0 = dialog.locator('slide').nth(0);
    await selectImageGroup(page, slide0, groupA.name);
    await expect(slide0.getByRole('combobox', { name: 'Status' })).toBeVisible();
    await expect(slide0.getByRole('combobox', { name: 'Paid By' })).toHaveCount(0);
    await slide0.getByRole('combobox', { name: 'Status' }).click();
    await page.getByRole('option', { name: 'Open', exact: true }).click();

    // Image 1 -> group B: paid-by shown+required, status hidden. Leave paid-by empty.
    await dialog.getByTestId('quick-scan-nav-right').getByRole('button').click();
    const slide1 = dialog.locator('slide').nth(1);
    await expect(slide1.getByRole('combobox', { name: 'Group' })).toBeVisible();
    await selectImageGroup(page, slide1, groupB.name);
    await expect(slide1.getByRole('combobox', { name: 'Paid By' })).toBeVisible();
    await expect(slide1.getByRole('combobox', { name: 'Status' })).toHaveCount(0);

    // Submit is blocked: image 1's required paid-by is empty -> error, no POST.
    await dialog.getByTestId('dialog-submit-button').click();
    await expect(
      page.getByText('Please fill in all required fields', { exact: false }),
    ).toBeVisible();
    await expect(dialog).toBeVisible(); // dialog stays open
    expect(posted).toBe(false);
  });
});
