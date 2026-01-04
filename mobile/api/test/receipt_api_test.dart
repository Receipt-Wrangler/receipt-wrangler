import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for ReceiptApi
void main() {
  final instance = Openapi().getReceiptApi();

  group(ReceiptApi, () {
    // Bulk receipt status update
    //
    // This will bulk update receipt statuses with the option of adding a comment to each [SYSTEM USER]
    //
    //Future<BuiltList<Receipt>> bulkReceiptStatusUpdate(BulkStatusUpdateCommand bulkStatusUpdateCommand) async
    test('test bulkReceiptStatusUpdate', () async {
      // TODO
    });

    // Create receipt
    //
    // This will create a receipt [SYSTEM USER]
    //
    //Future<Receipt> createReceipt(UpsertReceiptCommand upsertReceiptCommand) async
    test('test createReceipt', () async {
      // TODO
    });

    // Delete receipt
    //
    // This will delete a receipt by id [SYSTEM USER]
    //
    //Future deleteReceiptById(int receiptId) async
    test('test deleteReceiptById', () async {
      // TODO
    });

    // Duplicate receipt
    //
    // This will duplicate a receipt [SYSTEM USER]
    //
    //Future duplicateReceipt(int receiptId) async
    test('test duplicateReceipt', () async {
      // TODO
    });

    // Get receipt
    //
    // This will get a receipt by receipt id [SYSTEM USER]
    //
    //Future<Receipt> getReceiptById(int receiptId) async
    test('test getReceiptById', () async {
      // TODO
    });

    // Gets receipts
    //
    // This will return receipts with the option to sort and filter [SYSTEM USER]
    //
    //Future<PagedData> getReceiptsForGroup(int groupId, ReceiptPagedRequestCommand receiptPagedRequestCommand) async
    test('test getReceiptsForGroup', () async {
      // TODO
    });

    // Has access to receipt
    //
    // This will return whether or not the currently logged in user has access to the receipt
    //
    //Future hasAccessToReceipt(int receiptId, { String groupRole }) async
    test('test hasAccessToReceipt', () async {
      // TODO
    });

    // Quick scan a receipt
    //
    // This take an image and use magic fill to fill and save the receipt [SYSTEM USER]
    //
    //Future<BuiltList<Receipt>> quickScanReceipt(BuiltList<MultipartFile> files, BuiltList<int> groupIds, BuiltList<int> paidByUserIds, BuiltList<ReceiptStatus> statuses) async
    test('test quickScanReceipt', () async {
      // TODO
    });

    // Update receipt
    //
    // This will update a receipt by receipt id [SYSTEM USER]
    //
    //Future<Receipt> updateReceipt(int receiptId, UpsertReceiptCommand upsertReceiptCommand) async
    test('test updateReceipt', () async {
      // TODO
    });

  });
}
