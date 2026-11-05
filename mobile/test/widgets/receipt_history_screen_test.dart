import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/models/user_model.dart';
import 'package:receipt_wrangler_mobile/receipts/screens/receipt_history_screen.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_history_app_bar.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_history_entry_card.dart';
import 'package:receipt_wrangler_mobile/utils/receipt_history.dart';

import '../helpers/receipt_form_test_helpers.dart' show buildUserView;
import '../helpers/receipt_history_test_helpers.dart';
import '../helpers/widget_test_helpers.dart' show registerCustomCurrencyForTests;

class _MockOpenapi extends Mock implements api.Openapi {}

class _MockReceiptApi extends Mock implements api.ReceiptApi {}

const _receiptId = 5;

// Local times sent as UTC, so the day each entry lands on holds in any zone.
String _at(DateTime local) => local.toUtc().toIso8601String();

late _MockReceiptApi _receiptApi;

void _answerWith(Map<String, Object?> json) {
  when(() => _receiptApi.getReceiptHistory(receiptId: any(named: 'receiptId')))
      .thenAnswer((_) async => Response<api.ReceiptHistory>(
            requestOptions: RequestOptions(path: ''),
            statusCode: 200,
            data: historyFromJson(json),
          ));
}

Future<void> _pumpScreen(WidgetTester tester) async {
  await tester.pumpWidget(MultiProvider(
    providers: [
      ChangeNotifierProvider<UserModel>(
        create: (_) => UserModel()
          ..setUsers([buildUserView(id: 2, displayName: 'Maya Chen')]),
      ),
    ],
    child: const MaterialApp(home: ReceiptHistoryScreen(receiptId: _receiptId)),
  ));
  await tester.pumpAndSettle();
}

Finder _card(int systemTaskId) =>
    find.byKey(ValueKey('receipt-history-entry-$systemTaskId'));

Finder _inCard(int systemTaskId, Finder finder) =>
    find.descendant(of: _card(systemTaskId), matching: finder);

