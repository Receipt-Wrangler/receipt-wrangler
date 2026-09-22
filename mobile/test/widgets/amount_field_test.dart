import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_test/flutter_test.dart';

import '../helpers/widget_test_helpers.dart';

void main() {
  setUpAll(() {
    registerCustomCurrencyForTests();
  });

  const amountFieldKey = ValueKey('amount-field');

  Finder findInnerTextField() => find.descendant(
        of: find.byKey(amountFieldKey),
        matching: find.byType(FormBuilderTextField),
      );

  testWidgets(
    'preserves a negative initial amount through valueTransformer',
    (tester) async {
      final formKey = await pumpAmountField(
        tester,
        amountFieldKey: amountFieldKey,
        initialAmount: '-50.00',
      );

      // Negative controller text is non-empty → required validator passes,
      // valueTransformer round-trips the sign through exchangeCustomToUSD.
      expect(formKey.currentState!.saveAndValidate(), isTrue);
      expect(formKey.currentState!.value['amount'], '-50.00');
    },
  );

  testWidgets(
    'preserves a positive initial amount (regression baseline)',
    (tester) async {
      final formKey = await pumpAmountField(
        tester,
        amountFieldKey: amountFieldKey,
        initialAmount: '50.00',
      );

      expect(formKey.currentState!.saveAndValidate(), isTrue);
      expect(formKey.currentState!.value['amount'], '50.00');
    },
  );

  testWidgets(
    'transforms a zero initial amount to "0.00" but the required validator '
    'rejects the empty controller text',
    (tester) async {
      // Documents pre-existing CurrencyTextFieldController behavior:
      // initDoubleValue=0.0 produces empty controller text. The required
      // validator runs against the empty text and fails, while
      // valueTransformer still produces "0.00" for the form state.
      final formKey = await pumpAmountField(
        tester,
        amountFieldKey: amountFieldKey,
        initialAmount: '0.00',
      );

      expect(formKey.currentState!.saveAndValidate(), isFalse);
      expect(formKey.currentState!.value['amount'], '0.00');
    },
  );

  testWidgets(
    'falls back to zero when initial amount is empty',
    (tester) async {
      final formKey = await pumpAmountField(
        tester,
        amountFieldKey: amountFieldKey,
        initialAmount: '',
      );

      // Same as the zero case: getInitialAmount returns 0 via the
      // exchangeCustomToUSD null/empty guard, and the validator rejects
      // the empty controller text.
      expect(formKey.currentState!.saveAndValidate(), isFalse);
      expect(formKey.currentState!.value['amount'], '0.00');
    },
  );

  testWidgets(
    'uses signed-decimal keyboardType so users can enter a minus sign',
    (tester) async {
      await pumpAmountField(
        tester,
        amountFieldKey: amountFieldKey,
        initialAmount: '-1.00',
      );

      final textField = tester.widget<FormBuilderTextField>(findInnerTextField());
      expect(
        textField.keyboardType,
        const TextInputType.numberWithOptions(signed: true, decimal: true),
      );
    },
  );

  testWidgets(
    'required validator accepts a small non-zero amount',
    (tester) async {
      // Confirms the validator + valueTransformer chain works for a
      // genuinely valid input (the smallest non-zero amount we can have).
      final formKey = await pumpAmountField(
        tester,
        amountFieldKey: amountFieldKey,
        initialAmount: '0.01',
      );

      expect(formKey.currentState!.saveAndValidate(), isTrue);
      expect(formKey.currentState!.value['amount'], '0.01');
    },
  );

  group('the optional validator', () {
    testWidgets(
      'is required by default, which every receipt call site relies on',
      (tester) async {
        // A zero initial amount renders an EMPTY box, so the required validator
        // fails -- that is the existing contract and it must not move.
        final formKey = await pumpAmountField(
          tester,
          amountFieldKey: amountFieldKey,
          initialAmount: '0.00',
        );

        expect(formKey.currentState!.saveAndValidate(), isFalse);
      },
    );

    testWidgets(
      'an explicit validator replaces the required one',
      (tester) async {
        // The receipt filter needs this: an amount condition is only authored
        // when the user asks for one, so a required field would block the sheet.
        final formKey = await pumpAmountField(
          tester,
          amountFieldKey: amountFieldKey,
          initialAmount: '0.00',
          validator: (_) => null,
        );

        expect(formKey.currentState!.saveAndValidate(), isTrue);
      },
    );

    testWidgets(
      'an explicit validator still gets to reject',
      (tester) async {
        final formKey = await pumpAmountField(
          tester,
          amountFieldKey: amountFieldKey,
          initialAmount: '12.34',
          validator: (_) => 'nope',
        );

        expect(formKey.currentState!.saveAndValidate(), isFalse);
        await tester.pump();
        expect(find.text('nope'), findsOneWidget);
      },
    );

    testWidgets(
      'the transformed value is unchanged by the validator override',
      (tester) async {
        final formKey = await pumpAmountField(
          tester,
          amountFieldKey: amountFieldKey,
          initialAmount: '42.50',
          validator: (_) => null,
        );

        expect(formKey.currentState!.saveAndValidate(), isTrue);
        expect(formKey.currentState!.value['amount'], '42.50');
      },
    );
  });
}
