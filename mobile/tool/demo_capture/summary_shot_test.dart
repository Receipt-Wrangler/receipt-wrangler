import 'dart:io';

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:image/image.dart' as img;
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/receipt_summary_bar.dart';
import 'package:receipt_wrangler_mobile/models/system_settings_model.dart';
import 'package:receipt_wrangler_mobile/theme/app_theme.dart';

import '../../test/helpers/widget_test_helpers.dart';
import '../../test/helpers/receipt_summary_test_helpers.dart';
import 'capture.dart';

/// Captures a still of `ReceiptSummaryBar` at real phone width, for judging how it
/// LOOKS rather than how it behaves.
///
/// Not part of the test suite — `flutter test` with no arguments scans only `test/`,
/// so CI never runs this. Driven by `tool/record_summary_shot.sh`, which invokes it
/// once per revision being compared and stitches the panels.
///
/// A still, not the GIF: `writeSideBySideGif` quantises to an octree palette, which
/// bands badly on flat UI — the committed placement GIF shows stripes through the
/// status washes that are encoding artifacts rather than pixels the widget drew. A
/// PNG is what the widget actually renders.
///
/// `$SUMMARY_SHOT_OUT` is the output path, `$SUMMARY_SHOT_LABEL` the caption text, and
/// `$SUMMARY_SHOT_VARIANT` picks `typical` or `heavy` (see below). A caption of exactly
/// `BEFORE` is tinted red, for the one-off case of comparing against an older revision;
/// anything else is green.
void main() {
  const beforeColor = Color(0xFFB4232B);
  const afterColor = Color(0xFF15803D);

  /// Two currency fields, because one is already enough to clip the old grid at
  /// 390pt and two is what a group that actually uses the feature configures.
  List<api.ReceiptSummaryCustomFieldTotal> fields(String hst, String tip) => [
        buildSummaryCustomFieldTotal(customFieldId: 1, name: 'HST', total: hst),
        buildSummaryCustomFieldTotal(customFieldId: 2, name: 'Tip', total: tip),
      ];

  /// What most groups configure: a couple of statuses and no currency fields at all.
  /// Worth capturing separately, because it is the shape the block usually has and the
  /// one the heavy case below can hide.
  final typical = buildReceiptSummary(
    overall: buildSummaryRow(receiptCount: 8, total: '1419.31'),
    statuses: [
      buildSummaryRow(
          status: api.ReceiptStatus.OPEN, receiptCount: 5, total: '685.86'),
      buildSummaryRow(
          status: api.ReceiptStatus.RESOLVED, receiptCount: 3, total: '733.45'),
    ],
  );

  /// The overall row plus four broken-out statuses, one of which matched nothing —
  /// so the still shows the muted zero row too.
  final heavy = buildReceiptSummary(
    overall: buildSummaryRow(
      receiptCount: 8,
      total: '1419.31',
      customFieldTotals: fields('184.22', '96.40'),
    ),
    statuses: [
      buildSummaryRow(
        status: api.ReceiptStatus.OPEN,
        receiptCount: 4,
        total: '654.61',
        customFieldTotals: fields('85.10', '41.25'),
      ),
      buildSummaryRow(
        status: api.ReceiptStatus.RESOLVED,
        receiptCount: 3,
        total: '733.45',
        customFieldTotals: fields('95.32', '55.15'),
      ),
      buildSummaryRow(
        status: api.ReceiptStatus.NEEDS_ATTENTION,
        receiptCount: 1,
        total: '31.25',
        customFieldTotals: fields('3.80', '0.00'),
      ),
      buildSummaryRow(
        status: api.ReceiptStatus.DECLINED,
        receiptCount: 0,
        total: '0.00',
        customFieldTotals: fields('0.00', '0.00'),
      ),
    ],
  );

  testWidgets('receipt summary bar at phone width', (tester) async {
    final out = Platform.environment['SUMMARY_SHOT_OUT'] ?? '/tmp/summary.png';
    final variant = Platform.environment['SUMMARY_SHOT_VARIANT'] ?? 'heavy';
    final label = Platform.environment['SUMMARY_SHOT_LABEL'] ??
        (variant == 'typical'
            ? 'TYPICAL — two statuses, no currency fields'
            : 'HEAVY — every status, two currency fields');
    final isBefore = label == 'BEFORE';

    final summary = variant == 'typical' ? typical : heavy;

    addTearDown(tester.view.reset);
    registerCustomCurrencyForTests();
    await loadDemoFontsOrFail(tester);

    tester.view.physicalSize = const Size(demoPhoneWidth, demoPanelHeight);
    tester.view.devicePixelRatio = 1.0;
    tester.view.viewInsets = FakeViewPadding.zero;

    await tester.pumpWidget(buildDemoSurface(
      label: label,
      labelColor: isBefore ? beforeColor : afterColor,
      app: ChangeNotifierProvider<SystemSettingsModel>(
        create: (_) => SystemSettingsModel(),
        child: MaterialApp(
          debugShowCheckedModeBanner: false,
          theme: buildAppTheme(),
          home: Scaffold(
            // surfaceDim is the page canvas the receipts route actually sits on, so
            // the bar is judged against the background it really has.
            backgroundColor: buildAppTheme().colorScheme.surfaceDim,
            body: Column(
              children: [
                ReceiptSummaryBar(summary: summary, atTop: true),
                const Expanded(child: SizedBox.shrink()),
              ],
            ),
          ),
        ),
      ),
    ));
    await tester.pumpAndSettle();

    // Report the geometry the still is meant to show, so a panel that silently
    // failed to render the block is obvious in the run output, not only in the image.
    final bar = tester.getRect(find.byKey(const ValueKey('receipt-summary')));
    // ignore: avoid_print
    print('panel=$variant barHeight=${bar.height} barWidth=${bar.width}');

    final frame = await grabFrame(tester);
    File(out).writeAsBytesSync(img.encodePng(frame));
    // ignore: avoid_print
    print('wrote $out (${frame.width}x${frame.height})');
  }, timeout: const Timeout(Duration(minutes: 3)));
}
