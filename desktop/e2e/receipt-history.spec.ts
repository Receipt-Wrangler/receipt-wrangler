import { expect, test, type Locator, type Page } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiBulkStatusUpdate,
  apiCreateComment,
  apiCreateGroup,
  apiCreateReceipt,
  apiDeleteGroupById,
  apiGetGroupMembers,
  apiGetRoleId,
  apiGetUserId,
  apiReceiptHistory,
  apiSetGroupRoster,
  uniqueName,
  withAdminApi,
  withApiAs,
} from './helpers/provisioning';

// The receipt history drawer, end to end.
//
// One receipt goes through four recorded saves, each through the real API:
// created, edited (name + amount), commented on, and bulk-resolved. The comment
// and the bulk resolve are the two changes that recorded nothing before this
// feature, so their entries prove the new recording reaches the wire. The Jest
// specs build the history response by hand; this proves the drawer renders what
// the real server computes.
//
// The negative half adds e2e-user to the group as a Legacy Viewer: they can read
// the receipt (group.receipts.read) but their Legacy User app role lacks
// app.receipt-history.read, so the toggle is absent and the API refuses.

type HistoryChange = { field: string; kind: string; to?: { value?: string } };
type HistoryEntry = { systemTaskId: number; action: string; hasSnapshots: boolean; changes: HistoryChange[] };
type History = { entries: HistoryEntry[]; changeCount: number; saveCount: number; contributorCount: number };

test.use({ storageState: 'e2e/.auth/admin.json' });

