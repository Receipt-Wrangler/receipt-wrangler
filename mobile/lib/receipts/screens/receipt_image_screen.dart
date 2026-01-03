import 'package:dio/dio.dart';
import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/enums/form_state.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_images.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_image_app_bar.dart';
import 'package:rxdart/rxdart.dart';

import '../../client/client.dart';
import '../../models/receipt_model.dart';

class ReceiptImageScreen extends StatefulWidget {
  const ReceiptImageScreen({super.key, required this.formState});
  
  final WranglerFormState formState;

  @override
  State<ReceiptImageScreen> createState() => _ReceiptImageScreen();
}

class _ReceiptImageScreen extends State<ReceiptImageScreen> {
  late final receipt =
      Provider.of<ReceiptModel>(context, listen: false).receipt;
  late final future = getReceiptImageFutures(receipt);
  late final formState = widget.formState;
  late final receiptModel = Provider.of<ReceiptModel>(context, listen: false);
  var isLoadingBehaviorSubject = BehaviorSubject<bool>();

  @override
  void initState() {
    super.initState();
    isLoadingBehaviorSubject.add(true);
  }

  @override
  void didChangeDependencies() {
    super.didChangeDependencies();

    future.then((responses) {
      var images = <api.FileDataView?>[];

      responses.toList().forEach((response) {
        images.add(response.data);
      });
      receiptModel.imageBehaviorSubject.add(images);
      isLoadingBehaviorSubject.add(false);
    }).catchError((err) {
      isLoadingBehaviorSubject.add(false);
      return err;
    });
  }

  Future<List<Response<api.FileDataView?>>> getReceiptImageFutures(
      api.Receipt receipt) {
    if (formState == WranglerFormState.add) {
      return Future.wait([]);
    } else {
      List<Future<Response<api.FileDataView?>>> imageFutures = [];
      for (var image in receiptModel.receipt.imageFiles?.toList() ?? []) {
        var future = OpenApiClient.client
            .getReceiptImageApi()
            .getReceiptImageById(receiptImageId: image.id);
        imageFutures.add(future);
      }

      return Future.wait(imageFutures);
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: ReceiptImageAppBar(
        formState: formState,
        title: 'Receipt Images',
        onBack: () => Navigator.of(context).pop(),
        onEditMode: formState == WranglerFormState.view ? () {
          Navigator.of(context).pushReplacement(
            MaterialPageRoute(
              builder: (context) => ReceiptImageScreen(formState: WranglerFormState.edit),
            ),
          );
        } : null,
      ),
      body: StreamBuilder<bool>(
          stream: isLoadingBehaviorSubject.stream,
          builder: (context, snapshot) {
            Widget widget = ReceiptImages(
              imageStream: receiptModel.imageBehaviorSubject.stream,
              infiniteScrollController: receiptModel.infiniteScrollController,
            );

            if (snapshot.hasData && snapshot.data == true) {
              widget = const Center(child: CircularProgressIndicator());
            }

            return widget;
          }),
    );
  }
}
