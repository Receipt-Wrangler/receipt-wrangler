import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/utils/permissions.dart';

import '../helpers/channel_mocks.dart';

/// `ensureCameraAccess` is what stands between a Scan tap and either the camera
/// or the gallery fallback, so each OS permission state has to map to exactly
/// one outcome -- and a state the user cannot be re-prompted out of must not be
/// re-requested (the request resolves instantly with no dialog, which would look
/// to the user like the tap did nothing).
void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  tearDown(clearPermissionMocks);

  test('granted maps to granted without a request', () async {
    final calls =
        installPermissionMocks(status: PermissionStatusWire.granted);

    expect(await ensureCameraAccess(), CameraAccess.granted);
    expect(calls.requests, 0);
  });

  test('limited (iOS partial access) still permits capture', () async {
    installPermissionMocks(status: PermissionStatusWire.limited);

    expect(await ensureCameraAccess(), CameraAccess.granted);
  });

  test('provisional counts as granted', () async {
    installPermissionMocks(status: PermissionStatusWire.provisional);

    expect(await ensureCameraAccess(), CameraAccess.granted);
  });

  test('undecided prompts once and reports what the user chose', () async {
    final calls = installPermissionMocks(
      status: PermissionStatusWire.denied,
      requestStatus: PermissionStatusWire.granted,
    );

    expect(await ensureCameraAccess(), CameraAccess.granted);
    expect(calls.requests, 1);
  });

  test('a declined prompt reports denied, not permanently denied', () async {
    installPermissionMocks(
      status: PermissionStatusWire.denied,
      requestStatus: PermissionStatusWire.denied,
    );

    expect(await ensureCameraAccess(), CameraAccess.denied,
        reason: 'the user can still be asked again, so no settings detour');
  });

  test('permanently denied is reported without re-requesting', () async {
    final calls = installPermissionMocks(
        status: PermissionStatusWire.permanentlyDenied);

    expect(await ensureCameraAccess(), CameraAccess.permanentlyDenied);
    expect(calls.requests, 0,
        reason: 'a request here resolves instantly with no dialog, so it would '
            'read to the user as the tap doing nothing');
  });

  test('restricted (parental controls / MDM) behaves as permanently denied',
      () async {
    final calls =
        installPermissionMocks(status: PermissionStatusWire.restricted);

    expect(await ensureCameraAccess(), CameraAccess.permanentlyDenied,
        reason: 'the user cannot grant it themselves');
    expect(calls.requests, 0);
  });

  test('concurrent callers share one request', () async {
    final calls = installPermissionMocks(
      status: PermissionStatusWire.denied,
      requestStatus: PermissionStatusWire.granted,
    );

    // Two overlapping requests on the platform channel throw
    // ERROR_ALREADY_REQUESTING_PERMISSIONS, which is why the single-flight guard
    // exists on this path as well as on requestPermissions().
    final results = await Future.wait([
      ensureCameraAccess(),
      ensureCameraAccess(),
    ]);

    expect(results, [CameraAccess.granted, CameraAccess.granted]);
    expect(calls.requests, 1);
  });

  test('the debug override short-circuits the platform channel', () async {
    final calls =
        installPermissionMocks(status: PermissionStatusWire.granted);
    debugCameraAccessOverride = () async => CameraAccess.permanentlyDenied;
    addTearDown(() => debugCameraAccessOverride = null);

    expect(await ensureCameraAccess(), CameraAccess.permanentlyDenied);
    expect(calls.statusChecks, 0);
  });
}
