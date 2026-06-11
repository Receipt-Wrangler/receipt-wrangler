import 'package:flutter_test/flutter_test.dart';

import 'pump.dart';

/// From the GroupSelect screen, taps the [groupName] card to enter the group
/// context shell, returning once the group bottom nav ("Receipts" tab) is on
/// screen. A freshly-provisioned permission user belongs to exactly one group,
/// so the name is a unique match.
Future<void> enterGroup(WidgetTester tester, String groupName) async {
  await pumpUntilFound(tester, find.text(groupName));
  await tester.tap(find.text(groupName));
  // The group bottom nav ("Dashboards"/"Add"/"Receipts"/"Search") only mounts
  // inside the group context shell -- a stronger signal than a URL check.
  await pumpUntilFound(tester, find.text('Receipts'));
}

/// Enters [groupName] and opens its Receipts tab, returning once [receiptName]
/// is on screen.
Future<void> openGroupReceipts(
  WidgetTester tester,
  String groupName,
  String receiptName,
) async {
  await enterGroup(tester, groupName);
  await tester.tap(find.text('Receipts'));
  await pumpUntilFound(tester, find.text(receiptName));
}
