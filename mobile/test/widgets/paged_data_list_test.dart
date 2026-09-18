import 'package:built_collection/built_collection.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:one_of/any_of.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/shared/widgets/paged_data_list.dart';

import '../helpers/receipt_form_test_helpers.dart';

/// Regression guard for the refresh path the receipt filter leans on.
///
/// `_totalCount` is an instance field that gates paging (`getNextPageKey`
/// returns null once the loaded items reach it). It used to survive a refresh,
/// so a page that came back with `totalCount: 0` left `0 >= 0` permanently true
/// and NO page was ever requested again. Sorting can never reach a zero total,
/// which is why this only surfaces once the list can be filtered.
void main() {
  api.PagedDataDataInner wrapUser(int id) =>
      (api.PagedDataDataInnerBuilder()
            ..anyOf = AnyOf1<api.UserView>(
                value: buildUserView(id: id, displayName: "User $id")))
          .build();

  api.PagedData page(List<int> ids, int totalCount) =>
      (api.PagedDataBuilder()
            ..data = ListBuilder<api.PagedDataDataInner>(ids.map(wrapUser))
            ..totalCount = totalCount)
          .build();

  Response<api.PagedData> respond(api.PagedData data) =>
      Response(requestOptions: RequestOptions(path: "/"), data: data);

  testWidgets(
      "a refresh after an empty result can load rows again",
      (tester) async {
    var responses = <Response<api.PagedData>>[
      respond(page(const [], 0)),
    ];
    VoidCallback? refresh;

    await tester.pumpWidget(MaterialApp(
      home: Scaffold(
        body: Column(children: [
          PagedDataList(
            onRefreshCallbackSet: (callback) => refresh = callback,
            noItemsFoundText: "No receipts found",
            listItemBuilder: (context, item, index) => Text(
                (item.anyOf.values.values.whereType<api.UserView>().first).displayName),
            getPagedDataFuture: (_) async => responses.removeAt(0),
          ),
        ]),
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.text("No receipts found"), findsOneWidget);

    // The filter is loosened: the next fetch has rows again.
    responses = [respond(page(const [1, 2], 2))];
    refresh!();
    await tester.pumpAndSettle();

    expect(find.text("User 1"), findsOneWidget,
        reason: "a stale totalCount of 0 would stop the list ever paging again");
    expect(find.text("User 2"), findsOneWidget);
    expect(find.text("No receipts found"), findsNothing);
  });

  testWidgets("a refresh between two non-empty results swaps the rows",
      (tester) async {
    var responses = <Response<api.PagedData>>[
      respond(page(const [1], 1)),
    ];
    VoidCallback? refresh;

    await tester.pumpWidget(MaterialApp(
      home: Scaffold(
        body: Column(children: [
          PagedDataList(
            onRefreshCallbackSet: (callback) => refresh = callback,
            noItemsFoundText: "No receipts found",
            listItemBuilder: (context, item, index) => Text(
                (item.anyOf.values.values.whereType<api.UserView>().first).displayName),
            getPagedDataFuture: (_) async => responses.removeAt(0),
          ),
        ]),
      ),
    ));
    await tester.pumpAndSettle();

    expect(find.text("User 1"), findsOneWidget);

    responses = [respond(page(const [9], 1))];
    refresh!();
    await tester.pumpAndSettle();

    expect(find.text("User 9"), findsOneWidget);
    expect(find.text("User 1"), findsNothing);
  });
}
