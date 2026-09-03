import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_entry.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';

import 'pump.dart';
import 'receipt_test_helpers.dart';

/// Shared Quick Scan e2e actions, used by the config + submit specs.

/// Asserts [field] is genuinely on screen in the Quick Scan sheet once the slide
/// is scrolled as far down as it will go.
///
/// None of the obvious assertions work here:
///  * `findsOneWidget` passes for a widget clipped outside its viewport;
///  * `Finder.hitTestable` reports 0 under `IntegrationTestWidgetsFlutterBinding`
///    on desktop even for plainly visible widgets;
///  * `tester.ensureVisible` walks *every* ancestor scrollable, so it dragged the
///    field into view even when the sheet was broken -- which is precisely how
///    this shipped (`quick_scan_submit_test` reaches its comment that way).
///
/// So: scroll the field's own scroll view to `maxScrollExtent` -- the furthest a
/// user could ever get -- and then check geometry. If the field still is not
/// fully inside the window at that point, nothing the user does will reveal it.
/// The viewport check additionally catches a field left buried under the pinned
/// submit button.
///
/// No surface size is forced: this runs at the target's natural viewport
/// (1280x720 for the Linux desktop runner, a real handset for Android/iOS).
/// `setSurfaceSize` is a no-op on the desktop target, so pinning a "phone" size
/// here would only be misleading.
Future<void> expectQuickScanFieldOnScreen(
  WidgetTester tester,
  Finder field, {
  required String label,
}) async {
  expect(field, findsOneWidget, reason: '$label is in the tree');

  final scrollable = Scrollable.of(tester.element(field));
  scrollable.position.jumpTo(scrollable.position.maxScrollExtent);
  await tester.pump(const Duration(milliseconds: 200));

  final fieldRect = tester.getRect(field);
  final window =
      Offset.zero & tester.view.physicalSize / tester.view.devicePixelRatio;
  final viewport = tester.getRect(find.byWidget(scrollable.widget));

  expect(window.contains(fieldRect.topLeft), isTrue,
      reason: '$label starts inside the window ($fieldRect vs $window)');
  expect(window.contains(fieldRect.bottomRight), isTrue,
      reason: '$label ends inside the window ($fieldRect vs $window)');
  expect(fieldRect.bottom <= viewport.bottom, isTrue,
      reason: '$label is not buried under the pinned submit button '
          '($fieldRect vs viewport $viewport)');
}

/// Finds a `FormBuilderDropdown` in the Quick Scan form by its `name`.
Finder quickScanDropdown(String name) =>
    find.byWidgetPredicate((w) => w is FormBuilderDropdown && w.name == name);

/// Finds the Quick Scan form's comment text field.
Finder quickScanCommentField() => find.byWidgetPredicate(
    (w) => w is FormBuilderTextField && w.name == 'comment');

/// Injects [configure] onto the admin's first non-all group and returns the
/// configured group (its `name` selects it in the form's group dropdown). The
/// config is a live-`GroupModel` Provider mutation -- deterministic and
/// independent of whether the local API persists the quick-scan fields.
Future<api.Group> configureFirstGroup(
  WidgetTester tester,
  void Function(api.GroupBuilder) configure,
) async {
  final ctx = tester.element(find.byType(Scaffold).first);
  final groupModel = Provider.of<GroupModel>(ctx, listen: false);
  final target = groupModel.groups.firstWhere((g) => !g.isAllGroup);
  final configured = target.rebuild(configure);
  groupModel.setGroups(groupModel.groups
      .map((g) => g.id == target.id ? configured : g)
      .toList());
  await tester.pump();
  return configured;
}

/// Trims the live GroupModel to the all-group placeholder(s) plus the group with
/// [groupId] (returned), so the group dropdown is short and the target is on
/// screen regardless of how many groups the admin has accumulated. Presentation
/// only -- it does not touch the group's (persisted) config.
Future<api.Group> keepOnlyGroup(WidgetTester tester, int groupId) async {
  final ctx = tester.element(find.byType(Scaffold).first);
  final groupModel = Provider.of<GroupModel>(ctx, listen: false);
  final target = groupModel.groups.firstWhere((g) => g.id == groupId);
  final allGroups = groupModel.groups.where((g) => g.isAllGroup).toList();
  groupModel.setGroups([...allGroups, target]);
  await tester.pump();
  return target;
}

/// Taps the bottom-nav scan slot to capture one image via the mocked document
/// scanner, leaving the Quick Scan sheet open with its per-image form mounted
/// (the "Group" field is on screen).
///
/// The tap is the direct scan action, so the sheet arrives already seeded --
/// there is no menu step and no separate "add a photo" tap any more. Requires
/// the AI flag already enabled and `installDocumentScannerMock()` already
/// installed. The hitTestable + frame-drain hardening handles the sheet
/// slide-in (a deterministic tap-flake on iOS Cupertino transitions).
Future<void> openQuickScanImageForm(WidgetTester tester) async {
  await pumpUntilFound(tester, scanNavSlot());
  await tester.tap(scanNavSlot());
  for (int i = 0; i < 5; i++) {
    await tester.pump(const Duration(milliseconds: 100));
  }

  await pumpUntilFound(tester, find.text('Group'));
}

/// Taps the form's submit button and waits for the fix-error snackbar (a
/// required+empty field short-circuits the submit before any API call).
Future<void> expectQuickScanSubmitBlocked(WidgetTester tester) async {
  await pumpUntilFound(tester, find.byType(BottomSubmitButton).hitTestable());
  await tester.tap(find.byType(BottomSubmitButton));
  await pumpUntilFound(
    tester,
    find.textContaining('Please fix error on quick scan'),
    timeout: const Duration(seconds: 10),
  );
}

/// Taps the form's submit button and waits for the success snackbar (the scan
/// was accepted and queued server-side), then asserts the sheet confirms it
/// inline.
///
/// The snackbar alone is not enough: submitting disables every field and hides
/// the submit button, so once it fades the sheet would sit greyed out with
/// nothing saying why. This is the only positive assertion on
/// `quick-scan-queued-confirmation` -- the widget suite can only assert its
/// absence, since the sheet builds its `isCompletedSubject` internally and a
/// real submit is what sets it.
Future<void> expectQuickScanQueued(WidgetTester tester) async {
  await pumpUntilFound(tester, find.byType(BottomSubmitButton).hitTestable());
  await tester.tap(find.byType(BottomSubmitButton));
  await pumpUntilFound(
    tester,
    find.textContaining('Successfully queued'),
    timeout: const Duration(seconds: 15),
  );

  await pumpUntilFound(
    tester,
    find.byKey(const ValueKey('quick-scan-queued-confirmation')),
    timeout: const Duration(seconds: 10),
  );
  expect(find.text(quickScanQueuedMessage), findsOneWidget);
}
