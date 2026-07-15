import { expect, request as pwRequest, test } from '@playwright/test';
import { rmSync } from 'node:fs';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiDeleteRoleByName,
  apiDeleteUserByName,
  createRole,
  createUserWithRole,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import { gotoReports } from './helpers/reports';

// Saving a template is gated by the app-level app.reports.create permission. The
// happy path (admin, who holds it) lives in report-templates.spec.ts; this spec
// covers the NEGATIVE boundary: a user who can open the builder (app.reports.read)
// but lacks app.reports.create sees no Save Template control AND is refused by the
// endpoint (403). No seeded account fits, so an admin provisions a custom role +
// user through the UI, the tests run under that user's saved session, and the
// role/user are torn down through the admin API (user before role).

const NEW_USER_PASSWORD = 'a-really-secure-password';
const AUTH_FILE = 'e2e/.auth/report-template-nocreate.json';
const apiBaseUrl = (): string => process.env.E2E_BASE_URL ?? 'http://localhost:4200';

// A complete, valid ReportRequestCommand — it must pass loadReportCommand (and the
// non-empty-name check) to reach the permission gate, so the 403 proves the gate,
// not validation.
const VALID_TEMPLATE_BODY = {
  name: 'Denied Template',
  groupIds: ['1'],
  period: { preset: 'this_month' },
  detail: { mode: 'records' },
  columns: [{ kind: 'dimension', name: 'Name', label: 'Name', field: 'name' }],
  formats: ['csv'],
};

test.describe('Report Builder — Save Template permission gating', () => {
  test.use({ storageState: AUTH_FILE });

  let roleName: string;
  let username: string;

  test.beforeAll(async ({ browser }) => {
    roleName = uniqueName('report-nocreate-role');
    username = uniqueName('report-nocreate-user');

    // Admin provisions a role that grants app.reports.read (builder access) plus
    // the baselines a functional user needs, but NOT app.reports.create (enable the
    // whole Reports group, then switch off "Save Report Templates"). Only the create
    // permission is switched off — that's the gate under test — and it must be a
    // single disablePermissions entry: createRole toggles the panel open on each
    // entry, so two entries for the same panel would re-close it.
    const admin = await browser.newContext({ storageState: 'e2e/.auth/admin.json' });
    const adminPage = await admin.newPage();
    await stubTokenRefresh(adminPage);
    await createRole(adminPage, {
      name: roleName,
      type: 'Application role',
      preset: 'Start from scratch',
      enableCategories: ['Account', 'Notifications', 'Groups', 'Reports'],
      disablePermissions: [{ panelKey: 'app.reports', label: 'Save Report Templates' }],
    });
    await createUserWithRole(adminPage, { username, password: NEW_USER_PASSWORD, role: roleName });
    await admin.close();

    // Log in once as that user and persist the session (opt out of the describe's
    // AUTH_FILE, not created yet). Wait until app.reports.read is in the "auth"
    // localStorage slice so the saved session hydrates it on the next load.
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

  test('hides the Save Template control and the endpoint refuses a direct POST', async ({ page }) => {
    await gotoReports(page);

    // The generate bar renders (Generate is present), but the app.reports.create-
    // gated Save Template button is absent for this user.
    await expect(page.getByTestId('report-generate')).toBeVisible();
    await expect(page.getByTestId('report-save-template')).toHaveCount(0);

    // Server enforcement: a direct POST as this user is refused with 403 (the body
    // is valid, so this is the permission gate, not validation).
    const api = await pwRequest.newContext({ baseURL: apiBaseUrl(), storageState: AUTH_FILE });
    try {
      const res = await api.post('/api/report/template', { data: VALID_TEMPLATE_BODY });
      expect(res.status()).toBe(403);
    } finally {
      await api.dispose();
    }
  });
});
