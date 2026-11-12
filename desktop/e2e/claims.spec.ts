import { expect, type Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateGroup,
  apiCreateReceipt,
  apiDeleteGroupById,
  apiGetUserId,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import { openReceiptsOverflowMenu } from './helpers/receipts-table';

// Reimbursement claims, end to end (see desktop/CLAUDE.md → "Claims"). The Jest specs mock
// ReimbursementClaimService and the Go tests drive the handlers directly, so neither proves the two
// halves agree on the wire. This walks the main journey through the real UI and API:
//
//   - a claim created on the claim form, then filled from the receipts table's bulk action;
//   - a receipt already in another claim refused and named, not moved;
//   - a payment and a recoupment moving the server's figures;
//   - the receipt delete endpoints refusing a receipt that is in a claim, singly and in bulk;
//   - a receipt moved to another group on the receipt form leaving its claim.
//
// Runs as admin in two groups it provisions, so it touches no global setting and no shared
// fixture. The tests share one claim and run in order. Money is typed as whole numbers and asserted
// separator-tolerantly, because the currency format is a global System Setting this spec must not
// mutate.
test.use({ storageState: 'e2e/.auth/admin.json' });

test.describe('Reimbursement claims', () => {
  test.describe.configure({ mode: 'serial' });

  const claimName = uniqueName('claim-main');
  const secondClaimName = uniqueName('claim-second');
  const hotel = uniqueName('claim-hotel');
  const taxi = uniqueName('claim-taxi');
  const dinner = uniqueName('claim-dinner');

  let primary: { id: number; name: string };
  let other: { id: number; name: string };
  let taxiId: number;
  let claimId: string;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      primary = await apiCreateGroup(api, uniqueName('claim-primary'));
      other = await apiCreateGroup(api, uniqueName('claim-other'));
      const paidByUserId = await apiGetUserId(api, creds('admin').username);

      // 1000.25 + 120.50 = 1120.75, so the claimed total carries a thousands separator and cents.
      await apiCreateReceipt(api, { groupId: primary.id, paidByUserId, name: hotel, amount: '1000.25' });
      taxiId = await apiCreateReceipt(api, { groupId: primary.id, paidByUserId, name: taxi, amount: '120.50' });
      await apiCreateReceipt(api, { groupId: primary.id, paidByUserId, name: dinner, amount: '15.00' });
    });
  });

  test.afterAll(async () => {
    // Deleting a group cascades its claims, ledger entries and receipt links.
    try {
      await withAdminApi(async (api) => {
        for (const group of [primary, other]) {
          if (group) {
            await apiDeleteGroupById(api, String(group.id));
          }
        }
      });
    } catch (error) {
      console.warn(`claims teardown failed: ${error}`);
    }
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  /**
   * A money figure, tolerant of the configured symbol, its position and both separators. Pass the
   * whole part with a "," where a thousands separator falls. The leading guard and `(?!\d)` stop a
   * short value matching inside a longer one.
   */
  const money = (whole: string, cents: string) =>
    new RegExp(`(^|\\D)${whole.replace(',', '[.,\\s]?')}[.,]${cents}(?!\\d)`);

  const row = (page: Page, text: string) => page.getByRole('row').filter({ hasText: text });

  async function gotoGroupReceipts(page: Page, groupId: number): Promise<void> {
    await page.goto(`/receipts/group/${groupId}`);
    await expect(page.getByTestId('receipts-overflow-menu')).toBeVisible();
  }

  async function gotoClaim(page: Page): Promise<void> {
    await page.goto(`/claims/${claimId}/view`);
    await expect(page.getByTestId('claim-figures')).toBeVisible();
  }

  async function pickOption(page: Page, testId: string, option: string): Promise<void> {
    await page.getByTestId(testId).click();
    await page.getByRole('option', { name: option, exact: true }).click();
  }

  async function openAddToClaim(page: Page, names: string[]): Promise<void> {
    for (const name of names) {
      await row(page, name).locator('mat-checkbox').click();
    }
    await openReceiptsOverflowMenu(page);
    await page.getByTestId('receipts-add-to-claim').click();
  }

  test('creates a claim and fills it from the receipts table', async ({ page }) => {
    await page.goto(`/claims/group/${primary.id}`);
    await page.getByTestId('claim-add').click();
    await page.waitForURL(/\/claims\/add$/);

    await page.getByTestId('claim-name').getByRole('textbox').fill(claimName);
    await pickOption(page, 'claim-status-select', 'Submitted');
    await page.getByTestId('claim-submitted-date').getByRole('textbox').fill('9/20/2026');
    await page.getByTestId('claim-submit').click();

    await page.waitForURL(/\/claims\/\d+\/view$/);
    claimId = page.url().match(/\/claims\/(\d+)\/view/)![1];
    await expect(page.getByTestId('claim-status')).toHaveText('Submitted');
    // A calendar day: shown as the day entered, whatever the zones.
    await expect(page.getByTestId('claim-figure-waiting')).toContainText('Submitted Sep 20, 2026');
    await expect(page.getByTestId('claim-receipts-empty')).toBeVisible();

    await gotoGroupReceipts(page, primary.id);
    await openAddToClaim(page, [hotel, taxi]);
    const dialog = page.getByRole('dialog').filter({ hasText: 'Add to a claim' });
    // The claim select starts blank on purpose, so it is picked explicitly.
    await pickOption(page, 'add-to-claim-target', `${claimName} (Submitted)`);
    await dialog.getByTestId('dialog-submit-button').click();
    await expect(page.getByText(`Added 2 receipts to ${claimName}`)).toBeVisible();
    await expect(dialog).toBeHidden();

    await gotoClaim(page);
    const receipts = page.getByTestId('claim-receipts');
    await expect(receipts).toContainText(hotel);
    await expect(receipts).toContainText(taxi);
    const claimed = page.getByTestId('claim-figure-claimed');
    await expect(claimed).toContainText('Receipts total');
    await expect(claimed).toHaveText(money('1,120', '75'));
    await expect(page.getByTestId('claim-figure-outstanding')).toHaveText(money('1,120', '75'));
  });

  test('refuses a receipt that is already in another claim, and names it', async ({ page }) => {
    await gotoGroupReceipts(page, primary.id);
    await openAddToClaim(page, [taxi, dinner]);
    const dialog = page.getByRole('dialog').filter({ hasText: 'Add to a claim' });
    await pickOption(page, 'add-to-claim-target', 'New claim from these receipts');
    await dialog.getByTestId('add-to-claim-name').getByRole('textbox').fill(secondClaimName);
    await dialog.getByTestId('dialog-submit-button').click();

    const result = page.getByRole('dialog').filter({ hasText: 'Some receipts could not be added' });
    await expect(result.getByTestId('claim-receipts-result-summary')).toHaveText(
      `Added 1 of 2 receipts to ${secondClaimName}.`,
    );
    const failures = result.getByTestId('claim-receipts-result-failures');
    await expect(failures).toContainText(taxi);
    await expect(failures).toContainText('already in another claim');
    await expect(failures).not.toContainText(dinner);
    await result.getByRole('button', { name: 'Close' }).click();
    await expect(result).toBeHidden();

    // Refused, never moved: the first claim still holds it.
    await gotoClaim(page);
    await expect(page.getByTestId('claim-receipts')).toContainText(taxi);
  });

  test('records a payment and a recoupment, and the figures follow', async ({ page }) => {
    await gotoClaim(page);

    const addEntry = async (type: string | null, amount: string, note: string) => {
      await page.getByTestId('claim-ledger-add').click();
      const dialog = page.getByRole('dialog').filter({ hasText: 'Add ledger entry' });
      await expect(dialog).toBeVisible();
      if (type) {
        await pickOption(page, 'claim-ledger-entry-type', type);
      }
      await dialog.getByTestId('claim-ledger-entry-amount').getByRole('textbox').pressSequentially(amount);
      await dialog.getByTestId('claim-ledger-entry-note').getByRole('textbox').fill(note);
      await dialog.getByTestId('dialog-submit-button').click();
      await expect(dialog).toBeHidden();
    };

    // A new entry defaults to a payment dated today.
    await addEntry(null, '300', 'EFT 4471');
    await expect(page.getByTestId('claim-figure-paid')).toContainText(money('300', '00'));
    await expect(page.getByTestId('claim-figure-outstanding')).toHaveText(money('820', '75'));

    await addEntry('Recoupment (taken back)', '25', 'Per diem overpaid');
    // Paid is payments minus recoupments; outstanding is claimed minus paid.
    await expect(page.getByTestId('claim-figure-paid')).toContainText(money('275', '00'));
    await expect(page.getByTestId('claim-figure-outstanding')).toHaveText(money('845', '75'));

    const ledger = page.getByTestId('claim-ledger');
    await expect(ledger.getByTestId('claim-ledger-type-PAYMENT')).toBeVisible();
    await expect(ledger.getByTestId('claim-ledger-type-RECOUPMENT')).toBeVisible();
    // A recoupment reads as money going the other way.
    await expect(row(page, 'Per diem overpaid').getByTestId('claim-ledger-amount')).toHaveText(/^\s*−/);

    // The list reports the same server figures.
    await page.goto(`/claims/group/${primary.id}`);
    await expect(row(page, claimName)).toContainText(money('845', '75'));
  });

  test('refuses to delete a receipt that is in a claim', async ({ page }) => {
    await gotoGroupReceipts(page, primary.id);
    await row(page, hotel).getByTestId('receipt-delete').click();
    await page.getByRole('dialog').getByTestId('dialog-submit-button').click();
    // The interceptor toasts the server's message, which names the claim; the row stays.
    await expect(page.getByText(`remove it from claim '${claimName}' first`)).toBeVisible();
    await expect(row(page, hotel)).toBeVisible();

    await page.reload();
    await row(page, hotel).locator('mat-checkbox').click();
    await openReceiptsOverflowMenu(page);
    await page.getByTestId('receipts-bulk-delete').click();
    await page
      .getByRole('dialog')
      .filter({ hasText: 'Delete Receipts' })
      .getByTestId('dialog-submit-button')
      .click();

    const result = page.getByRole('dialog').filter({ hasText: 'Some receipts could not be deleted' });
    await expect(result.getByTestId('bulk-delete-result-failures')).toContainText(
      `remove it from claim '${claimName}' first`,
    );
    await expect(result.getByTestId('bulk-delete-result-failures')).toContainText(hotel);
    await result.getByRole('button', { name: 'Close' }).click();
    await expect(row(page, hotel)).toBeVisible();
  });

  test('a receipt moved to another group leaves its claim', async ({ page }) => {
    await page.goto(`/receipts/${taxiId}/edit`);
    const group = page.getByTestId('receipt-group');
    // A filled single-select autocomplete is readonly; clear it to change the group.
    await group.getByTestId('autocomplete-clear').getByRole('button').click();
    await group.getByRole('combobox').click();
    await page.getByRole('option', { name: other.name, exact: true }).click();
    // Changing the group clears Paid By.
    await page.getByLabel('Paid By').click();
    await page.getByRole('option').first().click();
    await page.getByRole('button', { name: 'Save', exact: true }).first().click({ force: true });
    await page.waitForURL(/\/receipts\/\d+\/view$/);

    await gotoClaim(page);
    const receipts = page.getByTestId('claim-receipts');
    await expect(receipts).toContainText(hotel);
    await expect(receipts).not.toContainText(taxi);
    await expect(page.getByTestId('claim-figure-claimed')).toHaveText(money('1,000', '25'));
    await expect(page.getByTestId('claim-figure-outstanding')).toHaveText(money('725', '25'));
  });
});
