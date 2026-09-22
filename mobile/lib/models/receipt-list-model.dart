import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/constants/receipts.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

class ReceiptListModel extends ChangeNotifier {
  String _orderBy = receiptSortOptions.first.columnName;

  String get orderBy => _orderBy;

  int _page = 1;

  int get page => _page;

  SortDirection _sortDirection = SortDirection.desc;

  SortDirection get sortDirection => _sortDirection;

  /// The applied filter, keyed by `ReceiptFilterField.key`.
  ///
  /// Only conditions the user has applied live here -- the filter sheet edits
  /// its own draft copy and writes back through [setFilter], so dismissing it
  /// discards rather than committing.
  Map<String, ReceiptFilterCondition> _filter = {};

  /// Unmodifiable so a caller cannot mutate the applied filter behind the
  /// model's back and leave the badge count disagreeing with the query.
  Map<String, ReceiptFilterCondition> get filter => Map.unmodifiable(_filter);

  /// The group [_filter] was authored in, or null when nothing is applied.
  ///
  /// The filter has to carry its own scope because this model outlives every
  /// screen that reads it, while `GroupReceiptsList` is rebuilt from scratch on
  /// almost every navigation -- a receipt round trip destroys and remounts it
  /// just as a group change does. A freshly-mounted list therefore cannot tell
  /// those two apart by remembering the group it last saw: it has no memory.
  /// Asking the filter which group it belongs to distinguishes them.
  String? _filterGroupId;

  String? get filterGroupId => _filterGroupId;

  /// How many conditions are narrowing the list, i.e. the app bar's badge.
  int get activeFilterCount => _filter.length;

  bool get hasActiveFilter => _filter.isNotEmpty;

  ReceiptPagedRequestCommand get receiptPagedRequestCommand =>
      (ReceiptPagedRequestCommandBuilder()
            ..page = _page
            ..pageSize = 10
            ..orderBy = _orderBy
            ..sortDirection = _sortDirection
            ..filter = buildReceiptPagedRequestFilter(_filter))
          .build();

  /// The same applied filter, shaped for the summary endpoint.
  ///
  /// A method rather than a getter because the synthetic All group has to name a
  /// configuration group. Deliberately carries NO page and NO sort: the summary
  /// aggregates the whole filtered set, which is exactly why paging and sorting must
  /// not refetch it.
  ///
  /// Built here rather than in the list widget so the table and the totals can never
  /// disagree about what is filtered -- and so nobody hands ONE
  /// `buildReceiptPagedRequestFilter` result to both commands, which share a mutable
  /// nested builder.
  ReceiptSummaryCommand receiptSummaryCommand({int? configurationGroupId}) =>
      (ReceiptSummaryCommandBuilder()
            ..filter = buildReceiptPagedRequestFilter(_filter)
            ..configurationGroupId = configurationGroupId)
          .build();

  int? _summaryConfigGroupId;

  int? get summaryConfigGroupId => _summaryConfigGroupId;

  /// Which member group's configuration shapes the All group's summary breakdown.
  ///
  /// Lives here, not in the list's State, for the reason [_filterGroupId] does: this
  /// model outlives every screen that reads it, while `GroupReceiptsList` is destroyed
  /// and remounted on every receipt round trip -- a widget-local pick would silently
  /// reset itself. Session-scoped only: mobile has no persisted slice equivalent to the
  /// desktop's NGXS `receiptTable`, so unlike there the pick does not survive a restart.
  ///
  /// [notify] exists for symmetry with the setters above; every call site passes false,
  /// for the same reason they do. Picking a configuration changes the breakdown's SHAPE,
  /// never the data, so notifying would make the list refetch for nothing.
  ///
  /// Reset on a GROUP change, by the list, not from [clearFilter]: that early-returns on
  /// an already-empty filter, so the pick would survive exactly the navigation that
  /// invalidates it -- and it is also what a user pressing "Reset" in the filter screen
  /// calls, where silently changing which group's configuration applies is a surprise.
  void setSummaryConfigGroupId(int? groupId, bool notify) {
    _summaryConfigGroupId = groupId;
    if (notify) {
      notifyListeners();
    }
  }

  void setOrderBy(String orderBy, bool notify) {
    _orderBy = orderBy;
    if (notify) {
      notifyListeners();
    }
  }

  void setPage(int page, bool notify) {
    _page = page;
    if (notify) {
      notifyListeners();
    }
  }

  void setSortDirection(SortDirection sortDirection, bool notify) {
    _sortDirection = sortDirection;
    if (notify) {
      notifyListeners();
    }
  }

  /// Replaces the applied filter.
  ///
  /// The receipts list listens for this notification to refresh itself. The sort
  /// setters above are deliberately called with `notify: false` (they refresh
  /// the list directly instead), so a notification from this model means the
  /// filter changed -- don't "tidy" those call sites into notifying.
  ///
  /// [groupId] is required rather than optional because a filter that does not
  /// know its group is one that can never be cleaned up: it silently follows the
  /// user into the next group, where its category, tag and user ids match
  /// nothing. Making the compiler ask is the point.
  void setFilter(
    Map<String, ReceiptFilterCondition> filter,
    bool notify, {
    required String groupId,
  }) {
    _filter = Map.of(filter);
    _filterGroupId = _filter.isEmpty ? null : groupId;
    if (notify) {
      notifyListeners();
    }
  }

  void clearFilter(bool notify) {
    if (_filter.isEmpty) {
      return;
    }

    _filter = {};
    _filterGroupId = null;
    if (notify) {
      notifyListeners();
    }
  }
}
