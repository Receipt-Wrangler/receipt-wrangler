import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';

/// Grants camera (**permission_handler**) and image-gallery (**gal**) access by
/// stubbing their method channels to a "granted"-equivalent response. Both
/// channels are touched by `main.dart`'s init-time `requestPermissions()`, and
/// permission_handler is *also* hit directly by `cunning_document_scanner`'s
/// `getPictures()` — so this is shared by [installLinuxDesktopMocks] (which has
/// no native backing for these on desktop) and [installDocumentScannerMock]
/// (which needs the grant on **every** platform to drive the scanner headlessly;
/// see that helper for the iOS/Android `ERROR_ALREADY_REQUESTING_PERMISSIONS`
/// race it avoids). Idempotent — `setMockMethodCallHandler` just replaces any
/// prior handler with an equivalent one.
///
/// `requestPermissions` returns an EMPTY map deliberately: permission_handler
/// decodes an absent entry as "denied", but our only consumers either ignore the
/// result (`main.dart`) or check solely for the *presence* of a denied value
/// (`cunning_document_scanner`), so an empty map reads as "nothing denied"
/// without pinning a specific `PermissionStatus` wire int.
void installCameraGalleryPermissionMocks() {
  final messenger =
      TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger;

  const permissions = MethodChannel('flutter.baseflow.com/permissions/methods');
  messenger.setMockMethodCallHandler(permissions, (call) async {
    switch (call.method) {
      case 'requestPermissions':
        return <int, int>{};
      case 'checkPermissionStatus':
      case 'checkServiceStatus':
        return 1;
      default:
        return null;
    }
  });

  const gal = MethodChannel('gal');
  messenger.setMockMethodCallHandler(gal, (call) async {
    switch (call.method) {
      case 'requestAccess':
      case 'hasAccess':
        return true;
      default:
        return null;
    }
  });
}

/// Installs mock [MethodChannel] handlers for plugins that don't have a working
/// Linux desktop implementation in this project's runtime environment:
///
/// * **permission_handler** / **gal** — see [installCameraGalleryPermissionMocks].
///   On Linux these method channels have no backing, so the fire-and-forget call
///   from `main.dart` surfaces as an unhandled async exception and fails the test.
/// * **flutter_secure_storage** — has a Linux backend (libsecret) but it
///   requires an unlocked keyring + dbus session. Headless CI/containers
///   don't have that, so reads/writes throw `Libsecret error, Failed to
///   unlock the keyring`. The mock below backs the plugin with an in-memory
///   map — good enough for the smoke test, where we only care that the UI
///   writes + reads tokens correctly within one process.
///
/// Call this from `setUpAll` (or at the top of `main()`) before any test
/// pumps the app. On Android/iOS targets these plugins work natively, so
/// only install the mocks when running on desktop — gate on `Platform.isLinux`.
void installLinuxDesktopMocks() {
  installCameraGalleryPermissionMocks();

  final messenger = TestDefaultBinaryMessengerBinding
      .instance.defaultBinaryMessenger;

  final storage = <String, String>{};
  const secureStorage =
      MethodChannel('plugins.it_nomads.com/flutter_secure_storage');
  messenger.setMockMethodCallHandler(secureStorage, (call) async {
    final args = (call.arguments as Map?)?.cast<String, dynamic>() ?? {};
    final key = args['key'] as String?;
    switch (call.method) {
      case 'write':
        if (key != null) storage[key] = args['value'] as String? ?? '';
        return null;
      case 'read':
        return key == null ? null : storage[key];
      case 'readAll':
        return Map<String, String>.from(storage);
      case 'delete':
        if (key != null) storage.remove(key);
        return null;
      case 'deleteAll':
        storage.clear();
        return null;
      case 'containsKey':
        return key != null && storage.containsKey(key);
      default:
        return null;
    }
  });
}
