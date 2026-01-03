import 'package:flutter/material.dart';
import 'package:money2/money2.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/currency.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';

import '../../interfaces/form_item.dart';
import '../../models/receipt_model.dart';
import '../../utils/currency.dart';

class ReceiptQuickActionsSubmitButton extends StatefulWidget {
  const ReceiptQuickActionsSubmitButton({super.key});

  @override
  State<ReceiptQuickActionsSubmitButton> createState() =>
      _ReceiptQuickActionsSubmitButton();
}

class _ReceiptQuickActionsSubmitButton
    extends State<ReceiptQuickActionsSubmitButton> {
  late final formKey =
      Provider.of<ReceiptModel>(context, listen: false).quickActionsFormKey;
  late final userModel = Provider.of<UserModel>(context, listen: false);
  late final receiptModel = Provider.of<ReceiptModel>(context, listen: false);
  late final shellContext =
      Provider.of<ContextModel>(context, listen: false).shellContext;

  List<api.UserView> getSelectedUsers() {
    return formKey.currentState!.value["users"] as List<api.UserView>;
  }

  List<FormItem> buildEvenSplitFormItems(String amount) {
    var selectedUsers = getSelectedUsers();

    var receiptAmount = Money.parse(amount.length == 0 ? "0" : amount,
        isoCode: customCurrencyISOCode);
    var divisor = Money.parse(selectedUsers.length.toString(),
        isoCode: customCurrencyISOCode);

    double equalSplit = receiptAmount.dividedBy(divisor);
    var items = [...receiptModel.items];

    selectedUsers.forEach((user) {
      var newItem = (api.ItemBuilder()
            ..name = "${user.displayName}'s Even Portion"
            ..amount = equalSplit.toString()
            ..chargedToUserId = user.id
            ..receiptId = receiptModel.receipt.id
            ..status = api.ItemStatus.OPEN)
          .build();

      items.add(FormItem.fromItem(newItem));
    });

    return items;
  }

  // TODO: don't forget about validation
  void splitEvenly() {
    if (formKey.currentState!.saveAndValidate()) {
      var formAmount =
          receiptModel.receiptFormKey.currentState!.fields["amount"]!.value;
      List<FormItem> items = buildEvenSplitFormItems(formAmount);
      receiptModel.setItems(items);
      Navigator.pop(shellContext as BuildContext);
    }
  }

  double calculateTotalPortionsInUSD() {
    var selectedUsers = getSelectedUsers();
    double totalUSD = 0.0;
    
    for (api.UserView user in selectedUsers) {
      var fieldValue = formKey.currentState?.fields["${user.id}"]?.value;
      String portionAmountDisplay = fieldValue?.toString() ?? "0";
      
      // Convert display currency amount to USD for calculation
      double portionValueUSD = exchangeCustomToUSD(portionAmountDisplay).toDouble();
      totalUSD += portionValueUSD;
    }
    
    return totalUSD;
  }

  void splitEvenlyWithPortions() {
    if (formKey.currentState!.saveAndValidate()) {
      // Additional validation for portions not exceeding receipt total
      String receiptAmountDisplay = receiptModel.receiptFormKey.currentState!.fields["amount"]!.value ?? "0";
      double receiptTotalUSD = exchangeCustomToUSD(receiptAmountDisplay).toDouble();
      double totalPortionsUSD = calculateTotalPortionsInUSD();
      
      if (totalPortionsUSD > receiptTotalUSD) {
        // Convert amounts back to display currency for error message
        String totalPortionsDisplay = formatCurrency(shellContext as BuildContext, totalPortionsUSD.toString()) ?? "\$0.00";
        String receiptTotalDisplay = formatCurrency(shellContext as BuildContext, receiptTotalUSD.toString()) ?? "\$0.00";
        
        ScaffoldMessenger.of(shellContext as BuildContext).showSnackBar(
          SnackBar(
            content: Text(
              "Total portions ($totalPortionsDisplay) cannot exceed receipt amount ($receiptTotalDisplay)",
            ),
            backgroundColor: Colors.red,
            duration: Duration(seconds: 3),
          ),
        );
        return;
      }
      List<FormItem> userPortionsItems = [];
      var selectedUsers = getSelectedUsers();
      var remainingAmount = Money.parse(
          receiptModel.receiptFormKey.currentState!.fields["amount"]!.value,
          isoCode: customCurrencyISOCode);

      selectedUsers.forEach((user) {
        var splitAmount =
            formKey.currentState!.fields["${user.id}"]!.value.toString();

        if (splitAmount.length > 0) {
          var splitCurrency =
              Money.parse(splitAmount, isoCode: customCurrencyISOCode);
          remainingAmount = remainingAmount - splitCurrency;

          var newItem = (api.ItemBuilder()
                ..name = "${user.displayName}'s Portion"
                ..amount = splitCurrency.toDouble().toString()
                ..chargedToUserId = user.id
                ..receiptId = receiptModel.receipt.id
                ..status = api.ItemStatus.OPEN)
              .build();

          userPortionsItems.add(FormItem.fromItem(newItem));
        }
      });

      receiptModel.setItems([
        ...buildEvenSplitFormItems(remainingAmount.toDouble().toString()),
        ...userPortionsItems
      ]);
      Navigator.pop(shellContext as BuildContext);
    }
  }

  double calculateTotalPercentage() {
    var selectedUsers = getSelectedUsers();
    double total = 0.0;
    
    for (api.UserView user in selectedUsers) {
      var userPercentageSelection = formKey.currentState!.fields["percentage_${user.id}"]?.value as String?;
      
      if (userPercentageSelection != null && userPercentageSelection.isNotEmpty) {
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
            var customPercentageStr = formKey.currentState!.fields["custom_percentage_${user.id}"]?.value?.toString();
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

  void splitByPercentage() {
    if (formKey.currentState!.saveAndValidate()) {
      // Check if total percentage exceeds 100%
      double totalPercentage = calculateTotalPercentage();
      if (totalPercentage > 100.0) {
        ScaffoldMessenger.of(shellContext as BuildContext).showSnackBar(
          SnackBar(
            content: Text(
              "Total percentage (${totalPercentage.toStringAsFixed(1)}%) cannot exceed 100%",
            ),
            backgroundColor: Colors.red,
            duration: Duration(seconds: 3),
          ),
        );
        return;
      }
      
      var selectedUsers = getSelectedUsers();
      
      // Check if all users have percentage selections
      for (api.UserView user in selectedUsers) {
        var userPercentageSelection = formKey.currentState!.fields["percentage_${user.id}"]?.value as String?;
        if (userPercentageSelection == null || userPercentageSelection.isEmpty) {
          ScaffoldMessenger.of(shellContext as BuildContext).showSnackBar(
            SnackBar(
              content: Text(
                "Please select a percentage for ${user.displayName}",
              ),
              backgroundColor: Colors.orange,
              duration: Duration(seconds: 3),
            ),
          );
          return;
        }
      }
      var receiptAmount = Money.parse(
          receiptModel.receiptFormKey.currentState!.fields["amount"]!.value.isNotEmpty 
              ? receiptModel.receiptFormKey.currentState!.fields["amount"]!.value 
              : "0",
          isoCode: customCurrencyISOCode);
      
      var items = [...receiptModel.items];
      
      selectedUsers.forEach((user) {
        // Get the percentage selection for this specific user
        var userPercentageSelection = formKey.currentState!.fields["percentage_${user.id}"]?.value as String?;
        
        if (userPercentageSelection == null || userPercentageSelection.isEmpty) {
          // Skip users without percentage selection
          return;
        }
        
        double percentage = 0.0;
        String percentageLabel = "";
        
        switch (userPercentageSelection) {
          case "25%":
            percentage = 0.25;
            percentageLabel = "25%";
            break;
          case "50%":
            percentage = 0.50;
            percentageLabel = "50%";
            break;
          case "75%":
            percentage = 0.75;
            percentageLabel = "75%";
            break;
          case "100%":
            percentage = 1.0;
            percentageLabel = "100%";
            break;
          case "Custom":
            // Get custom percentage value for this user
            var customPercentageStr = formKey.currentState!.fields["custom_percentage_${user.id}"]?.value?.toString();
            if (customPercentageStr != null && customPercentageStr.isNotEmpty) {
              percentage = double.parse(customPercentageStr) / 100.0;
              percentageLabel = "${customPercentageStr}%";
            }
            break;
        }
        
        if (percentage > 0) {
          var userAmount = receiptAmount * percentage;
          
          var newItem = (api.ItemBuilder()
                ..name = "${user.displayName}'s $percentageLabel Portion"
                ..amount = userAmount.toDouble().toString()
                ..chargedToUserId = user.id
                ..receiptId = receiptModel.receipt.id
                ..status = api.ItemStatus.OPEN)
              .build();
          
          items.add(FormItem.fromItem(newItem));
        }
      });
      
      receiptModel.setItems(items);
      Navigator.pop(shellContext as BuildContext);
    }
  }

  void onSubmitPressed() {
    if (formKey.currentState!.saveAndValidate()) {
      var quickActionSelection =
          formKey.currentState!.value["quickActionsSelection"] as List<bool>;
      var selected = quickActionSelection.indexOf(true);
      switch (selected) {
        case 0:
          splitEvenly();
        case 1:
          splitEvenlyWithPortions();
        case 2:
          splitByPercentage();
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    return BottomSubmitButton(onPressed: onSubmitPressed, buttonText: "Split");
  }
}
