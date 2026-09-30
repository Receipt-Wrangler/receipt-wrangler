import 'dart:typed_data';

import 'package:built_collection/built_collection.dart';
import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/interfaces/upload_multipart_file_data.dart';
import 'package:receipt_wrangler_mobile/shared/functions/receipt_upload.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';

// Creating a receipt is ONE multipart call carrying the command (comments
// included) and every staged image, so the server can check role-required
// images at create time and the create is atomic. These pin that shape: a
// single createReceiptWithFiles, no follow-up uploadReceiptImage, and files a
// retry can send again.

class _MockOpenapi extends Mock implements api.Openapi {}

class _MockReceiptApi extends Mock implements api.ReceiptApi {}

class _MockReceiptImageApi extends Mock implements api.ReceiptImageApi {}

UploadMultipartFileData _image(String name, List<int> bytes) =>
    UploadMultipartFileData(
      multipartFile: MultipartFile.fromBytes(bytes, filename: name),
      bytes: Uint8List.fromList(bytes),
    );

api.UpsertReceiptCommand _command() =>
    (api.UpsertReceiptCommandBuilder()
          ..name = 'Coffee'
          ..amount = '4.50'
          ..date = '2026-09-30T00:00:00Z'
          ..groupId = 7
          ..paidByUserId = 1
          ..status = api.ReceiptStatus.OPEN
          ..comments = ListBuilder<api.UpsertCommentCommand>([
            (api.UpsertCommentCommandBuilder()
                  ..comment = 'for the team'
                  ..receiptId = 0
                  ..userId = 1)
                .build(),
          ]))
        .build();

void main() {
  late _MockReceiptApi receiptApi;
  late _MockReceiptImageApi receiptImageApi;

  setUpAll(() {
    registerFallbackValue(_command());
  });

  setUp(() {
    final client = _MockOpenapi();
    receiptApi = _MockReceiptApi();
    receiptImageApi = _MockReceiptImageApi();
    when(() => client.getReceiptApi()).thenReturn(receiptApi);
    when(() => client.getReceiptImageApi()).thenReturn(receiptImageApi);
    OpenApiClient.client = client;
  });

  void stubCreate() {
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
  }

  BuiltList<MultipartFile>? capturedFiles() =>
      verify(
            () => receiptApi.createReceiptWithFiles(
              receipt: any(named: 'receipt'),
              files: captureAny(named: 'files'),
            ),
          ).captured.single
          as BuiltList<MultipartFile>?;

  test('sends the command and every staged image in one call', () async {
    stubCreate();
    final command = _command();

    final receipt = await createReceiptWithImages(command, [
      _image('a.png', [1, 2]),
      _image('b.pdf', [3]),
    ]);

    expect(receipt.id, 42);
    final captured =
        verify(
          () => receiptApi.createReceiptWithFiles(
            receipt: captureAny(named: 'receipt'),
            files: captureAny(named: 'files'),
          ),
        ).captured;
    expect(captured, hasLength(2), reason: 'exactly one create call');
    expect(captured[0], command);
    final files = captured[1] as BuiltList<MultipartFile>;
    expect(files.map((f) => f.filename), ['a.png', 'b.pdf']);
    expect(files.map((f) => f.length), [2, 1]);
    verifyZeroInteractions(receiptImageApi);
  });

  test('sends no files part when nothing is staged', () async {
    stubCreate();

    await createReceiptWithImages(_command(), []);

    expect(capturedFiles(), isNull);
    verifyZeroInteractions(receiptImageApi);
  });

  test('rebuilds each file, so a failed create can be retried', () async {
    // A dio MultipartFile can be finalized once; resending the staged instance
    // after a failure would throw before the retry left the device.
    final staged = _image('a.png', [1]);
    when(
      () => receiptApi.createReceiptWithFiles(
        receipt: any(named: 'receipt'),
        files: any(named: 'files'),
      ),
    ).thenThrow(DioException(requestOptions: RequestOptions(path: '')));

    await expectLater(
      createReceiptWithImages(_command(), [staged]),
      throwsA(isA<DioException>()),
    );
    final first = capturedFiles()!.single;

    stubCreate();
    await createReceiptWithImages(_command(), [staged]);
    final second = capturedFiles()!.single;

    expect(identical(first, staged.multipartFile), isFalse);
    expect(identical(second, first), isFalse);
    expect(second.filename, 'a.png');
  });
}
