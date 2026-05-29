import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:receipt_wrangler_mobile/utils/snackbar.dart';

DioException _exceptionWith({Response? response}) => DioException(
      requestOptions: RequestOptions(path: '/test'),
      response: response,
    );

Response _responseWith(Object? data) => Response(
      requestOptions: RequestOptions(path: '/test'),
      data: data,
    );

Future<void> _pumpAndFire(
  WidgetTester tester,
  DioException error,
) async {
  await tester.pumpWidget(MaterialApp(
    home: Scaffold(
      body: Builder(
        builder: (context) => ElevatedButton(
          onPressed: () => showApiErrorSnackbar(context, error),
          child: const Text('fire'),
        ),
      ),
    ),
  ));
  await tester.tap(find.text('fire'));
  await tester.pump();
}

void main() {
  testWidgets('null response → fallback message', (tester) async {
    await _pumpAndFire(tester, _exceptionWith(response: null));
    expect(find.text('An error occurred'), findsOneWidget);
  });

  testWidgets('response with data == null → fallback message', (tester) async {
    await _pumpAndFire(tester, _exceptionWith(response: _responseWith(null)));
    expect(find.text('An error occurred'), findsOneWidget);
  });

  testWidgets('non-Map data (string body) → fallback message', (tester) async {
    await _pumpAndFire(tester, _exceptionWith(response: _responseWith('oops')));
    expect(find.text('An error occurred'), findsOneWidget);
  });

  testWidgets('Map data without errorMsg key → fallback message',
      (tester) async {
    await _pumpAndFire(
      tester,
      _exceptionWith(response: _responseWith({'foo': 'bar'})),
    );
    expect(find.text('An error occurred'), findsOneWidget);
  });

  testWidgets('Map data with empty errorMsg → fallback message',
      (tester) async {
    await _pumpAndFire(
      tester,
      _exceptionWith(response: _responseWith({'errorMsg': ''})),
    );
    expect(find.text('An error occurred'), findsOneWidget);
  });

  testWidgets('Map data with valid errorMsg → shows backend message',
      (tester) async {
    await _pumpAndFire(
      tester,
      _exceptionWith(response: _responseWith({'errorMsg': 'Bad name'})),
    );
    expect(find.text('Bad name'), findsOneWidget);
    expect(find.text('An error occurred'), findsNothing);
  });
}
