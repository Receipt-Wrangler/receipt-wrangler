import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_date_filter.dart';

import 'helpers/api.dart';
import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/receipt_filter_actions.dart';

/// End-to-end coverage for the **quick date filter** -- the month stepper and
/// the date-field picker above the receipts list -- against the real Go API.
///
/// The widget suite already pins the whole control against a mocked
/// `ReceiptApi`: the month it writes, the single refetch, the field switch that
/// refetches nothing. What it cannot prove is that the `BETWEEN` the stepper
/// builds is accepted by the Go query builder and actually narrows real rows --
/// and, for the picker, that the *column* the server filters on really moves.
/// Every filter value is type-asserted server-side with no comma-ok, so a wrong
/// shape is a 500 rather than an ignored filter.
///
/// **Seeds relative to `DateTime.now()`, deliberately.** The sibling
/// `receipt_filter_test.dart` pins its rows to June 2026 and warns that any
/// month-relative case must seed its own data -- which this whole file is: the
/// stepper's months are computed from today, so fixed dates would drift out of
/// range the moment the calendar moved past them.
///
/// **Seed set** (one fixture group, two receipts):
///
/// | | name | date | status | resolved_date |
/// |---|---|---|---|---|
/// | LAST | `…-last` | the 15th of **last** month | RESOLVED | **now** (this month) |
/// | THIS | `…-this` | the 15th of **this** month | OPEN     | none |
///
/// The RESOLVED row is what makes the field picker testable without a second
/// seeding axis: the API stamps `resolved_date` with `time.Now().UTC()` on
/// create (`api/internal/repositories/receipts.go`), so LAST's receipt date and
/// its resolved date fall in **different months**. Filtering "this month" then
/// returns THIS on `date` and LAST on `resolvedDate` -- opposite rows for the
/// same range, which nothing but a real column change can produce.
///
/// The 15th on both sides keeps every seed clear of a month boundary, so no
/// case depends on which day the suite happens to run.
///
/// **Timezone: this spec needs a UTC host**, for the same reason
/// `receipt_filter_test.dart` does -- it seeds through the API with literal `Z`
/// strings while the client stamps local wall clock with a `Z`. The container
/// and CI are UTC.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('stepping back a month narrows the list to that month',
      (tester) async {
    final seed = await _seedAndEnterGroup(tester);

    expect(monthStepperLabel('All time'), findsOneWidget,
        reason: 'a fresh group carries no date filter');

    await settleTap(tester, _prevMonth());

    await pumpUntilGone(tester, receiptRow(seed.thisMonth));
    expect(receiptRow(seed.lastMonth), findsOneWidget,
        reason: 'a month must reach the server as a BETWEEN it accepts');
    expect(monthStepperLabel(filterMonthLabel(_lastMonth)), findsOneWidget);
    expect(filterBadge('1'), findsOneWidget);
  });

  testWidgets('the date-field picker moves the filter to another column',
      (tester) async {
    final seed = await _seedAndEnterGroup(tester);

    // Same range on both halves of this test -- only the column differs, so
    // the rows swapping is the proof.
    await _pickThisMonth(tester);
    await pumpUntilGone(tester, receiptRow(seed.lastMonth));
    expect(receiptRow(seed.thisMonth), findsOneWidget,
        reason: 'on Receipt Date, this month is the row dated this month');

    // Back to all time first: switching the field is non-destructive, so a
    // Receipt Date condition left applied would go on the wire alongside the
    // Resolved Date one and match nothing.
    await settleTap(tester, _clearMonth());
    await pumpUntilFound(tester, receiptRow(seed.lastMonth));

    await _chooseDateField(tester, 'resolvedDate');
    await _pickThisMonth(tester);

    await pumpUntilGone(tester, receiptRow(seed.thisMonth));
    expect(receiptRow(seed.lastMonth), findsOneWidget,
        reason: 'on Resolved Date the same month returns the OTHER row -- '
            'only a real column change can do that');
  });

  testWidgets('clearing the stepper returns every receipt', (tester) async {
    final seed = await _seedAndEnterGroup(tester);

    await settleTap(tester, _prevMonth());
    await pumpUntilGone(tester, receiptRow(seed.thisMonth));

    await settleTap(tester, _clearMonth());

    await pumpUntilFound(tester, receiptRow(seed.thisMonth));
    expect(receiptRow(seed.lastMonth), findsOneWidget);
    expect(monthStepperLabel('All time'), findsOneWidget);
    expect(filterBadge('1'), findsNothing);
  });
}

