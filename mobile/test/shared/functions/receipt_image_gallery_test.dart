import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:flutter/services.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' show Openapi, ReceiptImageApi;
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/shared/functions/receipt_image_gallery.dart';

class MockOpenapi extends Mock implements Openapi {}

class MockReceiptImageApi extends Mock implements ReceiptImageApi {}

void main() {
  TestWidgetsFlutterBinding.ensureInitialized();

  const galChannel = MethodChannel('gal');
  final messenger =
      TestDefaultBinaryMessengerBinding.instance.defaultBinaryMessenger;

  late MockOpenapi mockClient;
  late MockReceiptImageApi mockReceiptImageApi;
  late LoadingModel loadingModel;
  late bool grantAccess;
  late List<String> galCalls;

  setUp(() {
    mockClient = MockOpenapi();
    mockReceiptImageApi = MockReceiptImageApi();
    when(() => mockClient.getReceiptImageApi()).thenReturn(mockReceiptImageApi);
    // The default client builds a real Openapi with a relative baseUrl that dio
    // rejects off-Web, so install the mock for the test's lifetime (no restore).
    OpenApiClient.client = mockClient;

    loadingModel = LoadingModel();
    grantAccess = true;
    galCalls = [];
    messenger.setMockMethodCallHandler(galChannel, (call) async {
      galCalls.add(call.method);
      if (call.method == 'requestAccess') return grantAccess;
      if (call.method == 'hasAccess') return true;
      return null; // putImageBytes, open
    });
  });

  tearDown(() {
    messenger.setMockMethodCallHandler(galChannel, null);
  });

  void stubDownload({Uint8List? data, bool throwError = false}) {
    final call =
        when(() => mockReceiptImageApi.downloadReceiptImageById(receiptImageId: 5));
    if (throwError) {
      call.thenThrow(DioException(requestOptions: RequestOptions(path: '')));
    } else {
      call.thenAnswer((_) async => Response<Uint8List>(
            requestOptions: RequestOptions(path: ''),
            statusCode: 200,
            data: data,
          ));
    }
  }

  Widget host() => MaterialApp(
        home: Scaffold(
          body: Builder(
            builder: (ctx) => ElevatedButton(
              onPressed: () =>
                  saveReceiptImageToGallery(ctx, loadingModel, 5),
              child: const Text('go'),
            ),
          ),
        ),
      );

  Future<void> tapAndDrain(WidgetTester tester) async {
    await tester.tap(find.text('go'));
    await tester.pump(); // run the async body
    await tester.pump(const Duration(milliseconds: 750)); // snackbar entrance
  }

  testWidgets('saves the image and shows a success snackbar', (tester) async {
    stubDownload(data: Uint8List.fromList([1, 2, 3]));
    await tester.pumpWidget(host());

    await tapAndDrain(tester);

    expect(galCalls, contains('putImageBytes'));
    expect(find.text('Image saved to gallery'), findsOneWidget);
    expect(loadingModel.isLoading, isFalse);
  });

  testWidgets('bails with a snackbar when gallery access is denied',
      (tester) async {
    stubDownload(data: Uint8List.fromList([1, 2, 3]));
    grantAccess = false;
    await tester.pumpWidget(host());

    await tapAndDrain(tester);

    expect(find.text('Gallery access denied'), findsOneWidget);
    expect(galCalls.contains('putImageBytes'), isFalse);
    expect(loadingModel.isLoading, isFalse);
  });

  testWidgets('returns quietly when the image bytes are null', (tester) async {
    stubDownload(data: null);
    await tester.pumpWidget(host());

    await tapAndDrain(tester);

    expect(find.byType(SnackBar), findsNothing);
    expect(galCalls, isEmpty); // requestAccess is only reached after the null check
    expect(loadingModel.isLoading, isFalse);
  });

  testWidgets('shows an error snackbar when the download fails', (tester) async {
    stubDownload(throwError: true);
    await tester.pumpWidget(host());

    await tapAndDrain(tester);

    expect(find.text('Failed to save image to gallery'), findsOneWidget);
    expect(loadingModel.isLoading, isFalse);
  });
}
