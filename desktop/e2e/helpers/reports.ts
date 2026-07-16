import { expect, Locator, Page } from '@playwright/test';

// Shared Report Builder e2e flows. The reports area is gated by the app-level
// app.reports.read permission, so specs run as admin
// (test.use({ storageState: 'e2e/.auth/admin.json' })).

/** Navigate to the Reports list (the /reports landing page) and wait for it to render. */
export async function gotoReports(page: Page): Promise<void> {
  await page.goto('/reports');
  await expect(page.getByRole('heading', { name: 'Reports', level: 1 })).toBeVisible();
}

/** Navigate straight to the builder (New Report) and wait for it to render. */
export async function gotoReportBuilder(page: Page): Promise<void> {
  await page.goto('/reports/new');
  await expect(page.getByText('Report Builder')).toBeVisible();
}

/**
 * Add the first available group to the report scope via the Add Group dialog.
 * This is the precondition for a live preview (a report needs at least one group).
 */
export async function addFirstGroupToScope(page: Page): Promise<void> {
  await page.getByTestId('report-add-group').click();
  await page.getByTestId('add-group-select').first().click();
  await page.getByTestId('dialog-submit-button').click();
}

/** Add a specific (provisioned) group to the report scope by its display name. */
export async function addGroupToScopeByName(page: Page, name: string): Promise<void> {
  await page.getByTestId('report-add-group').click();
  await page.getByTestId('add-group-row').filter({ hasText: name }).getByTestId('add-group-select').click();
  await page.getByTestId('dialog-submit-button').click();
}

/**
 * Open a left-pane select and pick an option, retrying the open. The config panel
 * re-renders as the debounced preview refreshes, which can drop a freshly opened
 * mat-select panel before the option is clicked; the retry re-opens until the
 * option is visible. (Modal dialogs are stable and need no retry.)
 *
 * Chaining two picks on the *same* rapidly-re-rendering select can race the first
 * pick's debounced preview (a retry may re-open onto a stray option); settle with
 * waitForPreview between such picks.
 */
export async function openComboboxAndPick(page: Page, combobox: Locator, option: Locator): Promise<void> {
  await expect(async () => {
    await combobox.click();
    await expect(option).toBeVisible({ timeout: 2000 });
  }).toPass({ timeout: 15_000 });
  await option.click();
}

/** Resolve once the live preview responds — used to settle the debounced refresh. */
export function waitForPreview(page: Page) {
  return page.waitForResponse(
    (response) => response.url().includes('/api/report/preview') && response.request().method() === 'POST',
    { timeout: 20_000 },
  );
}

/**
 * Add a grouping level via the "Add grouping level…" picker and settle the
 * resulting debounced preview before returning, so a following pick on the same
 * (re-rendering) select doesn't race the refresh.
 */
export async function addGroupingLevel(page: Page, label: string): Promise<void> {
  const combobox = page.getByRole('combobox', { name: /Add grouping level/ });
  const option = page.getByRole('option', { name: label, exact: true });
  await Promise.all([waitForPreview(page), openComboboxAndPick(page, combobox, option)]);
}
