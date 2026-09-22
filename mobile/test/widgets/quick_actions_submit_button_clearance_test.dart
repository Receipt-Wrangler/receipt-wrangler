import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/quick_actions.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/quick_actions_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/utils/bottom_sheet.dart';

import '../helpers/receipt_form_test_helpers.dart';

/// The quick-actions split sheet's tail must stay clear of its "Split" button.
///
/// Same mechanism as the receipt form (see
/// `receipt_form_submit_button_clearance_test.dart`): the button is passed as
/// `bottomSheetWidget` **without** `bodyFillsSheet`, so it lands in
/// `Scaffold.bottomSheet` and floats over the body rather than reserving space.
///
/// Worse here than on the form: the column had **no** trailing slack at all
/// rather than 20px, and no `kDebugMode` widget padding the debug tree, so this
/// one was buried in every build mode.
///
/// **What this test does and does not prove.** In the default split-evenly mode
/// `buildSplitEvenlyTotal()` returns an empty list while no users are selected,
/// so the Users field is last and merely *touches* the button — nothing is
/// actually hidden. The content that really gets buried is the mode-dependent
/// totals rendered **below** it: in portions mode that is the display carrying
/// "Portions exceed receipt total", the message explaining why Split is
/// refusing to submit. Driving the sheet into that mode needs a user selection
/// made through a nested multi-select route, so instead this asserts the
/// **clearance those totals need** — the last field must sit at least the
/// button's own height above it. That is discriminating (it fails with
/// `submitButtonSpacing` removed) without pretending to exercise the portions
/// path itself.
///
/// The viewport is deliberately short. The sheet's default content is under
/// 200px tall, so at any realistic height `maxScrollExtent` is 0 and nothing
/// can be buried at all — the assertion would pass on the broken tree for the
/// wrong reason. 240px is the smallest height that forces the column to scroll.
void main() {
  const screen = Size(800, 240);

  /// `BottomSubmitButton` is a fixed 50px (`bottom_submit_button.dart`).
  const buttonHeight = 50.0;

  Future<void> pumpSheet(WidgetTester tester) async {
    tester.view.physicalSize = screen;
    tester.view.devicePixelRatio = 1.0;
    addTearDown(tester.view.reset);

    late BuildContext sheetContext;

    await tester.pumpWidget(
      MultiProvider(
        providers: [
          // ReceiptModel already seeds itself with getDefaultReceipt().
          ChangeNotifierProvider<ReceiptModel>(create: (_) => ReceiptModel()),
          ChangeNotifierProvider<UserModel>(create: (_) => UserModel()),
          ChangeNotifierProvider<GroupModel>(
            create: (_) => GroupModel()
              ..setGroups([buildGroup(id: 1, name: 'Household')]),
          ),
          ChangeNotifierProvider<ContextModel>(create: (_) => ContextModel()),
          // The sheet's TopAppBar reads AuthModel; BottomSubmitButton reads
          // LoadingModel.
          ChangeNotifierProvider<AuthModel>(create: (_) => AuthModel()),
          ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
        ],
        child: MaterialApp(
          home: Builder(builder: (context) {
            sheetContext = context;
            return const Scaffold(body: SizedBox.shrink());
          }),
        ),
      ),
    );
    await tester.pump();

    // Opened exactly as `receipt_form.dart` opens it — no `bodyFillsSheet`, so
    // the button goes to the floating slot.
    showFullscreenBottomSheet(
      sheetContext,
      const ReceiptQuickActions(groupId: 1),
      'Quick Actions',
      bottomSheetWidget: const ReceiptQuickActionsSubmitButton(),
    );
    await tester.pumpAndSettle();

    final scrollable = tester.state<ScrollableState>(find.byType(Scrollable).first);
    scrollable.position.jumpTo(scrollable.position.maxScrollExtent);
    await tester.pumpAndSettle();
  }

  testWidgets('the tail of the sheet clears the Split button at max scroll',
      (tester) async {
    await pumpSheet(tester);

    final button = tester.getRect(find.byType(BottomSubmitButton));
    final lastField = tester.getRect(find
        .ancestor(of: find.text('Users'), matching: find.byType(InputDecorator))
        .first);

    // Without submitButtonSpacing this reads exactly `button.top` -- the field
    // is flush against the button, leaving no room for the totals that render
    // below it in the portions and percentage modes.
    expect(lastField.bottom, lessThanOrEqualTo(button.top - buttonHeight),
        reason: 'the column must reserve the floating button\'s height so the '
            'mode-dependent totals below the Users field stay reachable');
  });

  testWidgets('the Split button still sits on the bottom of the sheet',
      (tester) async {
    await pumpSheet(tester);

    // The control: the clearance assertion above must not be satisfiable by a
    // button that drifted up the screen for some unrelated reason.
    expect(tester.getRect(find.byType(BottomSubmitButton)).bottom, screen.height);
  });
}
