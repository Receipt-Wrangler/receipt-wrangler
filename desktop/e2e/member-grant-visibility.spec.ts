import { type APIRequestContext, BrowserContext, expect, Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateCategory,
  apiCreateRole,
  apiCreateTag,
  apiDeleteCategoryById,
  apiDeleteGroupById,
  apiDeleteRoleByName,
  apiDeleteTagById,
  apiGetUserId,
  apiMemberCatalog,
  apiSetMemberGrants,
  apiUpdateRole,
  createGroupWithMember,
  uniqueName,
  withAdminApi,
  withApiAs,
} from './helpers/provisioning';

// Effective category/tag visibility for a group member, which composes TWO layers
// by INTERSECTION: the group role's grants (the ceiling) and the member's own
// individual assignment (which narrows within it).
//
// The per-layer rules are unit-tested in api/internal/services/group_member_grants_test.go;
// what these prove is the COMPOSED behavior end to end, through the real delivery
// path — appData's per-group catalogs, which is exactly the array the desktop's
// receipt-form pickers render from.
//
// The member is the standard e2e-user: the default app role (Legacy User)
// deliberately omits app.categories.read / app.tags.read, so they do NOT get the
// admin grant bypass and are genuinely restrictable.

test.describe('Per-member category/tag visibility', () => {
  // Serial: each test reconfigures the shared role + assignment before asserting.
  test.describe.configure({ mode: 'serial' });

  let adminContext: BrowserContext;
  let adminPage: Page;

  const roleName = uniqueName('grantvis-role');
  const groupName = uniqueName('grantvis-grp');

  let groupId: string;
  let roleId: number;
  let memberId: number;
  const category: Record<'a' | 'b' | 'c' | 'd', { id: number; name: string }> = {} as never;
  let tagOne: { id: number; name: string };
  let tagTwo: { id: number; name: string };

  const ROLE_PERMISSIONS = [
    'group.view',
    'group.receipts.read',
    'group.receipts.create',
    'group.receipts.update',
  ];

  /**
   * Sets the role ceiling and the member's individual assignment, then returns the
   * member's resulting catalog. Reconfiguring per test keeps each case
   * self-describing instead of depending on the previous test's leftovers.
   */
  async function configureAndRead(opts: {
    roleCategories?: number[];
    roleTags?: number[];
    requiresIndividualCategoryGrants?: boolean;
    requiresIndividualTagGrants?: boolean;
    memberCategories?: number[];
    memberTags?: number[];
  }): Promise<{ categories: string[]; tags: string[] }> {
    await withAdminApi(async (api) => {
      await apiUpdateRole(api, roleId, {
        name: roleName,
        scope: 'GROUP',
        permissions: ROLE_PERMISSIONS,
        categoryGrants: opts.roleCategories ?? [],
        tagGrants: opts.roleTags ?? [],
        requiresIndividualCategoryGrants: opts.requiresIndividualCategoryGrants ?? false,
        requiresIndividualTagGrants: opts.requiresIndividualTagGrants ?? false,
      });
      const res = await apiSetMemberGrants(api, groupId, memberId, {
        categoryGrants: opts.memberCategories ?? [],
        tagGrants: opts.memberTags ?? [],
      });
      expect(res.status()).toBe(200);
    });

    return withApiAs('user', (api: APIRequestContext) => apiMemberCatalog(api, groupId));
  }

  test.beforeAll(async ({ browser }) => {
    adminContext = await browser.newContext({ storageState: 'e2e/.auth/admin.json' });
    adminPage = await adminContext.newPage();
    await stubTokenRefresh(adminPage);

    await withAdminApi(async (api) => {
      memberId = await apiGetUserId(api, creds('user').username);
      for (const key of ['a', 'b', 'c', 'd'] as const) {
        category[key] = await apiCreateCategory(api, uniqueName(`grantvis-cat-${key}`));
      }
      tagOne = await apiCreateTag(api, uniqueName('grantvis-tag-1'));
      tagTwo = await apiCreateTag(api, uniqueName('grantvis-tag-2'));

      const role = await apiCreateRole(api, {
        name: roleName,
        scope: 'GROUP',
        permissions: ROLE_PERMISSIONS,
      });
      roleId = role.id;
    });

    groupId = await createGroupWithMember(adminPage, {
      groupName,
      memberDisplayName: 'E2E User',
      roleName,
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        // Group first (frees the role assignment), then the role, then the global
        // categories/tags — a leaked unique name breaks the next run's create.
        await apiDeleteGroupById(api, groupId);
        await apiDeleteRoleByName(api, roleName, 'GROUP');
        for (const key of ['a', 'b', 'c', 'd'] as const) {
          await apiDeleteCategoryById(api, category[key].id);
        }
        await apiDeleteTagById(api, tagOne.id);
        await apiDeleteTagById(api, tagTwo.id);
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a cleanup error.
    }
    await adminContext?.close();
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('role grants only: the member sees the role\'s full set', async () => {
    const seen = await configureAndRead({
      roleCategories: [category.a.id, category.b.id, category.c.id],
    });

    expect(seen.categories).toEqual(
      [category.a.name, category.b.name, category.c.name].sort(),
    );
  });

  test('member grants only: an unrestricted role imposes no ceiling', async () => {
    const seen = await configureAndRead({ memberCategories: [category.a.id] });

    expect(seen.categories).toEqual([category.a.name]);
  });

  test('both layers intersect — the member sees only their assignment', async () => {
    // The headline case: role allows A/B/C, the individual is assigned only B.
    const seen = await configureAndRead({
      roleCategories: [category.a.id, category.b.id, category.c.id],
      memberCategories: [category.b.id],
    });

    expect(seen.categories).toEqual([category.b.name]);
  });

  test('narrowing the role below an existing assignment fails closed', async () => {
    // Assign B within the ceiling, then narrow the role so B falls outside it.
    // Intersection must resolve to NOTHING rather than falling back to either
    // layer's set — the one case the write-side ceiling check cannot prevent,
    // because narrowing a role is a legitimate admin action.
    await configureAndRead({
      roleCategories: [category.a.id, category.b.id, category.c.id],
      memberCategories: [category.b.id],
    });

    await withAdminApi((api) =>
      apiUpdateRole(api, roleId, {
        name: roleName,
        scope: 'GROUP',
        permissions: ROLE_PERMISSIONS,
        categoryGrants: [category.a.id],
      }),
    );

    const seen = await withApiAs('user', (api) => apiMemberCatalog(api, groupId));
    expect(seen.categories).toEqual([]);
  });

  test('clearing the assignment hands the member back to their role', async () => {
    await configureAndRead({
      roleCategories: [category.a.id, category.b.id, category.c.id],
      memberCategories: [category.b.id],
    });

    const seen = await configureAndRead({
      roleCategories: [category.a.id, category.b.id, category.c.id],
      memberCategories: [],
    });

    expect(seen.categories).toEqual(
      [category.a.name, category.b.name, category.c.name].sort(),
    );
  });

  test('categories and tags are restricted independently', async () => {
    const seen = await configureAndRead({
      roleCategories: [category.a.id, category.b.id],
      memberCategories: [category.a.id],
      // Neither layer restricts tags.
    });

    expect(seen.categories).toEqual([category.a.name]);
    expect(seen.tags).toEqual([tagOne.name, tagTwo.name].sort());
  });

  test('restricting tags leaves categories untouched', async () => {
    const seen = await configureAndRead({
      roleTags: [tagOne.id, tagTwo.id],
      memberTags: [tagTwo.id],
    });

    expect(seen.tags).toEqual([tagTwo.name]);
    // Categories were never configured on either layer → unrestricted.
    expect(seen.categories).toEqual(
      expect.arrayContaining([category.a.name, category.d.name]),
    );
  });

  test('require-individual: an unassigned member sees nothing', async () => {
    // Without the toggle this member would inherit the role's full A/B/C set.
    const seen = await configureAndRead({
      roleCategories: [category.a.id, category.b.id, category.c.id],
      requiresIndividualCategoryGrants: true,
      memberCategories: [],
    });

    expect(seen.categories).toEqual([]);
  });

  test('require-individual: an assigned member sees exactly their assignment', async () => {
    const seen = await configureAndRead({
      roleCategories: [category.a.id, category.b.id, category.c.id],
      requiresIndividualCategoryGrants: true,
      memberCategories: [category.c.id],
    });

    expect(seen.categories).toEqual([category.c.name]);
  });

  test('require-individual on tags does not fail categories closed', async () => {
    const seen = await configureAndRead({
      roleCategories: [category.a.id, category.b.id],
      requiresIndividualTagGrants: true,
      memberCategories: [],
      memberTags: [],
    });

    expect(seen.categories).toEqual([category.a.name, category.b.name].sort());
    expect(seen.tags).toEqual([]);
  });

  test('creating a receipt with an out-of-grant category is rejected', async () => {
    // Read-hiding without a write block would be a hole: a crafted request could
    // still attach a category the member can't see.
    await configureAndRead({
      roleCategories: [category.a.id, category.b.id, category.c.id],
      memberCategories: [category.b.id],
    });

    await withApiAs('user', async (api) => {
      const denied = await api.post('/api/receipt', {
        data: {
          name: uniqueName('grantvis-receipt'),
          amount: '10.00',
          date: '2024-01-01T00:00:00Z',
          groupId: Number(groupId),
          paidByUserId: memberId,
          status: 'OPEN',
          categories: [{ id: category.a.id, name: category.a.name }],
        },
      });
      expect(denied.status()).toBe(403);

      const allowed = await api.post('/api/receipt', {
        data: {
          name: uniqueName('grantvis-receipt-ok'),
          amount: '10.00',
          date: '2024-01-01T00:00:00Z',
          groupId: Number(groupId),
          paidByUserId: memberId,
          status: 'OPEN',
          categories: [{ id: category.b.id, name: category.b.name }],
        },
      });
      expect(allowed.status()).toBe(200);
    });
  });

  test('the receipt form offers exactly the effective set', async ({ page }) => {
    // Proves the UI renders from the same per-group catalog the assertions above
    // read, rather than a parallel source. Driven from a receipt's EDIT form: the
    // group comes from the receipt itself, so there is no dependence on which
    // group happens to be selected in app state.
    await configureAndRead({
      roleCategories: [category.a.id, category.b.id, category.c.id],
      memberCategories: [category.b.id],
    });

    const receiptId = await withApiAs('user', async (api) => {
      const res = await api.post('/api/receipt', {
        data: {
          name: uniqueName('grantvis-edit'),
          amount: '10.00',
          date: '2024-01-01T00:00:00Z',
          groupId: Number(groupId),
          paidByUserId: memberId,
          status: 'OPEN',
        },
      });
      expect(res.status()).toBe(200);
      return ((await res.json()) as { id: number }).id;
    });

    await page.goto(`/receipts/${receiptId}/edit`);

    const categoryField = page
      .getByTestId('receipt-categories')
      .getByRole('combobox', { name: 'Categories' });
    await expect(categoryField).toBeVisible();
    await categoryField.click();

    // Only the intersected category is offered; the other two the ROLE allows are
    // not, because this member's individual assignment narrows past them.
    await expect(page.getByRole('option', { name: category.b.name, exact: true })).toBeVisible();
    await expect(page.getByRole('option', { name: category.a.name, exact: true })).toHaveCount(0);
    await expect(page.getByRole('option', { name: category.c.name, exact: true })).toHaveCount(0);
  });
});
