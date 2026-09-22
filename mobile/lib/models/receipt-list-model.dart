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
  void setQuickDateField(String key, bool notify) {
    if (_quickDateField == key) {
      return;
    }

    _quickDateField = key;
    if (notify) {
      notifyListeners();
    }
  }

  /// Drops every condition and re-points the quick date control at its default
  /// field, mirroring desktop's `ResetReceiptFilter`.
  ///
  /// The guard covers the field as well as the map: a filter cleared on a group
  /// change would otherwise leave the stepper pointed at `resolvedDate` in a
  /// group the user has only just arrived in.
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
