import { expect, request as pwRequest, test } from '@playwright/test';
import { rmSync } from 'node:fs';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiCreateReportTemplate,
  apiDeleteReportTemplateById,
  apiDeleteRoleByName,
  apiDeleteUserByName,
  createRole,
  createUserWithRole,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import { gotoReports } from './helpers/reports';

// The generate/duplicate row actions are gated by app.reports.generate /
// app.reports.duplicate. This spec covers the NEGATIVE boundary: a user who can see
// the list (app.reports.read) but holds neither of those two sees no Generate/
// Duplicate control AND is refused by both endpoints (403). No seeded account fits,
// so an admin provisions a custom role (Reports category on, those two switched off)
// + user through the UI; the tests run under that user's saved session; the role/
// user/seeded template are torn down through the admin API.

const NEW_USER_PASSWORD = `${uniqueName('pw')}-Aa1!`;
const AUTH_FILE = 'e2e/.auth/report-action-nocreate.json';
const apiBaseUrl = (): string => process.env.E2E_BASE_URL ?? 'http://localhost:4200';

// A complete, valid ReportRequestCommand — it must pass validation to reach the
// permission gate, so a 403 proves the gate, not a bad body.
const VALID_GENERATE_BODY = {
  name: 'Denied Generate',
  groupIds: ['1'],
  period: { preset: 'this_month' },
  detail: { mode: 'records' },
  columns: [{ kind: 'dimension', name: 'Name', label: 'Name', field: 'name' }],
  formats: ['csv'],
};

test.describe('Reports — generate/duplicate permission gating', () => {
  test.use({ storageState: AUTH_FILE });

  let roleName: string;
  let username: string;
  let seededId: number;
  let seededName: string;

  test.beforeAll(async ({ browser }) => {
    roleName = uniqueName('report-noactions-role');
    username = uniqueName('report-noactions-user');
    seededName = uniqueName('report-template-gated');

    // Admin provisions a role that grants app.reports.read (list access) plus the
    // baseline reads a functional user needs, but NOT app.reports.generate/duplicate:
    // enable the whole Reports category, then switch those two permissions off. Both
    // sit in the same app.reports panel — createRole opens that panel once and toggles
    // both.
    const admin = await browser.newContext({ storageState: 'e2e/.auth/admin.json' });
    const adminPage = await admin.newPage();
    await stubTokenRefresh(adminPage);
    await createRole(adminPage, {
      name: roleName,
      type: 'Application role',
      preset: 'Start from scratch',
      enableCategories: ['Account', 'Notifications', 'Groups', 'Reports'],
      disablePermissions: [
        { panelKey: 'app.reports', label: 'Duplicate Report Templates' },
        { panelKey: 'app.reports', label: 'Generate Reports' },
      ],
    });
    await createUserWithRole(adminPage, { username, password: NEW_USER_PASSWORD, role: roleName });

    // Seed a template (admin-owned) so the list has a row whose actions we inspect.
    // Use the admin context's own request (authenticated by admin.json) rather than
    // withAdminApi: the describe's storageState (AUTH_FILE) doesn't exist yet, and a
    // fresh request context would try to load it.
    const template = await apiCreateReportTemplate(admin.request, { name: seededName });
    seededId = template.id;
    await admin.close();

    // Log in once as the provisioned user and persist the session (opt out of the
    // describe's AUTH_FILE, not created yet). Wait until app.reports.read is in the
    // "auth" localStorage slice so the saved session hydrates it on the next load.
    const userContext = await browser.newContext({ storageState: undefined });
    const userPage = await userContext.newPage();
    await userPage.goto('/auth/login');
    await userPage.getByLabel('Username').fill(username);
    await userPage.getByLabel('Password').fill(NEW_USER_PASSWORD);
    await userPage.getByRole('button', { name: 'Login' }).click();
    await expect(userPage).toHaveURL(/\/dashboard\/group\/\d+/, { timeout: 15_000 });
    await userPage.waitForFunction(() =>
      (localStorage.getItem('auth') ?? '').includes('app.reports.read'),
    );
    await userContext.storageState({ path: AUTH_FILE });
    await userContext.close();
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        await apiDeleteReportTemplateById(api, seededId);
        // User first so the role becomes unassigned and deletable.
        await apiDeleteUserByName(api, username);
        await apiDeleteRoleByName(api, roleName, 'APP');
      });
    } catch {
      // Best-effort cleanup — don't mask a test failure with a teardown error.
    }
    rmSync(AUTH_FILE, { force: true });
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('hides the generate/duplicate row actions and the endpoints refuse a direct call', async ({ page }) => {
    await gotoReports(page);

    const row = page.getByRole('row').filter({ hasText: seededName });
    await expect(row).toBeVisible();

    // The generate/duplicate controls are absent for this user; edit/delete (read/
    // delete, which the role keeps) remain.
    await expect(row.getByTestId('report-template-generate')).toHaveCount(0);
    await expect(row.getByTestId('report-template-duplicate')).toHaveCount(0);
    await expect(row.getByTestId('report-template-edit')).toBeVisible();
    await expect(row.getByTestId('report-template-delete')).toBeVisible();

    // Server enforcement: direct calls as this user are refused with 403. Duplicate is
    // app-only gated, so its 403 is unambiguously app.reports.duplicate.
    const api = await pwRequest.newContext({ baseURL: apiBaseUrl(), storageState: AUTH_FILE });
    try {
      const duplicate = await api.post(`/api/report/template/${seededId}/duplicate`);
      expect(duplicate.status()).toBe(403);

      const generate = await api.post('/api/report/generate', { data: VALID_GENERATE_BODY });
      expect(generate.status()).toBe(403);
    } finally {
      await api.dispose();
    }
  });
});
