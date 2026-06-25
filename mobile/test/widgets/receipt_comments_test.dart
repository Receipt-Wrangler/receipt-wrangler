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

  api.Comment comment() => (api.CommentBuilder()
        ..id = 1
        ..comment = 'hello'
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

  Widget wrap(List<Permission> groupPermissions, List<api.Comment> comments) {
    final receiptModel = ReceiptModel();
    receiptModel.setReceipt(
      getDefaultReceipt().rebuild((b) => b
        ..id = 1
        ..groupId = groupId),
      false,
    );

    return MultiProvider(
      providers: [
        ChangeNotifierProvider.value(
          value: seededPermissions(group: {groupId: groupPermissions}),
        ),
        ChangeNotifierProvider.value(value: receiptModel),
        ChangeNotifierProvider.value(value: userModelWithTester()),
        ChangeNotifierProvider.value(value: AuthModel()),
      ],
      child: MaterialApp(
        home: Scaffold(
          body: ReceiptComments(
            comments: comments,
            formState: WranglerFormState.edit,
          ),
        ),
      ),
    );
  }

  SlidableWidget firstSlidable(WidgetTester tester) =>
      tester.widget<SlidableWidget>(find.byType(SlidableWidget).first);

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
}
