import 'package:flutter/material.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';

import 'api.dart';
import 'env.dart';
import 'permission_fixtures.dart';
import 'pump.dart';

// Shared drivers for the receipt form's Quick Actions split sheet, used by
// receipt_cost_split_test.dart (edit form) and
// receipt_add_quick_actions_split_test.dart (add form).

/// Creates a fresh group with admin as auto-Owner and `e2e-user` as Legacy
/// Editor, so the Quick Actions user picker (group-members-only) offers both
/// display names. Registers an `addTearDown` that deletes the group (which
/// cascades its receipts). Returns its id (for API checks) and name (for the groupId dropdown).
Future<({int id, String name})> provisionSharedSplitGroup() async {
  final jwt = await apiLogin(); // admin
  final userId = await userIdByUsername(E2eEnv.userUsername, jwt: jwt);
  final editorRoleId = await groupRoleIdByName('Legacy Editor', jwt: jwt);
  final groupName = 'e2e-split-grp-${DateTime.now().millisecondsSinceEpoch}';
  final groupId = await createGroupWithMember(
    name: groupName,
    memberUserId: userId,
    groupRoleId: editorRoleId,
    jwt: jwt,
  );
  addTearDown(() async => deleteGroup(groupId, jwt: await apiLogin()));
  return (id: groupId, name: groupName);
}

/// Locates the split-action IconButton on the receipt form. It's the
/// IconButton whose `icon` is an `SvgPicture` for `assets/icons/split.svg`
/// (receipt_form.dart:482-489). The neighboring "Add Share" IconButton
/// uses a Material `Icon`, not `SvgPicture`, so the predicate is unique.
Future<void> openQuickActionsSheet(WidgetTester tester) async {
  final splitButton = find.byWidgetPredicate(
    (w) => w is IconButton && w.icon is SvgPicture,
  );
  await pumpUntilFound(tester, splitButton);
  // The split-action button sits below the Shares row on the form,
  // off-screen on the 1280x900 test surface. Scroll it into view so
  // the tap lands -- otherwise the bottom sheet never opens.
  await tester.ensureVisible(splitButton);
  await tester.pumpAndSettle();
  await tester.tap(splitButton);
  // The fullscreen bottom sheet header is "Quick Actions"; wait for the
  // ToggleButtons row to render before driving the form. The text mounts on
  // the slide-in's first frame while the sheet is still rising -- a tap
  // computed then lands near the bottom edge and misses (observed as
  // "Offset(408.4, 852.0) ... would not hit test" on the By Percentage
  // toggle). Wait for hittability, then drain the slide-in.
  await pumpUntilFound(tester, find.text('Split Evenly').hitTestable());
  for (int i = 0; i < 5; i++) {
    await tester.pump(const Duration(milliseconds: 100));
  }
}

/// Opens the "Users" MultiSelectField, taps each ChoiceChip whose label
/// matches a display name in [displayNames], then taps the "Select"
/// confirm button. The MultiSelectField wraps its InputDecorator in an
/// opaque GestureDetector(onTap:) that fires `showUserMultiSelect`, so any
/// point inside the labeled "Users" field hits it.
Future<void> selectQuickActionsUsers(
  WidgetTester tester,
  List<String> displayNames,
) async {
  // The whole field is one tap target, so the decorator's center would work
  // too. The placeholder text is kept as the tap site because it is the
  // unambiguous locator for this field on a form full of decorators.
  await tester.tap(find.text('No Users selected'));
  await pumpUntilFound(tester, find.text('Select Users'));

  for (final name in displayNames) {
    // Wait for each chip -- the sheet's slide-in means the chips can be
    // mounted but mid-animation on the first loop iteration.
    final chip = find.widgetWithText(ChoiceChip, name);
    await pumpUntilFound(tester, chip);
    await tester.tap(chip);
    await tester.pump(const Duration(milliseconds: 200));
  }

  // The receipt-save success snackbar (root ScaffoldMessenger) renders ABOVE
  // the modal sheet and covers the bottom row where the Select button sits,
  // absorbing taps for its ~4s lifetime -- observed as "Offset(640.0, 875.0)
  // ... would not hit test on the specified widget". hitTestable() polls
  // until the snackbar departs so the tap actually lands.
  final selectButton = find.widgetWithText(BottomSubmitButton, 'Select');
  await pumpUntilFound(tester, selectButton.hitTestable());
  await tester.tap(selectButton.hitTestable());
  await tester.pumpAndSettle();
}
