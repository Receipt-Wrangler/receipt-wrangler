import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_slidable/flutter_slidable.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/slidable_widget.dart';

import 'helpers/api.dart';
import 'helpers/login.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/receipt_test_helpers.dart';

/// Exercises the comments feature end-to-end on an edit-mode receipt:
/// add two comments via the bottom-sheet input, then swipe-delete the
/// first one. Asserts against `Receipt.comments` from the API after each
/// mutation -- the API list is the source of truth. A UI-only assertion
/// would miss the documented swallow in `_submitCommentToApi`'s catch
/// (receipt_comment_screen.dart:166-168), which can keep stale state in
/// the UI even when the POST failed.
void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  // The second comment's send tap used to miss because the send button gets
  // nudged below the visible area once the first comment renders in the list.
  // `_submitComment` now `ensureVisible`s the send button before tapping, so
  // both submits land. See `_submitComment` below.
  testWidgets('admin can add, view, and delete receipt comments',
      (tester) async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));

    await loginAsAdmin(tester);

    final receiptName =
        'e2e-comments-${DateTime.now().millisecondsSinceEpoch}';
    final receiptId = await addManualReceiptViaUI(tester, receiptName);
    scheduleReceiptCleanup(receiptId);

    final jwt = await apiLogin();

    // Move from /view to /edit so the comment bottom sheet (and the
    // swipe-to-delete slidable) become interactive -- both gate on
    // `formState == edit`. Tap the popup menu's "Edit" item.
    //
    // The ReceiptEditPopupMenu is gated on canEditReceipt(), which checks the
    // user's group.receipts.update permission from PermissionsModel; on
    // cold-boot post-navigation, permissions may not be populated yet, so the
    // button isn't mounted immediately. Same pumpUntilFound pattern as
    // receipt_edit_test.dart:50.
    final menuButton = find.byType(PopupMenuButton<dynamic>);
    await pumpUntilFound(tester, menuButton);
    await tester.tap(menuButton);
    // The popup scales in; wait until "Edit" is actually hittable and drain
    // the open animation so the tap's computed center is its settled position
    // (same hardening as receipt_cost_split_test._navigateToEdit).
    await pumpUntilFound(tester, find.text('Edit').hitTestable());
    for (int i = 0; i < 5; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }
    await tester.tap(find.text('Edit').hitTestable());
    // /edit's destination-mounted marker: the bottom save button, which only
    // renders on edit/add paths -- find.text('Name') matches on /view too.
    await pumpUntilFound(tester, find.byType(BottomSubmitButton));

    // The Comments screen is pushed via Navigator (separate from GoRouter)
    // by tapping the "Comments" compact-action button on the edit form
    // (receipts/widgets/receipt_form.dart:391). The button is wrapped in
    // a Tooltip("View Comments") -- byTooltip is a single deterministic
    // match.
    final commentsButton = find.byTooltip('View Comments');
    await pumpUntilFound(tester, commentsButton);
    await tester.tap(commentsButton);
    await pumpUntilFound(tester, find.text('Receipt Comments'));

    const firstComment = 'e2e first comment';
    const secondComment = 'e2e second comment';

    await _submitComment(tester, firstComment);
    await pumpUntilFound(tester, find.text(firstComment));

    await _submitComment(tester, secondComment);
    await pumpUntilFound(tester, find.text(secondComment));

    // Both comments now present on screen -- and the API should agree.
    final afterAdds = await getReceipt(receiptId, jwt: jwt);
    final commentsAfterAdds =
        (afterAdds['comments'] as List).cast<Map<String, dynamic>>();
    expect(commentsAfterAdds.length, 2,
        reason: 'server should have 2 comments after two send taps; '
            "if the UI shows 2 but the API has fewer that's a real bug "
            'in _submitCommentToApi swallowing the POST error');
    expect(commentsAfterAdds.map((c) => c['comment']).toList(),
        containsAll(<String>[firstComment, secondComment]));

    // Swipe-delete the first comment. The slidable wraps a Column of the
    // comment row + a SizedBox spacer (receipt_comments.dart:42-47), so
    // we locate the slidable by walking up from the text we want to remove.
    final firstSlidable = find.ancestor(
      of: find.text(firstComment),
      matching: find.byType(SlidableWidget),
    );
    expect(firstSlidable, findsOneWidget);
    await tester.drag(firstSlidable, const Offset(-300, 0));
    await tester.pumpAndSettle();

    // After the drag, the end-action pane reveals exactly one
    // SlidableAction (the delete button). Tap it.
    final deleteAction = find.byType(SlidableAction);
    await pumpUntilFound(tester, deleteAction);
    await tester.tap(deleteAction);

    // Delete goes through `deleteComment` (await of API DELETE);
    // pump until the deleted text is gone from the visible tree.
    await pumpUntilGone(tester, find.text(firstComment));

    final afterDelete = await getReceipt(receiptId, jwt: jwt);
    final commentsAfterDelete =
        (afterDelete['comments'] as List).cast<Map<String, dynamic>>();
    expect(commentsAfterDelete.length, 1,
        reason: 'server should have 1 comment after swipe-delete');
    expect(commentsAfterDelete.single['comment'], secondComment);
  });
}

/// Types [comment] into the bottom-sheet comment field and taps send.
Future<void> _submitComment(WidgetTester tester, String comment) async {
  final commentField = find.byWidgetPredicate(
    (w) => w is FormBuilderTextField && w.name == 'comment',
  );
  await pumpUntilFound(tester, commentField);
  await tester.enterText(commentField, comment);
  await tester.pumpAndSettle(const Duration(seconds: 1));

  // The submit IconButton is gated on the textBehaviorSubject stream
  // (receipt_comment_screen.dart:82-85). Until the stream emits the new
  // value, IconButton.onPressed stays null and the tap is a no-op.
  // pumpUntil the *enabled* button (onPressed != null) -- not just the
  // widget's existence -- before tapping.
  final sendButton = find.byWidgetPredicate((w) =>
      w is IconButton &&
      w.icon is Icon &&
      (w.icon as Icon).icon == Icons.send &&
      w.onPressed != null);
  await pumpUntilFound(tester, sendButton);
  // After the first comment renders in the list, the send button can sit below
  // the visible area; scroll it back on-screen so the tap lands. ensureVisible
  // jumps the scroll position without relayout -- pump a frame so the tap
  // computes the post-scroll center, not the stale one.
  await tester.ensureVisible(sendButton);
  await tester.pump(const Duration(milliseconds: 100));
  // The "Receipt added/updated successfully" snackbar lingers over the send
  // row (bottom of the 900px surface) for ~4s and absorbs taps -- observed as
  // "Offset(936.0, 872.0) ... would not hit test on the specified widget".
  // hitTestable() filters the obscured button, so this wait resumes exactly
  // when the snackbar departs and the tap can land.
  await pumpUntilFound(tester, sendButton.hitTestable());
  await tester.tap(sendButton.hitTestable());
}
