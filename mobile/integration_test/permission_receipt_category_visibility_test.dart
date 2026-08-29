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
import 'helpers/receipt_test_helpers.dart';

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

  // The category and tag picker flows are identical except for the catalog
  // create/delete helpers, the placeholder text, and the asserted name — so the
  // shared "non-admin sees the per-group catalog in the picker" flow lives here.
  // A global catalog item is created BEFORE login so it lands in the non-admin's
  // AppData catalog; a Legacy Editor (app role Legacy User, no
  // app.categories.read / app.tags.read) then opens the add-receipt form and the
  // multiselect, and the item must be offered (it would be empty if the picker
  // read the flat admin-only list instead of the per-group catalog).
  Future<void> assertGroupCatalogPicker(
    WidgetTester tester, {
    required String namePrefix,
    required String placeholder,
    required Future<int> Function({required String name, required String jwt})
        createFn,
    required Future<void> Function(int id, {required String jwt}) deleteFn,
  }) async {
    final adminJwt = await apiLogin();
    final itemName = '$namePrefix-${DateTime.now().microsecondsSinceEpoch}';
    final itemId = await createFn(name: itemName, jwt: adminJwt);
    addTearDown(() async => deleteFn(itemId, jwt: await apiLogin()));

    final fixture = await provisionPermUser(roleName: 'Legacy Editor');
    await loginAs(
      tester,
      username: fixture.username,
      password: fixture.password,
    );
    await enterGroup(tester, fixture.groupName!);

    // Open the add-receipt form via the bottom-nav Add menu.
    await openManualReceiptForm(tester);
    await pumpUntilFound(tester, find.text('Name'));

    // Selecting the group makes the picker field visible (gated on the group's
    // receipt settings) and scopes the picker's source to that group's catalog.
    await selectDropdown(tester, 'groupId', fixture.groupName!);

    // Open the multiselect via its placeholder chip (a tappable GestureDetector).
    final field = find.text(placeholder);
    await pumpUntilFound(tester, field);
    await tester.ensureVisible(field);
    await tester.pump(const Duration(milliseconds: 200));
    await tester.tap(field);

    // The per-group catalog is sourced into the multiselect sheet, so the item
    // is offered. (Empty for a non-admin before the M2 fix.)
    await pumpUntilFound(tester, find.text(itemName));
    expect(find.text(itemName), findsWidgets);
  }

  testWidgets(
    'receipt form: a non-admin sees the group catalog in the category picker',
    (tester) async {
      await assertGroupCatalogPicker(
        tester,
        namePrefix: 'e2e-cat',
        placeholder: 'No Categories selected',
        createFn: createCategory,
        deleteFn: deleteCategory,
      );
    },
  );

  testWidgets(
    'receipt form: a non-admin sees the group catalog in the tag picker',
    (tester) async {
      await assertGroupCatalogPicker(
        tester,
        namePrefix: 'e2e-tag',
        placeholder: 'No Tags selected',
        createFn: createTag,
        deleteFn: deleteTag,
      );
    },
  );
}
