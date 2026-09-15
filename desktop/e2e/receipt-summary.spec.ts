import { expect, test, type Page } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateCustomField,
  apiCreateGroup,
  apiCreateReceipt,
  apiDeleteCustomFieldById,
  apiDeleteGroupById,
  apiGetUserId,
  apiSetGroupSummaryConfig,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import { gotoReceiptsTable } from './helpers/receipts-table';

// The receipt summary, end to end (see desktop/CLAUDE.md → "Receipt summary (the totals under the
// table)"). The Jest specs cover the aggregation rules and the refresh split against a MOCKED
// service, and the Go tests cover the endpoint — what only an e2e can prove is that the three
// pieces agree across the wire:
//
//   - the settings form's command actually persists (the Jest spec asserts the command's shape,
//     never the round-trip), so a field missing from the API's update assignment block fails here
//     and nowhere else;
//   - the configured set reaches a different consumer over a different endpoint — the receipts
//     table, via POST /receipt/group/{id}/summary;
//   - the "All" group's chip row sends a configurationGroupId the server honours.
//
// Runs as admin: the settings page needs group.update, and seeding needs
// app.custom-fields.create/delete.
test.use({ storageState: 'e2e/.auth/admin.json' });

test.describe('Receipt summary', () => {
  test.describe.configure({ mode: 'serial' });

  // Names end in distinct words: getByLabel and the autocomplete filter both match substrings, so a
  // name that is a prefix of another makes every "this field is offered" assertion ambiguous.
  const hstName = uniqueName('rs-fieldhst');
  const subtotalName = uniqueName('rs-fieldsubtotal');
  const notesName = uniqueName('rs-fieldnotes');

  const coffee = uniqueName('rs-coffee');
  const lunch = uniqueName('rs-lunch');
  const hotel = uniqueName('rs-hotel');
  const draft = uniqueName('rs-draft');
  const secondaryDraft = uniqueName('rs-secondarydraft');

  let hst: { id: number; name: string };
  let subtotal: { id: number; name: string };
  let notes: { id: number; name: string };
  let primary: { id: number; name: string };
  let secondary: { id: number; name: string };
  let off: { id: number; name: string };

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      hst = await apiCreateCustomField(api, { name: hstName, type: 'CURRENCY' });
      subtotal = await apiCreateCustomField(api, { name: subtotalName, type: 'CURRENCY' });
      // Only CURRENCY fields can be totalled; this one exists to prove the picker excludes it.
      notes = await apiCreateCustomField(api, { name: notesName, type: 'TEXT' });

      primary = await apiCreateGroup(api, uniqueName('rs-primary'));
      secondary = await apiCreateGroup(api, uniqueName('rs-secondary'));
      off = await apiCreateGroup(api, uniqueName('rs-off'));

      await apiSetGroupSummaryConfig(api, primary.id, {
        enabled: true,
        statuses: ['OPEN', 'RESOLVED'],
        currencyCustomFieldIds: [hst.id, subtotal.id],
      });
      // A deliberately different shape, so switching the chip on the All group changes what renders.
      await apiSetGroupSummaryConfig(api, secondary.id, {
        enabled: true,
        statuses: ['DRAFT'],
        currencyCustomFieldIds: [hst.id],
      });

      const userId = await apiGetUserId(api, creds('admin').username);

      // Amounts are chosen so a THOUSANDS SEPARATOR is load-bearing in three assertions: the
      // currency symbol and separators are global System Settings this spec must not mutate, so
      // money is asserted separator-tolerantly and a value without one would prove less.
      //   overall  11.30 + 22.60 + 1113.00 + 50.00 = 1196.90
      //   HST       1.30 +  2.60 +  113.00         =  116.90
      //   Subtotal 10.00 + 20.00 + 1000.00         = 1030.00
      const seed = (
        name: string,
        amount: string,
        status: string,
        values?: { hst: string; subtotal: string },
      ) =>
        apiCreateReceipt(api, {
          groupId: primary.id,
          paidByUserId: userId,
          name,
          amount,
          status,
          ...(values
            ? {
                customFields: [
                  { customFieldId: hst.id, currencyValue: values.hst },
                  { customFieldId: subtotal.id, currencyValue: values.subtotal },
                ],
              }
            : {}),
        });

      await seed(coffee, '11.30', 'OPEN', { hst: '1.30', subtotal: '10.00' });
      await seed(lunch, '22.60', 'OPEN', { hst: '2.60', subtotal: '20.00' });
      await seed(hotel, '1113.00', 'RESOLVED', { hst: '113.00', subtotal: '1000.00' });
      // DRAFT is NOT a configured status for this group: it must count toward the overall row and
      // get no row of its own.
      await seed(draft, '50.00', 'DRAFT');

      await apiCreateReceipt(api, {
        groupId: secondary.id,
        paidByUserId: userId,
        name: secondaryDraft,
        amount: '7.00',
        status: 'DRAFT',
      });
    });
  });

  test.afterAll(async () => {
    // Wrapped so a cleanup failure cannot mask a test failure, and per-id guarded so a beforeAll
    // that died partway still tears down what it did create.
    try {
      await withAdminApi(async (api) => {
        // Groups first: deleting a custom field destroys every value stored against it, so the
        // receipts holding those values have to be gone (cascaded with their group) first.
        for (const group of [primary, secondary, off]) {
          if (group) {
            await apiDeleteGroupById(api, String(group.id));
          }
        }
        for (const field of [hst, subtotal, notes]) {
          if (field) {
            await apiDeleteCustomFieldById(api, field.id);
          }
        }
      });
    } catch (error) {
      console.warn(`receipt-summary teardown failed: ${error}`);
    }
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  /** A selected chip in the summary's currency-field picker (the `chip` idiom from the sibling spec). */
  const summaryFieldChip = (page: Page, name: string) =>
    page
      .getByTestId('receipt-summary-custom-fields')
      .locator('mat-chip-row')
      .filter({ hasText: name });

  /**
   * Matches one figure in a summary row by NAME, then its value.
   *
   * `\D*` swallows the currency symbol wherever `currencySymbolPosition` puts it but cannot cross a
   * digit, so it can never slide onto the neighbouring figure; the separators are character classes
   * because the symbol and both separators are global System Settings this spec must not mutate; and
   * the trailing `(?!\d)` stops a short value matching inside a longer one.
   */
  const figure = (name: string, whole: string, cents: string) =>
    new RegExp(`${name.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')}\\D*${whole}[.,]${cents}(?!\\d)`);

  const totals = (page: Page) => page.getByTestId('receipt-totals');
  const overallRow = (page: Page) => page.getByTestId('receipt-totals-row-overall');
  const statusRow = (page: Page, status: string) =>
    page.getByTestId(`receipt-totals-row-${status}`);

  /** Opens a seeded group's receipts table. gotoReceiptsTable establishes the session first. */
  async function gotoGroupTable(page: Page, groupId: number): Promise<void> {
    await gotoReceiptsTable(page);
    await page.goto(`/receipts/group/${groupId}`);
    await expect(page.getByTestId('receipts-overflow-menu')).toBeVisible();
  }

  /** Counts POSTs to the summary endpoint from the moment it is installed. */
  function countSummaryRequests(page: Page): () => number {
    let count = 0;
    page.on('request', (request) => {
      if (/\/api\/receipt\/group\/\d+\/summary$/.test(request.url()) && request.method() === 'POST') {
        count += 1;
      }
    });

    return () => count;
  }

  /** Resolves when the table's own paged fetch lands — the anchor for "no summary request". */
  function waitForPagedFetch(page: Page) {
    return page.waitForResponse(
      (response) =>
        /\/api\/receipt\/group\/\d+$/.test(response.url()) &&
        response.request().method() === 'POST',
    );
  }

  // Gap 1: the Jest spec asserts the command's SHAPE against a mocked service. Only this proves the
  // three pieces survive a real PUT and come back out of GET /group/{id} through the resolver.
  test('round-trips the summary configuration through the settings form', async ({ page }) => {
    await page.goto(`/groups/${off.id}/receipt-settings/edit`);

    const enabled = page.getByTestId('receipt-summary-enabled').getByRole('checkbox');
    await expect(enabled).toBeVisible();
    await enabled.check();
    await page.getByTestId('receipt-summary-status-OPEN').getByRole('checkbox').check();
    await page.getByTestId('receipt-summary-status-RESOLVED').getByRole('checkbox').check();

    const picker = page.getByTestId('receipt-summary-custom-fields');
    const pickerInput = picker.getByRole('combobox');
    await pickerInput.click();

    // Only CURRENCY fields are selectable — the server 400s anything else, because only a currency
    // value can be summed. A TEXT field must never be offered.
    await expect(page.getByRole('option', { name: hstName, exact: true })).toBeVisible();
    await expect(page.getByRole('option', { name: notesName, exact: true })).toHaveCount(0);

    await pickerInput.fill(hstName);
    await page.getByRole('option', { name: hstName, exact: true }).click();
    await expect(summaryFieldChip(page, hstName)).toBeVisible();
    // The open overlay sits above Save and would intercept the click.
    await page.keyboard.press('Escape');

    await page.getByRole('button', { name: 'Save', exact: true }).first().click();
    await page.waitForURL(/\/receipt-settings\/view/);

    // Re-navigate rather than assert in place: the values have to survive a fresh resolver fetch.
    // A field missing from the API's update assignment block persists nothing and fails only here.
    await page.goto(`/groups/${off.id}/receipt-settings/edit`);
    await expect(page.getByTestId('receipt-summary-enabled').getByRole('checkbox')).toBeChecked();
    await expect(
      page.getByTestId('receipt-summary-status-OPEN').getByRole('checkbox'),
    ).toBeChecked();
    await expect(
      page.getByTestId('receipt-summary-status-RESOLVED').getByRole('checkbox'),
    ).toBeChecked();
    await expect(
      page.getByTestId('receipt-summary-status-DRAFT').getByRole('checkbox'),
    ).not.toBeChecked();
    await expect(summaryFieldChip(page, hstName)).toBeVisible();

    // Leave the group as the suite found it, so the "disabled renders nothing" test below is not
    // order-dependent on this one.
    await withAdminApi(async (api) => {
      await apiSetGroupSummaryConfig(api, off.id, { enabled: false });
    });
  });

  test('totals the whole filter result set, and agrees with the table', async ({ page }) => {
    await gotoGroupTable(page, primary.id);
    await expect(totals(page)).toBeVisible();

    // Money is separator-tolerant and symbol-free: the currency configuration is a global System
    // Setting this spec must not mutate.
    await expect(overallRow(page)).toContainText('4 receipts');
    await expect(overallRow(page)).toContainText(figure('Total', '1[,. ]196', '90'));
    await expect(overallRow(page)).toContainText(figure(hstName, '116', '90'));
    await expect(overallRow(page)).toContainText(figure(subtotalName, '1[,. ]030', '00'));

    await expect(statusRow(page, 'OPEN')).toContainText('2 receipts');
    await expect(statusRow(page, 'OPEN')).toContainText(figure('Total', '33', '90'));
    await expect(statusRow(page, 'OPEN')).toContainText(figure(hstName, '3', '90'));

    // Singular, and the RESOLVED receipt is the one carrying the thousands separator.
    await expect(statusRow(page, 'RESOLVED')).toContainText('1 receipt');
    await expect(statusRow(page, 'RESOLVED')).toContainText(figure('Total', '1[,. ]113', '00'));

    // DRAFT is not configured for this group: it counts toward the overall row (it IS in the filter
    // result — excluding it would make the total disagree with the table) but gets no row.
    await expect(statusRow(page, 'DRAFT')).toHaveCount(0);

    // The cheapest canary there is: the summary's count and the table's own rendering come from two
    // separate queries, and all four seeded receipts fit one page. If these ever diverge, the two
    // queries have drifted apart.
    await expect(page.getByRole('link', { name: /^e2e-rs-/ })).toHaveCount(4);
  });

  // The feature's headline rule: a configured status that matches nothing still renders, so the
  // block keeps its shape as the filter narrows.
  test('re-computes on a filter, keeping an unmatched status as a zero row', async ({ page }) => {
    await gotoGroupTable(page, primary.id);
    await expect(overallRow(page)).toContainText('4 receipts');

    await page.getByTestId('receipts-filter').click();
    const dialog = page.getByRole('dialog');
    await expect(dialog).toBeVisible();

    // The Status row is an autocomplete, not a select: it filters on the raw value (OPEN) and
    // displays the formatted one (Open). Its operation is auto-filled by the dialog, so only the
    // value needs setting.
    const status = dialog.getByRole('combobox', { name: 'Status' });
    await status.click();
    await status.fill('OPEN');
    await page.getByRole('option', { name: 'Open', exact: true }).click();
    // The panel stays open over Save and has to be dismissed, but NOT with Escape:
    // MatAutocomplete only consumes that while its panel is open, so on the runs where
    // the panel has already closed the same keypress reaches MatDialog and closes the
    // whole dialog (receipt-quick-date-filter.spec.ts asserts exactly that), detaching
    // the submit button mid-click. Tab blurs the input, which closes the panel if it is
    // open and is harmless if it is not.
    await page.keyboard.press('Tab');
    await expect(page.getByRole('option', { name: 'Open', exact: true })).toBeHidden();

    await dialog.getByTestId('dialog-submit-button').click();
    await expect(dialog).toBeHidden();
    await expect(page.getByTestId('receipt-filter-chip-status')).toBeVisible();

    await expect(overallRow(page)).toContainText('2 receipts');
    await expect(overallRow(page)).toContainText(figure('Total', '33', '90'));
    await expect(statusRow(page, 'OPEN')).toContainText('2 receipts');

    // Still rendered, reading zero — not dropped.
    await expect(statusRow(page, 'RESOLVED')).toBeVisible();
    await expect(statusRow(page, 'RESOLVED')).toContainText('0 receipts');
  });

  // The off state is a well-formed 200 with enabled:false, NOT a skipped request: the component's
  // no-request guard only fires on the All group when no member group has a summary. A real group
  // whose summary is off still asks, and renders nothing off the answer — deliberately, because
  // gating on the client's cached groupReceiptSettings would render a stale block.
  test('renders nothing off an enabled:false response when the summary is off', async ({ page }) => {
    await gotoReceiptsTable(page);

    // The body is captured through a route rather than read off waitForResponse: the navigation
    // that triggers the call also evicts the response body, so reading it afterwards fails with
    // "No resource with given identifier found". route.fetch + fulfill is the helpers' own idiom.
    const answers: { enabled: boolean }[] = [];
    await page.route(/\/api\/receipt\/group\/\d+\/summary$/, async (route) => {
      const response = await route.fetch();
      const body = (await response.json()) as { enabled: boolean };
      answers.push(body);
      await route.fulfill({ response, json: body });
    });

    await Promise.all([waitForPagedFetch(page), page.goto(`/receipts/group/${off.id}`)]);
    await expect(page.getByTestId('receipts-overflow-menu')).toBeVisible();

    await expect(totals(page)).toHaveCount(0);
    // A real group whose summary is off IS asked, and answers enabled:false at a normal 200 — not a
    // 404, and not a skipped request.
    expect(answers).toHaveLength(1);
    expect(answers[0].enabled).toBe(false);
  });

  // Gaps 2 and 3: the All group spans several groups and has no configuration of its own, so the
  // chip row picks whose applies. Only an e2e proves the configurationGroupId reaches the server.
  test('switches configuration from the All group chip row and remembers the pick', async ({
    page,
  }) => {
    await gotoReceiptsTable(page);

    // The All group is synthetic and its id varies per install, so it is resolved from the API's own
    // isAllGroup marker rather than assumed (single-group-default.spec.ts inlines the same fetch).
    const allGroupId = await withAdminApi(async (api) => {
      const groups = (await (await api.get('/api/group/')).json()) as {
        id: number;
        isAllGroup?: boolean;
      }[];
      const all = groups.find((group) => group.isAllGroup);
      if (!all) {
        throw new Error('no All group visible to the admin');
      }
      return all.id;
    });

    await page.goto(`/receipts/group/${allGroupId}`);
    await expect(page.getByTestId('receipts-overflow-menu')).toBeVisible();
    await expect(totals(page)).toBeVisible();

    // Assert the chips BY TESTID and never by count or position: the row lists every group the
    // admin belongs to whose summary is enabled, and a group leaked by a crashed earlier run would
    // otherwise break an exact-count assertion. Auto-select-alphabetically-first is covered
    // deterministically in the Jest util spec and deliberately not asserted here.
    const primaryChip = page.getByTestId(`receipt-totals-config-group-${primary.id}`);
    const secondaryChip = page.getByTestId(`receipt-totals-config-group-${secondary.id}`);
    await expect(primaryChip).toBeVisible();
    await expect(secondaryChip).toBeVisible();

    // Only the COUNT is comparable across the switch: the configurations differ in which currency
    // columns they render, so the row's full text legitimately changes.
    const receiptCount = (text: string | null) => text?.match(/\d+ receipts?/)?.[0];
    const countBefore = receiptCount(await overallRow(page).textContent());
    expect(countBefore).toBeDefined();

    // mat-chip-option's host is role="presentation"; the clickable element is the inner action.
    await secondaryChip.getByRole('option').click();

    // Secondary configures DRAFT only, so the rendered rows change...
    await expect(statusRow(page, 'DRAFT')).toBeVisible();
    await expect(statusRow(page, 'OPEN')).toHaveCount(0);
    await expect(statusRow(page, 'RESOLVED')).toHaveCount(0);
    // ...while the DATA still spans every group. This is the distinction the whole chip design
    // rests on: the chip picks a configuration, never a data scope. The figures themselves are not
    // assertable here — the All group aggregates every group the admin belongs to on a backend other
    // specs are mutating — so the count before/after is the assertion.
    await expect(overallRow(page)).toContainText(countBefore!);

    // The pick is persisted in the receiptTable slice, so it has to survive a reload.
    await page.reload();
    await expect(totals(page)).toBeVisible();
    await expect(statusRow(page, 'DRAFT')).toBeVisible();
    await expect(statusRow(page, 'OPEN')).toHaveCount(0);
  });

  // The summary rides its own refresh stream precisely so paging and sorting do not re-request an
  // unpaged aggregate. The Jest spec proves the split against a mocked service; only this proves it
  // on the wire.
  test('does not re-request the summary when only the sort changes', async ({ page }) => {
    await gotoGroupTable(page, primary.id);
    await expect(overallRow(page)).toContainText('4 receipts');

    const summaryRequests = countSummaryRequests(page);

    // Anchored to the paged fetch the sort DOES trigger, so a zero below means "the summary was
    // skipped", not "nothing happened".
    await Promise.all([
      waitForPagedFetch(page),
      page.getByRole('columnheader', { name: 'Amount' }).click(),
    ]);

    // The figures are unchanged by a sort, so they must still be on screen...
    await expect(overallRow(page)).toContainText('4 receipts');
    // ...and must not have been refetched to get there.
    expect(summaryRequests()).toBe(0);
  });
});
