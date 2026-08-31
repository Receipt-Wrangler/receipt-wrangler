import 'dart:math' as math;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/constants/colors.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_entry.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/quick_scan_unavailable_banner.dart';
import 'package:receipt_wrangler_mobile/shared/functions/receipt_entry_availability.dart';

/// The banner is the only thing telling a user why tapping Scan gave them a
/// form instead of a camera, so it has to name the right reason -- an
/// administrator-configuration problem and a missing permission need different
/// people to fix them.
void main() {
  const bannerKey = ValueKey('quick-scan-unavailable-banner');
  const dismissKey = ValueKey('quick-scan-unavailable-banner-dismiss');

  Widget wrap(Widget? banner) => MaterialApp(
        home: Scaffold(body: banner ?? const Text('no banner')),
      );

  testWidgets('the ai reason points at the server configuration',
      (tester) async {
    await tester.pumpWidget(wrap(const QuickScanUnavailableBanner(
        reason: QuickScanBlockedReason.aiDisabled)));

    expect(find.byKey(bannerKey), findsOneWidget);
    expect(find.text(quickScanUnavailableTitle), findsOneWidget);
    expect(find.text(quickScanAiDisabledMessage), findsOneWidget);
  });

  testWidgets('the permission reason names the group it applies to',
      (tester) async {
    await tester.pumpWidget(wrap(const QuickScanUnavailableBanner(
      reason: QuickScanBlockedReason.noPermission,
      groupName: 'Household',
    )));

    expect(find.text(quickScanNoPermissionMessageForGroup('Household')),
        findsOneWidget);
  });

  testWidgets('the permission reason stays generic with no single group',
      (tester) async {
    await tester.pumpWidget(wrap(const QuickScanUnavailableBanner(
        reason: QuickScanBlockedReason.noPermission)));

    expect(find.text(quickScanNoPermissionMessage), findsOneWidget);
  });

  testWidgets('reads as white on the dark notice surface', (tester) async {
    // The theme's ColorScheme never defines secondaryContainer, so the tokens
    // fall back to `secondary` (#8EA1AC) with black `onSecondary` -- muddy at
    // this type size. White on #8EA1AC would be worse still (2.7:1), so the
    // surface is darkened to carry white text at ~8.9:1.
    await tester.pumpWidget(wrap(const QuickScanUnavailableBanner(
        reason: QuickScanBlockedReason.aiDisabled)));

    final container = tester.widget<Container>(find.byKey(bannerKey));
    // The banner is rounded, so its colour rides on the decoration.
    expect((container.decoration as BoxDecoration).color, noticeSurface);

    for (final text in <String>[
      quickScanUnavailableTitle,
      quickScanAiDisabledMessage,
    ]) {
      expect(tester.widget<Text>(find.text(text)).style?.color, onNoticeSurface,
          reason: '"$text" must be white on the notice surface');
    }

    expect(tester.widget<Icon>(find.byIcon(Icons.info_outline)).color,
        onNoticeSurface);
  });

  testWidgets('the notice surface carries white text at readable contrast',
      (tester) async {
    // Guards the pair itself: a later palette tweak that drops contrast below
    // the WCAG AA threshold for normal text should fail here, not in review.
    double luminance(Color c) {
      double channel(double v) =>
          v <= 0.03928 ? v / 12.92 : math.pow((v + 0.055) / 1.055, 2.4) as double;
      return 0.2126 * channel(c.r) +
          0.7152 * channel(c.g) +
          0.0722 * channel(c.b);
    }

    final a = luminance(noticeSurface);
    final b = luminance(onNoticeSurface);
    final contrast =
        (math.max(a, b) + 0.05) / (math.min(a, b) + 0.05);

    expect(contrast, greaterThanOrEqualTo(4.5),
        reason: 'WCAG AA for normal text');
  });

  testWidgets('dismissing removes it', (tester) async {
    await tester.pumpWidget(wrap(const QuickScanUnavailableBanner(
        reason: QuickScanBlockedReason.aiDisabled)));

    await tester.tap(find.byKey(dismissKey));
    await tester.pump();

    expect(find.byKey(bannerKey), findsNothing);
  });

  group('fromRouteExtra', () {
    test('builds from a reason the entry point attached', () {
      final banner = QuickScanUnavailableBanner.fromRouteExtra({
        quickScanBlockedReasonExtraKey: QuickScanBlockedReason.noPermission,
        quickScanBlockedGroupExtraKey: 'Household',
      });

      expect(banner, isNotNull);
      expect(banner!.reason, QuickScanBlockedReason.noPermission);
      expect(banner.groupName, 'Household');
    });

    test('is absent for a deliberate Add Manual Receipt', () {
      // That path passes only the group id, so the form must stay quiet.
      expect(QuickScanUnavailableBanner.fromRouteExtra({'groupId': '5'}),
          isNull);
    });

    test('is absent when the route carries no extra at all', () {
      expect(QuickScanUnavailableBanner.fromRouteExtra(null), isNull);
    });

    test('ignores an extra of the wrong shape', () {
      expect(QuickScanUnavailableBanner.fromRouteExtra('not a map'), isNull);
      expect(
        QuickScanUnavailableBanner.fromRouteExtra(
            {quickScanBlockedReasonExtraKey: 'aiDisabled'}),
        isNull,
      );
    });
  });
}
