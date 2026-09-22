import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/constants/receipt_filter_fields.dart';
import 'package:receipt_wrangler_mobile/models/receipt-list-model.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_date_filter.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

import '../helpers/receipt_filter_test_helpers.dart';

/// The model owns the *applied* filter -- the filter screen edits its own draft
/// and commits through `setFilter`. Two contracts here are load-bearing beyond
/// this file: the request the list sends, and which mutations notify (the
/// receipts list treats a notification as "the filter changed").
void main() {
  late ReceiptListModel model;

  setUp(() => model = ReceiptListModel());

  /// The group every filter in this file is authored in. `setFilter` requires
  /// one so a filter always knows which group it is meaningful for.
  const groupId = "2";

  const nameCondition = ReceiptFilterCondition(
      operation: api.FilterOperation.CONTAINS, value: "Costco");
  const amountCondition = ReceiptFilterCondition(
      operation: api.FilterOperation.GREATER_THAN, value: 50.0);
  final monthCondition = monthFilterCondition(const FilterMonth(2026, 9));

  Map<String, dynamic> serializedFilterOf(ReceiptListModel model) {
    final command = model.receiptPagedRequestCommand;
    return Map<String, dynamic>.from(api.standardSerializers.serializeWith(
        api.ReceiptPagedRequestFilter.serializer, command.filter!) as Map);
  }

  group("defaults", () {
    test("starts unfiltered", () {
      expect(model.filter, isEmpty);
      expect(model.activeFilterCount, 0);
      expect(model.hasActiveFilter, isFalse);
    });

    test("an unfiltered command carries an empty filter", () {
      // The request an unfiltered list sends must not change.
      expect(serializedFilterOf(model), isEmpty);
    });
  });

  group("setFilter", () {
    test("stores the conditions and counts them", () {
      model.setFilter({"name": nameCondition, "amount": amountCondition}, false, groupId: groupId);

      expect(model.activeFilterCount, 2);
      expect(model.hasActiveFilter, isTrue);
      expect(model.filter.keys, containsAll(["name", "amount"]));
    });

    test("the conditions reach the command", () {
      model.setFilter({"name": nameCondition}, false, groupId: groupId);

      expect(serializedFilterOf(model),
          {"name": {"operation": "CONTAINS", "value": "Costco"}});
    });

    test("notifies when asked to", () {
      var notifications = 0;
      model.addListener(() => notifications++);

      model.setFilter({"name": nameCondition}, true, groupId: groupId);

      expect(notifications, 1);
    });

    test("stays silent when not", () {
      var notifications = 0;
      model.addListener(() => notifications++);

      model.setFilter({"name": nameCondition}, false, groupId: groupId);

      expect(notifications, 0);
    });

    test("keeps its own copy of the map it is handed", () {
      final conditions = {"name": nameCondition};
      model.setFilter(conditions, false, groupId: groupId);

      conditions["amount"] = amountCondition;

      expect(model.activeFilterCount, 1);
    });

    test("replaces rather than merges", () {
      model.setFilter({"name": nameCondition}, false, groupId: groupId);
      model.setFilter({"amount": amountCondition}, false, groupId: groupId);

      expect(model.filter.keys, ["amount"]);
    });
  });

  group("the filter remembers the group it was authored in", () {
    // The receipts list is rebuilt from scratch on almost every navigation, so
    // it cannot tell a group change from a round trip to a receipt by memory
    // alone. It asks the filter instead -- which is only possible because the
    // filter records its own scope here.
    test("setFilter records the group", () {
      model.setFilter({"name": nameCondition}, false, groupId: "7");

      expect(model.filterGroupId, "7");
    });

    test("an empty filter is scoped to no group", () {
      // Applying an emptied draft is how the screen's Reset commits. It has to
      // leave no scope behind, or the next group change would clear a filter
      // that is not there and refetch for nothing.
      model.setFilter({"name": nameCondition}, false, groupId: "7");
      model.setFilter({}, false, groupId: "7");

      expect(model.filterGroupId, isNull);
    });

    test("clearFilter drops the group with the conditions", () {
      model.setFilter({"name": nameCondition}, false, groupId: "7");
      model.clearFilter(false);

      expect(model.filterGroupId, isNull);
    });

    test("re-applying in another group re-scopes it", () {
      model.setFilter({"name": nameCondition}, false, groupId: "7");
      model.setFilter({"amount": amountCondition}, false, groupId: "9");

      expect(model.filterGroupId, "9");
    });
  });

  test("the applied filter cannot be mutated from outside", () {
    // Otherwise the badge count and the query could disagree.
    model.setFilter({"name": nameCondition}, false, groupId: groupId);

    expect(() => model.filter["amount"] = amountCondition, throwsUnsupportedError);
  });

  group("clearFilter", () {
    test("empties the filter and the command", () {
      model.setFilter({"name": nameCondition}, false, groupId: groupId);
      model.clearFilter(false);

      expect(model.filter, isEmpty);
      expect(model.activeFilterCount, 0);
      expect(serializedFilterOf(model), isEmpty);
    });

    test("notifies when asked to", () {
      model.setFilter({"name": nameCondition}, false, groupId: groupId);

      var notifications = 0;
      model.addListener(() => notifications++);
      model.clearFilter(true);

      expect(notifications, 1);
    });

    test("does nothing, and notifies nobody, when already empty", () {
      // It runs on every group change, most of which arrive already unfiltered.
      var notifications = 0;
      model.addListener(() => notifications++);

      model.clearFilter(true);

      expect(notifications, 0);
    });
  });

  group("the sort setters stay silent", () {
    // The receipts list refreshes itself after a sort change and listens to this
    // model only for the filter, so a notification from here means the filter
    // moved. Making these notify would refetch the list twice per sort.
    test("setOrderBy", () {
      var notifications = 0;
      model.addListener(() => notifications++);

      model.setOrderBy("amount", false);

      expect(notifications, 0);
      expect(model.orderBy, "amount");
    });

    test("setSortDirection", () {
      var notifications = 0;
      model.addListener(() => notifications++);

      model.setSortDirection(api.SortDirection.asc, false);

      expect(notifications, 0);
      expect(model.sortDirection, api.SortDirection.asc);
    });

    test("setPage", () {
      var notifications = 0;
      model.addListener(() => notifications++);

      model.setPage(3, false);

      expect(notifications, 0);
      expect(model.page, 3);
    });
  });

  group("the command composes every part", () {
    test("filter, sort, direction and page do not clobber each other", () {
      model.setOrderBy("amount", false);
      model.setSortDirection(api.SortDirection.asc, false);
      model.setPage(4, false);
      model.setFilter({"name": nameCondition}, false, groupId: groupId);

      final command = model.receiptPagedRequestCommand;

      expect(command.orderBy, "amount");
      expect(command.sortDirection, api.SortDirection.asc);
      expect(command.page, 4);
      expect(command.pageSize, 10);
      expect(serializedFilterOf(model).keys, ["name"]);
    });

    test("setting the filter leaves paging and sorting alone", () {
      // The list drives paging itself -- PagingController.refresh restarts from
      // the first page and setPage records it.
      model.setOrderBy("name", false);
      model.setPage(7, false);

      model.setFilter({"name": nameCondition}, false, groupId: groupId);

      expect(model.orderBy, "name");
      expect(model.page, 7);
    });

    test("a fresh filter builder is produced per request", () {
      // A built_value nested builder assigned into a command builder is a live
      // reference, so a cached one would be shared between requests.
      model.setFilter({"name": nameCondition}, false, groupId: groupId);

      final first = model.receiptPagedRequestCommand;
      final second = model.receiptPagedRequestCommand;

      expect(identical(first.filter, second.filter), isFalse);
      expect(first.filter, second.filter);
    });

    test("an invalid condition never reaches the request", () {
      model.setFilter({
        "date": const ReceiptFilterCondition(
            operation: api.FilterOperation.EQUALS)
      }, false, groupId: groupId);

      // It still counts toward the badge -- the user authored it -- but an empty
      // date value would become `date = ''` server-side and match nothing.
      expect(model.activeFilterCount, 1);
      expect(serializedFilterOf(model), isEmpty);
    });
  });

  test("a selection condition survives the round trip to the command", () {
    model.setFilter({
      "categories": ReceiptFilterCondition(
          operation: api.FilterOperation.CONTAINS,
          value: [buildCategory(3, "Fuel"), buildCategory(8, "Dining")])
    }, false, groupId: groupId);

    expect(serializedFilterOf(model)["categories"],
        {"operation": "CONTAINS", "value": [3, 8]});
  });

  group("setFilterField", () {
    test("adds a condition to an empty filter and records the group", () {
      model.setFilterField("date", monthCondition, groupId: groupId);

      expect(model.filter.keys, ["date"]);
      expect(model.filterGroupId, groupId);
    });

    test("replaces one field without disturbing the others", () {
      model.setFilter({"name": nameCondition}, false, groupId: groupId);

      model.setFilterField("date", monthCondition, groupId: groupId);

      expect(model.filter.keys, containsAll(["name", "date"]));
      expect(model.filter["name"], nameCondition);
    });

    test("a null condition removes the field", () {
      model.setFilter({"name": nameCondition, "date": monthCondition}, false,
          groupId: groupId);

      model.setFilterField("date", null, groupId: groupId);

      expect(model.filter.keys, ["name"]);
    });

    test("removing the last condition drops the group scope with it", () {
      model.setFilterField("date", monthCondition, groupId: groupId);

      model.setFilterField("date", null, groupId: groupId);

      expect(model.filter, isEmpty);
      expect(model.filterGroupId, isNull);
    });

    test("notifies, because the result set moved", () {
      var notifications = 0;
      model.addListener(() => notifications++);

      model.setFilterField("date", monthCondition, groupId: groupId);

      expect(notifications, 1);
    });
  });

  group("the quick date field", () {
    test("starts on Receipt Date", () {
      expect(model.quickDateField, defaultQuickDateFieldKey);
      expect(model.quickDateCondition, isNull);
    });

    test("reads the condition on whichever field it points at", () {
      model.setFilter({
        "date": nameCondition,
        "resolvedDate": monthCondition,
      }, false, groupId: groupId);

      model.setQuickDateField("resolvedDate", false);

      expect(model.quickDateCondition, monthCondition);
    });

    test("switching it disturbs no condition", () {
      // Non-destructive by design: it changes which condition the stepper
      // describes, not the filter itself.
      model.setFilter({"date": monthCondition}, false, groupId: groupId);

      model.setQuickDateField("createdAt", false);

      expect(model.filter.keys, ["date"]);
      expect(model.quickDateCondition, isNull);
    });

    test("honours its notify flag", () {
      // The receipts list treats a notification as "the filter changed", so a
      // notifying switch would refetch a result set that has not moved.
      var notifications = 0;
      model.addListener(() => notifications++);

      model.setQuickDateField("resolvedDate", false);
      expect(notifications, 0);

      model.setQuickDateField("createdAt", true);
      expect(notifications, 1);
    });

    test("is reset by clearFilter, alongside the conditions", () {
      model.setFilter({"resolvedDate": monthCondition}, false,
          groupId: groupId);
      model.setQuickDateField("resolvedDate", false);

      model.clearFilter(false);

      expect(model.filter, isEmpty);
      expect(model.quickDateField, defaultQuickDateFieldKey);
    });

    test("is reset even when no condition was ever applied", () {
      // The early return has to cover the field too, or a group change leaves
      // the stepper pointed at a field chosen in the group just left.
      model.setQuickDateField("createdAt", false);

      model.clearFilter(false);

      expect(model.quickDateField, defaultQuickDateFieldKey);
    });
  });
}
