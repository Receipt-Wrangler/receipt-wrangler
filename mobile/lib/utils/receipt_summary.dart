import 'package:openapi/openapi.dart' as api;

import '../models/group_model.dart';

/// Pure helpers for the receipt summary bar -- the Dart twin of the desktop's
/// `src/utils/receipt-summary.ts`. Kept out of the widget so the labelling and the
/// configuration-group fallback are unit-testable without pumping anything.

/// Groups whose summary configuration a viewer may choose between, alphabetically.
///
/// The synthetic "All" group is excluded: it is the view that NEEDS someone else's
/// configuration, never a candidate to supply one.
List<api.Group> summaryConfigGroups(GroupModel groupModel) {
  final enabled = groupModel.groupsWithoutAllGroup
      .where((group) => group.groupReceiptSettings.receiptSummaryEnabled == true)
      .toList();
  // Case-insensitive, to match the desktop twin's localeCompare. Dart's compareTo is
  // UTF-16 ORDINAL, so "Bravo team" sorts before "apple budget" ('B' 66 < 'a' 97) while
  // localeCompare puts them the other way round -- and since the fallback is "the first
  // enabled group", the two clients would default the All group to DIFFERENT
  // configurations, showing the same user different columns on phone and on desktop.
  // The tie-break keeps the order stable when two names differ only in case.
  enabled.sort((a, b) {
    final byName = a.name.toLowerCase().compareTo(b.name.toLowerCase());
    return byName != 0 ? byName : a.name.compareTo(b.name);
  });
  return enabled;
}

/// The stored pick when it still names an enabled group, else the first.
///
/// One branch covers every stale case: the picked group turned its summary off, the
/// user left it, nothing was ever picked, or `GroupModel` was replaced wholesale by an
/// AppData refresh while the screen was open.
api.Group? resolveSummaryConfigGroup(List<api.Group> enabledGroups, int? pickedGroupId) {
  for (final group in enabledGroups) {
    if (group.id == pickedGroupId) {
      return group;
    }
  }
  return enabledGroups.isEmpty ? null : enabledGroups.first;
}

/// The test id / key suffix for a row: `overall`, else the status name.
String summaryRowKey(api.ReceiptSummaryRow row, {required bool isOverall}) {
  return isOverall ? 'overall' : row.status.name;
}

/// "All Receipts" for the overall row, else "<Status> Receipts".
///
/// The generated [api.ReceiptSummaryRow.status] is non-nullable, so the overall row
/// arrives carrying [api.ReceiptStatus.empty] rather than null -- and
/// `receiptStatusLabel` renders that as "", which would read as a bare " Receipts".
String receiptSummaryRowLabel(
  api.ReceiptSummaryRow row,
  String Function(api.ReceiptStatus) statusLabel, {
  required bool isOverall,
}) {
  if (isOverall || row.status == api.ReceiptStatus.empty) {
    return 'All Receipts';
  }
  return '${statusLabel(row.status)} Receipts';
}

String receiptSummaryCountLabel(api.ReceiptSummaryRow row) {
  return row.receiptCount == 1 ? '1 receipt' : '${row.receiptCount} receipts';
}

/// Anything that is not TOP renders at the bottom -- including the `empty` member, which
/// exists only so an already-released build tolerates a value added later. That makes the
/// client's degraded case match the server's default rather than blanking the block.
bool isReceiptSummaryAtTop(api.ReceiptSummaryPosition? position) {
  return position == api.ReceiptSummaryPosition.TOP;
}
