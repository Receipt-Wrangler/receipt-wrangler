import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:image/image.dart' as img;
import 'package:mocktail/mocktail.dart';
import 'package:money2/money2.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/constants/currency.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/category_model.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/models/custom_field_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/models/system_settings_model.dart';
import 'package:receipt_wrangler_mobile/models/tag_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/receipts/nav/receipt_app_bar.dart';
import 'package:receipt_wrangler_mobile/receipts/nav/receipt_app_bar_action_builder.dart';
import 'package:receipt_wrangler_mobile/receipts/screens/receipt_history_screen.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_form.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_history_entry_card.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/receipt_edit_popup_menu.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/screen_wrapper.dart';
import 'package:receipt_wrangler_mobile/theme/app_theme.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';

import '../../test/helpers/permission_test_helpers.dart';
import '../../test/helpers/receipt_form_test_helpers.dart';
import '../../test/helpers/receipt_history_test_helpers.dart';
import 'capture.dart';

/// Records the demo GIF for the receipt history screen: the receipt's view
/// screen, its ⋮ menu, then History.
///
/// **Not** part of the test suite. `flutter test` with no arguments scans only
/// `test/`, so CI never runs this and never rewrites the committed GIF:
///
/// ```bash
/// cd mobile && ./tool/record_receipt_history_demo.sh
/// ```
///
/// Everything on screen is the shipping code: the receipt form in view mode,
/// the app bar the form screen builds (`ReceiptAppBar` +
/// `ReceiptAppBarActionBuilder`), the push to `ReceiptHistoryScreen`, and the
/// cards it renders. Two things stand in for the server: the history endpoint
/// is mocked, returning entries built through the real generated deserializer,
/// and the receipt is seeded into `ReceiptModel` rather than fetched (the host
/// mounts the form screen's body, not its fetch).
class _MockOpenapi extends Mock implements api.Openapi {}

class _MockReceiptApi extends Mock implements api.ReceiptApi {}

const _receiptId = 5;
const _groupId = 1;
const _maya = 2;
const _sam = 3;

/// Twice the logical size: the GIF is read, so text has to stay sharp.
const _pixelRatio = 2.0;

void main() {
  const surface = Size(demoPhoneWidth, demoPanelHeight);
  const captionColor = Color(0xFF0086D4);

  late List<({img.Image frame, int centis})> frames;

  setUp(() {
    _registerDollars();
    frames = [];
  });

  tearDown(() => OpenApiClient.client = _MockOpenapi());

  testWidgets('records the receipt history demo', (tester) async {
    await loadDemoFontsOrFail(tester);
    tester.view.physicalSize = surface;
    tester.view.devicePixelRatio = 1.0;
    addTearDown(tester.view.reset);

    _installHistoryApi();

    await tester.pumpWidget(buildDemoSurface(
      label: 'Receipt history',
      labelColor: captionColor,
      app: _app(),
    ));
    await tester.pumpAndSettle();

    /// Captures the settled screen and holds it for [centis] hundredths.
    Future<void> capture(int centis) async {
      frames.add((
        frame: await grabFrame(tester, pixelRatio: _pixelRatio),
        centis: centis,
      ));
    }

    Future<void> hold(int centis) async {
      await tester.pumpAndSettle();
      await capture(centis);
    }

    Future<void> tapAndHold(Finder finder, int centis) async {
      await tester.tap(finder);
      await hold(centis);
    }

    /// Scrolls the history list by [dy] in a few captured steps, so the GIF
    /// shows the list moving rather than jumping.
    Future<void> scrollBy(double dy, {int steps = 5}) async {
      final list = find
          .descendant(
            of: find.byType(ReceiptHistoryScreen),
            matching: find.byType(Scrollable),
          )
          .first;
      final position = tester.state<ScrollableState>(list).position;
      final start = position.pixels;
      final end = (start + dy)
          .clamp(position.minScrollExtent, position.maxScrollExtent);
      for (var i = 1; i <= steps; i++) {
        // Ease out, so the list decelerates into place like a fling.
        final t = 1 - (1 - i / steps) * (1 - i / steps);
        position.jumpTo(start + (end - start) * t);
        await tester.pump();
        await capture(6);
      }
    }

    final showMore = find
        .descendant(
          of: find.byKey(const ValueKey('receipt-history-entry-3')),
          matching: find.byKey(ReceiptHistoryEntryCard.showMoreKey),
        )
        .first;

    // 1. The receipt, in view mode.
    await hold(200);

    // 2. The ⋮ menu: Edit, then History.
    await tapAndHold(
      find.descendant(
        of: find.byType(ReceiptEditPopupMenu),
        matching: find.byType(PopupMenuButton<dynamic>),
      ),
      170,
    );

    // 3. History, pushed over the form.
    await tapAndHold(
        find.byKey(ReceiptAppBarActionBuilder.historyMenuItemKey), 300);

    // 4. Down past the comment and the failed save to the edit, which
    //    collapsed after three changes; expand it.
    await scrollBy(tester.getRect(showMore).bottom - (surface.height - 24));
    await hold(150);
    await tapAndHold(showMore, 260);

    // 5. On to the email upload that created the receipt.
    await scrollBy(1000);
    await hold(340);

    // Neural, not octree: the cards' 1px rounded borders speckle dark at the
    // corners under octree.
    writeGif(
      frames: frames,
      path: 'tool/receipt-history.gif',
      quantizerType: img.QuantizerType.neural,
    );
  }, timeout: const Timeout(Duration(minutes: 10)));
}

