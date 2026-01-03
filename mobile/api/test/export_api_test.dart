import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for ExportApi
void main() {
  final instance = Openapi().getExportApi();

  group(ExportApi, () {
    // Exports receipts
    //
    // This will export individual receipts [SYSTEM USER]
    //
    //Future<Uint8List> exportReceiptsById(ExportFormat format, { BuiltList<int> receiptIds }) async
    test('test exportReceiptsById', () async {
      // TODO
    });

    // Exports receipts
    //
    // This will export all receipts that belong to a group based on a filter [SYSTEM USER]
    //
    //Future<Uint8List> exportReceiptsForGroup(ExportFormat format, int groupId, ReceiptPagedRequestCommand receiptPagedRequestCommand) async
    test('test exportReceiptsForGroup', () async {
      // TODO
    });

  });
}
