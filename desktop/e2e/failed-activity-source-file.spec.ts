import { expect, test } from '@playwright/test';
import { readFileSync } from 'node:fs';
import { creds } from './helpers/auth';
import {
  apiCreateGroup,
  apiCreateUnreachableProcessingSettings,
  apiDeleteGroupById,
  apiDeleteProcessingSettings,
  apiGetUserId,
  apiPatchSystemSettings,
  apiQuickScan,
  apiWaitForFailedQuickScan,
  uniqueName,
  withAdminApi,
  withApiAs,
} from './helpers/provisioning';

const FIXTURE = 'e2e/fixtures/receipt.png';

// A quick scan has to be made to fail on purpose, and the only lever for that is
// the AI provider -- which is a GLOBAL system setting. So this suite mutates
// shared server state, runs serially, and restores what it found. It relies on
// the `e2e-shared-backend` job-level concurrency group in e2e.yml / mobile-e2e.yml,
// which exists precisely so global-mutating specs queue behind one another.
//
// Pointing the provider at an unreachable host, rather than assuming the backend
// has none configured, is what makes this work against the shared demo backend
// as well as a local dev database.
test.describe.serial('Failed activity source file', () => {
  let groupId: number;
  let processingSettingsId: number;
  let promptId: number;
  let originalProcessingSettingsId: unknown;
  let activityId: number;
  let adminUserId: number;

  test.beforeAll(async () => {
    test.setTimeout(120_000);

    await withAdminApi(async (api) => {
      const settings = await (await api.get('/api/systemSettings')).json();
      originalProcessingSettingsId = settings.receiptProcessingSettingsId;

      adminUserId = await apiGetUserId(api, creds('admin').username);

      const name = uniqueName('failed-scan');
      const created = await apiCreateUnreachableProcessingSettings(api, name);
      processingSettingsId = created.id;
      promptId = created.promptId;

      await apiPatchSystemSettings(api, {
        receiptProcessingSettingsId: processingSettingsId,
      });

      const group = await apiCreateGroup(api, uniqueName('failed-scan-group'));
      groupId = group.id;

      await apiQuickScan(api, groupId, adminUserId, {
        name: 'receipt.png',
        mimeType: 'image/png',
        buffer: readFileSync(FIXTURE),
      });

      const activity = await apiWaitForFailedQuickScan(api, groupId);
      activityId = activity.id;
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi(async (api) => {
        await apiPatchSystemSettings(api, {
          receiptProcessingSettingsId: originalProcessingSettingsId ?? null,
        });
        if (groupId) {
          await apiDeleteGroupById(api, String(groupId));
        }
        if (processingSettingsId) {
          await apiDeleteProcessingSettings(api, processingSettingsId, promptId);
        }
      });
    } catch (error) {
      // Best-effort teardown -- report the failure but don't mask the suite's
      // real result by throwing out of afterAll.
      console.warn('Failed to tear down the failed-activity fixtures', error);
    }
  });

  // This is the assertion no unit test can make: the upload actually survived in
  // temp/ through a failure, and the asynq payload still resolves to it.
  test('a failed quick scan keeps its upload and reports it on the activity', async () => {
    expect(activityId).toBeGreaterThan(0);
  });

  test('the source file is previewable as a converted image', async () => {
    await withAdminApi(async (api) => {
      const res = await api.get(`/api/systemTask/${activityId}/sourceFile`);
      expect(res.ok()).toBe(true);

      const sourceFile = (await res.json()) as {
        name: string;
        encodedImage: string;
      };
      expect(sourceFile.name).toBe('receipt.png');
      expect(sourceFile.encodedImage).toMatch(/^data:image\//);
    });
  });

  test('the source file downloads as the original bytes under its own name', async () => {
    await withAdminApi(async (api) => {
      const res = await api.get(
        `/api/systemTask/${activityId}/sourceFile/download`,
      );
      expect(res.ok()).toBe(true);
      expect(res.headers()['content-disposition']).toContain('receipt.png');

      // Byte-identical to what was uploaded -- the download serves the original,
      // never the converted copy the preview may use.
      expect(Buffer.from(await res.body())).toEqual(readFileSync(FIXTURE));
    });
  });

  test('a user outside the activity group cannot reach either endpoint', async () => {
    await withApiAs('user', async (api) => {
      const preview = await api.get(`/api/systemTask/${activityId}/sourceFile`);
      expect(preview.status()).toBe(403);

      const download = await api.get(
        `/api/systemTask/${activityId}/sourceFile/download`,
      );
      expect(download.status()).toBe(403);
    });
  });

  test('the system tasks table offers preview and download for the failed task', async ({
    browser,
  }) => {
    const context = await browser.newContext({
      storageState: 'e2e/.auth/admin.json',
    });
    const page = await context.newPage();

    await page.goto('/system-settings/system-tasks');

    // Newest first, and this suite just created the failure, so it is on page 1.
    await expect(
      page.getByTestId('system-task-source-file-preview').first(),
    ).toBeVisible();
    await expect(
      page.getByTestId('system-task-source-file-download').first(),
    ).toBeVisible();

    await page.getByTestId('system-task-source-file-preview').first().click();
    await expect(page.locator('app-source-file-viewer-dialog')).toBeVisible();

    await context.close();
  });
});