/// The receipt form screen's shape -- `ScreenWrapper`, the real receipt app
/// bar with the real action builder, and the form -- minus its fetch.
class _ReceiptViewHost extends StatelessWidget {
  const _ReceiptViewHost({required this.state});

  final GoRouterState state;

  @override
  Widget build(BuildContext context) {
    final receiptModel = Provider.of<ReceiptModel>(context, listen: false);
    final actionBuilder = ReceiptAppBarActionBuilder(context, receiptModel);
    Provider.of<ContextModel>(context, listen: false).setShellContext(context);

    return ScreenWrapper(
      appBarWidget: ReceiptAppBar(actions: actionBuilder.buildAppBarMenu(state)),
      child: const SingleChildScrollView(child: ReceiptForm()),
    );
  }
}

api.UserView _user(int id, String name, String color) =>
    buildUserView(id: id, displayName: name)
        .rebuild((b) => b..defaultAvatarColor = color);

Widget _app() {
  final receipt = getDefaultReceipt().rebuild((b) => b
    ..id = _receiptId
    ..name = 'Costco Wholesale'
    ..amount = '198.40'
    ..date = _daysAgo(3, 0, 0).toIso8601String()
    ..groupId = _groupId
    ..paidByUserId = _maya
    ..status = api.ReceiptStatus.RESOLVED);

  final permissions = seededPermissions(
    app: [api.Permission.appPeriodReceiptHistoryPeriodRead],
    group: {
      _groupId: [
        api.Permission.groupPeriodReceiptsPeriodRead,
        api.Permission.groupPeriodReceiptsPeriodUpdate,
      ],
    },
  );

  return MultiProvider(
    providers: [
      ChangeNotifierProvider<ReceiptModel>(
          create: (_) => ReceiptModel()..setReceipt(receipt, false)),
      ChangeNotifierProvider<GroupModel>(
        create: (_) => GroupModel()
          ..setGroups([
            buildGroup(id: _groupId, name: 'Household', members: [
              buildGroupMember(userId: _maya, groupId: _groupId),
              buildGroupMember(userId: _sam, groupId: _groupId),
            ]),
          ]),
      ),
      ChangeNotifierProvider<UserModel>(
        create: (_) => UserModel()
          ..setUsers([
            _user(_maya, 'Maya Chen', '#7C3AED'),
            _user(_sam, 'Sam Rivera', '#0F9D58'),
          ]),
      ),
      ChangeNotifierProvider<CustomFieldModel>(
          create: (_) => CustomFieldModel()),
      ChangeNotifierProvider<CategoryModel>(create: (_) => CategoryModel()),
      ChangeNotifierProvider<TagModel>(create: (_) => TagModel()),
      ChangeNotifierProvider<ContextModel>(create: (_) => ContextModel()),
      ChangeNotifierProvider<SystemSettingsModel>(
          create: (_) => SystemSettingsModel()),
      ChangeNotifierProvider<PermissionsModel>.value(value: permissions),
      ChangeNotifierProvider<AuthModel>(create: (_) => AuthModel()),
      ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
    ],
    child: MaterialApp.router(
      debugShowCheckedModeBanner: false,
      theme: buildAppTheme(),
      routerConfig: GoRouter(
        initialLocation: '/receipts/$_receiptId/view',
        routes: [
          GoRoute(
            path: '/receipts/:receiptId/view',
            builder: (context, state) => _ReceiptViewHost(state: state),
          ),
        ],
      ),
    ),
  );
}

// Local wall-clock times sent as UTC, so each entry lands on the day and hour
// written here in whatever zone the recording runs.
String _at(DateTime local) => local.toUtc().toIso8601String();

DateTime _daysAgo(int days, int hour, int minute) {
  final now = DateTime.now();
  return DateTime(now.year, now.month, now.day - days, hour, minute);
}

