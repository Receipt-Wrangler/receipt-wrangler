import { expect, type Page, test } from '@playwright/test';
import { createHash } from 'node:crypto';
import { readFileSync } from 'node:fs';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateGroup,
  apiCreateReceipt,
  apiDeleteGroupById,
  apiGetUserId,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import { RECEIPT_PNG } from './helpers/quick-scan';
import { openReceiptsOverflowMenu } from './helpers/receipts-table';
import { readZipEntries } from './helpers/zip';

// Reimbursement claims, end to end (see desktop/CLAUDE.md → "Claims"). The Jest specs mock
// ReimbursementClaimService and the Go tests drive the handlers directly, so neither proves the two
// halves agree on the wire. This walks the main journey through the real UI and API:
//
//   - a claim created on the claim form, then filled from the receipts table's bulk action;
//   - a receipt already in another claim named by the bulk dialog before submitting, then refused
//     and named, not moved;
//   - the claim page's receipt picker leaving out a receipt that is in another claim;
//   - a payment and a recoupment moving the server's figures;
//   - the claim packet downloaded and opened: its PDF, CSVs, the original image byte for byte, and
//     a manifest whose hashes match the files;
//   - the receipt delete endpoints refusing a receipt that is in a claim, singly and in bulk;
//   - a receipt moved to another group on the receipt form: its claim shown, a warning on picking
//     the new group, and the receipt leaving its claim on save.
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
  // Never put on a claim: the picker's proof that an unclaimed receipt is still offered.
  const lunch = uniqueName('claim-lunch');

  let primary: { id: number; name: string };
  let other: { id: number; name: string };
  let taxiId: number;
  const receiptImage = readFileSync(RECEIPT_PNG);
  let claimId: string;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      primary = await apiCreateGroup(api, uniqueName('claim-primary'));
      other = await apiCreateGroup(api, uniqueName('claim-other'));
      const paidByUserId = await apiGetUserId(api, creds('admin').username);

      // 1000.25 + 120.50 = 1120.75, so the claimed total carries a thousands separator and cents.
      const hotelId = await apiCreateReceipt(api, { groupId: primary.id, paidByUserId, name: hotel, amount: '1000.25' });
      // One image, so the packet carries an original to compare byte for byte.
      const upload = await api.post('/api/receiptImage/', {
        multipart: {
          file: { name: 'receipt.png', mimeType: 'image/png', buffer: receiptImage },
          receiptId: String(hotelId),
        },
      });
      expect(upload.ok(), `image upload: HTTP ${upload.status()}`).toBe(true);
      taxiId = await apiCreateReceipt(api, { groupId: primary.id, paidByUserId, name: taxi, amount: '120.50' });
      await apiCreateReceipt(api, { groupId: primary.id, paidByUserId, name: dinner, amount: '15.00' });
      await apiCreateReceipt(api, { groupId: primary.id, paidByUserId, name: lunch, amount: '9.00' });
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
    // Named up front, from the claimId/claimName the receipts list carries.
    const alreadyClaimed = dialog.getByTestId('add-to-claim-already-claimed');
    await expect(alreadyClaimed).toContainText(`${claimName}: ${taxi}`);
    await expect(alreadyClaimed).not.toContainText(dinner);
    await pickOption(page, 'add-to-claim-target', 'New claim from these receipts');
    await dialog.getByTestId('add-to-claim-name').getByRole('textbox').fill(secondClaimName);
    // Still sent, not dropped: the server's refusal is the record.
    await dialog.getByTestId('dialog-submit-button').click();

    const result = page.getByRole('dialog').filter({ hasText: 'Some receipts could not be added' });
    await expect(result.getByTestId('claim-receipts-result-summary')).toHaveText(
      `Added 1 of 2 receipts to ${secondClaimName}.`,
    );
    const failures = result.getByTestId('claim-receipts-result-failures');
    await expect(failures).toContainText(taxi);
    // The server names a claim the caller may see.
    await expect(failures).toContainText(`already in claim '${claimName}'`);
    await expect(failures).not.toContainText(dinner);
    await result.getByRole('button', { name: 'Close' }).click();
    await expect(result).toBeHidden();

    // Refused, never moved: the first claim still holds it.
    await gotoClaim(page);
    await expect(page.getByTestId('claim-receipts')).toContainText(taxi);
  });

  test("the claim page's picker leaves out a receipt that is in another claim", async ({ page }) => {
    await gotoClaim(page);
    await page.getByTestId('claim-receipts-add').click();
    const dialog = page.getByRole('dialog').filter({ hasText: 'Add receipts' });
    await expect(dialog).toBeVisible();

    // Dinner went into the second claim above; it is left out and counted.
    await expect(dialog.getByTestId('claim-receipt-picker-other-claims')).toHaveText(
      '1 receipt already in another claim is not listed.',
    );
    const input = dialog.getByTestId('claim-receipt-picker').getByRole('combobox');
    await input.click();
    await expect(page.getByRole('option').filter({ hasText: lunch })).toBeVisible();
    await expect(page.getByRole('option').filter({ hasText: dinner })).toHaveCount(0);
    // Searching for it by name finds nothing either.
    await input.fill(dinner);
    await expect(page.getByRole('option')).toHaveCount(0);

    await dialog.getByTestId('dialog-cancel-button').click();
    await expect(dialog).toBeHidden();
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

  test('downloads the claim packet, and it holds what it says', async ({ page }) => {
    await gotoClaim(page);

    await page.getByTestId('claim-download').click();
    const packetEntry = page.getByTestId('claim-download-packet');
    await expect(packetEntry).toContainText('SHA-256 manifest');
    // The plain CSV export is the menu's other entry.
    await expect(page.getByTestId('claim-export')).toBeVisible();

    const downloading = page.waitForEvent('download');
    await packetEntry.click();
    const download = await downloading;

    // The server's name: the claim name reduced to letters and digits, then today's date.
    const safeName = claimName.replace(/[^A-Za-z0-9]+/g, '_');
    expect(download.suggestedFilename()).toMatch(new RegExp(`^${safeName}_packet_\\d{4}-\\d{2}-\\d{2}\\.zip$`));
    await expect(page.getByText(`Claim packet downloaded: ${download.suggestedFilename()}`)).toBeVisible();
    await expect(page.getByTestId('claim-packet-busy')).toHaveCount(0);

    const zip = readFileSync((await download.path())!);
    const files = readZipEntries(zip);
    for (const name of ['packet.pdf', 'manifest.json', 'receipts.csv', 'items.csv']) {
      expect([...files.keys()]).toContain(name);
    }
    expect(files.get('packet.pdf')!.subarray(0, 5).toString('latin1')).toBe('%PDF-');
    expect(files.get('receipts.csv')!.toString('utf8')).toContain(hotel);

    const manifest = JSON.parse(files.get('manifest.json')!.toString('utf8')) as {
      files: { path: string; size: number; sha256: string }[];
    };
    expect(Array.isArray(manifest.files)).toBe(true);
    expect(manifest.files.length).toBeGreaterThan(0);
    const sha256 = (bytes: Buffer) => createHash('sha256').update(bytes).digest('hex');
    for (const file of manifest.files) {
      expect(file.sha256).toMatch(/^[0-9a-f]{64}$/);
      // Every hash is over the bytes as written into the zip.
      expect(sha256(files.get(file.path)!), file.path).toBe(file.sha256);
    }

    // The hotel's image is the original upload, untouched.
    const images = manifest.files.filter((file) => file.path.startsWith('images/'));
    expect(images).toHaveLength(1);
    expect(files.get(images[0].path)!.equals(receiptImage)).toBe(true);
    expect(images[0].sha256).toBe(sha256(receiptImage));
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
    // The receipt's claim, in the audit details, linking to the claim.
    await expect(page.getByTestId('receipt-claim')).toContainText(`In claim: ${claimName}`);
    await expect(page.getByTestId('receipt-claim').getByRole('link')).toHaveAttribute(
      'href',
      `/claims/${claimId}/view`,
    );
    const warning = page.getByTestId('receipt-claim-move-warning');
    await expect(warning).toHaveCount(0);

    const group = page.getByTestId('receipt-group');
    // A filled single-select autocomplete is readonly; clear it to change the group.
    await group.getByTestId('autocomplete-clear').getByRole('button').click();
    await group.getByRole('combobox').click();
    await page.getByRole('option', { name: other.name, exact: true }).click();
    // Warned, not blocked.
    await expect(warning).toHaveText(
      `Moving this receipt to another group will remove it from claim '${claimName}'.`,
    );
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
