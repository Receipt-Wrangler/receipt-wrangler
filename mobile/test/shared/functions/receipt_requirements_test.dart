import 'dart:typed_data';

import 'package:built_collection/built_collection.dart';
import 'package:dio/dio.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/interfaces/upload_multipart_file_data.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/shared/functions/receipt_requirements.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';

import '../../helpers/permission_test_helpers.dart';

api.Comment _comment(String text) =>
    (api.CommentBuilder()
          ..id = 1
          ..comment = text
          ..receiptId = 1
          ..userId = 1)
        .build();

api.FileData _fileData() =>
    (api.FileDataBuilder()
          ..id = 1
          ..createdAt = ''
          ..name = 'a.png'
          ..fileType = 'image/png'
          ..receiptId = 1
          ..size = 1)
        .build();

api.FileDataView _fileDataView() =>
    (api.FileDataViewBuilder()
          ..id = 1
          ..createdAt = ''
          ..encodedImage = ''
          ..name = 'a.png')
        .build();

UploadMultipartFileData _staged() => UploadMultipartFileData(
  multipartFile: MultipartFile.fromBytes(const [1], filename: 'a.png'),
  bytes: Uint8List.fromList(const [1]),
);

ReceiptModel _model({List<api.FileData>? imageFiles}) =>
    ReceiptModel()..setReceipt(
      getDefaultReceipt().rebuild(
        (b) =>
            b
              ..id = 1
              ..groupId = 7
              ..imageFiles =
                  imageFiles == null
                      ? null
                      : ListBuilder<api.FileData>(imageFiles),
      ),
      false,
    );

void main() {
  group('hasNonBlankComment', () {
    test('counts only comments with text, as the server does', () {
      expect(hasNonBlankComment([]), isFalse);
      expect(hasNonBlankComment([_comment('   ')]), isFalse);
      expect(hasNonBlankComment([_comment(' '), _comment('ok')]), isTrue);
      expect(
        countNonBlankComments([_comment(''), _comment('a'), _comment('b')]),
        2,
      );
    });
  });

  group('missingReceiptRequirementsMessage', () {
    test('null when nothing is required', () {
      expect(
        missingReceiptRequirementsMessage(
          receiptRequirements(),
          hasComment: false,
          hasImage: false,
        ),
        isNull,
      );
    });

    test('null when every requirement is met', () {
      expect(
        missingReceiptRequirementsMessage(
          receiptRequirements(comment: true, image: true),
          hasComment: true,
          hasImage: true,
        ),
        isNull,
      );
    });

    test('names exactly what is missing', () {
      final both = receiptRequirements(comment: true, image: true);
      expect(
        missingReceiptRequirementsMessage(
          both,
          hasComment: false,
          hasImage: false,
        ),
        contains('one image and one comment'),
      );
      expect(
        missingReceiptRequirementsMessage(
          both,
          hasComment: true,
          hasImage: false,
        ),
        allOf(contains('image'), isNot(contains('comment'))),
      );
      expect(
        missingReceiptRequirementsMessage(
          both,
          hasComment: false,
          hasImage: true,
        ),
        allOf(contains('comment'), isNot(contains('image'))),
      );
    });
  });

  group('receiptHasImage', () {
    test('add state reads the staged images', () {
      final model = _model(imageFiles: [_fileData()]);
      // A saved image does not count for an unsaved receipt.
      expect(receiptHasImage(model, WranglerFormState.add), isFalse);

      model.imagesToUploadBehaviorSubject.add([_staged()]);
      expect(receiptHasImage(model, WranglerFormState.add), isTrue);
    });

    test(
      'edit state reads the saved images before the images screen loads',
      () {
        expect(receiptHasImage(_model(), WranglerFormState.edit), isFalse);
        expect(
          receiptHasImage(
            _model(imageFiles: [_fileData()]),
            WranglerFormState.edit,
          ),
          isTrue,
        );
      },
    );

    test('edit state counts an image uploaded this session', () {
      final model = _model();
      model.imageBehaviorSubject.add([_fileDataView()]);

      expect(receiptHasImage(model, WranglerFormState.edit), isTrue);
    });

    test('edit state ignores staged images', () {
      final model = _model();
      model.imagesToUploadBehaviorSubject.add([_staged()]);

      expect(receiptHasImage(model, WranglerFormState.edit), isFalse);
    });
  });

  group('receiptSubmitRequirementsMessage', () {
    test('reads the model\'s comments and images', () {
      final model = _model();
      final requirements = receiptRequirements(comment: true, image: true);

      expect(
        receiptSubmitRequirementsMessage(
          requirements,
          receiptModel: model,
          formState: WranglerFormState.add,
        ),
        isNotNull,
      );

      model.setComments([_comment('looks right')]);
      model.imagesToUploadBehaviorSubject.add([_staged()]);

      expect(
        receiptSubmitRequirementsMessage(
          requirements,
          receiptModel: model,
          formState: WranglerFormState.add,
        ),
        isNull,
      );
    });
  });

  group('isLastRequiredItem', () {
    test('only blocks the last one while required', () {
      expect(isLastRequiredItem(required: true, remaining: 1), isTrue);
      expect(isLastRequiredItem(required: true, remaining: 2), isFalse);
      expect(isLastRequiredItem(required: false, remaining: 1), isFalse);
    });
  });
}
