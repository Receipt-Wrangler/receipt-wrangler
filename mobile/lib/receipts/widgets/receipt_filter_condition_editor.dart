import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/receipt_filter_fields.dart';
import 'package:receipt_wrangler_mobile/constants/spacing.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/models/group_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/shared/functions/multi_select_bottom_sheet.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/amount_field.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/category_select_field.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/multi-select-field.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/tag_select_field.dart';
import 'package:receipt_wrangler_mobile/utils/bottom_sheet.dart';
import 'package:receipt_wrangler_mobile/utils/date.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_filter_options.dart';

/// What the condition editor hands back to the filter sheet.
///
/// Three outcomes have to be distinguishable: saved (take this condition),
/// removed (drop the field entirely), and dismissed -- which returns null, so
/// the sheet leaves its draft exactly as it was.
class ReceiptFilterEditorResult {
  const ReceiptFilterEditorResult.saved(ReceiptFilterCondition this.condition)
      : removed = false;

  const ReceiptFilterEditorResult.removed()
      : condition = null,
        removed = true;

  final ReceiptFilterCondition? condition;

  final bool removed;
}

/// Opens the editor for one condition.
///
/// [existing] null means the field is being added, which is what hides the
/// Remove button.
Future<ReceiptFilterEditorResult?> showReceiptFilterConditionEditor(
  BuildContext context, {
  required ReceiptFilterField field,
  required String groupId,
  ReceiptFilterCondition? existing,
}) async {
  final result = await showFullscreenBottomSheet(
    context,
    ReceiptFilterConditionEditor(
      field: field,
      groupId: groupId,
      existing: existing,
    ),
    field.label,
  );

  return result is ReceiptFilterEditorResult ? result : null;
}

class ReceiptFilterConditionEditor extends StatefulWidget {
  const ReceiptFilterConditionEditor({
    super.key,
    required this.field,
    required this.groupId,
    this.existing,
  });

  final ReceiptFilterField field;

  /// The group being browsed, which scopes the category / tag / paid-by
  /// options. May be the synthetic "All" group -- see
  /// `utils/receipt_filter_options.dart`.
  final String groupId;

  final ReceiptFilterCondition? existing;

  @override
  State<ReceiptFilterConditionEditor> createState() =>
      ReceiptFilterConditionEditorState();
}

