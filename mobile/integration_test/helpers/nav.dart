import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';

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

/// Opens [receiptName] in [groupName] and moves it to the **edit** form via the
/// view's popup menu. Reaching edit state requires `group.receipts.update`.
/// Returns once the edit form's submit button is mounted -- the marker unique to
/// edit/add, since `find.text('Name')` also matches the view form.
Future<void> openReceiptEditForm(
  WidgetTester tester,
  String groupName,
  String receiptName,
) async {
  await openGroupReceipts(tester, groupName, receiptName);

  await tester.tap(find.text(receiptName));

  final menuButton = find.byType(PopupMenuButton<dynamic>);
  await pumpUntilFound(tester, menuButton);
  await tester.tap(menuButton);
  await pumpUntilFound(tester, find.text('Edit').hitTestable());
  await _drain(tester);
  await tester.tap(find.text('Edit').hitTestable());
  await pumpUntilFound(tester, find.byType(BottomSubmitButton));
}

/// [openReceiptEditForm], then the form's "View Comments" action -- the
/// edit-state comment screen, where adds and deletes hit the API immediately.
Future<void> openReceiptCommentsInEditMode(
  WidgetTester tester,
  String groupName,
  String receiptName,
) async {
  await openReceiptEditForm(tester, groupName, receiptName);

  final commentsButton = find.byTooltip('View Comments');
  await pumpUntilFound(tester, commentsButton);
  await tester.tap(commentsButton);
  await pumpUntilFound(tester, find.text('Receipt Comments'));
}

/// Drains the tail of a page/sheet transition so a follow-up tap computes its
/// center from settled geometry.
Future<void> _drain(WidgetTester tester) async {
  for (int i = 0; i < 5; i++) {
    await tester.pump(const Duration(milliseconds: 100));
  }
}
