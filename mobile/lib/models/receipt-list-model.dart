import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/constants/receipts.dart';

class ReceiptListModel extends ChangeNotifier {
  String _orderBy = receiptSortOptions.first.columnName;

  String get orderBy => _orderBy;

  int _page = 1;

  int get page => _page;

  SortDirection _sortDirection = SortDirection.desc;

  SortDirection get sortDirection => _sortDirection;

  ReceiptPagedRequestCommand get receiptPagedRequestCommand =>
      (ReceiptPagedRequestCommandBuilder()
            ..page = _page
            ..pageSize = 10
            ..orderBy = _orderBy
            ..sortDirection = _sortDirection
            ..filter = ReceiptPagedRequestFilterBuilder())
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
}
