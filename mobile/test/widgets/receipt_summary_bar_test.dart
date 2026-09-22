import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/receipt_summary_bar.dart';
import 'package:receipt_wrangler_mobile/models/system_settings_model.dart';
import 'package:receipt_wrangler_mobile/theme/app_theme.dart';

import '../helpers/receipt_form_test_helpers.dart';
import '../helpers/receipt_summary_test_helpers.dart';
import '../helpers/widget_test_helpers.dart';

/// The bar is presentational, so everything worth pinning is here rather than in the
/// list's tests: the labelling, the muted zero row, and -- the one that would otherwise
/// be found by a user -- that the figure columns stay locked to their headings.
void main() {
  setUpAll(registerCustomCurrencyForTests);

  Future<void> pumpBar(
    WidgetTester tester, {
    required api.ReceiptSummary summary,
    bool atTop = false,
    List<api.Group> configGroups = const [],
    int? selectedConfigGroupId,
    void Function(int)? onConfigGroupSelected,
    Size surface = const Size(390, 760),
  }) async {
    await tester.binding.setSurfaceSize(surface);
    addTearDown(() => tester.binding.setSurfaceSize(null));

    await tester.pumpWidget(
      ChangeNotifierProvider<SystemSettingsModel>(
        create: (_) => SystemSettingsModel(),
        child: MaterialApp(
          theme: buildAppTheme(),
          home: Scaffold(
            body: Column(
              children: [
                ReceiptSummaryBar(
                  summary: summary,
                  atTop: atTop,
                  configGroups: configGroups,
                  selectedConfigGroupId: selectedConfigGroupId,
                  onConfigGroupSelected: onConfigGroupSelected,
                ),
                const Expanded(child: SizedBox.shrink()),
              ],
            ),
          ),
        ),
      ),
    );
    await tester.pumpAndSettle();
  }

  final threeAndTwo = buildReceiptSummary(
    overall: buildSummaryRow(receiptCount: 3, total: '122.24'),
    statuses: [
      buildSummaryRow(
          status: api.ReceiptStatus.OPEN, receiptCount: 1, total: '22.24'),
      buildSummaryRow(
          status: api.ReceiptStatus.RESOLVED, receiptCount: 0, total: '0.00'),
    ],
  );

  testWidgets('renders the overall row with its count and total', (tester) async {
    await pumpBar(tester, summary: threeAndTwo);

    expect(find.byKey(const ValueKey('receipt-summary')), findsOneWidget);
    expect(find.byKey(const ValueKey('receipt-summary-row-overall')), findsOneWidget);
    expect(find.textContaining('All Receipts'), findsOneWidget);
    expect(find.textContaining('3 receipts'), findsOneWidget);
    // Matched on the digits, not the whole string: the symbol's side is a global
    // System Setting, and the test currency registers its own pattern.
    expect(find.textContaining('122.24'), findsOneWidget);
  });

  testWidgets('singularizes a one-receipt row', (tester) async {
    await pumpBar(tester, summary: threeAndTwo);

    expect(find.textContaining('(1 receipt)'), findsOneWidget);
  });

  // A configured status matching nothing still renders, so the block keeps its shape as
  // the filter narrows -- muted, so a legitimate zero does not read as a bug. Asserting
  // the COLOUR, not mere presence: presence passes against a tree with no mute at all.
  testWidgets('renders a zero row, muted', (tester) async {
    await pumpBar(tester, summary: threeAndTwo);

    final scheme = buildAppTheme().colorScheme;
    final zero = tester.widget<Text>(find.descendant(
      of: find.byKey(const ValueKey('receipt-summary-row-RESOLVED')),
      matching: find.byType(Text),
    ));
    final nonZero = tester.widget<Text>(find.descendant(
      of: find.byKey(const ValueKey('receipt-summary-row-OPEN')),
      matching: find.byType(Text),
    ));

    expect(zero.data, contains('0 receipts'));
    expect(zero.style?.color, scheme.onSurfaceVariant);
    expect(nonZero.style?.color, scheme.onSurface);
  });

  testWidgets('renders a column per currency custom field, in the API order',
      (tester) async {
    await pumpBar(
      tester,
      summary: buildReceiptSummary(
        overall: buildSummaryRow(
          receiptCount: 2,
          total: '100.00',
          customFieldTotals: [
            buildSummaryCustomFieldTotal(customFieldId: 7, name: 'HST', total: '13.00'),
            buildSummaryCustomFieldTotal(
                customFieldId: 4, name: 'Subtotal', total: '87.00'),
          ],
        ),
      ),
    );

    expect(find.text('HST'), findsOneWidget);
    expect(find.text('Subtotal'), findsOneWidget);
    expect(find.byKey(const ValueKey('receipt-summary-figure-overall-cf-7')),
        findsOneWidget);
    expect(find.byKey(const ValueKey('receipt-summary-figure-overall-cf-4')),
        findsOneWidget);

    // The headings are read off the overall row, so the id order is the API's.
    final hst = tester.getRect(find.text('HST'));
    final subtotal = tester.getRect(find.text('Subtotal'));
    expect(hst.left, lessThan(subtotal.left));
  });

  /// The regression guard for the layout decision: ONE horizontal scroll view for the
  /// whole grid. Per-row scroll views would let the rows desync under a drag and park a
  /// value under the wrong heading, and this case fails against that tree.
  testWidgets('every row scrolls together, and the labels do not move', (tester) async {
    final wide = buildReceiptSummary(
      overall: buildSummaryRow(
        receiptCount: 2,
        total: '100.00',
        customFieldTotals: [
          buildSummaryCustomFieldTotal(customFieldId: 7, name: 'HST', total: '13.00'),
          buildSummaryCustomFieldTotal(customFieldId: 4, name: 'Sub', total: '87.00'),
        ],
      ),
      statuses: [
        buildSummaryRow(
          status: api.ReceiptStatus.OPEN,
          receiptCount: 1,
          total: '40.00',
          customFieldTotals: [
            buildSummaryCustomFieldTotal(customFieldId: 7, name: 'HST', total: '5.00'),
            buildSummaryCustomFieldTotal(customFieldId: 4, name: 'Sub', total: '35.00'),
          ],
        ),
      ],
    );
    await pumpBar(tester, summary: wide);

    final labelBefore = tester.getRect(
        find.byKey(const ValueKey('receipt-summary-row-overall')));

    await tester.drag(
      find.byKey(const ValueKey('receipt-summary-figures-scroll')),
      const Offset(-120, 0),
    );
    await tester.pumpAndSettle();

    final overallCf = tester
        .getRect(find.byKey(const ValueKey('receipt-summary-figure-overall-cf-4')));
    final openCf = tester
        .getRect(find.byKey(const ValueKey('receipt-summary-figure-OPEN-cf-4')));
    expect(overallCf.left, openCf.left,
        reason: 'the rows share one scroll view, so a column cannot desync');

    expect(
      tester.getRect(find.byKey(const ValueKey('receipt-summary-row-overall'))),
      labelBefore,
      reason: 'the label column is outside the horizontal scroll view',
    );
  });

  // A chip row with one option is not a choice; naming the group says the same thing
  // without the false affordance.
  testWidgets('names the sole configuration group instead of offering a chip',
      (tester) async {
    await pumpBar(
      tester,
      summary: threeAndTwo,
      configGroups: [buildGroup(id: 2, name: 'Household')],
    );

    expect(find.byKey(const ValueKey('receipt-summary-config-note')), findsOneWidget);
    expect(find.textContaining('Household'), findsOneWidget);
    expect(find.byKey(const ValueKey('receipt-summary-config-chips')), findsNothing);
  });

  testWidgets('offers chips for several groups and reports the pick', (tester) async {
    int? picked;
    await pumpBar(
      tester,
      summary: threeAndTwo,
      configGroups: [
        buildGroup(id: 2, name: 'Household'),
        buildGroup(id: 3, name: 'Office'),
      ],
      selectedConfigGroupId: 2,
      onConfigGroupSelected: (id) => picked = id,
    );

    expect(find.byKey(const ValueKey('receipt-summary-config-chips')), findsOneWidget);
    final selected = tester
        .widget<ChoiceChip>(find.byKey(const ValueKey('receipt-summary-config-group-2')));
    expect(selected.selected, isTrue);

    await tester.tap(find.byKey(const ValueKey('receipt-summary-config-group-3')));
    await tester.pumpAndSettle();
    expect(picked, 3);
  });

  // The fixed-width label column is the shape that already overflowed once, in
  // ListItemTrailingStatus. All five statuses at phone width is the worst case.
  testWidgets('every configured status fits at phone width', (tester) async {
    await pumpBar(
      tester,
      summary: buildReceiptSummary(
        overall: buildSummaryRow(receiptCount: 9, total: '900.00'),
        statuses: [
          for (final status in [
            api.ReceiptStatus.OPEN,
            api.ReceiptStatus.NEEDS_ATTENTION,
            api.ReceiptStatus.RESOLVED,
            api.ReceiptStatus.DRAFT,
            api.ReceiptStatus.DECLINED,
          ])
            buildSummaryRow(status: status, receiptCount: 1, total: '1.00'),
        ],
      ),
    );

    expect(tester.takeException(), isNull);
    expect(find.byType(ErrorWidget), findsNothing);
  });

  // The divider goes on whichever edge faces the list, so the bar reads as attached to
  // it rather than floating.
  testWidgets('the divider follows the position', (tester) async {
    await pumpBar(tester, summary: threeAndTwo, atTop: true);
    var decoration = tester
        .widget<Container>(find.byKey(const ValueKey('receipt-summary')))
        .decoration as BoxDecoration;
    expect(decoration.border?.bottom.style, BorderStyle.solid);
    expect(decoration.border?.top.style, BorderStyle.none);

    await pumpBar(tester, summary: threeAndTwo, atTop: false);
    decoration = tester
        .widget<Container>(find.byKey(const ValueKey('receipt-summary')))
        .decoration as BoxDecoration;
    expect(decoration.border?.top.style, BorderStyle.solid);
    expect(decoration.border?.bottom.style, BorderStyle.none);
  });

  // formatCurrency parses with double.parse, which THROWS -- and from inside this bar
  // that takes down the whole receipts screen, not just the block.
  testWidgets('an unparseable total falls back to the raw text', (tester) async {
    await pumpBar(
      tester,
      summary: buildReceiptSummary(
        overall: buildSummaryRow(receiptCount: 1, total: 'not-a-number'),
      ),
    );

    expect(tester.takeException(), isNull);
    expect(find.text('not-a-number'), findsOneWidget);
  });
}
