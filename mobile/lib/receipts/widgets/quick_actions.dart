import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/spacing.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/amount_field.dart';
import 'package:receipt_wrangler_mobile/utils/users.dart';

import '../../models/context_model.dart';
import '../../models/group_model.dart';
import '../../models/receipt_model.dart';
import '../../shared/functions/multi_select_bottom_sheet.dart';
import '../../shared/widgets/multi-select-field.dart';
import '../../shared/widgets/total_display_widget.dart';
import '../../utils/currency.dart';

class ReceiptQuickActions extends StatefulWidget {
  const ReceiptQuickActions({super.key, required this.groupId});

  final int groupId;

  @override
  State<ReceiptQuickActions> createState() => _ReceiptQuickActions();
}

class _ReceiptQuickActions extends State<ReceiptQuickActions> {
  late final shellContext =
      Provider.of<ContextModel>(context, listen: false).shellContext;
  late final userModel = Provider.of<UserModel>(context, listen: false);
  late final groupModel = Provider.of<GroupModel>(context, listen: false);
  late final formKey =
      Provider.of<ReceiptModel>(context, listen: false).quickActionsFormKey;
  late final receiptModel = Provider.of<ReceiptModel>(context, listen: false);
  final formState = WranglerFormState.edit;

  var quickActions = [
    "Split Evenly",
    "With Portions",
    "By Percentage"
  ];
  var quickActionsSelection = [true, false, false];

  var percentageOptions = ["25%", "50%", "75%", "100%", "Custom"];

  double calculateTotalPercentage() {
    var users =
        formKey.currentState?.fields["users"]?.value as List<api.UserView>? ??
            [];
    double total = 0.0;

    for (api.UserView user in users) {
      var userPercentageSelection = formKey
          .currentState?.fields["percentage_${user.id}"]?.value as String?;

      if (userPercentageSelection != null &&
          userPercentageSelection.isNotEmpty) {
        switch (userPercentageSelection) {
          case "25%":
            total += 25.0;
            break;
          case "50%":
            total += 50.0;
            break;
          case "75%":
            total += 75.0;
            break;
          case "100%":
            total += 100.0;
            break;
          case "Custom":
            var customPercentageStr = formKey
                .currentState?.fields["custom_percentage_${user.id}"]?.value
                ?.toString();
            if (customPercentageStr != null && customPercentageStr.isNotEmpty) {
              final customPercentage = double.tryParse(customPercentageStr);
              if (customPercentage != null) {
                total += customPercentage;
              }
            }
            break;
        }
      }
    }

    return total;
  }

  bool isTotalPercentageValid() {
    return calculateTotalPercentage() <= 100.0;
  }

  String getReceiptAmount() {
    return receiptModel.receiptFormKey.currentState?.fields["amount"]?.value ??
        "0";
  }

  String formatCurrencyWithNegatives(BuildContext context, double amount) {
    if (amount < 0) {
      // Format positive amount and add negative sign
      String positiveAmount = formatCurrency(context, (-amount).toString()) ?? "\$0.00";
      return "-$positiveAmount";
    } else {
      return formatCurrency(context, amount.toString()) ?? "\$0.00";
    }
  }

  List<Widget> buildSplitEvenlyTotal() {
    List<Widget> fields = [];
    if (quickActionsSelection[0]) {
      var users =
          formKey.currentState?.fields["users"]?.value as List<api.UserView>? ??
              [];

      if (users.isNotEmpty) {
        // Get receipt amount in display currency and convert to USD for calculation
        String receiptAmountDisplay = getReceiptAmount();
        double receiptAmountUSD = exchangeCustomToUSD(receiptAmountDisplay).toDouble();
        
        // Calculate per-person amount in USD
        double perPersonAmountUSD = receiptAmountUSD / users.length;
        
        // Convert back to display currency for showing to user
        String perPersonAmountDisplay = formatCurrency(context, perPersonAmountUSD.toString()) ?? "\$0.00";

        fields.add(const SizedBox(height: 10));
        fields.add(TotalDisplayWidget(
          value: "${users.length} users × $perPersonAmountDisplay each",
          isValid: true,
        ));
      } else {
        fields.add(const SizedBox(height: 10));
        fields.add(Text(
          "Please select users to split the receipt evenly among them.",
          style: TextStyle(color: Colors.grey[600]),
        ));
      }
    }
    return fields;
  }

