import 'package:flutter/material.dart';
import 'package:infinite_carousel/infinite_carousel.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:provider/provider.dart';
import 'package:receipt_wrangler_mobile/models/receipt_model.dart';
import 'package:receipt_wrangler_mobile/receipts/widgets/receipt_image_carousel.dart';
import 'package:rxdart/rxdart.dart';

class ReceiptImages extends StatefulWidget {
  const ReceiptImages(
      {super.key,
      required this.imageStream,
      required this.infiniteScrollController});

  final Stream<List<api.FileDataView?>> imageStream;

  final InfiniteScrollController infiniteScrollController;

  @override
  State<ReceiptImages> createState() => _ReceiptImages();
}

class _ReceiptImages extends State<ReceiptImages> {
  late final receipt =
      Provider.of<ReceiptModel>(context, listen: false).receipt;
  late final receiptModel = Provider.of<ReceiptModel>(context, listen: false);
  late final imageStream = CombineLatestStream.combine2(
      widget.imageStream.asBroadcastStream(),
      receiptModel.imagesToUploadBehaviorSubject.stream.asBroadcastStream(),
      (a, b) => a).asBroadcastStream();

  @override
  Widget build(BuildContext context) {
    return StreamBuilder(
      stream: imageStream,
      builder: (context, snapshot) {
        if (snapshot.hasData) {
          return SizedBox(
            height: MediaQuery.of(context).size.height,
            width: MediaQuery.of(context).size.width,
            child: ReceiptImageCarousel(
              key: UniqueKey(),
              images: snapshot.data as List<api.FileDataView?>,
              imagesToUpload: receiptModel.imagesToUploadBehaviorSubject.value,
              infiniteScrollController: widget.infiniteScrollController,
            ),
          );
        } else {
          return const Center(child: CircularProgressIndicator());
        }
      },
    );
  }
}
