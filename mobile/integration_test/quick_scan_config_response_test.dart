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
}
