import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/add_receipt_filter_sheet.dart';

import '../helpers/receipt_filter_widget_helpers.dart';
import '../helpers/widget_test_helpers.dart';

void main() {
  setUpAll(registerCustomCurrencyForTests);

  group("unusedReceiptFilterFields", () {
    test("offers every field on a fresh filter", () {
      expect(unusedReceiptFilterFields({}, showGroupField: true).length, 10);
    });

    test("drops the fields already used", () {
      final keys = unusedReceiptFilterFields({"name", "amount"},
              showGroupField: true)
          .map((field) => field.key)
          .toList();

      expect(keys, isNot(contains("name")));
      expect(keys, isNot(contains("amount")));
      expect(keys.length, 8);
    });

    test("withholds Group unless the view spans groups", () {
      final withoutGroup =
          unusedReceiptFilterFields({}, showGroupField: false)
              .map((field) => field.key);

      expect(withoutGroup, isNot(contains("group")));
      expect(withoutGroup.length, 9);
    });

    test("keeps the shared table's order", () {
      expect(
          unusedReceiptFilterFields({"date"}, showGroupField: true)
              .take(3)
              .map((field) => field.key)
              .toList(),
          ["name", "paidBy", "group"]);
    });
  });

  group("the list", () {
    testWidgets("renders a row per field, with its hint", (tester) async {
      await pumpWithFilterHarness(
        tester,
        buildReceiptFilterHarness(),
        Scaffold(
          body: AddReceiptFilterList(
              fields: unusedReceiptFilterFields({}, showGroupField: false)),
        ),
      );

      expect(find.byKey(const ValueKey("add-receipt-filter-name")),
          findsOneWidget);
      expect(find.text("Merchant or receipt name"), findsOneWidget);
      expect(find.byKey(const ValueKey("add-receipt-filter-group")),
          findsNothing);
    });

    testWidgets("picking a row returns its key", (tester) async {
      String? picked;

      await pumpWithFilterHarness(
        tester,
        buildReceiptFilterHarness(),
        Builder(
          builder: (context) => Scaffold(
            body: ElevatedButton(
              onPressed: () async {
                picked = await Navigator.of(context).push<String>(
                  MaterialPageRoute(
                    builder: (_) => Scaffold(
                      body: AddReceiptFilterList(
                          fields: unusedReceiptFilterFields({},
                              showGroupField: false)),
                    ),
                  ),
                );
              },
              child: const Text("open"),
            ),
          ),
        ),
      );

      await tester.tap(find.text("open"));
      await tester.pumpAndSettle();
      await tester.tap(find.byKey(const ValueKey("add-receipt-filter-tags")));
      await tester.pumpAndSettle();

      expect(picked, "tags");
    });

    testWidgets("says so when every field is already filtered",
        (tester) async {
      await pumpWithFilterHarness(
        tester,
        buildReceiptFilterHarness(),
        const Scaffold(body: AddReceiptFilterList(fields: [])),
      );

      expect(find.byKey(const ValueKey("add-receipt-filter-empty")),
          findsOneWidget);
    });
  });
}
