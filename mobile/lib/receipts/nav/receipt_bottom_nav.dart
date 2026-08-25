import 'dart:async';

import 'package:built_collection/built_collection.dart';
import 'package:flutter/material.dart';
import 'package:go_router/go_router.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/routes.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/custom_field_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/shared/functions/custom_field_values.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_nav.dart';
import 'package:receipt_wrangler_mobile/utils/date.dart';

import '../../utils/forms.dart';

class ReceiptBottomNav extends StatefulWidget {
  const ReceiptBottomNav({super.key});

  @override
  State<ReceiptBottomNav> createState() => _ReceiptBottomNav();
}

class _ReceiptBottomNav extends State<ReceiptBottomNav> {
  late final receiptModel = Provider.of<ReceiptModel>(context, listen: false);
  late final customFieldModel = Provider.of<CustomFieldModel>(context, listen: false);
  var indexSelectedController = StreamController<int>();
  var imagesAddedController = StreamController<api.FileDataView>.broadcast();

  void updateModifiedReceipt() {
    var formState = getFormStateFromContext(context);
    // Defensive null check: the only current caller already guards on
    // `currentState != null` (build()'s onDestinationSelected, line 138),
    // but bare `!` here is the same inconsistency the receipt-submit
    // handler had at receipt_bottom_sheet_builder.dart:389 -- a future
    // refactor that calls this from a non-guarded path would crash with
    // "Null check operator used on a null value". No-op when the form
    // isn't currently attached to the model's key.
    final state = receiptModel.receiptFormKey.currentState;
    if (state == null) {
      return;
    }
    state.save();
    var form = {...state.value};
    var date = "";

    if (formState == WranglerFormState.view) {
      // The view-mode date field is registered under "dateDisplay" rather
      // than "date" (see receipt_form.dart:97-110) so the read-only
      // String value doesn't collide with the edit-mode DateTime value
      // FormBuilder would otherwise reuse from its instant-value map.
      // Fall back to "date" for older snapshots that still set it, and
      // bail out cleanly if neither is present rather than letting
      // convertDateFormatForForm crash on null.
      final rawDate = form["dateDisplay"] ?? form["date"];
      if (rawDate == null) {
        return;
      }
      date = convertDateFormatForForm(rawDate);
    } else {
      try {
        date = formatDate(zuluDateFormat, form["date"] as DateTime);
      } catch (e) {
        var zuluDate = convertDateFormatForForm(form["date"]);
        date = formatDate(zuluDateFormat, DateTime.parse(zuluDate));
      }
    }

    // Rebuild every attached custom field value from the live form -- shared
    // with the submit payload builder so the two can't drift (empty values are
    // kept; a value whose template isn't loaded keeps what's already stored).
    var customFieldValues = buildCustomFieldValues(
      attachedValues: receiptModel.modifiedReceipt.customFields,
      customFields: customFieldModel.customFields,
      form: form,
    );

    var modifiedReceipt = (api.ReceiptBuilder()
          ..id = receiptModel.receipt.id
          ..name = form["name"] ?? ""
          ..amount = form["amount"] ?? "0"
          ..date = date
          ..groupId = int.parse((form["groupId"] ?? "0").toString())
          ..paidByUserId = int.parse((form["paidByUserId"] ?? "0").toString())
          ..status = form["status"] as api.ReceiptStatus
          ..comments = ListBuilder(receiptModel.comments ?? [])
          ..categories = ListBuilder(List<api.Category>.from(
              (form["categories"] ?? []).map((item) => item as api.Category)))
          ..tags = ListBuilder(List<api.Tag>.from(
              (form["tags"] ?? []).map((item) => item as api.Tag)))
          ..customFields = ListBuilder(customFieldValues))
        .build();

    receiptModel.setModifiedReceipt(modifiedReceipt!);
  }

  @override
  Widget build(BuildContext context) {
    var formState = getFormStateFromContext(context);
    var formStateName = formState.name;

    onDestinationSelected(int indexSelected) {
      var receipt = receiptModel.receipt;

      if (receiptModel.receiptFormKey.currentState != null) {
        updateModifiedReceipt();
      }

      if (formState != WranglerFormState.add) {
        switch (indexSelected) {
          case 0:
            context.go("/receipts/${receipt.id}/${formStateName}");
            break;
          default:
            context.go("/groups");
        }
      } else {
        switch (indexSelected) {
          case 0:
            context.go("/receipts/${formStateName}");
            break;
          default:
            context.go("/groups");
        }
      }
    }

    setIndexSelected() {
      var fullPath = GoRouterState.of(context).fullPath ?? "";
      if (fullPath == fullReceiptViewPath) {
        return 0;
      }


      return 0;
    }

    const destinations = [
      NavigationDestination(
        icon: Icon(Icons.receipt),
        label: "Receipt",
      ),
    ];

    return BottomNav(
      key: const Key("receipt_bottom_nav"),
      destinations: destinations,
      onDestinationSelected: onDestinationSelected,
      getInitialSelectedIndex: setIndexSelected,
      indexSelectedController: indexSelectedController,
    );
  }

  @override
  void dispose() {
    indexSelectedController.close();
    imagesAddedController.close();
    super.dispose();
  }
}
