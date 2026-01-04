import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';

String getGroupId(BuildContext context) {
  var extraMap =
      (GoRouterState.of(context).extra ?? {}) as Map<dynamic, dynamic>;
  return GoRouterState.of(context).pathParameters["groupId"] ??
      extraMap["groupId"] ??
      "0";
}

String? getGroupByIdWithRouter(GoRouter router) {
  return router.routerDelegate.currentConfiguration.pathParameters["groupId"];
}
