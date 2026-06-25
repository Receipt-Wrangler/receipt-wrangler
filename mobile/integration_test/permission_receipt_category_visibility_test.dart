import 'dart:io' show Platform;

import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';

import 'helpers/api.dart';
import 'helpers/form_actions.dart';
import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

/// Regression coverage for the receipt-form category picker (M2).
///
/// Categories/tags are now delivered per-group on `AppData` (`groupCategories`),
/// filtered to the caller's group-role grants. The flat `categories` array is
/// admin-only, so a normal (non-admin) user must source the picker from the
/// per-group catalog — otherwise the picker is empty. This drives that path with
/// a **Legacy Editor** (a non-admin app role: Legacy User, which lacks
/// `app.categories.read`) and asserts a freshly-created category shows up in the
/// receipt form's category multiselect.
///
/// Before the fix the mobile receipt pickers read the flat admin-only list, so
/// this exact flow surfaced an empty picker — this test would have caught it.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets(
    'receipt form: a non-admin sees the group catalog in the category picker',
    (tester) async {
      // A global category exists (unrestricted group roles see the whole pool).
      // Create it BEFORE login so it lands in the non-admin's AppData catalog.
      final adminJwt = await apiLogin();
      final categoryName =
          'e2e-cat-${DateTime.now().microsecondsSinceEpoch}';
      final categoryId = await createCategory(name: categoryName, jwt: adminJwt);
      addTearDown(() async => deleteCategory(categoryId, jwt: await apiLogin()));

      // A Legacy Editor (app role: Legacy User — no app.categories.read) who can
      // create receipts in the fixture group.
      final fixture = await provisionPermUser(roleName: 'Legacy Editor');
      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );
      await enterGroup(tester, fixture.groupName!);

      // Open the add-receipt form via the bottom-nav Add menu.
      await tester.tap(find.text('Add').hitTestable());
      await pumpUntilFound(
          tester, find.text('Add Manual Receipt').hitTestable());
      for (int i = 0; i < 5; i++) {
        await tester.pump(const Duration(milliseconds: 100));
      }
      await tester.tap(find.text('Add Manual Receipt').hitTestable());
      await pumpUntilFound(tester, find.text('Name'));

      // Selecting the group makes the category field visible (its Visibility is
      // gated on the group's receipt settings) and scopes the picker's source.
      await selectDropdown(tester, 'groupId', fixture.groupName!);

      // Open the category multiselect. The field renders the placeholder chip
      // "No Categories selected" inside a tappable GestureDetector.
      final categoryField = find.text('No Categories selected');
      await pumpUntilFound(tester, categoryField);
      await tester.ensureVisible(categoryField);
      await tester.pump(const Duration(milliseconds: 200));
      await tester.tap(categoryField);

      // The per-group catalog is sourced into the multiselect sheet, so the
      // category is offered. (Empty for a non-admin before the M2 fix.)
      await pumpUntilFound(tester, find.text(categoryName));
      expect(find.text(categoryName), findsWidgets);
    },
  );
}
