import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_edit_popup_menu.dart';

import 'api.dart';
import 'form_actions.dart';
import 'pump.dart';
import 'users.dart';

/// Reads the current GoRouter URL by grabbing a context from inside the
/// routed tree (`MaterialApp` itself sits above the GoRouter scope, so
/// its element fails the inherited-widget lookup with "No GoRouter
/// found in context").
String currentUrl(WidgetTester tester) {
  final scaffold = find.byType(Scaffold).evaluate();
  if (scaffold.isEmpty) return '';
  return GoRouter.of(scaffold.first)
      .routerDelegate
      .currentConfiguration
      .uri
      .toString();
}

/// Pumps until the current GoRouter URL matches [pattern], or [timeout]
/// elapses.
///
/// PREFER `pumpUntilFound(tester, find.byType(SomeDestinationWidget))`
/// for asserting that a route has actually mounted -- a widget is the
/// stronger signal that the shell finished building (per `mobile/CLAUDE.md`
/// "Assert navigation by widget presence, not URL"). Use this URL helper
/// only when you also need to *extract* something from the URL after
/// arrival (e.g. the receipt id from `/receipts/<id>/view`); pair it with
/// a widget-presence wait first so the URL read isn't a race.
Future<String> pumpUntilUrl(
  WidgetTester tester,
  RegExp pattern, {
  Duration timeout = const Duration(seconds: 15),
}) async {
  final deadline = DateTime.now().add(timeout);
  while (DateTime.now().isBefore(deadline)) {
    await tester.pump(const Duration(milliseconds: 100));
    final url = currentUrl(tester);
    if (pattern.hasMatch(url)) return url;
  }
  throw StateError(
    'Timed out after ${timeout.inSeconds}s waiting for URL matching '
    '$pattern. Last seen: ${currentUrl(tester)}',
  );
}

/// Extracts the receipt id from a `/receipts/<id>/view` URL produced by
/// the production save handler.
int receiptIdFromUrl(String url) {
  final match = RegExp(r'/receipts/(\d+)/view').firstMatch(url);
  if (match == null) {
    throw StateError('Expected /receipts/<id>/view URL, got: $url');
  }
  return int.parse(match.group(1)!);
}

/// Best-effort cleanup: log in via the API and DELETE the receipt at
/// end-of-test. Wrapped via [addTearDown] so it runs even if the test
/// body throws after the receipt was created.
void scheduleReceiptCleanup(int receiptId) {
  addTearDown(() async {
    final jwt = await apiLogin();
    await deleteReceipt(receiptId, jwt: jwt);
  });
}

/// Drives the receipt-add UI from `/groups`: opens the bottom-nav Add
/// menu, fills the required fields, taps Submit, waits for navigation
/// to `/receipts/<id>/view`. Returns the new receipt's id.
///
/// Used by tests that need a baseline receipt to operate on (Flow #4
/// edits it, Flow #5 chains two of these). The same field-fill sequence
/// as Flow #1's smoke happy path -- if Flow #1 is green, this is too.
///
/// Pre-conditions: caller is logged in and currently on `/groups`.
Future<int> addManualReceiptViaUI(
  WidgetTester tester,
  String name, {
  String amount = '12.34',
}) async {
  await tester.tap(find.text('Add'));
  await pumpUntilFound(tester, find.text('Add Manual Receipt'));
  await tester.tap(find.text('Add Manual Receipt'));
  await pumpUntilFound(tester, find.text('Name'));

  await tester.enterText(formField('name'), name);
  await tester.enterText(formField('amount'), amount);
  await selectDropdown(tester, 'groupId', 'My Receipts');
  await selectDropdown(tester, 'paidByUserId', adminDisplayName(tester));

  // Drain the dropdown overlay teardown -- the popup-route's overlay
  // entry can otherwise leave the Scaffold's bottom-sheet area in an
  // Offstage state and the BottomSubmitButton tap silently misses.
  await tester.pumpAndSettle(const Duration(seconds: 3));

  await tester.tap(find.byType(BottomSubmitButton));
  // Assert /view shell has mounted via the ReceiptEditPopupMenu, which
  // only renders on /view (gated on canEditReceipt -- see
  // mobile/lib/shared/widgets/receipt_edit_popup_menu.dart:31). Then
  // extract the id from the URL.
  await pumpUntilFound(tester, find.byType(ReceiptEditPopupMenu));
  return receiptIdFromUrl(currentUrl(tester));
}
