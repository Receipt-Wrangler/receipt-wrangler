import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for ApiKeyApi
void main() {
  final instance = Openapi().getApiKeyApi();

  group(ApiKeyApi, () {
    // Create API key
    //
    // Create a new API key for the authenticated user
    //
    //Future<ApiKeyResult> createApiKey(UpsertApiKeyCommand upsertApiKeyCommand) async
    test('test createApiKey', () async {
      // TODO
    });

    // Delete API key
    //
    // Delete an API key. Admins can delete any API key, non-admins can only delete their own API keys.
    //
    //Future deleteApiKey(String id) async
    test('test deleteApiKey', () async {
      // TODO
    });

    // Get paged API keys
    //
    // This will return paged API keys for the authenticated user or all API keys for admins
    //
    //Future<PagedData> getPagedApiKeys(PagedApiKeyRequestCommand pagedApiKeyRequestCommand) async
    test('test getPagedApiKeys', () async {
      // TODO
    });

    // Update API key
    //
    // This will update an API key. Users can only update their own API keys.
    //
    //Future updateApiKey(String id, UpsertApiKeyCommand upsertApiKeyCommand) async
    test('test updateApiKey', () async {
      // TODO
    });

  });
}
