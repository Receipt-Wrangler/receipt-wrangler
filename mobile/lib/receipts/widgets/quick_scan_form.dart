import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/spacing.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/shared/classes/quick_scan_image.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/category_select_field.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/tag_select_field.dart';
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
  final void Function(
          int?, int?, api.ReceiptStatus?, List<api.Category>, List<api.Tag>)
      onFormChangeCallback;
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
      formValue["groupId"],
      formValue["paidByUserId"],
      formValue["status"],
      (formValue["categories"] as List?)?.cast<api.Category>() ?? const [],
      (formValue["tags"] as List?)?.cast<api.Tag>() ?? const [],
    );
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
        // The paid-by members and category/tag catalogs are group-scoped, so
        // clear them when the group changes. Fields may be hidden for either the
        // old or new group's config, so guard each access.
        widget.formKey.currentState?.fields["paidByUserId"]?.setValue(null);
        widget.formKey.currentState?.fields["categories"]
            ?.setValue(<api.Category>[]);
        widget.formKey.currentState?.fields["tags"]?.setValue(<api.Tag>[]);
        onValueChange();
        setState(() {
          groupId = value as int;
        });
      },
    );
  }

  Widget _buildUserDropDown(bool required) {
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
      validator: required ? FormBuilderValidators.required() : null,
      initialValue: valueExists ? initialValue : null,
      enabled: widget.enabled,
      // Set to null if value doesn't exist
      onChanged: (value) {
        onValueChange();
      },
    );
  }

  Widget _buildStatusDropdown(bool required) {
    api.ReceiptStatus? initialValue = widget.image.status;
    final items = buildReceiptStatusDropDownMenuItems();

    // Check if initialValue exists in the dropdown items
    bool valueExists = items.any((item) => item.value == initialValue);

    return FormBuilderDropdown(
      name: "status",
      decoration: const InputDecoration(labelText: "Status"),
      items: items,
      validator: required ? FormBuilderValidators.required() : null,
      initialValue: valueExists ? initialValue : null,
      enabled: widget.enabled,
      // Set to null if value doesn't exist
      onChanged: (value) {
        onValueChange();
      },
    );
  }

  Widget _buildCategoryField() {
    return CategorySelectField(
      fieldName: "categories",
      label: "Categories",
      groupId: groupId,
      initialCategories:
          (widget.formKey.currentState?.fields["categories"]?.value
                  as List<api.Category>?) ??
              widget.image.categories,
      formState: WranglerFormState.add,
      onCategoriesChanged: (categories) {
        setState(() {
          widget.formKey.currentState!.fields["categories"]!
              .setValue(categories);
        });
        onValueChange();
      },
    );
  }

  Widget _buildTagField() {
    return TagSelectField(
      fieldName: "tags",
      label: "Tags",
      groupId: groupId,
      initialTags: (widget.formKey.currentState?.fields["tags"]?.value
              as List<api.Tag>?) ??
          widget.image.tags,
      formState: WranglerFormState.add,
      onTagsChanged: (tags) {
        setState(() {
          widget.formKey.currentState!.fields["tags"]!.setValue(tags);
        });
        onValueChange();
      },
    );
  }

  @override
  Widget build(BuildContext context) {
    // Field visibility/requirement follows the selected group's quick-scan
    // config. When paid-by/status is not shown+required the server backfills a
    // configured default, so the field can be omitted here. Null (no group yet)
    // falls back to the backend defaults: paid-by/status shown, categories/tags
    // hidden.
    final settings = Provider.of<GroupModel>(context, listen: false)
        .getGroupReceiptSettings(groupId);

    final showPaidBy = settings?.quickScanPaidByEnabled ?? true;
    final requirePaidBy =
        showPaidBy && (settings?.quickScanPaidByRequired ?? true);
    final showStatus = settings?.quickScanStatusEnabled ?? true;
    final requireStatus =
        showStatus && (settings?.quickScanStatusRequired ?? true);
    final showCategories = settings?.quickScanCategoriesEnabled ?? false;
    final showTags = settings?.quickScanTagsEnabled ?? false;

    return FormBuilder(
        key: widget.formKey,
        child: Column(
          children: [
            textFieldSpacing,
            _buildGroupField(),
            Visibility(
              visible: showPaidBy,
              child: Column(children: [
                textFieldSpacing,
                _buildUserDropDown(requirePaidBy),
              ]),
            ),
            Visibility(
              visible: showStatus,
              child: Column(children: [
                textFieldSpacing,
                _buildStatusDropdown(requireStatus),
              ]),
            ),
            Visibility(
              visible: showCategories,
              child: Column(children: [
                textFieldSpacing,
                _buildCategoryField(),
              ]),
            ),
            Visibility(
              visible: showTags,
              child: Column(children: [
                textFieldSpacing,
                _buildTagField(),
              ]),
            ),
            submitButtonSpacing
          ],
        ));
  }
}
