import { expect, test, type Browser } from '@playwright/test';
import { rmSync } from 'node:fs';
import { stubTokenRefresh } from './helpers/auth';
import {
  apiCreateCustomField,
  apiDeleteCustomFieldById,
  apiDeleteRoleByName,
  apiDeleteUserByName,
  apiGetCustomFieldById,
  createRole,
  createUserWithRole,
  uniqueName,
  withAdminApi,
  withApiAsCreds,
} from './helpers/provisioning';

// Editing a custom field is gated on the app-scoped `app.custom-fields.update`:
// the row's pencil (`*hasAppPermission`) and the name link's edit-vs-view mode on
// the client, and `AppPermissions` on the PUT handler on the server. Legacy User
// deliberately does NOT hold it (it can create a field but not edit one), and no
// seeded account holds update-without-delete, so an admin context provisions two
// app roles through the real UI: an "editor" holding the whole Catalog resource,
// and a "viewer" identical but with Update Custom Fields switched off.
//
// The type is immutable for everyone -- a CustomFieldValue lives in a
// type-specific column, so re-typing would mis-column every stored value.

// Generated at runtime (no static secret in the repo). The -Aa1! suffix
// guarantees upper/lower/digit/symbol in case a password policy is ever added.
const PASSWORD = `${uniqueName('pw')}-Aa1!`;

const EDITOR_AUTH_FILE = 'e2e/.auth/custom-field-editor.json';
const VIEWER_AUTH_FILE = 'e2e/.auth/custom-field-viewer.json';

/**
 * Logs in as [username] once and persists the session, waiting until the held
 * permission has landed in the "auth" localStorage slice so the saved state
 * hydrates permissions synchronously on the next load.
 */
async function captureSession(
  browser: Browser,
  username: string,
  heldPermission: string,
  path: string,
): Promise<void> {
  const context = await browser.newContext({ storageState: undefined });
  const page = await context.newPage();
  await page.goto('/auth/login');
  await page.getByLabel('Username').fill(username);
  await page.getByLabel('Password').fill(PASSWORD);
  await page.getByRole('button', { name: 'Login' }).click();
  await expect(page).toHaveURL(/\/dashboard\/group\/\d+/, { timeout: 15_000 });
  await page.waitForFunction(
    (permission) => (localStorage.getItem('auth') ?? '').includes(permission),
    heldPermission,
  );
  await context.storageState({ path });
  await context.close();
}

