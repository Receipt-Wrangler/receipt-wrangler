import { BrowserContext, expect, Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateComment,
  apiCreateReceipt,
  apiDeleteGroupById,
  apiDeleteRoleByName,
  apiGetUserId,
  createGroupWithMember,
  createRole,
  uniqueName,
  withAdminApi,
  withApiAs,
} from './helpers/provisioning';

// Comment controls in the receipt form are permission-gated:
//   - the "Leave a Comment" composer renders only when the member holds
//     group.comments.create (canCreateComments()),
//   - the per-comment delete button renders only for the member's OWN comment
//     AND when they hold group.comments.delete (canDeleteComments()).
// Both controls only appear in EDIT mode (the receipt /:id/edit route, gated on
// group.receipts.update) — the read-only /view mode never shows them. So each
// fixture role starts from the "Receipt Editor" preset (which grants
// group.receipts.update + read + group.comments.create, but NOT
// group.comments.delete):
//   - the CREATE-deny member drops group.comments.create from that preset,
//   - the DELETE-deny member keeps the preset as-is (it already lacks delete),
//     and a comment owned by e2e-user is seeded via the API so the delete
//     control would otherwise render.
// Each member lives in its own group with its own seeded receipt; assertions
// run as e2e-user. A positive contrast (unmodified Receipt Editor → composer
// visible) proves the create gating is selective.

test.describe('Receipt comment create/delete gating', () => {
  test.describe.configure({ mode: 'serial' });

  let adminContext: BrowserContext;
  let adminPage: Page;

  // CREATE-deny fixture: Receipt Editor minus group.comments.create.
  const noCreateRole = uniqueName('no-comment-create-role');
  const noCreateGroup = uniqueName('no-comment-create-grp');
  let noCreateGroupId: string;
  let noCreateReceiptId: number;

  // CREATE-allow contrast fixture: unmodified Receipt Editor.
  const editorRole = uniqueName('comment-editor-role');
  const editorGroup = uniqueName('comment-editor-grp');
  let editorGroupId: string;
  let editorReceiptId: number;

  // DELETE-deny fixture: Receipt Editor (already lacks group.comments.delete)
  // with a comment owned by e2e-user seeded via the API.
  const noDeleteRole = uniqueName('no-comment-delete-role');
  const noDeleteGroup = uniqueName('no-comment-delete-grp');
  let noDeleteGroupId: string;
  let noDeleteReceiptId: number;

  test.beforeAll(async ({ browser }) => {
    adminContext = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    adminPage = await adminContext.newPage();
    await stubTokenRefresh(adminPage);

    // Receipt Editor minus the comment-create permission.
    await createRole(adminPage, {
      name: noCreateRole,
      type: 'Group role',
      preset: 'Receipt Editor',
      disablePermissions: [
        { panelKey: 'group.comments', label: 'Create Comments' },
      ],
    });
    noCreateGroupId = await createGroupWithMember(adminPage, {
      groupName: noCreateGroup,
      memberDisplayName: 'E2E User',
      roleName: noCreateRole,
    });

    // Unmodified Receipt Editor — holds comment-create (positive contrast).
    await createRole(adminPage, {
      name: editorRole,
      type: 'Group role',
      preset: 'Receipt Editor',
    });
    editorGroupId = await createGroupWithMember(adminPage, {
      groupName: editorGroup,
      memberDisplayName: 'E2E User',
      roleName: editorRole,
    });

    // Receipt Editor — holds create + update + read but NOT delete.
    await createRole(adminPage, {
      name: noDeleteRole,
      type: 'Group role',
      preset: 'Receipt Editor',
    });
    noDeleteGroupId = await createGroupWithMember(adminPage, {
      groupName: noDeleteGroup,
      memberDisplayName: 'E2E User',
      roleName: noDeleteRole,
    });

    // Seed one receipt per group (paid by e2e-user so it's plainly visible),
    // and a comment owned by e2e-user on the delete-fixture receipt.
    await withAdminApi(async (api) => {
      const userId = await apiGetUserId(api, creds('user').username);
      noCreateReceiptId = await apiCreateReceipt(api, {
        groupId: noCreateGroupId,
        paidByUserId: userId,
        name: uniqueName('no-create-rcpt'),
      });
      editorReceiptId = await apiCreateReceipt(api, {
        groupId: editorGroupId,
        paidByUserId: userId,
        name: uniqueName('editor-rcpt'),
      });
      noDeleteReceiptId = await apiCreateReceipt(api, {
        groupId: noDeleteGroupId,
        paidByUserId: userId,
        name: uniqueName('no-delete-rcpt'),
      });
    });

    // Seed the comment AS e2e-user so it's owned by them (the delete control
    // only renders for the logged-in user's own comments).
    await withApiAs('user', async (api) => {
      await apiCreateComment(api, {
        receiptId: noDeleteReceiptId,
        comment: 'e2e-user own comment',
      });
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        // Group delete frees the members' role assignment; then the roles delete.
        await apiDeleteGroupById(api, noCreateGroupId);
        await apiDeleteGroupById(api, editorGroupId);
        await apiDeleteGroupById(api, noDeleteGroupId);
        await apiDeleteRoleByName(api, noCreateRole, 'GROUP');
        await apiDeleteRoleByName(api, editorRole, 'GROUP');
        await apiDeleteRoleByName(api, noDeleteRole, 'GROUP');
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
    await adminContext?.close();
  });

  // Assertions run as the default project user (e2e-user = the member).
  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('without group.comments.create the comment composer is absent', async ({
    page,
  }) => {
    await page.goto(`/receipts/${noCreateReceiptId}/edit`);
    // Edit mode loads (the member holds group.receipts.update).
    await expect(page.getByLabel('Name')).toBeVisible();
    // The "Leave a Comment" composer is gated on group.comments.create.
    await expect(page.getByText('Leave a Comment')).toHaveCount(0);
  });

  test('with group.comments.create the comment composer renders (contrast)', async ({
    page,
  }) => {
    await page.goto(`/receipts/${editorReceiptId}/edit`);
    await expect(page.getByLabel('Name')).toBeVisible();
    // The Receipt Editor preset grants group.comments.create.
    await expect(page.getByText('Leave a Comment')).toBeVisible();
  });

  test('without group.comments.delete the own-comment delete control is absent', async ({
    page,
  }) => {
    await page.goto(`/receipts/${noDeleteReceiptId}/edit`);
    await expect(page.getByLabel('Name')).toBeVisible();
    // The seeded comment (owned by e2e-user) is shown...
    await expect(page.getByText('e2e-user own comment')).toBeVisible();
    // ...but its delete control is gated on group.comments.delete (not held).
    await expect(page.getByTestId('comment-delete')).toHaveCount(0);
  });
});
