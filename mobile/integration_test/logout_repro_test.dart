import 'dart:io' show Platform;

import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:integration_test/integration_test.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/groups/screens/group_select.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/top_app_bar.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/user_avatar.dart';

import 'helpers/login.dart';
import 'helpers/platform_mocks.dart';
import 'helpers/pump.dart';

/// Diagnostic repro: does ONE tap on Logout actually log out?
void main() {
  IntegrationTestWidgetsFlutterBinding.ensureInitialized();

  setUp(() {
    if (Platform.isLinux) {
      installLinuxDesktopMocks();
    }
  });

  Future<void> openAvatarMenu(WidgetTester tester) async {
    final avatar = find.descendant(
      of: find.byType(TopAppBar),
      matching: find.byType(UserAvatar),
    );
    await pumpUntilFound(tester, avatar);
    await tester.tap(avatar);
    await pumpUntilFound(tester, find.text('Logout').hitTestable());
    for (var i = 0; i < 5; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }
  }

  Future<void> report(WidgetTester tester, AuthModel auth, String label) async {
    // Let any in-flight logout request / navigation finish.
    for (var i = 0; i < 40; i++) {
      await tester.pump(const Duration(milliseconds: 100));
    }
    final jwt = await auth.getJwt();
    // ignore: avoid_print
    print('[$label] '
        'loginScreen=${find.text('Log In').evaluate().isNotEmpty} '
        'groupSelect=${find.byType(GroupSelect).evaluate().isNotEmpty} '
        'menuOpen=${find.text('User Profile').evaluate().isNotEmpty} '
        'jwtPresent=${jwt != null && jwt.isNotEmpty}');
  }

  testWidgets('one tap on the Logout TEXT', (tester) async {
    await loginAsAdmin(tester);
    final auth = Provider.of<AuthModel>(
        tester.element(find.byType(GroupSelect)),
        listen: false);

    await openAvatarMenu(tester);
    await tester.tap(find.text('Logout'));
    await report(tester, auth, 'text tap #1');

    expect(find.text('Log In'), findsOneWidget,
        reason: 'one tap on the Logout label should reach the login screen');
  });

  testWidgets('one tap on the Logout ROW, right of the label', (tester) async {
    await loginAsAdmin(tester);
    final auth = Provider.of<AuthModel>(
        tester.element(find.byType(GroupSelect)),
        listen: false);

    await openAvatarMenu(tester);
    final row = find.ancestor(
      of: find.text('Logout'),
      matching: find.byWidgetPredicate((w) => w is PopupMenuItem),
    );
    final rowRect = tester.getRect(row);
    final labelRect = tester.getRect(find.ancestor(
        of: find.text('Logout'), matching: find.byType(TextButton)));
    // ignore: avoid_print
    print('row=$rowRect textButton=$labelRect');

    await tester.tapAt(rowRect.centerRight - const Offset(12, 0));
    await report(tester, auth, 'row tap #1');

    expect(find.text('Log In'), findsOneWidget,
        reason: 'one tap on the Logout row should reach the login screen');
  });
}