Map<String, Object?> _change(
  String field, {
  Map<String, Object?>? from,
  Map<String, Object?>? to,
  int? userId,
}) =>
    {
      'field': field,
      'kind': from != null && to != null
          ? 'CHANGED'
          : (to != null ? 'ADDED' : 'REMOVED'),
      if (userId != null) 'userId': userId,
      if (from != null) 'from': from,
      if (to != null) 'to': to,
    };

Map<String, Object?> _value(String type, String value) =>
    {'type': type, 'value': value};

String _day(DateTime date) =>
    '${date.year.toString().padLeft(4, '0')}-'
    '${date.month.toString().padLeft(2, '0')}-'
    '${date.day.toString().padLeft(2, '0')}';

/// The receipt's life: emailed in, corrected, a save the role refused, a
/// comment, then resolved.
api.ReceiptHistory _history() {
  final resolvedAt = _daysAgo(0, 10, 42);
  return historyFromJson(historyJson([
    entryJson(
      systemTaskId: 6,
      action: 'RESOLVED',
      recordedAt: _at(resolvedAt),
      actor: {'type': 'USER', 'userId': _maya},
      changes: [
        _change('STATUS',
            from: _value('STATUS', 'OPEN'), to: _value('STATUS', 'RESOLVED')),
        _change('RESOLVED_DATE', to: _value('INSTANT', _at(resolvedAt))),
      ],
    ),
    entryJson(
      systemTaskId: 5,
      action: 'IMAGES_COMMENTS_UPDATED',
      recordedAt: _at(_daysAgo(1, 18, 5)),
      actor: {'type': 'USER', 'userId': _sam},
      changes: [
        _change('COMMENT',
            userId: _sam,
            to: _value('TEXT', 'Paid my half of the membership separately')),
      ],
    ),
    entryJson(
      systemTaskId: 4,
      action: 'FAILED',
      recordedAt: _at(_daysAgo(1, 17, 51)),
      actor: {'type': 'USER', 'userId': _maya},
      hasSnapshots: false,
      errorMessage: 'Your role requires at least one image.',
    ),
    entryJson(
      systemTaskId: 3,
      action: 'UPDATED',
      recordedAt: _at(_daysAgo(1, 17, 48)),
      actor: {'type': 'USER', 'userId': _maya},
      changes: [
        _change('NAME',
            from: _value('TEXT', 'COSTCO WHSE #1042'),
            to: _value('TEXT', 'Costco Wholesale')),
        _change('AMOUNT',
            from: _value('CURRENCY', '212.04'),
            to: _value('CURRENCY', '198.40')),
        _change('DATE',
            from: _value('DATE', _day(_daysAgo(2, 0, 0))),
            to: _value('DATE', _day(_daysAgo(3, 0, 0)))),
        _change('CATEGORIES', to: _value('TEXT', 'Groceries')),
        _change('PAID_BY',
            from: {'type': 'USER', 'userId': _sam},
            to: {'type': 'USER', 'userId': _maya}),
      ],
    ),
    entryJson(
      systemTaskId: 1,
      action: 'CREATED',
      recordedAt: _at(_daysAgo(2, 8, 14)),
      actor: {'type': 'EMAIL'},
      changes: [
        _change('EMAIL_INBOX', to: _value('TEXT', 'test@test.com')),
        _change('EMAIL_SENDER', to: _value('TEXT', 'sent@sent.com')),
        _change('NAME', to: _value('TEXT', 'COSTCO WHSE #1042')),
      ],
    ),
  ], changeCount: 11, saveCount: 4, contributorCount: 3));
}

void _installHistoryApi() {
  final client = _MockOpenapi();
  final receiptApi = _MockReceiptApi();
  when(() => client.getReceiptApi()).thenReturn(receiptApi);
  when(() => receiptApi.getReceiptHistory(receiptId: any(named: 'receiptId')))
      .thenAnswer((_) async => Response<api.ReceiptHistory>(
            requestOptions: RequestOptions(path: ''),
            statusCode: 200,
            data: _history(),
          ));
  OpenApiClient.client = client;
}

/// The app's currency, as `registerCustomCurrency` builds it from System
/// Settings -- here a leading `$`. The shared test registration puts the
/// symbol last, which reads as a bug in a picture.
void _registerDollars() {
  if (Currencies().find(customCurrencyISOCode) != null) {
    return;
  }
  Currencies().register(Currency.create(
    customCurrencyISOCode,
    2,
    name: customCurrencyISOCode,
    symbol: r'$',
    groupSeparator: ',',
    decimalSeparator: '.',
    pattern: 'S###,###.00',
  ));
}
