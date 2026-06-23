// TODO: fix delete, and fix adding new items, fix broken amount owed
import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/shared/functions/status_field.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/amount_field.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/category_select_field.dart';
import 'package:receipt_wrangler_mobile/utils/currency.dart';

import '../../interfaces/form_item.dart';
import '../../models/group_model.dart';
import '../../models/receipt_model.dart';
import '../../shared/widgets/tag_select_field.dart';
import '../../utils/forms.dart';

class ReceiptItemItems extends StatefulWidget {
  const ReceiptItemItems({
    super.key,
    required this.items,
    required this.groupId,
  });

  final List<FormItem> items;

  final int groupId;

  @override
  State<ReceiptItemItems> createState() => _ReceiptItemItems();
}

class _ReceiptItemItems extends State<ReceiptItemItems> {
  late final GroupModel groupModel =
      Provider.of<GroupModel>(context, listen: false);
  late final WranglerFormState formState = getFormStateFromContext(context);
  late final ReceiptModel receiptModel =
      Provider.of<ReceiptModel>(context, listen: false);
  late final GlobalKey<FormBuilderState> formKey =
      Provider.of<ReceiptModel>(context, listen: false).receiptFormKey;

  Map<int, List<FormItem>> getUserItemMap() {
    final Map<int, List<FormItem>> itemMap = <int, List<FormItem>>{};
    for (final FormItem item in widget.items) {
      itemMap.putIfAbsent(item.chargedToUserId, () => <FormItem>[]).add(item);
    }
    return itemMap;
  }

  String _safeItemName(FormItem item) {
    // Defensive against malformed/null values from generated client churn,
    // even though FormItem currently types name as String.
    final dynamic raw = item.name;
    if (raw is String) {
      final String trimmed = raw.trim();
      if (trimmed.isNotEmpty) {
        return trimmed;
      }
    }
    return "(unnamed item)";
  }

  String _safeItemAmount(FormItem item) {
    final dynamic raw = item.amount;
    if (raw is String && double.tryParse(raw) != null) {
      return raw;
    }
    if (raw is num) {
      return raw.toString();
    }
    return "0";
  }

  String _formatSubtotal(List<FormItem> items) {
    if (items.isEmpty) {
      return exchangeUSDToCustom("0").toString();
    }
    return items
        .map((FormItem item) => exchangeUSDToCustom(_safeItemAmount(item)))
        .reduce((value, element) => value + element)
        .toString();
  }

