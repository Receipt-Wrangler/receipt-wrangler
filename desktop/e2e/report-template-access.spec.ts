import { expect, test } from '@playwright/test';
import { rmSync } from 'node:fs';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiCreateReportTemplate,
  apiCreateRole,
  apiDeleteGroupById,
  apiDeleteReportTemplateById,
  apiDeleteRoleByName,
  apiDeleteUserByName,
  apiUpdateRole,
  createGroupWithMember,
  createUserWithRole,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import { gotoReports } from './helpers/reports';

// End-to-end coverage of the report-template access model: the per-template access
// MATRIX on the group-role form (authoring + hydration), and its enforcement — a
// group member restricted by the matrix sees only the granted template and only its
// granted row actions. The exhaustive per-permission decision matrix is unit-tested
// in api/internal/services/report_authz_test.go; these prove the UI wiring.

test.describe('Reports — access matrix authoring', () => {
  test.use({ storageState: 'e2e/.auth/admin.json' });

  let templateName: string;
  let templateId: number;
  const roleName = uniqueName('matrix-role');

  test.beforeAll(async () => {
    templateName = uniqueName('matrix-template');
    const template = await withAdminApi((api) => apiCreateReportTemplate(api, { name: templateName }));
    templateId = template.id;
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        await apiDeleteReportTemplateById(api, templateId);
        await apiDeleteRoleByName(api, roleName, 'GROUP');
      });
    } catch {
      // Best-effort teardown.
    }
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('authors the matrix on a group role and rehydrates it on re-open', async ({ page }) => {
    await page.goto('/roles');
    await page.getByRole('button', { name: 'Add Role' }).first().click();
    await expect(page).toHaveURL(/\/roles\/new$/);
    await page.getByRole('button', { name: /Group role/ }).click();
    await page.getByLabel('Role Name').fill(roleName);
    await page.getByRole('button', { name: 'Start from scratch' }).click();
    // A realistic group role — grant group.reports.read so it could actually reach reports.
    await page.getByRole('button', { name: 'Toggle all Reports permissions' }).click();

    // The matrix lists the seeded template (options load from GET /report/template/options).
    // Grant it View + Generate; leave Edit/Delete/Duplicate off.
    const viewToggle = page.getByRole('button', { name: `View ${templateName}`, exact: true });
    await expect(viewToggle).toBeVisible();
    await viewToggle.click();
    await page.getByRole('button', { name: `Generate ${templateName}`, exact: true }).click();

    // Save and capture the new role id from the create call, so we can re-open it directly.
    const [createResponse] = await Promise.all([
      page.waitForResponse((r) => r.url().endsWith('/api/role') && r.request().method() === 'POST'),
      page.getByRole('button', { name: 'Save Role' }).click(),
    ]);
    expect(createResponse.status()).toBe(200);
    const created = (await createResponse.json()) as { id: number };
    await expect(page).toHaveURL(/\/roles$/);

    // Re-open in edit mode → the matrix rehydrates: View + Generate on, the rest off.
    await page.goto(`/roles/${created.id}/edit?scope=group`);
    await expect(page.getByRole('button', { name: `View ${templateName}`, exact: true })).toHaveClass(/\bon\b/);
    await expect(page.getByRole('button', { name: `Generate ${templateName}`, exact: true })).toHaveClass(/\bon\b/);
    await expect(page.getByRole('button', { name: `Edit ${templateName}`, exact: true })).not.toHaveClass(/\bon\b/);
    await expect(page.getByRole('button', { name: `Delete ${templateName}`, exact: true })).not.toHaveClass(/\bon\b/);
    await expect(page.getByRole('button', { name: `Duplicate ${templateName}`, exact: true })).not.toHaveClass(/\bon\b/);
  });
});

