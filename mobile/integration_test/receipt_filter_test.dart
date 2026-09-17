import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';

import 'helpers/api.dart';
import 'helpers/form_actions.dart';
import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/receipt_filter_actions.dart';

/// End-to-end coverage for **receipt filtering**, against the real Go API.
///
/// The unit and widget suites already pin the filter exhaustively on the client
/// side -- the encoder through the real serializer, the draft-vs-applied split,
/// the save gate, the count header, the card's gesture arena. All of it runs
/// against a mocked or absent API, so what none of it can prove is the part most
/// likely to break: that the encoded filter is **accepted by the Go query
/// builder and actually narrows a real result set**. Every value is
/// type-asserted server-side with no comma-ok, so a wrong shape is a 500 rather
/// than an ignored filter (see mobile/CLAUDE.md > "Receipt filtering").
///
/// These tests therefore drive the real pickers -- the multiselect sheet, the
/// per-group catalog, the date range picker, the currency field -- none of which
/// any widget test taps, and assert over rows the server sent back.
///
/// **Seed set** (one fixture group, three receipts, one category):
///
/// | | name | amount | date | status | category |
/// |---|---|---|---|---|---|
/// | R1 | `…-alpha`   | 12.34 | 2026-06-02T09:00Z | OPEN     | -- |
/// | R2 | `…-bravo`   | 99.99 | 2026-06-13T18:30Z | RESOLVED | yes |
/// | R3 | `…-charlie` | 12.34 | 2026-06-20T09:00Z | OPEN     | -- |
///
/// Every test gets one match and two non-matches out of those three rows.
///
/// **R2's 18:30 is load-bearing.** A BETWEEN of 06/11-06/13 encodes to
/// `['2026-06-11T00:00:00Z', '2026-06-13T23:59:59Z']` because the client expands
/// the range to start-of-day/end-of-day, and the server compares against a
/// datetime column. Were R2 stored at midnight, a *broken* `<= 06-13T00:00:00Z`
/// upper bound would still match it and the date test would pass against the
/// very bug it exists to catch.
///
/// The dates sit in June 2026 deliberately: `WITHIN_CURRENT_MONTH` is out of
/// scope here (the server pins that range to the current month itself), and
/// `date EQUALS` is avoided because a datetime column only equals an exact
/// midnight. Anyone adding either case must seed relative to `DateTime.now()`.
///
/// **Timezone.** `formatDate` converts to local time and then stamps a literal
/// `Z`, so the seeds and the picked range only agree on a UTC host. This
/// container and CI are UTC; a non-UTC runner would break this spec -- and
/// production date filtering with it.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('a name filter narrows the list, and reset restores it',
      (tester) async {
    final seed = await _seedAndEnterGroup(tester);

    await openFilterScreen(tester);
    expect(filterEmptyState(), findsOneWidget,
        reason: 'a fresh group carries no filter');

    await _addNameCondition(tester, '${seed.stamp}-alpha');
    expect(find.text('1 CONDITION'), findsOneWidget);
    await applyFilter(tester);

    // Absence first: every row is already on screen before the filter is
    // applied, so asserting the match first would pass without the filter ever
    // having run. Waiting for an excluded row to GO is the proof of the round
    // trip.
    await pumpUntilGone(tester, receiptRow(seed.bravo));
    expect(receiptRow(seed.alpha), findsOneWidget);
    expect(receiptRow(seed.charlie), findsNothing);
    expect(filterBadge('1'), findsOneWidget);

    // Drive it to zero matches THROUGH THE API. This is what arms the paging
    // bug below: PagedDataList stops requesting pages once the loaded count
    // reaches the total, and a total of zero satisfies that forever.
    await openFilterScreen(tester);
    await editCondition(tester, 'name', _textValueField());
    await tester.enterText(formField('value'), 'zzz-${seed.stamp}');
    await drainFrames(tester);
    await saveCondition(tester, 'name');
    await applyFilter(tester);
    await pumpUntilFound(tester, find.text('No receipts match this filter'));

    // Reset and re-apply. Without `_totalCount = null` in PagedDataList's
    // refresh callback, the stale zero total means no page is ever requested
    // again and the list stays empty forever.
    await openFilterScreen(tester);
    await settleTap(tester, resetButton());
    await pumpUntilFound(tester, filterEmptyState());
    await applyFilter(tester);

    await pumpUntilFound(tester, receiptRow(seed.charlie));
    expect(receiptRow(seed.alpha), findsOneWidget);
    expect(receiptRow(seed.bravo), findsOneWidget);
    expect(filterBadge('1'), findsNothing,
        reason: 'the badge clears with the filter');
  });

  testWidgets('a status filter picked in the real sheet reaches the server',
      (tester) async {
    final seed = await _seedAndEnterGroup(tester);

    await openFilterScreen(tester);
    const statusField = ValueKey('receipt-filter-status');
    await addCondition(tester, 'status', find.byKey(statusField));

    // List-type fields offer CONTAINS only, and it is preselected -- there is
    // no operation to choose, which is itself worth pinning.
    expect(operationChip('CONTAINS'), findsOneWidget);
    expect(find.text('No Statuses selected'), findsOneWidget);

    await settleTap(tester, find.byKey(statusField));
    await pumpUntilFound(tester, find.widgetWithText(ChoiceChip, 'Resolved'));
    await pickFromMultiselect(tester, 'Resolved');

    // The sheet's result must be written back before Save is poked; the applied
    // selection renders as an InputChip, the sheet's options as ChoiceChips.
    await pumpUntilFound(tester, selectedChip('Resolved'));
    await saveCondition(tester, 'status');
    await applyFilter(tester);

    await pumpUntilGone(tester, receiptRow(seed.alpha));
    expect(receiptRow(seed.bravo), findsOneWidget,
        reason: 'the status wire name RESOLVED must reach `status IN (?)`');
    expect(receiptRow(seed.charlie), findsNothing);
    expect(filterBadge('1'), findsOneWidget);
  });

  testWidgets('a category filter uses the real per-group catalog',
      (tester) async {
    final seed = await _seedAndEnterGroup(tester);

    await openFilterScreen(tester);
    const categoriesField = ValueKey('receipt-filter-categories');
    await addCondition(tester, 'categories', find.byKey(categoriesField));

    await settleTap(tester, find.byKey(categoriesField));
    // The fixture user is not an app admin, so the catalog it sees here came
    // from the group-scoped AppData path -- which no widget test exercises
    // (they all hand the editor a hand-built harness catalog).
    await pumpUntilFound(tester, formField('filter'));
    await pickFromMultiselect(tester, seed.categoryName, useFilterBox: true);

    await pumpUntilFound(tester, selectedChip(seed.categoryName));
    await saveCondition(tester, 'categories');
    await applyFilter(tester);

    await pumpUntilGone(tester, receiptRow(seed.alpha));
    expect(receiptRow(seed.bravo), findsOneWidget,
        reason: 'category ids must reach the receipt_categories subquery');
    expect(receiptRow(seed.charlie), findsNothing);
    expect(filterBadge('1'), findsOneWidget);
  });

  testWidgets('a date range includes a receipt late on its final day',
      (tester) async {
    final seed = await _seedAndEnterGroup(tester);

    await openFilterScreen(tester);
    await addCondition(
        tester, 'date', find.byKey(const ValueKey('receipt-filter-date')));

    await settleTap(tester, operationChip('BETWEEN'));
    final rangeField = find.byKey(const ValueKey('receipt-filter-date-range'));
    await pumpUntilFound(tester, rangeField);
    await drainFrames(tester);
    await tester.tap(rangeField);

    await _pickDateRange(tester, start: '06/11/2026', end: '06/13/2026');

    // Prove the picked range reached the editor's state BEFORE any network
    // call. If this passes and the row assertions below fail, the server
    // disagreed -- which is a real finding, not a broken picker.
    expect(
      find.descendant(
          of: rangeField, matching: find.text('06/11/2026 - 06/13/2026')),
      findsOneWidget,
    );

    await saveCondition(tester, 'date');
    await applyFilter(tester);

    await pumpUntilGone(tester, receiptRow(seed.charlie));
    expect(
      receiptRow(seed.bravo),
      findsOneWidget,
      reason: 'R2 is dated 18:30 on the range\'s LAST day -- it matches only '
          'because the client expands the upper bound to end-of-day. A bare '
          '<= 2026-06-13T00:00:00Z would exclude it.',
    );
    expect(receiptRow(seed.alpha), findsNothing,
        reason: 'R1 (06/02) is before the range');
    expect(filterBadge('1'), findsOneWidget);
  });

  testWidgets('an amount filter reaches the server as a number',
      (tester) async {
    final seed = await _seedAndEnterGroup(tester);

    await openFilterScreen(tester);
    // NOT find.byKey: AmountField forwards its key to the FormBuilderTextField
    // it builds, so the keyed finder matches two widgets and throws.
    await addCondition(tester, 'amount', formField('value'));

    await settleTap(tester, operationChip('GREATER_THAN'));
    // Switching operation remounts AmountField under a new key, with a fresh
    // controller seeded at zero.
    await pumpUntilFound(tester, formField('value'));
    await drainFrames(tester);
    // CurrencyTextFieldController reads keystrokes as cents, so the full cents
    // form is required: '50' would enter 0.50.
    await tester.enterText(formField('value'), '50.00');
    await drainFrames(tester);

    await saveCondition(tester, 'amount');
    await applyFilter(tester);

    await pumpUntilGone(tester, receiptRow(seed.alpha));
    expect(receiptRow(seed.bravo), findsOneWidget,
        reason: 'amount must be sent as a JSON number -- a string is a 500');
    expect(receiptRow(seed.charlie), findsNothing);
    expect(filterBadge('1'), findsOneWidget);
  });
}

