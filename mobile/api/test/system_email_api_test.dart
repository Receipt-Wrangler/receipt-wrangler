import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for SystemEmailApi
void main() {
  final instance = Openapi().getSystemEmailApi();

  group(SystemEmailApi, () {
    // Check system email connectivity
    //
    // This will check system email connectivity
    //
    //Future<SystemTask> checkSystemEmailConnectivity(CheckEmailConnectivityCommand checkEmailConnectivityCommand) async
    test('test checkSystemEmailConnectivity', () async {
      // TODO
    });

    // Create system email
    //
    // This will create a system email
    //
    //Future<SystemEmail> createSystemEmail(UpsertSystemEmailCommand upsertSystemEmailCommand) async
    test('test createSystemEmail', () async {
      // TODO
    });

    // Delete system email by id
    //
    // This will delete a system email by id
    //
    //Future deleteSystemEmailById(int id) async
    test('test deleteSystemEmailById', () async {
      // TODO
    });

    // Gets paged system emails
    //
    // This will return paged and sorted system emails
    //
    //Future<PagedData> getPagedSystemEmails(PagedRequestCommand pagedRequestCommand) async
    test('test getPagedSystemEmails', () async {
      // TODO
    });

    // Get system email by id
    //
    // This will get a system email by id
    //
    //Future<SystemEmail> getSystemEmailById(int id) async
    test('test getSystemEmailById', () async {
      // TODO
    });

    // Update system email by id
    //
    // This will update a system email by id
    //
    //Future<SystemEmail> updateSystemEmailById(int id, bool updatePassword, UpsertSystemEmailCommand upsertSystemEmailCommand) async
    test('test updateSystemEmailById', () async {
      // TODO
    });

  });
}
