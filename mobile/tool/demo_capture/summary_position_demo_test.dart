import 'package:built_collection/built_collection.dart';
import 'package:one_of/any_of.dart';
import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:image/image.dart' as img;
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
import 'package:receipt_wrangler_mobile/shared/widgets/screen_wrapper.dart';
import 'package:receipt_wrangler_mobile/theme/app_theme.dart';

import '../../test/helpers/receipt_form_test_helpers.dart';
import '../../test/helpers/receipt_summary_test_helpers.dart';
import '../../test/helpers/widget_test_helpers.dart';
import 'capture.dart';

/// Records the receipt-summary placement GIF.
///
/// **Not** part of the test suite. `flutter test` with no arguments scans only `test/`,
/// so CI never runs this and never rewrites the committed GIF; an explicit path is used
/// verbatim, which is how it is invoked:
///
/// ```bash
/// cd mobile && ./tool/record_summary_position_demo.sh
/// ```
///
/// **This is the first demo here that needed NO production seam**, and that is a property
/// of the feature rather than luck: the position arrives as *data on the summary
/// response*, so the two panels are the identical real tree fed two responses differing
/// in one field. Better than a debug flag — there is no flag to leave switched on, no
/// source to restore under a `trap`, and if placement regresses the two panels simply
/// become identical (the self-invalidating property the keyboard demo has).
///
/// Only `OpenApiClient.client` is mocked. Everything else is the shipping tree: the real
/// theme, the real `ScreenWrapper` with the receipts route's zeroed body padding, the
/// real providers, and the real `GroupReceiptsList` with its real `PagedDataList`.
void main() {
  const surface = Size(demoPhoneWidth, demoPanelHeight);

  const topColor = Color(0xFF15803D);
  const bottomColor = Color(0xFF1D4ED8);

  const groupId = 7;

  /// The list is scrolled rather than the keyboard raised, so the ramp is its own.
  /// A beat at rest, a dozen scroll steps, then a long hold on the outcome -- the frame
  /// a reviewer actually reads. Buy hold time with DURATIONS, never duplicate frames:
  /// GifEncoder writes every frame in full, so a held frame costs a moving one.
  const scrollSteps = 7;
  List<int> durations() => const [
        170,
        ...[8, 8, 8, 8, 8, 8, 8],
        300,
      ];

  api.Receipt receipt(int id, String name, String amount, api.ReceiptStatus status) =>
      (api.ReceiptBuilder()
            ..id = id
            ..createdAt = '2026-09-0${(id % 9) + 1}T12:00:00Z'
            ..name = name
            ..amount = amount
            ..date = '2026-09-0${(id % 9) + 1}T12:00:00Z'
            ..paidByUserId = 10
            ..groupId = groupId
            ..status = status)
          .build();

  final rows = <api.Receipt>[
    receipt(1, 'Costco run', '184.32', api.ReceiptStatus.OPEN),
    receipt(2, 'Office chairs', '612.50', api.ReceiptStatus.RESOLVED),
    receipt(3, 'Team lunch', '96.40', api.ReceiptStatus.OPEN),
    receipt(4, 'Printer toner', '78.15', api.ReceiptStatus.RESOLVED),
    receipt(5, 'Client dinner', '243.90', api.ReceiptStatus.OPEN),
    receipt(6, 'Whiteboard markers', '31.25', api.ReceiptStatus.DRAFT),
    receipt(7, 'Monitor arm', '129.99', api.ReceiptStatus.OPEN),
    receipt(8, 'Coffee beans', '42.80', api.ReceiptStatus.RESOLVED),
  ];

  api.ReceiptSummary summaryAt(api.ReceiptSummaryPosition position) =>
      buildReceiptSummary(
        position: position,
        overall: buildSummaryRow(
          receiptCount: 8,
          total: '1419.31',
          customFieldTotals: [
            buildSummaryCustomFieldTotal(
                customFieldId: 1, name: 'HST', total: '184.51'),
          ],
        ),
        statuses: [
          buildSummaryRow(
            status: api.ReceiptStatus.OPEN,
            receiptCount: 4,
            total: '654.61',
            customFieldTotals: [
              buildSummaryCustomFieldTotal(
                  customFieldId: 1, name: 'HST', total: '85.10'),
            ],
          ),
          buildSummaryRow(
            status: api.ReceiptStatus.RESOLVED,
            receiptCount: 3,
            total: '733.45',
            customFieldTotals: [
              buildSummaryCustomFieldTotal(
                  customFieldId: 1, name: 'HST', total: '95.35'),
            ],
          ),
        ],
      );

  Widget buildApp() {
    final groupModel = GroupModel()
      ..setGroups([
        buildGroup(
          id: groupId,
          name: 'Summary Demo',
          receiptSummaryEnabled: true,
          members: [buildGroupMember(userId: 10, groupId: groupId)],
        ),
      ]);

    // ReceiptListItem resolves every row's payer by id and `firstWhere`s with no
    // orElse, so an unseeded UserModel throws "Bad state: No element" per row.
    final userModel = UserModel()
      ..setUsers([buildUserView(id: 10, displayName: 'Noah Hall')]);

    final router = GoRouter(
      initialLocation: '/groups/$groupId/receipts',
      routes: [
        GoRoute(
          path: '/groups/:groupId/receipts',
          // Mirrors lib/main.dart: the receipts route zeroes the body padding so
          // rows sit edge to edge, which is why the bar supplies its own 16.
          builder: (context, state) => const ScreenWrapper(
            bodyPadding: EdgeInsets.all(0),
            child: GroupReceiptsList(),
          ),
        ),
      ],
    );

    return MultiProvider(
      providers: [
        ChangeNotifierProvider<ReceiptListModel>(create: (_) => ReceiptListModel()),
        ChangeNotifierProvider<GroupModel>.value(value: groupModel),
        ChangeNotifierProvider<CategoryModel>(create: (_) => CategoryModel()),
        ChangeNotifierProvider<TagModel>(create: (_) => TagModel()),
        ChangeNotifierProvider<UserModel>.value(value: userModel),
        ChangeNotifierProvider<ContextModel>(create: (_) => ContextModel()),
        ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
        ChangeNotifierProvider<PermissionsModel>(create: (_) => PermissionsModel()),
        ChangeNotifierProvider<SystemSettingsModel>(
            create: (_) => SystemSettingsModel()),
      ],
      child: MaterialApp.router(
        theme: buildAppTheme(),
        // Material's red ribbon reads as a render-overflow stripe at panel scale and
        // sends the reviewer hunting a layout bug that is not there.
        debugShowCheckedModeBanner: false,
        routerConfig: router,
      ),
    );
  }

  /// Records one panel: mount fresh at [position], then scroll the list under it.
  Future<List<img.Image>> recordPanel(
    WidgetTester tester,
    api.ReceiptSummaryPosition position,
  ) async {
    _stubClient(rows, summaryAt(position));

    // Tear the previous panel down explicitly: `pumpWidget` UPDATES a structurally
    // identical tree rather than replacing it, and `capture.dart`'s boundary key is a
    // single file-private GlobalKey, so two surfaces cannot be mounted at once.
    await tester.pumpWidget(const SizedBox.shrink());
    await tester.pump();

    tester.view.physicalSize = surface;
    tester.view.devicePixelRatio = 1.0;
    tester.view.viewInsets = FakeViewPadding.zero;

    final atTop = position == api.ReceiptSummaryPosition.TOP;
    await tester.pumpWidget(buildDemoSurface(
      app: buildApp(),
      label: atTop ? 'TOP — totals above the list' : 'BOTTOM — totals under the list',
      labelColor: atTop ? topColor : bottomColor,
    ));
    await tester.pumpAndSettle();

    final frames = <img.Image>[
      // A beat at rest before anything moves.
      await grabFrame(tester),
    ];

    // An explicit gesture with fixed pumps, never `drag` + `pumpAndSettle`: ballistic
    // scroll physics produce a frame count that is not reproducible run to run, and
    // pumpAndSettle is independently unsafe here (PagedDataList's first-load spinner
    // animates, and TopAppBar's progress bar is permanently mounted with a ten-minute
    // settle timeout -- that regression hangs rather than fails).
    final gesture =
        await tester.startGesture(tester.getCenter(find.byType(PagedDataList)));
    for (var i = 0; i < scrollSteps; i++) {
      await gesture.moveBy(const Offset(0, -38));
      await tester.pump(const Duration(milliseconds: 40));
      frames.add(await grabFrame(tester));
    }
    await gesture.up();
    await tester.pump();

    // A hold: the list has moved, the bar has not, which is the whole point of the
    // recording. There used to be a sideways swipe of the figure grid after this, to
    // show the clipped currency column was scrollable rather than broken; the rows
    // wrap now, so there is nothing off-screen to reveal.
    frames.add(await grabFrame(tester));
    return frames;
  }

  testWidgets('summary position', (tester) async {
    addTearDown(tester.view.reset);
    addTearDown(() => OpenApiClient.client = _MockOpenapi());
    registerCustomCurrencyForTests();
    registerFallbackValue((api.ReceiptPagedRequestCommandBuilder()
          ..page = 1
          ..pageSize = 10)
        .build());
    registerFallbackValue(api.ReceiptSummaryCommandBuilder().build());
    await tester.runAsync(loadDemoFonts);

    final top = await recordPanel(tester, api.ReceiptSummaryPosition.TOP);
    final bottom = await recordPanel(tester, api.ReceiptSummaryPosition.BOTTOM);

    writeSideBySideGif(
      before: top,
      after: bottom,
      path: 'tool/receipt-summary-position.gif',
      durationsCentis: durations(),
      // The status chips are a gradient; at 128 the octree banding reads as black
      // stripes through them.
      numColors: 192,
    );
  });
}

