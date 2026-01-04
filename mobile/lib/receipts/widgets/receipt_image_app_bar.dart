import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:gal/gal.dart';
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
import '../../utils/scan.dart';
import '../../utils/snackbar.dart';

class ReceiptImageAppBar extends StatelessWidget implements PreferredSizeWidget {
  const ReceiptImageAppBar({
    super.key,
    required this.formState,
    required this.title,
    required this.onBack,
    this.onEditMode,
  });

  final WranglerFormState formState;
  final String title;
  final VoidCallback onBack;
  final VoidCallback? onEditMode;

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight);

  @override
  Widget build(BuildContext context) {
    final receiptModel = Provider.of<ReceiptModel>(context, listen: false);
    
    return AppBar(
      title: Text(title),
      leading: IconButton(
        icon: const Icon(Icons.arrow_back),
        onPressed: onBack,
      ),
      actions: [_buildImageActions(context, receiptModel)],
    );
  }

  Widget _buildImageActions(BuildContext context, ReceiptModel receiptModel) {
    List<PopupMenuEntry> options = [];
    
    if (_isEditMode()) {
      options = _buildEditModeActions(context, receiptModel);
    } else {
      options = _buildViewModeActions(context, receiptModel);
    }

    options.add(_buildImageDownloadButton(context, receiptModel));
    options.add(_buildViewInFullScreenButton(context, receiptModel));

    var combinedStream = Rx.merge([
      receiptModel.imagesToUploadBehaviorSubject.stream,
      receiptModel.imageBehaviorSubject.stream
    ]).asBroadcastStream();

    return StreamBuilder(
        stream: combinedStream,
        builder: (context, snapshot) {
          return ReceiptEditPopupMenu(
              groupId: receiptModel.receipt.groupId,
              popupMenuChildren: options,
              formState: formState);
        });
  }

  bool _isEditMode() {
    return formState == WranglerFormState.edit || formState == WranglerFormState.add;
  }

  List<PopupMenuEntry> _buildViewModeActions(BuildContext context, ReceiptModel receiptModel) {
    return [
      PopupMenuItem(
          value: "edit",
          child: const Text("Edit"),
          onTap: () {
            if (onEditMode != null) {
              onEditMode!();
            }
          }),
    ];
  }

  List<PopupMenuEntry> _buildEditModeActions(BuildContext context, ReceiptModel receiptModel) {
    List<PopupMenuEntry> popupMenuEntries = [
      _buildUploadFromCameraButton(context, receiptModel),
      _buildUploadFromGalleryButton(context, receiptModel),
    ];

    if (receiptModel.imageBehaviorSubject.value.isNotEmpty ||
        receiptModel.imagesToUploadBehaviorSubject.value.isNotEmpty) {
      popupMenuEntries.add(
        _buildDeleteButton(context, receiptModel),
      );
    }

    return popupMenuEntries;
  }

  PopupMenuItem _buildViewInFullScreenButton(BuildContext context, ReceiptModel receiptModel) {
    return PopupMenuItem(
        child: Text("View in full screen"),
        enabled: formState == WranglerFormState.add ? _areImagesToUpload(receiptModel) : true,
        onTap: () {
          var selectedIndex = receiptModel.infiniteScrollController.selectedItem;
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

  PopupMenuItem _buildImageDownloadButton(BuildContext context, ReceiptModel receiptModel) {
    return PopupMenuItem(
      child: const Text("Download"),
      enabled: formState == WranglerFormState.add ? _areImagesToUpload(receiptModel) : true,
      onTap: () async => await _downloadImage(context, receiptModel),
    );
  }

  Future _downloadImage(BuildContext context, ReceiptModel receiptModel) async {
    final loadingModel = Provider.of<LoadingModel>(context, listen: false);
    loadingModel.setIsLoading(true);
    
    var receiptImage = _getCurrentlySelectedImage(receiptModel);
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

  PopupMenuEntry _buildDeleteButton(BuildContext context, ReceiptModel receiptModel) {
    return PopupMenuItem(
      child: const Text("Delete Image"),
      value: "delete",
      onTap: () async => await _deleteImage(context, receiptModel),
    );
  }

  Future<void> _deleteImage(BuildContext context, ReceiptModel receiptModel) async {
    if (formState == WranglerFormState.add) {
      _deleteImageFromModel(receiptModel);
      return;
    }

    if (formState == WranglerFormState.edit) {
      await _deleteImageViaApi(context, receiptModel);
      return;
    }
  }

  void _deleteImageFromModel(ReceiptModel receiptModel) {
    var index = receiptModel.infiniteScrollController.selectedItem;
    var currentImages = List<UploadMultipartFileData>.from(
        receiptModel.imagesToUploadBehaviorSubject.value);
    currentImages.removeAt(index);
    receiptModel.imagesToUploadBehaviorSubject.add(currentImages);
  }

  Future<void> _deleteImageViaApi(BuildContext context, ReceiptModel receiptModel) async {
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

  PopupMenuEntry _buildUploadFromGalleryButton(BuildContext context, ReceiptModel receiptModel) {
    return PopupMenuItem(
        value: "gallery",
        child: const Text("Upload from Gallery"),
        onTap: () async => await _getImages(context, receiptModel, UploadMethod.gallery));
  }

  PopupMenuEntry _buildUploadFromCameraButton(BuildContext context, ReceiptModel receiptModel) {
    return PopupMenuItem(
        value: "camera",
        child: const Text("Upload from Camera"),
        onTap: () async => await _getImages(context, receiptModel, UploadMethod.camera));
  }

  Future<void> _getImages(BuildContext context, ReceiptModel receiptModel, UploadMethod method) async {
    List<UploadMultipartFileData> imagesToUpload;

    if (method == UploadMethod.camera) {
      imagesToUpload = await scanImagesMultiPart(1);
    } else {
      imagesToUpload = await getGalleryImages(multiple: false);
    }

    if (formState == WranglerFormState.add) {
      _addImagesToModel(receiptModel, imagesToUpload);
    } else if (formState == WranglerFormState.edit) {
      await _uploadImages(context, receiptModel, imagesToUpload);
    }
  }

  void _addImagesToModel(ReceiptModel receiptModel, List<UploadMultipartFileData> imagesToUpload) {
    for (var image in imagesToUpload) {
      var currentList = receiptModel.imagesToUploadBehaviorSubject.value;
      receiptModel.imagesToUploadBehaviorSubject.add([...currentList, image]);
    }
  }

  Future<void> _uploadImages(BuildContext context, ReceiptModel receiptModel,
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

  bool _areImagesToUpload(ReceiptModel receiptModel) {
    return !receiptModel.imagesToUploadBehaviorSubject.value.isEmpty;
  }

  api.FileDataView? _getCurrentlySelectedImage(ReceiptModel receiptModel) {
    var index = receiptModel.infiniteScrollController.selectedItem;
    return receiptModel.imageBehaviorSubject.value[index];
  }
}