test.describe('Receipt history', () => {
  test.describe.configure({ mode: 'serial' });

  const groupName = uniqueName('history-group');
  const oldName = uniqueName('history-before');
  const newName = uniqueName('history-after');
  const comment = 'Took the bag fee off, we brought our own.';
  let groupId: number;
  let receiptId: number;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      const adminId = await apiGetUserId(api, creds('admin').username);
      const userId = await apiGetUserId(api, creds('user').username);
      groupId = (await apiCreateGroup(api, groupName)).id;

      receiptId = await apiCreateReceipt(api, {
        groupId,
        paidByUserId: adminId,
        name: oldName,
        amount: '10.00',
      });

      const edit = await api.put(`/api/receipt/${receiptId}`, {
        data: {
          name: newName,
          amount: '12.50',
          date: '2024-01-01T00:00:00Z',
          groupId,
          paidByUserId: adminId,
          status: 'OPEN',
        },
      });
      expect(edit.ok()).toBe(true);

      await apiCreateComment(api, { receiptId, comment });
      await apiBulkStatusUpdate(api, { receiptIds: [receiptId], status: 'RESOLVED' });

      // e2e-user joins as a Legacy Viewer, keeping the admin's own membership.
      const members = await apiGetGroupMembers(api, groupId);
      const adminGroupRoleId = members.find((m) => m.userId === adminId)!.groupRoleId;
      const viewerRoleId = await apiGetRoleId(api, 'Legacy Viewer', 'GROUP');
      const roster = await apiSetGroupRoster(api, groupId, {
        name: groupName,
        members: [
          { userId: adminId, groupRoleId: adminGroupRoleId },
          { userId, groupRoleId: viewerRoleId },
        ],
      });
      expect(roster.status()).toBe(200);
    });
  });

  test.afterAll(async () => {
    await withAdminApi(async (api) => {
      try {
        // Cascades the receipt and frees e2e-user's group role.
        await apiDeleteGroupById(api, String(groupId));
      } catch {
        // Best effort; a leaked group does not affect other specs.
      }
    });
  });

  const readHistory = (): Promise<History> =>
    withAdminApi(async (api) => {
      const res = await apiReceiptHistory(api, receiptId);
      expect(res.status()).toBe(200);
      return (await res.json()) as History;
    });

  const openDrawer = async (page: Page): Promise<Locator> => {
    await stubTokenRefresh(page);
    await page.goto(`/receipts/${receiptId}/view`);
    await page.getByTestId('receipt-history-toggle').click();
    const drawer = page.getByTestId('receipt-history-drawer');
    await expect(drawer).toBeVisible();
    return drawer;
  };

  const entryFor = (drawer: Locator, action: string): Locator =>
    drawer.locator(`[data-testid="receipt-history-entry"][data-action="${action}"]`);

  test('the API records every save, including the comment and the bulk resolve', async () => {
    const history = await readHistory();

    expect(history.entries.map((entry) => entry.action)).toEqual([
      'RESOLVED',
      'IMAGES_COMMENTS_UPDATED',
      'UPDATED',
      'CREATED',
    ]);
    const changes = history.entries.flatMap((entry) => entry.changes);
    const commentAdded = changes.find((change) => change.field === 'COMMENT' && change.kind === 'ADDED');
    expect(commentAdded?.to?.value).toBe(comment);
    expect(history.saveCount).toBe(4);
    expect(history.contributorCount).toBe(1);
  });

  test('the drawer lists every save newest first, under the summary', async ({ page }) => {
    const history = await readHistory();
    const drawer = await openDrawer(page);

    const entries = drawer.getByTestId('receipt-history-entry');
    await expect(entries).toHaveCount(4);
    const verbs = [
      ['RESOLVED', 'resolved this receipt'],
      ['IMAGES_COMMENTS_UPDATED', 'updated images and comments'],
      ['UPDATED', 'edited this receipt'],
      ['CREATED', 'created this receipt'],
    ];
    for (const [index, [action, verb]] of verbs.entries()) {
      await expect(entries.nth(index)).toHaveAttribute('data-action', action);
      await expect(entries.nth(index).getByTestId('receipt-history-verb')).toHaveText(verb);
    }

    const counted = (count: number, word: string) => `${count} ${word}${count === 1 ? '' : 's'}`;
    await expect(drawer.getByTestId('receipt-history-summary')).toHaveText(
      [
        counted(history.changeCount, 'change'),
        counted(history.saveCount, 'save'),
        counted(history.contributorCount, 'contributor'),
      ].join(' · '),
    );
  });

  test('the edit shows the old and the new name', async ({ page }) => {
    const drawer = await openDrawer(page);

    const edited = entryFor(drawer, 'UPDATED');
    const nameRow = edited.locator('[data-testid="receipt-history-change"][data-field="NAME"]');
    await expect(nameRow).toHaveAttribute('data-kind', 'CHANGED');
    await expect(nameRow.getByTestId('receipt-history-from')).toHaveText(oldName);
    await expect(nameRow.getByTestId('receipt-history-to')).toHaveText(newName);
    await expect(edited.locator('[data-testid="receipt-history-change"][data-field="AMOUNT"]')).toHaveCount(1);
  });

  test('the comment shows as an added comment', async ({ page }) => {
    const drawer = await openDrawer(page);

    const commentRow = entryFor(drawer, 'IMAGES_COMMENTS_UPDATED').locator(
      '[data-testid="receipt-history-change"][data-field="COMMENT"]',
    );
    await expect(commentRow).toHaveAttribute('data-kind', 'ADDED');
    await expect(commentRow.getByTestId('receipt-history-to')).toContainText(comment);
    await expect(commentRow.getByTestId('receipt-history-from')).toHaveCount(0);
  });

  test('View JSON loads the snapshots behind an entry', async ({ page }) => {
    const drawer = await openDrawer(page);
    const edited = entryFor(drawer, 'UPDATED');
    const toggle = edited.getByTestId('receipt-history-json-toggle');

    const snapshots = page.waitForResponse((res) => /\/api\/receipt\/\d+\/history\/\d+\/snapshots$/.test(res.url()));
    await toggle.click();
    expect((await snapshots).status()).toBe(200);

    const json = edited.getByTestId('receipt-history-json');
    await expect(json).toContainText(`"name": "${newName}"`);
    await expect(json).toContainText(`"name": "${oldName}"`);
    await expect(toggle).toHaveText('Hide JSON');

    await toggle.click();
    await expect(json).toHaveCount(0);
    await expect(toggle).toHaveText('View JSON');
  });

  test('Escape closes the drawer and returns focus to the toggle', async ({ page }) => {
    const drawer = await openDrawer(page);
    // Opening the drawer moves focus into it.
    await expect(drawer.getByRole('button', { name: 'Close receipt history' })).toBeFocused();

    await page.keyboard.press('Escape');

    await expect(drawer).toHaveCount(0);
    await expect(page.getByRole('button', { name: 'Receipt history', exact: true })).toBeFocused();
  });

  test.describe('without the permission', () => {
    test.use({ storageState: 'e2e/.auth/user.json' });

    test('a Legacy User sees the receipt but no history toggle', async ({ page }) => {
      await stubTokenRefresh(page);
      await page.goto(`/receipts/${receiptId}/view`);

      // The header's view-mode buttons have rendered, so the toggle's absence is real.
      await expect(page.getByTestId('receipt-duplicate')).toBeVisible();
      await expect(page.getByTestId('receipt-history-toggle')).toHaveCount(0);
    });

    test('the API refuses a Legacy User the history and its snapshots', async () => {
      const history = await readHistory();
      const taskId = history.entries.find((entry) => entry.hasSnapshots)!.systemTaskId;

      await withApiAs('user', async (api) => {
        expect((await apiReceiptHistory(api, receiptId)).status()).toBe(403);
        const snapshots = await api.get(`/api/receipt/${receiptId}/history/${taskId}/snapshots`);
        expect(snapshots.status()).toBe(403);
      });
    });
  });
});
