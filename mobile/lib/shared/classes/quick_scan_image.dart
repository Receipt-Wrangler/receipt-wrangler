import 'dart:typed_data';

import 'package:dio/dio.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/interfaces/upload_multipart_file_data.dart';

class QuickScanImage extends UploadMultipartFileData {
  QuickScanImage(
      {required this.multipartFile,
      required this.bytes,
      required this.formKey,
      this.groupId,
      this.paidByUserId,
      this.status})
      : super(multipartFile: multipartFile, bytes: bytes);

  final MultipartFile multipartFile;

  final Uint8List bytes;

  final GlobalKey<FormBuilderState> formKey;

  int? groupId;

  int? paidByUserId;

  ReceiptStatus? status;

  static QuickScanImage fromUploadMultipartFileData(
      UploadMultipartFileData data,
      int? initialGroupId,
      int? initialPaidByUserId,
      ReceiptStatus? initialStatus) {
    return QuickScanImage(
        multipartFile: data.multipartFile,
        bytes: data.bytes,
        formKey: GlobalKey<FormBuilderState>(),
        groupId: initialGroupId,
        paidByUserId: initialPaidByUserId,
        status: initialStatus);
  }
}
