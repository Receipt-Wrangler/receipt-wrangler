import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_edit_popup_menu.dart';

import 'helpers/api.dart';
import 'helpers/form_actions.dart';
import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/receipt_filter_actions.dart';

/// Lifecycle coverage for the receipt filter: how long an applied filter lives,
/// and when it must stop applying.
///
/// Neither behaviour is reachable from a widget test. The filter lives on
/// `ReceiptListModel`, an app-level provider, while the list that reads it is
/// rebuilt from scratch on every navigation -- so what these tests exercise is
/// the seam between a provider that outlives the screen and a screen that has to
/// re-derive its query on mount. Only a full `buildApp()` under the real
/// go_router has that seam.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('an applied filter survives a round trip to a receipt',
      (tester) async {
    final seed = await _seedTwoReceipts(tester);

    await openFilterScreen(tester);
    await _addNameCondition(tester, '${seed.stamp}-keep');
    await applyFilter(tester);

    await pumpUntilGone(tester, receiptRow(seed.dropName));
    expect(receiptRow(seed.keepName), findsOneWidget);
    expect(filterBadge('1'), findsOneWidget);

    // `/receipts/:id/view` is a TOP-LEVEL route, so the whole group shell --
    // the list, its PagedDataList, the paging controller, the total count --
    // is torn down. Only ReceiptListModel carries the filter across, and the
    // freshly-mounted list has to refetch *with* it.
    await settleTap(tester, receiptRow(seed.keepName));
    await pumpUntilFound(tester, find.byType(ReceiptEditPopupMenu));

    await _tapBackArrow(tester);
    await pumpUntilFound(tester, receiptRow(seed.keepName));

    expect(receiptRow(seed.dropName), findsNothing,
        reason: 'the applied filter must still be narrowing the rebuilt list');
    expect(filterBadge('1'), findsOneWidget,
        reason: 'and the app bar must still say so');

    // Stronger than the badge: the model still holds the condition, so
    // reopening the screen re-seeds its draft from it rather than showing the
    // empty state.
    await openFilterScreen(tester);
    expect(conditionCard('name'), findsOneWidget);
    expect(filterEmptyState(), findsNothing);
    await settleTap(
        tester, find.byKey(const ValueKey('receipt-filter-close')));
  });

  testWidgets('switching groups clears the filter', (tester) async {
    final jwt = await apiLogin();
    final stamp = DateTime.now().microsecondsSinceEpoch.toString();
    final twoGroups = await provisionPermUserWithTwoGroups();
    final fixture = twoGroups.fixture;

    final inGroupA = 'e2e-flt2-$stamp-ga';
    final inGroupB = 'e2e-flt2-$stamp-gb';
    await createReceipt(
      groupId: fixture.groupId!,
      paidByUserId: fixture.userId,
      jwt: jwt,
      name: inGroupA,
    );
    await createReceipt(
      groupId: twoGroups.secondGroupId,
      paidByUserId: fixture.userId,
      jwt: jwt,
      name: inGroupB,
    );

    await loginAs(
      tester,
      username: fixture.username,
      password: fixture.password,
    );
    await openGroupReceipts(tester, fixture.groupName!, inGroupA);

    await openFilterScreen(tester);
    await _addNameCondition(tester, '$stamp-ga');
    await applyFilter(tester);
    await pumpUntilFound(tester, receiptRow(inGroupA));
    expect(filterBadge('1'), findsOneWidget);

    // The real journey, and the only one the app offers: there is no lateral
    // group switch. The app-bar arrow goes to /groups, group cards go to
    // /groups/<id>/dashboards, and the bottom nav only ever navigates within
    // the current group.
    await _tapBackArrow(tester);
    await openGroupReceipts(tester, twoGroups.secondGroupName, inGroupB);

    // Group B's receipt is on screen, so the filter is not applying here. Had
    // it followed the user across, "$stamp-ga" would match nothing in B and
    // this list would render "No receipts match this filter" instead.
    expect(receiptRow(inGroupB), findsOneWidget);
    expect(find.text('No receipts match this filter'), findsNothing,
        reason: "group A's filter must not follow the user into group B");
    expect(receiptRow(inGroupA), findsNothing,
        reason: "and group A's receipt is not in group B");
    expect(filterBadge('1'), findsNothing,
        reason: 'the badge must not count conditions that no longer apply');

    await openFilterScreen(tester);
    expect(filterEmptyState(), findsOneWidget);
    expect(conditionCard('name'), findsNothing);
  });
}

// ---------------------------------------------------------------------------
// Seeding
// ---------------------------------------------------------------------------

class _TwoReceipts {
  _TwoReceipts({
    required this.stamp,
    required this.keepName,
    required this.dropName,
  });

  final String stamp;

  /// Matches the `-keep` name condition these tests author.
  final String keepName;

  /// Filtered out by it.
  final String dropName;
}

Future<_TwoReceipts> _seedTwoReceipts(WidgetTester tester) async {
  final jwt = await apiLogin();
  final stamp = DateTime.now().microsecondsSinceEpoch.toString();
  final fixture = await provisionPermUser(roleName: 'Legacy Editor');

  final keepName = 'e2e-flt-$stamp-keep';
  final dropName = 'e2e-flt-$stamp-drop';
  // No per-receipt teardown: the fixture group's delete cascades its receipts.
  await createReceipt(
    groupId: fixture.groupId!,
    paidByUserId: fixture.userId,
    jwt: jwt,
    name: keepName,
  );
  await createReceipt(
    groupId: fixture.groupId!,
    paidByUserId: fixture.userId,
    jwt: jwt,
    name: dropName,
  );

  await loginAs(
    tester,
    username: fixture.username,
    password: fixture.password,
  );
  await openGroupReceipts(tester, fixture.groupName!, keepName);
  await pumpUntilFound(tester, receiptRow(dropName));

  return _TwoReceipts(stamp: stamp, keepName: keepName, dropName: dropName);
}

// ---------------------------------------------------------------------------
// Local drivers
// ---------------------------------------------------------------------------

Future<void> _addNameCondition(WidgetTester tester, String value) async {
  await addCondition(
      tester, 'name', find.byKey(const ValueKey('receipt-filter-text-value')));
  await tester.enterText(formField('value'), value);
  await drainFrames(tester);
  await saveCondition(tester, 'name');
}

/// Taps the visible back arrow. `.hitTestable()` picks the single tappable one:
/// the `/view` route can leave a stale offstage `ReceiptFormScreen` (and so a
/// second AppBar) behind the visible screen -- see receipt_navigation_test.dart.
Future<void> _tapBackArrow(WidgetTester tester) async {
  final back = find.byIcon(Icons.arrow_back).hitTestable();
  await pumpUntilFound(tester, back);
  await drainFrames(tester);
  await tester.tap(back);
}
