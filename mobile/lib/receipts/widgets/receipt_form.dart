import 'package:built_collection/built_collection.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_svg/svg.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/spacing.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/quick_actions.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/quick_actions_submit_button.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_item_list.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/amount_field.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/tag_select_field.dart';
import 'package:receipt_wrangler_mobile/utils/bottom_sheet.dart';
import 'package:receipt_wrangler_mobile/utils/date.dart';
import 'package:receipt_wrangler_mobile/utils/forms.dart';

import '../../interfaces/form_item.dart';
import '../../models/category_model.dart';
import '../../models/context_model.dart';
import '../../models/custom_field_model.dart';
import '../../models/group_model.dart';
import '../../models/receipt_model.dart';
import '../../models/tag_model.dart';
import '../../shared/functions/forms.dart';
import '../../shared/functions/status_field.dart';
import '../../shared/widgets/audit_detail_section.dart';
import '../../shared/widgets/category_select_field.dart';
import '../../shared/widgets/custom_field_widget.dart';
import '../screens/receipt_image_screen.dart';
import '../screens/receipt_comment_screen.dart';

class ReceiptForm extends StatefulWidget {
  const ReceiptForm({super.key});

  @override
  State<ReceiptForm> createState() => _ReceiptForm();
}

class _ReceiptForm extends State<ReceiptForm> {
  late final receipt =
      Provider.of<ReceiptModel>(context, listen: false).receipt;
  late final receiptModel = Provider.of<ReceiptModel>(context, listen: false);
  late final formKey =
      Provider.of<ReceiptModel>(context, listen: false).receiptFormKey;
  late final groupModel = Provider.of<GroupModel>(context, listen: false);
  late final formState = getFormStateFromContext(context);
  late final categoryModel = Provider.of<CategoryModel>(context, listen: false);
  late final tagModel = Provider.of<TagModel>(context, listen: false);
  late final customFieldModel = Provider.of<CustomFieldModel>(context, listen: false);
  late final shellContext =
      Provider.of<ContextModel>(context, listen: false).shellContext;
  final addSharesFormKey = GlobalKey<FormBuilderState>();
  int groupId = 0;
  bool isAddingShare = false;

  @override
  void initState() {
    super.initState();
    Provider.of<ReceiptModel>(context, listen: false);

    groupId = modifiedReceipt.groupId;
  }

  Widget buildAuditDetailSection() {
    if (formState == WranglerFormState.add) {
      return SizedBox.shrink();
    }

    return AuditDetailSection(entity: receipt);
  }

  Widget buildNameField() {
    return FormBuilderTextField(
      name: "name",
      decoration: const InputDecoration(labelText: "Name"),
      initialValue: modifiedReceipt.name,
      validator: FormBuilderValidators.required(),
      readOnly: isFieldReadOnly(formState),
    );
  }

  Widget buildAmountField() {
    return AmountField(
        label: "Amount",
        fieldName: "amount",
        initialAmount: modifiedReceipt.amount.toString(),
        formState: formState);
  }

  Widget buildDateField() {
    if (formState == WranglerFormState.view) {
      var formattedDate =
          formatDate(defaultDateFormat, DateTime.parse(modifiedReceipt.date));
      return FormBuilderTextField(
          name: "date",
          decoration: const InputDecoration(labelText: "Date"),
          initialValue: formattedDate,
          readOnly: true);
    } else {
      return FormBuilderDateTimePicker(
        name: "date",
        decoration: const InputDecoration(labelText: "Date"),
        validator: FormBuilderValidators.required(),
        initialValue: DateTime.parse(modifiedReceipt.date),
        inputType: InputType.date,
      );
    }
  }

  Widget buildGroupField() {
    int? initialValue = modifiedReceipt.groupId;
    if (formState == WranglerFormState.add && initialValue == 0) {
      initialValue = null;
    }

    return FormBuilderDropdown(
      name: "groupId",
      decoration: const InputDecoration(labelText: "Group"),
      items: buildGroupDropDownMenuItems(context),
      initialValue: initialValue,
      enabled: !isFieldReadOnly(formState),
      validator: FormBuilderValidators.required(),
      onChanged: (value) {
        setState(() {
          formKey.currentState!.fields["paidByUserId"]!.setValue(null);
          groupId = value as int;
        });
      },
    );
  }

  Widget buildPaidByField() {
    List<DropdownMenuItem> items = [];
    var initialValue = null;
    if (groupId == modifiedReceipt.groupId) {
      initialValue = modifiedReceipt.paidByUserId;
    }

    if (groupId > 0) {
      items = buildGroupMemberDropDownMenuItems(context, groupId.toString());
    }

    if (formState == WranglerFormState.add && initialValue == 0) {
      initialValue = null;
    }

    return FormBuilderDropdown(
      name: "paidByUserId",
      decoration: const InputDecoration(labelText: "Paid By"),
      items: items.toList(),
      initialValue: initialValue,
      validator: FormBuilderValidators.required(),
      enabled: !isFieldReadOnly(formState),
    );
  }

