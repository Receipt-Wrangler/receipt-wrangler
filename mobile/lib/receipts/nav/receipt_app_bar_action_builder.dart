import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:gal/gal.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/full_screen_image_viewer.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';
import 'package:rxdart/rxdart.dart';

import '../../client/client.dart';
import '../../enums/upload_method.dart';
import '../../interfaces/upload_multipart_file_data.dart';
import '../../models/loading_model.dart';
import '../../models/receipt_model.dart';
import '../../shared/widgets/receipt_edit_popup_menu.dart';
import '../../utils/forms.dart';
import '../../utils/scan.dart';
import '../../utils/snackbar.dart';

class ReceiptAppBarActionBuilder {
  late final ReceiptModel receiptModel;

  late final BuildContext context;

  late final loadingModel = Provider.of<LoadingModel>(context, listen: false);

  late final formState = getFormStateFromContext(context);

  ReceiptAppBarActionBuilder(BuildContext context, ReceiptModel receiptModel) {
    this.context = context;
    this.receiptModel = receiptModel;
  }

  bool areImagesToUpload() {
    return !receiptModel.imagesToUploadBehaviorSubject.value.isEmpty;
  }

  api.FileDataView? getCurrentlySelectedImage() {
    var index = receiptModel.infiniteScrollController.selectedItem;
    return receiptModel.imageBehaviorSubject.value[index];
  }

  List<Widget> buildAppBarMenu(GoRouterState state) {
    if (state.fullPath!.contains("images")) {
      return [_buildImagesAppBarMenu(state.fullPath!)];
    } else if (state.fullPath!.contains("comments")) {
      return [_buildCommentsAppBarMenu(state.fullPath!)];
    } else {
      return [_buildReceiptAppBarMenu(state.fullPath!)];
    }
  }

  Widget _buildReceiptAppBarMenu(String fullPath) {
    if (fullPath.contains("view")) {
      return ReceiptEditPopupMenu(
          groupId: receiptModel.receipt.groupId,
          popupMenuChildren: [
            PopupMenuItem(
              child: Text("Edit"),
              onTap: () {
                context.go("/receipts/${receiptModel.receipt.id}/edit");
              },
            )
          ]);
    }

    return SizedBox.shrink();
  }

  Widget _buildCommentsAppBarMenu(String fullPath) {
    if (fullPath.contains("view")) {
      return ReceiptEditPopupMenu(
          groupId: receiptModel.receipt.groupId,
          popupMenuChildren: [
            PopupMenuItem(
              child: Text("Edit"),
              onTap: () {
                context
                    .go("/receipts/${receiptModel.receipt.id}/comments/edit");
              },
            )
          ]);
    }

    return SizedBox.shrink();
  }

  Widget _buildImagesAppBarMenu(String fullPath) {
    List<PopupMenuEntry> options = [];
    if (isEditingBasedOnFullPath(fullPath)) {
      options = buildEditAppBarMenuOptions();
    } else {
      options = _buildViewAppBarMenuOptions();
    }

    options.add(buildImageDownloadButton());
    options.add(viewInFullScreenButton());

    var combinedStream = Rx.merge([
      receiptModel.imagesToUploadBehaviorSubject.stream,
      receiptModel.imageBehaviorSubject.stream
    ]).asBroadcastStream();

    return StreamBuilder(
        stream: combinedStream,
        builder: (context, snapshot) {
          return ReceiptEditPopupMenu(
              groupId: receiptModel.receipt.groupId,
              popupMenuChildren: options);
        });
  }

  PopupMenuItem viewInFullScreenButton() {
    return PopupMenuItem(
        child: Text("View in full screen"),
        enabled:
            formState == WranglerFormState.add ? areImagesToUpload() : true,
        onTap: () {
          var selectedIndex =
              receiptModel.infiniteScrollController.selectedItem;
          var image = receiptModel
                  .imageBehaviorSubject.value[selectedIndex]?.encodedImage ??
              "";

          var bytes = getBytesFromEncodedImage(image);
          Navigator.push(
            context,
            MaterialPageRoute(
                builder: (context) =>
                    FullScreenImageViewer(image: Image.memory(bytes))),
          );
        });
  }

  PopupMenuItem buildImageDownloadButton() {
    return PopupMenuItem(
      child: const Text("Download"),
      enabled: formState == WranglerFormState.add ? areImagesToUpload() : true,
      onTap: () async => await downloadImage(),
    );
  }

  Future downloadImage() async {
    loadingModel.setIsLoading(true);
    var receiptImage = getCurrentlySelectedImage();
    if (receiptImage == null) {
      return;
    }

    var imageResponse = await OpenApiClient.client
        .getReceiptImageApi()
        .downloadReceiptImageById(receiptImageId: receiptImage.id);

    var imageBytes = imageResponse.data;

    if (imageBytes == null) {
      loadingModel.setIsLoading(false);
      return;
    }

    await Gal.putImageBytes(imageBytes).then((value) {
      var snackbarAction = SnackBarAction(
        label: "Open",
        onPressed: () async => await Gal.open(),
      );

      showSuccessSnackbar(context, "Image saved to gallery",
          action: snackbarAction);
      loadingModel.setIsLoading(false);
    }).catchError((e) {
      showErrorSnackbar(context, "Failed to save image to gallery");
      loadingModel.setIsLoading(false);
    });
  }