void main() {
  setUpAll(registerCustomCurrencyForTests);

  setUp(() {
    // Never read the real OpenApiClient.client: it builds a Dio with a relative
    // baseUrl, which dio rejects off-Web.
    final client = _MockOpenapi();
    _receiptApi = _MockReceiptApi();
    when(() => client.getReceiptApi()).thenReturn(_receiptApi);
    OpenApiClient.client = client;
  });

  testWidgets('loads the receipt history and summarizes it in the app bar', (tester) async {
    _answerWith(historyJson([
      entryJson(
        systemTaskId: 3,
        action: 'UPDATED',
        recordedAt: _at(DateTime(2026, 10, 13, 16)),
        changes: [textChangeJson('NAME', from: 'COSTCO WHSE', to: 'Costco Wholesale')],
      ),
      entryJson(
        systemTaskId: 1,
        action: 'CREATED',
        recordedAt: _at(DateTime(2026, 10, 12, 9)),
        actor: {'type': 'EMAIL'},
        changes: [textChangeJson('EMAIL_INBOX', to: 'test@test.com')],
      ),
    ], changeCount: 1, saveCount: 2, contributorCount: 2));

    await _pumpScreen(tester);

    verify(() => _receiptApi.getReceiptHistory(receiptId: _receiptId)).called(1);
    expect(find.text('Receipt History'), findsOneWidget);
    expect(
      tester.widget<Text>(find.byKey(ReceiptHistoryAppBar.summaryKey)).data,
      '1 change · 2 saves · 2 contributors',
    );
    expect(
      find.byKey(ReceiptHistoryScreen.dayKey).evaluate().map((e) => (e.widget as Text).data),
      ['TUESDAY, OCT 13', 'MONDAY, OCT 12'],
    );

    expect(_inCard(3, find.textContaining('Maya Chen', findRichText: true)), findsOneWidget);
    expect(_inCard(3, find.textContaining('edited this receipt', findRichText: true)), findsOneWidget);
    expect(_inCard(3, find.text('COSTCO WHSE')), findsOneWidget);
    expect(_inCard(3, find.text('Costco Wholesale')), findsOneWidget);

    expect(_inCard(1, find.textContaining('Email Upload', findRichText: true)), findsOneWidget);
    expect(_inCard(1, find.text('Inbox')), findsOneWidget);
    expect(_inCard(1, find.text('test@test.com')), findsOneWidget);
  });

  testWidgets('collapses an entry after three changes, and expands it again', (tester) async {
    _answerWith(historyJson([
      entryJson(systemTaskId: 7, changes: [
        for (var i = 0; i < 5; i++) textChangeJson('TAGS', to: 'Tag $i'),
      ]),
    ]));

    await _pumpScreen(tester);

    expect(_inCard(7, find.textContaining('Tag ')), findsNWidgets(historyCollapseAfter));
    final showMore = _inCard(7, find.byKey(ReceiptHistoryEntryCard.showMoreKey));
    expect(_inCard(7, find.text('Show 2 more')), findsOneWidget);

    await tester.tap(showMore);
    await tester.pumpAndSettle();
    expect(_inCard(7, find.textContaining('Tag ')), findsNWidgets(5));
    expect(_inCard(7, find.text('Show less')), findsOneWidget);

    await tester.tap(showMore);
    await tester.pumpAndSettle();
    expect(_inCard(7, find.textContaining('Tag ')), findsNWidgets(historyCollapseAfter));
  });

  testWidgets('shows a failed save, and an older entry\'s notice', (tester) async {
    _answerWith(historyJson([
      entryJson(
        systemTaskId: 9,
        action: 'FAILED',
        hasSnapshots: false,
        errorMessage: 'Your role requires at least one image.',
      ),
      entryJson(systemTaskId: 8, beforeState: 'INCOMPLETE'),
    ]));

    await _pumpScreen(tester);

    expect(_inCard(9, find.byKey(ReceiptHistoryEntryCard.failedKey)), findsOneWidget);
    expect(
      _inCard(9, find.text('Save failed: Your role requires at least one image.')),
      findsOneWidget,
    );
    expect(_inCard(9, find.textContaining('tried to save changes', findRichText: true)), findsOneWidget);
    expect(_inCard(8, find.byKey(ReceiptHistoryEntryCard.noticeKey)), findsOneWidget);
    expect(_inCard(8, find.byKey(ReceiptHistoryEntryCard.failedKey)), findsNothing);
  });

  testWidgets('names a user it cannot resolve as unknown rather than failing', (tester) async {
    _answerWith(historyJson([
      entryJson(systemTaskId: 4, actor: {'type': 'USER', 'userId': 99}),
    ]));

    await _pumpScreen(tester);

    expect(_inCard(4, find.textContaining(unknownUserName, findRichText: true)), findsOneWidget);
    expect(tester.takeException(), isNull);
  });

  testWidgets('pulling down loads the history again', (tester) async {
    _answerWith(historyJson([entryJson(systemTaskId: 2)], saveCount: 1, contributorCount: 1));
    await _pumpScreen(tester);

    await tester.fling(_card(2), const Offset(0, 400), 1000);
    await tester.pumpAndSettle();

    verify(() => _receiptApi.getReceiptHistory(receiptId: _receiptId)).called(2);
    expect(tester.takeException(), isNull);
    expect(_card(2), findsOneWidget);
  });

  testWidgets('says so when there is no history yet', (tester) async {
    _answerWith(historyJson([]));

    await _pumpScreen(tester);

    expect(find.byKey(ReceiptHistoryScreen.emptyKey), findsOneWidget);
  });

  testWidgets('reports a failed load and retries it', (tester) async {
    when(() => _receiptApi.getReceiptHistory(receiptId: any(named: 'receiptId')))
        .thenAnswer((_) async => throw Exception('offline'));

    await _pumpScreen(tester);
    expect(find.byKey(ReceiptHistoryScreen.errorKey), findsOneWidget);
    expect(find.byKey(ReceiptHistoryAppBar.summaryKey), findsNothing);

    _answerWith(historyJson([entryJson(systemTaskId: 2)], saveCount: 1, contributorCount: 1));
    await tester.tap(find.byKey(ReceiptHistoryScreen.retryKey));
    await tester.pumpAndSettle();

    expect(find.byKey(ReceiptHistoryScreen.errorKey), findsNothing);
    expect(_card(2), findsOneWidget);
  });
}
