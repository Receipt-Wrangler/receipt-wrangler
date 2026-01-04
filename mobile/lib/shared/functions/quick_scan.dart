import 'dart:async';

import 'package:built_collection/built_collection.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:infinite_carousel/infinite_carousel.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/user_preferences_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/quick_scan.dart';
import 'package:receipt_wrangler_mobile/shared/classes/quick_scan_image.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/delete_button.dart';
import 'package:receipt_wrangler_mobile/utils/has_feature.dart';
import 'package:receipt_wrangler_mobile/utils/scan.dart';
import 'package:receipt_wrangler_mobile/utils/snackbar.dart';
import 'package:rxdart/rxdart.dart';

import '../../client/client.dart';
import '../../utils/bottom_sheet.dart';

Widget _getUploadIcon(
    context,
    BehaviorSubject<List<QuickScanImage>> imageSubject,
    BehaviorSubject<bool> isCompletedSubject) {
  return StreamBuilder<bool>(
    stream: isCompletedSubject.stream,
    builder: (context, snapshot) {
      final isCompleted = snapshot.hasData && snapshot.data == true;

      return IconButton(
        icon: const Icon(Icons.add_a_photo),
        onPressed: isCompleted
            ? null
            : () async {
                var uploadedImages = await scanImagesMultiPart(100);
                if (uploadedImages.isNotEmpty) {
                  List<QuickScanImage> quickScanImages = [];
                  var initialQuickScanValues = _getInitialQuickScanValues(context);
                  for (var image in uploadedImages) {
                    var quickScanImage = QuickScanImage.fromUploadMultipartFileData(
                        image,
                        initialQuickScanValues.groupId,
                        initialQuickScanValues.paidByUserId,
                        initialQuickScanValues.status);
                    quickScanImages.add(quickScanImage);
                  }
                  imageSubject.add(imageSubject.value + quickScanImages);
                }
              },
      );
    },
  );
}

Widget _getGalleryUploadImage(
    context,
    BehaviorSubject<List<QuickScanImage>> imageSubject,
    BehaviorSubject<bool> isCompletedSubject) {
  return StreamBuilder<bool>(
    stream: isCompletedSubject.stream,
    builder: (context, snapshot) {
      final isCompleted = snapshot.hasData && snapshot.data == true;

      return IconButton(
        icon: const Icon(Icons.upload_file_rounded),
        onPressed: isCompleted
            ? null
            : () async {
                var uploadedImages = await getGalleryImages();
                if (uploadedImages.isNotEmpty) {
                  List<QuickScanImage> quickScanImages = [];
                  var initialQuickScanValues = _getInitialQuickScanValues(context);
                  for (var image in uploadedImages) {
                    var quickScanImage = QuickScanImage.fromUploadMultipartFileData(
                        image,
                        initialQuickScanValues.groupId,
                        initialQuickScanValues.paidByUserId,
                        initialQuickScanValues.status);
                    quickScanImages.add(quickScanImage);
                  }

                  imageSubject.add(imageSubject.value + quickScanImages);
                }
              },
      );
    },
  );
}

({int? groupId, int? paidByUserId, api.ReceiptStatus? status})
    _getInitialQuickScanValues(BuildContext context) {
  late final userPreferenceModel =
      Provider.of<UserPreferencesModel>(context, listen: false);
  final userPreferences = userPreferenceModel.userPreferences;
  return (
    groupId: userPreferences.quickScanDefaultGroupId,
    paidByUserId: userPreferences.quickScanDefaultPaidById,
    status: userPreferences.quickScanDefaultStatus,
  );
}

Widget _getSubmitButton(
    BuildContext context,
    BehaviorSubject<List<QuickScanImage>> imageSubject,
    BehaviorSubject<bool> isCompletedSubject) {
  return StreamBuilder<bool>(
    stream: isCompletedSubject.stream,
    builder: (context, snapshot) {
      if (snapshot.hasData && snapshot.data == true) {
        return const SizedBox.shrink();
      }

      return BottomSubmitButton(
        onPressed: () async {
          final loadingModel = Provider.of<LoadingModel>(context, listen: false);
          await _submitQuickScan(context, imageSubject.value, loadingModel, isCompletedSubject);
        },
      );
    },
  );
}

