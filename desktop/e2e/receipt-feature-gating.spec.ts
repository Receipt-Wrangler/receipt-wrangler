import { BrowserContext, expect, Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateReceipt,
  apiDeleteGroupById,
  apiDeleteRoleByName,
  apiGetUserId,
  createGroupWithMember,
  createRole,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';

// The AI-powered receipt features are permission-gated in the UI:
//   - Quick Scan  -> group.receipts.quick-scan (sidebar + receipts-table header),
//   - Poll Email  -> group.email.poll (receipts-table header),
//   - Magic Fill  -> group.receipts.magic-fill (receipt form image toolbar).
// A "Viewer" group role holds none of these, so an admin context provisions a
// Viewer group with e2e-user + a seeded receipt and asserts (as e2e-user) that
// none of the three controls render.
//
// IMPORTANT CONFOUND: all three controls ALSO sit behind the `aiPoweredReceipts`
// feature flag (`*appFeature="'aiPoweredReceipts'"` for Quick Scan/Poll Email,
// and `aiPoweredReceipts()` for Magic Fill). That flag is server-config
// (FeatureConfig) and is currently `false` in the dev/CI API, so these controls
// are hidden for EVERYONE regardless of permission. These tests therefore prove
// the permission-denied member sees no control (the negative axis), but they
// cannot prove the positive contrast (a permitted member WOULD see it) without
// enabling the feature flag on the API — which this harness must not modify. The
// positive-contrast tests below are left as test.fixme for that reason.

test.describe('Receipt feature gating (quick-scan / poll-email / magic-fill)', () => {
  test.describe.configure({ mode: 'serial' });

  let adminContext: BrowserContext;
  let adminPage: Page;
  const roleName = uniqueName('no-ai-role');
  const groupName = uniqueName('no-ai-grp');
  let groupId: string;
  let receiptId: number;

  test.beforeAll(async ({ browser }) => {
    adminContext = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    adminPage = await adminContext.newPage();
    await stubTokenRefresh(adminPage);

    // Viewer: holds group.receipts.read (so the list/form load) but none of
    // quick-scan / magic-fill / email.poll.
    await createRole(adminPage, {
      name: roleName,
      type: 'Group role',
      preset: 'Viewer',
    });

    groupId = await createGroupWithMember(adminPage, {
      groupName,
      memberDisplayName: 'E2E User',
      roleName,
    });

    await withAdminApi(async (api) => {
      const userId = await apiGetUserId(api, creds('user').username);
      receiptId = await apiCreateReceipt(api, {
        groupId,
        paidByUserId: userId,
        name: uniqueName('no-ai-rcpt'),
      });
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        await apiDeleteGroupById(api, groupId);
        await apiDeleteRoleByName(api, roleName, 'GROUP');
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
    await adminContext?.close();
  });

  // Assertions run as the default project user (e2e-user = the Viewer).
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('the receipts header shows no Quick Scan button', async ({ page }) => {
    await page.goto(`/receipts/group/${groupId}`);
    // The table renders for a member who holds group.receipts.read.
    await expect(page.getByTestId('configure-columns')).toBeVisible();
    // Quick Scan (group.receipts.quick-scan, also feature-flagged) is absent.
    await expect(
      page.getByRole('button', { name: 'Quick Scan' }),
    ).toHaveCount(0);
  });

  test('the receipts header shows no Poll Email button', async ({ page }) => {
    await page.goto(`/receipts/group/${groupId}`);
    await expect(page.getByTestId('configure-columns')).toBeVisible();
    // Poll Email (group.email.poll, also feature-flagged + needs email
    // integration enabled) is absent.
    await expect(
      page.getByRole('button', { name: 'Poll email(s)' }),
    ).toHaveCount(0);
  });

  test('the receipt form shows no Magic Fill button', async ({ page }) => {
    await page.goto(`/receipts/${receiptId}/view`);
    // The receipt form renders for a member who holds group.receipts.read.
    await expect(page.getByLabel('Name')).toBeVisible();
    // Magic Fill (group.receipts.magic-fill, also feature-flagged) is absent.
    await expect(page.getByRole('button', { name: 'Magic fill' })).toHaveCount(
      0,
    );
  });

  // Positive contrasts require the aiPoweredReceipts feature flag enabled on the
  // API, which this harness must not modify (FeatureConfig is server config and
  // is `false` in dev/CI). Left as fixme so the gap is explicit rather than a
  // false green: with the flag on, a member holding the permission WOULD see the
  // control while a Viewer still would not.
  test.fixme(
    'positive contrast: a member with the permission sees the control (needs aiPoweredReceipts enabled)',
    async () => {
      // Cannot run: aiPoweredReceipts is false in the dev/CI API and the harness
      // must not change API config. Enabling it would let us assert e.g. a
      // Receipt Editor (group.receipts.magic-fill) sees the Magic Fill button
      // while the Viewer does not.
    },
  );
});
