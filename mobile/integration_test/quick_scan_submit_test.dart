// Quick Scan SUBMIT outcomes -- the per-image form actually POSTs to the API.
//
// Unlike the config-response spec (which asserts visibility / client-side
// blocking without hitting the backend), a successful submit is validated by the
// backend's resolveQuickScanFields against the group's PERSISTED config. So these
// tests PERSIST the group's quick-scan config via the API (client reads it from
// AppData at login; backend validates the same -- exactly as in production),
// rather than mutating GroupModel client-side. Config is restored on teardown.
//
// Paid-by/status stay required (making them optional needs a persisted default
// the backend enforces), so both tests fill them; the CATEGORY axis carries the
// variation:
//   * categories shown+OPTIONAL, left empty, still queues  -> optional-empty OK,
//   * categories shown+REQUIRED, a picked category, queues  -> the category
//     picker path (which crashed before the ContextModel.resolveSheetContext fix).
//
// Successful submits queue an async receipt with no deterministic id, so (like
// quick_scan_test) there is no receipt cleanup -- only the seeded category and
// the persisted config are torn down.

import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/category_select_field.dart';

import 'helpers/api.dart';
import 'helpers/document_scanner_mock.dart';
import 'helpers/feature_flags.dart';
import 'helpers/form_actions.dart';
import 'helpers/login.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/quick_scan_actions.dart';
import 'helpers/users.dart';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('required fields filled + optional category empty queues the scan',
      (tester) async {
    await enableAiPoweredReceiptsForTest();
    await installDocumentScannerMock();
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    // Persist: paid-by/status required, categories shown+OPTIONAL.
    final jwt = await apiLogin();
    final g = await firstNonAllGroup(jwt);
    await setGroupQuickScanConfig(groupId: g.id, jwt: jwt, overrides: const {
      'quickScanPaidByEnabled': true,
      'quickScanPaidByRequired': true,
      'quickScanStatusEnabled': true,
      'quickScanStatusRequired': true,
      'quickScanCategoriesEnabled': true,
      'quickScanCategoriesRequired': false,
      'quickScanTagsEnabled': false,
      'quickScanTagsRequired': false,
    });

    await loginAsAdmin(tester);
    final target = await keepOnlyGroup(tester, g.id);

    await openQuickScanImageForm(tester);
    await selectDropdown(tester, 'groupId', target.name);

    // Fill the required fields; leave the optional categories empty.
    await selectDropdown(tester, 'paidByUserId', adminDisplayName(tester));
    await selectDropdown(tester, 'status', 'Open');
    expect(find.byType(CategorySelectField), findsOneWidget,
        reason: 'categories shown (optional)');

    // Required filled + optional left empty → the scan queues.
    await expectQuickScanQueued(tester);
  });

  testWidgets('required category selected via the picker queues the scan',
      (tester) async {
    await enableAiPoweredReceiptsForTest();
    await installDocumentScannerMock();
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    // Seed a category (unrestricted → in the admin's per-group catalog) and
    // persist: paid-by/status required, categories shown+REQUIRED.
    final jwt = await apiLogin();
    final g = await firstNonAllGroup(jwt);
    final categoryName = 'e2e-qs-cat-${DateTime.now().millisecondsSinceEpoch}';
    final categoryId = await createCategory(name: categoryName, jwt: jwt);
    addTearDown(() async => deleteCategory(categoryId, jwt: await apiLogin()));
    await setGroupQuickScanConfig(groupId: g.id, jwt: jwt, overrides: const {
      'quickScanPaidByEnabled': true,
      'quickScanPaidByRequired': true,
      'quickScanStatusEnabled': true,
      'quickScanStatusRequired': true,
      'quickScanCategoriesEnabled': true,
      'quickScanCategoriesRequired': true,
      'quickScanTagsEnabled': false,
      'quickScanTagsRequired': false,
    });

    await loginAsAdmin(tester);
    final target = await keepOnlyGroup(tester, g.id);

    await openQuickScanImageForm(tester);
    await selectDropdown(tester, 'groupId', target.name);
    await selectDropdown(tester, 'paidByUserId', adminDisplayName(tester));
    await selectDropdown(tester, 'status', 'Open');

    // Open the category picker (the path that crashed pre-fix via
    // Navigator.of(null)), filter to the seeded category, select it, confirm.
    final placeholder = find.text('No Categories selected');
    await tester.ensureVisible(placeholder);
    await tester.pump(const Duration(milliseconds: 200));
    await tester.tap(placeholder);
    await pumpUntilFound(tester, find.text('Select Categories'));

    // Filter to the (unique) seeded category so its chip is the only one laid
    // out -- the admin's catalog can hold many categories off-screen.
    final filter = find.byWidgetPredicate(
        (w) => w is FormBuilderTextField && w.name == 'filter');
    await tester.enterText(filter, categoryName);
    await tester.pump(const Duration(milliseconds: 300));

    final chip = find.widgetWithText(ChoiceChip, categoryName);
    await pumpUntilFound(tester, chip.hitTestable());
    await tester.tap(chip);
    await tester.pump(const Duration(milliseconds: 200));

    final selectBtn = find.widgetWithText(BottomSubmitButton, 'Select');
    await pumpUntilFound(tester, selectBtn.hitTestable());
    await tester.tap(selectBtn.hitTestable());

    // Drain the picker's exit animation so only the form's submit remains.
    for (int i = 0; i < 6; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }

    await expectQuickScanQueued(tester);
  });
}
