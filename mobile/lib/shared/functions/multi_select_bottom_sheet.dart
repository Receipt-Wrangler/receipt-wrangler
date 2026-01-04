import 'package:flutter/material.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';

import '../../utils/bottom_sheet.dart';
import '../widgets/filter_multiselect.dart';

Future<dynamic> showMultiselectBottomSheet(
    context,
    String label,
    String buttonText,
    List<dynamic> options,
    List<dynamic> initialSelectedOptions,
    String Function(dynamic) itemDisplayName) {
  var navigator = Navigator.of(context);
  var multiselectKey = GlobalKey<FilterMultiSelectState<dynamic>>();

  return showFullscreenBottomSheet(
          context,
          FilterMultiSelect<dynamic>(
            key: multiselectKey,
            options: options,
            initialSelectedOptions: initialSelectedOptions,
            itemDisplayName: itemDisplayName,
          ),
          label,
          actions: [],
          bottomSheetWidget: BottomSubmitButton(
            onPressed: () {
              navigator.pop(multiselectKey.currentState!.selectedOptions ?? []);
            },
            buttonText: buttonText,
          ),
          bodyPadding: EdgeInsets.zero)
      .then((value) {
    return value;
  });
}
