import 'dart:async';

import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_nav.dart';

import '../../../constants/search.dart';
import '../../../shared/functions/show_add_menu.dart';

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
    onDestinationSelected(int indexSelected) {
      switch (indexSelected) {
        case 0:
          context.go("/groups");
          break;
        case 1:
          showAddMenu(context, addButtonKey);
          break;
        case 2:
          context.go("/search", extra: {"from": fromGroupSelectBottomNav});
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

      if (uri.contains("/groups")) {
        index = 0;
      } else if (uri.contains("/add")) {
        index = 1;
      } else if (uri.contains("/search")) {
        index = 2;
      }

      return index;
    }

    var destinations = [
      const NavigationDestination(
        icon: Icon(Icons.group),
        label: "Groups",
      ),
      NavigationDestination(
        key: addButtonKey,
        icon: Icon(Icons.add),
        label: "Add",
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
