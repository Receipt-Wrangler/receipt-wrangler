import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/receipts.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/receipt_list_item.dart';
import 'package:receipt_wrangler_mobile/models/receipt-list-model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/paged_data_list.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';

import '../../client/client.dart';

class GroupReceiptsList extends StatefulWidget {
  const GroupReceiptsList({super.key});

  @override
  State<GroupReceiptsList> createState() => _GroupReceiptsList();
}

class _GroupReceiptsList extends State<GroupReceiptsList> {
  VoidCallback? _refreshCallback;

  late final ReceiptListModel _receiptListModel =
      Provider.of<ReceiptListModel>(context, listen: false);

  @override
  void initState() {
    super.initState();
    // The sort setters deliberately pass notify: false and refresh the list
    // themselves, so a notification from this model means the applied filter
    // changed -- which is the one state change the list cannot see for itself.
    _receiptListModel.addListener(_refreshForFilterChange);
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();

    // getGroupId reads GoRouterState, an inherited widget, so this cannot be
    // done from initState.
    final groupId = getGroupId(context);
    final filterGroupId = _receiptListModel.filterGroupId;
    if (filterGroupId != null && filterGroupId != groupId) {
      // A filter holds the previous group's category, tag and user ids, which
      // match nothing here -- and would leave the badge counting conditions the
      // user cannot see. Clearing everything is the predictable rule: switching
      // groups shows that group's receipts.
      //
      // The scope is read off the FILTER, not off a group this widget
      // remembered. The app offers no lateral group switch -- the app-bar arrow
      // goes to /groups, group cards go to /groups/<id>/dashboards -- so every
      // real group change leaves the group shell and this State is rebuilt with
      // no memory of where the user came from. A widget-local "last group I
      // saw" is null on arrival and the clear never fires, which is exactly the
      // leak that shipped. Asking the filter also keeps the case this must NOT
      // clear: a round trip to a receipt remounts this list just the same, but
      // the filter still belongs to the group we are returning to.
      //
      // Silently: this runs during a build, where notifying would rebuild the
      // badge mid-frame. The refresh below is the visible half -- and is null on
      // a first mount, where the clear simply lands before the first fetch.
      _receiptListModel.clearFilter(false);
      _refreshCallback?.call();
    }
  }

  @override
  void dispose() {
    _receiptListModel.removeListener(_refreshForFilterChange);
    super.dispose();
  }

  void _refreshForFilterChange() {
    if (!mounted) {
      return;
    }

    // setState as well as refetch: the empty-state text below reads the applied
    // filter, and the list is otherwise built with listen: false. Safe to call
    // here because the only notifying writer is the filter screen's Apply, a
    // user event -- the group reset writes silently, during a build.
    setState(() {});
    _refreshCallback?.call();
  }

  Widget buildSortFilterBar() {
    return Row(
      children: [
        buildSortChip(),
        SizedBox(
          width: 4,
        ),
        buildSortDirectionChip()
      ],
    );
  }

  Widget buildSortDirectionChip() {
    var direction =
        Provider.of<ReceiptListModel>(context, listen: false).sortDirection;
    return PopupMenuButton(
      child: Chip(
        label: Text(
            direction == api.SortDirection.asc ? "Ascending" : "Descending"),
      ),
      itemBuilder: (context) {
        return [
          PopupMenuItem(
            child: Text("Sort Ascending"),
            value: api.SortDirection.asc,
          ),
          PopupMenuItem(
            child: Text("Sort Descending"),
            value: api.SortDirection.desc,
          ),
        ];
      },
      onSelected: (value) {
        var model = Provider.of<ReceiptListModel>(context, listen: false);
        model.setSortDirection(value, false);
        _refreshCallback?.call();
      },
    );
  }

  Widget buildSortChip() {
    return PopupMenuButton(
      child: Chip(
        label: Text(getSortChipText()),
      ),
      itemBuilder: (context) => receiptSortOptions.map((option) {
        return PopupMenuItem(
          child: Text("Sort by ${option.displayLabel}"),
          value: option.columnName,
        );
      }).toList(),
      onSelected: (value) {
        var model = Provider.of<ReceiptListModel>(context, listen: false);
        model.setOrderBy(value, false);
        _refreshCallback?.call();
      },
    );
  }

  String getSortChipText() {
    var model = Provider.of<ReceiptListModel>(context, listen: false);
    var orderBy = model.orderBy;

    var option = receiptSortOptions.firstWhere(
      (element) => element.columnName == orderBy,
      orElse: () => receiptSortOptions.first,
    );

    return option.displayLabel;
  }

  @override
  Widget build(BuildContext context) {
    return Column(
      children: [
        buildSortFilterBar(),
        PagedDataList(
          onRefreshCallbackSet: (callback) {
            _refreshCallback = callback;
          },
          noItemsFoundText: _receiptListModel.hasActiveFilter
              ? "No receipts match this filter"
              : "No receipts found",
          listItemBuilder: (context, receipt, index) {
            return ReceiptListItem(
                receipt: receipt.anyOf.values[0] as api.Receipt);
          },
          getPagedDataFuture: (pageKey) {
            var model = Provider.of<ReceiptListModel>(context, listen: false);
            model.setPage(pageKey, false);
            var command = model.receiptPagedRequestCommand;

            return OpenApiClient.client.getReceiptApi().getReceiptsForGroup(
                  groupId: int.parse(getGroupId(context)),
                  receiptPagedRequestCommand: command,
                );
          },
        ),
      ],
    );
  }
}
