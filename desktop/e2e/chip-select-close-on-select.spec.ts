import { expect, test, type Page } from '@playwright/test';
import { rmSync } from 'node:fs';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiCreateCategory,
  apiDeleteCategoryById,
  apiDeleteUserByName,
  createUserWithRole,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';

// The "close the option list after each selection" user preference, end to end.
//
// The Jest specs drive AutocomleteComponent against a mocked store, so they
// prove the branch but nothing about the wire. Only an e2e covers the whole
// chain: the checkbox on the User Preferences form -> PUT /userPreferences ->
// the stored column -> AppData on the next navigation -> AuthState -> the
// behavior of a chip picker in a completely different feature.
//
// It provisions its OWN account. The preference is per-user and global, and the
// suite runs fullyParallel, so flipping it on a shared e2e account would change
// behavior under every other spec running as that account at the same time.

const NEW_USER_PASSWORD = `${uniqueName('pw')}-Aa1!`;
const AUTH_FILE = 'e2e/.auth/chip-select-user.json';

test.describe('Close chip select on select preference', () => {
  // Serial: the second test changes stored state the first one depends on.
  test.describe.configure({ mode: 'serial' });

  test.use({ storageState: AUTH_FILE });

  let username: string;
  let categoryA: { id: number; name: string };
  let categoryB: { id: number; name: string };

  test.beforeAll(async ({ browser }) => {
    username = uniqueName('chip-select-user');

    const admin = await browser.newContext({ storageState: 'e2e/.auth/admin.json' });
    const adminPage = await admin.newPage();
    await stubTokenRefresh(adminPage);
    await createUserWithRole(adminPage, {
      username,
      password: NEW_USER_PASSWORD,
      role: 'Legacy User',
    });
    await admin.close();

    // Log in once and persist the session (AUTH_FILE does not exist yet, so
    // this context opts out of it).
    const userContext = await browser.newContext({ storageState: undefined });
    const userPage = await userContext.newPage();
    await userPage.goto('/auth/login');
    await userPage.getByLabel('Username').fill(username);
    await userPage.getByLabel('Password').fill(NEW_USER_PASSWORD);
    await userPage.getByRole('button', { name: 'Login' }).click();
    await expect(userPage).toHaveURL(/\/dashboard\/group\/\d+/, { timeout: 15_000 });
    await userContext.storageState({ path: AUTH_FILE });
    await userContext.close();

    // Categories are a global pool, so seeding two here guarantees the picker
    // has options regardless of what the rest of the suite has created. This
    // has to come AFTER the login above: withAdminApi opens a request context,
    // which inherits the describe's storageState, and AUTH_FILE does not exist
    // until that login writes it.
    await withAdminApi(async (api) => {
      categoryA = await apiCreateCategory(api, uniqueName('chipsel-cata'));
      categoryB = await apiCreateCategory(api, uniqueName('chipsel-catb'));
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        await apiDeleteUserByName(api, username);
        await apiDeleteCategoryById(api, categoryA.id);
        await apiDeleteCategoryById(api, categoryB.id);
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a teardown error.
    }
    rmSync(AUTH_FILE, { force: true });
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  /**
   * Opens the Categories picker on a fresh add-receipt form and picks the first
   * option. Navigating in rather than reusing a page is deliberate: AppData is
   * refetched per navigation, so the preference is read off the server on every
   * call rather than from whatever the previous test left in the store.
   */
  async function pickFirstCategory(page: Page) {
    await page.goto('/receipts/add');
    await expect(page.getByLabel('Name')).toBeVisible();

    // getByRole, not getByLabel: once the panel opens, the overlay's listbox is
    // labelled by the same mat-label, so getByLabel('Categories') matches two
    // elements and every later assertion dies on strict mode.
    const categories = page.getByRole('combobox', { name: 'Categories' });
    await categories.click();
    await expect(page.getByRole('listbox')).toBeVisible();
    await page.getByRole('option').first().click();

    return categories;
  }

  test('leaves the option list open by default', async ({ page }) => {
    const categories = await pickFirstCategory(page);

    // A chip is what says the pick landed; without it the panel assertions
    // below would be describing a click that did nothing.
    await expect(page.getByTestId('receipt-categories').locator('mat-chip-row')).toHaveCount(1);
    await expect(categories).toBeFocused();
    await expect(page.getByRole('listbox')).toBeVisible();
  });

  test('closes the option list once the preference is on', async ({ page }) => {
    await page.goto('/settings/user-preferences/edit');
    const checkbox = page.getByLabel('Close the option list after each selection?');
    await expect(checkbox).toBeVisible();
    await checkbox.check();
    await page.getByRole('button', { name: 'Save', exact: true }).first().click();
    await expect(page).toHaveURL(/\/settings\/user-preferences\/view/);

    const categories = await pickFirstCategory(page);

    await expect(page.getByTestId('receipt-categories').locator('mat-chip-row')).toHaveCount(1);
    // Assert the blur BEFORE the panel: Material closes the panel on selection
    // by itself, so "hidden" alone would pass even with the preference ignored.
    // Losing focus is what only the close-on-select path produces.
    await expect(categories).not.toBeFocused();
    await expect(page.getByRole('listbox')).toBeHidden();
  });
});
