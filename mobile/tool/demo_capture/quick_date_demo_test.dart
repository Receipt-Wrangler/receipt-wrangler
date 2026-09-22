import 'package:built_collection/built_collection.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:image/image.dart' as img;
import 'package:mocktail/mocktail.dart';
import 'package:one_of/any_of.dart';
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
import 'package:receipt_wrangler_mobile/theme/app_theme.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_date_filter.dart';

import '../../test/helpers/receipt_filter_widget_helpers.dart';
import '../../test/helpers/widget_test_helpers.dart';
import 'capture.dart';

/// Records the demo GIF for the quick date filter — the month stepper and the
/// date-field picker above the receipts list.
///
/// **Not** part of the test suite. `flutter test` with no arguments scans only
/// `test/`, so CI never runs this and never rewrites the committed GIF; an
/// explicit path is used verbatim, which is how it is invoked:
///
/// ```bash
/// cd mobile && ./tool/record_quick_date_demo.sh
/// ```
///
/// **One panel, not a before/after pair.** The other demos in this directory
/// record the same screen twice across a `debugDisable*` seam, because they
/// document *fixes* — a bug on the left, its absence on the right. This is a
/// new control, so there is no seam to flip and nothing to compare against;
/// what the viewer needs to see is the interaction itself.
///
/// **Route 1 could not record this.** Driving the real control needs the Go
/// API behind it, which the sandbox cannot bring up (ImageMagick 7 from source
/// — see the root `CLAUDE.md`). Here the list is the **real**
/// `GroupReceiptsList` with the real stepper, sheet and picker, and only the
/// receipts endpoint is stubbed — by a fake that decodes the filter the client
/// actually sent and returns the rows matching it. So the rows in the GIF
/// narrow because the encoded `BETWEEN` selected them, not because a script
/// said so: break the encoder and the demo visibly stops filtering.
class _MockOpenapi extends Mock implements api.Openapi {}

class _MockReceiptApi extends Mock implements api.ReceiptApi {}

/// One seeded receipt. [monthsAgo] is relative to today, so the demo's months
/// are always the ones a stepper starting at "All time" can reach.
typedef _Row = ({String name, int monthsAgo, int day, String amount});

const _seed = <_Row>[
  (name: "Whole Foods Market", monthsAgo: 0, day: 8, amount: "84.21"),
  (name: "Shell", monthsAgo: 0, day: 3, amount: "46.10"),
  (name: "Costco", monthsAgo: 1, day: 22, amount: "212.04"),
  (name: "Trader Joe's", monthsAgo: 1, day: 9, amount: "38.77"),
  (name: "Delta Air Lines", monthsAgo: 2, day: 17, amount: "418.60"),
];

void main() {
  // One panel, so the window is the phone box plus its caption strip.
  const surface = Size(demoPhoneWidth, demoPanelHeight);
  const captionColor = Color(0xFF0086D4);

  late List<({img.Image frame, int centis})> frames;

  setUp(() {
    registerCustomCurrencyForTests();
    frames = [];
    registerFallbackValue((api.ReceiptPagedRequestCommandBuilder()
          ..page = 1
          ..pageSize = 10)
        .build());
  });

  tearDown(() => OpenApiClient.client = _MockOpenapi());

  testWidgets("records the quick date filter demo", (tester) async {
    await loadDemoFontsOrFail(tester);
    tester.view.physicalSize = surface;
    tester.view.devicePixelRatio = 1.0;
    addTearDown(tester.view.reset);

    _installFilterAwareApi();

    final harness = buildReceiptFilterHarness();
    await tester.pumpWidget(buildDemoSurface(
      label: "Quick date filter",
      labelColor: captionColor,
      app: _app(harness),
    ));
    await tester.pumpAndSettle();

    /// Captures the settled screen and holds it for [centis] hundredths.
    Future<void> hold(int centis) async {
      await tester.pumpAndSettle();
      frames.add((frame: await grabFrame(tester), centis: centis));
    }

    /// Taps [finder], lets the frame settle, then holds the result.
    Future<void> tapAndHold(Finder finder, int centis) async {
      await tester.tap(finder);
      await hold(centis);
    }

    // 1. All time: every receipt, across three months.
    await hold(220);

    // 2. Step back a month, twice. The rows narrow because the BETWEEN the
    //    client encoded selected them.
    await tapAndHold(find.byKey(const ValueKey("receipt-month-prev")), 190);
    await tapAndHold(find.byKey(const ValueKey("receipt-month-prev")), 190);

    // 3. Open the month sheet and jump straight back with a shortcut.
    await tapAndHold(find.byKey(const ValueKey("receipt-month-label")), 200);
    await tapAndHold(
        find.byKey(const ValueKey("month-picker-this-month")), 200);

    // 4. Re-point the control at another date column. Non-destructive, so the
    //    list does not move and the Receipt Date condition stays applied — the
    //    stepper just has nothing to describe on the column it now reads, so it
    //    falls back to "All time" and drops its clear button.
    await tapAndHold(
        find.byKey(const ValueKey("receipt-quick-date-field")), 170);
    await tapAndHold(
        find.byKey(const ValueKey("receipt-quick-date-field-resolvedDate")),
        230);

    // 5. Switch back, and the month is still there — which is the proof that
    //    re-pointing the control discarded nothing.
    await tapAndHold(
        find.byKey(const ValueKey("receipt-quick-date-field")), 170);
    await tapAndHold(
        find.byKey(const ValueKey("receipt-quick-date-field-date")), 200);

    // 6. Clear back to all time.
    await tapAndHold(find.byKey(const ValueKey("receipt-month-clear")), 240);

    writeGif(frames: frames, path: "tool/quick-date-filter.gif");
  }, timeout: const Timeout(Duration(minutes: 10)));
}

