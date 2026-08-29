import 'dart:async';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/search.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_nav.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/scan_nav_item.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';

const _dashboardsId = "dashboards";
const _receiptsId = "receipts";
const _searchId = "search";

class GroupBottomNav extends StatefulWidget {
  const GroupBottomNav({super.key});

  @override
  State<GroupBottomNav> createState() => _GroupBottomNav();
}

class _GroupBottomNav extends State<GroupBottomNav> {
  var indexSelectedController = StreamController<int>();
  final addButtonKey = GlobalKey();

  @override
  Widget build(BuildContext context) {
    final permissionsModel =
        Provider.of<PermissionsModel>(context, listen: false);
    final canSearch = permissionsModel
        .hasAppPermission(Permission.appPeriodReceiptsPeriodSearch);

    // Both the scan slot and Search are permission-gated, so a destination can be
    // missing and every index after it shifts. Everything below therefore keys
    // off the item ids rather than a hardcoded position.
    final scanItem = buildScanNavItem(context, addButtonKey);

    var items = <NavDestinationItem>[
      const NavDestinationItem(
        id: _dashboardsId,
        destination: NavigationDestination(
          icon: Icon(Icons.dashboard),
          label: "Dashboards",
        ),
      ),
      if (scanItem != null) scanItem,
      const NavDestinationItem(
        id: _receiptsId,
        destination: NavigationDestination(
          icon: Icon(Icons.receipt),
          label: "Receipts",
        ),
      ),
      if (canSearch)
        const NavDestinationItem(
          id: _searchId,
          destination: NavigationDestination(
            icon: Icon(Icons.search),
            label: "Search",
          ),
        ),
    ];

    onDestinationSelected(int indexSelected) {
      var groupId = getGroupId(context);

      switch (items[indexSelected].id) {
        case _dashboardsId:
          context.go("/groups/$groupId/dashboards");
        case scanNavDestinationId:
          onScanNavItemSelected(context);
        case _receiptsId:
          context.go("/groups/$groupId/receipts");
        case _searchId:
          context.go("/search",
              extra: {"from": fromGroupBottomNav, "groupId": groupId});
        default:
          context.go("/groups");
      }

      indexSelectedController.add(indexSelected);
    }

    setIndexSelected() {
      var uri =
          GoRouter.of(context).routeInformationProvider.value.uri.toString();
      var id = _dashboardsId;

      if (uri.contains("dashboards")) {
        id = _dashboardsId;
      } else if (uri.contains("/add")) {
        id = scanNavDestinationId;
      } else if (uri.contains("receipts")) {
        id = _receiptsId;
      } else if (uri.contains("/search")) {
        id = _searchId;
      }

      // A route whose slot is hidden (e.g. /receipts/add reached by a deep link
      // after the permission was revoked) has no destination to highlight, so
      // fall back to the first one.
      final index = items.indexWhere((item) => item.id == id);
      return index < 0 ? 0 : index;
    }

    return BottomNav(
      items: items,
      onDestinationSelected: onDestinationSelected,
      getInitialSelectedIndex: setIndexSelected,
      indexSelectedController: indexSelectedController,
    );
  }
}
