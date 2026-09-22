import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_filter_fields.dart';
import 'package:receipt_wrangler_mobile/constants/receipts.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/receipt_list_item.dart';
import 'package:receipt_wrangler_mobile/models/receipt-list-model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_month_picker_sheet.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_month_stepper.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/paged_data_list.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_date_filter.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter_options.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_summary.dart';

import '../../client/client.dart';
import 'receipt_summary_bar.dart';

class GroupReceiptsList extends StatefulWidget {
  const GroupReceiptsList({super.key});

  @override
  State<GroupReceiptsList> createState() => _GroupReceiptsList();
}

/// What [_GroupReceiptsList._resolveSummaryRequest] decided: whether to ask at all, and
/// if so whose configuration to ask for. A bare `int?` cannot express this, because
/// `configurationGroupId: null` is MEANINGFUL on the wire -- it means "use the path
/// group" -- and would be read as "skip" at the call site.
typedef _SummaryRequest = ({bool shouldRequest, int? configurationGroupId});

class _GroupReceiptsList extends State<GroupReceiptsList> {
  VoidCallback? _refreshCallback;

  api.ReceiptSummary? _summary;

  /// The group [_summary] belongs to. didChangeDependencies fires for ANY inherited
  /// widget change and getGroupId makes GoRouterState one, so without this guard an
  /// unpaged aggregate -- the expensive request -- fires on incidental rebuilds.
  String? _summaryGroupId;

  /// Last-request-wins, the manual equivalent of the desktop's switchMap. Two filter
  /// applies in quick succession can land backwards and paint figures for a filter the
  /// user has already moved past.
  int _summaryRequestSeq = 0;

  /// Keeps the paged list's State across a reparent.
  ///
  /// The summary bar renders in one of two slots either side of this list, so the list's
  /// position among its siblings changes when the group's configured position does.
  /// Element.updateChildren walks the old and new child lists inward from both ends while
  /// Widget.canUpdate holds and reuses only KEYED children in the middle it cannot match
  /// -- so without a key PagedDataList is rebuilt from scratch, discarding its paging
  /// controller, every loaded page and _totalCount, and silently refetching page 1.
  ///
  /// A key rather than padding the Column with always-present placeholder slots: the key
  /// travels with the widget it protects, so adding any other conditional child here
  /// later cannot reintroduce the bug. Pinned by "flipping the position does not reset
  /// the paged list", which asserts State identity -- a request count does not catch it,
  /// because the reconstructed State refetches immediately and that looks exactly like
  /// the filter's own refetch.
  final GlobalKey _pagedListKey = GlobalKey();

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

    // Guarded on the group, not run unconditionally: this method fires for any inherited
    // widget change, and the summary is an UNPAGED aggregate.
    if (_summaryGroupId != groupId) {
      // The All-group configuration pick belongs to the group being browsed, so it goes
      // with the group -- NOT inside the filter branch above, which is gated on there
      // having been a filter at all. A user who never filtered would otherwise carry a
      // pick from the All group into and back out of a real one.
      _receiptListModel.setSummaryConfigGroupId(null, false);
      _summaryGroupId = groupId;
      _fetchSummary(groupId);
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
    // here because every notifying writer is a user EVENT -- the filter screen's
    // Apply and the month stepper's tap (setFilterField) -- while the group reset
    // and the date-field re-point write silently.
    setState(() {});
    _refreshCallback?.call();
    // The filter changed, so the figures did too. Stepping the month reaches here and
    // must, since it narrows the result set. What does NOT is the sort setters and the
    // date-field re-point, which call _refreshCallback directly (or nothing at all) with
    // notify: false -- and that is what makes "only a real filter change refetches the
    // summary" true by construction rather than by a rule someone has to remember.
    // Paging never reaches here either.
    _fetchSummary(getGroupId(context));
  }

  /// Which configuration the summary should be asked for, or that it should not be
  /// asked at all.
  ///
  /// The cached group settings decide only WHETHER to send the request -- the one thing
  /// that has to be decided before there is a response. Everything the block renders
  /// (enabled, position, the rows) is read off the response, which stays authoritative.
  /// A group that never opted in therefore costs nothing, exactly as on desktop.
  _SummaryRequest _resolveSummaryRequest(String groupId) {
    final groupModel = Provider.of<GroupModel>(context, listen: false);

    if (!isAllGroupId(groupModel, groupId)) {
      // A real group configures itself, so the command OMITS configurationGroupId --
      // naming a different group is a 400, and naming its own is merely redundant.
      final settings =
          groupModel.getGroupReceiptSettings(int.tryParse(groupId) ?? 0);
      return (
        shouldRequest: settings?.receiptSummaryEnabled == true,
        configurationGroupId: null,
      );
    }

    // The All group has no configuration of its own, so it borrows one.
    final enabled = summaryConfigGroups(groupModel);
    final resolved =
        resolveSummaryConfigGroup(enabled, _receiptListModel.summaryConfigGroupId);
    return (shouldRequest: resolved != null, configurationGroupId: resolved?.id);
  }