  Widget buildStatusField() {
    return receiptStatusField(
        "Status", "status", modifiedReceipt.status, formState);
  }

  Widget buildCategoryField() {
    return CategorySelectField(
      fieldName: "categories",
      label: "Categories",
      initialCategories: formKey.currentState?.fields["categories"]?.value ??
          modifiedReceipt.categories!.toList(),
      formState: formState,
      onCategoriesChanged: (categories) => {
        setState(() {
          formKey.currentState!.fields["categories"]!.setValue(categories);
        }),
      },
    );
  }

  Widget buildTagField() {
    return TagSelectField(
        label: "Tags",
        fieldName: "tags",
        initialTags: formKey.currentState?.fields["tags"]?.value ??
            modifiedReceipt.tags!.toList(),
        formState: formState,
        onTagsChanged: (tags) => {
              setState(() {
                formKey.currentState!.fields["tags"]!.setValue(tags);
              })
            });
  }

  Widget buildCustomFieldsSection() {
    return Column(
      children: [
        // Render existing custom fields
        ...modifiedReceipt.customFields.map((customFieldValue) {
          final customField = customFieldModel.customFields
              .where((cf) => cf.id == customFieldValue.customFieldId)
              .firstOrNull;
          
          // Show loading placeholder if custom field template is not found but still loading
          if (customField == null) {
            if (customFieldModel.isLoading) {
              return Container(
                margin: const EdgeInsets.only(bottom: 16),
                child: const Row(
                  children: [
                    SizedBox(
                      width: 16,
                      height: 16,
                      child: CircularProgressIndicator(strokeWidth: 2),
                    ),
                    SizedBox(width: 8),
                    Text('Loading custom field...'),
                  ],
                ),
              );
            }
            return const SizedBox.shrink();
          }

          return Column(
            children: [
              CustomFieldWidget(
                customField: customField,
                existingValue: customFieldValue,
                formState: formState,
                onRemove: formState != WranglerFormState.view
                    ? () {
                        _removeCustomField(customFieldValue.customFieldId);
                      }
                    : null,
              ),
              textFieldSpacing,
            ],
          );
        }).toList(),
        
        // Add Custom Field button
        if (formState != WranglerFormState.view)
          buildAddCustomFieldButton(),
      ],
    );
  }

  Widget buildAddCustomFieldButton() {
    // Get custom field IDs that are already added
    final addedCustomFieldIds = modifiedReceipt.customFields
        .map((cfv) => cfv.customFieldId)
        .toSet();
    
    // Get available custom fields that haven't been added yet
    final availableCustomFields = customFieldModel.customFields
        .where((cf) => !addedCustomFieldIds.contains(cf.id))
        .toList();

    if (availableCustomFields.isEmpty) {
      return SizedBox.shrink();
    }

    return Align(
      alignment: Alignment.centerLeft,
      child: TextButton.icon(
        onPressed: () {
          showModalBottomSheet(
            context: context,
            builder: (context) => buildCustomFieldSelectionSheet(availableCustomFields),
          );
        },
        icon: Icon(Icons.add, color: Theme.of(context).primaryColor),
        label: Text(
          'Add Custom Field',
          style: TextStyle(color: Theme.of(context).primaryColor),
        ),
      ),
    );
  }

