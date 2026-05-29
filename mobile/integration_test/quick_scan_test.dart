// Flow #8 -- Quick Scan happy path.
//
// Quick Scan is the AI-assisted bulk-receipt entry flow. The user
// taps "Quick Scan" from the bottom-nav Add menu, picks one or more
// images (camera or gallery), fills a per-image form (group, paid
// by, status), and submits. The backend queues each image as an
// async OCR/AI extraction job that materializes into a receipt.
//
// PRECONDITION: the demo and local backends both have
// `featureConfig.aiPoweredReceipts: true`. If false, the in-app
// flow shows an error snackbar instead of the bottom sheet (see
// mobile/lib/shared/functions/quick_scan.dart:227-231).
//
// Skipped on Linux: scan.dart's gallery path throws "Unsupported
// platform" for Linux/macOS/Windows desktop. Runs on Android
// emulator + iOS simulator in CI.

import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';

import 'helpers/file_selector_mock.dart';
import 'helpers/form_actions.dart';
import 'helpers/login.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/users.dart';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('quick scan from gallery: pick image, fill form, submit',
      // Same Linux skip as Flow #2 / Flow B -- gallery picker only
      // supports Android/iOS in scan.dart.
      skip: Platform.isLinux,
      (tester) async {
    await installFileSelectorMock();
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    await loginAsAdmin(tester);

    // Open the bottom-nav Add menu and pick "Quick Scan".
    // The featureConfig.aiPoweredReceipts gate is enforced inside
    // showQuickScanBottomSheet (a snackbar fires if it's false), but
    // the menu item is unconditionally present in showAddMenu (see
    // mobile/lib/shared/functions/show_add_menu.dart:50-54). So the
    // gate's verification is implicit: if the bottom sheet opens,
    // the flag is true.
    await tester.tap(find.text('Add'));
    await pumpUntilFound(tester, find.text('Quick Scan'));
    await tester.tap(find.text('Quick Scan'));

    // Wait for the bottom sheet to mount. The title "Quick Scan"
    // shows in the sheet header; the gallery upload action shows as
    // an Icons.upload_file_rounded IconButton (quick_scan.dart:67).
    await pumpUntilFound(tester, find.byIcon(Icons.upload_file_rounded));

    // Trigger the gallery picker -- the mock returns one 1x1 PNG.
    await tester.tap(find.byIcon(Icons.upload_file_rounded));
    // Mock resolves immediately; the imageSubject emits and the
    // QuickScanForm card mounts with three dropdowns (groupId,
    // paidByUserId, status). Wait for the form to appear.
    await pumpUntilFound(tester, find.text('Group'));

    // Fill the per-image form. e2e-admin's quickScan user prefs are
    // null, so all three fields need to be set explicitly.
    await selectDropdown(tester, 'groupId', 'My Receipts');
    await selectDropdown(tester, 'paidByUserId', adminDisplayName(tester));
    await selectDropdown(tester, 'status', 'Open');

    // Drain the dropdown overlay teardown before tapping Submit.
    await tester.pumpAndSettle(const Duration(seconds: 3));

    await tester.tap(find.byType(BottomSubmitButton));

    // Success message text is from quick_scan.dart:179-182:
    // "Successfully queued $imageWord for processing!" where
    // imageWord is "image" for 1 file.
    await pumpUntilFound(
      tester,
      find.textContaining('Successfully queued'),
      timeout: const Duration(seconds: 15),
    );

    // No cleanup needed: Quick Scan queues an async backend job;
    // the resulting receipt is created by the worker and won't be
    // assigned a deterministic id we can DELETE inline. The unique
    // PNG bytes used by the mock are unlikely to produce a
    // recognizable extracted name -- if these accumulate on demo
    // they're identifiable as low-info OCR results from the test
    // window.
  });
}
