import { expect, Locator, Page, test } from '@playwright/test';
import { stubTokenRefresh } from './helpers/auth';

// The "Generate password" button lives in the shared app-input's suffix, next to
// the visibility eye, and is opted into per field. Clicking it fills the control
// with a generated password, reveals it, copies it to the clipboard and toasts.
// These specs are deliberately non-mutating: every dialog is escaped, never
// submitted, so no user is created and nothing needs tearing down.

const GENERATED_PASSWORD_LENGTH = 20;
const COPIED_MESSAGE = 'Password generated and copied to clipboard';

function generateButton(scope: Page | Locator): Locator {
  // Unlike the app-button-wrapped test IDs, this one sits on the native button.
  return scope.getByTestId('password-generate');
}

async function openCreateUserDialog(page: Page): Promise<Locator> {
  await page.goto('/users');
  await page.getByTestId('user-add').click();
  const dialog = page.getByRole('dialog').filter({ hasText: 'Create User' });
  await expect(dialog).toBeVisible();
  return dialog;
}

test.describe('Generate password — Create User', () => {
  // Managing users is admin-only, and reading the clipboard back needs the
  // permission granted explicitly (no other spec in the suite uses it yet).
  test.use({
    storageState: 'e2e/.auth/admin.json',
    permissions: ['clipboard-read', 'clipboard-write'],
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('fills, reveals and copies a generated password', async ({ page }) => {
    const dialog = await openCreateUserDialog(page);
    const password = dialog.getByLabel('Password', { exact: true });

    await expect(password).toHaveAttribute('type', 'password');
    await generateButton(dialog).click();

    // The field is written synchronously on the click; the clipboard write and
    // its toast follow as a detached promise.
    await expect(password).toHaveAttribute('type', 'text');
    const generated = await password.inputValue();
    expect(generated).toHaveLength(GENERATED_PASSWORD_LENGTH);

    // A clipboard failure would surface the error toast instead, so the success
    // message also proves the write resolved.
    await expect(page.getByText(COPIED_MESSAGE)).toBeVisible();
    expect(await page.evaluate(() => navigator.clipboard.readText())).toEqual(
      generated,
    );

    await page.keyboard.press('Escape');
  });

  test('generates a different password on each click', async ({ page }) => {
    const dialog = await openCreateUserDialog(page);
    const password = dialog.getByLabel('Password', { exact: true });

    await generateButton(dialog).click();
    const first = await password.inputValue();

    await generateButton(dialog).click();
    await expect(password).not.toHaveValue(first);

    await page.keyboard.press('Escape');
  });

  test('re-masks a generated password with the visibility eye', async ({
    page,
  }) => {
    const dialog = await openCreateUserDialog(page);
    const password = dialog.getByLabel('Password', { exact: true });

    await generateButton(dialog).click();
    await expect(password).toHaveAttribute('type', 'text');

    await dialog.locator('button.visibility-eye-button').click();
    await expect(password).toHaveAttribute('type', 'password');
    // Masking is display-only — the generated value is still submitted.
    await expect(password).not.toHaveValue('');

    await page.keyboard.press('Escape');
  });

  test('disables the generate button for a dummy user', async ({ page }) => {
    const dialog = await openCreateUserDialog(page);

    // Ticking "Is Dummy User?" clears and disables the password control, since a
    // dummy user has no credentials — the generate button follows it.
    await dialog.getByText('Is Dummy User?').click();
    await expect(generateButton(dialog)).toBeDisabled();

    await dialog.getByText('Is Dummy User?').click();
    await expect(generateButton(dialog)).toBeEnabled();

    await page.keyboard.press('Escape');
  });
});

test.describe('Generate password — Set Password dialog', () => {
  test.use({
    storageState: 'e2e/.auth/admin.json',
    permissions: ['clipboard-read', 'clipboard-write'],
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  test('fills, reveals and copies a generated password', async ({ page }) => {
    await page.goto('/users');
    // Any row will do — the dialog is escaped, never submitted.
    await page.getByTestId('user-set-password').first().click();
    const dialog = page.getByRole('dialog').filter({ hasText: 'Set Password' });
    await expect(dialog).toBeVisible();

    const password = dialog.getByLabel('Password', { exact: true });
    await generateButton(dialog).click();

    await expect(password).toHaveAttribute('type', 'text');
    expect(await password.inputValue()).toHaveLength(GENERATED_PASSWORD_LENGTH);
    await expect(page.getByText(COPIED_MESSAGE)).toBeVisible();

    await page.keyboard.press('Escape');
  });
});

test.describe('password fields without a generate affordance', () => {
  // Logged out: the login form is an eye-only password field. Guards against the
  // shared input's showGeneratePassword default being flipped on.
  test.use({ storageState: { cookies: [], origins: [] } });

  test('the login form offers the eye but no generate button', async ({
    page,
  }) => {
    await page.goto('/auth/login');
    await expect(page.getByLabel('Password', { exact: true })).toBeVisible();

    await expect(page.locator('button.visibility-eye-button')).toHaveCount(1);
    await expect(generateButton(page)).toHaveCount(0);
  });
});
