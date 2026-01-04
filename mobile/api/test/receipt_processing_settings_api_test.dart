import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for ReceiptProcessingSettingsApi
void main() {
  final instance = Openapi().getReceiptProcessingSettingsApi();

  group(ReceiptProcessingSettingsApi, () {
    // Check receipt processing settings connectivity
    //
    //Future<SystemTask> checkReceiptProcessingSettingsConnectivity(CheckReceiptProcessingSettingsConnectivityCommand checkReceiptProcessingSettingsConnectivityCommand) async
    test('test checkReceiptProcessingSettingsConnectivity', () async {
      // TODO
    });

    // Create receipt processing settings
    //
    // This will create receipt processing settings
    //
    //Future<ReceiptProcessingSettings> createReceiptProcessingSettings(UpsertReceiptProcessingSettingsCommand upsertReceiptProcessingSettingsCommand) async
    test('test createReceiptProcessingSettings', () async {
      // TODO
    });

    // Delete receipt processing settings by id
    //
    // This will delete receipt processing settings by id
    //
    //Future deleteReceiptProcessingSettingsById(int id) async
    test('test deleteReceiptProcessingSettingsById', () async {
      // TODO
    });

    // Gets paged processing settings
    //
    // This will return paged processing settings
    //
    //Future<PagedData> getPagedProcessingSettings(PagedRequestCommand pagedRequestCommand) async
    test('test getPagedProcessingSettings', () async {
      // TODO
    });

    // Get receipt processing settings by id
    //
    // This will get receipt processing settings by id
    //
    //Future<ReceiptProcessingSettings> getReceiptProcessingSettingsById(int id) async
    test('test getReceiptProcessingSettingsById', () async {
      // TODO
    });

    // Update receipt processing settings by id
    //
    // This will update receipt processing settings by id
    //
    //Future<ReceiptProcessingSettings> updateReceiptProcessingSettingsById(int id, bool updateKey, UpsertReceiptProcessingSettingsCommand upsertReceiptProcessingSettingsCommand) async
    test('test updateReceiptProcessingSettingsById', () async {
      // TODO
    });

  });
}
