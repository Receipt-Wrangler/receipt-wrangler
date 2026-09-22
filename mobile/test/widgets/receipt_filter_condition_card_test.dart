import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/constants/receipt_filter_fields.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_filter_condition_card.dart';
import 'package:receipt_wrangler_mobile/theme/app_theme.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

import '../helpers/widget_test_helpers.dart';

void main() {
  setUpAll(registerCustomCurrencyForTests);

  final amountField = receiptFilterFieldByKey("amount")!;
  const amountCondition = ReceiptFilterCondition(
      operation: api.FilterOperation.GREATER_THAN, value: 25.0);

  Future<_Counts> pumpCard(
    WidgetTester tester, {
    ReceiptFilterField? field,
    ReceiptFilterCondition condition = amountCondition,
  }) async {
    final counts = _Counts();

    await tester.pumpWidget(MaterialApp(
      theme: buildAppTheme(),
      home: Scaffold(
        body: ReceiptFilterConditionCard(
          field: field ?? amountField,
          condition: condition,
          onTap: () => counts.taps++,
          onRemove: () => counts.removes++,
        ),
      ),
    ));
    await tester.pump();

    return counts;
  }

  testWidgets("reads its field, operation and value", (tester) async {
    await pumpCard(tester);

    expect(find.text("Amount"), findsOneWidget);
    expect(find.text("Greater than"), findsOneWidget);
    expect(
        find.text(receiptFilterValueLabel(amountField, amountCondition)),
        findsOneWidget);
  });

  testWidgets("carries no field icon -- those belong to the add sheet",
      (tester) async {
    await pumpCard(tester);

    expect(find.byIcon(amountField.icon), findsNothing);
    // The affordances it does carry.
    expect(find.byIcon(Icons.chevron_right), findsOneWidget);
    expect(find.byIcon(Icons.close), findsOneWidget);
  });

  testWidgets("tapping the body reopens the condition", (tester) async {
    final counts = await pumpCard(tester);

    await tester.tap(find.text("Amount"));
    await tester.pump();

    expect(counts.taps, 1);
    expect(counts.removes, 0);
  });

  testWidgets("the remove X removes WITHOUT also reopening", (tester) async {
    // The whole card is one InkWell now, so the X sits inside the tap target
    // it must beat. If it loses the gesture arena, removing a condition also
    // opens the editor for the field that was just dropped.
    final counts = await pumpCard(tester);

    await tester
        .tap(find.byKey(const ValueKey("receipt-filter-remove-amount")));
    await tester.pump();

    expect(counts.removes, 1);
    expect(counts.taps, 0);
  });

  testWidgets("a long value is ellipsised rather than wrapping the card",
      (tester) async {
    final categories = receiptFilterFieldByKey("categories")!;
    await pumpCard(
      tester,
      field: categories,
      condition: ReceiptFilterCondition(
        operation: api.FilterOperation.CONTAINS,
        value: List.generate(
            12,
            (i) => (api.CategoryBuilder()
                  ..id = i
                  ..name = "A rather long category name $i")
                .build()),
      ),
    );

    final value = tester.widget<Text>(find.textContaining("A rather long"));
    expect(value.overflow, TextOverflow.ellipsis);
    expect(tester.takeException(), isNull);
  });
}

class _Counts {
  int taps = 0;
  int removes = 0;
}
