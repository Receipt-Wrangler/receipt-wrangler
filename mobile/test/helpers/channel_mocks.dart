import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:permission_handler/permission_handler.dart' show Permission;

/// Platform-channel stubs for the plugins the receipt-scan paths reach.
///
/// These plugins are statics with no injectable seam, so a test controls them by
/// answering their method channels. Shared by the widget suite (`test/`) and the
/// e2e suite (`integration_test/helpers/platform_mocks.dart` re-exports them) so
/// both drive the same permission states through the same wire values.

const _permissionsChannel =
    MethodChannel('flutter.baseflow.com/permissions/methods');
const _galChannel = MethodChannel('gal');
const _documentScannerChannel = MethodChannel('cunning_document_scanner');

/// The wire integers `permission_handler` decodes into `PermissionStatus`.
///
/// They are the enum's *declaration order* in
/// `permission_handler_platform_interface/lib/src/permission_status.dart`, which
/// is the plugin's actual wire contract. Named here so a package bump that
/// reorders them is a one-line fix rather than a hunt through mock literals.
class PermissionStatusWire {
  static const denied = 0;
  static const granted = 1;
  static const restricted = 2;
  static const limited = 3;
  static const permanentlyDenied = 4;
  static const provisional = 5;
}

/// Records what a permission mock was asked to do, so a test can assert that
/// e.g. an already-permanently-denied permission was never re-requested.
class PermissionMockCalls {
  int statusChecks = 0;
  int requests = 0;
  int settingsOpens = 0;
}

/// Answers the `permission_handler` channel with [status] for every permission.
///
/// [requestStatus] is what a *request* resolves to, when it differs from the
/// initial [status] — i.e. the user answering the OS dialog. Defaults to
/// [status], which models a user who doesn't change their mind.
///
/// Returns a recorder for the assertions above. Callers should
/// `addTearDown(clearPermissionMocks)`.
PermissionMockCalls installPermissionMocks({
  int status = PermissionStatusWire.granted,
  int? requestStatus,
}) {
  final calls = PermissionMockCalls();
  final messenger =
      TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger;

  messenger.setMockMethodCallHandler(_permissionsChannel, (call) async {
    switch (call.method) {
      case 'checkPermissionStatus':
        calls.statusChecks++;
        return status;
      case 'requestPermissions':
        calls.requests++;
        return <int, int>{Permission.camera.value: requestStatus ?? status};
      case 'checkServiceStatus':
        return PermissionStatusWire.granted;
      case 'openAppSettings':
        calls.settingsOpens++;
        return true;
      default:
        return null;
    }
  });

  return calls;
}

/// Grants camera (**permission_handler**) and image-gallery (**gal**) access.
///
/// `requestPermissions` returns an EMPTY map deliberately for the gal-paired
/// callers: permission_handler decodes an absent entry as "denied", but the
/// consumers that rely on this variant either ignore the result (`main.dart`) or
/// check solely for the *presence* of a denied value (`cunning_document_scanner`),
/// so an empty map reads as "nothing denied" without pinning a wire int.
/// Idempotent — `setMockMethodCallHandler` replaces any prior handler.
void installCameraGalleryPermissionMocks() {
  final messenger =
      TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger;

  messenger.setMockMethodCallHandler(_permissionsChannel, (call) async {
    switch (call.method) {
      case 'requestPermissions':
        return <int, int>{};
      case 'checkPermissionStatus':
      case 'checkServiceStatus':
        return PermissionStatusWire.granted;
      default:
        return null;
    }
  });

  messenger.setMockMethodCallHandler(_galChannel, (call) async {
    switch (call.method) {
      case 'requestAccess':
      case 'hasAccess':
        return true;
      default:
        return null;
    }
  });
}

/// Answers `cunning_document_scanner.getPictures` with [paths].
///
/// An empty list is what the plugin returns when the user cancels the scanner,
/// so it is the way to drive that branch.
void installDocumentScannerChannelMock(List<String> paths) {
  TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger
      .setMockMethodCallHandler(_documentScannerChannel, (call) async {
    if (call.method == 'getPictures') {
      return paths;
    }
    return null;
  });
}

/// Detaches every handler installed above.
void clearPermissionMocks() {
  final messenger =
      TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger;
  messenger.setMockMethodCallHandler(_permissionsChannel, null);
  messenger.setMockMethodCallHandler(_galChannel, null);
  messenger.setMockMethodCallHandler(_documentScannerChannel, null);
}
