// Flow #5 -- two consecutive add-receipt flows in one session.
//
// Login -> add receipt 1 (UI) -> back to /groups -> add receipt 2 (UI)
// -> verify via API that receipt 2's name is the second name we typed
// (no leakage from receipt 1's state via ReceiptModel.resetModel()
// at the /receipts/add redirect in mobile/lib/main.dart:124).
//
// Own file for the GoRouter-persistence reason (see
// receipt_add_share_test.dart's header comment).

import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';

import 'helpers/api.dart';
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

  testWidgets('two consecutive adds produce two distinct receipts',
      (tester) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    await loginAsAdmin(tester);

    final ts = DateTime.now().millisecondsSinceEpoch;
    final name1 = 'e2e-consec-1-$ts';
    final name2 = 'e2e-consec-2-$ts';

    // Add receipt 1.
    final id1 = await addManualReceiptViaUI(tester, name1, amount: '11.11');
    scheduleReceiptCleanup(id1);

    // Wait for the view screen's app bar (back arrow) to mount fully
    // -- it depends on receiptModel having the loaded receipt's
    // groupId for the back URL.
    await pumpUntilFound(tester, find.byIcon(Icons.arrow_back));
    // Two arrow_back icons exist in the view tree (the app bar's
    // leading button + something in the deeper widget chain). Tap
    // the first -- in widget order it's the AppBar's leading slot.
    final backButton = find.byIcon(Icons.arrow_back).first;
    // Tap the back arrow. ReceiptAppBar's onLeadingArrowPressed calls
    // receiptModel.resetModel() BEFORE navigating, so the form's
    // FormBuilderDateTimePicker doesn't try to rebuild with stale
    // String date values during the navigation transition (which it
    // does if we GoRouter.go() directly -- a known production bug:
    // "type 'String' is not a subtype of type 'DateTime?' of 'value'"
    // from FormBuilderFieldState.setValue when the field's stored
    // value carries over from receipt 1's save).
    await tester.tap(backButton);
    // Lands on /groups/<groupId>/receipts which uses GroupBottomNav.
    // GroupBottomNav also has an "Add" destination -- waiting on that
    // text is both the destination-mounted check and the target for
    // the next addManualReceiptViaUI call.
    await pumpUntilFound(tester, find.text('Add'));

    // Add receipt 2.
    final id2 = await addManualReceiptViaUI(tester, name2, amount: '22.22');
    scheduleReceiptCleanup(id2);

    expect(id1 == id2, isFalse,
        reason: 'Consecutive adds should produce distinct ids');

    // Verify via API that each receipt has the correct name -- this
    // is the actual resetModel() regression check. If state leaked
    // between the two adds, receipt 2 would carry receipt 1's name
    // (the name field's initial value) rather than what the test typed.
    final jwt = await apiLogin();
    final receipt2 = await getReceipt(id2, jwt: jwt);
    expect(receipt2['name'], name2,
        reason: 'Receipt 2 should carry name2 ($name2), not name1 '
            '(${receipt2['name']}). State leak from receipt 1.');

    // Sanity check on receipt 1 too.
    final receipt1 = await getReceipt(id1, jwt: jwt);
    expect(receipt1['name'], name1);
  });
}
