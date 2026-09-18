import 'package:built_collection/built_collection.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/groups/widgets/group_receipts_list.dart';
import 'package:receipt_wrangler_mobile/models/category_model.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt-list-model.dart';
import 'package:receipt_wrangler_mobile/models/system_settings_model.dart';
import 'package:receipt_wrangler_mobile/models/tag_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

import '../helpers/receipt_filter_widget_helpers.dart';
import '../helpers/widget_test_helpers.dart';

class _MockOpenapi extends Mock implements api.Openapi {}

class _MockReceiptApi extends Mock implements api.ReceiptApi {}

/// The list is the only thing that knows how to refetch, and it learns that the
/// filter changed from a ReceiptListModel notification. These cases pin that
/// wiring plus the group-change reset, both of which are invisible until they
/// go wrong.
void main() {
  setUpAll(() {
    registerCustomCurrencyForTests();
    registerFallbackValue((api.ReceiptPagedRequestCommandBuilder()
          ..page = 1
          ..pageSize = 10)
        .build());
  });

  const nameCondition = ReceiptFilterCondition(
      operation: api.FilterOperation.CONTAINS, value: "Costco");

  late _MockOpenapi mockClient;
  late _MockReceiptApi mockReceiptApi;
  late List<api.ReceiptPagedRequestCommand> requests;

  setUp(() {
    // Deliberately no save-and-restore of the real client: reading
    // OpenApiClient.client triggers its lazy `Openapi()` construction, whose
    // "/api" baseUrl Dio rejects off-web. Each test file gets its own isolate,
    // so swapping in a fresh mock per test is enough isolation.
    mockClient = _MockOpenapi();
    mockReceiptApi = _MockReceiptApi();
    requests = [];

    when(() => mockClient.getReceiptApi()).thenReturn(mockReceiptApi);
    when(() => mockReceiptApi.getReceiptsForGroup(
          groupId: any(named: "groupId"),
          receiptPagedRequestCommand:
              any(named: "receiptPagedRequestCommand"),
        )).thenAnswer((invocation) async {
      requests.add(invocation.namedArguments[#receiptPagedRequestCommand]
          as api.ReceiptPagedRequestCommand);

      return Response(
        requestOptions: RequestOptions(path: "/"),
        data: (api.PagedDataBuilder()
              ..data = ListBuilder<api.PagedDataDataInner>()
              ..totalCount = 0)
            .build(),
      );
    });

    OpenApiClient.client = mockClient;
  });

  tearDown(() => OpenApiClient.client = _MockOpenapi());

  Map<String, dynamic> filterOf(api.ReceiptPagedRequestCommand command) =>
      Map<String, dynamic>.from(api.standardSerializers.serializeWith(
          api.ReceiptPagedRequestFilter.serializer, command.filter!) as Map);

  Future<ReceiptFilterHarness> pumpList(
    WidgetTester tester, {
    required GoRouter router,
    /// Pass one to seed an applied filter *before* the list mounts -- the only
    /// way to reach the "mounted straight into another group" case, which is
    /// how every real group switch arrives (see the group-change tests below).
    ReceiptFilterHarness? harness,
  }) async {
    harness ??= buildReceiptFilterHarness();

    await tester.pumpWidget(MultiProvider(
      providers: [
        ChangeNotifierProvider<ReceiptListModel>.value(
            value: harness.receiptListModel),
        ChangeNotifierProvider<GroupModel>.value(value: harness.groupModel),
        ChangeNotifierProvider<CategoryModel>.value(
            value: harness.categoryModel),
        ChangeNotifierProvider<TagModel>.value(value: harness.tagModel),
        ChangeNotifierProvider<UserModel>.value(value: harness.userModel),
        ChangeNotifierProvider<ContextModel>(create: (_) => ContextModel()),
        ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
        ChangeNotifierProvider<PermissionsModel>(
            create: (_) => PermissionsModel()),
        ChangeNotifierProvider<SystemSettingsModel>(
            create: (_) => SystemSettingsModel()),
      ],
      child: MaterialApp.router(routerConfig: router),
    ));
    await tester.pumpAndSettle();

    return harness;
  }

  GoRouter routerFor(String location) => GoRouter(
        initialLocation: location,
        routes: [
          GoRoute(
            path: "/groups/:groupId/receipts",
            builder: (context, state) =>
                const Scaffold(body: GroupReceiptsList()),
          ),
        ],
      );

  testWidgets("the first fetch carries an empty filter", (tester) async {
    // An unfiltered list must send the request it always did.
    await pumpList(tester,
        router: routerFor("/groups/${ReceiptFilterHarness.householdId}/receipts"));

    expect(requests, hasLength(1));
    expect(filterOf(requests.single), isEmpty);
  });

  testWidgets("applying a filter refetches, once, with the conditions",
      (tester) async {
    final harness = await pumpList(tester,
        router: routerFor("/groups/${ReceiptFilterHarness.householdId}/receipts"));
    requests.clear();

    harness.receiptListModel.setFilter({"name": nameCondition}, true, groupId: "${ReceiptFilterHarness.householdId}");
    await tester.pumpAndSettle();

    expect(requests, hasLength(1),
        reason: "one notification must mean one refetch");
    expect(filterOf(requests.single),
        {"name": {"operation": "CONTAINS", "value": "Costco"}});
  });

  testWidgets("a silent filter write does not refetch", (tester) async {
    // notify: false is how the group reset clears without a rebuild mid-frame.
    final harness = await pumpList(tester,
        router: routerFor("/groups/${ReceiptFilterHarness.householdId}/receipts"));
    requests.clear();

    harness.receiptListModel.setFilter({"name": nameCondition}, false, groupId: "${ReceiptFilterHarness.householdId}");
    await tester.pumpAndSettle();

    expect(requests, isEmpty);
  });

  testWidgets("a sort change still refetches through its own path",
      (tester) async {
    await pumpList(tester,
        router: routerFor("/groups/${ReceiptFilterHarness.householdId}/receipts"));
    requests.clear();

    await tester.tap(find.text("Added At"));
    await tester.pumpAndSettle();
    await tester.tap(find.text("Sort by Amount"));
    await tester.pumpAndSettle();

    expect(requests, hasLength(1),
        reason: "the sort setters pass notify: false, so this must not "
            "also fire the filter listener");
    expect(requests.single.orderBy, "amount");
  });

  testWidgets("changing group clears the filter and refetches", (tester) async {
    final router =
        routerFor("/groups/${ReceiptFilterHarness.householdId}/receipts");
    final harness = await pumpList(tester, router: router);

    harness.receiptListModel.setFilter({"name": nameCondition}, true, groupId: "${ReceiptFilterHarness.householdId}");
    await tester.pumpAndSettle();
    requests.clear();

    router.go("/groups/${ReceiptFilterHarness.officeId}/receipts");
    await tester.pumpAndSettle();

    // A filter holding the previous group's category, tag and user ids matches
    // nothing here, and would leave the badge counting invisible conditions.
    expect(harness.receiptListModel.hasActiveFilter, isFalse);
    expect(requests, isNotEmpty);
    expect(filterOf(requests.last), isEmpty);
  });

  testWidgets(
      "a list mounted into another group clears the filter before it fetches",
      (tester) async {
    // The case the shipped guard missed. The app offers no lateral group
    // switch -- every real one goes out through /groups and back in, which
    // destroys this widget's State. So the list arrives in the new group with
    // no memory of the old one, and a widget-local "last group I saw" is null
    // exactly when the clear is needed. The filter's own scope is what makes
    // this reachable.
    final harness = buildReceiptFilterHarness();
    harness.receiptListModel.setFilter({"name": nameCondition}, false,
        groupId: "${ReceiptFilterHarness.householdId}");

    await pumpList(tester,
        harness: harness,
        router: routerFor("/groups/${ReceiptFilterHarness.officeId}/receipts"));

    expect(harness.receiptListModel.hasActiveFilter, isFalse);
    expect(requests, hasLength(1),
        reason: "the clear must land before the first fetch, not cause a "
            "second one");
    expect(filterOf(requests.single), isEmpty);
  });

  testWidgets("a list remounted into the SAME group keeps the filter",
      (tester) async {
    // The mirror image, and why the scope cannot simply be "clear on mount":
    // a round trip to a receipt tears this list down and rebuilds it too.
    final harness = buildReceiptFilterHarness();
    harness.receiptListModel.setFilter({"name": nameCondition}, false,
        groupId: "${ReceiptFilterHarness.householdId}");

    await pumpList(tester,
        harness: harness,
        router:
            routerFor("/groups/${ReceiptFilterHarness.householdId}/receipts"));

    expect(harness.receiptListModel.hasActiveFilter, isTrue);
    expect(filterOf(requests.single),
        {"name": {"operation": "CONTAINS", "value": "Costco"}});
  });

  testWidgets("the empty state names the filter when one is applied",
      (tester) async {
    final harness = await pumpList(tester,
        router: routerFor("/groups/${ReceiptFilterHarness.householdId}/receipts"));

    expect(find.text("No receipts found"), findsOneWidget);

    harness.receiptListModel.setFilter({"name": nameCondition}, true, groupId: "${ReceiptFilterHarness.householdId}");
    await tester.pumpAndSettle();

    expect(find.text("No receipts match this filter"), findsOneWidget);
  });
}
