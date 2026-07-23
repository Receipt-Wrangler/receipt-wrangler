import { expect, Page, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';

// API-key visibility is gated on two app permissions added with the modern
// permission system:
//   - app.api-keys.read-any  -> see EVERY user's keys (vs only your own)
//   - app.api-keys.delete-any -> delete any user's key (vs only your own)
// A Legacy User holds the base create/read/update/delete (own keys only); a
// Legacy Admin additionally holds read-any + delete-any. These specs lock in
// that boundary: a Legacy User can only ever see/manage their own keys, while
// an admin can view all keys and delete one owned by another user.

const API_KEYS_URL = '/settings/api-keys/view';

function uniqueName(tag: string) {
  return `e2e-${tag}-${Date.now()}`;
}

// Switches the table to "All API Keys" via the filter dialog. Only callable
// when the caller holds read-any (otherwise the filter button does not render).
async function filterToAllKeys(page: Page) {
  await page.getByTestId('api-key-filter').click();
  const dialog = page.getByRole('dialog').filter({ hasText: 'Filter API Keys' });
  await expect(dialog).toBeVisible();
  // The option panel is an overlay attached to the page root, not the dialog.
  await dialog.getByRole('combobox', { name: 'API Keys to View' }).click();
  await page.getByRole('option', { name: 'All API Keys', exact: true }).click();
  await dialog.getByTestId('dialog-submit-button').click();
  await expect(dialog).toBeHidden();
}

// Creates an API key via the table-header "add" dialog. The add button is the
// first button in the header, so it's addressed positionally (as
// ensureCustomFieldExists does in receipts.spec.ts).
// Submitting reveals the one-time secret view, which we dismiss with Close —
// that closes with a truthy result, refreshing the table with the new key.
async function createApiKey(page: Page, name: string) {
  await page.goto(API_KEYS_URL);
  await page.locator('app-table-header').getByRole('button').first().click();
  const dialog = page.getByRole('dialog').filter({ hasText: 'Create API Key' });
  await expect(dialog).toBeVisible();
  await dialog.getByLabel('Name').fill(name);
  await dialog.getByTestId('dialog-submit-button').click();
  await dialog.getByRole('button', { name: 'Close' }).click();
  await expect(dialog).toBeHidden();
}

test.describe('API key visibility — Legacy User', () => {
  // Default project storageState = e2e-user (Legacy User).
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('sees only their own keys, with no view-all affordance', async ({
    page,
  }) => {
    await page.goto(API_KEYS_URL);
    // The header is scoped to "My API Keys" and there is no filter button to
    // switch to all keys (filter is gated on app.api-keys.read-any).
    await expect(
      page.getByRole('heading', { name: 'My API Keys' }),
    ).toBeVisible();
    await expect(page.getByTestId('api-key-filter')).toHaveCount(0);
  });
});

test.describe('API key visibility — Admin (read-any)', () => {
  test.use({ storageState: 'e2e/.auth/admin.json' });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test("can filter to view all users' keys", async ({ page }) => {
    await page.goto(API_KEYS_URL);
    await expect(page.getByTestId('api-key-filter')).toBeVisible();
    await filterToAllKeys(page);
    await expect(
      page.getByRole('heading', { name: 'All API Keys' }),
    ).toBeVisible();
  });
});

test.describe('API key ownership and delete-any', () => {
  // Create once (as the owner), then assert the admin can delete it.
  test.describe.configure({ mode: 'serial' });

  // Default project storageState = e2e-user (the owner).
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  let keyName: string;

  test('owner can create a key and sees edit + delete on it', async ({
    page,
  }) => {
    keyName = uniqueName('apikey');
    await createApiKey(page, keyName);

    const row = page.getByRole('row').filter({ hasText: keyName }).first();
    await expect(row).toBeVisible();
    // Owner holds api-keys.update + delete -> both row actions render.
    await expect(row.getByTestId('api-key-edit')).toBeVisible();
    await expect(row.getByTestId('api-key-delete')).toBeVisible();
  });

  test("an admin can view and delete another user's key", async ({
    browser,
  }) => {
    const admin = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const adminPage = await admin.newPage();
    await stubTokenRefresh(adminPage);

    await adminPage.goto(API_KEYS_URL);
    await filterToAllKeys(adminPage);

    const row = adminPage
      .getByRole('row')
      .filter({ hasText: keyName })
      .first();
    // read-any: the admin can see a key owned by e2e-user...
    await expect(row).toBeVisible();
    // ...and delete-any: the delete action renders on a key they don't own.
    await expect(row.getByTestId('api-key-delete')).toBeVisible();
    await row.getByTestId('api-key-delete').click();

    const confirm = adminPage.getByRole('dialog');
    await expect(confirm).toBeVisible();
    await confirm.getByTestId('dialog-submit-button').click();
    await expect(
      adminPage.getByRole('row').filter({ hasText: keyName }),
    ).toHaveCount(0);

    await admin.close();
  });

  // Safety net: if the delete-any test didn't run, an admin removes the orphan.
  test.afterAll(async ({ browser }) => {
    if (!keyName) return;
    const admin = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const adminPage = await admin.newPage();
    try {
      await stubTokenRefresh(adminPage);
      await adminPage.goto(API_KEYS_URL);
      await filterToAllKeys(adminPage);
      const row = adminPage
        .getByRole('row')
        .filter({ hasText: keyName })
        .first();
      if ((await row.count()) > 0) {
        await row.getByTestId('api-key-delete').click();
        await adminPage
          .getByRole('dialog')
          .getByTestId('dialog-submit-button')
          .click();
      }
    } catch {
      // Best-effort cleanup — don't mask a test failure.
    } finally {
      await admin.close();
    }
  });
});
