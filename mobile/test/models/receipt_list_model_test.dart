import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/models/receipt-list-model.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

import '../helpers/receipt_filter_test_helpers.dart';

/// The model owns the *applied* filter -- the filter screen edits its own draft
/// and commits through `setFilter`. Two contracts here are load-bearing beyond
/// this file: the request the list sends, and which mutations notify (the
/// receipts list treats a notification as "the filter changed").
void main() {
  late ReceiptListModel model;

  setUp(() => model = ReceiptListModel());

  const nameCondition = ReceiptFilterCondition(
      operation: api.FilterOperation.CONTAINS, value: "Costco");
  const amountCondition = ReceiptFilterCondition(
      operation: api.FilterOperation.GREATER_THAN, value: 50.0);

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
      model.setFilter({"name": nameCondition, "amount": amountCondition}, false);

      expect(model.activeFilterCount, 2);
      expect(model.hasActiveFilter, isTrue);
      expect(model.filter.keys, containsAll(["name", "amount"]));
    });

    test("the conditions reach the command", () {
      model.setFilter({"name": nameCondition}, false);

      expect(serializedFilterOf(model),
          {"name": {"operation": "CONTAINS", "value": "Costco"}});
    });

    test("notifies when asked to", () {
      var notifications = 0;
      model.addListener(() => notifications++);

      model.setFilter({"name": nameCondition}, true);

      expect(notifications, 1);
    });

    test("stays silent when not", () {
      var notifications = 0;
      model.addListener(() => notifications++);

      model.setFilter({"name": nameCondition}, false);

      expect(notifications, 0);
    });

    test("keeps its own copy of the map it is handed", () {
      final conditions = {"name": nameCondition};
      model.setFilter(conditions, false);

      conditions["amount"] = amountCondition;

      expect(model.activeFilterCount, 1);
    });

    test("replaces rather than merges", () {
      model.setFilter({"name": nameCondition}, false);
      model.setFilter({"amount": amountCondition}, false);

      expect(model.filter.keys, ["amount"]);
    });
  });

  test("the applied filter cannot be mutated from outside", () {
    // Otherwise the badge count and the query could disagree.
    model.setFilter({"name": nameCondition}, false);

    expect(() => model.filter["amount"] = amountCondition, throwsUnsupportedError);
  });

  group("clearFilter", () {
    test("empties the filter and the command", () {
      model.setFilter({"name": nameCondition}, false);
      model.clearFilter(false);

      expect(model.filter, isEmpty);
      expect(model.activeFilterCount, 0);
      expect(serializedFilterOf(model), isEmpty);
    });

    test("notifies when asked to", () {
      model.setFilter({"name": nameCondition}, false);

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
      model.setFilter({"name": nameCondition}, false);

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

      model.setFilter({"name": nameCondition}, false);

      expect(model.orderBy, "name");
      expect(model.page, 7);
    });

    test("a fresh filter builder is produced per request", () {
      // A built_value nested builder assigned into a command builder is a live
      // reference, so a cached one would be shared between requests.
      model.setFilter({"name": nameCondition}, false);

      final first = model.receiptPagedRequestCommand;
      final second = model.receiptPagedRequestCommand;

      expect(identical(first.filter, second.filter), isFalse);
      expect(first.filter, second.filter);
    });

    test("an invalid condition never reaches the request", () {
      model.setFilter({
        "date": const ReceiptFilterCondition(
            operation: api.FilterOperation.EQUALS)
      }, false);

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
    }, false);

    expect(serializedFilterOf(model)["categories"],
        {"operation": "CONTAINS", "value": [3, 8]});
  });
}
