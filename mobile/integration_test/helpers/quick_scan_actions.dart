import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';

import 'pump.dart';
import 'receipt_test_helpers.dart';

/// Shared Quick Scan e2e actions, used by the config + submit specs.

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
/// was accepted and queued server-side).
Future<void> expectQuickScanQueued(WidgetTester tester) async {
  await pumpUntilFound(tester, find.byType(BottomSubmitButton).hitTestable());
  await tester.tap(find.byType(BottomSubmitButton));
  await pumpUntilFound(
    tester,
    find.textContaining('Successfully queued'),
    timeout: const Duration(seconds: 15),
  );
}
