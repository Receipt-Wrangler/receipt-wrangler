import 'package:flutter_test/flutter_test.dart';

import 'pump.dart';

/// From the GroupSelect screen, taps the [groupName] card to enter the group
/// context shell, returning once the group bottom nav ("Receipts" tab) is on
/// screen. A freshly-provisioned permission user belongs to exactly one group,
/// so the name is a unique match.
///
/// Both taps wait for *hittability*, not bare existence: on iOS the Cupertino
/// page transitions (~400ms) keep the destination sliding after its widgets
/// mount, and a tap computed mid-slide misses (observed deterministically on
/// the iOS simulator as "Offset(423.9, 796.0) ... would not hit test" against
/// the Receipts tab). See mobile/CLAUDE.md "Three tap-flake patterns".
Future<void> enterGroup(WidgetTester tester, String groupName) async {
  await pumpUntilFound(tester, find.text(groupName).hitTestable());
  await _drain(tester);
  await tester.tap(find.text(groupName).hitTestable());
  // The group bottom nav ("Dashboards"/"Add"/"Receipts"/"Search") only mounts
  // inside the group context shell -- a stronger signal than a URL check.
  await pumpUntilFound(tester, find.text('Receipts').hitTestable());
}

/// Enters [groupName] and opens its Receipts tab, returning once [receiptName]
/// is on screen.
Future<void> openGroupReceipts(
  WidgetTester tester,
  String groupName,
  String receiptName,
) async {
  await enterGroup(tester, groupName);
  await _drain(tester);
  await tester.tap(find.text('Receipts').hitTestable());
  await pumpUntilFound(tester, find.text(receiptName));
}

/// Drains the tail of a page/sheet transition so a follow-up tap computes its
/// center from settled geometry.
Future<void> _drain(WidgetTester tester) async {
  for (int i = 0; i < 5; i++) {
    await tester.pump(const Duration(milliseconds: 100));
  }
}
