import { test as setup } from '@playwright/test';
import { mkdirSync } from 'node:fs';
import { dirname } from 'node:path';
import { loginViaUi } from './helpers/auth';

// One state file per role. Tests default to the user state via the chromium
// project config; admin-scoped tests override with:
//   test.use({ storageState: 'e2e/.auth/admin.json' });
export const USER_AUTH_FILE = 'e2e/.auth/user.json';
export const ADMIN_AUTH_FILE = 'e2e/.auth/admin.json';

setup('authenticate as regular user', async ({ page }) => {
  mkdirSync(dirname(USER_AUTH_FILE), { recursive: true });
  await loginViaUi(page, 'user');
  await page.context().storageState({ path: USER_AUTH_FILE });
});

// Admin-scoped tests (e.g. roles.spec.ts) override the default user state with:
//   test.use({ storageState: 'e2e/.auth/admin.json' });
setup('authenticate as admin', async ({ page }) => {
  mkdirSync(dirname(ADMIN_AUTH_FILE), { recursive: true });
  await loginViaUi(page, 'admin');
  await page.context().storageState({ path: ADMIN_AUTH_FILE });
});
