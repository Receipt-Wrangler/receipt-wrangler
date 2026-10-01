import 'package:built_collection/built_collection.dart';
import 'package:flutter/material.dart';
import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/models/loading_model.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_image_app_bar.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';

import '../helpers/permission_test_helpers.dart';

// While the receipt's group role requires an image, its last saved image cannot
// be deleted from the edit-state image menu — the server would refuse that with
// a 400. Uploading a second image first re-enables it.
void main() {
  const groupId = 7;

  api.FileDataView image(int id) =>
      (api.FileDataViewBuilder()
            ..id = id
            ..createdAt = ''
            ..encodedImage = ''
            ..name = 'image-$id.png')
          .build();

  Future<ReceiptModel> pump(
    WidgetTester tester, {
    required List<api.FileDataView> images,
    required api.ReceiptRequirements requirements,
  }) async {
    final receiptModel =
        ReceiptModel()..setReceipt(
          getDefaultReceipt().rebuild(
            (b) =>
                b
                  ..id = 1
                  ..groupId = groupId
                  ..imageFiles = ListBuilder<api.FileData>(),
          ),
          false,
        );
    receiptModel.imageBehaviorSubject.add(images);

    await tester.pumpWidget(
      MultiProvider(
        providers: [
          ChangeNotifierProvider.value(value: receiptModel),
          ChangeNotifierProvider(
            create:
                (_) => seededPermissions(
                  group: {
                    groupId: [api.Permission.groupPeriodReceiptsPeriodUpdate],
                  },
                  receiptRequirements: {groupId: requirements},
                ),
          ),
          ChangeNotifierProvider(create: (_) => LoadingModel()),
        ],
        child: MaterialApp(
          home: Scaffold(
            appBar: ReceiptImageAppBar(
              formState: WranglerFormState.edit,
              title: 'Images',
              onBack: () {},
            ),
          ),
        ),
      ),
    );
    await tester.pump();
    return receiptModel;
  }

  Future<PopupMenuItem> openDelete(WidgetTester tester) async {
    await tester.tap(find.byType(PopupMenuButton));
    await tester.pumpAndSettle();
    return tester.widget<PopupMenuItem>(
      find.ancestor(
        of: find.text('Delete Image'),
        matching: find.byType(PopupMenuItem),
      ),
    );
  }

  testWidgets('the only image cannot be deleted while one is required', (
    tester,
  ) async {
    await pump(
      tester,
      images: [image(1)],
      requirements: receiptRequirements(image: true),
    );

    expect((await openDelete(tester)).enabled, isFalse);
  });

  testWidgets('with a second image, delete is enabled', (tester) async {
    await pump(
      tester,
      images: [image(1), image(2)],
      requirements: receiptRequirements(image: true),
    );

    expect((await openDelete(tester)).enabled, isTrue);
  });

  testWidgets('an upload re-enables delete on the same menu', (tester) async {
    final model = await pump(
      tester,
      images: [image(1)],
      requirements: receiptRequirements(image: true),
    );

    // The entries follow the image list rather than being fixed at build.
    model.imageBehaviorSubject.add([image(1), image(2)]);
    await tester.pumpAndSettle();

    expect((await openDelete(tester)).enabled, isTrue);
  });

  testWidgets('without the requirement the last image is deletable', (
    tester,
  ) async {
    await pump(
      tester,
      images: [image(1)],
      requirements: receiptRequirements(comment: true),
    );

    expect((await openDelete(tester)).enabled, isTrue);
  });
}
