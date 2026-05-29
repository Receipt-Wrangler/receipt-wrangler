import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';

import 'env.dart';

/// Returns the display string the paidByUserId / chargedToUserId
/// dropdowns render for the e2e-admin user, looked up by username
/// against `UserModel.users`. The app already loads all users at
/// login (via /user/appData), so this is an in-process lookup with
/// no network round-trip.
///
/// Use this anywhere a test passes the admin's display string to
/// `selectDropdown`. Hard-coded display strings (e.g. "ee") are
/// environment-specific -- a user's displayName at sign-up differs
/// between local dev and the demo backend -- and break the suite when
/// run against any backend other than the one the test was written
/// against.
String adminDisplayName(WidgetTester tester) =>
    _displayNameForUsername(tester, E2eEnv.adminUsername);

/// Same shape as [adminDisplayName] but for the e2e-user, so multi-user
/// flows (group share, cost split) don't have to hard-code a displayName.
String userDisplayName(WidgetTester tester) =>
    _displayNameForUsername(tester, E2eEnv.userUsername);

String _displayNameForUsername(WidgetTester tester, String username) {
  final ctx = tester.element(find.byType(Scaffold).first);
  final userModel = Provider.of<UserModel>(ctx, listen: false);
  final user = userModel.users.firstWhere(
    (u) => u.username == username,
    orElse: () => throw StateError(
      'No user with username "$username" in UserModel. '
      'Available usernames: '
      '${userModel.users.map((u) => u.username).toList()}',
    ),
  );
  return user.displayName;
}
