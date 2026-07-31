import { expect, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import { withAdminApi } from './helpers/provisioning';

// A clean URL (no spaces / reserved chars) so Go's url.QueryEscape and JS's
// encode/decodeURIComponent agree byte-for-byte on the fragment.
const SERVER_URL = 'https://e2e-mobile.example.com/api';
const QR_IMG_NAME = 'Scan to set up the Receipt Wrangler mobile app';

// `showLoginQr` is a GLOBAL system setting, so these tests mutate shared server
// state and must run serially. afterAll force-reverts it so the QR can't leak
// onto the parallel suite's login page.
test.describe.serial('Login QR (System Settings → login page)', () => {
  let original: Record<string, unknown> | undefined;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      original = await (await api.get('/api/systemSettings')).json();
    });
  });

  test.afterAll(async () => {
    if (!original) {
      return;
    }
    // Restore the full settings object (it's an upsert with required currency /
    // taskConcurrency fields), forcing the login QR back off.
    await withAdminApi(async (api) => {
      await api.put('/api/systemSettings', {
        data: { ...original, showLoginQr: false, mobileServerUrl: '' },
      });
    });
  });

  test('admin enables the login QR in System Settings and it persists', async ({
    browser,
  }) => {
    const context = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const page = await context.newPage();
    await stubTokenRefresh(page);

    await page.goto('/system-settings/settings/edit');
    // Enabling the toggle makes the URL required, so fill it before saving.
    await page.getByLabel('Show login QR code').check();
    await page.getByLabel('Mobile Server URL').fill(SERVER_URL);
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page).toHaveURL(/\/system-settings\/settings\/view/);

    await context.close();

    // The setting persisted server-side.
    await withAdminApi(async (api) => {
      const settings = await (await api.get('/api/systemSettings')).json();
      expect(settings.showLoginQr).toBe(true);
      expect(settings.mobileServerUrl).toBe(SERVER_URL);
    });
  });

  test('the QR shows on the login page and encodes the server URL', async ({
    browser,
  }) => {
    // The login page is pre-auth: use a fresh context with no session so the
    // app fetches GET /featureConfig (rather than riding on appData).
    const context = await browser.newContext({
      storageState: { cookies: [], origins: [] },
    });
    const page = await context.newPage();

    const featureConfig = page.waitForResponse((res) =>
      res.url().includes('/api/featureConfig'),
    );
    await page.goto('/auth/login');

    // The backend composed the deep link, and the encoded fragment round-trips
    // back to the exact server URL an admin set.
    const body = await (await featureConfig).json();
    expect(body.loginQrUrl).toMatch(
      /^https:\/\/receiptwrangler\.io\/app\/setup#url=/,
    );
    const encoded = String(body.loginQrUrl).split('#url=')[1];
    expect(decodeURIComponent(encoded)).toBe(SERVER_URL);

    // And the QR image actually renders (generated client-side from loginQrUrl).
    await expect(page.getByRole('img', { name: QR_IMG_NAME })).toBeVisible();
    await expect(page.getByText('Set up the mobile app')).toBeVisible();

    await context.close();
  });

  test('disabling the toggle hides the QR on the login page', async ({
    browser,
  }) => {
    // Turn it off through the admin UI.
    const adminContext = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const adminPage = await adminContext.newPage();
    await stubTokenRefresh(adminPage);
    await adminPage.goto('/system-settings/settings/edit');
    await adminPage.getByLabel('Show login QR code').uncheck();
    await adminPage.getByRole('button', { name: 'Save' }).click();
    await expect(adminPage).toHaveURL(/\/system-settings\/settings\/view/);
    await adminContext.close();

    // The login page no longer renders the QR.
    const context = await browser.newContext({
      storageState: { cookies: [], origins: [] },
    });
    const page = await context.newPage();
    await page.goto('/auth/login');
    await expect(page.getByRole('img', { name: QR_IMG_NAME })).toHaveCount(0);

    await context.close();
  });
});