  Widget buildUserPanels(Map<int, List<FormItem>> userItemMap) {
    final List<Widget> panels = <Widget>[];
    for (final MapEntry<int, List<FormItem>> entry in userItemMap.entries) {
      panels.add(buildUserExpansionTile(entry.key, entry.value));
      panels.add(const SizedBox(height: 8));
    }
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: panels,
    );
  }

  Widget buildUserExpansionTile(int userId, List<FormItem> items) {
    final UserModel userModel =
        Provider.of<UserModel>(context, listen: false);
    final api.UserView? user = _safeGetUser(userModel, userId);
    final String displayName =
        (user?.displayName ?? "").trim().isEmpty ? "Unknown" : user!.displayName;
    final String subtotal = _formatSubtotal(items);

    return Card(
      margin: EdgeInsets.zero,
      child: ExpansionTile(
        initiallyExpanded: true,
        tilePadding: const EdgeInsets.symmetric(horizontal: 16),
        childrenPadding:
            const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        title: Text(displayName),
        subtitle: Text("${items.length} items • $subtotal"),
        children: <Widget>[
          ...items.map((FormItem item) => _buildItemEntry(item)),
          if (formState != WranglerFormState.view) buildAddItemButton(items),
        ],
      ),
    );
  }

  api.UserView? _safeGetUser(UserModel userModel, int userId) {
    try {
      return userModel.getUserById(userId.toString());
    } catch (_) {
      // UserModel.getUserById throws when the id is unknown (firstWhere
      // without orElse). Fall back to null instead of crashing the panel.
      return null;
    }
  }

  Widget _buildItemEntry(FormItem item) {
    final int itemIndex =
        widget.items.indexWhere((FormItem candidate) => candidate.formId == item.formId);
    if (formState == WranglerFormState.view) {
      return buildItemDisplayRow(item);
    }
    return buildItemRow(item, itemIndex);
  }

  Widget buildAddItemButton(List<FormItem> items) {
    return Padding(
      padding: const EdgeInsets.only(top: 8),
      child: ElevatedButton(
        onPressed: () {
          formKey.currentState?.save();
          final List<FormItem> newItems = <FormItem>[...widget.items];
          final api.Item newItem = (api.ItemBuilder()
                ..name = ""
                ..amount = "0.00"
                ..chargedToUserId = items.first.chargedToUserId
                ..receiptId = receiptModel.receipt.id
                ..status = api.ItemStatus.OPEN)
              .build();
          newItems.add(FormItem.fromItem(newItem));

          setState(() {
            receiptModel.setItems(newItems);
          });
        },
        child: const Text("Add Share"),
      ),
    );
  }

  Widget buildItemDisplayRow(FormItem item) {
    final ThemeData theme = Theme.of(context);
    final String name = _safeItemName(item);
    final String amount = exchangeUSDToCustom(_safeItemAmount(item)).toString();
    final List<Widget> chips = <Widget>[
      buildStatusChip(item.status),
    ];

    final dynamic rawCategories = item.categories;
    if (rawCategories is List) {
      for (final dynamic category in rawCategories) {
        if (category is api.Category) {
          final String label = (category.name ?? "").trim();
          if (label.isNotEmpty) {
            chips.add(buildMetaChip(label, theme.colorScheme.secondaryContainer));
          }
        }
      }
    }

    final dynamic rawTags = item.tags;
    if (rawTags is List) {
      for (final dynamic tag in rawTags) {
        if (tag is api.Tag) {
          final String label = tag.name.trim();
          if (label.isNotEmpty) {
            chips.add(buildMetaChip(label, theme.colorScheme.tertiaryContainer));
          }
        }
      }
    }

    return Padding(
      padding: const EdgeInsets.symmetric(vertical: 6),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: <Widget>[
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: <Widget>[
              Expanded(
                child: Text(
                  name,
                  style: theme.textTheme.bodyLarge,
                  overflow: TextOverflow.ellipsis,
                ),
              ),
              const SizedBox(width: 8),
              Text(amount, style: theme.textTheme.bodyLarge),
            ],
          ),
          const SizedBox(height: 6),
          Wrap(
            spacing: 6,
            runSpacing: 4,
            children: chips,
          ),
        ],
      ),
    );
  }

  Widget buildStatusChip(api.ItemStatus status) {
    final ThemeData theme = Theme.of(context);
    Color background;
    String label;
    switch (status) {
      case api.ItemStatus.OPEN:
        background = theme.colorScheme.primaryContainer;
        label = "Open";
        break;
      case api.ItemStatus.RESOLVED:
        background = theme.colorScheme.secondaryContainer;
        label = "Resolved";
        break;
      case api.ItemStatus.DRAFT:
        background = theme.colorScheme.secondaryContainer;
        label = "Draft";
        break;
      default:
        background = theme.colorScheme.secondaryContainer;
        label = status.name;
    }
    return buildMetaChip(label, background);
  }

  Widget buildMetaChip(String label, Color background) {
    return Chip(
      label: Text(label),
      backgroundColor: background,
      visualDensity: VisualDensity.compact,
      materialTapTargetSize: MaterialTapTargetSize.shrinkWrap,
    );
  }

  Widget buildItemRow(FormItem item, int index) {
    final String itemName = FormItem.buildItemNameName(item);
    final String amountName = FormItem.buildItemAmountName(item);
    final String statusName = FormItem.buildItemStatusName(item);
    final String categoryName = FormItem.buildItemCategoryName(item);
    final String tagName = FormItem.buildItemTagName(item);

    // FormBuilder field values come back as `dynamic`; guard each one with a
    // runtime type-check before passing into a typed field. This protects
    // against generated-client churn where a wire field changes shape and
    // would otherwise crash with a CastError on a hot reload.
    final dynamic rawName =
        formKey.currentState?.fields[itemName]?.value ?? item.name;
    final String? initialName = rawName is String ? rawName : null;

    final dynamic rawAmount =
        formKey.currentState?.fields[amountName]?.value ?? item.amount;
    final String initialAmount = rawAmount is String ? rawAmount : "0.00";

    final dynamic rawStatus =
        formKey.currentState?.fields[statusName]?.value ?? item.status;
    final api.ItemStatus initialStatus =
        rawStatus is api.ItemStatus ? rawStatus : api.ItemStatus.OPEN;

    final dynamic rawCategories =
        formKey.currentState?.fields[categoryName]?.value ?? item.categories;
    final List<api.Category> initialCategories = rawCategories is List
        ? rawCategories.whereType<api.Category>().toList()
        : <api.Category>[];

    final dynamic rawTags =
        formKey.currentState?.fields[tagName]?.value ?? item.tags;
    final List<api.Tag> initialTags = rawTags is List
        ? rawTags.whereType<api.Tag>().toList()
        : <api.Tag>[];

    Widget iconButton = const SizedBox.shrink();
    if (!isFieldReadOnly(formState)) {
      iconButton = IconButton(
        icon: const Icon(Icons.delete, color: Colors.red),
        onPressed: () {
          final List<FormItem> newItems = <FormItem>[...widget.items];
          newItems.removeAt(index);

          setState(() {
            receiptModel.setItems(newItems);
          });
        },
      );
    }
    // TODO: need to fix new item data being wiped out when adding
    return Padding(
      padding: const EdgeInsets.only(bottom: 12),
      child: Column(
        children: <Widget>[
          Row(
            children: <Widget>[
              Expanded(
                key: Key(itemName),
                child: FormBuilderTextField(
                  name: itemName,
                  initialValue: initialName,
                  decoration: const InputDecoration(label: Text("Name")),
                  readOnly: isFieldReadOnly(formState),
                ),
              ),
              Expanded(
                key: Key(amountName),
                child: AmountField(
                  label: "Amount",
                  fieldName: amountName,
                  initialAmount: initialAmount,
                  formState: formState,
                ),
              ),
              Expanded(
                key: Key(statusName),
                child: itemStatusField(
                  "Status",
                  statusName,
                  initialStatus,
                  formState,
                ),
              ),
              iconButton,
            ],
          ),
          Visibility(
            visible: groupModel
                    .getGroupReceiptSettings(widget.groupId)
                    ?.hideItemCategories ==
                false,
            child: CategorySelectField(
              label: "Categories",
              fieldName: categoryName,
              initialCategories: initialCategories,
              formState: formState,
              onCategoriesChanged: (categories) {
                setState(() {
                  formKey.currentState?.fields[categoryName]
                      ?.setValue(categories);
                });
              },
            ),
          ),
          Visibility(
            visible: groupModel
                    .getGroupReceiptSettings(widget.groupId)
                    ?.hideItemTags ==
                false,
            child: TagSelectField(
              label: "Tags",
              fieldName: tagName,
              initialTags: initialTags,
              formState: formState,
              onTagsChanged: (tags) {
                setState(() {
                  formKey.currentState?.fields[tagName]?.setValue(tags);
                });
              },
            ),
          ),
        ],
      ),
    );
  }

  @override
  Widget build(BuildContext context) {
    final Map<int, List<FormItem>> userItemMap = getUserItemMap();

    if (userItemMap.isEmpty) {
      return const Padding(
        padding: EdgeInsets.symmetric(vertical: 8),
        child: Text("No items on this receipt"),
      );
    }

    return buildUserPanels(userItemMap);
  }
}
