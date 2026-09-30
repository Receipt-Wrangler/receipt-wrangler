import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_comments.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/slidable_widget.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';

import '../helpers/permission_test_helpers.dart';

/// Widget-level coverage for the comment swipe-to-delete gate
/// (`receipt_comments.dart` → `SlidableWidget.slideEnabled` =
/// `canCommentDelete`, i.e. `group.comments.delete`, in edit state).
/// Complements the e2e in `integration_test/permission_comments_test.dart`.
void main() {
  const groupId = 7;

  api.Comment comment({int id = 1, String text = 'hello'}) => (api.CommentBuilder()
        ..id = id
        ..comment = text
        ..receiptId = 1
        ..userId = 1
        ..createdAt = DateTime.now().toIso8601String())
      .build();

  UserModel userModelWithTester() {
    final model = UserModel();
    model.setUsers([
      (api.UserViewBuilder()
            ..id = 1
            ..username = 'tester'
            ..displayName = 'Tester'
            ..isDummyUser = false)
          .build(),
    ]);
    return model;
  }

  // Key the widget under test so the slidable is located via find.descendant
  // rather than a brittle bare find.byType (house rule in mobile/CLAUDE.md).
  const commentsKey = ValueKey('comments-under-test');

  Widget wrap(
    List<Permission> groupPermissions,
    List<api.Comment> comments, {
    api.ReceiptRequirements? requirements,
  }) {
    return MultiProvider(
      providers: [
        // Test-owned notifiers are injected with create: so the provider owns
        // their lifecycle (mobile/CLAUDE.md).
        ChangeNotifierProvider(
          create: (_) => seededPermissions(
            group: {groupId: groupPermissions},
            receiptRequirements:
                requirements == null ? const {} : {groupId: requirements},
          ),
        ),
        ChangeNotifierProvider(
          create: (_) => ReceiptModel()
            ..setReceipt(
              getDefaultReceipt().rebuild((b) => b
                ..id = 1
                ..groupId = groupId),
              false,
            ),
        ),
        ChangeNotifierProvider(create: (_) => userModelWithTester()),
        ChangeNotifierProvider(create: (_) => AuthModel()),
      ],
      child: MaterialApp(
        home: Scaffold(
          body: ReceiptComments(
            key: commentsKey,
            comments: comments,
            formState: WranglerFormState.edit,
          ),
        ),
      ),
    );
  }

  List<SlidableWidget> slidables(WidgetTester tester) => tester
      .widgetList<SlidableWidget>(find.descendant(
        of: find.byKey(commentsKey),
        matching: find.byType(SlidableWidget),
      ))
      .toList();

  SlidableWidget firstSlidable(WidgetTester tester) => slidables(tester).first;

  group('ReceiptComments swipe-to-delete gate (edit state)', () {
    testWidgets('enabled with group.comments.delete', (tester) async {
      await tester.pumpWidget(
          wrap([Permission.groupPeriodCommentsPeriodDelete], [comment()]));
      await tester.pump();

      expect(firstSlidable(tester).slideEnabled, isTrue);
    });

    testWidgets('disabled without group.comments.delete', (tester) async {
      // Holds create but not delete (a Legacy Editor minus delete) — proves the
      // delete gate is independent of the create gate.
      await tester.pumpWidget(
          wrap([Permission.groupPeriodCommentsPeriodCreate], [comment()]));
      await tester.pump();

      expect(firstSlidable(tester).slideEnabled, isFalse);
    });
  });
  group('ReceiptComments keeps a role-required comment (edit state)', () {
    const canDelete = [Permission.groupPeriodCommentsPeriodDelete];

    testWidgets('the only comment cannot be swiped away', (tester) async {
      await tester.pumpWidget(wrap(canDelete, [comment()],
          requirements: receiptRequirements(comment: true)));
      await tester.pump();

      expect(firstSlidable(tester).slideEnabled, isFalse);
    });

    testWidgets('with two, either can go', (tester) async {
      await tester.pumpWidget(wrap(
          canDelete, [comment(id: 1), comment(id: 2, text: 'second')],
          requirements: receiptRequirements(comment: true)));
      await tester.pump();

      expect(slidables(tester).map((s) => s.slideEnabled), [true, true]);
    });

    testWidgets('a blank comment does not count, and stays deletable',
        (tester) async {
      // The server counts only non-blank comments, so the real one is the last
      // that satisfies the requirement.
      await tester.pumpWidget(wrap(
          canDelete, [comment(id: 1, text: '  '), comment(id: 2)],
          requirements: receiptRequirements(comment: true)));
      await tester.pump();

      expect(slidables(tester).map((s) => s.slideEnabled), [true, false]);
    });

    testWidgets('without the requirement the last comment is deletable',
        (tester) async {
      await tester.pumpWidget(wrap(canDelete, [comment()],
          requirements: receiptRequirements(image: true)));
      await tester.pump();

      expect(firstSlidable(tester).slideEnabled, isTrue);
    });
  });
}
