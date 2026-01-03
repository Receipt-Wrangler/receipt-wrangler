import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for ReceiptImageApi
void main() {
  final instance = Openapi().getReceiptImageApi();

  group(ReceiptImageApi, () {
    // Converts a receipt image to jpg
    //
    // This will convert a receipt image to jpg, [SYSTEM USER]
    //
    //Future<EncodedImage> convertToJpg(MultipartFile file) async
    test('test convertToJpg', () async {
      // TODO
    });

    // Delete receipt image
    //
    // This will delete a receipt image by id [SYSTEM USER]
    //
    //Future deleteReceiptImageById(int receiptImageId) async
    test('test deleteReceiptImageById', () async {
      // TODO
    });

    // Download receipt image
    //
    // This will download a receipt image by id, [SYSTEM USER]
    //
    //Future<Uint8List> downloadReceiptImageById(int receiptImageId) async
    test('test downloadReceiptImageById', () async {
      // TODO
    });

    // Get receipt image
    //
    // This will get a receipt image by id, [SYSTEM USER]
    //
    //Future<FileDataView> getReceiptImageById(int receiptImageId) async
    test('test getReceiptImageById', () async {
      // TODO
    });

    // Reads a receipt image and returns the parsed results
    //
    // This will parse and read a receipt image, [SYSTEM USER]
    //
    //Future<Receipt> magicFillReceipt({ int receiptImageId, MultipartFile file }) async
    test('test magicFillReceipt', () async {
      // TODO
    });

    // Uploads a receipt image
    //
    // This will upload a receipt image, [SYSTEM USER]
    //
    //Future<FileDataView> uploadReceiptImage(MultipartFile file, int receiptId, { String encodedImage }) async
    test('test uploadReceiptImage', () async {
      // TODO
    });

  });
}
