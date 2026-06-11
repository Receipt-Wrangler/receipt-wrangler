import 'dart:io' show Platform;

import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';

import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

/// Permission-gating coverage for the add menu (`show_add_menu.dart`).
///
/// "Add Manual Receipt" is gated on `group.receipts.create`. Inside a group the
/// menu checks that permission for THAT group; on the group-select screen it
/// falls back to "held in any group". When no item is permitted it shows the
/// error snackbar instead of a menu.
///
/// Notes:
///   - Every account owns a personal "My Receipts" group (create allowed), so
///     the "held in any group" check is always satisfiable — we exercise that
///     branch positively from group-select and the *deny* path per-group via a
///     Legacy Viewer (read but not create).
///   - Inside a group two bottom navs are mounted (the group-select shell sits
///     under the group shell), so two "Add" destinations match; `.hitTestable()`
///     taps the visible one.
///   - Quick Scan is not asserted: it is additionally gated on the
///     `aiPoweredReceipts` feature flag (off here), so its absence is
///     flag-driven, not permission-driven.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  const noPermissionSnack = "You don't have permission to add receipts here.";

  Finder addButton() => find.text('Add').hitTestable();

  testWidgets(
    'add menu: a Legacy Viewer cannot add a receipt in their group',
    (tester) async {
      final fixture = await provisionPermUser(roleName: 'Legacy Viewer');
      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );
      await enterGroup(tester, fixture.groupName!);

      // Inside the group the menu checks group.receipts.create for THIS group;
      // a Viewer lacks it, so the menu is empty -> error snackbar.
      await tester.tap(addButton());
      await pumpUntilFound(tester, find.text(noPermissionSnack));
      // Presence, not an exact count: two bottom-nav shells are mounted in a
      // group, so the snackbar can render once per ScaffoldMessenger. The gate
      // outcome (denied -> snackbar, no menu) is what matters.
      expect(find.text(noPermissionSnack), findsWidgets);
      expect(find.text('Add Manual Receipt'), findsNothing);
    },
  );

  testWidgets(
    'add menu: a Legacy Editor can add a manual receipt in their group',
    (tester) async {
      final fixture = await provisionPermUser(roleName: 'Legacy Editor');
      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );
      await enterGroup(tester, fixture.groupName!);

      // An Editor holds group.receipts.create -> the menu offers the item.
      await tester.tap(addButton());
      await pumpUntilFound(tester, find.text('Add Manual Receipt'));
      expect(find.text('Add Manual Receipt'), findsOneWidget);
      expect(find.text(noPermissionSnack), findsNothing);
    },
  );

  testWidgets(
    'add menu: group-select offers add when the user can create in any group',
    (tester) async {
      // A fresh account is the owner of its personal "My Receipts" group, so
      // the group-select "held in any group" check passes and the menu offers
      // the item (no current group is selected here).
      final fixture = await provisionPermUser();
      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );

      await tester.tap(addButton());
      await pumpUntilFound(tester, find.text('Add Manual Receipt'));
      expect(find.text('Add Manual Receipt'), findsOneWidget);
      expect(find.text(noPermissionSnack), findsNothing);
    },
  );
}
