// The camera-denied fallback, end to end.
//
// A denied camera used to produce nothing at all -- getPictures returned an
// empty list (or threw) and the tap looked like it had done nothing. It now
// explains itself and continues with the gallery, and the permanently-denied
// case additionally offers the OS settings screen, which is the only way out of
// that state.
//
// The permission states are driven through permission_handler's method channel
// rather than a real OS dialog: the suite uses the first-party integration_test
// package, which cannot drive native permission dialogs (that would be Patrol).

import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_entry.dart';
import 'package:receipt_wrangler_mobile/utils/permissions.dart';

import 'helpers/feature_flags.dart';
import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/receipt_test_helpers.dart';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  tearDown(() => debugCameraAccessOverride = null);

  /// Logs in a Legacy Editor with the AI flag on, inside their group, with the
  /// camera reporting [access].
  Future<void> loginWithCamera(
    WidgetTester tester,
    CameraAccess access,
  ) async {
    await enableAiPoweredReceiptsForTest();
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    // Set before the tap rather than by swapping the channel mock: the login
    // bootstrap also touches permission_handler, and this pins only the branch
    // under test.
    debugCameraAccessOverride = () async => access;

    final fixture = await provisionPermUser(roleName: 'Legacy Editor');
    await loginAs(
      tester,
      username: fixture.username,
      password: fixture.password,
    );
    await enterGroup(tester, fixture.groupName!);
  }

  testWidgets('a denied camera explains itself and moves on', (tester) async {
    await loginWithCamera(tester, CameraAccess.denied);

    await pumpUntilFound(tester, scanNavSlot());
    await tester.tap(scanNavSlot());

    await pumpUntilFound(tester, find.text(cameraDeniedFallbackMessage));
    expect(find.widgetWithText(SnackBarAction, 'Settings'), findsNothing,
        reason: 'the user can still be re-prompted, so no settings detour');
  });

  testWidgets('a permanently denied camera offers the settings escape hatch',
      (tester) async {
    await loginWithCamera(tester, CameraAccess.permanentlyDenied);

    await pumpUntilFound(tester, scanNavSlot());
    await tester.tap(scanNavSlot());

    await pumpUntilFound(tester, find.text(cameraDeniedFallbackMessage));
    expect(find.widgetWithText(SnackBarAction, 'Settings'), findsOneWidget,
        reason: 'a re-request resolves with no dialog here, so without this '
            'the tap would read as doing nothing');
  });
}
