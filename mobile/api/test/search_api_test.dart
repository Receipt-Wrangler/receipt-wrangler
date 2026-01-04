import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for SearchApi
void main() {
  final instance = Openapi().getSearchApi();

  group(SearchApi, () {
    // Receipt Search
    //
    // This will search for receipts based on a search term
    //
    //Future<BuiltList<SearchResult>> receiptSearch(String searchTerm) async
    test('test receiptSearch', () async {
      // TODO
    });

  });
}
