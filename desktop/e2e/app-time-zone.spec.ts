import { expect, type Locator, type Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateGroup,
  apiCreateReceipt,
  apiDeleteGroupById,
  apiGetUserId,
  apiPatchSystemSettings,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import { addGroupToScopeByName, gotoReportBuilder, openComboboxAndPick } from './helpers/reports';

// The app time zone System Setting, end to end.
//
// The Jest specs cover the picker, the payload and the appDate pipe against a
// mocked store; what only an e2e can prove is the wire: the zone picked in the
// typeahead is persisted by the server, comes back on AppData, and every
// instant the app renders moves with it — without a server restart, and in
// both directions.
//
// The assertion is a receipt's "Added at" moment on its view page, rendered
// with appDate's "medium" format. A receipt's created_at cannot be set through
// the API, so rather than back-dating one into the late-Sep-30-Eastern window
// (the Go report tests pin that case), this compares the rendered time against
// the receipt's real created_at formatted in each zone. New York is never on
// UTC, so the two strings always differ in the hour.
//
// `timeZone` is a GLOBAL system setting, so this suite mutates shared server
// state, runs serially, and restores what it found in afterAll. It relies on
// the `e2e-shared-backend` job-level concurrency group in e2e.yml /
// mobile-e2e.yml to keep other e2e jobs off the backend meanwhile.

test.use({ storageState: 'e2e/.auth/admin.json' });

const NEW_YORK = 'America/New_York';
const UTC = 'UTC';

/**
 * [instant] as Angular's en-US "medium" format renders it in [timeZone]:
 * `MMM d, y, h:mm:ss a`, e.g. "Oct 8, 2026, 2:15:30 PM". Assembled from parts
 * rather than Intl's own joiners, which vary by ICU version (", " vs " at ").
 */
function mediumInZone(instant: string, timeZone: string): string {
  const parts = new Intl.DateTimeFormat('en-US', {
    timeZone,
    month: 'short',
    day: 'numeric',
    year: 'numeric',
    hour: 'numeric',
    minute: '2-digit',
    second: '2-digit',
    hour12: true,
  }).formatToParts(new Date(instant));
  const part = (type: Intl.DateTimeFormatPartTypes) =>
    parts.find((p) => p.type === type)?.value ?? '';

  return `${part('month')} ${part('day')}, ${part('year')}, ${part('hour')}:${part('minute')}:${part('second')} ${part('dayPeriod')}`;
}

/** Collapses every whitespace run — Angular's CLDR data puts a U+202F before AM/PM. */
function normalize(text: string | null): string {
  return (text ?? '').replace(/\s+/g, ' ').trim();
}

function timeZonePicker(page: Page): Locator {
  return page.getByTestId('system-settings-time-zone');
}

/** The "Added at: …" line of the receipt view's Audit Details section. */
function addedAtLine(page: Page): Locator {
  return page.getByText(/^\s*Added at:/);
}

async function expectAddedAt(page: Page, createdAt: string, timeZone: string): Promise<void> {
  await expect
    .poll(async () => normalize(await addedAtLine(page).textContent()))
    .toBe(normalize(`Added at: ${mediumInZone(createdAt, timeZone)}`));
}

async function getStoredTimeZone(): Promise<string> {
  return withAdminApi(async (api) => {
    const res = await api.get('/api/systemSettings');
    if (!res.ok()) {
      throw new Error(`GET /api/systemSettings failed: HTTP ${res.status()}`);
    }
    return (await res.json()).timeZone;
  });
}

test.describe.serial('App time zone (System Settings → every rendered instant)', () => {
  let originalTimeZone: string | undefined;
  let group: { id: number; name: string };
  let receiptId: number;
  let createdAt: string;

  test.beforeAll(async () => {
    originalTimeZone = await getStoredTimeZone();

    await withAdminApi(async (api) => {
      // Start from UTC so the first test's switch is a real change whatever
      // the backend was left on.
      await apiPatchSystemSettings(api, { timeZone: UTC });

      const adminId = await apiGetUserId(api, creds('admin').username);
      group = await apiCreateGroup(api, uniqueName('tz-group'));
      receiptId = await apiCreateReceipt(api, {
        groupId: group.id,
        paidByUserId: adminId,
        name: uniqueName('tz-receipt'),
      });

      const res = await api.get(`/api/receipt/${receiptId}`);
      if (!res.ok()) {
        throw new Error(`GET /api/receipt/${receiptId} failed: HTTP ${res.status()}`);
      }
      createdAt = (await res.json()).createdAt;
    });

    // Guards the premise: the two renderings must differ or the suite proves nothing.
    expect(mediumInZone(createdAt, NEW_YORK)).not.toEqual(mediumInZone(createdAt, UTC));
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        // Restore the captured value, never a hardcoded default — an environment
        // that configured its own zone must get it back.
        await apiPatchSystemSettings(api, { timeZone: originalTimeZone || UTC });
        if (group) {
          await apiDeleteGroupById(api, String(group.id));
        }
      });
    } catch (error) {
      // Best-effort teardown — report it, but don't mask the suite's real result.
      console.warn('Failed to restore the app time zone / delete the e2e group', error);
    }
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('renders instants in UTC by default', async ({ page }) => {
    await page.goto(`/receipts/${receiptId}/view`);
    await expectAddedAt(page, createdAt, UTC);
  });

  test('the admin picks America/New_York in the typeahead and it persists', async ({ page }) => {
    await page.goto('/system-settings/settings/edit');
    const picker = timeZonePicker(page);
    const input = picker.getByRole('combobox');
    await expect(input).toHaveValue(UTC);

    // A single-select autocomplete goes readonly once it holds a value, so the
    // current zone has to be cleared before the typeahead accepts input.
    await picker.getByTestId('autocomplete-clear').click();
    await input.click();
    await input.fill('New_York');
    await page.getByRole('option', { name: NEW_YORK, exact: true }).click();
    await expect(input).toHaveValue(NEW_YORK);

    await page.getByRole('button', { name: 'Save' }).click();
    await expect(page).toHaveURL(/\/system-settings\/settings\/view/);

    expect(await getStoredTimeZone()).toBe(NEW_YORK);

    // ...and the view page reads it back out of a fresh resolver fetch.
    await page.reload();
    await expect(timeZonePicker(page).getByRole('combobox')).toHaveValue(NEW_YORK);
  });

  test('a receipt\'s Added At and the report hint follow the New York zone', async ({ page }) => {
    await page.goto(`/receipts/${receiptId}/view`);
    await expectAddedAt(page, createdAt, NEW_YORK);

    // An instant period field names the zone it is read in.
    await gotoReportBuilder(page);
    await addGroupToScopeByName(page, group.name);
    await openComboboxAndPick(
      page,
      page.getByRole('combobox', { name: 'Date field', exact: true }),
      page.getByRole('option', { name: 'Added At', exact: true }),
    );
    await expect(page.getByText(/Resolves to .* on Added At \(America\/New_York\)/)).toBeVisible();
  });

  test('switching back to UTC moves every rendered instant back, no restart', async ({ page }) => {
    // Through the API this time: the browser only learns the new zone from the
    // AppData fetched on load, which is the path every other user relies on.
    await withAdminApi((api) => apiPatchSystemSettings(api, { timeZone: UTC }));

    await page.goto(`/receipts/${receiptId}/view`);
    await expectAddedAt(page, createdAt, UTC);

    await gotoReportBuilder(page);
    await addGroupToScopeByName(page, group.name);
    await openComboboxAndPick(
      page,
      page.getByRole('combobox', { name: 'Date field', exact: true }),
      page.getByRole('option', { name: 'Added At', exact: true }),
    );
    await expect(page.getByText(/Resolves to .* on Added At \(UTC\)/)).toBeVisible();
  });
});
