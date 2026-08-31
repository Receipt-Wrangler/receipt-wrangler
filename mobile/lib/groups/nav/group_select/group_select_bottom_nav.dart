import 'dart:async';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' show Permission;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/permissions_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_nav.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/scan_nav_item.dart';

import '../../../constants/search.dart';

const _groupsId = "groups";
const _searchId = "search";

class GroupSelectBottomNav extends StatefulWidget {
  const GroupSelectBottomNav({super.key});

  @override
  State<GroupSelectBottomNav> createState() => _GroupSelectBottomNav();
}

class _GroupSelectBottomNav extends State<GroupSelectBottomNav> {
  var indexSelectedController = StreamController<int>();
  final addButtonKey = GlobalKey();

  @override
  Widget build(BuildContext context) {
    final permissionsModel =
        Provider.of<PermissionsModel>(context, listen: false);
    final canSearch = permissionsModel
        .hasAppPermission(Permission.appPeriodReceiptsPeriodSearch);

    // Same id-keyed layout as the group nav: the scan slot and Search are both
    // permission-gated, so positions are not stable.
    final scanItem = buildScanNavItem(context, addButtonKey);

    var items = <NavDestinationItem>[
      const NavDestinationItem(
        id: _groupsId,
        destination: NavigationDestination(
          icon: Icon(Icons.group),
          label: "Groups",
        ),
      ),
      if (scanItem != null) scanItem,
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
      switch (items[indexSelected].id) {
        case _groupsId:
          context.go("/groups");
        case scanNavDestinationId:
          onScanNavItemSelected(context);
        case _searchId:
          context.go("/search", extra: {"from": fromGroupSelectBottomNav});
        default:
          context.go("/groups");
      }

      indexSelectedController.add(indexSelected);
    }

    setIndexSelected() {
      var uri =
          GoRouter.of(context).routeInformationProvider.value.uri.toString();
      var id = _groupsId;

      if (uri.contains("/groups")) {
        id = _groupsId;
      } else if (uri.contains("/add")) {
        id = scanNavDestinationId;
      } else if (uri.contains("/search")) {
        id = _searchId;
      }

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
