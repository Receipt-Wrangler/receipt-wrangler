import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/constants/receipt_filter_fields.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_filter_condition_editor.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/amount_field.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/multi-select-field.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

import '../helpers/receipt_filter_test_helpers.dart';
import '../helpers/receipt_filter_widget_helpers.dart';
import '../helpers/widget_test_helpers.dart';

/// The editor is where a condition's shape is decided, so these cases pin the
/// operation-to-input mapping, the save gate, and the three outcomes the filter
/// screen has to tell apart (saved / removed / dismissed).
void main() {
  setUpAll(registerCustomCurrencyForTests);

  const saveKey = ValueKey("receipt-filter-condition-save");
  const removeKey = ValueKey("receipt-filter-condition-remove");

  ReceiptFilterField field(String key) => receiptFilterFieldByKey(key)!;

  Future<void> pumpEditor(
    WidgetTester tester, {
    required String fieldKey,
    ReceiptFilterCondition? existing,
    String groupId = "${ReceiptFilterHarness.householdId}",
  }) async {
    await pumpWithFilterHarness(
      tester,
      buildReceiptFilterHarness(),
      Scaffold(
        body: ReceiptFilterConditionEditor(
          field: field(fieldKey),
          groupId: groupId,
          existing: existing,
        ),
      ),
    );
  }

  bool saveEnabled(WidgetTester tester) =>
      tester.widget<MaterialButton>(find.byKey(saveKey)).onPressed != null;

  Finder operationChip(api.FilterOperation operation) =>
      find.byKey(ValueKey("receipt-filter-operation-${operation.name}"));

  group("the operations offered", () {
    test("come from the shared table, so the editor cannot invent one", () {
      // Pinned exhaustively in receipt_filter_fields_test; this just records
      // that the editor is not a second source.
      expect(operationsForFilterField(field("date")).length, 5);
      expect(operationsForFilterField(field("name")).length, 2);
      expect(operationsForFilterField(field("tags")).length, 1);
    });

    testWidgets("a date field offers all five", (tester) async {
      await pumpEditor(tester, fieldKey: "date");

      for (final operation in operationsForFilterField(field("date"))) {
        expect(operationChip(operation), findsOneWidget);
      }
    });

    testWidgets("a text field offers Contains and Equals", (tester) async {
      await pumpEditor(tester, fieldKey: "name");

      expect(operationChip(api.FilterOperation.CONTAINS), findsOneWidget);
      expect(operationChip(api.FilterOperation.EQUALS), findsOneWidget);
      expect(operationChip(api.FilterOperation.BETWEEN), findsNothing);
    });

    testWidgets("a list field still renders its single Contains chip",
        (tester) async {
      // Rendering it keeps every editor the same shape and says what the
      // condition means; hiding it makes list fields look half-built.
      await pumpEditor(tester, fieldKey: "tags");

      expect(operationChip(api.FilterOperation.CONTAINS), findsOneWidget);
    });
  });

  group("the value input follows the operation", () {
    testWidgets("a single amount becomes a pair on BETWEEN, and back",
        (tester) async {
      await pumpEditor(tester, fieldKey: "amount");
      expect(find.byType(AmountField), findsOneWidget);

      await tester.tap(operationChip(api.FilterOperation.BETWEEN));
      await tester.pump();
      expect(find.byType(AmountField), findsNWidgets(2));
      expect(find.text("and"), findsOneWidget);

      await tester.tap(operationChip(api.FilterOperation.EQUALS));
      await tester.pump();
      expect(find.byType(AmountField), findsOneWidget);
    });

    testWidgets("a date swaps its single picker for a range picker",
        (tester) async {
      await pumpEditor(tester, fieldKey: "date");
      expect(find.byKey(const ValueKey("receipt-filter-date")), findsOneWidget);

      await tester.tap(operationChip(api.FilterOperation.BETWEEN));
      await tester.pump();

      expect(find.byKey(const ValueKey("receipt-filter-date")), findsNothing);
      expect(find.byKey(const ValueKey("receipt-filter-date-range")),
          findsOneWidget);
    });

    testWidgets("WITHIN_CURRENT_MONTH renders no input at all", (tester) async {
      await pumpEditor(tester, fieldKey: "date");

      await tester.tap(operationChip(api.FilterOperation.WITHIN_CURRENT_MONTH));
      await tester.pump();

      expect(find.byKey(const ValueKey("receipt-filter-date")), findsNothing);
      expect(find.byKey(const ValueKey("receipt-filter-date-range")),
          findsNothing);
      expect(
          find.byKey(
              const ValueKey("receipt-filter-within-current-month-note")),
          findsOneWidget);
    });

    testWidgets("text renders one field", (tester) async {
      await pumpEditor(tester, fieldKey: "name");

      expect(find.byType(FormBuilderTextField), findsOneWidget);
    });

    testWidgets("a list field renders the shared chip field", (tester) async {
      await pumpEditor(tester, fieldKey: "status");

      expect(find.byType(MultiSelectField<api.ReceiptStatus>), findsOneWidget);
    });

    testWidgets("the amount input is the shared currency field, not required",
        (tester) async {
      // An amount condition is only authored when the user asks for one, so a
      // required validator would block the editor before anything is typed.
      await pumpEditor(tester, fieldKey: "amount");

      final amountField = tester.widget<AmountField>(find.byType(AmountField));
      expect(amountField.validator, isNotNull);
      expect(amountField.validator!(""), isNull);
    });
  });

  group("the save gate", () {
    testWidgets("an untouched text condition cannot be saved", (tester) async {
      await pumpEditor(tester, fieldKey: "name");

      expect(saveEnabled(tester), isFalse);
    });

    testWidgets("typing enables it", (tester) async {
      await pumpEditor(tester, fieldKey: "name");

      await tester.enterText(find.byType(FormBuilderTextField), "Costco");
      await tester.pump();

      expect(saveEnabled(tester), isTrue);
    });

    testWidgets("whitespace alone does not", (tester) async {
      await pumpEditor(tester, fieldKey: "name");

      await tester.enterText(find.byType(FormBuilderTextField), "   ");
      await tester.pump();

      expect(saveEnabled(tester), isFalse);
    });

    testWidgets("an untouched amount cannot be saved", (tester) async {
      // AmountField seeds an EMPTY box at zero, and its transformer maps empty
      // to "0.00" -- so emptiness has to come from the raw value, or adding an
      // Amount condition would silently mean `amount = 0`.
      await pumpEditor(tester, fieldKey: "amount");

      expect(saveEnabled(tester), isFalse);
    });

    testWidgets("an empty selection cannot be saved", (tester) async {
      await pumpEditor(tester, fieldKey: "status");

      expect(saveEnabled(tester), isFalse);
    });

    testWidgets("WITHIN_CURRENT_MONTH is saveable with no value at all",
        (tester) async {
      await pumpEditor(tester, fieldKey: "date");

      await tester.tap(operationChip(api.FilterOperation.WITHIN_CURRENT_MONTH));
      await tester.pump();

      expect(saveEnabled(tester), isTrue);
    });

    testWidgets("an existing condition opens already saveable", (tester) async {
      await pumpEditor(
        tester,
        fieldKey: "name",
        existing: const ReceiptFilterCondition(
            operation: api.FilterOperation.EQUALS, value: "Costco"),
      );

      expect(saveEnabled(tester), isTrue);
    });
  });

  group("seeding from an existing condition", () {
    testWidgets("pre-selects its operation", (tester) async {
      await pumpEditor(
        tester,
        fieldKey: "name",
        existing: const ReceiptFilterCondition(
            operation: api.FilterOperation.EQUALS, value: "Costco"),
      );

      final chip = tester.widget<ChoiceChip>(
          operationChip(api.FilterOperation.EQUALS));
      expect(chip.selected, isTrue);

      final other = tester.widget<ChoiceChip>(
          operationChip(api.FilterOperation.CONTAINS));
      expect(other.selected, isFalse);
    });

    testWidgets("pre-fills its value", (tester) async {
      await pumpEditor(
        tester,
        fieldKey: "name",
        existing: const ReceiptFilterCondition(
            operation: api.FilterOperation.CONTAINS, value: "Trader Joe's"),
      );

      expect(find.text("Trader Joe's"), findsOneWidget);
    });

    testWidgets("pre-fills a selection as chips", (tester) async {
      await pumpEditor(
        tester,
        fieldKey: "status",
        existing: const ReceiptFilterCondition(
            operation: api.FilterOperation.CONTAINS,
            value: [api.ReceiptStatus.RESOLVED]),
      );

      expect(find.text("Resolved"), findsOneWidget);
    });
  });

  group("Remove", () {
    testWidgets("is absent when the field is being added", (tester) async {
      await pumpEditor(tester, fieldKey: "name");

      expect(find.byKey(removeKey), findsNothing);
    });

    testWidgets("is present when an existing condition is being edited",
        (tester) async {
      await pumpEditor(
        tester,
        fieldKey: "name",
        existing: const ReceiptFilterCondition(
            operation: api.FilterOperation.CONTAINS, value: "Costco"),
      );

      expect(find.byKey(removeKey), findsOneWidget);
    });
  });

  group("what the editor hands back", () {
    /// Pumps the editor behind a button that opens it as a route, so the popped
    /// result is observable the way the filter screen sees it.
    Future<ReceiptFilterEditorResult?> runEditor(
      WidgetTester tester, {
      required String fieldKey,
      ReceiptFilterCondition? existing,
      required Future<void> Function(WidgetTester tester) act,
    }) async {
      ReceiptFilterEditorResult? captured;
      var completed = false;

      await pumpWithFilterHarness(
        tester,
        buildReceiptFilterHarness(),
        Scaffold(
          body: Builder(
            builder: (context) => ElevatedButton(
              onPressed: () async {
                captured = await Navigator.of(context)
                    .push<ReceiptFilterEditorResult>(MaterialPageRoute(
                  builder: (_) => Scaffold(
                    body: ReceiptFilterConditionEditor(
                      field: field(fieldKey),
                      groupId: "${ReceiptFilterHarness.householdId}",
                      existing: existing,
                    ),
                  ),
                ));
                completed = true;
              },
              child: const Text("open"),
            ),
          ),
        ),
      );

      await tester.tap(find.text("open"));
      await tester.pumpAndSettle();

      await act(tester);
      await tester.pumpAndSettle();

      expect(completed, isTrue, reason: "the editor never popped");
      return captured;
    }

    testWidgets("Save returns the authored condition", (tester) async {
      final result = await runEditor(
        tester,
        fieldKey: "name",
        act: (tester) async {
          await tester.enterText(find.byType(FormBuilderTextField), "Costco");
          await tester.pump();
          await tester.tap(find.byKey(saveKey));
        },
      );

      expect(result, isNotNull);
      expect(result!.removed, isFalse);
      expect(result.condition!.operation, api.FilterOperation.CONTAINS);
      expect(result.condition!.value, "Costco");
    });

    testWidgets("Remove says so rather than returning a condition",
        (tester) async {
      final result = await runEditor(
        tester,
        fieldKey: "name",
        existing: const ReceiptFilterCondition(
            operation: api.FilterOperation.CONTAINS, value: "Costco"),
        act: (tester) async => tester.tap(find.byKey(removeKey)),
      );

      expect(result!.removed, isTrue);
      expect(result.condition, isNull);
    });

    testWidgets("backing out returns nothing, which means no change",
        (tester) async {
      // Distinct from Remove: the draft must be left exactly as it was.
      final result = await runEditor(
        tester,
        fieldKey: "name",
        existing: const ReceiptFilterCondition(
            operation: api.FilterOperation.CONTAINS, value: "Costco"),
        act: (tester) async {
          final NavigatorState navigator =
              tester.state(find.byType(Navigator).last);
          navigator.pop();
        },
      );

      expect(result, isNull);
    });

    testWidgets("the chosen operation rides along", (tester) async {
      final result = await runEditor(
        tester,
        fieldKey: "name",
        act: (tester) async {
          await tester.tap(operationChip(api.FilterOperation.EQUALS));
          await tester.pump();
          await tester.enterText(find.byType(FormBuilderTextField), "Costco");
          await tester.pump();
          await tester.tap(find.byKey(saveKey));
        },
      );

      expect(result!.condition!.operation, api.FilterOperation.EQUALS);
    });
  });

  group("switching operations", () {
    testWidgets("drops a value whose shape no longer fits", (tester) async {
      // A lone amount carried into BETWEEN would be encoded as a range.
      await pumpEditor(
        tester,
        fieldKey: "amount",
        existing: const ReceiptFilterCondition(
            operation: api.FilterOperation.EQUALS, value: 42.0),
      );
      expect(saveEnabled(tester), isTrue);

      await tester.tap(operationChip(api.FilterOperation.BETWEEN));
      await tester.pump();

      expect(saveEnabled(tester), isFalse);
    });

    testWidgets("keeps a value the new operation can still use",
        (tester) async {
      await pumpEditor(
        tester,
        fieldKey: "name",
        existing: const ReceiptFilterCondition(
            operation: api.FilterOperation.CONTAINS, value: "Costco"),
      );

      await tester.tap(operationChip(api.FilterOperation.EQUALS));
      await tester.pump();

      expect(find.text("Costco"), findsOneWidget);
      expect(saveEnabled(tester), isTrue);
    });

    testWidgets("a selection survives, since both operations are CONTAINS",
        (tester) async {
      await pumpEditor(
        tester,
        fieldKey: "status",
        existing: const ReceiptFilterCondition(
            operation: api.FilterOperation.CONTAINS,
            value: [api.ReceiptStatus.OPEN]),
      );

      expect(saveEnabled(tester), isTrue);
      expect(find.text("Open"), findsOneWidget);
    });
  });

  testWidgets("a selection made in the picker reaches the saved condition",
      (tester) async {
    // The chip field is the shared MultiSelectField, which does not write its
    // own value -- the editor is the source of truth and must write it back.
    await pumpEditor(
      tester,
      fieldKey: "status",
      existing: const ReceiptFilterCondition(
          operation: api.FilterOperation.CONTAINS,
          value: [api.ReceiptStatus.OPEN, api.ReceiptStatus.DRAFT]),
    );

    final multiSelect = tester.widget<MultiSelectField<api.ReceiptStatus>>(
        find.byType(MultiSelectField<api.ReceiptStatus>));

    // Removing a chip goes out through onRemove with the remaining list.
    multiSelect.onRemove!([api.ReceiptStatus.DRAFT]);
    await tester.pump();

    expect(find.text("Open"), findsNothing);
    expect(find.text("Draft"), findsOneWidget);
  });

  testWidgets("the categories editor uses the group's own catalog",
      (tester) async {
    // Sourced from CategoryModel.categoriesForGroup, never the flat admin-only
    // list -- a normal user would see an empty picker otherwise.
    await pumpWithFilterHarness(
      tester,
      buildReceiptFilterHarness(
          householdCategories: [buildCategory(1, "Groceries")]),
      Scaffold(
        body: ReceiptFilterConditionEditor(
          field: field("categories"),
          groupId: "${ReceiptFilterHarness.householdId}",
          existing: ReceiptFilterCondition(
              operation: api.FilterOperation.CONTAINS,
              value: [buildCategory(1, "Groceries")]),
        ),
      ),
    );

    expect(find.text("Groceries"), findsOneWidget);
    expect(saveEnabled(tester), isTrue);
  });
}
