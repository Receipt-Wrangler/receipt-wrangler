import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for UpsertSystemSettingsCommand
void main() {
  final instance = UpsertSystemSettingsCommandBuilder();
  // TODO add properties to the builder and call build()

  group(UpsertSystemSettingsCommand, () {
    // Whether local sign up is enabled
    // bool enableLocalSignUp
    test('to test the property `enableLocalSignUp`', () async {
      // TODO
    });

    // Currency display
    // String currencyDisplay
    test('to test the property `currencyDisplay`', () async {
      // TODO
    });

    // bool debugOcr
    test('to test the property `debugOcr`', () async {
      // TODO
    });

    // Number of workers to use
    // int numWorkers (default value: 1)
    test('to test the property `numWorkers`', () async {
      // TODO
    });

    // Email polling interval
    // int emailPollingInterval
    test('to test the property `emailPollingInterval`', () async {
      // TODO
    });

    // Receipt processing settings foreign key
    // int receiptProcessingSettingsId
    test('to test the property `receiptProcessingSettingsId`', () async {
      // TODO
    });

    // Fallback receipt processing settings foreign key
    // int fallbackReceiptProcessingSettingsId
    test('to test the property `fallbackReceiptProcessingSettingsId`', () async {
      // TODO
    });

  });
}
