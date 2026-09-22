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
import 'package:receipt_wrangler_mobile/shared/widgets/paged_data_list.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';

import '../helpers/receipt_filter_widget_helpers.dart';
import '../helpers/receipt_form_test_helpers.dart';
import '../helpers/receipt_summary_test_helpers.dart';
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
    registerFallbackValue(api.ReceiptSummaryCommandBuilder().build());
  });

  const nameCondition = ReceiptFilterCondition(
      operation: api.FilterOperation.CONTAINS, value: "Costco");

  late _MockOpenapi mockClient;
  late _MockReceiptApi mockReceiptApi;
  late List<api.ReceiptPagedRequestCommand> requests;
  late List<({int groupId, api.ReceiptSummaryCommand command})> summaryRequests;
  late api.ReceiptSummary summaryResponse;

  setUp(() {
    // Deliberately no save-and-restore of the real client: reading
    // OpenApiClient.client triggers its lazy `Openapi()` construction, whose
    // "/api" baseUrl Dio rejects off-web. Each test file gets its own isolate,
    // so swapping in a fresh mock per test is enough isolation.
    mockClient = _MockOpenapi();
    mockReceiptApi = _MockReceiptApi();
    requests = [];
    summaryRequests = [];
    summaryResponse = buildReceiptSummary(
      overall: buildSummaryRow(receiptCount: 3, total: '30.00'),
    );

    when(() => mockReceiptApi.getReceiptSummaryForGroup(
          groupId: any(named: "groupId"),
          receiptSummaryCommand: any(named: "receiptSummaryCommand"),
        )).thenAnswer((invocation) async {
      summaryRequests.add((
        groupId: invocation.namedArguments[#groupId] as int,
        command: invocation.namedArguments[#receiptSummaryCommand]
            as api.ReceiptSummaryCommand,
      ));

      return Response(
        requestOptions: RequestOptions(path: "/"),
        data: summaryResponse,
      );
    });

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

  /// A harness whose Household group has opted into the summary. The default harness
  /// leaves it off, which is what makes "no request when not opted in" the baseline that
  /// every other test in this file already asserts by never stubbing the endpoint.
  ReceiptFilterHarness harnessWithSummary({
    bool householdEnabled = true,
    bool officeEnabled = false,
    bool allGroupEnabled = false,
  }) {
    final harness = buildReceiptFilterHarness();
    harness.groupModel.setGroups([
      buildGroup(
          id: ReceiptFilterHarness.allGroupId,
          name: "All",
          isAllGroup: true,
          receiptSummaryEnabled: allGroupEnabled),
      buildGroup(
          id: ReceiptFilterHarness.householdId,
          name: "Household",
          receiptSummaryEnabled: householdEnabled),
      buildGroup(
          id: ReceiptFilterHarness.officeId,
          name: "Office",
          receiptSummaryEnabled: officeEnabled),
    ]);
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

  /// The summary rides the same refresh contract as the desktop's, and the point of it
  /// is what it does NOT do: sorting and paging change neither the filter nor the
  /// figures, so an unpaged aggregate must not be refetched for either.
  group('the receipt summary', () {
    String householdRoute() =>
        "/groups/${ReceiptFilterHarness.householdId}/receipts";

    testWidgets('is fetched once on mount, with an empty filter', (tester) async {
      await pumpList(tester,
          router: routerFor(householdRoute()), harness: harnessWithSummary());

      expect(summaryRequests, hasLength(1));
      expect(summaryRequests.single.groupId, ReceiptFilterHarness.householdId);
      // A real group configures itself, so the key is OMITTED -- naming another group
      // is a 400 and naming its own is merely redundant.
      expect(summaryRequests.single.command.configurationGroupId, isNull);
      expect(find.byKey(const ValueKey('receipt-summary')), findsOneWidget);
    });

    testWidgets('is not fetched for a group that never opted in', (tester) async {
      await pumpList(tester,
          router: routerFor(householdRoute()),
          harness: harnessWithSummary(householdEnabled: false));

      expect(summaryRequests, isEmpty);
      expect(find.byKey(const ValueKey('receipt-summary')), findsNothing);
    });

    testWidgets('applying a filter refetches it, once, with the conditions',
        (tester) async {
      final harness = await pumpList(tester,
          router: routerFor(householdRoute()), harness: harnessWithSummary());
      summaryRequests.clear();

      harness.receiptListModel.setFilter({"name": nameCondition}, true,
          groupId: "${ReceiptFilterHarness.householdId}");
      await tester.pumpAndSettle();

      expect(summaryRequests, hasLength(1));
      final filter = Map<String, dynamic>.from(api.standardSerializers.serializeWith(
          api.ReceiptPagedRequestFilter.serializer,
          summaryRequests.single.command.filter!) as Map);
      expect(filter, {
        "name": {"operation": "CONTAINS", "value": "Costco"}
      });
    });

    // The single most valuable case here. The sort setters pass notify: false and call
    // the list's own refresh callback, so the summary must not be on that path at all.
    testWidgets('a sort change refetches the list but NOT the summary', (tester) async {
      await pumpList(tester,
          router: routerFor(householdRoute()), harness: harnessWithSummary());
      requests.clear();
      summaryRequests.clear();

      await tester.tap(find.text("Added At"));
      await tester.pumpAndSettle();
      await tester.tap(find.text("Sort by Amount"));
      await tester.pumpAndSettle();

      expect(requests, hasLength(1), reason: "the list still refetches");
      expect(summaryRequests, isEmpty,
          reason: "sorting changes the order of the result set, not its membership");
    });

    testWidgets('a group change refetches it for the new group', (tester) async {
      final router = routerFor(householdRoute());
      await pumpList(tester,
          router: router,
          harness: harnessWithSummary(officeEnabled: true));
      summaryRequests.clear();

      router.go("/groups/${ReceiptFilterHarness.officeId}/receipts");
      await tester.pumpAndSettle();

      expect(summaryRequests, hasLength(1));
      expect(summaryRequests.single.groupId, ReceiptFilterHarness.officeId);
    });

    testWidgets('renders nothing for an enabled: false response', (tester) async {
      summaryResponse = buildReceiptSummary(enabled: false);
      await pumpList(tester,
          router: routerFor(householdRoute()), harness: harnessWithSummary());

      // The request still goes out -- gating on the client's cached settings would render
      // a stale block when an admin has just changed the configuration.
      expect(summaryRequests, hasLength(1));
      expect(find.byKey(const ValueKey('receipt-summary')), findsNothing);
    });

    testWidgets('a failed summary leaves the list working', (tester) async {
      when(() => mockReceiptApi.getReceiptSummaryForGroup(
            groupId: any(named: "groupId"),
            receiptSummaryCommand: any(named: "receiptSummaryCommand"),
          )).thenThrow(DioException(requestOptions: RequestOptions(path: "/")));

      await pumpList(tester,
          router: routerFor(householdRoute()), harness: harnessWithSummary());

      expect(tester.takeException(), isNull);
      expect(find.byType(PagedDataList), findsOneWidget);
      expect(find.byKey(const ValueKey('receipt-summary')), findsNothing);
    });

    group('placement', () {
      Future<ReceiptFilterHarness> pumpAt(
          WidgetTester tester, api.ReceiptSummaryPosition position) async {
        summaryResponse = buildReceiptSummary(
          position: position,
          overall: buildSummaryRow(receiptCount: 3, total: '30.00'),
        );
        return pumpList(tester,
            router: routerFor(householdRoute()), harness: harnessWithSummary());
      }

      // Geometry, not a slot key: it tests the thing rather than a label for it.
      testWidgets('BOTTOM puts the bar under the list', (tester) async {
        await pumpAt(tester, api.ReceiptSummaryPosition.BOTTOM);

        expect(
          tester.getTopLeft(find.byKey(const ValueKey('receipt-summary'))).dy,
          greaterThan(tester.getTopLeft(find.byType(PagedDataList)).dy),
        );
      });

      testWidgets('TOP puts the bar above the list', (tester) async {
        await pumpAt(tester, api.ReceiptSummaryPosition.TOP);

        expect(
          tester.getTopLeft(find.byKey(const ValueKey('receipt-summary'))).dy,
          lessThan(tester.getTopLeft(find.byType(PagedDataList)).dy),
        );
      });

      /// The four-slot Column guard. Column matches children by index and runtime type,
      /// so moving one bar between slots would shift PagedDataList's index, fail
      /// Widget.canUpdate and silently discard its State -- losing every loaded page and
      /// refetching page 1. This fails against that implementation.
      testWidgets('flipping the position does not reset the paged list',
          (tester) async {
        final harness = await pumpAt(tester, api.ReceiptSummaryPosition.TOP);
        final listRequestsBefore = requests.length;
        final pagedStateBefore = tester.state(find.byType(PagedDataList));

        summaryResponse = buildReceiptSummary(
          position: api.ReceiptSummaryPosition.BOTTOM,
          overall: buildSummaryRow(receiptCount: 3, total: '30.00'),
        );
        // A filter apply is the cheapest way to drive a fresh summary response through
        // the real path; it refetches the list exactly once, which is accounted for.
        harness.receiptListModel.setFilter({"name": nameCondition}, true,
            groupId: "${ReceiptFilterHarness.householdId}");
        await tester.pumpAndSettle();

        expect(
          tester.getTopLeft(find.byKey(const ValueKey('receipt-summary'))).dy,
          greaterThan(tester.getTopLeft(find.byType(PagedDataList)).dy),
          reason: 'the bar moved to the other slot',
        );
        // State IDENTITY, not a request count: a discarded State is reconstructed
        // immediately and its refetch is easy to mistake for the filter's own.
        expect(identical(tester.state(find.byType(PagedDataList)), pagedStateBefore),
            isTrue,
            reason: 'PagedDataList kept its index, so it kept its State -- and with it '
                'the paging controller, every loaded page and _totalCount');
        expect(requests.length, listRequestsBefore + 1,
            reason: 'the filter refetch, and nothing else');
      });
    });

    group('on the All group', () {
      String allRoute() => "/groups/${ReceiptFilterHarness.allGroupId}/receipts";

      testWidgets('borrows the first enabled group\'s configuration', (tester) async {
        await pumpList(tester,
            router: routerFor(allRoute()),
            harness: harnessWithSummary(officeEnabled: true));

        expect(summaryRequests, hasLength(1));
        expect(summaryRequests.single.groupId, ReceiptFilterHarness.allGroupId);
        // Household sorts before Office.
        expect(summaryRequests.single.command.configurationGroupId,
            ReceiptFilterHarness.householdId);
      });

      testWidgets('asks for nothing when no member group has one', (tester) async {
        await pumpList(tester,
            router: routerFor(allRoute()),
            harness: harnessWithSummary(householdEnabled: false));

        expect(summaryRequests, isEmpty);
      });

      testWidgets('a chip pick refetches the summary only', (tester) async {
        await pumpList(tester,
            router: routerFor(allRoute()),
            harness: harnessWithSummary(officeEnabled: true));
        requests.clear();
        summaryRequests.clear();

        await tester.tap(
            find.byKey(const ValueKey('receipt-summary-config-group-'
                '${ReceiptFilterHarness.officeId}')));
        await tester.pumpAndSettle();

        expect(summaryRequests, hasLength(1));
        expect(summaryRequests.single.command.configurationGroupId,
            ReceiptFilterHarness.officeId);
        expect(requests, isEmpty,
            reason: 'a configuration pick changes the breakdown shape, not the data');
      });
    });
  });
}
