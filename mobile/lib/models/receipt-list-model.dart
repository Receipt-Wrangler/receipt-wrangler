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
  void setFilter(Map<String, ReceiptFilterCondition> filter, bool notify) {
    _filter = Map.of(filter);
    if (notify) {
      notifyListeners();
    }
  }

  void clearFilter(bool notify) {
    if (_filter.isEmpty) {
      return;
    }

    _filter = {};
    if (notify) {
      notifyListeners();
    }
  }
}
