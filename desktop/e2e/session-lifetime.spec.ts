import { expect, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import { withAdminApi } from './helpers/provisioning';

// 2 days. Chosen because it divides evenly into days, so it also pins the
// hours <-> unit round trip: the form must render 48 stored hours back as
// "2 Days", not "48 Hours".
const LIFETIME_DAYS = 2;
const LIFETIME_HOURS = LIFETIME_DAYS * 24;

// Allow generous slack on the expiry assertion -- the value is minted when the
// server handles the login, which is some way after this suite reads the clock.
const SLACK_MS = 10 * 60 * 1000;

// `refreshTokenValidForHours` is a GLOBAL system setting, so these tests mutate
// shared server state and must run serially. afterAll puts the field back to
// whatever it was before the suite -- never a hardcoded default, which would
// silently retune sessions on an environment that had configured it.
test.describe.serial('Session lifetime (System Settings → refresh cookie)', () => {
  let originalRefreshTokenValidForHours: unknown;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      const res = await api.get('/api/systemSettings');
      if (!res.ok()) {
        throw new Error(`GET /api/systemSettings failed: HTTP ${res.status()}`);
      }
      originalRefreshTokenValidForHours = (await res.json())
        .refreshTokenValidForHours;
    });
  });

  test.afterAll(async () => {
    // Re-read the live settings and overlay only the captured field, so a
    // concurrent change to an unrelated setting isn't clobbered by a stale
    // snapshot (the PUT is an upsert needing the full object).
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
            refreshTokenValidForHours: originalRefreshTokenValidForHours ?? 24,
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
      console.warn('Failed to restore session lifetime system setting', error);
    }
  });

  test('admin sets the session lifetime in days and it persists as hours', async ({
    browser,
  }) => {
    const context = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const page = await context.newPage();
    await stubTokenRefresh(page);

    await page.goto('/system-settings/settings/edit');
    await page.getByLabel('Stay signed in for').fill(String(LIFETIME_DAYS));
    await page
      .locator('app-form-section')
      .filter({ hasText: 'Stay signed in for' })
      .getByLabel('Unit')
      .click();
    await page.getByRole('option', { name: 'Days' }).click();
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page).toHaveURL(/\/system-settings\/settings\/view/);

    // The unit is presentation only -- the wire and the DB store whole hours.
    await withAdminApi(async (api) => {
      const settings = await (await api.get('/api/systemSettings')).json();
      expect(settings.refreshTokenValidForHours).toBe(LIFETIME_HOURS);
    });

    // ...and the view page renders those hours back in the friendlier unit.
    await page.reload();
    await expect(page.getByLabel('Stay signed in for')).toHaveValue(
      String(LIFETIME_DAYS),
    );

    await context.close();
  });

  test('a fresh login gets a refresh cookie sized to the configured lifetime', async ({
    browser,
  }) => {
    // A fresh, unauthenticated context so the cookies come from this login and
    // not from a stored session.
    const context = await browser.newContext({
      storageState: { cookies: [], origins: [] },
    });
    const page = await context.newPage();

    const { username, password } = creds('admin');
    await page.goto('/auth/login');
    await page.getByLabel('Username').fill(username);
    await page.getByLabel('Password').fill(password);

    const before = Date.now();
    await page.getByRole('button', { name: 'Login' }).click();
    await expect(page).toHaveURL(/\/dashboard\/group\/\d+/, { timeout: 15_000 });

    const cookies = await context.cookies();
    const refreshToken = cookies.find((c) => c.name === 'refresh_token');
    const jwt = cookies.find((c) => c.name === 'jwt');

    expect(refreshToken, 'refresh_token cookie should be set').toBeDefined();
    expect(jwt, 'jwt cookie should be set').toBeDefined();

    // This is the assertion the unit tests cannot make: it proves the value
    // survived model → command → DB → JWT → Set-Cookie.
    const refreshExpiresMs = refreshToken!.expires * 1000;
    const expected = before + LIFETIME_HOURS * 60 * 60 * 1000;
    expect(refreshExpiresMs).toBeGreaterThan(expected - SLACK_MS);
    expect(refreshExpiresMs).toBeLessThan(expected + SLACK_MS);

    // The access token is deliberately NOT configurable and stays at 20
    // minutes -- both clients size their refresh timer against that window.
    const jwtExpiresMs = jwt!.expires * 1000;
    expect(jwtExpiresMs).toBeGreaterThan(before);
    expect(jwtExpiresMs).toBeLessThan(before + 25 * 60 * 1000);

    await context.close();
  });

  test('the form rejects a lifetime above the 30 day maximum', async ({
    browser,
  }) => {
    const context = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const page = await context.newPage();
    await stubTokenRefresh(page);

    await page.goto('/system-settings/settings/edit');
    await page.getByLabel('Stay signed in for').fill('31');
    await page.getByLabel('Stay signed in for').blur();

    // The message has to be visible -- Validators.max has no message mapping in
    // BaseInputComponent and would render an empty mat-error.
    await expect(page.getByText('Must be at most 30.')).toBeVisible();

    await context.close();
  });
});
