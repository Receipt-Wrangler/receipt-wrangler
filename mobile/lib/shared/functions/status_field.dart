import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:openapi/openapi.dart' as api;

import '../../enums/form_state.dart';
import '../../utils/forms.dart';

Widget receiptStatusField(
  String label,
  String fieldName,
  api.ReceiptStatus initialValue,
  WranglerFormState formState,
) {
  return FormBuilderDropdown(
    name: fieldName,
    decoration: InputDecoration(labelText: label),
    items: buildReceiptStatusDropDownMenuItems(),
    initialValue: initialValue,
    enabled: !isFieldReadOnly(formState),
    validator: FormBuilderValidators.required(),
  );
}

Widget itemStatusField(String label, String fieldName,
    api.ItemStatus initialValue, WranglerFormState formState,
    {Key? key}) {
  return FormBuilderDropdown(
    key: key,
    name: fieldName,
    decoration: InputDecoration(labelText: label),
    items: buildItemStatusDropDownMenuItems(),
    initialValue: initialValue,
    enabled: !isFieldReadOnly(formState),
    validator: FormBuilderValidators.required(),
  );
}
