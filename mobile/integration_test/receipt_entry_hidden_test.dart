// A user who can neither quick-scan nor create gets no receipt-entry affordance
// at all -- not a disabled one, and not one that refuses on tap.
//
// This is the case the design did not cover: it assumes manual entry is always
// available, but group.receipts.create is a real, separately-enforced permission
// that a role can lack.

import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_entry.dart';

import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/receipt_test_helpers.dart';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('a Legacy Viewer sees no scan slot and no overflow entries',
      (tester) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    final fixture = await provisionPermUser(roleName: 'Legacy Viewer');
    await loginAs(
      tester,
      username: fixture.username,
      password: fixture.password,
    );
    await enterGroup(tester, fixture.groupName!);
    await tester.tap(find.text('Receipts').hitTestable());
    await tester.pumpAndSettle(const Duration(milliseconds: 500));

    expect(scanNavSlot(), findsNothing,
        reason: 'a Viewer holds neither create nor quick-scan in this group');
    expect(find.text('Scan').hitTestable(), findsNothing);

    // The other destinations are unaffected -- the nav resolves by destination
    // id, so removing a middle slot must not disturb them.
    expect(find.text('Dashboards'), findsWidgets);
    expect(find.text('Receipts'), findsWidgets);

    expect(find.byKey(const ValueKey('receipt-entry-overflow-menu')),
        findsNothing,
        reason: 'an overflow that could only ever say "no" is worse than none');
    expect(find.text(addManualReceiptLabel), findsNothing);
  });
}
