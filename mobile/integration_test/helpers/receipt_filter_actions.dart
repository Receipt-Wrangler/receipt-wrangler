import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/receipt_list_item.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';

import 'form_actions.dart';
import 'pump.dart';

/// Shared drivers for the receipt-filter e2e specs.
///
/// The filter is authored across three stacked surfaces -- a pushed route, the
/// "Add a filter" sheet and the condition editor sheet -- and each transition
/// needs the suite's tap-flake handling (see mobile/CLAUDE.md "Three tap-flake
/// patterns"). Seven tests restating that by hand is how one of them ends up
/// subtly different and flakes only in CI, so the sequences live here.
///
/// **Two finder rules these helpers exist to enforce:**
///
///  * **Wait for *enabled*, not hittable.** A disabled `FilledButton` still hit
///    tests, so `tester.tap` on one silently does nothing and the next
///    `pumpUntilFound` fails ten seconds later pointing at the wrong thing.
///    "Apply Filter" and the picker's "Select" are gated on
///    `LoadingModel.isLoading`; "Save condition" is gated on the editor's own
///    validity. [enabledFilledButton] / [enabledSubmitButton] cover all three.
///  * **Never `find.byKey` an amount field.** `AmountField` forwards its
///    `widget.key` onto the `FormBuilderTextField` it builds, so
///    `find.byKey(ValueKey("receipt-filter-amount-EQUALS"))` matches *two*
///    widgets and any tap or `enterText` against it throws. Use
///    `formField("value")` -- that is deliberate, not an oversight.

Finder filterButton() => find.byKey(const ValueKey("receipt-filter-button"));

Finder addFilterButton() => find.byKey(const ValueKey("receipt-filter-add"));

Finder applyButton() => find.byKey(const ValueKey("receipt-filter-apply"));

Finder resetButton() => find.byKey(const ValueKey("receipt-filter-reset"));

Finder filterEmptyState() =>
    find.byKey(const ValueKey("receipt-filter-empty"));

Finder conditionCard(String fieldKey) =>
    find.byKey(ValueKey("receipt-filter-card-$fieldKey"));

Finder operationChip(String operationName) =>
    find.byKey(ValueKey("receipt-filter-operation-$operationName"));

/// A receipts-list row for [receiptName].
///
/// Scoped to [ReceiptListItem] rather than a bare `find.text`: the filter screen
/// renders the same string on the condition card (a name condition's value), and
/// the list stays mounted underneath the pushed route.
Finder receiptRow(String receiptName) =>
    find.widgetWithText(ReceiptListItem, receiptName);

/// The filter action's badge count, scoped to the filter button so no other
/// `Badge` in the shell can satisfy it.
Finder filterBadge(String count) => find.descendant(
      of: find.ancestor(of: filterButton(), matching: find.byType(Badge)),
      matching: find.text(count),
    );

/// A [FilledButton] carrying [key] whose `onPressed` is non-null.
Finder enabledFilledButton(Key key) => find.byWidgetPredicate(
      (w) => w is FilledButton && w.key == key && w.onPressed != null,
    );

/// The enabled [FilledButton] *inside* a [BottomSubmitButton] keyed [key].
/// The key sits on the wrapper, the enablement on the button it builds.
Finder enabledSubmitButton(Key key) => find.descendant(
      of: find.byKey(key),
      matching: find.byWidgetPredicate(
          (w) => w is FilledButton && w.onPressed != null),
    );

/// Drains frames so a follow-up tap computes its centre from settled geometry.
/// A single `pump` advances one frame; a sheet or route transition needs more.
Future<void> drainFrames(WidgetTester tester,
    {int frames = 5, int milliseconds = 100}) async {
  for (var i = 0; i < frames; i++) {
    await tester.pump(Duration(milliseconds: milliseconds));
  }
}

/// Wait for hittability, drain the animation, then tap -- the suite's standard
/// sequence for anything on an animating surface.
Future<void> settleTap(WidgetTester tester, Finder finder) async {
  await pumpUntilFound(tester, finder.hitTestable());
  await drainFrames(tester);
  await tester.tap(finder.hitTestable());
}

