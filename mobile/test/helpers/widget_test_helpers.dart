import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:money2/money2.dart';
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/constants/currency.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/system_settings_model.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/amount_field.dart';

/// Registers the custom currency that `exchangeCustomToUSD` /
/// `exchangeUSDToCustom` rely on. Safe to call multiple times — money2's
/// Currencies registry is a process-wide singleton.
void registerCustomCurrencyForTests() {
  if (Currencies().find(customCurrencyISOCode) != null) {
    return;
  }
  Currencies().register(Currency.create(
    customCurrencyISOCode,
    2,
    name: customCurrencyISOCode,
    symbol: '\$',
    groupSeparator: ',',
    decimalSeparator: '.',
    pattern: '###,###.00S',
  ));
}

/// Pumps an [AmountField] inside the minimal widget tree it needs:
/// [MaterialApp] for theming, [ChangeNotifierProvider] for the
/// [SystemSettingsModel], and a [FormBuilder] parent so the field can register
/// itself. Returns the form key so tests can call `saveAndValidate()` and
/// inspect `currentState!.value`.
Future<GlobalKey<FormBuilderState>> pumpAmountField(
  WidgetTester tester, {
  required Key amountFieldKey,
  String initialAmount = '0.00',
  WranglerFormState formState = WranglerFormState.add,
  String? Function(String?)? validator,
}) async {
  final formKey = GlobalKey<FormBuilderState>();

  await tester.pumpWidget(
    ChangeNotifierProvider<SystemSettingsModel>(
      create: (_) => SystemSettingsModel(),
      child: MaterialApp(
        home: Scaffold(
          body: FormBuilder(
            key: formKey,
            child: AmountField(
              key: amountFieldKey,
              label: 'Amount',
              fieldName: 'amount',
              initialAmount: initialAmount,
              formState: formState,
              validator: validator,
            ),
          ),
        ),
      ),
    ),
  );
  await tester.pump();

  return formKey;
}
