import 'dart:typed_data';
import 'package:dio/dio.dart';

class UploadMultipartFileData {
  MultipartFile multipartFile;
  Uint8List bytes;

  UploadMultipartFileData({required this.multipartFile, required this.bytes});
}