  List<PopupMenuEntry> _buildViewAppBarMenuOptions() {
    return [
      PopupMenuItem(
          value: "edit",
          child: const Text("Edit"),
          onTap: () =>
              context.go("/receipts/${receiptModel.receipt.id}/images/edit")),
    ];
  }

  List<PopupMenuEntry> buildEditAppBarMenuOptions() {
    List<PopupMenuEntry> popupMenuEntries = [
      buildUploadFromCameraButton(),
      buildUploadFromGalleryButton(),
    ];

    if (receiptModel.imageBehaviorSubject.value.isNotEmpty ||
        receiptModel.imagesToUploadBehaviorSubject.value.isNotEmpty) {
      popupMenuEntries.add(
        buildDeleteButton(),
      );
    }

    return popupMenuEntries;
  }

  PopupMenuEntry buildDeleteButton() {
    return PopupMenuItem(
      child: const Text("Delete Image"),
      value: "delete",
      onTap: () async => await deleteImage(),
    );
  }

  Future<void> deleteImage() async {
    var formState = getFormStateFromContext(context);
    if (formState == WranglerFormState.add) {
      deleteImageFromModel();
      return;
    }

    if (formState == WranglerFormState.edit) {
      await deleteImageViaApi();
      return;
    }
  }

  void deleteImageFromModel() {
    var index = receiptModel.infiniteScrollController.selectedItem;
    var currentImages = List<UploadMultipartFileData>.from(
        receiptModel.imagesToUploadBehaviorSubject.value);
    currentImages.removeAt(index);
    receiptModel.imagesToUploadBehaviorSubject.add(currentImages);
  }

  Future<void> deleteImageViaApi() async {
    try {
      var index = receiptModel.infiniteScrollController.selectedItem;
      await OpenApiClient.client.getReceiptImageApi().deleteReceiptImageById(
          receiptImageId: receiptModel.imageBehaviorSubject.value[index]!.id);
      var currentImages =
          List<api.FileDataView?>.from(receiptModel.imageBehaviorSubject.value);
      currentImages.removeAt(index);
      receiptModel.imageBehaviorSubject.add(currentImages);

      showSuccessSnackbar(context, "Successfully deleted image");
    } catch (e) {
      showApiErrorSnackbar(context, e as DioException);
    }
  }

  PopupMenuEntry buildUploadFromGalleryButton() {
    return PopupMenuItem(
        value: "gallery",
        child: const Text("Upload from Gallery"),
        onTap: () async => await getImages(UploadMethod.gallery));
  }

  PopupMenuEntry buildUploadFromCameraButton() {
    return PopupMenuItem(
        value: "camera",
        child: const Text("Upload from Camera"),
        onTap: () async => await getImages(UploadMethod.camera));
  }

  Future<void> getImages(UploadMethod method) async {
    List<UploadMultipartFileData> imagesToUpload;
    var formState = getFormStateFromContext(context);

    if (method == UploadMethod.camera) {
      imagesToUpload = await scanImagesMultiPart(1);
    } else {
      imagesToUpload = await getGalleryImages(multiple: false);
    }

    if (formState == WranglerFormState.add) {
      addImagesToModel(imagesToUpload);
    } else if (formState == WranglerFormState.edit) {
      await uploadImages(imagesToUpload);
    }
  }

  void addImagesToModel(List<UploadMultipartFileData> imagesToUpload) {
    for (var image in imagesToUpload) {
      var currentList = receiptModel.imagesToUploadBehaviorSubject.value;
      receiptModel.imagesToUploadBehaviorSubject.add([...currentList, image]);
    }
  }

  Future<void> uploadImages(
      List<UploadMultipartFileData> imagesToUpload) async {
    var successMessage = "Successfully uploaded image";
    try {
      if (imagesToUpload.isNotEmpty) {
        Provider.of<LoadingModel>(context, listen: false).setIsLoading(true);
      }

      for (var image in imagesToUpload) {
        var uploadedImage = await OpenApiClient.client
            .getReceiptImageApi()
            .uploadReceiptImage(
                file: image.multipartFile, receiptId: receiptModel.receipt.id);
        var oldImages = List<api.FileDataView?>.from(
            receiptModel.imageBehaviorSubject.value);
        oldImages.add(uploadedImage.data);
        receiptModel.imageBehaviorSubject.add(oldImages);
      }
      Provider.of<LoadingModel>(context, listen: false).setIsLoading(false);

      if (imagesToUpload.isEmpty) {
        return;
      }
      if (imagesToUpload.length > 1) {
        successMessage =
            "Successfully uploaded ${imagesToUpload.length} images";
        return;
      }

      showSuccessSnackbar(context, successMessage);
    } catch (e) {
      print(e);
      Provider.of<LoadingModel>(context, listen: false).setIsLoading(false);
      showApiErrorSnackbar(context, e as DioException);
    }
  }
}
