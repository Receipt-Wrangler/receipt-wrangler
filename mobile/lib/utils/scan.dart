import 'dart:io';

import 'package:cunning_document_scanner/cunning_document_scanner.dart';
import 'package:file_selector/file_selector.dart';
import 'package:dio/dio.dart' as dio;
import 'package:http/http.dart';
import 'package:receipt_wrangler_mobile/interfaces/upload_multipart_file_data.dart';

const multipleFileFieldName = "files";
const singularFileFieldName = "file";

Future<List<String>> scanImages(int numberOfPages) async {
  return await CunningDocumentScanner.getPictures(noOfPages: numberOfPages) ??
      [];
}

Future<List<UploadMultipartFileData>> scanImagesMultiPart(
    int numberOfPages) async {
  var files = <UploadMultipartFileData>[];
  var filePaths =
      await CunningDocumentScanner.getPictures(noOfPages: numberOfPages);

  if (filePaths!.isEmpty) {
    return files;
  }

  var multiple = numberOfPages > 1;
  for (var filePath in filePaths) {
    var multipartFile =
        await MultipartFile.fromPath(getFieldName(multiple), filePath);
    var bytes = await File(filePath).readAsBytes();

    var dioMultipartFile =
        dio.MultipartFile.fromBytes(bytes, filename: multipartFile.filename);

    files.add(
        UploadMultipartFileData(multipartFile: dioMultipartFile, bytes: bytes));
  }

  return files;
}

Future<List<UploadMultipartFileData>> getGalleryImages(
    {multiple = true}) async {
  var files = <UploadMultipartFileData>[];

  const XTypeGroup typeGroup = XTypeGroup();
  List<XFile> openedFiles = [];

  switch (Platform.operatingSystem) {
    case "android":
      openedFiles = await openAndroidGallery();
      break;
    case "ios":
      openedFiles = await openIOSGallery();
      break;
    default:
      throw Exception("Unsupported platform");
  }

  for (var file in openedFiles) {
    var bytes = await file.readAsBytes();
    var multipartFile = await MultipartFile.fromBytes(
        getFieldName(multiple), bytes,
        filename: file.name);

    var dioMultipartFile =
        dio.MultipartFile.fromBytes(bytes, filename: multipartFile.filename);
    files.add(
        UploadMultipartFileData(multipartFile: dioMultipartFile, bytes: bytes));
  }

  return files;
}

getFieldName(bool multiple) {
  return multiple ? multipleFileFieldName : singularFileFieldName;
}

Future<List<XFile>> openIOSGallery() async {
  const typeGroup = XTypeGroup(
    uniformTypeIdentifiers: ["public.image", "com.adobe.pdf"],
  );
  return await openFiles(acceptedTypeGroups: [typeGroup]);
}

Future<List<XFile>> openAndroidGallery() async {
  const typeGroup = XTypeGroup(
    mimeTypes: ["image/*", "application/pdf"],
  );
  return await openFiles(
    acceptedTypeGroups: [typeGroup],
  );
}
