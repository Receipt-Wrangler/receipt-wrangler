// Quick Scan per-image PREFILL from user preferences, and how it interacts with
// the group's config.
//
// Each added image is seeded from userPreferences.quickScanDefault{GroupId,
// PaidById,Status} (quick_scan.dart `_getInitialQuickScanValues`, delivered via
// AppData at login). The correctness under test: a preset that the selected
// group HIDES must "fall off" -- the field isn't shown, and _submitQuickScan
// sends the sentinel for it (quick_scan.dart:173-179), so the preset never
// reaches the receipt. A preset for a SHOWN field is honored.
//
// This is the real prefs -> AppData -> form path, so the config is PERSISTED (not
// client-injected): PUT /group/{id}/groupReceiptSettings + PUT /userPreferences,
// both restored on teardown. The submit is only asserted as "queued" (the hidden
// paid-by is backfilled server-side from the UPLOADER default; the queued receipt
// has no deterministic id, so the payload isn't observable here -- the widget
// tests assert the field falls off deterministically).

import 'dart:io' show Platform;

import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:flutter/widgets.dart';
import 'package:integration_test/integration_test.dart';
import 'package:openapi/openapi.dart' as api;

import 'helpers/api.dart';
import 'helpers/document_scanner_mock.dart';
import 'helpers/env.dart';
import 'helpers/feature_flags.dart';
import 'helpers/login.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/quick_scan_actions.dart';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets(
      'a preset paid-by falls off when the group hides it; a preset status is kept',
      (tester) async {
    await enableAiPoweredReceiptsForTest();
    await installDocumentScannerMock();
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    final jwt = await apiLogin();
    final g = await firstNonAllGroup(jwt);
    final adminId = await userIdByUsername(E2eEnv.adminUsername, jwt: jwt);

    // Group G HIDES paid-by (needs an UPLOADER default so the backend backfills)
    // but SHOWS status (required); categories/tags hidden.
    await setGroupQuickScanConfig(groupId: g.id, jwt: jwt, overrides: const {
      'quickScanPaidByEnabled': false,
      'quickScanDefaultPaidByType': 'UPLOADER',
      'quickScanStatusEnabled': true,
      'quickScanStatusRequired': true,
      'quickScanCategoriesEnabled': false,
      'quickScanCategoriesRequired': false,
      'quickScanTagsEnabled': false,
      'quickScanTagsRequired': false,
    });

    // The user has a preset paid-by (the admin) + status (OPEN) + default group G.
    await setUserQuickScanPrefs(
        jwt: jwt, groupId: g.id, paidById: adminId, status: 'OPEN');

    await loginAsAdmin(tester);

    // Adding an image prefills group=G / paid-by=admin / status=OPEN; the group is
    // pre-selected from the prefill, so the form reads G's config immediately.
    await openQuickScanImageForm(tester);

    // Paid-by is hidden by G -> the preset admin fell off the form.
    expect(quickScanDropdown('paidByUserId'), findsNothing,
        reason: 'preset paid-by fell off (group hides paid-by)');

    // Status is shown -> the preset OPEN is honored.
    expect(quickScanDropdown('status'), findsOneWidget, reason: 'status shown');
    final status =
        tester.widget(quickScanDropdown('status')) as FormBuilderDropdown;
    expect(status.initialValue, api.ReceiptStatus.OPEN,
        reason: 'preset status prefilled');

    // Submit succeeds: status satisfied by the preset, paid-by backfilled from
    // the group's UPLOADER default.
    await expectQuickScanQueued(tester);
  });
}