// ---------------------------------------------------------------------------
// Local drivers
// ---------------------------------------------------------------------------

final _thisMonth = monthOfDate(DateTime.now());
final _lastMonth = shiftMonth(_thisMonth, -1);

Finder _prevMonth() => find.byKey(const ValueKey('receipt-month-prev'));

Finder _clearMonth() => find.byKey(const ValueKey('receipt-month-clear'));

/// The stepper's label, scoped so it cannot match the same text elsewhere.
Finder monthStepperLabel(String text) => find.descendant(
      of: find.byKey(const ValueKey('receipt-month-label')),
      matching: find.text(text),
    );

/// Opens the month sheet and takes its "This month" shortcut.
Future<void> _pickThisMonth(WidgetTester tester) async {
  await settleTap(tester, find.byKey(const ValueKey('receipt-month-label')));
  await pumpUntilFound(
      tester, find.byKey(const ValueKey('month-picker-this-month')));
  await settleTap(
      tester, find.byKey(const ValueKey('month-picker-this-month')));
  await pumpUntilGone(
      tester, find.byKey(const ValueKey('month-picker-this-month')));
}

/// Re-points the stepper at [fieldKey] through the chip's popup menu.
Future<void> _chooseDateField(WidgetTester tester, String fieldKey) async {
  await settleTap(
      tester, find.byKey(const ValueKey('receipt-quick-date-field')));
  final item = find.byKey(ValueKey('receipt-quick-date-field-$fieldKey'));
  await pumpUntilFound(tester, item);
  await settleTap(tester, item);
  await pumpUntilGone(tester, item);
}

// ---------------------------------------------------------------------------
// Seeding
// ---------------------------------------------------------------------------

class _Seed {
  _Seed({required this.fixture, required this.lastMonth, required this.thisMonth});

  final PermFixture fixture;
  final String lastMonth;
  final String thisMonth;
}

/// The 15th of [month], as the literal zulu string the seeding API takes.
String _seedDate(FilterMonth month) {
  final padded = month.month.toString().padLeft(2, '0');
  return '${month.year}-$padded-15T09:00:00Z';
}

/// Seeds the two receipts described in this file's header, logs in as their
/// group's member and opens the receipts list with both rows on screen.
Future<_Seed> _seedAndEnterGroup(WidgetTester tester) async {
  // One admin JWT for every write. apiLogin throws on a non-200, and a throw
  // inside a teardown aborts the rest of that closure.
  final jwt = await apiLogin();
  final stamp = DateTime.now().microsecondsSinceEpoch.toString();

  final fixture = await provisionPermUser(roleName: 'Legacy Editor');
  final groupId = fixture.groupId!;

  final lastMonthName = 'e2e-qdf-$stamp-last';
  final thisMonthName = 'e2e-qdf-$stamp-this';

  // No per-receipt teardown: DELETE /group/{id} cascades its receipts, and the
  // fixture group's delete is already registered.
  //
  // RESOLVED is load-bearing, not decoration: it is what makes the API stamp
  // resolved_date with *now*, putting this row's two date columns in different
  // months.
  await createReceipt(
    groupId: groupId,
    paidByUserId: fixture.userId,
    jwt: jwt,
    name: lastMonthName,
    date: _seedDate(_lastMonth),
    status: 'RESOLVED',
  );
  await createReceipt(
    groupId: groupId,
    paidByUserId: fixture.userId,
    jwt: jwt,
    name: thisMonthName,
    date: _seedDate(_thisMonth),
  );

  await loginAs(
    tester,
    username: fixture.username,
    password: fixture.password,
  );
  await openGroupReceipts(tester, fixture.groupName!, lastMonthName);
  // Both landed before any filter is authored.
  await pumpUntilFound(tester, receiptRow(thisMonthName));

  return _Seed(
    fixture: fixture,
    lastMonth: lastMonthName,
    thisMonth: thisMonthName,
  );
}
