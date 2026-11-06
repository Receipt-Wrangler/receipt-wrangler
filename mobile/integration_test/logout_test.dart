import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/groups/screens/group_select.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';

import 'helpers/login.dart';
import 'helpers/nav.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

/// Regression guard for "Logout takes two taps".
///
/// The avatar menu's entries used to be TextButtons nested in PopupMenuItems
/// with no `onTap`, so only the label's box ran `_logout()`. A tap anywhere
/// else on the row closed the menu and left the user logged in, and they had to
/// reopen the menu and aim for the word. Measured on Linux desktop the button
/// was 71.5x32 of a 168x48 row.
///
/// This taps the row's far edge, clear of the label -- the tap the old layout
/// dropped -- once, and expects a real logout against the real API: the login
/// screen, and no stored JWT. The widget-level twin is
/// test/widgets/top_app_bar_avatar_menu_test.dart, which also covers the other
/// two entries.
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  testWidgets('one tap anywhere on the Logout row logs out', (tester) async {
    await loginAsAdmin(tester);
    final auth = Provider.of<AuthModel>(
      tester.element(find.byType(GroupSelect)),
      listen: false,
    );

    await openAvatarMenu(tester);

    final label = find.text('Logout');
    final rowRect = tester.getRect(find.ancestor(
      of: label,
      matching: find.byWidgetPredicate((w) => w is PopupMenuItem),
    ));
    final point = rowRect.centerRight - const Offset(12, 0);
    expect(tester.getRect(label).contains(point), isFalse,
        reason: 'the tap must land on the row, clear of the label');

    await tester.tapAt(point);

    await pumpUntilFound(tester, find.text('Log In'));
    // The logged-in page slides out after the login page mounts, so wait for
    // it to leave rather than asserting mid-transition.
    await pumpUntilGone(tester, find.byType(GroupSelect));
    expect(await auth.getJwt(), isNull);
  });
}
