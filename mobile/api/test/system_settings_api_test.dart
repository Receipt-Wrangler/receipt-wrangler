import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for SystemSettingsApi
void main() {
  final instance = Openapi().getSystemSettingsApi();

  group(SystemSettingsApi, () {
    // Get system settings
    //
    // This will get system settings
    //
    //Future<SystemSettings> getSystemSettings() async
    test('test getSystemSettings', () async {
      // TODO
    });

    // Update system settings
    //
    // This will update system settings
    //
    //Future<SystemSettings> updateSystemSettings(UpsertSystemSettingsCommand upsertSystemSettingsCommand) async
    test('test updateSystemSettings', () async {
      // TODO
    });

  });
}
