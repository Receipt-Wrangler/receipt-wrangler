import { BrowserContext, expect, Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateCategory,
  apiCreateRole,
  apiDeleteCategoryById,
  apiDeleteGroupById,
  apiDeleteRoleByName,
  apiGetGroupMembers,
  apiGetUserId,
  apiMemberCatalog,
  apiSetGroupRoster,
  apiSetMemberGrants,
  apiUpdateRole,
  createGroupWithMember,
  uniqueName,
  withAdminApi,
  withApiAs,
} from './helpers/provisioning';

// The failure modes of per-member category/tag grants that are SILENT — they
// don't throw, they just quietly show someone categories they were never
// assigned. Two of these (rename resetting the restriction, and grants
// resurrecting when a removed member rejoins) were real bugs found while
// building the feature, and neither was caught by a unit test that exercised the
// helper instead of the real code path.

test.describe('Per-member grants — escalation & lifecycle', () => {
  // Serial: the tests mutate a shared role, roster and assignment.
  test.describe.configure({ mode: 'serial' });

  let adminContext: BrowserContext;
  let adminPage: Page;

  const roleName = uniqueName('grantsec-role');
  const groupName = uniqueName('grantsec-grp');

  let groupId: string;
  let roleId: number;
  let memberId: number;
  let adminId: number;
  let adminGroupRoleId: number;
  let catA: { id: number; name: string };
  let catB: { id: number; name: string };
  let catOutside: { id: number; name: string };

  // Holds group.members.update but NOT group.members.grants.update — the whole
  // point of keeping the grant permission separate.
  const ROSTER_PERMISSIONS = [
    'group.view',
    'group.receipts.read',
    'group.members.create',
    'group.members.update',
    'group.members.delete',
  ];

  /** Resets the role ceiling and the member's assignment to a known state. */
  async function reset(opts: {
    permissions?: string[];
    roleCategories?: number[];
    memberCategories?: number[];
  }): Promise<void> {
    await withAdminApi(async (api) => {
      await apiUpdateRole(api, roleId, {
        name: roleName,
        scope: 'GROUP',
        permissions: opts.permissions ?? ROSTER_PERMISSIONS,
        categoryGrants: opts.roleCategories ?? [],
      });
      const res = await apiSetMemberGrants(api, groupId, memberId, {
        categoryGrants: opts.memberCategories ?? [],
        tagGrants: [],
      });
      expect(res.status()).toBe(200);
    });
  }

  test.beforeAll(async ({ browser }) => {
    adminContext = await browser.newContext({ storageState: 'e2e/.auth/admin.json' });
    adminPage = await adminContext.newPage();
    await stubTokenRefresh(adminPage);

    await withAdminApi(async (api) => {
      memberId = await apiGetUserId(api, creds('user').username);
      adminId = await apiGetUserId(api, creds('admin').username);
      catA = await apiCreateCategory(api, uniqueName('grantsec-cat-a'));
      catB = await apiCreateCategory(api, uniqueName('grantsec-cat-b'));
      catOutside = await apiCreateCategory(api, uniqueName('grantsec-cat-out'));

      const role = await apiCreateRole(api, {
        name: roleName,
        scope: 'GROUP',
        permissions: ROSTER_PERMISSIONS,
      });
      roleId = role.id;
    });

    groupId = await createGroupWithMember(adminPage, {
      groupName,
      memberDisplayName: 'E2E User',
      roleName,
    });

    // The creating admin keeps their own (owner) group role — capture it so the
    // roster rewrites below don't accidentally strip the admin's own
    // group.members.* permissions and 403 the very calls under test.
    await withAdminApi(async (api) => {
      const members = await apiGetGroupMembers(api, groupId);
      adminGroupRoleId = members.find((m) => m.userId === adminId)!.groupRoleId;
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        await apiDeleteGroupById(api, groupId);
        await apiDeleteRoleByName(api, roleName, 'GROUP');
        await apiDeleteCategoryById(api, catA.id);
        await apiDeleteCategoryById(api, catB.id);
        await apiDeleteCategoryById(api, catOutside.id);
      });
    } catch {
      // Best-effort cleanup.
    }
    await adminContext?.close();
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  // --- Self-escalation --------------------------------------------------------

  test('a member who can manage the roster still cannot edit grants', async () => {
    // THE escalation this feature must not have: group.members.update alone must
    // not let a restricted member lift their own restriction.
    await reset({
      roleCategories: [catA.id, catB.id],
      memberCategories: [catA.id],
    });

    await withApiAs('user', async (api) => {
      const denied = await apiSetMemberGrants(api, groupId, memberId, {
        categoryGrants: [catA.id, catB.id],
      });
      expect(denied.status()).toBe(403);
    });

    // ...and the attempt changed nothing.
    const seen = await withApiAs('user', (api) => apiMemberCatalog(api, groupId));
    expect(seen.categories).toEqual([catA.name]);
  });

  test('granting group.members.grants.update lets the same member through', async () => {
    // Positive contrast, so the 403 above is proven to be about THIS permission
    // rather than some unrelated denial.
    await reset({
      permissions: [...ROSTER_PERMISSIONS, 'group.members.grants.update'],
      roleCategories: [catA.id, catB.id],
      memberCategories: [catA.id],
    });

    await withApiAs('user', async (api) => {
      const allowed = await apiSetMemberGrants(api, groupId, memberId, {
        categoryGrants: [catB.id],
      });
      expect(allowed.status()).toBe(200);
    });

    const seen = await withApiAs('user', (api) => apiMemberCatalog(api, groupId));
    expect(seen.categories).toEqual([catB.name]);
  });

  // --- Write validation -------------------------------------------------------

  test('an assignment outside the role ceiling is rejected and persists nothing', async () => {
    await reset({
      roleCategories: [catA.id, catB.id],
      memberCategories: [catA.id],
    });

    const body = await withAdminApi(async (api) => {
      const res = await apiSetMemberGrants(api, groupId, memberId, {
        categoryGrants: [catOutside.id],
      });
      expect(res.status()).toBe(400);
      return res.text();
    });

    // The error names the offending id, so an admin can see what to fix.
    expect(body).toContain(String(catOutside.id));

    const seen = await withApiAs('user', (api) => apiMemberCatalog(api, groupId));
    expect(seen.categories).toEqual([catA.name]);
  });

  test('a partially-invalid selection is rejected whole', async () => {
    // One good id + one out-of-ceiling id must not half-apply.
    await reset({
      roleCategories: [catA.id, catB.id],
      memberCategories: [catA.id],
    });

    await withAdminApi(async (api) => {
      const res = await apiSetMemberGrants(api, groupId, memberId, {
        categoryGrants: [catB.id, catOutside.id],
      });
      expect(res.status()).toBe(400);
    });

    const seen = await withApiAs('user', (api) => apiMemberCatalog(api, groupId));
    expect(seen.categories).toEqual([catA.name]);
  });

  test('a non-existent category id is rejected', async () => {
    await reset({ roleCategories: [], memberCategories: [] });

    await withAdminApi(async (api) => {
      const res = await apiSetMemberGrants(api, groupId, memberId, {
        categoryGrants: [99_999_999],
      });
      expect(res.status()).toBe(400);
    });
  });

  test('assigning grants to a non-member is rejected', async () => {
    await withAdminApi(async (api) => {
      const res = await apiSetMemberGrants(api, groupId, 99_999_999, {
        categoryGrants: [],
      });
      expect(res.status()).toBe(404);
    });
  });

  test('the URL identifies the membership, not the request body', async () => {
    // A body carrying contradicting ids must not redirect the write to another
    // membership — the URL is what the caller was authorized against.
    await reset({
      roleCategories: [catA.id, catB.id],
      memberCategories: [catA.id],
    });

    await withAdminApi(async (api) => {
      const res = await apiSetMemberGrants(api, groupId, memberId, {
        categoryGrants: [catB.id],
        // Ignored — present only to prove they are.
        groupId: 99_999_999,
        userId: 99_999_999,
      });
      expect(res.status()).toBe(200);
    });

    const seen = await withApiAs('user', (api) => apiMemberCatalog(api, groupId));
    expect(seen.categories).toEqual([catB.name]);
  });

  // --- Lifecycle: the silent ones ---------------------------------------------

  test('a group rename preserves the assignment AND its restriction', async () => {
    // Regression: UpdateGroup rewrites member rows from the request command, which
    // carry the restriction flags at their zero value. Before the fix, renaming a
    // group cleared every member's restriction and silently widened them back to
    // their role's full set — invisible in the UI.
    await reset({
      roleCategories: [catA.id, catB.id],
      memberCategories: [catA.id],
    });

    await withAdminApi(async (api) => {
      const res = await apiSetGroupRoster(api, groupId, {
        name: `${groupName}-renamed`,
        members: [
          { userId: adminId, groupRoleId: adminGroupRoleId },
          { userId: memberId, groupRoleId: roleId },
        ],
      });
      expect(res.status()).toBe(200);
    });

    // Still narrowed to their assignment — NOT widened to the role's [A, B].
    const seen = await withApiAs('user', (api) => apiMemberCatalog(api, groupId));
    expect(seen.categories).toEqual([catA.name]);

    // And the grant row itself survived the roster replace.
    const members = await withAdminApi((api) => apiGetGroupMembers(api, groupId));
    expect(members.find((m) => m.userId === memberId)?.categoryGrants).toEqual([catA.id]);
  });

  test('a removed member does not get their grants back on rejoin', async () => {
    // The highest-risk hazard: grant rows are keyed (user, group) with no FK to
    // group_members, and the teardown paths use raw deletes. Orphaned rows would
    // be re-adopted on rejoin — a removed foster parent coming back with a child
    // silently on their list again.
    await reset({
      roleCategories: [catA.id, catB.id],
      memberCategories: [catA.id],
    });

    await withAdminApi(async (api) => {
      const removed = await apiSetGroupRoster(api, groupId, {
        name: `${groupName}-renamed`,
        members: [{ userId: adminId, groupRoleId: adminGroupRoleId }],
      });
      expect(removed.status()).toBe(200);

      const readded = await apiSetGroupRoster(api, groupId, {
        name: `${groupName}-renamed`,
        members: [
          { userId: adminId, groupRoleId: adminGroupRoleId },
          { userId: memberId, groupRoleId: roleId },
        ],
      });
      expect(readded.status()).toBe(200);
    });

    const members = await withAdminApi((api) => apiGetGroupMembers(api, groupId));
    const rejoined = members.find((m) => m.userId === memberId);
    expect(rejoined).toBeDefined();
    expect(rejoined?.categoryGrants ?? []).toEqual([]);

    // They come back unrestricted (falling back to the role), not carrying the
    // revoked assignment.
    const seen = await withApiAs('user', (api) => apiMemberCatalog(api, groupId));
    expect(seen.categories).toEqual([catA.name, catB.name].sort());
  });
});
