import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/slidable_widget.dart';

import 'helpers/api.dart';
import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

/// Permission-gating coverage for receipt comments on the edit-state comment
/// screen:
///   - the comment-input bottom sheet (`receipt_comment_screen.dart`) is hidden
///     without `group.comments.create`, and
///   - the comment swipe-to-delete (`receipt_comments.dart`,
///     `SlidableWidget.slideEnabled`) is disabled without `group.comments.delete`.
///
/// Both gates only apply in **edit** state, which requires `group.receipts.update`
/// to reach — so the members are provisioned from the "Legacy Editor" baseline
/// (which holds update + both comment perms) minus the single permission under
/// test. The existing `receipt_comments_test.dart` already proves the positive
/// path (an admin can add + delete), so this spec covers the deny paths.
void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets(
    'member without group.comments.create sees no comment input',
    (tester) async {
      await binding.setSurfaceSize(const Size(1280, 900));
      addTearDown(() => binding.setSurfaceSize(null));

      final fixture = await provisionGroupMemberWithoutPermission(
        'group.comments.create',
        baselineRole: 'Legacy Editor',
        withReceipt: true,
      );
      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );
      await openReceiptCommentsInEditMode(
        tester,
        fixture.groupName!,
        fixture.receiptName!,
      );

      // The comment-input bottom sheet renders SizedBox.shrink without
      // group.comments.create, so the field never mounts.
      expect(_commentField(), findsNothing);
    },
  );

  testWidgets(
    'member without group.comments.delete cannot swipe-delete a comment',
    (tester) async {
      await binding.setSurfaceSize(const Size(1280, 900));
      addTearDown(() => binding.setSurfaceSize(null));

      final fixture = await provisionGroupMemberWithoutPermission(
        'group.comments.delete',
        baselineRole: 'Legacy Editor',
        withReceipt: true,
      );
      // Seed a comment (as admin) so the swipe-to-delete gate has a row to show.
      const seededComment = 'e2e-seeded-comment';
      await createComment(
        receiptId: fixture.receiptId!,
        jwt: await apiLogin(),
        comment: seededComment,
      );

      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );
      await openReceiptCommentsInEditMode(
        tester,
        fixture.groupName!,
        fixture.receiptName!,
      );

      await pumpUntilFound(tester, find.text(seededComment));

      // The comment renders, but swipe-to-delete is disabled without
      // group.comments.delete.
      final slidable = tester.widget<SlidableWidget>(
        find.ancestor(
          of: find.text(seededComment),
          matching: find.byType(SlidableWidget),
        ),
      );
      expect(slidable.slideEnabled, isFalse);

      // Sanity: this member still holds group.comments.create, so the input is
      // present — proving the delete gate is independent of the create gate.
      expect(_commentField(), findsOneWidget);
    },
  );
}

Finder _commentField() => find.byWidgetPredicate(
      (w) => w is FormBuilderTextField && w.name == 'comment',
    );