// ---------------------------------------------------------------------------
// Seeding
// ---------------------------------------------------------------------------

class _Seed {
  _Seed({
    required this.fixture,
    required this.stamp,
    required this.categoryName,
    required this.alpha,
    required this.bravo,
    required this.charlie,
  });

  final PermFixture fixture;
  final String stamp;
  final String categoryName;
  final String alpha;
  final String bravo;
  final String charlie;
}

/// Seeds the group described in this file's header, logs in as its member and
/// opens the receipts list with all three rows on screen.
Future<_Seed> _seedAndEnterGroup(WidgetTester tester) async {
  // One admin JWT for every write AND for our own teardowns. apiLogin throws on
  // a non-200, and a throw inside a teardown aborts the rest of that closure --
  // re-logging-in there is how fixtures leak when the login endpoint is busy.
  final jwt = await apiLogin();
  final stamp = DateTime.now().microsecondsSinceEpoch.toString();

  final categoryName = 'e2e-flt-cat-$stamp';
  final categoryId = await createCategory(name: categoryName, jwt: jwt);
  // Registered FIRST so LIFO runs it LAST -- after the group cascade has
  // removed the receipt_categories join rows that reference it.
  addTearDown(() async => deleteCategory(categoryId, jwt: jwt));

  final fixture = await provisionPermUser(roleName: 'Legacy Editor');
  final groupId = fixture.groupId!;

  final alpha = 'e2e-flt-$stamp-alpha';
  final bravo = 'e2e-flt-$stamp-bravo';
  final charlie = 'e2e-flt-$stamp-charlie';

  // No per-receipt teardown: DELETE /group/{id} cascades its receipts, and the
  // fixture group's delete is already registered.
  await createReceipt(
    groupId: groupId,
    paidByUserId: fixture.userId,
    jwt: jwt,
    name: alpha,
    amount: '12.34',
    date: '2026-06-02T09:00:00Z',
  );
  await createReceipt(
    groupId: groupId,
    paidByUserId: fixture.userId,
    jwt: jwt,
    name: bravo,
    amount: '99.99',
    date: '2026-06-13T18:30:00Z',
    status: 'RESOLVED',
    categories: [(id: categoryId, name: categoryName)],
  );
  await createReceipt(
    groupId: groupId,
    paidByUserId: fixture.userId,
    jwt: jwt,
    name: charlie,
    amount: '12.34',
    date: '2026-06-20T09:00:00Z',
  );

  await loginAs(
    tester,
    username: fixture.username,
    password: fixture.password,
  );
  await openGroupReceipts(tester, fixture.groupName!, alpha);
  // All three landed before any filter is authored.
  await pumpUntilFound(tester, receiptRow(charlie));

  return _Seed(
    fixture: fixture,
    stamp: stamp,
    categoryName: categoryName,
    alpha: alpha,
    bravo: bravo,
    charlie: charlie,
  );
}