/// The app under the demo surface: the real receipts list, on a real route, in
/// the real theme and provider tree.
Widget _app(ReceiptFilterHarness harness) {
  return MultiProvider(
    providers: [
      ChangeNotifierProvider<ReceiptListModel>.value(
          value: harness.receiptListModel),
      ChangeNotifierProvider<GroupModel>.value(value: harness.groupModel),
      ChangeNotifierProvider<CategoryModel>.value(value: harness.categoryModel),
      ChangeNotifierProvider<TagModel>.value(value: harness.tagModel),
      ChangeNotifierProvider<UserModel>.value(value: harness.userModel),
      ChangeNotifierProvider<ContextModel>(create: (_) => ContextModel()),
      ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
      ChangeNotifierProvider<PermissionsModel>(
          create: (_) => PermissionsModel()),
      ChangeNotifierProvider<SystemSettingsModel>(
          create: (_) => SystemSettingsModel()),
    ],
    child: MaterialApp.router(
      // MaterialApp draws its red ribbon across the top-right of every captured
      // frame otherwise, and at this scale it reads as a render-overflow stripe
      // rather than a banner.
      debugShowCheckedModeBanner: false,
      theme: buildAppTheme(),
      routerConfig: GoRouter(
        initialLocation: "/groups/${ReceiptFilterHarness.householdId}/receipts",
        routes: [
          GoRoute(
            path: "/groups/:groupId/receipts",
            builder: (context, state) => const Scaffold(
              body: Padding(
                padding: EdgeInsets.symmetric(horizontal: 12),
                child: GroupReceiptsList(),
              ),
            ),
          ),
        ],
      ),
    ),
  );
}

/// Stubs `getReceiptsForGroup` with a fake that **honours the filter**.
///
/// It decodes the `date` / `resolvedDate` / `createdAt` condition out of the
/// command the client built and returns only the seeded rows inside it — the
/// same `BETWEEN` bounds the Go query builder would compare against. That is
/// what makes the GIF evidence rather than choreography.
void _installFilterAwareApi() {
  final mockClient = _MockOpenapi();
  final mockReceiptApi = _MockReceiptApi();

  when(() => mockClient.getReceiptApi()).thenReturn(mockReceiptApi);
  when(() => mockReceiptApi.getReceiptsForGroup(
        groupId: any(named: "groupId"),
        receiptPagedRequestCommand: any(named: "receiptPagedRequestCommand"),
      )).thenAnswer((invocation) async {
    final command = invocation.namedArguments[#receiptPagedRequestCommand]
        as api.ReceiptPagedRequestCommand;
    final rows = _matching(command.filter);

    return Response(
      requestOptions: RequestOptions(path: "/"),
      data: (api.PagedDataBuilder()
            ..data = ListBuilder<api.PagedDataDataInner>(rows)
            ..totalCount = rows.length)
          .build(),
    );
  });

  OpenApiClient.client = mockClient;
}

/// The seeded rows the encoded [filter] selects.
List<api.PagedDataDataInner> _matching(api.ReceiptPagedRequestFilter? filter) {
  final serialized = filter == null
      ? const <String, dynamic>{}
      : Map<String, dynamic>.from(api.standardSerializers
          .serializeWith(api.ReceiptPagedRequestFilter.serializer, filter)
          as Map);

  // Only the receipt-date column narrows the seed set here: nothing in it is
  // resolved, so a resolvedDate condition legitimately matches nothing, which
  // is exactly what the demo's field switch is meant to show.
  DateTimeRange? range;
  for (final key in const ["date", "resolvedDate", "createdAt"]) {
    final field = serialized[key];
    if (field == null) {
      continue;
    }
    final value = Map<String, dynamic>.from(field as Map)["value"];
    if (value is List && value.length == 2) {
      range = DateTimeRange(
        start: DateTime.parse(value[0] as String),
        end: DateTime.parse(value[1] as String),
      );
    }
  }

  final matched = <api.PagedDataDataInner>[];
  for (var i = 0; i < _seed.length; i++) {
    final row = _seed[i];
    final date = _dateOf(row);
    if (range != null &&
        (date.isBefore(range.start) || date.isAfter(range.end))) {
      continue;
    }
    matched.add((api.PagedDataDataInnerBuilder()
          ..anyOf = AnyOf1<api.Receipt>(value: _receipt(i, row, date)))
        .build());
  }
  return matched;
}

DateTime _dateOf(_Row row) {
  final month = shiftMonth(monthOfDate(DateTime.now()), -row.monthsAgo);
  return DateTime(month.year, month.month, row.day, 9);
}

api.Receipt _receipt(int index, _Row row, DateTime date) {
  final stamp = date.toIso8601String();
  return (api.ReceiptBuilder()
        ..id = index + 1
        ..name = row.name
        ..amount = row.amount
        ..date = stamp
        ..createdAt = stamp
        ..updatedAt = stamp
        ..groupId = ReceiptFilterHarness.householdId
        ..paidByUserId = 10
        ..status = api.ReceiptStatus.OPEN)
      .build();
}
