import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
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
