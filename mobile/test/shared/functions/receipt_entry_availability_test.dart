import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/shared/functions/receipt_entry_availability.dart';

import '../../helpers/permission_test_helpers.dart';
import '../../helpers/receipt_entry_test_helpers.dart';
import '../../helpers/receipt_form_test_helpers.dart';

/// The gate matrix behind every receipt-entry affordance.
///
/// Three independent inputs decide what the Scan/Add slot offers -- the
/// install-wide `aiPoweredReceipts` flag, `group.receipts.quick-scan` and
/// `group.receipts.create` -- and the backend enforces the two permissions
/// separately, so every combination is reachable in production.
void main() {
  const quickScan = api.Permission.groupPeriodReceiptsPeriodQuickScan;
  const create = api.Permission.groupPeriodReceiptsPeriodCreate;

  /// Resolves availability inside group 5, or (when [inGroup] is false) from a
  /// route with no group -- the group-select / all-groups case, where the check
  /// falls back to "held in any group".
  Future<ReceiptEntryAvailability> resolve(
    WidgetTester tester, {
    required bool aiEnabled,
    required List<api.Permission> permissions,
    bool inGroup = true,
    bool groupIsAllGroup = false,
  }) async {
    late ReceiptEntryAvailability result;

    final router = GoRouter(
      initialLocation: inGroup ? '/groups/5/receipts' : '/',
      routes: [
        GoRoute(
          path: '/',
          builder: (context, state) => _Probe((c) => result = c),
        ),
        GoRoute(
          path: '/groups/:groupId/receipts',
          builder: (context, state) => _Probe((c) => result = c),
        ),
      ],
    );

    await tester.pumpWidget(MultiProvider(
      providers: [
        ChangeNotifierProvider<AuthModel>.value(
            value: authModelWithAi(aiEnabled)),
        ChangeNotifierProvider<GroupModel>.value(
          value: groupModelWith([
            buildGroup(id: 5, name: 'Household', isAllGroup: groupIsAllGroup),
          ]),
        ),
        ChangeNotifierProvider<PermissionsModel>.value(
          value: seededPermissions(group: {5: permissions}),
        ),
      ],
      child: MaterialApp.router(routerConfig: router),
    ));
    await tester.pump();

    return result;
  }

  group('resolveReceiptEntryAvailability inside a group', () {
    testWidgets('ai on + both permissions: scan and manual both offered',
        (tester) async {
      final a = await resolve(tester,
          aiEnabled: true, permissions: [quickScan, create]);

      expect(a.canQuickScan, isTrue);
      expect(a.canCreateManual, isTrue);
      expect(a.blockedReason, isNull);
      expect(a.isVisible, isTrue);
      expect(a.groupName, 'Household');
    });

    testWidgets('ai off blocks quick scan even with the permission',
        (tester) async {
      final a = await resolve(tester,
          aiEnabled: false, permissions: [quickScan, create]);

      expect(a.canQuickScan, isFalse);
      expect(a.blockedReason, QuickScanBlockedReason.aiDisabled);
      expect(a.canCreateManual, isTrue);
      expect(a.isVisible, isTrue);
    });

    testWidgets('ai on without the quick-scan permission blocks on permission',
        (tester) async {
      final a = await resolve(tester, aiEnabled: true, permissions: [create]);

      expect(a.canQuickScan, isFalse);
      expect(a.blockedReason, QuickScanBlockedReason.noPermission);
      expect(a.canCreateManual, isTrue);
    });

    testWidgets(
        'the ai flag wins when both are missing: it is the install-wide reason',
        (tester) async {
      final a = await resolve(tester, aiEnabled: false, permissions: [create]);

      expect(a.blockedReason, QuickScanBlockedReason.aiDisabled,
          reason: 'sending the user to chase a permission they still could not '
              'use would point them at the wrong person');
    });

    testWidgets('quick scan without create: scan offered, manual is not',
        (tester) async {
      // The backend enforces the two permissions separately, so this role is
      // real -- the user may scan but must not be offered the manual form.
      final a =
          await resolve(tester, aiEnabled: true, permissions: [quickScan]);

      expect(a.canQuickScan, isTrue);
      expect(a.canCreateManual, isFalse);
      expect(a.isVisible, isTrue);
    });

    testWidgets('create without quick scan: manual offered, scan is not',
        (tester) async {
      final a = await resolve(tester, aiEnabled: true, permissions: [create]);

      expect(a.canQuickScan, isFalse);
      expect(a.canCreateManual, isTrue);
      expect(a.isVisible, isTrue);
    });

    testWidgets('neither permission: nothing to offer, so nothing is visible',
        (tester) async {
      final a = await resolve(tester,
          aiEnabled: true,
          permissions: [api.Permission.groupPeriodReceiptsPeriodRead]);

      expect(a.canQuickScan, isFalse);
      expect(a.canCreateManual, isFalse);
      expect(a.isVisible, isFalse);
    });

    testWidgets('ai off and no create: still nothing visible', (tester) async {
      final a = await resolve(tester,
          aiEnabled: false, permissions: [quickScan]);

      expect(a.isVisible, isFalse,
          reason: 'the quick-scan permission alone is unusable with ai off');
    });
  });

  group('resolveReceiptEntryAvailability with no single group', () {
    testWidgets('group-select falls back to "held in any group"',
        (tester) async {
      final a = await resolve(tester,
          aiEnabled: true, permissions: [quickScan, create], inGroup: false);

      expect(a.canQuickScan, isTrue);
      expect(a.canCreateManual, isTrue);
      expect(a.groupName, isNull,
          reason: 'no single group to name in the banner copy');
    });

    testWidgets('the all-groups view uses the same fallback', (tester) async {
      final a = await resolve(tester,
          aiEnabled: true,
          permissions: [quickScan, create],
          groupIsAllGroup: true);

      expect(a.canQuickScan, isTrue);
      expect(a.groupName, isNull);
    });

    testWidgets('group-select with the permission held nowhere hides the slot',
        (tester) async {
      final a = await resolve(tester,
          aiEnabled: true,
          permissions: [api.Permission.groupPeriodReceiptsPeriodRead],
          inGroup: false);

      expect(a.isVisible, isFalse);
    });
  });
}

/// Captures the resolved availability from a real element in the tree, so the
/// Provider and GoRouter lookups run exactly as they do in the app.
class _Probe extends StatelessWidget {
  const _Probe(this.onResolved);

  final void Function(ReceiptEntryAvailability) onResolved;

  @override
  Widget build(BuildContext context) {
    onResolved(resolveReceiptEntryAvailability(context));
    return const Scaffold(body: SizedBox.shrink());
  }
}
