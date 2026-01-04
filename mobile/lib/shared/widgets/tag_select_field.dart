import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/models/tag_model.dart';
import 'package:receipt_wrangler_mobile/shared/functions/multi_select_bottom_sheet.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/multi-select-field.dart';

class TagSelectField extends StatefulWidget {
  const TagSelectField({
    super.key,
    required this.label,
    required this.fieldName,
    required this.initialTags,
    required this.formState,
    required this.onTagsChanged,
  });

  final String label;

  final String fieldName;

  final List<api.Tag> initialTags;

  final WranglerFormState formState;

  final Function(List<api.Tag>)? onTagsChanged;

  @override
  State<TagSelectField> createState() => _TagSelectField();
}

class _TagSelectField extends State<TagSelectField> {
  late final tagModel = Provider.of<TagModel>(context, listen: false);
  late final contextModel = Provider.of<ContextModel>(context, listen: false);

  void showTagMultiSelect() {
    showMultiselectBottomSheet(
        contextModel.shellContext,
        "Select Tags",
        "Select",
        tagModel.tags,
        widget.initialTags,
        (tag) => tag.name).then((value) {
      if (value != null) {
        var tags = List<api.Tag>.from(value.map((item) => item as api.Tag));

        if (widget.onTagsChanged != null) {
          widget.onTagsChanged!(tags);
        }
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return MultiSelectField<api.Tag>(
        name: widget.fieldName,
        label: widget.label,
        initialValue: widget.initialTags,
        itemDisplayName: (tag) => tag.name ?? "",
        itemName: "Tags",
        onTap: widget.formState == WranglerFormState.view
            ? null
            : showTagMultiSelect);
  }
}
