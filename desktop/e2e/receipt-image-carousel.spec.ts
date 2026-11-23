import { expect, type Locator, type Page, test } from '@playwright/test';
import { creds, stubTokenRefresh } from './helpers/auth';
import {
  apiCreateGroup,
  apiCreateReceipt,
  apiDeleteGroupById,
  apiGetUserId,
  apiUploadReceiptImage,
  uniqueName,
  withAdminApi,
} from './helpers/provisioning';
import { RECEIPT_PNG } from './helpers/quick-scan';

// The receipt form's image carousel (app-carousel on the shared
// app-slide-carousel): moving between a receipt's images with the arrows, the
// indicators and the keyboard, and the fullscreen view opening on the image
// being looked at. The unit specs run under jsdom, which has no focus model,
// layout or visibility; this runs them against a real receipt with three images.

test.use({ storageState: 'e2e/.auth/admin.json' });

test.describe('Receipt image carousel', () => {
  test.describe.configure({ mode: 'serial' });

  let groupId: number;
  let receiptId: number;

  test.beforeAll(async () => {
    await withAdminApi(async (api) => {
      const group = await apiCreateGroup(api, uniqueName('carousel'));
      groupId = group.id;
      receiptId = await apiCreateReceipt(api, {
        groupId,
        paidByUserId: await apiGetUserId(api, creds('admin').username),
        name: uniqueName('carousel-receipt'),
      });
      for (let i = 0; i < 3; i++) {
        await apiUploadReceiptImage(api, { receiptId, filePath: RECEIPT_PNG });
      }
    });
  });

  test.afterAll(async () => {
    try {
      await withAdminApi((api) => apiDeleteGroupById(api, String(groupId)));
    } catch {
      // Best-effort cleanup — don't mask a test failure with a teardown error.
    }
  });

  test.beforeEach(async ({ page }) => {
    await stubTokenRefresh(page);
  });

  /** Opens the receipt and returns the inline image stage once all three images are in. */
  async function openReceipt(page: Page): Promise<Locator> {
    await page.goto(`/receipts/${receiptId}/view`);
    const stage = page.getByTestId('receipt-image-stage');
    await expect(stage.getByTestId('receipt-image-slide')).toHaveCount(3);
    return stage;
  }

  /** Asserts that [index] is the one image shown, and the others are hidden. */
  async function expectShowing(scope: Locator, index: number): Promise<void> {
    const slides = scope.getByTestId('receipt-image-slide');
    for (let i = 0; i < 3; i++) {
      if (i === index) {
        await expect(slides.nth(i)).toBeVisible();
      } else {
        await expect(slides.nth(i)).toBeHidden();
      }
    }
  }

  test('the arrows and indicators move between the images, wrapping at both ends', async ({
    page,
  }) => {
    const stage = await openReceipt(page);
    await expectShowing(stage, 0);

    await stage.getByTestId('carousel-next').click();
    await expectShowing(stage, 1);
    await stage.getByTestId('carousel-next').click();
    await stage.getByTestId('carousel-next').click();
    await expectShowing(stage, 0);

    await stage.getByTestId('carousel-previous').click();
    await expectShowing(stage, 2);

    const indicators = stage.getByTestId('carousel-indicator');
    await indicators.nth(1).click();
    await expectShowing(stage, 1);
    await expect(indicators.nth(1)).toHaveAttribute('aria-current', 'true');
    await expect(indicators.nth(0)).not.toHaveAttribute('aria-current');
  });

  test('the keyboard moves between images from the carousel, and controls keep their own keys', async ({
    page,
  }) => {
    const stage = await openReceipt(page);

    await stage.getByTestId('carousel').focus();
    await page.keyboard.press('ArrowRight');
    await expectShowing(stage, 1);
    await page.keyboard.press('ArrowLeft');
    await expectShowing(stage, 0);
    await page.keyboard.press('Enter');
    await expectShowing(stage, 1);

    // Enter on Previous goes back. The carousel used to treat every Enter
    // inside it as "next", so this went forward.
    await stage.getByTestId('carousel-previous').focus();
    await page.keyboard.press('Enter');
    await expectShowing(stage, 0);
  });

  test('fullscreen opens on the image being shown', async ({ page }) => {
    const stage = await openReceipt(page);
    await stage.getByTestId('carousel-indicator').nth(2).click();
    await expectShowing(stage, 2);

    await page.getByTestId('receipt-image-fullscreen').locator('button').click();
    const dialog = page.getByRole('dialog');
    await expect(dialog.getByTestId('receipt-image-slide')).toHaveCount(3);
    await expectShowing(dialog, 2);

    await dialog.getByTestId('carousel-next').click();
    await expectShowing(dialog, 0);
  });
});
