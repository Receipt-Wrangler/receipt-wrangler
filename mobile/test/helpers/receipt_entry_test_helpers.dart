import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/auth_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';

/// Shared fixtures for the receipt-entry (Scan/Add) tests.
///
/// House style is real models seeded through their public setters rather than
/// mocks -- neither of these does I/O, and a real model proves the gate reads
/// the same field the app does.

/// An [AuthModel] whose feature config has the `aiPoweredReceipts` flag set to
/// [aiPoweredReceipts]. Every other flag is left at a fixed value so a test can
/// never accidentally depend on one.
AuthModel authModelWithAi(bool aiPoweredReceipts) {
  final model = AuthModel();
  model.setFeatureConfig((api.FeatureConfigBuilder()
        ..aiPoweredReceipts = aiPoweredReceipts
        ..enableLocalSignUp = false)
      .build());
  // The shared app bar renders a UserAvatar off these, and it dereferences them
  // without a null guard.
  model.setClaims(api.Claims((b) => b..displayName = 'Test User'));
  return model;
}

/// A [GroupModel] holding [groups].
GroupModel groupModelWith(List<api.Group> groups) {
  final model = GroupModel();
  model.setGroups(groups);
  return model;
}

/// Providers every receipt-entry affordance reads, wrapped around [router].
///
/// Kept here because the Scan slot, its menu, the overflow menu and the
/// orchestration functions all need the same three models, and a test that
/// forgets one fails with a `ProviderNotFoundException` rather than a useful
/// assertion.
Widget pumpReceiptEntryApp({
  required GoRouter router,
  required bool aiEnabled,
  required PermissionsModel permissions,
  List<api.Group> groups = const [],
}) {
  return MultiProvider(
    providers: [
      ChangeNotifierProvider<AuthModel>.value(value: authModelWithAi(aiEnabled)),
      ChangeNotifierProvider<GroupModel>.value(value: groupModelWith(groups)),
      ChangeNotifierProvider<PermissionsModel>.value(value: permissions),
      // Not read by the receipt-entry gates themselves, but the shared app bar
      // and bottom nav sit in the same trees these tests build.
      ChangeNotifierProvider<LoadingModel>(create: (_) => LoadingModel()),
      ChangeNotifierProvider<UserModel>(create: (_) => UserModel()),
    ],
    child: MaterialApp.router(routerConfig: router),
  );
}
