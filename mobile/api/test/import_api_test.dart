import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for ImportApi
void main() {
  final instance = Openapi().getImportApi();

  group(ImportApi, () {
    // Import config json
    //
    // This will import a config json
    //
    //Future importConfigJson(MultipartFile file) async
    test('test importConfigJson', () async {
      // TODO
    });

  });
}
