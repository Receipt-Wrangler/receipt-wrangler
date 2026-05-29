import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';

/// Installs mock [MethodChannel] handlers for plugins that don't have a working
/// Linux desktop implementation in this project's runtime environment:
///
/// * **permission_handler** — ships Android/iOS natives. On Linux the
///   `requestPermissions`/`checkPermissionStatus` method channel has no
///   backing, so the fire-and-forget call from `main.dart` surfaces as an
///   unhandled async exception and fails the test.
/// * **gal** — image-gallery access (mobile-only). Called by the same
///   `requestPermissions()` helper as permission_handler. Stubbed to grant.
/// * **flutter_secure_storage** — has a Linux backend (libsecret) but it
///   requires an unlocked keyring + dbus session. Headless CI/containers
///   don't have that, so reads/writes throw `Libsecret error, Failed to
///   unlock the keyring`. The mock below backs the plugin with an in-memory
///   map — good enough for the smoke test, where we only care that the UI
///   writes + reads tokens correctly within one process.
///
/// Call this from `setUpAll` (or at the top of `main()`) before any test
/// pumps the app. On Android/iOS targets these plugins work natively, so
/// only install the mocks when running on desktop — today that's always the
/// case for our local integration_test runs; when Android lands, gate on
/// `Platform.isLinux`.
void installLinuxDesktopMocks() {
  final messenger = TestDefaultBinaryMessengerBinding
      .instance.defaultBinaryMessenger;

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
