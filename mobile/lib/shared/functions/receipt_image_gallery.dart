import 'package:flutter/material.dart';
import 'package:gal/gal.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/utils/snackbar.dart';

/// Downloads receipt image [receiptImageId] and saves it to the device gallery,
/// driving the [loadingModel] spinner and surfacing success / access-denied /
/// failure via snackbars. Shared by the two receipt image app bars so the flow
/// stays in one place.
Future<void> saveReceiptImageToGallery(
    BuildContext context, LoadingModel loadingModel, int receiptImageId) async {
  loadingModel.setIsLoading(true);
  try {
    final imageResponse = await OpenApiClient.client
        .getReceiptImageApi()
        .downloadReceiptImageById(receiptImageId: receiptImageId);

    final imageBytes = imageResponse.data;
    if (imageBytes == null) {
      return;
    }

    // Gal.requestAccess() returns false when the user denies gallery access;
    // bail out instead of falling through to a putImageBytes that would throw.
    if (!await Gal.requestAccess()) {
      if (context.mounted) {
        showErrorSnackbar(context, "Gallery access denied");
      }
      return;
    }

    await Gal.putImageBytes(imageBytes);
    if (context.mounted) {
      showSuccessSnackbar(context, "Image saved to gallery",
          action: SnackBarAction(
              label: "Open", onPressed: () async => await Gal.open()));
    }
  } catch (e) {
    if (context.mounted) {
      showErrorSnackbar(context, "Failed to save image to gallery");
    }
  } finally {
    loadingModel.setIsLoading(false);
  }
}
