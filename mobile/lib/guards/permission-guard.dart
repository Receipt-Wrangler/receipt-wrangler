import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';

import '../constants/search.dart';
import '../models/permissions_model.dart';

/// go_router route-level redirects that gate permission-scoped screens, mirroring
/// the desktop route guards (`desktop/src/dashboard/guards/group-dashboard-read.guard.ts`).
/// The server re-checks permissions on every request, so these are UI hints that
/// keep users off pages they can't use (which would otherwise 403); they read
/// [PermissionsModel] from the provider tree the router is built under.

/// Allows the group dashboards route only when the caller holds
/// `group.dashboards.read` for the route's group. On deny, redirects to that
/// group's receipt list so the user lands on a page they can use, instead of
/// letting the dashboard fetch 403. The group id (for both the check and the
/// redirect) comes from the route's `:groupId` param, mirroring the desktop guard.
String? groupDashboardReadRedirect(BuildContext context, GoRouterState state) {
  final permissionsModel =
      Provider.of<PermissionsModel>(context, listen: false);
  final groupId = state.pathParameters["groupId"];
  final parsedGroupId = int.tryParse(groupId ?? "");

  if (parsedGroupId != null &&
      permissionsModel.hasGroupPermission(
          parsedGroupId, Permission.groupPeriodDashboardsPeriodRead)) {
    return null;
  }

  return "/groups/$groupId/receipts";
}

/// Allows the search route only when the caller holds `app.receipts.search`.
/// The search entry points are hidden for users without it (see the bottom
/// navs), so this is defense-in-depth against a deep link. On deny, lands on the
/// originating group's receipts when known (the group bottom nav passes a
/// `groupId` in `extra`), else the group-select screen. Reads `extra`
/// defensively so a no-`extra` deep link can't reach the search shell builder's
/// `state.extra as Map` cast.
String? receiptsSearchRedirect(BuildContext context, GoRouterState state) {
  final permissionsModel =
      Provider.of<PermissionsModel>(context, listen: false);

  if (permissionsModel
      .hasAppPermission(Permission.appPeriodReceiptsPeriodSearch)) {
    return null;
  }

  final extra = state.extra;
  if (extra is Map &&
      extra["from"] == fromGroupBottomNav &&
      extra["groupId"] != null) {
    return "/groups/${extra["groupId"]}/receipts";
  }
  return "/groups";
}
