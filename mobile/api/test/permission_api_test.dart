import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for PermissionApi
void main() {
  final instance = Openapi().getPermissionApi();

  group(PermissionApi, () {
    // List all permission descriptors
    //
    // Returns the full catalog of permission strings the API recognizes, with metadata for UI rendering.
    //
    //Future<BuiltList<PermissionDescriptor>> getPermissions() async
    test('test getPermissions', () async {
      // TODO
    });

  });
}
