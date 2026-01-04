import 'package:currency_textfield/currency_textfield.dart';
import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:form_builder_validators/form_builder_validators.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/utils/currency.dart';

import '../../constants/currency.dart';
import '../../models/system_settings_model.dart';
import '../../utils/forms.dart';

class AmountField extends StatefulWidget {
  const AmountField({
    super.key,
    required this.label,
    required this.fieldName,
    required this.initialAmount,
    required this.formState,
    this.suffixIcon,
  });

  final String label;

  final String fieldName;

  final String initialAmount;

  final WranglerFormState formState;

  final Widget? suffixIcon;

  @override
  State<AmountField> createState() => _AmountField();
}

class _AmountField extends State<AmountField> {
  late final systemSettingsModel = Provider.of<SystemSettingsModel>(context);
  late final controller = CurrencyTextFieldController(
      decimalSymbol: systemSettingsModel.currencyHideDecimalPlaces
          ? ""
          : getCurrencySeparatorLiteral(
              systemSettingsModel.currencyDecimalSeparator),
      thousandSymbol: getCurrencySeparatorLiteral(
          systemSettingsModel.currencyThousandSeparator),
      initDoubleValue: getInitialAmount(),
      numberOfDecimals: systemSettingsModel.currencyHideDecimalPlaces ? 0 : 2,
      currencySymbol: "",
      startWithSeparator: true,
      forceCursorToEnd: false);

  double getInitialAmount() {
    var doubleAmount = 0.0;

    try {
      doubleAmount = double.parse(widget.initialAmount);
    } catch (e) {
      var exchanged = exchangeCustomToUSD(widget.initialAmount);
      doubleAmount = double.parse(exchanged.format(numberFormatWithoutSymbol));
    }

    return doubleAmount;
  }

  @override
  Widget build(BuildContext context) {
    var systemSettingsModel = Provider.of<SystemSettingsModel>(context);

    var prefix = systemSettingsModel.currencySymbolPosition ==
            api.CurrencySymbolPosition.START
        ? systemSettingsModel.currencyDisplay
        : null;

    var suffix = systemSettingsModel.currencySymbolPosition ==
            api.CurrencySymbolPosition.END
        ? systemSettingsModel.currencyDisplay
        : null;

    return FormBuilderTextField(
      key: widget.key,
      name: widget.fieldName,
      decoration: InputDecoration(
        labelText: widget.label,
        prefixText: prefix,
        suffixText: suffix,
        suffixIcon: widget.suffixIcon,
      ),
      keyboardType: TextInputType.number,
      validator: FormBuilderValidators.compose([
        FormBuilderValidators.required(),
      ]),
      readOnly: isFieldReadOnly(widget.formState),
      controller: controller,
      valueTransformer: (value) {
        return exchangeCustomToUSD(value)
            .format(numberFormatWithoutSymbolOrGroupSeparator);
      },
    );
  }
}