test.describe('Edit custom fields (app.custom-fields.update)', () => {
  test.describe.configure({ mode: 'serial' });

  let editorRole: string;
  let viewerRole: string;
  let editorUser: string;
  let viewerUser: string;
  // Two fields so the mutating test can't disturb the read-only assertions.
  let editableField: { id: number; name: string };
  let untouchedField: { id: number; name: string };
  let selectField: { id: number; name: string };

  test.beforeAll(async ({ browser }) => {
    editorRole = uniqueName('custom-field-editor-role');
    viewerRole = uniqueName('custom-field-viewer-role');
    editorUser = uniqueName('custom-field-editor');
    viewerUser = uniqueName('custom-field-viewer');

    const admin = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const adminPage = await admin.newPage();
    await stubTokenRefresh(adminPage);

    // The whole Custom Fields resource: create/read/update/delete. (The role
    // form groups permissions by RESOURCE key, so the toggle is "Custom Fields",
    // not the registry's "Catalog" category.)
    await createRole(adminPage, {
      name: editorRole,
      type: 'Application role',
      preset: 'Start from scratch',
      enableCategories: ['Account', 'Notifications', 'Custom Fields'],
    });
    // Identical, minus the one permission under test.
    await createRole(adminPage, {
      name: viewerRole,
      type: 'Application role',
      preset: 'Start from scratch',
      enableCategories: ['Account', 'Notifications', 'Custom Fields'],
      disablePermissions: [
        { panelKey: 'app.custom-fields', label: 'Update Custom Fields' },
      ],
    });

    await createUserWithRole(adminPage, {
      username: editorUser,
      password: PASSWORD,
      role: editorRole,
    });
    await createUserWithRole(adminPage, {
      username: viewerUser,
      password: PASSWORD,
      role: viewerRole,
    });
    await admin.close();

    await withAdminApi(async (api) => {
      editableField = await apiCreateCustomField(api, uniqueName('editable-field'));
      untouchedField = await apiCreateCustomField(api, uniqueName('leftover-field'));
      selectField = await apiCreateCustomField(
        api,
        uniqueName('select-field'),
        'SELECT',
        ['Option 1', 'Option 2'],
      );
    });

    await captureSession(
      browser,
      editorUser,
      'app.custom-fields.update',
      EDITOR_AUTH_FILE,
    );
    await captureSession(
      browser,
      viewerUser,
      'app.custom-fields.read',
      VIEWER_AUTH_FILE,
    );
  });

  test.afterAll(async () => {
    // Cleanup stays best-effort -- it must never mask a test failure -- but each
    // step is isolated so one failure doesn't strand the rest. A leaked role
    // blocks nothing by name (names are unique per run) but grows the role list,
    // and an assigned role can't be deleted later by hand.
    const cleanUp = async (what: string, run: () => Promise<void>) => {
      try {
        await run();
      } catch (err) {
        console.warn(`teardown: ${what} failed —`, err);
      }
    };

    try {
      await withAdminApi(async (api) => {
        for (const field of [editableField, untouchedField, selectField]) {
          if (field) {
            await cleanUp(`delete custom field ${field.name}`, () =>
              apiDeleteCustomFieldById(api, field.id),
            );
          }
        }
        // Users first -- deleting them frees the role assignments, without which
        // the roles can't be removed.
        for (const username of [editorUser, viewerUser]) {
          await cleanUp(`delete user ${username}`, () =>
            apiDeleteUserByName(api, username),
          );
        }
        for (const role of [editorRole, viewerRole]) {
          await cleanUp(`delete role ${role}`, () =>
            apiDeleteRoleByName(api, role, 'APP'),
          );
        }
      });
    } catch (err) {
      // The admin API context itself failed (e.g. login) — nothing to clean up with.
      console.warn('teardown: admin API unavailable —', err);
    }
    rmSync(EDITOR_AUTH_FILE, { force: true });
    rmSync(VIEWER_AUTH_FILE, { force: true });
  });

  // The sessions are written in beforeAll, which runs after Playwright would
  // resolve a `test.use({ storageState })` option — so each test opens its own
  // context from the file instead of declaring it as a fixture.
  test('an app.custom-fields.update holder renames a field through the dialog', async ({
    browser,
  }) => {
    const context = await browser.newContext({ storageState: EDITOR_AUTH_FILE });
    const page = await context.newPage();
    const renamed = `${editableField.name}-renamed`;
    try {
      await stubTokenRefresh(page);
      await page.goto('/custom-fields');

      const row = page
        .getByRole('row')
        .filter({ hasText: editableField.name })
        .first();
      await expect(row).toBeVisible();

      await row.getByTestId('custom-field-edit').click();

      const dialog = page.getByRole('dialog');
      await expect(dialog).toBeVisible();
      await expect(dialog.getByText(`Edit ${editableField.name}`)).toBeVisible();

      await dialog.getByLabel('Name').fill(renamed);
      await dialog.getByLabel('Description').fill('Set by the edit dialog');
      await dialog.getByTestId('dialog-submit-button').click();
      await expect(dialog).toBeHidden();

      // The table refetches after a successful save.
      await expect(
        page.getByRole('row').filter({ hasText: renamed }),
      ).toHaveCount(1);
    } finally {
      await context.close();
    }

    // The rename is what persisted, and the type is untouched.
    await withAdminApi(async (api) => {
      const field = await apiGetCustomFieldById(api, editableField.id);
      expect(field.name).toBe(renamed);
      expect(field.description).toBe('Set by the edit dialog');
      expect(field.type).toBe('TEXT');
    });
    editableField = { ...editableField, name: renamed };
  });

  test('the edit dialog locks the type and offers no delete on a saved option', async ({
    browser,
  }) => {
    const context = await browser.newContext({ storageState: EDITOR_AUTH_FILE });
    const page = await context.newPage();
    try {
      await stubTokenRefresh(page);
      await page.goto('/custom-fields');

      const row = page
        .getByRole('row')
        .filter({ hasText: selectField.name })
        .first();
      await row.getByTestId('custom-field-edit').click();

      const dialog = page.getByRole('dialog');
      await expect(dialog).toBeVisible();

      // A CustomFieldValue lives in a type-specific column, so the type can
      // never change once values may exist. In readonly mode app-select swaps
      // its mat-select for a plain non-editable input.
      await expect(dialog.getByLabel('Type')).not.toBeEditable();
      // Both options are already saved, so neither may be removed -- their ids
      // are what CustomFieldValue.SelectValue points at.
      await expect(dialog.getByTestId('custom-field-option-delete')).toHaveCount(0);
      // ...but a new option can still be appended, and that one is removable.
      await dialog.getByTestId('custom-field-option-add').click();
      await expect(dialog.getByTestId('custom-field-option-delete')).toHaveCount(1);
    } finally {
      await context.close();
    }
  });

  test('a viewer without the update permission sees no edit action', async ({
    browser,
  }) => {
    const context = await browser.newContext({ storageState: VIEWER_AUTH_FILE });
    const page = await context.newPage();
    try {
      await stubTokenRefresh(page);
      await page.goto('/custom-fields');

      const row = page
        .getByRole('row')
        .filter({ hasText: untouchedField.name })
        .first();
      // app.custom-fields.read alone still lists it...
      await expect(row).toBeVisible();
      // ...but the update gate is not held.
      await expect(row.getByTestId('custom-field-edit')).toHaveCount(0);
    } finally {
      await context.close();
    }
  });

  test('a viewer without the update permission is denied by the server', async () => {
    // The hidden button is a UI hint; the endpoint is the real gate.
    await withApiAsCreds(viewerUser, PASSWORD, async (api) => {
      const res = await api.put(`/api/customField/${untouchedField.id}`, {
        data: { name: 'Hijacked', type: 'TEXT', description: '' },
      });
      expect(res.status()).toBe(403);
    });
  });

  test('the server refuses a type change even from an update holder', async () => {
    await withApiAsCreds(editorUser, PASSWORD, async (api) => {
      const res = await api.put(`/api/customField/${untouchedField.id}`, {
        data: { name: untouchedField.name, type: 'CURRENCY', description: '' },
      });
      expect(res.status()).toBe(400);
    });

    await withAdminApi(async (api) => {
      const field = await apiGetCustomFieldById(api, untouchedField.id);
      expect(field.type).toBe('TEXT');
    });
  });
});
