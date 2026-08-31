import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:infinite_carousel/infinite_carousel.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/constants/receipt_entry.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/quick_scan.dart';
import 'package:receipt_wrangler_mobile/shared/functions/quick_scan.dart';
import 'package:rxdart/rxdart.dart';

import '../helpers/permission_test_helpers.dart';
import '../helpers/receipt_entry_test_helpers.dart';
import '../helpers/receipt_form_test_helpers.dart';

/// The Quick Scan sheet is now reachable from the nav tap, the long-press menu
/// and the overflow menu, so it re-checks its own gates rather than trusting the
/// caller -- and it offers a way out to manual entry only for users who could
/// actually save one.
void main() {
  const quickScan = api.Permission.groupPeriodReceiptsPeriodQuickScan;
  const create = api.Permission.groupPeriodReceiptsPeriodCreate;
  const manualLinkKey = ValueKey('quick-scan-manual-entry-link');
  const pageIndicatorKey = ValueKey('quick-scan-page-indicator');

  late BuildContext capturedContext;

  Future<void> pumpSheet(
    WidgetTester tester, {
    required bool aiEnabled,
    required List<api.Permission> permissions,
  }) async {
    final router = GoRouter(
      initialLocation: '/groups/5/receipts',
      routes: [
        GoRoute(
          path: '/groups/:groupId/receipts',
          builder: (context, state) {
            capturedContext = context;
            return const Scaffold(body: SizedBox.shrink());
          },
        ),
        GoRoute(
          path: '/receipts/add',
          builder: (context, state) => const Scaffold(body: Text('Add screen')),
        ),
      ],
    );

    await tester.pumpWidget(pumpReceiptEntryApp(
      router: router,
      aiEnabled: aiEnabled,
      permissions: seededPermissions(group: {5: permissions}),
      groups: [buildGroup(id: 5, name: 'Household')],
    ));
    await tester.pump();

    showQuickScanBottomSheet(capturedContext);
    await tester.pumpAndSettle();
  }

  /// Mounts the carousel with [pageCount] pages already loaded.
  ///
  /// The sheet builds its own `imageSubject` internally and seeds it empty, so
  /// [pumpSheet] can never reach a loaded state -- it can only prove the counter
  /// is absent. Mounting [QuickScan] directly is what lets the counter be
  /// asserted present, which is the whole point of keying it.
  Future<void> pumpCarousel(
    WidgetTester tester, {
    required int pageCount,
  }) async {
    final images = List.generate(pageCount, (_) => buildQuickScanImage(groupId: 5));

    final router = GoRouter(
      initialLocation: '/',
      routes: [
        GoRoute(
          path: '/',
          builder: (context, state) => Scaffold(
            body: QuickScan(
              imageSubject: BehaviorSubject.seeded(images),
              infiniteScrollController: InfiniteScrollController(),
              isCompletedSubject: BehaviorSubject.seeded(false),
            ),
          ),
        ),
      ],
    );

    await tester.pumpWidget(pumpReceiptEntryApp(
      router: router,
      aiEnabled: true,
      permissions: seededPermissions(group: {5: [quickScan, create]}),
      groups: [buildGroup(id: 5, name: 'Household')],
    ));
    await tester.pump();
  }

  testWidgets('opens for a user who holds the quick-scan permission',
      (tester) async {
    await pumpSheet(tester, aiEnabled: true, permissions: [quickScan, create]);

    expect(find.text(quickScanLabel), findsOneWidget);
    expect(find.text('Scan or upload an image to get started'), findsOneWidget);
  });

  testWidgets('refuses to open without the permission, naming the group',
      (tester) async {
    await pumpSheet(tester, aiEnabled: true, permissions: [create]);

    expect(find.text('Scan or upload an image to get started'), findsNothing);
    expect(find.text(quickScanNoPermissionMessageForGroup('Household')),
        findsOneWidget);
  });

  testWidgets('refuses to open with the ai flag off, naming that reason',
      (tester) async {
    await pumpSheet(tester,
        aiEnabled: false, permissions: [quickScan, create]);

    expect(find.text('Scan or upload an image to get started'), findsNothing);
    expect(find.text(quickScanAiDisabledMessage), findsOneWidget);
  });

  testWidgets('offers manual entry to a user who can create receipts',
      (tester) async {
    await pumpSheet(tester, aiEnabled: true, permissions: [quickScan, create]);

    expect(find.byKey(manualLinkKey), findsOneWidget);
    expect(find.text(enterDetailsManuallyLabel), findsOneWidget);
  });

  testWidgets('hides the manual link without the create permission',
      (tester) async {
    await pumpSheet(tester, aiEnabled: true, permissions: [quickScan]);

    expect(find.byKey(manualLinkKey), findsNothing,
        reason: 'the manual form would only reject their save');
  });

  testWidgets('nothing confirms a queued scan until one is submitted',
      (tester) async {
    await pumpSheet(tester, aiEnabled: true, permissions: [quickScan, create]);

    expect(find.byKey(const ValueKey('quick-scan-queued-confirmation')),
        findsNothing);
    expect(find.text(quickScanQueuedMessage), findsNothing);
  });

  testWidgets('a single-page scan gets no page counter', (tester) async {
    await pumpCarousel(tester, pageCount: 1);

    expect(find.byKey(pageIndicatorKey), findsNothing,
        reason: 'nothing to swipe to');
  });

  testWidgets('a multi-page scan counts its pages', (tester) async {
    await pumpCarousel(tester, pageCount: 3);

    expect(find.byKey(pageIndicatorKey), findsOneWidget);
    expect(find.text('1 of 3 · swipe for the next'), findsOneWidget,
        reason: 'a scan carries up to 100 pages and the carousel gives no '
            'other hint the later ones exist');
  });

  testWidgets('the manual link closes the sheet and opens the form',
      (tester) async {
    await pumpSheet(tester, aiEnabled: true, permissions: [quickScan, create]);

    await tester.tap(find.byKey(manualLinkKey));
    await tester.pumpAndSettle();

    expect(find.text('Add screen'), findsOneWidget);
    expect(find.text(quickScanLabel), findsNothing,
        reason: 'leaving the modal above the form would hide the screen the '
            'user just asked for');
  });
}
