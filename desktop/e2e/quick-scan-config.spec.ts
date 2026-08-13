import { expect, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiCreateGroup,
  apiDeleteGroupById,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';

// Drives the group receipt settings "Quick Scan" config UI as an admin (group owner → group.update).
// Asserts the section renders, the conditional default controls appear/disappear per the Show/Require
// toggles, and that the form's validation gates Save when an optional field has no configured default.
// It does NOT assert the persisted round-trip — the config page save is exercised for its UI behavior
// (backend persistence is covered by the API repo tests).

test.use({ storageState: 'e2e/.auth/admin.json' });

test.describe('Quick scan group receipt settings', () => {
  test.describe.configure({ mode: 'serial' });

  let groupId: string;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      const group = await apiCreateGroup(api, uniqueName('qs-cfg'));
      groupId = String(group.id);
    });
  });

  test.afterAll(async () => {
    await withAdminApi(async (api) => {
      await apiDeleteGroupById(api, groupId);
    });
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
    await page.goto(`/groups/${groupId}/receipt-settings/edit`);
    // The Quick Scan section has loaded once its (uniquely identified) controls are present.
    await expect(page.getByTestId('quick-scan-paidby-show').getByRole('checkbox')).toBeVisible();
  });

  test('reveals the paid-by/status defaults only when the field is optional', async ({ page }) => {
    const paidByShow = page.getByTestId('quick-scan-paidby-show').getByRole('checkbox');
    const paidByRequire = page.getByTestId('quick-scan-paidby-require').getByRole('checkbox');
    const defaultPaidBy = page.getByRole('combobox', { name: 'Default Paid By', exact: true });
    const defaultPaidByUser = page.getByRole('combobox', { name: 'Default Paid By User' });

    // Shown + required → no default needed → default controls hidden.
    await paidByShow.check();
    await paidByRequire.check();
    await expect(defaultPaidBy).toHaveCount(0);

    // Optional → the default paid-by selector appears.
    await paidByRequire.uncheck();
    await expect(defaultPaidBy).toBeVisible();

    // "Specific user" reveals the member autocomplete; "Uploader" hides it again.
    await defaultPaidBy.click();
    await page.getByRole('option', { name: 'Specific user', exact: true }).click();
    await expect(defaultPaidByUser).toBeVisible();

    await defaultPaidBy.click();
    await page.getByRole('option', { name: 'Uploader', exact: true }).click();
    await expect(defaultPaidByUser).toHaveCount(0);

    // Status behaves the same way.
    const statusRequire = page.getByTestId('quick-scan-status-require').getByRole('checkbox');
    const defaultStatus = page.getByRole('combobox', { name: 'Default Status' });
    await page.getByTestId('quick-scan-status-show').getByRole('checkbox').check();
    await statusRequire.check();
    await expect(defaultStatus).toHaveCount(0);
    await statusRequire.uncheck();
    await expect(defaultStatus).toBeVisible();
  });

  test('blocks save until an optional paid-by has a default', async ({ page }) => {
    await page.getByTestId('quick-scan-paidby-show').getByRole('checkbox').check();
    await page.getByTestId('quick-scan-paidby-require').getByRole('checkbox').uncheck();

    // "Specific user" with no user chosen → form invalid → Save is a no-op (stays on /edit).
    await page.getByRole('combobox', { name: 'Default Paid By', exact: true }).click();
    await page.getByRole('option', { name: 'Specific user', exact: true }).click();
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page).toHaveURL(/\/receipt-settings\/edit/);

    // Switch to "Uploader" (a complete default) → Save succeeds and navigates to the view page.
    await page.getByRole('combobox', { name: 'Default Paid By', exact: true }).click();
    await page.getByRole('option', { name: 'Uploader', exact: true }).click();
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page).toHaveURL(/\/receipt-settings\/view/);
  });

  test('greys out the comment toggles while comments are hidden, keeping their values', async ({ page }) => {
    const commentShow = page.getByTestId('quick-scan-comment-show').getByRole('checkbox');
    const commentRequire = page.getByTestId('quick-scan-comment-require').getByRole('checkbox');
    const hideComments = page.getByText('Hide Comments', { exact: true });

    await commentShow.check();
    await commentRequire.check();
    await expect(commentShow).toBeEnabled();

    // Hiding comments group-wide hides the quick-scan comment too — the toggles grey out but keep
    // their values (derived state, nothing is cleared).
    await hideComments.click();
    await expect(commentShow).toBeDisabled();
    await expect(commentRequire).toBeDisabled();
    await expect(commentShow).toBeChecked();
    await expect(commentRequire).toBeChecked();

    // Un-hiding restores them exactly as configured.
    await hideComments.click();
    await expect(commentShow).toBeEnabled();
    await expect(commentShow).toBeChecked();
    await expect(commentRequire).toBeChecked();
  });

  test('persists the comment configuration', async ({ page }) => {
    await page.getByTestId('quick-scan-comment-show').getByRole('checkbox').check();
    await page.getByTestId('quick-scan-comment-require').getByRole('checkbox').check();
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page).toHaveURL(/\/receipt-settings\/view/);

    // Reload the edit page: the toggles must come back checked. A field missing from the API's
    // update assignment block persists nothing and fails only here.
    await page.goto(`/groups/${groupId}/receipt-settings/edit`);
    await expect(page.getByTestId('quick-scan-comment-show').getByRole('checkbox')).toBeChecked();
    await expect(page.getByTestId('quick-scan-comment-require').getByRole('checkbox')).toBeChecked();
  });
});
