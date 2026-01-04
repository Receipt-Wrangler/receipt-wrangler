import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for SystemTaskApi
void main() {
  final instance = Openapi().getSystemTaskApi();

  group(SystemTaskApi, () {
    // Gets paged system tasks
    //
    // This will return paged system tasks
    //
    //Future<PagedData> getPagedSystemTasks(GetSystemTaskCommand getSystemTaskCommand) async
    test('test getPagedSystemTasks', () async {
      // TODO
    });

  });
}
