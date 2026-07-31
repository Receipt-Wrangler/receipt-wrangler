import { expect, type Locator, type Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateCategory,
  apiCreateRole,
  apiDeleteCategoryById,
  apiDeleteGroupById,
  apiDeleteRoleByName,
  apiGetGroupMembers,
  apiGetUserId,
  apiSetMemberGrants,
  createGroupWithMember,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';

// Authoring per-member category/tag assignments through the UI.
//
// The same assignment is reachable from two places — the user form (one section
// per group the user belongs to) and the group-member form (just that group's) —
// and they are the SAME membership record, written through the same endpoint.
// That equivalence is easy to break and easy to misunderstand, so it is asserted
// directly rather than assumed.

test.describe('Per-member grant assignment (UI)', () => {
  // Serial: the tests read back each other's saved assignment.
  test.describe.configure({ mode: 'serial' });
  test.use({ storageState: 'e2e/.auth/admin.json' });

  const cappedRoleName = uniqueName('grantui-capped');
  const openRoleName = uniqueName('grantui-open');
  const cappedGroupName = uniqueName('grantui-capped-grp');
  const openGroupName = uniqueName('grantui-open-grp');

  let cappedGroupId: string;
  let openGroupId: string;
  let memberId: number;
  let catA: { id: number; name: string };
  let catB: { id: number; name: string };
  let catC: { id: number; name: string };
  let catOutside: { id: number; name: string };

  const PERMISSIONS = ['group.view', 'group.receipts.read'];

  /** Opens the Edit User dialog for the e2e member. */
  async function openUserForm(page: Page): Promise<Locator> {
    await page.goto('/users');
    const row = page.getByRole('row').filter({ hasText: creds('user').username });
    await row.getByTestId('user-edit').click();

    const dialog = page.getByRole('dialog').filter({ hasText: 'Edit' });
    await expect(dialog).toBeVisible();
    return dialog;
  }

  /** The assignment block for one group inside the user form. */
  function grantSection(scope: Locator, groupId: string): Locator {
    return scope.getByTestId(`grant-section-${groupId}`);
  }

  /** Names currently selected in a picker (its chips). */
  async function selectedNames(picker: Locator): Promise<string[]> {
    const text = await picker.getByTestId('grant-picker-categories').innerText();
    return [catA, catB, catC, catOutside]
      .filter((c) => text.includes(c.name))
      .map((c) => c.name)
      .sort();
  }

  /** Opens a picker's option list and returns the offered category names. */
  async function offeredNames(page: Page, picker: Locator): Promise<string[]> {
    const field = picker
      .getByTestId('grant-picker-categories')
      .getByRole('combobox', { name: 'Categories' });
    await field.click();
    await expect(page.getByRole('listbox')).toBeVisible();

    const offered: string[] = [];
    for (const candidate of [catA, catB, catC, catOutside]) {
      const count = await page
        .getByRole('option', { name: candidate.name, exact: true })
        .count();
      if (count > 0) {
        offered.push(candidate.name);
      }
    }
    await page.keyboard.press('Escape');
    return offered.sort();
  }

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      memberId = await apiGetUserId(api, creds('user').username);
      catA = await apiCreateCategory(api, uniqueName('grantui-cat-a'));
      catB = await apiCreateCategory(api, uniqueName('grantui-cat-b'));
      catC = await apiCreateCategory(api, uniqueName('grantui-cat-c'));
      catOutside = await apiCreateCategory(api, uniqueName('grantui-cat-out'));

      // One role capped to A/B/C, one granting nothing (unrestricted).
      await apiCreateRole(api, {
        name: cappedRoleName,
        scope: 'GROUP',
        permissions: PERMISSIONS,
        categoryGrants: [catA.id, catB.id, catC.id],
      });
      await apiCreateRole(api, {
        name: openRoleName,
        scope: 'GROUP',
        permissions: PERMISSIONS,
      });
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        await apiDeleteGroupById(api, cappedGroupId);
        await apiDeleteGroupById(api, openGroupId);
        await apiDeleteRoleByName(api, cappedRoleName, 'GROUP');
        await apiDeleteRoleByName(api, openRoleName, 'GROUP');
        for (const category of [catA, catB, catC, catOutside]) {
          await apiDeleteCategoryById(api, category.id);
        }
      });
    } catch {
      // Best-effort cleanup.
    }
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('creates the two groups the member belongs to', async ({ page }) => {
    // Provisioning through the real group form; the later tests assert the user
    // form shows one section per group.
    cappedGroupId = await createGroupWithMember(page, {
      groupName: cappedGroupName,
      memberDisplayName: 'E2E User',
      roleName: cappedRoleName,
    });
    openGroupId = await createGroupWithMember(page, {
      groupName: openGroupName,
      memberDisplayName: 'E2E User',
      roleName: openRoleName,
    });

    expect(cappedGroupId).not.toBe(openGroupId);
  });

  test('the user form shows one assignment section per group', async ({ page }) => {
    const dialog = await openUserForm(page);

    await expect(dialog.getByText('Category & tag access')).toBeVisible();
    await expect(grantSection(dialog, cappedGroupId)).toBeVisible();
    await expect(grantSection(dialog, openGroupId)).toBeVisible();

    // Each block names its group and the role that governs it.
    await expect(grantSection(dialog, cappedGroupId)).toContainText(cappedGroupName);
    await expect(grantSection(dialog, cappedGroupId)).toContainText(cappedRoleName);
    await expect(grantSection(dialog, openGroupId)).toContainText(openGroupName);
  });

  test('a capped role narrows the offered pool and says why', async ({ page }) => {
    const dialog = await openUserForm(page);
    const section = grantSection(dialog, cappedGroupId);

    await expect(section.getByTestId('grant-picker-category-ceiling')).toContainText(
      `Limited to the 3 categories allowed by role ${cappedRoleName}`,
    );
    expect(await offeredNames(page, section)).toEqual(
      [catA.name, catB.name, catC.name].sort(),
    );
  });

  test('a role granting nothing offers the whole pool with no hint', async ({ page }) => {
    // Pins the "empty grants = unrestricted" convention in the UI.
    const dialog = await openUserForm(page);
    const section = grantSection(dialog, openGroupId);

    await expect(section.getByTestId('grant-picker-category-ceiling')).toHaveCount(0);
    expect(await offeredNames(page, section)).toEqual(
      expect.arrayContaining([catA.name, catOutside.name]),
    );
  });

  test('assigning a category persists it', async ({ page }) => {
    const dialog = await openUserForm(page);
    const section = grantSection(dialog, cappedGroupId);

    const field = section
      .getByTestId('grant-picker-categories')
      .getByRole('combobox', { name: 'Categories' });
    await field.click();
    await field.fill(catB.name);
    await page.getByRole('option', { name: catB.name, exact: true }).click();
    await page.keyboard.press('Escape');

    const [grantsResponse] = await Promise.all([
      page.waitForResponse(
        (r) =>
          /\/api\/group\/\d+\/member\/\d+\/grants$/.test(r.url()) &&
          r.request().method() === 'PUT',
      ),
      dialog.locator('app-submit-button button').click(),
    ]);
    expect(grantsResponse.status()).toBe(200);
    expect(grantsResponse.url()).toContain(`/group/${cappedGroupId}/member/${memberId}/`);

    const members = await withAdminApi((api) => apiGetGroupMembers(api, cappedGroupId));
    expect(members.find((m) => m.userId === memberId)?.categoryGrants).toEqual([catB.id]);
  });

  test('re-opening pre-selects the saved assignment', async ({ page }) => {
    const dialog = await openUserForm(page);
    expect(await selectedNames(grantSection(dialog, cappedGroupId))).toEqual([catB.name]);
  });

  test('saving without touching the pickers writes no grants request', async ({ page }) => {
    // Only CHANGED rows are written. A blanket rewrite would also risk tripping
    // the ceiling check on a role narrowed since the assignment was made.
    const dialog = await openUserForm(page);

    let grantsRequests = 0;
    page.on('request', (request) => {
      if (
        /\/api\/group\/\d+\/member\/\d+\/grants$/.test(request.url()) &&
        request.method() === 'PUT'
      ) {
        grantsRequests += 1;
      }
    });

    await dialog.getByLabel('Displayname').fill('E2E User');
    await Promise.all([
      page.waitForResponse(
        (r) => /\/api\/user\/\d+$/.test(r.url()) && r.request().method() === 'PUT',
      ),
      dialog.locator('app-submit-button button').click(),
    ]);

    expect(grantsRequests).toBe(0);
  });

  test('the Create User form has no assignment section', async ({ page }) => {
    // Grants hang off a MEMBERSHIP; a user being created has none yet.
    await page.goto('/users');
    await page.getByTestId('user-add').click();

    const dialog = page.getByRole('dialog').filter({ hasText: 'Create User' });
    await expect(dialog).toBeVisible();
    await expect(dialog.getByText('Category & tag access')).toHaveCount(0);
  });

  test('the group-member form edits the SAME assignment as the user form', async ({
    page,
  }) => {
    // The "one record, two doors" guarantee: the group side must show what the
    // user side saved, and vice versa.
    // The roster is only editable on the edit route; /details/view is read-only.
    await page.goto(`/groups/${cappedGroupId}/details/edit`);
    await expect(page.getByTestId('group-member-edit').first()).toBeVisible();
    const row = page.getByRole('row').filter({ hasText: 'E2E User' });
    await row.getByTestId('group-member-edit').click();

    const memberDialog = page.getByRole('dialog').filter({ hasText: 'Edit Group Member' });
    await expect(memberDialog).toBeVisible();
    const section = memberDialog.getByTestId('member-grant-section');
    await expect(section).toBeVisible();

    // The dialog must be for the member we clicked, not another row's.
    await expect(memberDialog.getByRole('combobox', { name: 'User' })).toHaveValue('E2E User');

    // The section only renders once every async input has landed (see
    // GroupMemberFormComponent.grantsReady), so it is stable as soon as it appears.
    await expect(section.getByTestId('grant-picker-category-ceiling')).toBeVisible();

    // It shows what the USER form saved earlier.
    await expect(section).toContainText(catB.name);

    // Change it here, and the user form must reflect the change.
    const field = section
      .getByTestId('grant-picker-categories')
      .getByRole('combobox', { name: 'Categories' });
    await field.click();
    await field.fill(catC.name);
    await page.getByRole('option', { name: catC.name, exact: true }).click();
    await page.keyboard.press('Escape');

    // dispatchEvent rather than click(): adding a chip grows the dialog, and
    // MatDialog re-centres itself as it does, so the footer button never settles
    // long enough to satisfy the actionability check. The button's presence is
    // already asserted by resolving the locator.
    await memberDialog
      .getByTestId('dialog-submit-button')
      .locator('button')
      .dispatchEvent('click');
    await expect(memberDialog).toBeHidden();

    await expect
      .poll(
        async () => {
          const members = await withAdminApi((api) =>
            apiGetGroupMembers(api, cappedGroupId),
          );
          return [...(members.find((m) => m.userId === memberId)?.categoryGrants ?? [])].sort();
        },
        { timeout: 15_000 },
      )
      .toEqual([catB.id, catC.id].sort());

    const userDialog = await openUserForm(page);
    expect(await selectedNames(grantSection(userDialog, cappedGroupId))).toEqual(
      [catB.name, catC.name].sort(),
    );
  });

  test('the Add Group Member form has no assignment section', async ({ page }) => {
    // The membership does not exist until the group form is saved.
    await page.goto(`/groups/${cappedGroupId}/details/edit`);
    await page.getByTestId('add-group-member').click();

    const dialog = page.getByRole('dialog').filter({ hasText: 'Add Group Member' });
    await expect(dialog).toBeVisible();
    await expect(dialog.getByTestId('member-grant-section')).toHaveCount(0);
  });

  test('the role form authors a ceiling and rehydrates it', async ({ page }) => {
    // The shared picker is also the role editor's grant control; an untouched
    // save must not wipe the role's existing grants (it seeds itself silently).
    await page.goto('/roles');
    const row = page.getByRole('row').filter({ hasText: cappedRoleName });
    await row.getByTestId('role-edit').click();
    await expect(page).toHaveURL(/\/roles\/\d+/);

    const picker = page.locator('app-grant-picker');
    await expect(picker.getByTestId('grant-picker-categories')).toContainText(catA.name);
    await expect(picker.getByTestId('grant-picker-categories')).toContainText(catB.name);
    await expect(picker.getByTestId('grant-picker-categories')).toContainText(catC.name);

    // Save without touching the picker — the grants must survive.
    const [saveResponse] = await Promise.all([
      page.waitForResponse(
        (r) => /\/api\/role\/\d+$/.test(r.url()) && r.request().method() === 'PUT',
      ),
      page.getByRole('button', { name: 'Save Role' }).click(),
    ]);
    expect(saveResponse.status()).toBe(200);

    const saved = (await saveResponse.json()) as { categoryGrants: number[] };
    expect([...saved.categoryGrants].sort()).toEqual([catA.id, catB.id, catC.id].sort());
  });
});
