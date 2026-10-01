// A group role can require its members to keep at least one comment and/or one
// image on the group's receipts. The server resolves that per user and group and
// delivers it on AppData.groupReceiptRequirements; the mobile client reads it to
// refuse a submit early, to disable deleting the last required item, and to
// show + require the Quick Scan comment. The server enforces all of it anyway.
//
// The widget suite injects the requirements into a mocked PermissionsModel, so
// it proves nothing about the wire. These specs run against a real role:
// the flag set on the role -> resolved into AppData at login -> read by the UI
// -> the single-call create the server accepts.
//
// Every member is a fresh user whose group role is a copy of Legacy Editor plus
// the requirement (`provisionMemberWithReceiptRequirements`); the admin that
// owns the group keeps Legacy Owner, so seeding over the admin API is unaffected.

import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/slidable_widget.dart';

import 'helpers/api.dart';
import 'helpers/document_scanner_mock.dart';
import 'helpers/feature_flags.dart';
import 'helpers/file_selector_mock.dart';
import 'helpers/form_actions.dart';
import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/permission_fixtures.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';
import 'helpers/quick_scan_actions.dart';
import 'helpers/receipt_test_helpers.dart';

void main() {
  final binding = IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  Future<void> useTallSurface() async {
    await binding.setSurfaceSize(const Size(1280, 900));
    addTearDown(() => binding.setSurfaceSize(null));
  }

  testWidgets(
    'add: blocked until an image and a comment are attached, then created '
    'in one call carrying both',
    (tester) async {
      await useTallSurface();
      await installFileSelectorMock();
      final fixture = await provisionMemberWithReceiptRequirements(
        comment: true,
        image: true,
      );
      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );

      final receiptName =
          'e2e-req-add-${DateTime.now().millisecondsSinceEpoch}';
      await openManualReceiptForm(tester);
      await tester.enterText(formField('name'), receiptName);
      await tester.enterText(formField('amount'), '12.34');
      await selectDropdown(tester, 'groupId', fixture.groupName!);
      await selectDropdown(tester, 'paidByUserId', fixture.displayName);
      await tester.pumpAndSettle(const Duration(seconds: 3));

      // Nothing attached: the client refuses before any request.
      await tester.tap(find.byType(BottomSubmitButton));
      await pumpUntilFound(
        tester,
        find.textContaining('one image and one comment'),
      );
      final before = await listReceiptsForGroup(
        fixture.groupId!,
        jwt: await apiLogin(),
      );
      expect(before.where((r) => r['name'] == receiptName), isEmpty);

      // Attach both; they ride the create.
      await attachFileFromReceiptForm(tester);
      await addCommentFromReceiptForm(tester, 'e2e required comment');

      // The blocked-submit snackbar can still cover the button; wait it out.
      await pumpUntilFound(
        tester,
        find.byType(BottomSubmitButton).hitTestable(),
        timeout: const Duration(seconds: 10),
      );
      final receiptId = await submitManualReceiptForm(tester);
      scheduleReceiptCleanup(receiptId);

      final receipt = await getReceipt(receiptId, jwt: await apiLogin());
      expect((receipt['imageFiles'] as List?)?.length ?? 0, 1);
      expect(
        ((receipt['comments'] as List?) ?? const []).map(
          (c) => (c as Map)['comment'],
        ),
        ['e2e required comment'],
      );
    },
  );

  testWidgets('edit: a receipt without the required comment cannot be saved', (
    tester,
  ) async {
    await useTallSurface();
    final fixture = await provisionMemberWithReceiptRequirements(
      comment: true,
      withReceipt: true,
    );
    await loginAs(
      tester,
      username: fixture.username,
      password: fixture.password,
    );

    await openReceiptEditForm(tester, fixture.groupName!, fixture.receiptName!);
    await tester.tap(find.byType(BottomSubmitButton));

    await pumpUntilFound(
      tester,
      find.textContaining('requires at least one comment'),
    );
    // Still on the edit form: the save never went out.
    expect(find.byType(BottomSubmitButton), findsOneWidget);
  });

  testWidgets('edit: the last required comment cannot be swiped away', (
    tester,
  ) async {
    await useTallSurface();
    final fixture = await provisionMemberWithReceiptRequirements(
      comment: true,
      withReceipt: true,
    );
    const onlyComment = 'e2e the only comment';
    await createComment(
      receiptId: fixture.receiptId!,
      jwt: await apiLogin(),
      comment: onlyComment,
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
    await pumpUntilFound(tester, find.text(onlyComment));

    // The member holds group.comments.delete (Legacy Editor), so the only thing
    // disabling the swipe is the role's requirement.
    final slidable = tester.widget<SlidableWidget>(
      find.ancestor(
        of: find.text(onlyComment),
        matching: find.byType(SlidableWidget),
      ),
    );
    expect(slidable.slideEnabled, isFalse);
  });

  testWidgets(
    'quick scan: a role-required comment is shown and required although the '
    'group config leaves it off',
    (tester) async {
      // A fresh fixture group's quick-scan config has the comment OFF, so the
      // field can only come from the role.
      await enableAiPoweredReceiptsForTest();
      await installDocumentScannerMock();
      final fixture = await provisionMemberWithReceiptRequirements(
        comment: true,
      );
      await loginAs(
        tester,
        username: fixture.username,
        password: fixture.password,
      );

      await openQuickScanImageForm(tester);
      await selectDropdown(tester, 'groupId', fixture.groupName!);

      await pumpUntilFound(tester, quickScanCommentField());
      // Required and empty: the submit short-circuits before any request.
      await expectQuickScanSubmitBlocked(tester);
    },
  );
}