/// Opens the filter screen from the receipts app bar.
Future<void> openFilterScreen(WidgetTester tester) async {
  await settleTap(tester, filterButton());
  // "Add filter" is on the screen in both states (empty and populated), so it
  // is the landing marker. The screen is a `Navigator.push`, not a go_router
  // location, so there is no URL to assert against.
  await pumpUntilFound(tester, addFilterButton());
  await drainFrames(tester);
}

/// "+ Add filter" -> pick [fieldKey] -> returns with the condition editor
/// mounted, waiting on [editorMarker] (the field's own value input).
Future<void> addCondition(
  WidgetTester tester,
  String fieldKey,
  Finder editorMarker,
) async {
  await settleTap(tester, addFilterButton());
  await settleTap(tester, find.byKey(ValueKey("add-receipt-filter-$fieldKey")));
  // The add sheet pops and the editor sheet opens in one step -- picking a
  // field goes straight into editing it.
  await pumpUntilFound(tester, editorMarker);
  await drainFrames(tester);
}

/// Reopens an existing condition by tapping its card.
Future<void> editCondition(
  WidgetTester tester,
  String fieldKey,
  Finder editorMarker,
) async {
  await settleTap(tester, conditionCard(fieldKey));
  await pumpUntilFound(tester, editorMarker);
  await drainFrames(tester);
}

/// Taps "Save condition" and returns once the condition's card is on the filter
/// screen. Waiting for the button to be *enabled* first doubles as proof the
/// authored value reached the editor's state.
Future<void> saveCondition(WidgetTester tester, String fieldKey) async {
  const key = ValueKey("receipt-filter-condition-save");
  await pumpUntilFound(tester, enabledFilledButton(key));
  await drainFrames(tester);
  await tester.tap(find.byKey(key));
  await pumpUntilFound(tester, conditionCard(fieldKey));
  await drainFrames(tester);
}

/// Commits the draft and returns once the filter screen has popped.
Future<void> applyFilter(WidgetTester tester) async {
  await pumpUntilFound(
      tester, enabledSubmitButton(const ValueKey("receipt-filter-apply")));
  await drainFrames(tester);
  await tester.tap(applyButton());
  await pumpUntilGone(tester, applyButton());
  await drainFrames(tester);
}

/// Picks [label] in the multiselect sheet a list/users field opens, then
/// confirms with "Select".
///
/// [useFilterBox] types [label] into the sheet's own "Filter" box first. The
/// category and tag catalogs are install-wide, so a shared dev backend can hold
/// enough of them to push a freshly-seeded one below the fold -- filtering the
/// grid to a single chip is deterministic where scrolling it is not (see
/// receipt_category_picker_scroll_test.dart for why the grid genuinely scrolls).
Future<void> pickFromMultiselect(
  WidgetTester tester,
  String label, {
  bool useFilterBox = false,
}) async {
  if (useFilterBox) {
    await tester.enterText(formField("filter"), label);
    await drainFrames(tester);
  }

  // Scoped to ChoiceChip: the sheet's options are ChoiceChips, while the
  // selections already made render as InputChips and the operation row above
  // uses ChoiceChips with operation labels ("Contains", "Between", ...).
  final chip = find.widgetWithText(ChoiceChip, label);
  await pumpUntilFound(tester, chip);
  await drainFrames(tester);
  await tester.ensureVisible(chip);
  // ensureVisible does not pump; without this the tap is computed against
  // pre-scroll geometry.
  await tester.pump(const Duration(milliseconds: 100));
  await tester.tap(chip);
  await drainFrames(tester);

  final select = find.widgetWithText(BottomSubmitButton, "Select");
  await pumpUntilFound(tester, select.hitTestable());
  await drainFrames(tester);
  await tester.tap(select);
  await drainFrames(tester);
}

/// The chip a [MultiSelectField] renders for an applied selection -- an
/// `InputChip`, not the `ChoiceChip` that was tapped in the sheet. Waiting on it
/// proves the sheet's result was written back before "Save condition" is poked.
Finder selectedChip(String label) => find.widgetWithText(InputChip, label);
