import { expect, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import { withAdminApi } from './helpers/provisioning';

// A clean URL (no spaces / reserved chars) so Go's url.QueryEscape and JS's
// encode/decodeURIComponent agree byte-for-byte on the fragment.
const SERVER_URL = 'https://e2e-mobile.example.com/api';
const QR_IMG_NAME = 'Scan to set up the Receipt Wrangler mobile app';

// `showLoginQr` / `mobileServerUrl` are GLOBAL system settings, so these tests
// mutate shared server state and must run serially. afterAll puts both fields
// back to whatever they were before the suite -- never a hardcoded "off", which
// would silently disable the QR on an environment that had it enabled.
test.describe.serial('Login QR (System Settings → login page)', () => {
  let originalShowLoginQr: unknown;
  let originalMobileServerUrl: unknown;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      const res = await api.get('/api/systemSettings');
      if (!res.ok()) {
        throw new Error(`GET /api/systemSettings failed: HTTP ${res.status()}`);
      }
      const settings = await res.json();
      originalShowLoginQr = settings.showLoginQr;
      originalMobileServerUrl = settings.mobileServerUrl;
    });
  });

  test.afterAll(async () => {
    // Re-read the live settings and overlay only the two captured fields, so a
    // concurrent change to an unrelated setting isn't clobbered by a stale
    // snapshot (the PUT is an upsert needing the full object). Mirrors
    // mobile/integration_test/helpers/login_qr_fixtures.dart.
    try {
      await withAdminApi(async (api) => {
        // Playwright's request methods resolve on 4xx/5xx, so both calls need
        // an explicit status check or a failed restore passes silently.
        const getResponse = await api.get('/api/systemSettings');
        if (!getResponse.ok()) {
          throw new Error(
            `GET /api/systemSettings failed: HTTP ${getResponse.status()}`,
          );
        }
        const current = await getResponse.json();

        const putResponse = await api.put('/api/systemSettings', {
          data: {
            ...current,
            showLoginQr: originalShowLoginQr ?? false,
            mobileServerUrl: originalMobileServerUrl ?? '',
          },
        });
        if (!putResponse.ok()) {
          throw new Error(
            `PUT /api/systemSettings failed: HTTP ${putResponse.status()}`,
          );
        }
      });
    } catch (error) {
      // Best-effort teardown -- report the failure but don't mask the suite's
      // real result by throwing out of afterAll.
      console.warn('Failed to restore login QR system settings', error);
    }
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

  test('the QR also shows on the About dialog while logged in', async ({
    browser,
  }) => {
    // The whole point of the About QR: a signed-in user can reach it without
    // logging out. This context rides on appData (not GET /featureConfig).
    const context = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const page = await context.newPage();
    await stubTokenRefresh(page);
    // The sidebar holding the About menu defaults to closed and its open/closed
    // flag is persisted, so seed it rather than blind-clicking the toggle (which
    // would CLOSE it if a saved session already had it open).
    await page.addInitScript(() => {
      const layout = JSON.parse(localStorage.getItem('layout') ?? '{}');
      localStorage.setItem(
        'layout',
        JSON.stringify({ ...layout, isSidebarOpen: true }),
      );
    });

    await page.goto('/');
    await page.getByTestId('sidebar-avatar-menu').click();
    await page.getByRole('menuitem', { name: 'About' }).click();

    const dialog = page.getByRole('dialog');
    await expect(dialog.getByText('Mobile App')).toBeVisible();
    await expect(dialog.getByRole('img', { name: QR_IMG_NAME })).toBeVisible();

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
