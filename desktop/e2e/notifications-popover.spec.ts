import { expect, test } from '@playwright/test';
import { rmSync } from 'node:fs';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateGroup,
  apiCreateReceipt,
  apiDeleteGroupById,
  apiDeleteUserByName,
  apiGetRoleId,
  apiGetUserId,
  apiSetGroupRoster,
  createUserWithRole,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';

// The header bell's notifications panel is a CDK connected overlay
// (header.component). The unit spec drives it under jsdom, which cannot tell
// whether a click outside the panel still reaches what is under it, and has no
// real focus or Tab order. This spec covers those against a real notification.
//
// The recipient is a provisioned account so its notification list holds exactly
// what this spec put there: shared e2e accounts accumulate notifications from
// every spec that creates a receipt in one of their groups.

const NEW_USER_PASSWORD = `${uniqueName('pw')}-Aa1!`;
const AUTH_FILE = 'e2e/.auth/notifications-popover-user.json';

test.describe('Header notifications popover', () => {
  // One provisioned account for the whole file, so one worker.
  test.describe.configure({ mode: 'serial' });
  test.use({ storageState: AUTH_FILE });

  let username: string;
  let groupId: number;
  let receiptName: string;

  test.beforeAll(async ({ browser }) => {
    username = uniqueName('notif-user');
    receiptName = uniqueName('notif-receipt');

    const admin = await browser.newContext({ storageState: 'e2e/.auth/admin.json' });
    const adminPage = await admin.newPage();
    await stubTokenRefresh(adminPage);
    await createUserWithRole(adminPage, {
      username,
      password: NEW_USER_PASSWORD,
      role: 'Legacy User',
    });
    await admin.close();

    // Saved before any API call: request contexts opened inside a test inherit
    // the describe's storageState, so AUTH_FILE has to exist first.
    const userContext = await browser.newContext({ storageState: undefined });
    const userPage = await userContext.newPage();
    await userPage.goto('/auth/login');
    await userPage.getByLabel('Username').fill(username);
    await userPage.getByLabel('Password').fill(NEW_USER_PASSWORD);
    await userPage.getByRole('button', { name: 'Login' }).click();
    await expect(userPage).toHaveURL(/\/dashboard\/group\/\d+/, { timeout: 15_000 });
    await userContext.storageState({ path: AUTH_FILE });
    await userContext.close();

    // A receipt uploaded to a group notifies every member except its creator,
    // so the admin creates one in a group the new user belongs to.
    await withAdminApi(async (api) => {
      const group = await apiCreateGroup(api, uniqueName('notif-group'));
      groupId = group.id;
      const adminId = await apiGetUserId(api, creds('admin').username);
      const userId = await apiGetUserId(api, username);
      const roster = await apiSetGroupRoster(api, groupId, {
        name: group.name,
        members: [
          { userId: adminId, groupRoleId: await apiGetRoleId(api, 'Legacy Owner', 'GROUP') },
          { userId, groupRoleId: await apiGetRoleId(api, 'Legacy Viewer', 'GROUP') },
        ],
      });
      expect(roster.ok(), 'group roster saved').toBeTruthy();
      await apiCreateReceipt(api, { groupId, paidByUserId: adminId, name: receiptName });
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        await apiDeleteGroupById(api, String(groupId));
        await apiDeleteUserByName(api, username);
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a teardown error.
    }
    rmSync(AUTH_FILE, { force: true });
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
    await page.goto('/');
    await expect(page).toHaveURL(/\/dashboard\/group\/\d+/);
  });

  test('the bell opens the panel on the notification, and toggles it shut', async ({ page }) => {
    const bell = page.getByTestId('header-notifications-button');
    const panel = page.getByTestId('header-notifications-panel');

    await bell.locator('button').click();
    await expect(panel).toBeVisible();
    await expect(panel.getByText(receiptName)).toBeVisible();
    await expect(bell.locator('button')).toHaveAttribute('aria-expanded', 'true');
    // Opened with the mouse, the focused panel shows no ring, as the popover never did.
    await expect(panel).toBeFocused();
    await expect(panel).toHaveCSS('outline-style', 'none');

    await bell.locator('button').click();
    await expect(panel).toBeHidden();
    await expect(bell.locator('button')).toHaveAttribute('aria-expanded', 'false');
  });

  test('a click elsewhere closes the panel and still lands', async ({ page }) => {
    const panel = page.getByTestId('header-notifications-panel');

    await page.getByTestId('header-notifications-button').locator('button').click();
    await expect(panel).toBeVisible();

    // No backdrop: the first click on another header control both closes the
    // panel and does what that control does.
    await page.getByTestId('header-receipts-link').click();
    await expect(panel).toBeHidden();
    await expect(page).toHaveURL(/\/receipts\/group\/\d+/);
  });

  test('keyboard: focus moves in, Tab stays in, Escape returns to the bell', async ({ page }) => {
    const bellButton = page.getByTestId('header-notifications-button').locator('button');
    const panel = page.getByTestId('header-notifications-panel');

    await bellButton.focus();
    await page.keyboard.press('Enter');
    await expect(panel).toBeVisible();
    await expect(panel).toBeFocused();
    // Opened from the keyboard, the focused panel shows a ring.
    await expect(panel).toHaveCSS('outline-style', 'solid');
    // The list loads after the panel opens. Until it has, there is nothing to
    // Tab to, and focus rightly stays on the panel itself.
    await expect(panel.getByText(receiptName)).toBeVisible();

    // The overlay renders at the end of <body>; without the trap, Tab would
    // leave the panel for whatever follows it in the document.
    for (let i = 0; i < 4; i++) {
      await page.keyboard.press('Tab');
      await expect(panel.locator(':focus')).toHaveCount(1);
    }

    // Tab 3 reached the notification's delete button, which shows a tooltip on
    // keyboard focus; Tab 4 moved on and it fades out. While it is still up,
    // Material's tooltip takes Escape to close itself, so wait for it to go.
    // (Its accessible name comes from aria-label; the visible text is the
    // tooltip's.)
    await expect(page.getByText('Cancel', { exact: true })).toBeHidden();
    await page.keyboard.press('Escape');
    await expect(panel).toBeHidden();
    await expect(bellButton).toBeFocused();
  });

  // Last on purpose: it empties this account's notifications.
  test('keyboard: once the list is emptied, Tab stays on the panel', async ({ page }) => {
    const bellButton = page.getByTestId('header-notifications-button').locator('button');
    const panel = page.getByTestId('header-notifications-panel');

    await bellButton.focus();
    await page.keyboard.press('Enter');
    await expect(panel).toBeFocused();

    const markAllRead = panel.getByRole('button', { name: 'Mark all as read' });
    // Enabled once the list has loaded with something to mark read.
    await expect(markAllRead).toBeEnabled();
    await page.keyboard.press('Tab');
    await expect(markAllRead).toBeFocused();
    await page.keyboard.press('Enter');
    await expect(panel.getByText('All caught up!')).toBeVisible();
    await expect(markAllRead).toBeDisabled();

    // Nothing in the panel is tabbable now, so the focus trap alone would let
    // Tab walk out to the page behind it.
    for (let i = 0; i < 2; i++) {
      await page.keyboard.press('Tab');
      await expect(panel).toBeFocused();
    }
  });
});
