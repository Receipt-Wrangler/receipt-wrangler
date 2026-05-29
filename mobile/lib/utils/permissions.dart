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
