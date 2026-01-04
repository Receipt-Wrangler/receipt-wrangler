import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';

class SystemSettingsModel extends ChangeNotifier {
  String currencyDisplay = '';

  String get getCurrencyDisplay => currencyDisplay;

  CurrencySeparator currencyDecimalSeparator = CurrencySeparator.period;

  CurrencySeparator get getCurrencyDecimalSeparator => currencyDecimalSeparator;

  CurrencySeparator currencyThousandSeparator = CurrencySeparator.comma;

  CurrencySeparator get getCurrencyThousandSeparator =>
      currencyThousandSeparator;

  CurrencySymbolPosition currencySymbolPosition = CurrencySymbolPosition.END;

  CurrencySymbolPosition get getCurrencySymbolPosition =>
      currencySymbolPosition;

  bool currencyHideDecimalPlaces = false;

  bool get getCurrencyHideDecimalPlaces => currencyHideDecimalPlaces;

  void setCurrencyDisplay(String currencyDisplay) {
    this.currencyDisplay = currencyDisplay;

    notifyListeners();
  }

  void setCurrencyDecimalSeparator(CurrencySeparator currencyDecimalSeparator) {
    this.currencyDecimalSeparator = currencyDecimalSeparator;

    notifyListeners();
  }

  void setCurrencyThousandSeparator(
      CurrencySeparator currencyThousandSeparator) {
    this.currencyThousandSeparator = currencyThousandSeparator;

    notifyListeners();
  }

  void setCurrencySymbolPosition(
      CurrencySymbolPosition currencySymbolPosition) {
    this.currencySymbolPosition = currencySymbolPosition;

    notifyListeners();
  }

  void setCurrencyHideDecimalPlaces(bool currencyHideDecimalPlaces) {
    this.currencyHideDecimalPlaces = currencyHideDecimalPlaces;

    notifyListeners();
  }
}
