import { expect, test } from '@playwright/test';
import { creds, loginViaUi } from './helpers/auth';

// These tests exercise the login UI itself, so they must start unauthenticated
// even though the chromium project defaults to a logged-in storage state.
test.use({ storageState: { cookies: [], origins: [] } });

test.describe('authentication', () => {
  test('admin can log in and reach the dashboard', async ({ page }) => {
    await loginViaUi(page, 'admin');
    await expect(page).toHaveURL(/\/dashboard\/group\/\d+/);
  });

  test('regular user can log in and reach the dashboard', async ({ page }) => {
    await loginViaUi(page, 'user');
    await expect(page).toHaveURL(/\/dashboard\/group\/\d+/);
  });

  // The only test covering both halves of the failed-login fix on the wire: the
  // API returns a bad password as a 500 carrying errorMsg, which the interceptor
  // used to overwrite with a generic message, behind a spinner that never left.
  test('a wrong password reports it and leaves the form filled in', async ({
    page,
  }) => {
    const { username } = creds('user');
    await page.goto('/auth/login');
    await page.getByLabel('Username').fill(username);
    await page.getByLabel('Password').fill('definitely-not-the-password');
    await page.getByRole('button', { name: 'Login' }).click();

    // Assert the toast first — it auto-dismisses after 3s.
    await expect(page.getByText('Invalid credentials.')).toBeVisible();
    await expect(page.getByText(/Http failure response/)).toHaveCount(0);

    // Still on the login form, with what was typed.
    await expect(page).toHaveURL(/\/auth\/login/);
    await expect(page.locator('.loading-container')).toHaveCount(0);
    await expect(page.getByLabel('Username')).toHaveValue(username);
    await expect(page.getByLabel('Password')).toHaveValue(
      'definitely-not-the-password',
    );
  });
});
