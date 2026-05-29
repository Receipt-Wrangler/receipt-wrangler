import 'dart:io';
import 'dart:typed_data';

import 'package:file_selector_platform_interface/file_selector_platform_interface.dart';
import 'package:flutter/services.dart' show rootBundle;
import 'package:plugin_platform_interface/plugin_platform_interface.dart';

/// Replaces [FileSelectorPlatform.instance] with a stub that returns a
/// fixed file on disk from [openFiles]. Lets tests exercise the receipt
/// form's "Upload from Gallery" path without driving a native picker.
///
/// The platform-interface swap works on every target Flutter supports, so
/// this is the same code path on Linux desktop, Android emulator, and
/// iOS simulator. Unlike a per-platform method-channel mock, we don't
/// have to discover three different channel names.
///
/// We materialize a real temp file on disk and return `XFile(path, ...)`
/// instead of `XFile.fromData(bytes, ...)`. Reason: on iOS Simulator,
/// `XFile.fromData()` returns an XFile whose `.name` getter is empty,
/// which then propagates through `MultipartFile(filename: ...)` as an
/// empty filename. The Go API's multipart parser closes the connection
/// on empty-filename uploads ("Connection closed before full header was
/// received"), dio surfaces that as a generic `DioException [unknown]`,
/// and the upload swallows it -- yielding a saved receipt with zero
/// imageFiles. Backing the XFile with a real file path makes both
/// `.name` and `.readAsBytes()` work identically on Android and iOS,
/// matching the production picker shape.
///
/// `MockPlatformInterfaceMixin` is what lets a test class extend
/// `FileSelectorPlatform` -- without it, `PlatformInterface.verify()`
/// rejects the swap.
class _FakeFileSelector extends FileSelectorPlatform
    with MockPlatformInterfaceMixin {
  _FakeFileSelector(this._path, this._name);

  final String _path;
  final String _name;

  @override
  Future<List<XFile>> openFiles({
    List<XTypeGroup>? acceptedTypeGroups,
    String? initialDirectory,
    String? confirmButtonText,
  }) async {
    return <XFile>[
      XFile(_path, name: _name, mimeType: 'image/png'),
    ];
  }
}

/// Loads `assets/test/sample.png`, writes it to a fresh tempdir under
/// `Directory.systemTemp`, and installs a fake file-selector that returns
/// an XFile backed by that on-disk path. The bundled asset (16x16 RGBA
/// PNG, ~300B) is shipped with the app so `rootBundle.load` works at
/// runtime on every target. Pass `bytes` to override with a different
/// fixture. The tempdir is process-scoped; the OS cleans it up.
Future<void> installFileSelectorMock({
  Uint8List? bytes,
  String name = 'sample.png',
}) async {
  final pngBytes = bytes ??
      (await rootBundle.load('assets/test/sample.png'))
          .buffer
          .asUint8List();
  final tempDir = await Directory.systemTemp.createTemp('file_selector_mock_');
  final tempFile = File('${tempDir.path}/$name');
  await tempFile.writeAsBytes(pngBytes, flush: true);
  FileSelectorPlatform.instance = _FakeFileSelector(tempFile.path, name);
}
