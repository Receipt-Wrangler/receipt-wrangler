import 'dart:async';

import 'package:built_collection/built_collection.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:infinite_carousel/infinite_carousel.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/user_preferences_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/quick_scan.dart';
import 'package:receipt_wrangler_mobile/shared/classes/quick_scan_image.dart';
import 'package:receipt_wrangler_mobile/shared/functions/quick_scan_field_config.dart';
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
  List<String> categoryIds = [];
  List<String> tagIds = [];
  List<MultipartFile> files = [];

  if (images.isEmpty) {
    showErrorSnackbar(context, "Please upload at least one image");
    return;
  }

  final groupModel = Provider.of<GroupModel>(context, listen: false);

  var errored = false;
  for (var (index, image) in images.indexed) {
    final groupId = image.groupId ?? 0;
    if (groupId <= 0) {
      errored = true;
      showErrorSnackbar(
          context, "Please fix error on quick scan ${index + 1} to continue");
      break;
    }

    // Resolve each field against the group's quick-scan config, mirroring the
    // backend's resolveQuickScanFields. A hidden/optional field is sent as the
    // "unset" sentinel (0 / empty) so the server backfills the configured
    // default; a shown+required field must have a value.
    final settings = groupModel.getGroupReceiptSettings(groupId);

    final config = resolveQuickScanFieldConfig(settings);
    final showPaidBy = config.showPaidBy;
    final requirePaidBy = config.requirePaidBy;
    final showStatus = config.showStatus;
    final requireStatus = config.requireStatus;
    final showCategories = config.showCategories;
    final requireCategories = config.requireCategories;
    final showTags = config.showTags;
    final requireTags = config.requireTags;

    int paidBy = 0;
    if (showPaidBy) {
      paidBy = image.paidByUserId ?? 0;
      if (requirePaidBy && paidBy <= 0) {
        errored = true;
      }
    }

    api.ReceiptStatus status = api.ReceiptStatus.empty;
    if (showStatus) {
      status = image.status ?? api.ReceiptStatus.empty;
      if (requireStatus && status == api.ReceiptStatus.empty) {
        errored = true;
      }
    }

    List<int> catIds = [];
    if (showCategories) {
      catIds =
          image.categories.map((c) => c.id ?? 0).where((id) => id > 0).toList();
      if (requireCategories && catIds.isEmpty) {
        errored = true;
      }
    }

    List<int> tgIds = [];
    if (showTags) {
      tgIds = image.tags.map((t) => t.id ?? 0).where((id) => id > 0).toList();
      if (requireTags && tgIds.isEmpty) {
        errored = true;
      }
    }

    if (errored) {
      showErrorSnackbar(
          context, "Please fix error on quick scan ${index + 1} to continue");
      break;
    }

    files.add(image.multipartFile);
    groupIds.add(groupId);
    paidByUserIds.add(paidBy);
    statuses.add(status);
    categoryIds.add(catIds.join(","));
    tagIds.add(tgIds.join(","));
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
        statuses: statuses.toBuiltList(),
        categoryIds: categoryIds.toBuiltList(),
        tagIds: tagIds.toBuiltList());

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