  Widget buildUserField() {
    return MultiSelectField<api.UserView>(
        name: "users",
        label: "Users",
        required: true,
        initialValue: [],
        itemDisplayName: (user) => user.displayName ?? "",
        itemName: "Users",
        onTap: showUserMultiSelect);
  }

  void showUserMultiSelect() {
    showMultiselectBottomSheet(
        shellContext,
        "Select Users",
        "Select",
        getUsersInGroup(userModel, groupModel, widget.groupId.toString()),
        [],
        (user) => user.displayName).then((value) {
      if (value != null) {
        var users =
            List<api.UserView>.from(value.map((item) => item as api.UserView));

        setState(() {
          formKey.currentState!.fields["users"]!.setValue(users);
        });
      }
    });
  }

  Widget buildToggleButtons() {
    return ToggleButtons(
      children: [
        ...quickActions.map((action) => Padding(
              padding: EdgeInsets.all(4),
              child: Text(action),
            )),
      ],
      isSelected: quickActionsSelection,
      onPressed: (selected) => setState(() {
        if (quickActionsSelection[selected]) {
          return;
        }

        var newState = !quickActionsSelection[selected];

        quickActionsSelection = [false, false, false];
        quickActionsSelection[selected] = newState;

        formKey.currentState!.fields["quickActionsSelection"]!
            .setValue(quickActionsSelection);
      }),
    );
  }

  double calculateTotalPortionsInUSD() {
    var users =
        formKey.currentState?.fields["users"]?.value as List<api.UserView>? ??
            [];
    double totalUSD = 0.0;

    for (api.UserView user in users) {
      var fieldValue = formKey.currentState?.fields["${user.id}"]?.value;
      String portionAmountDisplay = fieldValue?.toString() ?? "0";
      
      // Convert display currency amount to USD for calculation
      double portionValueUSD = exchangeCustomToUSD(portionAmountDisplay).toDouble();
      totalUSD += portionValueUSD;
    }

    return totalUSD;
  }

  List<Widget> buildSplitWithPortionsField() {
    List<Widget> fields = [];
    if (quickActionsSelection[1]) {
      var users =
          formKey.currentState!.fields["users"]!.value as List<api.UserView>;

      users.forEach((user) {
        fields.add(const SizedBox(
          height: 10,
          width: 10,
        ));
        fields.add(AmountField(
            label: "${user.displayName}'s Portion",
            fieldName: "${user.id}",
            initialAmount: "0",
            formState: formState));
      });

      // Add total display for portions
      if (users.isNotEmpty) {
        fields.add(const SizedBox(height: 10));

        // Get receipt amount and convert to USD for calculation
        String receiptAmountDisplay = getReceiptAmount();
        double receiptValueUSD = exchangeCustomToUSD(receiptAmountDisplay).toDouble();
        
        // Calculate total portions in USD
        double totalPortionsUSD = calculateTotalPortionsInUSD();
        double remainingAmountUSD = receiptValueUSD - totalPortionsUSD;
        bool isValid = remainingAmountUSD >= 0;

        // Convert amounts back to display currency for showing to user
        String totalPortionsDisplay = formatCurrencyWithNegatives(context, totalPortionsUSD);
        String remainingAmountDisplay = formatCurrencyWithNegatives(context, remainingAmountUSD);

        fields.add(TotalDisplayWidget(
          value: "Portions: $totalPortionsDisplay, Remaining: $remainingAmountDisplay",
          isValid: isValid,
          errorMessage: !isValid ? "Portions exceed receipt total" : null,
        ));
      } else {
        fields.add(const SizedBox(height: 10));
        fields.add(Text(
          "Please select users to set custom portion amounts. The remaining balance will be split evenly.",
          style: TextStyle(color: Colors.grey[600]),
        ));
      }
    }

    return fields;
  }

  Widget buildPercentageButtonsForUser(api.UserView user) {
    String fieldName = "percentage_${user.id}";

    return Column(
      crossAxisAlignment: CrossAxisAlignment.start,
      children: [
        Text(
          "${user.displayName}'s Percentage:",
          style: TextStyle(fontWeight: FontWeight.w500),
        ),
        const SizedBox(height: 8),
        Wrap(
          spacing: 8,
          children: [
            ...percentageOptions.asMap().entries.map((entry) {
              int index = entry.key;
              String option = entry.value;

              return FilterChip(
                label: Text(option),
                selected:
                    formKey.currentState?.fields[fieldName]?.value == option,
                onSelected: (selected) {
                  if (selected) {
                    setState(() {
                      formKey.currentState!.fields[fieldName]!.setValue(option);
                      // Clear custom percentage field if not custom
                      if (option != "Custom") {
                        formKey.currentState!
                            .fields["custom_percentage_${user.id}"]
                            ?.setValue("");
                      }
                    });
                  }
                },
              );
            }).toList(),
          ],
        ),
      ],
    );
  }

