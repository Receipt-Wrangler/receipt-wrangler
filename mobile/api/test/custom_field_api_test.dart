import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for CustomFieldApi
void main() {
  final instance = Openapi().getCustomFieldApi();

  group(CustomFieldApi, () {
    // Create custom field
    //
    // This will create a custom field
    //
    //Future<CustomField> createCustomField(UpsertCustomFieldCommand upsertCustomFieldCommand) async
    test('test createCustomField', () async {
      // TODO
    });

    // Delete custom field
    //
    // This will delete a custom field by id
    //
    //Future deleteCustomField(int customFieldId) async
    test('test deleteCustomField', () async {
      // TODO
    });

    // Get custom field
    //
    // This will get a custom field by id
    //
    //Future<CustomField> getCustomFieldById(int customFieldId) async
    test('test getCustomFieldById', () async {
      // TODO
    });

    // Get paged custom fields
    //
    // This will return paged custom fields
    //
    //Future<PagedData> getPagedCustomFields(PagedRequestCommand pagedRequestCommand) async
    test('test getPagedCustomFields', () async {
      // TODO
    });

  });
}