  Widget buildCustomFieldSelectionSheet(List<api.CustomField> availableCustomFields) {
    return Container(
      padding: const EdgeInsets.all(16),
      child: Column(
        mainAxisSize: MainAxisSize.min,
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Text(
            'Select Custom Field',
            style: Theme.of(context).textTheme.headlineSmall,
          ),
          const SizedBox(height: 16),
          Flexible(
            child: ListView.builder(
              shrinkWrap: true,
              itemCount: availableCustomFields.length,
              itemBuilder: (context, index) {
                final customField = availableCustomFields[index];
                return ListTile(
                  title: Text(customField.name),
                  subtitle: customField.description != null
                      ? Text(customField.description!)
                      : null,
                  trailing: Text(customField.type.name),
                  onTap: () {
                    _addCustomField(customField.id);
                    Navigator.of(context).pop();
                  },
                );
              },
            ),
          ),
        ],
      ),
    );
  }

  void _addCustomField(int customFieldId) {
    // Create a new empty custom field value with all required fields
    final newCustomFieldValue = (api.CustomFieldValueBuilder()
          ..id = 0  // Use 0 for new custom field values
          ..customFieldId = customFieldId
          ..receiptId = modifiedReceipt.id
          ..createdAt = DateTime.now().toIso8601String()  // Set current timestamp
          ..createdBy = 0  // Placeholder for user ID
          ..createdByString = ''  // Empty string placeholder
          ..updatedAt = '')  // Empty string placeholder
        .build();

    // Add it to the modified receipt
    final updatedCustomFields = [
      ...modifiedReceipt.customFields,
      newCustomFieldValue,
    ];

    final updatedReceipt = modifiedReceipt.rebuild((b) => b
      ..customFields = ListBuilder(updatedCustomFields));

    receiptModel.setModifiedReceipt(updatedReceipt);
  }

  void _removeCustomField(int customFieldId) {
    // Remove the custom field from the modified receipt
    final updatedCustomFields = modifiedReceipt.customFields
        .where((cfv) => cfv.customFieldId != customFieldId)
        .toList();

    final updatedReceipt = modifiedReceipt.rebuild((b) => b
      ..customFields = ListBuilder(updatedCustomFields));

    receiptModel.setModifiedReceipt(updatedReceipt);
  }

  Widget buildReceiptItemList() {
    return ReceiptItemField(
      groupId: groupId,
    );
  }

  Widget buildDetailsHeader() {
    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        buildHeaderText("Details"),
        Row(
          mainAxisSize: MainAxisSize.min,
          children: [
            _buildCompactActionButton(
              icon: Icons.image_outlined,
              label: "Images",
              onTap: () {
                Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (context) => ReceiptImageScreen(formState: formState),
                  ),
                );
              },
            ),
            const SizedBox(width: 12), // Increased spacing to prevent mis-taps
            _buildCompactActionButton(
              icon: Icons.chat_bubble_outline,
              label: "Comments", 
              onTap: () {
                Navigator.of(context).push(
                  MaterialPageRoute(
                    builder: (context) => ReceiptCommentScreen(formState: formState),
                  ),
                );
              },
            ),
          ],
        ),
      ],
    );
  }

  Widget _buildCompactActionButton({
    required IconData icon,
    required String label,
    required VoidCallback onTap,
  }) {
    return Tooltip(
      message: 'View $label',
      child: Material(
        color: Colors.transparent,
        child: InkWell(
          onTap: onTap,
          borderRadius: BorderRadius.circular(24),
          child: Container(
            // Minimum 48dp touch target for mobile accessibility
            constraints: const BoxConstraints(minHeight: 48, minWidth: 48),
            padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
            decoration: BoxDecoration(
              color: Theme.of(context).colorScheme.surface,
              borderRadius: BorderRadius.circular(24),
              border: Border.all(
                color: Theme.of(context).colorScheme.outline.withValues(alpha: 0.3),
                width: 1,
              ),
              // Add subtle shadow for better depth perception
              boxShadow: [
                BoxShadow(
                  color: Theme.of(context).colorScheme.shadow.withValues(alpha: 0.1),
                  blurRadius: 2,
                  offset: const Offset(0, 1),
                ),
              ],
            ),
            child: Row(
              mainAxisSize: MainAxisSize.min,
              children: [
                Icon(
                  icon,
                  size: 20, // Increased from 16 for better visibility
                  color: Theme.of(context).colorScheme.onSurfaceVariant,
                ),
                const SizedBox(width: 6), // Increased spacing
                Text(
                  label,
                  style: TextStyle(
                    fontSize: 14, // Increased from 12 for better readability
                    color: Theme.of(context).colorScheme.onSurfaceVariant,
                    fontWeight: FontWeight.w500,
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }


  Widget buildSharesHeader() {
    var rowChildren = [
      buildHeaderText("Shares"),
    ];

    if (formState != WranglerFormState.view) {
      rowChildren.add(Row(
        children: [
          IconButton(
            icon: Icon(Icons.add, color: Theme.of(context).primaryColor),
            onPressed: (isAddingShare || groupId == 0)
                ? null
                : () {
                    setState(() {
                      isAddingShare = true;
                    });
                  },
          ),
          IconButton(
              onPressed: openQuickActionsBottomSheet,
              icon: SvgPicture.asset(
                "assets/icons/split.svg",
                width: 24,
                height: 24,
              ))
        ],
      ));
    }

    return Row(
      mainAxisAlignment: MainAxisAlignment.spaceBetween,
      children: [
        ...rowChildren,
      ],
    );
  }

  void openQuickActionsBottomSheet() {
    receiptModel.resetQuickActionsFormKey();
    showFullscreenBottomSheet(
        shellContext as BuildContext,
        ReceiptQuickActions(
          groupId: groupId,
        ),
        "Quick Actions",
        bottomSheetWidget: const ReceiptQuickActionsSubmitButton());
  }

  Widget buildAddSharesCard() {
    if (isAddingShare) {
      return Card(
        color: Colors.white,
        surfaceTintColor: Colors.white,
        child: Padding(
          padding: const EdgeInsets.all(8.0),
          child: Column(
            children: [
              Row(children: [
                Text(
                  "Add Share",
                  style: TextStyle(fontSize: 15, fontWeight: FontWeight.bold),
                ),
              ]),
              textFieldSpacing,
              FormBuilder(
                key: addSharesFormKey,
                child: Column(
                  children: [
                    FormBuilderDropdown(
                      name: "chargedToUserId",
                      decoration:
                          const InputDecoration(labelText: "Shared With"),
                      items: buildGroupMemberDropDownMenuItems(
                          context, groupId.toString()),
                      initialValue: "",
                      validator: FormBuilderValidators.required(),
                      enabled: !isFieldReadOnly(formState),
                    ),
                    textFieldSpacing,
                    FormBuilderTextField(
                      name: "name",
                      decoration: InputDecoration(labelText: "Name"),
                      validator: FormBuilderValidators.required(),
                    ),
                    textFieldSpacing,
                    AmountField(
                        label: "Amount",
                        fieldName: "amount",
                        initialAmount: "0.00",
                        formState: formState),
                    Row(
                      mainAxisSize: MainAxisSize.max,
                      children: [
                        Expanded(
                          child: ElevatedButton(
                            child: Text("Add Share"),
                            onPressed: () {
                              if (!addSharesFormKey.currentState!
                                  .saveAndValidate()) {
                                return;
                              }
                              var form = addSharesFormKey.currentState!.value;

                              var items = [...receiptModel.items];
                              var newItem = (api.ItemBuilder()
                                    ..name = form["name"]
                                    ..amount = form["amount"]
                                    ..chargedToUserId = form["chargedToUserId"]
                                    ..receiptId = receipt?.id ?? 0
                                    ..status = api.ItemStatus.OPEN)
                                  .build();

                              items.add(FormItem.fromItem(newItem));

                              receiptModel.setItems(items);
                              setState(() {
                                isAddingShare = false;
                              });
                            },
                          ),
                        )
                      ],
                    ),
                    Row(
                      children: [
                        Expanded(
                          child: TextButton(
                            child: Text("Cancel"),
                            onPressed: () {
                              setState(() {
                                isAddingShare = false;
                              });
                            },
                          ),
                        )
                      ],
                    ),
                  ],
                ),
              )
              // user field,
              // item name
              // amount field,
            ],
          ),
        ),
      );
    } else {
      return SizedBox.shrink();
    }
  }

  // Get the current modified receipt from the model
  api.Receipt get modifiedReceipt => receiptModel.modifiedReceipt;

  @override
  Widget build(BuildContext context) {
    return Consumer<ReceiptModel>(
      builder: (context, receiptModel, child) {
        return FormBuilder(
          key: formKey,
          child: Column(
        mainAxisAlignment: MainAxisAlignment.start,
        children: [
          buildAuditDetailSection(),
          textFieldSpacing,
          buildDetailsHeader(),
          textFieldSpacing,
          buildNameField(),
          textFieldSpacing,
          buildAmountField(),
          textFieldSpacing,
          buildDateField(),
          textFieldSpacing,
          buildGroupField(),
          textFieldSpacing,
          buildPaidByField(),
          textFieldSpacing,
          buildStatusField(),
          textFieldSpacing,
          buildCustomFieldsSection(),
          textFieldSpacing,
          Visibility(
              visible: groupModel
                      .getGroupReceiptSettings(groupId)
                      ?.hideReceiptCategories ==
                  false,
              child: Column(children: [
                buildCategoryField(),
                textFieldSpacing,
              ])),
          Visibility(
              visible: groupModel
                      .getGroupReceiptSettings(groupId)
                      ?.hideReceiptTags ==
                  false,
              child: Column(children: [
                buildTagField(),
                textFieldSpacing,
              ])),
          buildSharesHeader(),
          textFieldSpacing,
          buildAddSharesCard(),
          textFieldSpacing,
          buildReceiptItemList(),
          textFieldSpacing,
          kDebugMode
              ? ElevatedButton(
                  onPressed: () => {
                        if (formKey.currentState!.saveAndValidate())
                          {print(formKey.currentState!.value)}
                      },
                  child: Text("Check form value"))
              : SizedBox.shrink(),
        ],
          ),
        );
      },
    );
  }
}