class ReceiptFilterConditionEditorState
    extends State<ReceiptFilterConditionEditor> {
  final _formKey = GlobalKey<FormBuilderState>();

  late api.FilterOperation _operation;

  dynamic _value;

  List<api.FilterOperation> get _operations =>
      operationsForFilterField(widget.field);

  bool get _isEditing => widget.existing != null;

  bool get _isValid => isReceiptFilterConditionValid(
      widget.field, ReceiptFilterCondition(operation: _operation, value: _value));

  @override
  void initState() {
    super.initState();

    _operation = widget.existing?.operation ??
        (_operations.isEmpty ? api.FilterOperation.CONTAINS : _operations.first);
    _value = defaultReceiptFilterValue(
        widget.field, _operation, widget.existing?.value);
  }

  /// Switching operation can change the value's *shape* (a single amount vs. a
  /// pair), so anything that no longer fits is dropped rather than carried into
  /// the encoder. [defaultReceiptFilterValue] owns that rule.
  void _pickOperation(api.FilterOperation operation) {
    setState(() {
      _operation = operation;
      _value = defaultReceiptFilterValue(widget.field, operation, _value);
    });
  }

  void _setValue(dynamic value) {
    setState(() => _value = value);
  }

  /// Re-reads the amount inputs whenever the form reports a change.
  ///
  /// Only the keys the current operation actually renders are read: FormBuilder
  /// runs with `clearValueOnUnregister: false`, so a field left behind by the
  /// previous operation still has an entry in the value map.
  void _syncAmountsFromForm() {
    final state = _formKey.currentState;
    if (state == null) {
      return;
    }

    if (_operation == api.FilterOperation.BETWEEN) {
      final start = _readAmount(state, _startKey);
      final end = _readAmount(state, _endKey);
      if (start == null || end == null) {
        _setValue(null);
        return;
      }
      _setValue([start, end]);
      return;
    }

    _setValue(_readAmount(state, _valueKey));
  }

  /// The amount in [key], or null when the user has not entered one.
  ///
  /// AmountField's `valueTransformer` maps an empty box to "0.00", so the
  /// transformed value cannot tell "untouched" from "typed zero" -- and zero is
  /// a filter the API honours (`amount = 0`), so guessing wrong silently hides
  /// every receipt. Emptiness therefore comes from the RAW value and the
  /// transformed one is only parsed once there is something to parse.
  double? _readAmount(FormBuilderState state, String key) {
    final raw = state.getRawValue<dynamic>(key);
    if (raw == null || (raw is String && raw.trim().isEmpty)) {
      return null;
    }
    return _parseAmount(state.instantValue[key]);
  }

  void _syncTextFromForm() {
    final state = _formKey.currentState;
    if (state == null) {
      return;
    }
    _setValue(state.instantValue[_valueKey] as String? ?? "");
  }

  /// AmountField's `valueTransformer` hands back a USD decimal string, the same
  /// unit the API stores amounts in.
  double? _parseAmount(dynamic raw) {
    if (raw is double) return raw;
    if (raw is num) return raw.toDouble();
    if (raw is String) return double.tryParse(raw);
    return null;
  }

  void _save() {
    if (!_isValid) {
      return;
    }
    Navigator.of(context).pop(ReceiptFilterEditorResult.saved(
        ReceiptFilterCondition(operation: _operation, value: _value)));
  }

  void _remove() {
    Navigator.of(context).pop(const ReceiptFilterEditorResult.removed());
  }

  @override
  Widget build(BuildContext context) {
    return FormBuilder(
      key: _formKey,
      onChanged: _onFormChanged,
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          ..._buildOperationSection(),
          _buildValueInput(),
          headerSpacing,
          _buildActions(),
        ],
      ),
    );
  }

  void _onFormChanged() {
    switch (widget.field.type) {
      case ReceiptFilterFieldType.number:
        _syncAmountsFromForm();
        break;
      case ReceiptFilterFieldType.text:
        _syncTextFromForm();
        break;
      default:
        // Dates and selections report through their own callbacks, which
        // already hold the typed value -- re-reading the form would only
        // round-trip it through a less precise representation.
        break;
    }
  }

  /// The operation chips.
  ///
  /// Rendered even for the list and users fields, which offer `CONTAINS` alone:
  /// a selected, inert "Contains" chip says what the condition means, and every
  /// editor keeps the same shape. Hiding it makes those fields look half-built.
  List<Widget> _buildOperationSection() {
    if (_operations.isEmpty) {
      return const [];
    }

    return [
      Text("Operation", style: Theme.of(context).textTheme.labelLarge),
      const SizedBox(height: 8),
      Wrap(
        spacing: 8,
        runSpacing: 8,
        children: _operations
            .map((operation) => ChoiceChip(
                  key: ValueKey("receipt-filter-operation-${operation.name}"),
                  label: Text(filterOperationLabels[operation] ?? operation.name),
                  selected: _operation == operation,
                  // The same selected treatment the shared MultiSelectField
                  // gives its chips, so the two chip rows on this sheet do not
                  // read as different kinds of control. M3's default resolves
                  // to the secondary slate, which looks disabled next to them.
                  selectedColor: Theme.of(context).primaryColor,
                  showCheckmark: false,
                  onSelected: (_) => _pickOperation(operation),
                ))
            .toList(),
      ),
      textFieldSpacing,
    ];
  }

  Widget _buildValueInput() {
    switch (widget.field.type) {
      case ReceiptFilterFieldType.text:
        return _buildTextInput();
      case ReceiptFilterFieldType.number:
        return _buildAmountInput();
      case ReceiptFilterFieldType.date:
        return _buildDateInput();
      case ReceiptFilterFieldType.list:
      case ReceiptFilterFieldType.users:
        return _buildSelectionInput();
    }
  }

  Widget _buildTextInput() {
    return FormBuilderTextField(
      key: const ValueKey("receipt-filter-text-value"),
      name: _valueKey,
      initialValue: _value as String? ?? "",
      decoration: InputDecoration(
        labelText: widget.field.label,
        hintText: widget.field.hint,
      ),
    );
  }

  Widget _buildAmountInput() {
    // The amount inputs are the shared currency field, so a filter is typed in
    // the same units, symbol and separators as every other amount in the app.
    // It seeds at zero rather than blank (CurrencyTextFieldController has no
    // empty state), which is why an amount condition reads as valid from the
    // moment it is added -- zero is itself a filter the API honours.
    if (_operation == api.FilterOperation.BETWEEN) {
      final pair = _value is List && (_value as List).length == 2
          ? _value as List
          : null;

      return Row(
        crossAxisAlignment: CrossAxisAlignment.center,
        children: [
          Expanded(
            child: AmountField(
              key: ValueKey("receipt-filter-amount-start-${_operation.name}"),
              label: "From",
              fieldName: _startKey,
              initialAmount: _amountText(pair?[0]),
              formState: WranglerFormState.edit,
              validator: _noValidation,
            ),
          ),
          const Padding(
            padding: EdgeInsets.symmetric(horizontal: 10),
            child: Text("and"),
          ),
          Expanded(
            child: AmountField(
              key: ValueKey("receipt-filter-amount-end-${_operation.name}"),
              label: "To",
              fieldName: _endKey,
              initialAmount: _amountText(pair?[1]),
              formState: WranglerFormState.edit,
              validator: _noValidation,
            ),
          ),
        ],
      );
    }

    return AmountField(
      key: ValueKey("receipt-filter-amount-${_operation.name}"),
      label: widget.field.label,
      fieldName: _valueKey,
      initialAmount: _amountText(_value),
      formState: WranglerFormState.edit,
      validator: _noValidation,
    );
  }

  Widget _buildDateInput() {
    if (_operation == api.FilterOperation.WITHIN_CURRENT_MONTH) {
      // The API pins this range itself (month start through today), so there is
      // nothing to pick -- say so rather than showing a dead field.
      return const Padding(
        key: ValueKey("receipt-filter-within-current-month-note"),
        padding: EdgeInsets.symmetric(vertical: 8),
        child: Text(
            "This month, from the 1st through today. The server sets the range."),
      );
    }

    if (_operation == api.FilterOperation.BETWEEN) {
      return _TapField(
        fieldKey: const ValueKey("receipt-filter-date-range"),
        label: widget.field.label,
        text: receiptFilterValueLabel(widget.field,
            ReceiptFilterCondition(operation: _operation, value: _value)),
        onTap: _pickDateRange,
      );
    }

    return _TapField(
      fieldKey: const ValueKey("receipt-filter-date"),
      label: widget.field.label,
      text: receiptFilterValueLabel(widget.field,
          ReceiptFilterCondition(operation: _operation, value: _value)),
      onTap: _pickDate,
    );
  }

  Future<void> _pickDate() async {
    final current = _value is DateTime ? _value as DateTime : DateTime.now();
    final picked = await showDatePicker(
      context: context,
      initialDate: current,
      firstDate: _firstSelectableDate,
      lastDate: _lastSelectableDate,
    );

    if (picked != null) {
      _setValue(picked);
    }
  }

  Future<void> _pickDateRange() async {
    final existing = _value is List && (_value as List).length == 2
        ? DateTimeRange(
            start: (_value as List)[0] as DateTime,
            end: (_value as List)[1] as DateTime)
        : null;

    final picked = await showDateRangePicker(
      context: context,
      initialDateRange: existing,
      firstDate: _firstSelectableDate,
      lastDate: _lastSelectableDate,
    );

    if (picked != null) {
      _setValue([picked.start, picked.end]);
    }
  }

  /// Categories and tags go through their own wrappers, which already source
  /// the caller's per-group catalog for this group id -- including the All
  /// group, which is a real row with a catalog of its own. The other three have
  /// no wrapper, so they pair the shared chip field with the shared picker.
  Widget _buildSelectionInput() {
    final groupModel = Provider.of<GroupModel>(context, listen: false);

    switch (widget.field.key) {
      case "categories":
        return CategorySelectField(
          key: const ValueKey("receipt-filter-categories"),
          label: widget.field.label,
          fieldName: _valueKey,
          initialCategories: _selection<api.Category>(),
          formState: WranglerFormState.edit,
          groupId: int.tryParse(widget.groupId) ?? 0,
          onCategoriesChanged: (categories) => _applySelection(categories),
        );

      case "tags":
        return TagSelectField(
          key: const ValueKey("receipt-filter-tags"),
          label: widget.field.label,
          fieldName: _valueKey,
          initialTags: _selection<api.Tag>(),
          formState: WranglerFormState.edit,
          groupId: int.tryParse(widget.groupId) ?? 0,
          onTagsChanged: (tags) => _applySelection(tags),
        );

      case "paidBy":
        return _buildMultiSelect<api.UserView>(
          widgetKey: const ValueKey("receipt-filter-paid-by"),
          sheetTitle: "Select ${widget.field.label}",
          itemName: "People",
          options: filterPaidByOptions(
              Provider.of<UserModel>(context, listen: false),
              groupModel,
              widget.groupId),
        );

      case "group":
        return _buildMultiSelect<api.Group>(
          widgetKey: const ValueKey("receipt-filter-group"),
          sheetTitle: "Select ${widget.field.label}",
          itemName: "Groups",
          options: filterGroupOptions(groupModel),
        );

      case "status":
        return _buildMultiSelect<api.ReceiptStatus>(
          widgetKey: const ValueKey("receipt-filter-status"),
          sheetTitle: "Select ${widget.field.label}",
          itemName: "Statuses",
          options: filterStatusOptions(),
        );
    }

    return const SizedBox.shrink();
  }

  /// The shared chip field plus the shared full-screen picker, for the three
  /// list fields that have no dedicated wrapper of their own.
  Widget _buildMultiSelect<T>({
    required Key widgetKey,
    required String sheetTitle,
    // The plural the empty state reads as "No <itemName> selected", which the
    // field label alone does not always give ("No Paid By selected").
    required String itemName,
    required List<T> options,
  }) {
    final selected = _selection<T>();

    void openPicker() {
      final contextModel = Provider.of<ContextModel>(context, listen: false);

      showMultiselectBottomSheet(
        contextModel.resolveSheetContext(context),
        sheetTitle,
        "Select",
        options,
        selected,
        (option) => receiptFilterOptionLabel(option),
      ).then((value) {
        // A dismissed sheet returns null, which means "no change" -- an empty
        // list is the legitimate "everything removed" value.
        if (value != null) {
          _applySelection(List<T>.from(value.map((item) => item as T)));
        }
      });
    }

    return MultiSelectField<T>(
      key: widgetKey,
      name: _valueKey,
      label: widget.field.label,
      initialValue: selected,
      itemDisplayName: (option) => receiptFilterOptionLabel(option),
      itemName: itemName,
      onTap: openPicker,
      onRemove: (remaining) => _applySelection(remaining),
    );
  }

  List<T> _selection<T>() {
    final value = _value;
    return value is List ? List<T>.from(value.whereType<T>()) : <T>[];
  }

  /// MultiSelectField does not write its own value -- the owner is the source of
  /// truth and answers with `setValue` plus a rebuild (mobile/CLAUDE.md,
  /// "Category / Tag / Users pickers").
  void _applySelection<T>(List<T> selection) {
    setState(() {
      _value = selection;
      _formKey.currentState?.fields[_valueKey]?.setValue(selection);
    });
  }

  Widget _buildActions() {
    return Row(
      children: [
        if (_isEditing) ...[
          OutlinedButton(
            key: const ValueKey("receipt-filter-condition-remove"),
            onPressed: _remove,
            child: const Text("Remove"),
          ),
          const SizedBox(width: 10),
        ],
        Expanded(
          child: SizedBox(
            height: 48,
            child: MaterialButton(
              key: const ValueKey("receipt-filter-condition-save"),
              onPressed: _isValid ? _save : null,
              color: Theme.of(context).primaryColor,
              child: const Text("Save condition",
                  style: TextStyle(color: Colors.white)),
            ),
          ),
        ),
      ],
    );
  }

  String _amountText(dynamic amount) {
    return amount is double ? amount.toString() : "0";
  }

  static String? _noValidation(String? _) => null;

  static const _valueKey = "value";
  static const _startKey = "valueStart";
  static const _endKey = "valueEnd";
}

/// Bounds wide enough that a receipt date is never unreachable, without letting
/// the picker scroll to the year 275760.
final DateTime _firstSelectableDate = DateTime(2000);
final DateTime _lastSelectableDate = DateTime(startOfDay(DateTime.now()).year + 5);

/// A read-only, tappable field that reads like the app's other inputs.
///
/// The same shape MultiSelectField uses: an opaque GestureDetector wrapping the
/// InputDecorator, so the label, the border and the padding gutters are all one
/// tap target rather than just the text.
class _TapField extends StatelessWidget {
  const _TapField({
    required this.fieldKey,
    required this.label,
    required this.text,
    required this.onTap,
  });

  final Key fieldKey;

  final String label;

  final String text;

  final VoidCallback onTap;

  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      key: fieldKey,
      behavior: HitTestBehavior.opaque,
      onTap: onTap,
      child: InputDecorator(
        decoration: InputDecoration(
          labelText: label,
          suffixIcon: const Icon(Icons.calendar_today),
        ),
        child: Text(text),
      ),
    );
  }
}