// ---------------------------------------------------------------------------
// Local drivers
// ---------------------------------------------------------------------------

Finder _textValueField() =>
    find.byKey(const ValueKey('receipt-filter-text-value'));

/// Adds a `name CONTAINS [value]` condition and saves it.
Future<void> _addNameCondition(WidgetTester tester, String value) async {
  await addCondition(tester, 'name', _textValueField());
  await tester.enterText(formField('value'), value);
  await drainFrames(tester);
  await saveCondition(tester, 'name');
}

/// Fills Flutter's date range picker via its **text entry** mode.
///
/// The calendar grid is not drivable in a test: `firstDate` is 2000 and
/// `lastDate` is five years out, so it is ~370 lazily-built month items whose
/// day cells are bare unkeyed `Text('11')` repeating every month. The picker's
/// own text mode is deterministic instead -- and it is a real user affordance,
/// so this still exercises the picker rather than going around it.
///
/// The app registers no `flutter_localizations` delegates, so
/// `DefaultMaterialLocalizations` is in force: dates parse as US `mm/dd/yyyy`,
/// the toggle is tooltipped "Switch to input", and the confirm button in text
/// mode is "OK" (it is "Save" in calendar mode).
Future<void> _pickDateRange(
  WidgetTester tester, {
  required String start,
  required String end,
}) async {
  final dialog = find.byType(DateRangePickerDialog);
  await pumpUntilFound(tester, dialog);
  await drainFrames(tester);

  await settleTap(
    tester,
    find.descendant(of: dialog, matching: find.byTooltip('Switch to input')),
  );
  // The dialog resizes over ~200ms on the mode swap.
  await drainFrames(tester, frames: 8, milliseconds: 50);

  Finder field(String label) =>
      find.ancestor(of: find.text(label), matching: find.byType(TextField));

  // The dialog header also renders the literal "Start Date" while nothing is
  // selected, but that Text has no TextField ancestor, so this stays unique.
  await pumpUntilFound(tester, field('Start Date'));
  await tester.enterText(field('Start Date'), start);
  await tester.pump();
  await tester.enterText(field('End Date'), end);
  await tester.pump();

  await settleTap(
    tester,
    find.descendant(
        of: dialog, matching: find.widgetWithText(TextButton, 'OK')),
  );
  await pumpUntilGone(tester, dialog);
  await drainFrames(tester);
}
