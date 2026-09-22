import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/paged_data_list.dart';

import 'helpers/api.dart';
import 'helpers/form_actions.dart';
import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/receipt_filter_actions.dart';

/// End-to-end coverage for the **receipt summary bar**, against the real Go API.
///
/// The widget suite pins the bar exhaustively -- the labelling, the muted zero row, the
/// locked figure columns, the refresh contract -- but all of it against a mocked
/// `ReceiptApi` and a `GroupModel` seeded by hand. What none of it can prove is the two
/// things most likely to break:
///
///  * that the **encoded filter reaches the summary endpoint** and narrows the figures
///    the same way it narrows the rows above them, and
///  * that **`receiptSummaryPosition` survives** model -> command -> DB -> response ->
///    render, which is the whole path the new setting travels.
///
/// **Seed set** (one fixture group, three receipts):
///
/// | | name | amount | status |
/// |---|---|---|---|
/// | R1 | `…-alpha`   | 10.00 | OPEN     |
/// | R2 | `…-bravo`   | 25.50 | RESOLVED |
/// | R3 | `…-charlie` | 4.50  | OPEN     |
///
/// Overall is 3 receipts / 40.00; OPEN is 2 / 14.50. The amounts are chosen so the
/// overall total cannot be produced by summing any two of them -- a summary that
/// silently dropped a row would still have to report a different number.
///
/// One spec, several steps, deliberately: the e2e suite's cost is dominated by
/// provisioning and login, so six `testWidgets` would be six logins.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  Finder summaryBar() => find.byKey(const ValueKey('receipt-summary'));
  Finder overallRow() => find.byKey(const ValueKey('receipt-summary-row-overall'));
  Finder openRow() => find.byKey(const ValueKey('receipt-summary-row-OPEN'));

  /// The rendered text of a row's label cell, e.g. "All Receipts\n(3 receipts)".
  String rowText(WidgetTester tester, Finder row) =>
      tester.widget<Text>(find.descendant(of: row, matching: find.byType(Text))).data!;

  testWidgets('the summary totals the whole filter, and honours its position',
      (tester) async {
    await tester.binding.setSurfaceSize(const Size(430, 932));
    addTearDown(() => tester.binding.setSurfaceSize(null));

    final jwt = await apiLogin();
    final stamp = DateTime.now().microsecondsSinceEpoch.toString();
    final fixture = await provisionPermUser(roleName: 'Legacy Editor');
    final groupId = fixture.groupId!;

    final alpha = 'e2e-sum-$stamp-alpha';
    final bravo = 'e2e-sum-$stamp-bravo';
    final charlie = 'e2e-sum-$stamp-charlie';

    // No per-receipt teardown: DELETE /group/{id} cascades its receipts, and the
    // fixture group's delete is already registered.
    await createReceipt(
        groupId: groupId,
        paidByUserId: fixture.userId,
        jwt: jwt,
        name: alpha,
        amount: '10.00');
    await createReceipt(
        groupId: groupId,
        paidByUserId: fixture.userId,
        jwt: jwt,
        name: bravo,
        amount: '25.50',
        status: 'RESOLVED');
    await createReceipt(
        groupId: groupId,
        paidByUserId: fixture.userId,
        jwt: jwt,
        name: charlie,
        amount: '4.50');

    // OPEN broken out, at the historical position. setGroupSummaryConfig restores all
    // four keys on teardown -- an omitted key means "leave unchanged", so replaying the
    // base command would leave the summary switched on for every later spec.
    await setGroupSummaryConfig(
      groupId: groupId,
      jwt: jwt,
      enabled: true,
      position: 'BOTTOM',
      statuses: ['OPEN'],
    );

    await loginAs(tester,
        username: fixture.username, password: fixture.password);
    await openGroupReceipts(tester, fixture.groupName!, alpha);

    // --- 1. the figures come off a real decimal fold over every matching receipt
    await pumpUntilFound(tester, summaryBar());
    expect(rowText(tester, overallRow()), contains('3 receipts'));
    expect(find.textContaining('40.00'), findsWidgets);
    expect(rowText(tester, openRow()), contains('2 receipts'));

    // --- 2. it is BELOW the list, which is where it rendered before the setting existed
    expect(
      tester.getTopLeft(summaryBar()).dy,
      greaterThan(tester.getTopLeft(find.byType(PagedDataList)).dy),
    );

    // --- 3. a filter recomputes every row, not just the rows above
    await openFilterScreen(tester);
    await addCondition(tester, 'name', _textValueField());
    await tester.enterText(formField('value'), alpha);
    await drainFrames(tester);
    await saveCondition(tester, 'name');
    await applyFilter(tester);

    // Absence first: the unfiltered figures are already on screen, so asserting the
    // narrowed ones straight away could pass without the request ever having run.
    await pumpUntilGone(tester, receiptRow(bravo));
    await pumpUntilFound(tester, find.textContaining('10.00'));
    expect(rowText(tester, overallRow()), contains('1 receipt'),
        reason: 'the encoded filter must reach the summary endpoint too, '
            'and the count must singularize');
    expect(rowText(tester, openRow()), contains('1 receipt'));

    // --- 4. a configured status matching nothing keeps its row, so the block does not
    // change shape as the filter narrows.
    await openFilterScreen(tester);
    await editCondition(tester, 'name', _textValueField());
    await tester.enterText(formField('value'), bravo);
    await drainFrames(tester);
    await saveCondition(tester, 'name');
    await applyFilter(tester);

    await pumpUntilFound(tester, receiptRow(bravo));
    expect(openRow(), findsOneWidget,
        reason: 'a configured status with no matches still renders, as a zero row');
    expect(rowText(tester, openRow()), contains('0 receipts'));
  });

  // The new setting's own path: it is written on the group and read back off the SUMMARY
  // RESPONSE, never off the client's cached group settings -- so this is the only level
  // that proves model -> command -> DB -> response -> render.
  testWidgets('the position setting moves the bar above the list', (tester) async {
    await tester.binding.setSurfaceSize(const Size(430, 932));
    addTearDown(() => tester.binding.setSurfaceSize(null));

    final jwt = await apiLogin();
    final stamp = DateTime.now().microsecondsSinceEpoch.toString();
    final fixture = await provisionPermUser(roleName: 'Legacy Editor');
    final name = 'e2e-sumpos-$stamp';

    await createReceipt(
        groupId: fixture.groupId!,
        paidByUserId: fixture.userId,
        jwt: jwt,
        name: name,
        amount: '7.25');

    await setGroupSummaryConfig(
        groupId: fixture.groupId!, jwt: jwt, enabled: true, position: 'TOP');

    await loginAs(tester,
        username: fixture.username, password: fixture.password);
    await openGroupReceipts(tester, fixture.groupName!, name);

    await pumpUntilFound(tester, summaryBar());
    expect(
      tester.getTopLeft(summaryBar()).dy,
      lessThan(tester.getTopLeft(find.byType(PagedDataList)).dy),
      reason: 'TOP must survive the whole round trip, not just the settings form',
    );
  });
}

/// The text condition editor's value field, keyed rather than found by name so the
/// finder does not depend on the editor's internal form field naming.
Finder _textValueField() =>
    find.byKey(const ValueKey('receipt-filter-text-value'));
