import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/category_model.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/shared/functions/multi_select_bottom_sheet.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/multi-select-field.dart';

class CategorySelectField extends StatefulWidget {
  const CategorySelectField({
    super.key,
    required this.label,
    required this.fieldName,
    required this.initialCategories,
    required this.formState,
    required this.onCategoriesChanged,
  });

  final String label;

  final String fieldName;

  final List<api.Category> initialCategories;

  final WranglerFormState formState;

  final Function(List<api.Category>)? onCategoriesChanged;

  @override
  State<CategorySelectField> createState() => _CategorySelectField();
}

class _CategorySelectField extends State<CategorySelectField> {
  late final categoryModel = Provider.of<CategoryModel>(context, listen: false);
  late final contextModel = Provider.of<ContextModel>(context, listen: false);

  void showCategoryMultiSelect() {
    showMultiselectBottomSheet(
        contextModel.shellContext,
        "Select Categories",
        "Select",
        categoryModel.categories,
        widget.initialCategories,
        (category) => category.name).then((value) {
      if (value != null) {
        var categories =
            List<api.Category>.from(value.map((item) => item as api.Category));

        if (widget.onCategoriesChanged != null) {
          widget.onCategoriesChanged!(categories);
        }
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return MultiSelectField<api.Category>(
        name: widget.fieldName,
        label: widget.label,
        initialValue: widget.initialCategories,
        itemDisplayName: (category) => category.name ?? "",
        itemName: "Categories",
        onTap: widget.formState == WranglerFormState.view
            ? null
            : showCategoryMultiSelect);
  }
}
