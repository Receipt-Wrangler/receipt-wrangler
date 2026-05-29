import 'dart:io' show Platform;

import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';

import 'helpers/login.dart';
import 'helpers/platform_mocks.dart';

void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  // Per-test install so each gets a fresh in-memory secure-storage map
  // (a leaked JWT from the previous test would skip the login screens
  // and break this test's "Server URL" assertion).
  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('admin can log in from a cold boot', (tester) async {
    await loginAsAdmin(tester);
  });
}
