import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/spacing.dart';
import 'package:receipt_wrangler_mobile/shared/classes/quick_scan_image.dart';
import 'package:receipt_wrangler_mobile/utils/forms.dart';

import '../../models/user_preferences_model.dart';

class QuickScanForm extends StatefulWidget {
  const QuickScanForm(
      {super.key,
      required this.formKey,
      required this.image,
      required this.index,
      required this.onFormChangeCallback,
      this.enabled = true});

  final GlobalKey<FormBuilderState> formKey;
  final QuickScanImage image;
  final int index;
  final void Function(int?, int?, api.ReceiptStatus?) onFormChangeCallback;
  final bool enabled;

  @override
  State<QuickScanForm> createState() => _QuickScanForm();
}

class _QuickScanForm extends State<QuickScanForm> {
  late final userPreferences =
      Provider.of<UserPreferencesModel>(context, listen: false).userPreferences;
  int groupId = 0;

  @override
  initState() {
    super.initState();
    groupId = widget.image.groupId ?? 0;
  }

  void onValueChange() {
    widget.formKey.currentState!.save();
    var formValue = widget.formKey.currentState!.value;
    widget.onFormChangeCallback(
        formValue["groupId"], formValue["paidByUserId"], formValue["status"]);
  }

  // TODO: refactor to a common Widget to use in receipt form
  Widget _buildGroupField() {
    // Get the list of groups for dropdown
    final dropdownItems = buildGroupDropDownMenuItems(context);
    int? initialValue = widget.image.groupId;

    // Check if initialValue exists in the dropdown items
    bool valueExists = dropdownItems.any((item) => item.value == initialValue);

    return FormBuilderDropdown(
      name: "groupId",
      decoration: const InputDecoration(labelText: "Group"),
      items: dropdownItems,
      validator: FormBuilderValidators.required(),
      initialValue: valueExists ? initialValue : null,
      enabled: widget.enabled,
      // Set to null if value doesn't exist
      onChanged: (value) {
        onValueChange();
        setState(() {
          widget.formKey.currentState!.fields["paidByUserId"]!.setValue(null);
          groupId = value as int;
        });
      },
    );
  }

  Widget _buildUserDropDown() {
    List<DropdownMenuItem> items = [];
    int? initialValue = widget.image.paidByUserId;

    if (groupId > 0) {
      items = buildGroupMemberDropDownMenuItems(context, groupId.toString());
    }

    // Check if initialValue exists in the dropdown items
    bool valueExists = items.any((item) => item.value == initialValue);

    return FormBuilderDropdown(
      name: "paidByUserId",
      decoration: const InputDecoration(labelText: "Paid By"),
      items: items,
      validator: FormBuilderValidators.required(),
      initialValue: valueExists ? initialValue : null,
      enabled: widget.enabled,
      // Set to null if value doesn't exist
      onChanged: (value) {
        onValueChange();
      },
    );
  }

  Widget _buildStatusDropdown() {
    api.ReceiptStatus? initialValue = widget.image.status;
    final items = buildReceiptStatusDropDownMenuItems();

    // Check if initialValue exists in the dropdown items
    bool valueExists = items.any((item) => item.value == initialValue);

    return FormBuilderDropdown(
      name: "status",
      decoration: const InputDecoration(labelText: "Status"),
      items: items,
      validator: FormBuilderValidators.required(),
      initialValue: valueExists ? initialValue : null,
      enabled: widget.enabled,
      // Set to null if value doesn't exist
      onChanged: (value) {
        onValueChange();
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    return FormBuilder(
        key: widget.formKey,
        child: Column(
          children: [
            textFieldSpacing,
            _buildGroupField(),
            textFieldSpacing,
            _buildUserDropDown(),
            textFieldSpacing,
            _buildStatusDropdown(),
            submitButtonSpacing
          ],
        ));
  }
}