  Future<void> _fetchSummary(String groupId) async {
    // Bumped BEFORE the skip branch, not only before a request: a decision not to ask
    // still supersedes whatever is in flight. Without it, leaving a summary-enabled group
    // for one without a summary lets the old group's response land after the block was
    // cleared and repaint its figures over the new group's list.
    final seq = ++_summaryRequestSeq;

    final request = _resolveSummaryRequest(groupId);
    if (!request.shouldRequest) {
      if (mounted && _summary != null) {
        setState(() => _summary = null);
      }
      return;
    }

    try {
      final response =
          await OpenApiClient.client.getReceiptApi().getReceiptSummaryForGroup(
                groupId: int.parse(groupId),
                receiptSummaryCommand: _receiptListModel.receiptSummaryCommand(
                  configurationGroupId: request.configurationGroupId,
                ),
              );
      // A response that is no longer the latest is dropped rather than painted.
      if (!mounted || seq != _summaryRequestSeq) {
        return;
      }
      setState(() => _summary = response.data);
    } catch (_) {
      // Swallowed on purpose, mirroring the desktop's catchError: an unhandled
      // DioException from a setState-driven fetch takes the whole receipts screen down,
      // and the right degradation for a block of totals is to keep the last good
      // figures rather than to interrupt someone browsing receipts.
      //
      // The group guard is released so a failure is not recorded as "this group is done".
      // In practice every navigation off this screen remounts the State anyway, so there
      // is no demonstrable case where the latch bites -- hence no test for it. It is kept
      // because a guard that means "already fetched" should not be set by a fetch that
      // did not happen.
      if (mounted && seq == _summaryRequestSeq && _summaryGroupId == groupId) {
        _summaryGroupId = null;
      }
    }
  }

  /// Changing the configuration changes the breakdown's shape, never the data -- so this
  /// refreshes the summary ONLY, and must not go through _refreshCallback.
  void _onConfigGroupSelected(int groupId) {
    _receiptListModel.setSummaryConfigGroupId(groupId, false);
    _fetchSummary(getGroupId(context));
  }

  /// The quick date control: a month stepper over the date field the chip row
  /// below points it at.
  ///
  /// Mirrors desktop's `.receipts-quick-date`, which reads left to right as
  /// "September 2026 ... on Receipt Date"; on a phone the two stack instead, so
  /// it reads top to bottom.
  Widget buildMonthStepper() {
    final month = monthFromCondition(_receiptListModel.quickDateCondition);

    return ReceiptMonthStepper(
      month: month,
      label: buildMonthStepperLabel(month),
      // Not `month != null`: a condition the stepper cannot describe is still
      // one the user needs a way out of.
      hasCondition: _receiptListModel.quickDateCondition != null,
      onMonthSelected: applyMonth,
      onAllTimeSelected: () => applyMonth(null),
      onLabelPressed: openMonthPicker,
    );
  }

  /// "September 2026" for a whole calendar month, "Custom" for a condition the
  /// stepper cannot describe, "All time" for no condition at all.
  ///
  /// A condition it cannot describe still has to read as a filter -- the
  /// stepper never hides and never silently drops one.
  String buildMonthStepperLabel(FilterMonth? month) {
    if (month != null) {
      return filterMonthLabel(month);
    }

    return _receiptListModel.quickDateCondition == null ? "All time" : "Custom";
  }

  /// Writes [month] onto the field the control is pointed at, or clears it.
  ///
  /// The quick control IS its field's filter, so it overwrites whatever was
  /// there. `setFilterField` notifies, which refetches once through
  /// [_refreshForFilterChange].
  void applyMonth(FilterMonth? month) {
    _receiptListModel.setFilterField(
      _receiptListModel.quickDateField,
      month == null ? null : monthFilterCondition(month),
      groupId: getGroupId(context),
    );
  }

