import 'dart:async';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:receipt_wrangler_mobile/constants/search.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_nav.dart';
import 'package:receipt_wrangler_mobile/utils/group.dart';

import '../../../shared/functions/show_add_menu.dart';

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
    onDestinationSelected(int indexSelected) {
      var groupId = getGroupId(context);

      switch (indexSelected) {
        case 0:
          context.go("/groups/$groupId/dashboards");
          break;
        case 1:
          showAddMenu(context, addButtonKey);
          break;
        case 2:
          context.go("/groups/$groupId/receipts");
          break;
        case 3:
          context.go("/search",
              extra: {"from": fromGroupBottomNav, "groupId": groupId});
          break;
        default:
          context.go("/groups");
      }

      indexSelectedController.add(indexSelected);
    }

    setIndexSelected() {
      var uri =
          GoRouter.of(context).routeInformationProvider.value.uri.toString();
      var index = 0;

      if (uri.contains("dashboards")) {
        index = 0;
      } else if (uri.contains("/add")) {
        index = 1;
      } else if (uri.contains("receipts")) {
        index = 2;
      } else if (uri.contains("/search")) {
        index = 3;
      }

      return index;
    }

    var destinations = [
      const NavigationDestination(
        icon: Icon(Icons.dashboard),
        label: "Dashboards",
      ),
      NavigationDestination(
        key: addButtonKey,
        icon: Icon(Icons.add),
        label: "Add",
      ),
      const NavigationDestination(
        icon: Icon(Icons.receipt),
        label: "Receipts",
      ),
      const NavigationDestination(
        icon: Icon(Icons.search),
        label: "Search",
      ),
    ];

    return BottomNav(
      destinations: destinations,
      onDestinationSelected: onDestinationSelected,
      getInitialSelectedIndex: setIndexSelected,
      indexSelectedController: indexSelectedController,
    );
  }
}
