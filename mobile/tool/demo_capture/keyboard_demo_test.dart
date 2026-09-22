import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:go_router/go_router.dart';
import 'package:image/image.dart' as img;
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/shared/functions/multi_select_bottom_sheet.dart';
import 'package:receipt_wrangler_mobile/shared/functions/quick_scan.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/screen_wrapper.dart';

import '../../test/helpers/permission_test_helpers.dart';
import '../../test/helpers/receipt_entry_test_helpers.dart';
import '../../test/helpers/receipt_form_test_helpers.dart';
import 'capture.dart';

/// Records the before/after GIFs for the keyboard-inset fix.
///
/// **Not** part of the test suite. `flutter test` with no arguments scans only
/// `test/`, so CI never runs this and never rewrites the committed GIFs; an
/// explicit path is used verbatim, which is how it is invoked:
///
/// ```bash
/// cd mobile && ./tool/record_keyboard_demo.sh
/// ```
///
/// Each demo records the **same real screen twice**, differing only by
/// `debugDisableKeyboardLift`, and stitches the two ramps side by side. So the
/// "before" panel is the production tree with the fix switched off rather than
/// a hand-copied imitation that could drift, and the "after" panel drives the
/// real widgets — if the fix regresses, re-running this makes both panels
/// identical.
///
/// Note the two phases each pump their **own** freshly built tree. Mounting a
/// tree and then re-pumping it wrapped in the demo surface instead re-parents
/// the router subtree, which deadlocks — hence the `buildApp` factory rather
/// than a pre-built widget.
void main() {
  // One panel per capture, so the window must match the phone box:
  // `ScreenWrapper`'s body sizes its Container to `MediaQuery.size.width`.
  const surface = Size(demoPhoneWidth, demoPanelHeight);

  const beforeColor = Color(0xFFB4232B);
  const afterColor = Color(0xFF15803D);

  setUp(() {
    // The caret blinks forever otherwise, so frames would differ run to run.
    EditableText.debugDeterministicCursor = true;
  });

  tearDown(() {
    EditableText.debugDeterministicCursor = false;
    // ignore: invalid_use_of_visible_for_testing_member
    debugDisableKeyboardLift = false;
  });

  /// Records one panel: mount fresh with the lift on or off, then run the ramp.
  Future<List<img.Image>> recordPanel(
    WidgetTester tester, {
    required bool legacy,
    required Widget Function() buildApp,
    Future<void> Function(WidgetTester)? open,
  }) async {
    // ignore: invalid_use_of_visible_for_testing_member
    debugDisableKeyboardLift = legacy;

    // Tear the previous panel down explicitly. `pumpWidget` *updates* a
    // structurally identical tree rather than replacing it, so without this the
    // first phase's Navigator survives into the second and `open` pushes a
    // second sheet on top of the first one.
    await tester.pumpWidget(const SizedBox.shrink());
    await tester.pump();

    tester.view.physicalSize = surface;
    tester.view.devicePixelRatio = 1.0;
    tester.view.viewInsets = FakeViewPadding.zero;

    await tester.pumpWidget(buildDemoSurface(
      app: buildApp(),
      label: legacy ? 'BEFORE — keyboard covers it' : 'AFTER — rides above it',
      labelColor: legacy ? beforeColor : afterColor,
    ));
    await tester.pump();
    if (open != null) {
      await open(tester);
    }
    await tester.pumpAndSettle();

    return recordInsetRamp(tester);
  }

  /// Records both panels and writes the stitched GIF.
  Future<void> record(
    WidgetTester tester, {
    required String outPath,
    required Widget Function() buildApp,
    Future<void> Function(WidgetTester)? open,
  }) async {
    addTearDown(tester.view.reset);
    await loadDemoFontsOrFail(tester);

    final before =
        await recordPanel(tester, legacy: true, buildApp: buildApp, open: open);
    final after =
        await recordPanel(tester, legacy: false, buildApp: buildApp, open: open);

    writeSideBySideGif(before: before, after: after, path: outPath);
  }

  testWidgets('quick scan sheet', (tester) async {
    late BuildContext captured;

    Widget buildApp() {
      final router = GoRouter(
        initialLocation: '/groups/5/receipts',
        routes: [
          GoRoute(
            path: '/groups/:groupId/receipts',
            builder: (context, state) {
              captured = context;
              return const Scaffold(body: SizedBox.shrink());
            },
          ),
          GoRoute(
            path: '/receipts/add',
            builder: (context, state) => const Scaffold(body: SizedBox.shrink()),
          ),
        ],
      );

      return pumpReceiptEntryApp(
        router: router,
        aiEnabled: true,
        permissions: seededPermissions(group: {
          5: [
            api.Permission.groupPeriodReceiptsPeriodQuickScan,
            api.Permission.groupPeriodReceiptsPeriodCreate,
          ]
        }),
        groups: [buildGroup(id: 5, name: 'Household')],
        showDebugBanner: false,
      );
    }

    await record(
      tester,
      outPath: 'tool/quick-scan-keyboard.gif',
      buildApp: buildApp,
      open: (t) async {
        showQuickScanBottomSheet(captured,
            initialImages: [buildQuickScanImage(groupId: 5)]);
        await t.pumpAndSettle();
      },
    );
  }, timeout: const Timeout(Duration(minutes: 5)));

  testWidgets('category / tag / users picker', (tester) async {
    late BuildContext sheetContext;

    Widget buildApp() => MultiProvider(
          providers: [
            // The sheet's TopAppBar reads AuthModel; BottomSubmitButton reads
            // LoadingModel.
            ChangeNotifierProvider<AuthModel>(create: (_) => AuthModel()),
            ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
          ],
          child: MaterialApp(
            debugShowCheckedModeBanner: false,
            theme: ThemeData(fontFamily: 'Raleway', useMaterial3: true),
            home: Builder(builder: (context) {
              sheetContext = context;
              return const Scaffold(body: SizedBox.shrink());
            }),
          ),
        );

    await record(
      tester,
      outPath: 'tool/multi-select-keyboard.gif',
      buildApp: buildApp,
      open: (t) async {
        showMultiselectBottomSheet(
          sheetContext,
          'Categories',
          'Select',
          List.generate(40, (i) => 'Category ${i + 1}'),
          <String>[],
          (option) => option as String,
        );
        await t.pumpAndSettle();
        // Focus the Filter field, so the keyboard in the recording is the one
        // this sheet actually raises rather than an abstract rectangle.
        await t.tap(find.byWidgetPredicate((w) =>
            w is TextField && w.decoration?.labelText == 'Filter'));
        await t.pump();
      },
    );
  }, timeout: const Timeout(Duration(minutes: 5)));

  testWidgets('search shell', (tester) async {
    // /search is the one screen that fills BOTH slots (`lib/main.dart:234`,
    // `:238`): the nav bar in bottomNavigationBar, WranglerSearchBar in
    // bottomSheet. Lifting the bar therefore restacks the search field above
    // it — the deliberate consequence of covering every bottom bar rather than
    // only the two sheets, which is what this panel exists to show.
    Widget buildApp() => MaterialApp(
          debugShowCheckedModeBanner: false,
          theme: ThemeData(fontFamily: 'Raleway', useMaterial3: true),
          home: ScreenWrapper(
            appBarWidget: AppBar(title: const Text('Search')),
            bottomNavigationBarWidget: NavigationBar(
              selectedIndex: 2,
              destinations: const [
                NavigationDestination(icon: Icon(Icons.groups), label: 'Groups'),
                NavigationDestination(
                    icon: Icon(Icons.receipt_long), label: 'Receipts'),
                NavigationDestination(icon: Icon(Icons.search), label: 'Search'),
              ],
            ),
            bottomSheetWidget: const Padding(
              padding: EdgeInsets.all(8),
              child: TextField(
                decoration: InputDecoration(
                  labelText: 'Search',
                  border: OutlineInputBorder(),
                ),
              ),
            ),
            child: ListView(
              children: [
                for (var i = 1; i <= 12; i++)
                  ListTile(
                    leading: const Icon(Icons.receipt),
                    title: Text('Receipt $i'),
                    subtitle: const Text('Household'),
                  ),
              ],
            ),
          ),
        );

    await record(
      tester,
      outPath: 'tool/search-nav-keyboard.gif',
      buildApp: buildApp,
      open: (t) async {
        await t.tap(find.byType(TextField));
        await t.pump();
      },
    );
  }, timeout: const Timeout(Duration(minutes: 5)));
}
