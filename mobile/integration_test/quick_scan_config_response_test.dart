// Verifies the Quick Scan per-image FORM responds to the selected group's
// quick-scan configuration (GroupReceiptSettings.quickScan*) -- no scan is sent.
//
// Runs headless on Linux desktop. Two facts shape the approach:
//   * Quick Scan is gated on featureConfig.aiPoweredReceipts (off locally), so
//     we flip it on server-side before login (enableAiPoweredReceiptsForTest).
//   * The sheet's gallery-upload icon throws "Unsupported platform" on desktop,
//     so we feed an image through the document-scanner icon via a mocked
//     `cunning_document_scanner` channel (installDocumentScannerMock) -- the
//     camera permission it requests is already granted by installLinuxDesktopMocks.
//
// The group config is injected the same way quick_scan_disabled_test injects the
// feature flag: a Provider mutation on the live GroupModel. This is deterministic
// and independent of whether the local API persists the new fields.

import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/category_select_field.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/tag_select_field.dart';

import 'helpers/document_scanner_mock.dart';
import 'helpers/feature_flags.dart';
import 'helpers/form_actions.dart';
import 'helpers/login.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

Finder _dropdown(String name) =>
    find.byWidgetPredicate((w) => w is FormBuilderDropdown && w.name == name);

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('quick scan form shows/hides + requires fields per group config',
      (tester) async {
    await enableAiPoweredReceiptsForTest();
    await installDocumentScannerMock();
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    await loginAsAdmin(tester);

    // Inject a distinctive quick-scan config onto a group the admin owns:
    // Paid By hidden, Status shown+optional, Categories shown+required, Tags hidden.
    final ctx = tester.element(find.byType(Scaffold).first);
    final groupModel = Provider.of<GroupModel>(ctx, listen: false);
    final target = groupModel.groups.firstWhere((g) => !g.isAllGroup);
    final configured = target.rebuild((b) => b
      ..groupReceiptSettings.quickScanPaidByEnabled = false
      ..groupReceiptSettings.quickScanPaidByRequired = false
      ..groupReceiptSettings.quickScanStatusEnabled = true
      ..groupReceiptSettings.quickScanStatusRequired = false
      ..groupReceiptSettings.quickScanCategoriesEnabled = true
      ..groupReceiptSettings.quickScanCategoriesRequired = true
      ..groupReceiptSettings.quickScanTagsEnabled = false
      ..groupReceiptSettings.quickScanTagsRequired = false);
    groupModel.setGroups(groupModel.groups
        .map((g) => g.id == target.id ? configured : g)
        .toList());
    await tester.pump();

    // Open the bottom-nav Add menu and pick "Quick Scan" (the AI flag is on, so
    // the menu item renders). Wait for hittability + drain the slide-in.
    await tester.tap(find.text('Add').hitTestable());
    await pumpUntilFound(tester, find.text('Quick Scan').hitTestable());
    for (int i = 0; i < 5; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }
    await tester.tap(find.text('Quick Scan').hitTestable());

    // Add an image via the document-scanner icon (mock returns one PNG); the
    // per-image form mounts once an image is present.
    await pumpUntilFound(tester, find.byIcon(Icons.add_a_photo));
    await tester.tap(find.byIcon(Icons.add_a_photo));
    await pumpUntilFound(tester, find.text('Group'));

    // Select the configured group so the form reads its config.
    await selectDropdown(tester, 'groupId', target.name);

    // Fields now reflect the injected config.
    expect(_dropdown('paidByUserId'), findsNothing, reason: 'paid-by hidden');
    expect(_dropdown('status'), findsOneWidget, reason: 'status shown');
    expect(find.byType(CategorySelectField), findsOneWidget,
        reason: 'categories shown');
    expect(find.byType(TagSelectField), findsNothing, reason: 'tags hidden');

    // Categories is required and empty -> submitting surfaces the fix-error
    // snackbar and queues nothing (the submit returns before the API call).
    await pumpUntilFound(tester, find.byType(BottomSubmitButton).hitTestable());
    await tester.tap(find.byType(BottomSubmitButton));
    await pumpUntilFound(
      tester,
      find.textContaining('Please fix error on quick scan'),
      timeout: const Duration(seconds: 10),
    );
  });

  testWidgets('quick scan form re-reads config when the group changes',
      (tester) async {
    await enableAiPoweredReceiptsForTest();
    await installDocumentScannerMock();
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    await loginAsAdmin(tester);

    // Two groups with opposite configs. Group 1 is the admin's real group
    // rebuilt with config A; group 2 is a clone with a fresh id + name and the
    // inverse config B. Injected via the live GroupModel (deterministic, no
    // reliance on the API persisting the new fields).
    final ctx = tester.element(find.byType(Scaffold).first);
    final groupModel = Provider.of<GroupModel>(ctx, listen: false);
    final real = groupModel.groups.firstWhere((g) => !g.isAllGroup);
    final group2Id =
        groupModel.groups.fold<int>(0, (m, g) => g.id > m ? g.id : m) + 1;

    // Config A: paid-by shown+optional, status hidden, categories shown+required,
    // tags hidden.
    final group1 = real.rebuild((b) => b
      ..groupReceiptSettings.quickScanPaidByEnabled = true
      ..groupReceiptSettings.quickScanPaidByRequired = false
      ..groupReceiptSettings.quickScanStatusEnabled = false
      ..groupReceiptSettings.quickScanStatusRequired = false
      ..groupReceiptSettings.quickScanCategoriesEnabled = true
      ..groupReceiptSettings.quickScanCategoriesRequired = true
      ..groupReceiptSettings.quickScanTagsEnabled = false
      ..groupReceiptSettings.quickScanTagsRequired = false);

    // Config B: paid-by hidden, status shown+required, categories hidden,
    // tags shown+optional. Distinct id/name so it's a separate dropdown entry.
    final group2 = real.rebuild((b) => b
      ..id = group2Id
      ..name = 'E2E Switch Group B'
      ..groupReceiptSettings.id = group2Id
      ..groupReceiptSettings.groupId = group2Id
      ..groupReceiptSettings.quickScanPaidByEnabled = false
      ..groupReceiptSettings.quickScanPaidByRequired = false
      ..groupReceiptSettings.quickScanStatusEnabled = true
      ..groupReceiptSettings.quickScanStatusRequired = true
      ..groupReceiptSettings.quickScanCategoriesEnabled = false
      ..groupReceiptSettings.quickScanCategoriesRequired = false
      ..groupReceiptSettings.quickScanTagsEnabled = true
      ..groupReceiptSettings.quickScanTagsRequired = false);

    // Keep only the all-group placeholder(s) plus our two configured groups so
    // the group dropdown is short and both entries are on-screen. The seeded
    // admin can accumulate many groups; an appended entry would land below the
    // dropdown menu's viewport, where `find.text` can't reach it.
    final allGroups = groupModel.groups.where((g) => g.isAllGroup).toList();
    groupModel.setGroups([...allGroups, group1, group2]);
    await tester.pump();

    // Open Add → Quick Scan and add an image so the per-image form mounts.
    await tester.tap(find.text('Add').hitTestable());
    await pumpUntilFound(tester, find.text('Quick Scan').hitTestable());
    for (int i = 0; i < 5; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }
    await tester.tap(find.text('Quick Scan').hitTestable());

    await pumpUntilFound(tester, find.byIcon(Icons.add_a_photo));
    await tester.tap(find.byIcon(Icons.add_a_photo));
    await pumpUntilFound(tester, find.text('Group'));

    // Select group 1 → config-A field set.
    await selectDropdown(tester, 'groupId', real.name);
    expect(_dropdown('paidByUserId'), findsOneWidget, reason: 'A: paid-by shown');
    expect(_dropdown('status'), findsNothing, reason: 'A: status hidden');
    expect(find.byType(CategorySelectField), findsOneWidget,
        reason: 'A: categories shown');
    expect(find.byType(TagSelectField), findsNothing, reason: 'A: tags hidden');

    // Switch to group 2 → the field set flips to config B.
    await selectDropdown(tester, 'groupId', group2.name);
    expect(_dropdown('paidByUserId'), findsNothing, reason: 'B: paid-by hidden');
    expect(_dropdown('status'), findsOneWidget, reason: 'B: status shown');
    expect(find.byType(CategorySelectField), findsNothing,
        reason: 'B: categories hidden');
    expect(find.byType(TagSelectField), findsOneWidget, reason: 'B: tags shown');

    // Group 2's status is shown+required and empty → submit is blocked with the
    // fix-error snackbar (a required field other than the categories case above).
    await pumpUntilFound(tester, find.byType(BottomSubmitButton).hitTestable());
    await tester.tap(find.byType(BottomSubmitButton));
    await pumpUntilFound(
      tester,
      find.textContaining('Please fix error on quick scan'),
      timeout: const Duration(seconds: 10),
    );
  });
}
