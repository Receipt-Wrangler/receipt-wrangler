import 'dart:io';
import 'dart:typed_data';

import 'package:flutter/services.dart' show MethodChannel, rootBundle;
import 'package:flutter_test/flutter_test.dart';

/// Stubs the `cunning_document_scanner` method channel so Quick Scan's
/// "add a photo" action (`scanImagesMultiPart` → `CunningDocumentScanner.getPictures`)
/// returns a fixed on-disk image instead of driving the native camera scanner.
///
/// This is the *only* way to feed an image into the Quick Scan sheet on Linux
/// desktop: the sheet's other upload icon (`getGalleryImages`) hard-throws
/// "Unsupported platform" via a `Platform.operatingSystem` switch before it ever
/// reaches `file_selector`, so the `file_selector` mock can't help there. The
/// scanner path only touches this method channel plus a `permission_handler`
/// camera request (already granted by `installLinuxDesktopMocks`), so mocking the
/// channel reaches the real per-image form headlessly, with no production change.
///
/// A real temp file is written to disk (not `XFile.fromData`) so `MultipartFile.fromPath`
/// / `File.readAsBytes` behave exactly like the production picker. The bundled
/// `assets/test/sample.png` (16x16 RGBA PNG) ships with the app, so `rootBundle.load`
/// works on every target; the process-scoped tempdir is cleaned up by the OS.
Future<void> installDocumentScannerMock({
  Uint8List? bytes,
  String name = 'sample.png',
}) async {
  final pngBytes = bytes ??
      (await rootBundle.load('assets/test/sample.png')).buffer.asUint8List();
  final tempDir = await Directory.systemTemp.createTemp('doc_scanner_mock_');
  final tempFile = File('${tempDir.path}/$name');
  await tempFile.writeAsBytes(pngBytes, flush: true);

  const channel = MethodChannel('cunning_document_scanner');
  TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger
      .setMockMethodCallHandler(channel, (call) async {
    if (call.method == 'getPictures') {
      return <String>[tempFile.path];
    }
    return null;
  });
}
