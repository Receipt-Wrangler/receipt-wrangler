import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/screen_wrapper.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/top_app_bar.dart';

showFullscreenBottomSheet(BuildContext context, Widget child, String label,
    {List<Widget>? actions,
    Widget? bottomSheetWidget,
    EdgeInsets? bodyPadding}) {
  return showModalBottomSheet(
    context: context,
    enableDrag: true,
    isDismissible: true,
    useSafeArea: true,
    isScrollControlled: true,
    showDragHandle: true,
    constraints: BoxConstraints(),
    builder: (BuildContext context) {
      return StatefulBuilder(builder:
          (BuildContext context, void Function(void Function()) setState) {
        return ScreenWrapper(
            bodyPadding: bodyPadding,
            appBarWidget: TopAppBar(
              titleText: label,
              actions: actions,
              hideAvatar: true,
              surfaceTintColor: Colors.white,
            ),
            bottomSheetWidget: bottomSheetWidget,
            child: SingleChildScrollView(
              child: child,
            ));
      });
    },
  );
}
