import { BrowserContext, expect, Page, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';

// A group VIEWER (low-tier member) must be denied the group-scoped edit actions
// — editing the group, its receipt settings, and creating receipts — while
// still able to VIEW the group (proving the denial is selective, not a blanket
// lockout). The rest of the suite runs as a Legacy Admin or as e2e-user acting
// as a group OWNER, so this is the first coverage of the group-scoped *deny*
// path (a Viewer/Editor tier being held back from what an Owner can do).
//
// e2e-user has no low-tier membership by default (creating a group makes them
// the Owner), so an admin browser context provisions a group and adds e2e-user
// as a "Legacy Viewer" through the real group-member form — the same flow
// exercised by role-assignment.spec.ts. The assertions then run as the default
// project user (e2e-user = the Viewer) against that group.

function uniqueName(tag: string) {
  return `e2e-${tag}-${Date.now()}`;
}

// Deletes a group from the list by name. The delete button is an icon-only
// control (data-testid), as is the confirmation dialog's submit. Mirrors the
// helper in role-defaults.spec.ts.
async function deleteGroupByName(page: Page, name: string) {
  await page.goto('/groups');
  const row = page.getByRole('row').filter({ hasText: name }).first();
  await expect(row).toBeVisible();
  await row.getByTestId('group-delete').click();

  const dialog = page.getByRole('dialog');
  await expect(dialog).toBeVisible();
  await dialog.getByTestId('dialog-submit-button').click();

  await expect(page.getByRole('row').filter({ hasText: name })).toHaveCount(0);
}

test.describe('Group Viewer is denied group/receipt edits', () => {
  // Provision the viewer group once, then run the independent assertions
  // against it in order.
  test.describe.configure({ mode: 'serial' });

  let adminContext: BrowserContext;
  let adminPage: Page;
  let viewerGroupId: string;
  let groupName: string;

  test.beforeAll(async ({ browser }) => {
    groupName = uniqueName('viewer-grp');
    adminContext = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    adminPage = await adminContext.newPage();
    await stubTokenRefresh(adminPage);

    // Create the group and add e2e-user (display name "E2E User") as a Viewer.
    await adminPage.goto('/groups/create');
    await expect(adminPage.getByLabel('Group Name')).toBeVisible();
    await adminPage.getByLabel('Group Name').fill(groupName);

    await adminPage.getByTestId('add-group-member').click();
    const memberDialog = adminPage
      .getByRole('dialog')
      .filter({ hasText: 'Add Group Member' });
    await expect(memberDialog).toBeVisible();

    // The user autocomplete filters/displays by display name. Target the
    // combobox role: while the option panel is open the listbox shares the
    // field's aria label, so getByLabel would match two elements.
    const userField = memberDialog.getByRole('combobox', { name: 'User' });
    await userField.fill('E2E User');
    await adminPage
      .getByRole('option', { name: 'E2E User', exact: true })
      .click();

    // Role defaults to "Legacy Owner" — switch it to "Legacy Viewer".
    const roleSelect = memberDialog.getByRole('combobox', { name: 'Role' });
    await roleSelect.click();
    await adminPage
      .getByRole('option', { name: 'Legacy Viewer', exact: true })
      .click();

    await memberDialog.getByTestId('dialog-submit-button').click();
    await expect(memberDialog).toBeHidden();

    // Submit the group form. The page-level save is an icon-only button; the
    // app-form submits on Enter from a field (as in role-defaults.spec.ts).
    // The owner check is skipped in add-mode, and the backend assigns the
    // creating admin the default (Owner) role, so the group keeps an owner.
    await adminPage.getByLabel('Group Name').press('Enter');
    await adminPage.waitForURL(/\/groups\/\d+\/details\/view/);
    viewerGroupId = adminPage.url().match(/\/groups\/(\d+)\//)![1];
  });

  test.afterAll(async () => {
    try {
      await deleteGroupByName(adminPage, groupName);
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
    await adminContext?.close();
  });

  // Assertions run as the default project user (e2e-user = the Viewer).
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('cannot open the group edit route', async ({ page }) => {
    await page.goto(`/groups/${viewerGroupId}/details/edit`);
    // groupPermissionGuard denies (group.update) and redirects to the dashboard.
    await expect(page).toHaveURL(/\/dashboard\/group\/\d+/);
  });

  test('cannot open the receipt-settings edit route', async ({ page }) => {
    await page.goto(`/groups/${viewerGroupId}/receipt-settings/edit`);
    await expect(page).toHaveURL(/\/dashboard\/group\/\d+/);
  });

  test('group row exposes no edit or delete action', async ({ page }) => {
    await page.goto('/groups');
    const row = page.getByRole('row').filter({ hasText: groupName }).first();
    await expect(row).toBeVisible();
    // Neither gate (group.update / group.delete) is held for this group.
    await expect(row.getByTestId('group-edit')).toHaveCount(0);
    await expect(row.getByTestId('group-delete')).toHaveCount(0);
  });

  test('cannot add a receipt in the group', async ({ page }) => {
    await page.goto(`/receipts/group/${viewerGroupId}`);
    // Add Receipt is gated on group.receipts.* which a Viewer lacks.
    await expect(
      page.getByRole('button', { name: 'Add Receipt' }),
    ).toHaveCount(0);
  });

  test('can still view the group (denial is selective)', async ({ page }) => {
    await page.goto(`/groups/${viewerGroupId}/details/view`);
    await expect(page).toHaveURL(
      new RegExp(`/groups/${viewerGroupId}/details/view`),
    );
    // A Viewer holds group.view — the read path is allowed, only writes deny.
    await expect(page.getByLabel('Group Name')).toHaveValue(groupName);
  });
});
