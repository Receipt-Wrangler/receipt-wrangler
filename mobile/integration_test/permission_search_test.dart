import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';

import 'helpers/login.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

/// Permission-gating coverage for the Search bottom-nav destination, gated on the
/// app-scoped `app.receipts.search` (`group_bottom_nav.dart` /
/// `group_select_bottom_nav.dart`). A caller without it sees no Search
/// destination at all (and the `/search` route is additionally redirected by
/// `receiptsSearchRedirect`, but the nav destination is the user-facing surface).
///
/// We assert on the group-select bottom nav, which mounts immediately on the
/// post-login `GroupSelect` landing — no group needed, and "Search" text appears
/// only in the bottom nav on mobile. The deny case needs a custom app role
/// (every seeded role holds `app.receipts.search`), provisioned and torn down by
/// `provisionUserWithoutAppPermission`.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets(
    'search nav: hidden for a user without app.receipts.search',
    (tester) async {
      final fixture =
          await provisionUserWithoutAppPermission('app.receipts.search');
      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );

      // The group-select bottom nav is up, but the Search destination is gated
      // out for a caller missing app.receipts.search.
      await pumpUntilFound(tester, find.byType(NavigationBar));
      expect(find.text('Search'), findsNothing);
      expect(find.byIcon(Icons.search), findsNothing);
    },
  );

  testWidgets(
    'search nav: visible for a Legacy User (holds app.receipts.search)',
    (tester) async {
      // The default e2e-user is a Legacy User, which holds app.receipts.search.
      await loginAsUser(tester);

      await pumpUntilFound(tester, find.byType(NavigationBar));
      expect(find.text('Search'), findsOneWidget);
      expect(find.byIcon(Icons.search), findsOneWidget);
    },
  );
}