class _MockOpenapi extends Mock implements api.Openapi {}

class _MockReceiptApi extends Mock implements api.ReceiptApi {}

/// Serves one page of [rows] and a fixed [summary]. The only mock in the whole demo.
void _stubClient(List<api.Receipt> rows, api.ReceiptSummary summary) {
  final client = _MockOpenapi();
  final receiptApi = _MockReceiptApi();
  when(() => client.getReceiptApi()).thenReturn(receiptApi);

  when(() => receiptApi.getReceiptsForGroup(
        groupId: any(named: 'groupId'),
        receiptPagedRequestCommand: any(named: 'receiptPagedRequestCommand'),
      )).thenAnswer((_) async => Response(
        requestOptions: RequestOptions(path: '/'),
        data: (api.PagedDataBuilder()
              ..data = ListBuilder<api.PagedDataDataInner>(rows.map((r) =>
                  (api.PagedDataDataInnerBuilder()
                        ..anyOf = AnyOf1<api.Receipt>(value: r))
                      .build()))
              ..totalCount = rows.length)
            .build(),
      ));

  when(() => receiptApi.getReceiptSummaryForGroup(
        groupId: any(named: 'groupId'),
        receiptSummaryCommand: any(named: 'receiptSummaryCommand'),
      )).thenAnswer((_) async => Response(
        requestOptions: RequestOptions(path: '/'),
        data: summary,
      ));

  OpenApiClient.client = client;
}
