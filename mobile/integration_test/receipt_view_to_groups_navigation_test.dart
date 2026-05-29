// Reproduces and verifies a fix for a production bug in the receipt
// form's date field handling.
//
// SYMPTOM: navigating away from /receipts/<id>/view via
// GoRouter.go('/groups') (i.e. NOT through the app bar back arrow,
// which calls receiptModel.resetModel() first) throws an uncaught
// widgets-library exception:
//
//   type 'String' is not a subtype of type 'DateTime?' of 'value'
//   (FormBuilderFieldState.setValue, called from
//    _FormBuilderDateTimePickerState.initState ->
//    FormBuilderState.registerField)
//
// CAUSE (suspected, confirmed by this test): the receipt form names
// both the view-mode date field (FormBuilderTextField, value: String)
// and the edit/add-mode date field (FormBuilderDateTimePicker, value:
// DateTime) "date". When the form is unmounted/remounted across these
// modes, FormBuilder's internal value map carries the previous
// String value, and the DateTimePicker's setValue rejects it.
//
// REPRO: this test asserts that no widget exception fires during a
// post-add navigation to /groups. Without the fix, the typer error
// fires during the rebuild and tester.takeException() returns it.

import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/groups/screens/group_select.dart';

import 'helpers/login.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/receipt_test_helpers.dart';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets(
      'navigating /view -> /groups via GoRouter does not throw '
      'FormBuilderDateTimePicker String/DateTime type error',
      (tester) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    await loginAsAdmin(tester);

    // Add a receipt so we land on /receipts/<id>/view.
    final receiptName =
        'e2e-nav-${DateTime.now().millisecondsSinceEpoch}';
    final receiptId = await addManualReceiptViaUI(tester, receiptName);
    scheduleReceiptCleanup(receiptId);

    // Now on /view. Wait briefly so the FutureBuilder has a chance
    // to load the receipt and call setReceipt -- this is when the
    // bug would surface in the field initStates.
    await tester.pump(const Duration(milliseconds: 500));

    // Navigate AWAY without going through the app bar back arrow
    // (which would call receiptModel.resetModel() and side-step the
    // bug). Programmatic GoRouter navigation forces the original
    // unmount/remount race.
    final scaffolds = find.byType(Scaffold).evaluate();
    expect(scaffolds, isNotEmpty);
    GoRouter.of(scaffolds.first).go('/groups');
    // /groups mounts GroupSelect; that's the destination-mounted marker
    // per the project convention. The follow-on pump(500ms) keeps the
    // breathing room the original flow had after the route landed.
    await pumpUntilFound(tester, find.byType(GroupSelect));
    await tester.pump(const Duration(milliseconds: 500));

    // Now drive a SECOND add through the bottom nav (matching the
    // original failing flow in receipt_add_consecutive_test.dart).
    // The form rebuild triggered by the second /receipts/add mount
    // is when the type error originally fires.
    final receiptName2 =
        'e2e-nav-2-${DateTime.now().millisecondsSinceEpoch}';
    final id2 = await addManualReceiptViaUI(tester, receiptName2);
    scheduleReceiptCleanup(id2);

    final ex = tester.takeException();
    expect(
      ex,
      isNull,
      reason: 'A second consecutive add via GoRouter.go(/groups) '
          'should not throw. Got: $ex',
    );
  });
}
