import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_summary.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';

import '../helpers/receipt_form_test_helpers.dart';
import '../helpers/receipt_summary_test_helpers.dart';

/// The labelling and the configuration-group fallback live in pure functions precisely
/// so they can be pinned without pumping a widget.
void main() {
  GroupModel modelWith(List<api.Group> groups) => GroupModel()..setGroups(groups);

  group('summaryConfigGroups', () {
    test('offers every enabled group alphabetically', () {
      final model = modelWith([
        buildGroup(id: 2, name: 'Zulu', receiptSummaryEnabled: true),
        buildGroup(id: 3, name: 'Alpha', receiptSummaryEnabled: true),
      ]);

      expect(summaryConfigGroups(model).map((g) => g.name), ['Alpha', 'Zulu']);
    });

    // Case-INsensitively, to match the desktop twin's localeCompare. Dart's compareTo is
    // UTF-16 ordinal, so this list comes back in the opposite order without the fix --
    // and because the fallback is "the first enabled group", the two clients would default
    // the All group to different configurations for the same user.
    test('sorts case-insensitively, like the desktop twin', () {
      final model = modelWith([
        buildGroup(id: 2, name: 'Bravo team', receiptSummaryEnabled: true),
        buildGroup(id: 3, name: 'apple budget', receiptSummaryEnabled: true),
      ]);

      expect(summaryConfigGroups(model).map((g) => g.name),
          ['apple budget', 'Bravo team']);
      expect(resolveSummaryConfigGroup(summaryConfigGroups(model), null)?.id, 3);
    });

    test('excludes a group whose summary is off', () {
      final model = modelWith([
        buildGroup(id: 2, name: 'On', receiptSummaryEnabled: true),
        buildGroup(id: 3, name: 'Off', receiptSummaryEnabled: false),
      ]);

      expect(summaryConfigGroups(model).map((g) => g.name), ['On']);
    });

    // The All group is the view that NEEDS a configuration, never one that supplies it.
    test('excludes the All group even with the summary enabled', () {
      final model = modelWith([
        buildGroup(id: 1, name: 'All', isAllGroup: true, receiptSummaryEnabled: true),
        buildGroup(id: 2, name: 'Household', receiptSummaryEnabled: true),
      ]);

      expect(summaryConfigGroups(model).map((g) => g.name), ['Household']);
    });
  });

  group('resolveSummaryConfigGroup', () {
    final alpha = buildGroup(id: 2, name: 'Alpha', receiptSummaryEnabled: true);
    final bravo = buildGroup(id: 3, name: 'Bravo', receiptSummaryEnabled: true);

    test('prefers the pick', () {
      expect(resolveSummaryConfigGroup([alpha, bravo], 3)?.id, 3);
    });

    test('falls back to the first with no pick', () {
      expect(resolveSummaryConfigGroup([alpha, bravo], null)?.id, 2);
    });

    // The picked group turned its summary off, or the user left it -- one branch covers
    // every stale case, including an AppData refresh replacing GroupModel mid-session.
    test('falls back when the pick is stale', () {
      expect(resolveSummaryConfigGroup([alpha, bravo], 99)?.id, 2);
    });

    test('resolves to nothing when no group qualifies', () {
      expect(resolveSummaryConfigGroup([], 2), isNull);
    });
  });

  group('labels', () {
    test('the overall row reads All Receipts, not a bare " Receipts"', () {
      // ReceiptSummaryRow.status is non-nullable, so the overall row arrives carrying
      // the empty member -- which receiptStatusLabel renders as "".
      final row = buildSummaryRow(receiptCount: 3);

      expect(receiptSummaryRowLabel(row, receiptStatusLabel, isOverall: true),
          'All Receipts');
      expect(receiptSummaryRowLabel(row, receiptStatusLabel, isOverall: false),
          'All Receipts');
    });

    test('a status row names its status', () {
      final row = buildSummaryRow(status: api.ReceiptStatus.NEEDS_ATTENTION);

      expect(receiptSummaryRowLabel(row, receiptStatusLabel, isOverall: false),
          'Needs Attention Receipts');
    });

    test('the count singularizes at exactly one', () {
      expect(receiptSummaryCountLabel(buildSummaryRow(receiptCount: 0)), '0 receipts');
      expect(receiptSummaryCountLabel(buildSummaryRow(receiptCount: 1)), '1 receipt');
      expect(receiptSummaryCountLabel(buildSummaryRow(receiptCount: 2)), '2 receipts');
    });

    test('the row key is overall for the overall row and the status name otherwise', () {
      expect(summaryRowKey(buildSummaryRow(), isOverall: true), 'overall');
      expect(
        summaryRowKey(buildSummaryRow(status: api.ReceiptStatus.OPEN), isOverall: false),
        'OPEN',
      );
    });
  });

  group('isReceiptSummaryAtTop', () {
    test('only TOP is top', () {
      expect(isReceiptSummaryAtTop(api.ReceiptSummaryPosition.TOP), isTrue);
      expect(isReceiptSummaryAtTop(api.ReceiptSummaryPosition.BOTTOM), isFalse);
    });

    // The empty member exists only so an already-released build tolerates a value added
    // later. Degrading to the bottom matches the server default; blanking the block
    // would not.
    test('the empty member and null degrade to the bottom', () {
      expect(isReceiptSummaryAtTop(api.ReceiptSummaryPosition.empty), isFalse);
      expect(isReceiptSummaryAtTop(null), isFalse);
    });
  });
}