  Future<void> openMonthPicker() async {
    final result = await showReceiptMonthPickerSheet(
      context,
      selected: monthFromCondition(_receiptListModel.quickDateCondition),
    );

    // A dismissal changes nothing. `applyMonth` reads the route for the group
    // id, so bail if the sheet outlived this widget.
    if (result == null || !mounted) {
      return;
    }

    applyMonth(result.isAllTime ? null : result.month);
  }

  /// Horizontally scrollable: three chips do not fit a narrow phone, and an
  /// unscrollable `Row` overflows rather than truncating.
  Widget buildSortFilterBar() {
    return SingleChildScrollView(
      scrollDirection: Axis.horizontal,
      child: Row(
        children: [
          buildQuickDateFieldChip(),
          SizedBox(
            width: 4,
          ),
          buildSortChip(),
          SizedBox(
            width: 4,
          ),
          buildSortDirectionChip()
        ],
      ),
    );
  }

  /// Which date field the stepper reads and writes.
  ///
  /// Switching is deliberately non-destructive and refetches nothing: it
  /// changes no condition, only which one the stepper describes, so a Receipt
  /// Date filter survives a move to Resolved Date and stays applied. The
  /// `setState` is only to repaint this bar and the label above it.
  ///
  /// The label is prefixed with "On" because [receiptSortOptions] names the
  /// same three columns identically -- sorting by Receipt Date would otherwise
  /// put two chips reading "Receipt Date" side by side, meaning different
  /// things. Desktop needs no prefix: its picker sits against the stepper and
  /// reads "September 2026 ... on Receipt Date" in one line, which the two rows
  /// here cannot.
  Widget buildQuickDateFieldChip() {
    final selected = _receiptListModel.quickDateField;
    final label = receiptDateFilterFields
        .firstWhere((field) => field.key == selected,
            orElse: () => receiptDateFilterFields.first)
        .label;

    return PopupMenuButton<String>(
      key: const ValueKey("receipt-quick-date-field"),
      tooltip: "Choose which date field to filter on",
      child: Chip(
        avatar: Icon(Icons.event, size: 18),
        label: Text("On $label"),
      ),
      itemBuilder: (context) => receiptDateFilterFields.map((field) {
        return PopupMenuItem(
          key: ValueKey("receipt-quick-date-field-${field.key}"),
          value: field.key,
          child: Row(
            children: [
              Icon(field.key == selected ? Icons.check : null, size: 18),
              SizedBox(width: 8),
              Text(field.label),
            ],
          ),
        );
      }).toList(),
      onSelected: (value) {
        _receiptListModel.setQuickDateField(value, false,
            groupId: getGroupId(context));
        setState(() {});
      },
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

  /// Null unless the server says the configuration group has the summary turned on.
  /// `enabled: false` arrives as a normal 200 with zeroed rows, so a client with stale
  /// group settings renders nothing rather than surfacing an error.
  Widget? _buildSummaryBar() {
    final summary = _summary;
    if (summary == null || !summary.enabled) {
      return null;
    }

    final groupId = getGroupId(context);
    final groupModel = Provider.of<GroupModel>(context, listen: false);
    // Chips only where there is a choice to make: the All group, which borrows a
    // configuration. A real group configures itself.
    final configGroups =
        isAllGroupId(groupModel, groupId) ? summaryConfigGroups(groupModel) : <api.Group>[];

    return ReceiptSummaryBar(
      summary: summary,
      atTop: isReceiptSummaryAtTop(summary.position),
      configGroups: configGroups,
      // The pick, not the last response's configurationGroupId: otherwise a tap leaves the
      // previous chip highlighted for the whole round trip, and stays there forever if the
      // request fails -- with the model and the UI silently disagreeing.
      selectedConfigGroupId:
          _receiptListModel.summaryConfigGroupId ?? summary.configurationGroupId,
      onConfigGroupSelected: _onConfigGroupSelected,
    );
  }

  @override
  Widget build(BuildContext context) {
    final bar = _buildSummaryBar();
    final atTop = bar != null && isReceiptSummaryAtTop(_summary?.position);

    return Column(
      children: [
        buildMonthStepper(),
        buildSortFilterBar(),
        // Both slots are PINNED for free: PagedDataList returns an Expanded, so it takes
        // the remaining height and the bar never scrolls with the list. That is what
        // makes the bottom position reachable at all under infinite scroll -- and it is
        // why _pagedListKey exists; see its doc comment.
        if (atTop) bar,
        PagedDataList(
          key: _pagedListKey,
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
        if (bar != null && !atTop) bar,
      ],
    );
  }
}