  Widget buildCustomPercentageFieldForUser(api.UserView user) {
    String percentageFieldName = "percentage_${user.id}";
    String customFieldName = "custom_percentage_${user.id}";

    // Only show if "Custom" is selected for this user
    bool showCustomField =
        formKey.currentState?.fields[percentageFieldName]?.value == "Custom";

    if (showCustomField) {
      return Column(
        children: [
          const SizedBox(height: 8),
          FormBuilderTextField(
            name: customFieldName,
            decoration: InputDecoration(
              labelText: "${user.displayName}'s Custom Percentage",
              suffixText: "%",
              border: OutlineInputBorder(),
            ),
            keyboardType: TextInputType.number,
            validator: (value) {
              if (value == null || value.isEmpty) {
                return "Please enter a percentage";
              }
              final percentage = double.tryParse(value);
              if (percentage == null) {
                return "Please enter a valid number";
              }
              if (percentage < 0 || percentage > 100) {
                return "Percentage must be between 0 and 100";
              }
              return null;
            },
            onChanged: (value) {
              // Trigger rebuild to update total percentage display
              setState(() {});
            },
            initialValue: "0",
          ),
        ],
      );
    }
    return SizedBox.shrink();
  }

  List<Widget> buildPercentageFields() {
    List<Widget> fields = [];
    if (quickActionsSelection[2]) {
      var users =
          formKey.currentState?.fields["users"]?.value as List<api.UserView>? ??
              [];

      if (users.isNotEmpty) {
        fields.add(const SizedBox(height: 10));
        fields.add(Text(
          "Select percentage for each user:",
          style: TextStyle(fontSize: 16, fontWeight: FontWeight.w600),
        ));

        // Add total percentage display and validation
        fields.add(const SizedBox(height: 8));
        double totalPercentage = calculateTotalPercentage();
        bool isValid = isTotalPercentageValid();

        fields.add(TotalDisplayWidget(
          value: "${totalPercentage.toStringAsFixed(1)}%",
          isValid: isValid,
          errorMessage: !isValid ? "Exceeds 100%" : null,
        ));
        fields.add(const SizedBox(height: 10));

        for (api.UserView user in users) {
          // Add hidden form fields for this user's percentage selection
          fields.add(FormBuilderField<String>(
            name: "percentage_${user.id}",
            initialValue: "",
            builder: (FormFieldState<String> field) {
              return SizedBox.shrink();
            },
          ));

          fields.add(buildPercentageButtonsForUser(user));
          fields.add(buildCustomPercentageFieldForUser(user));
          fields.add(const SizedBox(height: 16));
        }
      } else {
        fields.add(const SizedBox(height: 10));
        fields.add(Text(
          "Please select users first to set their percentages.",
          style: TextStyle(color: Colors.grey[600]),
        ));
      }
    }
    return fields;
  }

  @override
  Widget build(BuildContext context) {
    return FormBuilder(
        key: formKey,
        onChanged: () {
          // Trigger rebuild when any form field changes
          if (quickActionsSelection[1]) {
            // Only rebuild for portions to avoid unnecessary updates
            Future.microtask(() {
              if (mounted) {
                setState(() {});
              }
            });
          }
        },
        child: Column(
          children: [
            FormBuilderField<List<bool>>(
              name: "quickActionsSelection",
              initialValue: quickActionsSelection,
              builder: (FormFieldState<List<bool>> field) {
                return SizedBox.shrink();
              },
            ),
            SingleChildScrollView(
              scrollDirection: Axis.horizontal,
              child: Row(
                crossAxisAlignment: CrossAxisAlignment.center,
                mainAxisAlignment: MainAxisAlignment.center,
                children: [buildToggleButtons()],
              ),
            ),
            textFieldSpacing,
            buildUserField(),
            ...buildSplitEvenlyTotal(),
            ...buildSplitWithPortionsField(),
            ...buildPercentageFields()
          ],
        ));
  }
}