test.describe('Reports — matrix restricts a group member', () => {
  const AUTH_FILE = 'e2e/.auth/report-matrix-member.json';
  test.use({ storageState: AUTH_FILE });

  const password = `${uniqueName('pw')}-Aa1!`;
  const appRoleName = uniqueName('matrix-app-role');
  const groupRoleName = uniqueName('matrix-group-role');
  const username = uniqueName('matrix-user');
  const templateAName = uniqueName('matrix-tmplA');
  const templateBName = uniqueName('matrix-tmplB');
  let groupId: string;
  let templateAId: number;
  let templateBId: number;

  test.beforeAll(async ({ browser }) => {
    // Provision through the admin browser context's own request (authenticated by
    // admin.json). withAdminApi opens a FRESH request context that inherits this
    // describe's storageState (AUTH_FILE), which doesn't exist yet during beforeAll.
    const admin = await browser.newContext({ storageState: 'e2e/.auth/admin.json' });
    const adminPage = await admin.newPage();
    await stubTokenRefresh(adminPage);

    // App role: every base report action (so the ONLY thing that can restrict edit/
    // delete/duplicate is the matrix, not a missing app permission) + the baseline a
    // functional user needs (or the 403 interceptor logs it out). No "*All" bypasses.
    // Group role: group.reports.read; its matrix is set after the templates exist.
    await apiCreateRole(admin.request, {
      name: appRoleName,
      scope: 'APP',
      permissions: [
        'app.reports.read', 'app.reports.generate', 'app.reports.update',
        'app.reports.delete', 'app.reports.duplicate',
        'app.account.read', 'app.notifications.read',
      ],
    });
    const groupRole = await apiCreateRole(admin.request, {
      name: groupRoleName, scope: 'GROUP', permissions: ['group.reports.read'],
    });

    // UI: the user (assigned the app role) and a group with the user as the group-role member.
    await createUserWithRole(adminPage, { username, password, role: appRoleName });
    groupId = await createGroupWithMember(adminPage, {
      groupName: uniqueName('matrix-group'), memberDisplayName: username, roleName: groupRoleName,
    });

    // Both templates cover the member's group; the matrix grants only template A
    // (read + generate). Template B is granted nothing → invisible to the member.
    const a = await apiCreateReportTemplate(admin.request, { name: templateAName, groupIds: [groupId] });
    const b = await apiCreateReportTemplate(admin.request, { name: templateBName, groupIds: [groupId] });
    templateAId = a.id;
    templateBId = b.id;
    await apiUpdateRole(admin.request, groupRole.id, {
      name: groupRoleName, scope: 'GROUP', permissions: ['group.reports.read'],
      reportTemplateGrants: [{ reportTemplateId: a.id, permissions: ['read', 'generate'] }],
    });
    await admin.close();

    // Log in once as the member and persist the session (wait for app.reports.read to
    // land in the "auth" localStorage slice so it hydrates on the next load).
    const userContext = await browser.newContext({ storageState: undefined });
    const userPage = await userContext.newPage();
    await userPage.goto('/auth/login');
    await userPage.getByLabel('Username').fill(username);
    await userPage.getByLabel('Password').fill(password);
    await userPage.getByRole('button', { name: 'Login' }).click();
    await expect(userPage).toHaveURL(/\/dashboard\/group\/\d+/, { timeout: 15_000 });
    await userPage.waitForFunction(() => (localStorage.getItem('auth') ?? '').includes('app.reports.read'));
    await userContext.storageState({ path: AUTH_FILE });
    await userContext.close();
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        await apiDeleteReportTemplateById(api, templateAId);
        await apiDeleteReportTemplateById(api, templateBId);
        // User + group first so the roles become unassigned and deletable.
        await apiDeleteUserByName(api, username);
        await apiDeleteGroupById(api, groupId);
        await apiDeleteRoleByName(api, appRoleName, 'APP');
        await apiDeleteRoleByName(api, groupRoleName, 'GROUP');
      });
    } catch {
      // Best-effort teardown.
    }
    rmSync(AUTH_FILE, { force: true });
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('lists only the granted template with only its granted actions', async ({ page }) => {
    await gotoReports(page);

    // Template A is visible; the matrix grants read + generate only, so Generate shows
    // while Edit/Delete/Duplicate stay hidden — even though the app role holds all of them.
    const rowA = page.getByRole('row').filter({ hasText: templateAName });
    await expect(rowA).toBeVisible();
    await expect(rowA.getByTestId('report-template-generate')).toBeVisible();
    await expect(rowA.getByTestId('report-template-edit')).toHaveCount(0);
    await expect(rowA.getByTestId('report-template-delete')).toHaveCount(0);
    await expect(rowA.getByTestId('report-template-duplicate')).toHaveCount(0);

    // Template B isn't granted at all → the matrix hides it from the list entirely.
    await expect(page.getByRole('row').filter({ hasText: templateBName })).toHaveCount(0);
  });

  test('generate runs the per-template enforcing endpoint', async ({ page }) => {
    await gotoReports(page);
    const rowA = page.getByRole('row').filter({ hasText: templateAName });

    const [response] = await Promise.all([
      page.waitForResponse(
        (r) => /\/api\/report\/template\/\d+\/generate$/.test(r.url()) && r.request().method() === 'POST',
      ),
      rowA.getByTestId('report-template-generate').click(),
    ]);
    expect(response.status()).toBe(200);
  });
});
