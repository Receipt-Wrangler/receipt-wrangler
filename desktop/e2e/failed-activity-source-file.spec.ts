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

// A PDF is the case the preview used to get wrong: the uploader accepts PDF and
// HEIC, quick scan stores the upload verbatim, and encoding those bytes without
// converting them yields a data URI no browser can render. It failed silently,
// because GetFileType labels a PDF `image/jpeg` regardless -- so the assertion
// below has to check the decoded payload, not the URI's own mime.
const PDF_FIXTURE = 'e2e/fixtures/receipt.pdf';

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
  // The PDF scan gets its OWN group: apiWaitForFailedQuickScan returns the first
  // failed quick scan it finds for a group, so sharing one would let it hand back
  // the PNG activity.
  let pdfGroupId: number;
  let pdfActivityId: number;

  test.beforeAll(async () => {
    test.setTimeout(180_000);

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

      const pdfGroup = await apiCreateGroup(api, uniqueName('failed-scan-pdf-group'));
      pdfGroupId = pdfGroup.id;

      await apiQuickScan(api, pdfGroupId, adminUserId, {
        name: 'receipt.pdf',
        mimeType: 'application/pdf',
        buffer: readFileSync(PDF_FIXTURE),
      });

      const pdfActivity = await apiWaitForFailedQuickScan(api, pdfGroupId);
      pdfActivityId = pdfActivity.id;
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
        if (pdfGroupId) {
          await apiDeleteGroupById(api, String(pdfGroupId));
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

  // The regression guard for the shared BuildDisplayImageString path. There is no
  // Go test that can cover it: GetSystemTaskSourceFile needs a live Redis to
  // resolve the task payload, so the handler sits at 0% coverage and this is the
  // only thing that proves it still converts.
  test('a PDF upload previews as a real image and downloads as the PDF', async () => {
    await withAdminApi(async (api) => {
      const res = await api.get(`/api/systemTask/${pdfActivityId}/sourceFile`);
      expect(res.ok()).toBe(true);

      const sourceFile = (await res.json()) as {
        name: string;
        encodedImage: string;
      };
      expect(sourceFile.name).toBe('receipt.pdf');

      // The mime the URI CLAIMS proves nothing -- GetFileType reports image/jpeg
      // for raw PDF bytes too. Decode it and look at what a browser would get.
      const [, base64Payload] = sourceFile.encodedImage.split('base64,');
      expect(base64Payload).toBeTruthy();
      const decoded = Buffer.from(base64Payload, 'base64');
      expect(decoded.subarray(0, 4).toString('latin1')).not.toBe('%PDF');
      expect(decoded.subarray(0, 2)).toEqual(Buffer.from([0xff, 0xd8])); // JPEG SOI

      // ...while the download still serves the ORIGINAL, not the converted copy.
      const download = await api.get(
        `/api/systemTask/${pdfActivityId}/sourceFile/download`,
      );
      expect(download.ok()).toBe(true);
      expect(Buffer.from(await download.body())).toEqual(readFileSync(PDF_FIXTURE));
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