Future<void> _submitQuickScan(
    BuildContext context,
    List<QuickScanImage> images,
    LoadingModel loadingModel,
    BehaviorSubject<bool> isCompletedSubject) async {
  List<int> groupIds = [];
  List<int> paidByUserIds = [];
  List<api.ReceiptStatus> statuses = [];
  List<MultipartFile> files = [];

  if (images.isEmpty) {
    showErrorSnackbar(context, "Please upload at least one image");
    return;
  }

  var errored = false;
  for (var (index, image) in images.indexed) {
    files.add(image.multipartFile);
    var isGroupIdValid = image.groupId != null && (image.groupId ?? 0) > 0;
    var isPaidByUserIdValid =
        image.paidByUserId != null && (image.paidByUserId ?? 0) > 0;
    var isStatusValid =
        image.status != null && image.status != api.ReceiptStatus.empty;

    if (isGroupIdValid && isPaidByUserIdValid && isStatusValid) {
      groupIds.add(image.groupId as int);
      paidByUserIds.add(image.paidByUserId as int);
      statuses.add(image.status as api.ReceiptStatus);
    } else {
      errored = true;
      showErrorSnackbar(
          context,
          "Please fix error on quick scan " +
              (index + 1).toString() +
              " to continue");
      break;
    }
  }

  if (errored) {
    return;
  }

  loadingModel.setIsLoading(true);

  try {
    await OpenApiClient.client.getReceiptApi().quickScanReceipt(
        files: files.toBuiltList(),
        groupIds: groupIds.toBuiltList(),
        paidByUserIds: paidByUserIds.toBuiltList(),
        statuses: statuses.toBuiltList());

    var imageWord = images.length > 1 ? "images" : "image";

    showSuccessSnackbar(
      context,
      "Successfully queued $imageWord for processing!",
    );

    isCompletedSubject.add(true);
  } catch (e) {
    print(e);
    showApiErrorSnackbar(context, e as dynamic);
  } finally {
    loadingModel.setIsLoading(false);
  }

  return;
}

Widget _getDeleteIcon(
    InfiniteScrollController infiniteScrollController,
    BehaviorSubject<List<QuickScanImage>> imageSubject,
    BehaviorSubject<bool> isCompletedSubject) {
  return StreamBuilder<List<QuickScanImage>>(
    stream: imageSubject.stream.asBroadcastStream(),
    builder: (context, imageSnapshot) {
      if (imageSnapshot.hasData && imageSnapshot.data!.isNotEmpty) {
        return StreamBuilder<bool>(
          stream: isCompletedSubject.stream,
          builder: (context, completedSnapshot) {
            final isCompleted = completedSnapshot.hasData && completedSnapshot.data == true;

            return DeleteButton(
              onPressed: isCompleted
                  ? null
                  : () {
                      var images = imageSubject.value;
                      images.removeAt(infiniteScrollController.selectedItem);
                      imageSubject.add(images);
                    },
            );
          },
        );
      } else {
        return const SizedBox();
      }
    },
  );
}

showQuickScanBottomSheet(context) {
  if (!hasAiPoweredReceipts(context)) {
    showErrorSnackbar(context,
        "A configured Receipt Processing Settings is required to use Quick Scan. Contact your administrator for more information.");
    return;
  }

  var infiniteScrollController = InfiniteScrollController();
  BehaviorSubject<List<QuickScanImage>> imageSubject =
      BehaviorSubject<List<QuickScanImage>>.seeded([]);
  BehaviorSubject<bool> isCompletedSubject =
      BehaviorSubject<bool>.seeded(false);

  List<Widget> actions = [
    _getUploadIcon(context, imageSubject, isCompletedSubject),
    _getGalleryUploadImage(context, imageSubject, isCompletedSubject),
    _getDeleteIcon(infiniteScrollController, imageSubject, isCompletedSubject),
  ];

  showFullscreenBottomSheet(
      context,
      QuickScan(
        imageSubject: imageSubject,
        infiniteScrollController: infiniteScrollController,
        isCompletedSubject: isCompletedSubject,
      ),
      "Quick Scan",
      actions: actions,
      bodyPadding: EdgeInsets.zero,
      bottomSheetWidget: _getSubmitButton(context, imageSubject, isCompletedSubject));
}
