import 'dart:typed_data';

import 'package:built_collection/built_collection.dart';
import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:mocktail/mocktail.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/client/client.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/interfaces/upload_multipart_file_data.dart';
import 'package:receipt_wrangler_mobile/shared/widgets/bottom_submit_button.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';

import '../helpers/permission_test_helpers.dart';
import '../helpers/receipt_form_test_helpers.dart';

// The receipt form's submit refuses, client-side, a receipt missing what the
// caller's group role requires (AppData.groupReceiptRequirements), and an add
// that passes goes out as ONE createReceiptWithFiles call. Driven through the
// real ReceiptBottomSheetBuilder button on the real form.

class _MockOpenapi extends Mock implements api.Openapi {}

class _MockReceiptApi extends Mock implements api.ReceiptApi {}

class _MockReceiptImageApi extends Mock implements api.ReceiptImageApi {}

const _groupId = 7;
const _userId = 1;

api.Receipt _receipt({int id = 0, List<api.Comment> comments = const []}) =>
    getDefaultReceipt().rebuild(
      (b) =>
          b
            ..id = id
            ..name = 'Coffee'
            ..amount = '12.34'
            ..groupId = _groupId
            ..paidByUserId = _userId
            ..comments = ListBuilder<api.Comment>(comments),
    );

api.Comment _comment(String text) =>
    (api.CommentBuilder()
          ..id = 1
          ..comment = text
          ..receiptId = 5
          ..userId = _userId)
        .build();

UploadMultipartFileData _staged() => UploadMultipartFileData(
  multipartFile: MultipartFile.fromBytes(const [1], filename: 'a.png'),
  bytes: Uint8List.fromList(const [1]),
);

void main() {
  late _MockReceiptApi receiptApi;
  late _MockReceiptImageApi receiptImageApi;

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
    receiptImageApi = _MockReceiptImageApi();
    when(() => client.getReceiptApi()).thenReturn(receiptApi);
    when(() => client.getReceiptImageApi()).thenReturn(receiptImageApi);
    OpenApiClient.client = client;

    Response<api.Receipt> ok() => Response<api.Receipt>(
      requestOptions: RequestOptions(path: ''),
      statusCode: 200,
      data: _receipt(id: 42),
    );
    when(
      () => receiptApi.createReceiptWithFiles(
        receipt: any(named: 'receipt'),
        files: any(named: 'files'),
      ),
    ).thenAnswer((_) async => ok());
    when(
      () => receiptApi.updateReceipt(
        receiptId: any(named: 'receiptId'),
        upsertReceiptCommand: any(named: 'upsertReceiptCommand'),
      ),
    ).thenAnswer((_) async => ok());
  });

  Future<ReceiptFormHarness> pump(
    WidgetTester tester, {
    required WranglerFormState formState,
    required api.ReceiptRequirements requirements,
    api.Receipt? receipt,
  }) {
    return pumpReceiptForm(
      tester,
      formState: formState,
      receipt:
          receipt ?? _receipt(id: formState == WranglerFormState.add ? 0 : 5),
      groups: [
        buildGroup(
          id: _groupId,
          members: [buildGroupMember(userId: _userId, groupId: _groupId)],
        ),
      ],
      users: [buildUserView(id: _userId)],
      realSubmitButton: true,
      permissionsModel: seededPermissions(
        receiptRequirements: {_groupId: requirements},
      ),
    );
  }

  Future<void> submit(WidgetTester tester) async {
    await tester.tap(find.byType(BottomSubmitButton));
    await tester.pump();
    await tester.pump();
  }

  void verifyNoCreate() => verifyNever(
    () => receiptApi.createReceiptWithFiles(
      receipt: any(named: 'receipt'),
      files: any(named: 'files'),
    ),
  );

  group('add', () {
    testWidgets('blocks a receipt missing a required image and comment', (
      tester,
    ) async {
      await pump(
        tester,
        formState: WranglerFormState.add,
        requirements: receiptRequirements(comment: true, image: true),
      );

      await submit(tester);

      expect(find.textContaining('one image and one comment'), findsOneWidget);
      verifyNoCreate();
    });

    testWidgets('blocks a receipt with only whitespace comments', (
      tester,
    ) async {
      final harness = await pump(
        tester,
        formState: WranglerFormState.add,
        requirements: receiptRequirements(comment: true),
      );
      harness.receiptModel.setComments([_comment('   ')]);

      await submit(tester);

      expect(find.textContaining('one comment'), findsOneWidget);
      verifyNoCreate();
    });

    testWidgets(
      'submits once, carrying the staged images and comments, when met',
      (tester) async {
        final harness = await pump(
          tester,
          formState: WranglerFormState.add,
          requirements: receiptRequirements(comment: true, image: true),
        );
        harness.receiptModel.setComments([_comment('for the team')]);
        harness.receiptModel.imagesToUploadBehaviorSubject.add([_staged()]);

        await submit(tester);
        await tester.pumpAndSettle();

        final captured =
            verify(
              () => receiptApi.createReceiptWithFiles(
                receipt: captureAny(named: 'receipt'),
                files: captureAny(named: 'files'),
              ),
            ).captured;
        final command = captured[0] as api.UpsertReceiptCommand;
        final files = captured[1] as BuiltList<MultipartFile>;
        expect(command.comments?.map((c) => c.comment), ['for the team']);
        expect(files.map((f) => f.filename), ['a.png']);
        // The create is atomic; there is no follow-up upload.
        verifyZeroInteractions(receiptImageApi);
        expect(find.text('receipt 42 view'), findsOneWidget);
      },
    );

    testWidgets('a group with no requirements submits without images', (
      tester,
    ) async {
      await pump(
        tester,
        formState: WranglerFormState.add,
        requirements: receiptRequirements(),
      );

      await submit(tester);
      await tester.pumpAndSettle();

      verify(
        () => receiptApi.createReceiptWithFiles(
          receipt: any(named: 'receipt'),
          files: null,
        ),
      ).called(1);
    });
  });

  group('edit', () {
    testWidgets('blocks a save while the receipt has no required comment', (
      tester,
    ) async {
      await pump(
        tester,
        formState: WranglerFormState.edit,
        requirements: receiptRequirements(comment: true),
      );

      await submit(tester);

      expect(find.textContaining('one comment'), findsOneWidget);
      verifyNever(
        () => receiptApi.updateReceipt(
          receiptId: any(named: 'receiptId'),
          upsertReceiptCommand: any(named: 'upsertReceiptCommand'),
        ),
      );
    });

    testWidgets('saves when the stored receipt already meets it', (
      tester,
    ) async {
      await pump(
        tester,
        formState: WranglerFormState.edit,
        requirements: receiptRequirements(comment: true),
        receipt: _receipt(id: 5, comments: [_comment('noted')]),
      );

      await submit(tester);
      await tester.pumpAndSettle();

      verify(
        () => receiptApi.updateReceipt(
          receiptId: 5,
          upsertReceiptCommand: any(named: 'upsertReceiptCommand'),
        ),
      ).called(1);
    });
  });
}
