import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:flutter_svg/flutter_svg.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/context_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_form.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';

import '../helpers/receipt_form_test_helpers.dart';

class _MockOpenapi extends Mock implements api.Openapi {}

class _MockReceiptApi extends Mock implements api.ReceiptApi {}

const _groupId = 7;

Finder _textField(String name) =>
    find.byWidgetPredicate((w) => w is FormBuilderTextField && w.name == name);

Finder _dropdown(String name) =>
    find.byWidgetPredicate((w) => w is FormBuilderDropdown && w.name == name);

void main() {
  late _MockReceiptApi receiptApi;

  setUpAll(() {
    registerFallbackValue(
      (api.UpsertReceiptCommandBuilder()
            ..name = ''
            ..amount = '0'
            ..date = ''
            ..groupId = 0
            ..paidByUserId = 0
            ..status = api.ReceiptStatus.OPEN)
          .build(),
    );
  });

  setUp(() {
    final client = _MockOpenapi();
    receiptApi = _MockReceiptApi();
    when(() => client.getReceiptApi()).thenReturn(receiptApi);
    OpenApiClient.client = client;

    when(
      () => receiptApi.createReceiptWithFiles(
        receipt: any(named: 'receipt'),
        files: any(named: 'files'),
      ),
    ).thenAnswer(
      (_) async => Response<api.Receipt>(
        requestOptions: RequestOptions(path: ''),
        statusCode: 200,
        data: getDefaultReceipt().rebuild((b) => b..id = 42),
      ),
    );
  });

  final users = [
    buildUserView(id: 1, displayName: 'Alice'),
    buildUserView(id: 2, displayName: 'Bob'),
  ];

  Future<ReceiptFormHarness> pumpBlankAddForm(WidgetTester tester) async {
    tester.view.physicalSize = const Size(1080, 2400);
    tester.view.devicePixelRatio = 2.0;
    addTearDown(tester.view.reset);

    // What the server actually sends: the item pickers are shown.
    final settings = buildGroupReceiptSettings(groupId: _groupId).rebuild(
      (b) =>
          b
            ..hideItemCategories = false
            ..hideItemTags = false,
    );

    // A blank add form. The user's only group seeds the Group field.
    final harness = await pumpReceiptForm(
      tester,
      formState: WranglerFormState.add,
      groups: [
        buildGroup(
          id: _groupId,
          receiptSettings: settings,
          members: [
            buildGroupMember(userId: 1, groupId: _groupId),
            buildGroupMember(userId: 2, groupId: _groupId),
          ],
        ),
      ],
      users: users,
      realSubmitButton: true,
    );

    // ReceiptFormScreen sets this in the app; the sheet's Split button pops
    // through it.
    final formContext = tester.element(find.byType(ReceiptForm));
    Provider.of<ContextModel>(
      formContext,
      listen: false,
    ).setShellContext(formContext);

    return harness;
  }

  Future<void> fillReceipt(WidgetTester tester) async {
    await tester.enterText(_textField('name'), 'Dinner');
    await tester.enterText(_textField('amount'), '100.00');
    await tester.pumpAndSettle();

    await tester.ensureVisible(_dropdown('paidByUserId'));
    await tester.pumpAndSettle();
    await tester.tap(_dropdown('paidByUserId'));
    await tester.pumpAndSettle();
    await tester.tap(find.text('Alice').last);
    await tester.pumpAndSettle();
  }

  Future<void> splitEvenlyThroughTheSheet(WidgetTester tester) async {
    await tester.ensureVisible(find.byType(SvgPicture));
    await tester.pumpAndSettle();
    await tester.tap(find.byType(SvgPicture));
    await tester.pumpAndSettle();
    expect(find.text('Quick Actions'), findsOneWidget);

    // The real Users picker: a nested sheet of chips plus a Select button.
    await tester.tap(find.text('Users'), warnIfMissed: false);
    await tester.pumpAndSettle();
    expect(find.text('Select Users'), findsOneWidget);
    await tester.tap(find.text('Alice').last);
    await tester.pump();
    await tester.tap(find.text('Bob').last);
    await tester.pump();
    await tester.tap(find.text('Select'));
    await tester.pumpAndSettle();

    await tester.tap(find.text('Split'));
    await tester.pumpAndSettle();
    expect(
      find.text('Quick Actions'),
      findsNothing,
      reason: 'the sheet should close after Split',
    );
  }

  Future<void> submitAndExpectCreate(
    WidgetTester tester,
    ReceiptFormHarness harness,
  ) async {
    expect(harness.receiptModel.items.length, 2);

    await tester.tap(find.byType(BottomSubmitButton));
    await tester.pumpAndSettle();

    // Null once a successful submit has navigated off the form.
    final errors = harness.formKey.currentState?.errors ?? const {};
    expect(
      errors,
      isEmpty,
      reason: 'the submit was refused by field errors the user cannot see',
    );
    final captured =
        verify(
          () => receiptApi.createReceiptWithFiles(
            receipt: captureAny(named: 'receipt'),
            files: any(named: 'files'),
          ),
        ).captured;
    final command = captured.single as api.UpsertReceiptCommand;
    expect(command.receiptItems?.length, 2);
    expect(find.text('receipt 42 view'), findsOneWidget);
  }

  testWidgets('fill the form, then Split Evenly, then submit', (tester) async {
    final harness = await pumpBlankAddForm(tester);
    await fillReceipt(tester);
    await splitEvenlyThroughTheSheet(tester);
    await submitAndExpectCreate(tester, harness);
  });

  // Reproduces the bug: split while the receipt amount is still blank and every
  // share gets 0.0, which the currency field renders as an empty box. Each
  // share's required Amount validator then fails inside its collapsed "Shared
  // With" panel, and the submit returns silently -- "nothing happens".
  testWidgets(
    'Split Evenly first, then fill the form, then submit',
    skip: true, // Known bug, unskip with the fix.
    (tester) async {
      final harness = await pumpBlankAddForm(tester);
      await splitEvenlyThroughTheSheet(tester);
      await fillReceipt(tester);
      await submitAndExpectCreate(tester, harness);
    },
  );
}
