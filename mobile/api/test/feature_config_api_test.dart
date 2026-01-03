import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for FeatureConfigApi
void main() {
  final instance = Openapi().getFeatureConfigApi();

  group(FeatureConfigApi, () {
    // Get feature config
    //
    // This will get the server's feature config
    //
    //Future<FeatureConfig> getFeatureConfig() async
    test('test getFeatureConfig', () async {
      // TODO
    });

  });
}
