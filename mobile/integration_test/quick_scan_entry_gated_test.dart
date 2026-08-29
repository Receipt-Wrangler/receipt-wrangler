// What the scan slot does when Quick Scan cannot run.
//
// The tap falls through to the manual receipt form, carrying a banner that names
// the reason -- and the two reasons need different people to fix them, so they
// must not be conflated. The permission half is provisioned as a real role
// (Legacy Editor minus quick-scan) so the gating is proved against the backend's
// own role definitions rather than a seeded model.

import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_entry.dart';

import 'helpers/feature_flags.dart';
import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/receipt_test_helpers.dart';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  const bannerKey = ValueKey('quick-scan-unavailable-banner');
  const dismissKey = ValueKey('quick-scan-unavailable-banner-dismiss');

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('without the quick-scan permission, the tap explains itself',
      (tester) async {
    await enableAiPoweredReceiptsForTest();
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    // Legacy Editor MINUS quick-scan: still holds receipts.create, so the slot
    // stays and manual entry is the right fallback.
    final fixture = await provisionGroupMemberWithoutPermission(
      'group.receipts.quick-scan',
      baselineRole: 'Legacy Editor',
    );
    await loginAs(
      tester,
      username: fixture.username,
      password: fixture.password,
    );
    await enterGroup(tester, fixture.groupName!);

    expect(find.text('Add').hitTestable(), findsWidgets,
        reason: 'the slot must not advertise a scanner this user cannot use');

    await tester.tap(scanNavSlot());
    await pumpUntilFound(tester, find.text('Name'));

    await pumpUntilFound(tester, find.byKey(bannerKey));
    expect(
      find.text(quickScanNoPermissionMessageForGroup(fixture.groupName!)),
      findsOneWidget,
      reason: 'the permission is per-group, so the banner names the group',
    );

    await tester.tap(find.byKey(dismissKey));
    await tester.pumpAndSettle();
    expect(find.byKey(bannerKey), findsNothing);
  });

  testWidgets('with the ai flag off, the banner points at the server config',
      (tester) async {
    // The local backend ships with no receipt-processing settings, so
    // aiPoweredReceipts is off by default -- no fixture needed.
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    final fixture = await provisionPermUser(roleName: 'Legacy Editor');
    await loginAs(
      tester,
      username: fixture.username,
      password: fixture.password,
    );
    await enterGroup(tester, fixture.groupName!);

    await tester.tap(scanNavSlot());
    await pumpUntilFound(tester, find.text('Name'));

    await pumpUntilFound(tester, find.byKey(bannerKey));
    expect(find.text(quickScanAiDisabledMessage), findsOneWidget,
        reason: 'this user holds quick-scan; the install is what is missing, '
            'so the banner must send them to an administrator');
  });

  testWidgets('a deliberate Add Manual Receipt shows no banner',
      (tester) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    final fixture = await provisionPermUser(roleName: 'Legacy Editor');
    await loginAs(
      tester,
      username: fixture.username,
      password: fixture.password,
    );
    await enterGroup(tester, fixture.groupName!);

    await openManualReceiptForm(tester);

    expect(find.byKey(bannerKey), findsNothing,
        reason: 'the user chose the form; nothing about Quick Scan surprised '
            'them');
  });
}
