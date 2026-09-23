import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_filter_fields.dart';
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

  /// The group this model's quick-filter **state** belongs to, or null when it
  /// holds none.
  ///
  /// That state is the conditions *and* the quick date field: both are authored
  /// against one group, and both have to be dropped on the way into another.
  /// Scoping only the conditions left the field un-scoped, so changing it
  /// without applying a condition kept `_filterGroupId` null, `didChangeDependencies`
  /// skipped its reset, and the field followed the user into the next group.
  ///
  /// The filter has to carry its own scope because this model outlives every
  /// screen that reads it, while `GroupReceiptsList` is rebuilt from scratch on
  /// almost every navigation -- a receipt round trip destroys and remounts it
  /// just as a group change does. A freshly-mounted list therefore cannot tell
  /// those two apart by remembering the group it last saw: it has no memory.
  /// Asking the filter which group it belongs to distinguishes them.
  String? _filterGroupId;

  String? get filterGroupId => _filterGroupId;

  /// The date field the quick date control (the month stepper) reads and
  /// writes -- one of [receiptDateFilterFields].
  ///
  /// It lives here rather than in `GroupReceiptsList`'s State because that
  /// widget is torn down and rebuilt on almost every navigation, including a
  /// round trip to a receipt. A widget-local field would silently snap back to
  /// "Receipt Date" while the month it wrote still sat on `resolvedDate`,
  /// leaving the stepper describing a condition it does not own.
  String _quickDateField = defaultQuickDateFieldKey;

  String get quickDateField => _quickDateField;

  /// The condition the quick date control currently owns, if any.
  ReceiptFilterCondition? get quickDateCondition => _filter[_quickDateField];

  /// Whether anything here belongs to one group, and so has to be dropped when
  /// the user arrives in another: conditions, a non-default date field, or both.
  ///
  /// Recomputed on every write rather than tracked, so a field toggled away
  /// from the default and back leaves nothing group-scoped behind.
  bool get _hasGroupScopedState =>
      _filter.isNotEmpty || _quickDateField != defaultQuickDateFieldKey;

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
    _filterGroupId = _hasGroupScopedState ? groupId : null;
    if (notify) {
      notifyListeners();
    }
  }

  /// Replaces, or with a null [condition] removes, a single field's condition.
  ///
  /// The quick date control's write path, and the analogue of desktop's
  /// `SetReceiptFilterField`. It notifies, so the list refetches once -- unlike
  /// the filter screen, which authors a whole draft and commits it through
  /// [setFilter], the stepper applies on the tap.
  ///
  /// [groupId] is required for the same reason it is on [setFilter]: a filter
  /// that does not know its group can never be cleaned up.
  void setFilterField(
    String key,
    ReceiptFilterCondition? condition, {
    required String groupId,
  }) {
    final next = Map.of(_filter);
    if (condition == null) {
      next.remove(key);
    } else {
      next[key] = condition;
    }

    setFilter(next, true, groupId: groupId);
  }

  /// Re-points the quick date control at another date field.
  ///
  /// Deliberately non-destructive: it changes no condition, only which one the
  /// stepper describes, so a Receipt Date filter authored on the filter screen
  /// survives a move to Resolved Date and stays applied. Callers therefore pass
  /// `notify: false` -- a notification here means "the filter changed" and
  /// would refetch a result set that has not moved.
  ///
  /// [groupId] is required for the same reason it is on [setFilter], and the
  /// compiler asking is the point: the chosen field is authored against one
  /// group and has to be dropped on the way into another. It is recorded even
  /// with no conditions applied, which is the case that made the field outlive
  /// its group -- and dropped again when the field returns to its default,
  /// since nothing group-scoped is left at that point.
  void setQuickDateField(String key, bool notify, {required String groupId}) {
    if (_quickDateField == key) {
      return;
    }

    _quickDateField = key;
    _filterGroupId = _hasGroupScopedState ? groupId : null;
    if (notify) {
      notifyListeners();
    }
  }

  /// Drops every condition and re-points the quick date control at its default
  /// field, mirroring desktop's `ResetReceiptFilter`.
  ///
  /// The guard covers the field as well as the map, and [_hasGroupScopedState]
  /// is what gets this *called* for a field-only change: together they stop the
  /// stepper arriving in a new group still pointed at `resolvedDate`.
  void clearFilter(bool notify) {
    if (_filter.isEmpty && _quickDateField == defaultQuickDateFieldKey) {
      return;
    }

    _filter = {};
    _filterGroupId = null;
    _quickDateField = defaultQuickDateFieldKey;
    if (notify) {
      notifyListeners();
    }
  }
}
