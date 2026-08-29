import 'package:flutter/foundation.dart';
import 'package:gal/gal.dart';
import 'package:permission_handler/permission_handler.dart';

Future<void>? _inFlight;

/// Requests camera + gallery permissions.
///
/// Concurrent callers share the same in-flight Future. Without this, two
/// invocations on the same process (e.g. integration_test runs that boot
/// `app.main()` twice when chaining tests) race on the platform channel
/// and the second one throws `PlatformException(ERROR_ALREADY_REQUESTING_PERMISSIONS)`.
Future<void> requestPermissions() {
  return _inFlight ??= () async {
    try {
      await Permission.camera.request();
      await Gal.requestAccess();
    } finally {
      _inFlight = null;
    }
  }();
}

/// The camera permission outcome the scan entry points branch on.
///
/// [permanentlyDenied] is separated from [denied] because a re-request in that
/// state resolves instantly without showing a dialog — the only way out is the
/// OS settings screen, so the UI must offer that instead of silently retrying.
enum CameraAccess { granted, denied, permanentlyDenied }

/// Test seam for [ensureCameraAccess].
///
/// The permission plugins are statics, so there is no constructor to inject.
/// This mirrors the two existing precedents in the codebase: the settable
/// `OpenApiClient.client` static, and the `debugForce*` flags on
/// `QrScannerScreen`. Widget tests set this to drive each branch without the
/// platform channel; production leaves it null.
@visibleForTesting
Future<CameraAccess> Function()? debugCameraAccessOverride;

Future<CameraAccess>? _cameraRequestInFlight;

/// Resolves whether the camera may be used, requesting permission once if the
/// user has not decided yet.
///
/// Called lazily, from the scan action — **never at launch**. A launch-time
/// request was deliberately removed for the iOS 26.x render-pause freeze
/// (see `lib/main.dart` and GitHub #617).
///
/// Concurrent callers share one in-flight request for the same reason
/// [requestPermissions] does: two overlapping requests on the platform channel
/// throw `ERROR_ALREADY_REQUESTING_PERMISSIONS`.
Future<CameraAccess> ensureCameraAccess() async {
  final override = debugCameraAccessOverride;
  if (override != null) {
    return override();
  }

  final current = _mapStatus(await Permission.camera.status);
  // Only an undecided permission is worth a request: granted needs nothing, and
  // a permanently-denied request resolves instantly without a dialog.
  if (current != CameraAccess.denied) {
    return current;
  }

  return _cameraRequestInFlight ??= () async {
    try {
      return _mapStatus(await Permission.camera.request());
    } finally {
      _cameraRequestInFlight = null;
    }
  }();
}

CameraAccess _mapStatus(PermissionStatus status) {
  // `limited` (iOS partial access) still permits capture, so it counts as
  // granted. `restricted` (parental controls / MDM) cannot be requested away by
  // the user, so it behaves like permanently denied.
  if (status.isGranted || status.isLimited || status.isProvisional) {
    return CameraAccess.granted;
  }
  if (status.isPermanentlyDenied || status.isRestricted) {
    return CameraAccess.permanentlyDenied;
  }
  return CameraAccess.denied;
}
