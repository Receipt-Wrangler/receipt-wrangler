import { expect, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';
import { withAdminApi } from './helpers/provisioning';

// 45 days. Divides evenly into days, so it also pins the hours <-> unit round
// trip: the form must render 1080 stored hours back as "45 Days", not
// "1080 Hours".
const RETENTION_DAYS = 45;
const RETENTION_HOURS = RETENTION_DAYS * 24;

// `tempFileRetentionHours` is a GLOBAL system setting, so these tests mutate
// shared server state and must run serially. afterAll puts the field back to
// whatever it was before the suite -- never a hardcoded default, which would
// silently retune an environment that had configured it.
test.describe.serial('Temp file retention (System Settings)', () => {
  let originalTempFileRetentionHours: unknown;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      const res = await api.get('/api/systemSettings');
      if (!res.ok()) {
        throw new Error(`GET /api/systemSettings failed: HTTP ${res.status()}`);
      }
      originalTempFileRetentionHours = (await res.json())
        .tempFileRetentionHours;
    });
  });

  test.afterAll(async () => {
    // Re-read the live settings and overlay only the captured field, so a
    // concurrent change to an unrelated setting isn't clobbered by a stale
    // snapshot (the PUT is an upsert needing the full object).
    try {
      await withAdminApi(async (api) => {
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
            tempFileRetentionHours: originalTempFileRetentionHours ?? 720,
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
      console.warn('Failed to restore temp file retention system setting', error);
    }
  });

  test('admin sets the retention in days and it persists as hours', async ({
    browser,
  }) => {
    const context = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const page = await context.newPage();
    await stubTokenRefresh(page);

    await page.goto('/system-settings/settings/edit');
    await page.getByLabel('Keep failed uploads for').fill(String(RETENTION_DAYS));
    await page
      .locator('app-form-section')
      .filter({ hasText: 'Keep failed uploads for' })
      .getByLabel('Unit')
      .click();
    await page.getByRole('option', { name: 'Days' }).click();
    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page).toHaveURL(/\/system-settings\/settings\/view/);

    // The unit is presentation only -- the wire and the DB store whole hours.
    await withAdminApi(async (api) => {
      const settings = await (await api.get('/api/systemSettings')).json();
      expect(settings.tempFileRetentionHours).toBe(RETENTION_HOURS);
    });

    // ...and the view page renders those hours back in the friendlier unit.
    await page.reload();
    await expect(page.getByLabel('Keep failed uploads for')).toHaveValue(
      String(RETENTION_DAYS),
    );

    await context.close();
  });

  // The server floors this at 24 hours. Without the client-side floor the entry
  // passes validation and comes back as a bare 400 with nothing on the field.
  test('the form rejects a retention below the 24 hour minimum', async ({
    browser,
  }) => {
    const context = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const page = await context.newPage();
    await stubTokenRefresh(page);

    await page.goto('/system-settings/settings/edit');

    const section = page
      .locator('app-form-section')
      .filter({ hasText: 'Keep failed uploads for' });

    await section.getByLabel('Unit').click();
    await page.getByRole('option', { name: 'Hours' }).click();
    await page.getByLabel('Keep failed uploads for').fill('23');
    await page.getByLabel('Keep failed uploads for').blur();

    await expect(section.getByText('Must be at least 24.')).toBeVisible();

    await context.close();
  });
});
