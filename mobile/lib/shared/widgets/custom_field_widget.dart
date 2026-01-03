import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/amount_field.dart';
import 'package:receipt_wrangler_mobile/utils/date.dart';

class CustomFieldWidget extends StatelessWidget {
  final api.CustomField customField;
  final api.CustomFieldValue? existingValue;
  final WranglerFormState formState;
  final VoidCallback? onRemove;

  const CustomFieldWidget({
    super.key,
    required this.customField,
    this.existingValue,
    required this.formState,
    this.onRemove,
  });

  String get fieldName => "customField_${customField.id}";

  bool get isReadOnly => formState == WranglerFormState.view;

  dynamic get initialValue {
    if (existingValue == null) return null;

    switch (customField.type) {
      case api.CustomFieldType.TEXT:
        return existingValue!.stringValue;
      case api.CustomFieldType.DATE:
        return existingValue!.dateValue != null
            ? DateTime.parse(existingValue!.dateValue!)
            : null;
      case api.CustomFieldType.SELECT:
        return existingValue!.selectValue;
      case api.CustomFieldType.CURRENCY:
        return existingValue!.currencyValue;
      case api.CustomFieldType.BOOLEAN:
        return existingValue!.booleanValue;
      default:
        return null;
    }
  }

  Widget buildTextField() {
    return FormBuilderTextField(
      name: fieldName,
      decoration: InputDecoration(
        labelText: customField.name,
        hintText: customField.description,
        suffixIcon: onRemove != null && !isReadOnly
            ? IconButton(
                icon: const Icon(Icons.remove_circle_outline),
                onPressed: onRemove,
              )
            : null,
      ),
      initialValue: initialValue as String?,
      readOnly: isReadOnly,
    );
  }

  Widget buildDateField() {
    if (isReadOnly && initialValue != null) {
      var formattedDate = formatDate(defaultDateFormat, initialValue as DateTime);
      return FormBuilderTextField(
        name: fieldName,
        decoration: InputDecoration(
          labelText: customField.name,
          hintText: customField.description,
          suffixIcon: onRemove != null
              ? IconButton(
                  icon: const Icon(Icons.remove_circle_outline),
                  onPressed: onRemove,
                )
              : null,
        ),
        initialValue: formattedDate,
        readOnly: true,
      );
    } else {
      return FormBuilderDateTimePicker(
        name: fieldName,
        decoration: InputDecoration(
          labelText: customField.name,
          hintText: customField.description,
          suffixIcon: onRemove != null && !isReadOnly
              ? IconButton(
                  icon: const Icon(Icons.remove_circle_outline),
                  onPressed: onRemove,
                )
              : null,
        ),
        initialValue: initialValue as DateTime?,
        inputType: InputType.date,
      );
    }
  }

  Widget buildSelectField() {
    List<DropdownMenuItem<int>> items = [];
    if (customField.options != null) {
      items = customField.options!
          .map((option) => DropdownMenuItem<int>(
                value: option.id,
                child: Text(option.value ?? ''),
              ))
          .toList();
    }

    return FormBuilderDropdown<int>(
      name: fieldName,
      decoration: InputDecoration(
        labelText: customField.name,
        hintText: customField.description,
        suffixIcon: onRemove != null && !isReadOnly
            ? IconButton(
                icon: const Icon(Icons.remove_circle_outline),
                onPressed: onRemove,
              )
            : null,
      ),
      items: items,
      initialValue: initialValue as int?,
      enabled: !isReadOnly,
    );
  }

  Widget buildCurrencyField() {
    return AmountField(
      label: customField.name,
      fieldName: fieldName,
      initialAmount: (initialValue as String?) ?? "0.00",
      formState: formState,
      suffixIcon: onRemove != null && !isReadOnly
          ? IconButton(
              icon: const Icon(Icons.remove_circle_outline),
              onPressed: onRemove,
            )
          : null,
    );
  }

  Widget buildBooleanField() {
    return FormBuilderCheckbox(
      name: fieldName,
      title: Text(customField.name),
      subtitle: customField.description != null
          ? Text(customField.description!)
          : null,
      initialValue: (initialValue as bool?) ?? false,
      enabled: !isReadOnly,
      decoration: InputDecoration(
        border: InputBorder.none,
        suffixIcon: onRemove != null && !isReadOnly
            ? IconButton(
                icon: const Icon(Icons.remove_circle_outline),
                onPressed: onRemove,
              )
            : null,
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    switch (customField.type) {
      case api.CustomFieldType.TEXT:
        return buildTextField();
      case api.CustomFieldType.DATE:
        return buildDateField();
      case api.CustomFieldType.SELECT:
        return buildSelectField();
      case api.CustomFieldType.CURRENCY:
        return buildCurrencyField();
      case api.CustomFieldType.BOOLEAN:
        return buildBooleanField();
      default:
        return Container();
    }
  }
}