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
/// list's tests: the labelling, the muted zero row, and -- the one that WAS found by a
/// user -- that nothing the server sends is ever clipped out of sight.
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
      buildSummaryRow(status: api.ReceiptStatus.OPEN, receiptCount: 1, total: '22.24'),
      buildSummaryRow(status: api.ReceiptStatus.RESOLVED, receiptCount: 0, total: '0.00'),
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

    expect(find.text('1 receipt'), findsOneWidget);
  });

  // A configured status matching nothing still renders, so the block keeps its shape as
  // the filter narrows -- muted, so a legitimate zero does not read as a bug. Asserting
  // the COLOUR, not mere presence: presence passes against a tree with no mute at all.
  testWidgets('renders a zero row, muted', (tester) async {
    await pumpBar(tester, summary: threeAndTwo);

    final scheme = buildAppTheme().colorScheme;

    // Read off the row LABEL by name rather than the row's first Text: a row now holds
    // a label, a count and a figure per column, and the count is muted in EVERY row, so
    // grabbing whichever Text comes first would pass against a tree with no mute at all.
    Text labelOf(String rowKey, String text) => tester.widget<Text>(
      find.descendant(
        of: find.byKey(ValueKey('receipt-summary-row-$rowKey')),
        matching: find.text(text),
      ),
    );

    expect(
      find.descendant(
        of: find.byKey(const ValueKey('receipt-summary-row-RESOLVED')),
        matching: find.text('0 receipts'),
      ),
      findsOneWidget,
    );
    expect(
      labelOf('RESOLVED', 'Resolved Receipts').style?.color,
      scheme.onSurfaceVariant,
    );
    expect(labelOf('OPEN', 'Open Receipts').style?.color, scheme.onSurface);
  });

  testWidgets('renders a column per currency custom field, in the API order', (
    tester,
  ) async {
    await pumpBar(
      tester,
      summary: buildReceiptSummary(
        overall: buildSummaryRow(
          receiptCount: 2,
          total: '100.00',
          customFieldTotals: [
            buildSummaryCustomFieldTotal(customFieldId: 7, name: 'HST', total: '13.00'),
            buildSummaryCustomFieldTotal(
              customFieldId: 4,
              name: 'Subtotal',
              total: '87.00',
            ),
          ],
        ),
      ),
    );

    expect(find.text('HST'), findsOneWidget);
    expect(find.text('Subtotal'), findsOneWidget);
    expect(
      find.byKey(const ValueKey('receipt-summary-figure-overall-cf-7')),
      findsOneWidget,
    );
    expect(
      find.byKey(const ValueKey('receipt-summary-figure-overall-cf-4')),
      findsOneWidget,
    );

    // The figures are emitted in the API's id order. Asserted in READING order rather
    // than by x alone: figures wrap, so a later one can legitimately sit further left
    // on the following line -- which is exactly what a wide text scale produces.
    final hst = tester.getRect(find.text('HST'));
    final subtotal = tester.getRect(find.text('Subtotal'));
    expect(
      hst.top < subtotal.top || (hst.top == subtotal.top && hst.left < subtotal.left),
      isTrue,
      reason: 'HST (id 7) is sent first, so it must render first',
    );
  });

  /// The regression guard this widget exists for.
  ///
  /// The first version was a frozen label column beside fixed-width figure columns in a
  /// horizontal scroll view. At 390pt that leaves 210pt for figures and two columns
  /// wanted 232, so a group with a SINGLE currency field already had its second column
  /// chopped mid-number -- and the test that was meant to cover it only asserted no
  /// exception was thrown, which a clipped column passes happily.
  ///
  /// So this asserts the thing that actually matters: every figure the server sent is
  /// fully within the bar horizontally. Vertical overflow is deliberately NOT asserted
  /// -- the 35% cap scrolls, with a visible thumb, which is a real affordance; there is
  /// no such thing as scrolling a number back from beyond the right edge.
  testWidgets('clips nothing horizontally, however much is configured', (tester) async {
    List<api.ReceiptSummaryCustomFieldTotal> fields(String a, String b) => [
      buildSummaryCustomFieldTotal(customFieldId: 7, name: 'HST', total: a),
      buildSummaryCustomFieldTotal(customFieldId: 4, name: 'Subtotal', total: b),
    ];

    await pumpBar(
      tester,
      summary: buildReceiptSummary(
        overall: buildSummaryRow(
          receiptCount: 9,
          total: '1419.31',
          customFieldTotals: fields('184.22', '1235.09'),
        ),
        statuses: [
          for (final status in [
            api.ReceiptStatus.OPEN,
            // The longest label in the set, and the one whose two-line wrap used to
            // swallow its own receipt count through the label's ellipsis.
            api.ReceiptStatus.NEEDS_ATTENTION,
            api.ReceiptStatus.RESOLVED,
            api.ReceiptStatus.DRAFT,
            api.ReceiptStatus.DECLINED,
          ])
            buildSummaryRow(
              status: status,
              receiptCount: 1,
              total: '1.00',
              customFieldTotals: fields('0.13', '0.87'),
            ),
        ],
      ),
    );

    expect(tester.takeException(), isNull);
    expect(find.byType(ErrorWidget), findsNothing);

    final bar = tester.getRect(find.byKey(const ValueKey('receipt-summary')));

    for (final rowKey in [
      'overall',
      'OPEN',
      'NEEDS_ATTENTION',
      'RESOLVED',
      'DRAFT',
      'DECLINED',
    ]) {
      for (final figure in ['total', 'cf-7', 'cf-4']) {
        final finder = find.byKey(ValueKey('receipt-summary-figure-$rowKey-$figure'));
        expect(finder, findsOneWidget, reason: '$rowKey/$figure should be rendered');

        final rect = tester.getRect(finder);
        expect(
          rect.left,
          greaterThanOrEqualTo(bar.left),
          reason: '$rowKey/$figure starts left of the bar',
        );
        expect(
          rect.right,
          lessThanOrEqualTo(bar.right),
          reason: '$rowKey/$figure runs past the right edge and is clipped',
        );
      }

      // And the count survives beside even the longest label, rather than being
      // ellipsised away with it.
      expect(
        find.descendant(
          of: find.byKey(ValueKey('receipt-summary-row-$rowKey')),
          matching: find.textContaining('receipt'),
        ),
        findsWidgets,
        reason: '$rowKey lost its receipt count',
      );
    }
  });

  // A chip row with one option is not a choice; naming the group says the same thing
  // without the false affordance.
  testWidgets('names the sole configuration group instead of offering a chip', (
    tester,
  ) async {
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
    final selected = tester.widget<ChoiceChip>(
      find.byKey(const ValueKey('receipt-summary-config-group-2')),
    );
    expect(selected.selected, isTrue);

    await tester.tap(find.byKey(const ValueKey('receipt-summary-config-group-3')));
    await tester.pumpAndSettle();
    expect(picked, 3);
  });

  // The fixed-width label column is the shape that already overflowed once, in
  // ListItemTrailingStatus. All five statuses at phone width is the worst case.
  // The divider goes on whichever edge faces the list, so the bar reads as attached to
  // it rather than floating.
  testWidgets('the divider follows the position', (tester) async {
    await pumpBar(tester, summary: threeAndTwo, atTop: true);
    var decoration =
        tester.widget<Container>(find.byKey(const ValueKey('receipt-summary'))).decoration
            as BoxDecoration;
    expect(decoration.border?.bottom.style, BorderStyle.solid);
    expect(decoration.border?.top.style, BorderStyle.none);

    await pumpBar(tester, summary: threeAndTwo, atTop: false);
    decoration =
        tester.widget<Container>(find.byKey(const ValueKey('receipt-summary'))).decoration
            as BoxDecoration;
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
