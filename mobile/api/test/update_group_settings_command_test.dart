import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for UpdateGroupSettingsCommand
void main() {
  final instance = UpdateGroupSettingsCommandBuilder();
  // TODO add properties to the builder and call build()

  group(UpdateGroupSettingsCommand, () {
    // System email foreign key
    // int systemEmailId
    test('to test the property `systemEmailId`', () async {
      // TODO
    });

    // Whether email integration is enabled
    // bool emailIntegrationEnabled
    test('to test the property `emailIntegrationEnabled`', () async {
      // TODO
    });

    // Subject line regexes
    // BuiltList<SubjectLineRegex> subjectLineRegexes
    test('to test the property `subjectLineRegexes`', () async {
      // TODO
    });

    // Email white list
    // BuiltList<GroupSettingsWhiteListEmail> emailWhiteList
    test('to test the property `emailWhiteList`', () async {
      // TODO
    });

    // ReceiptStatus emailDefaultReceiptStatus
    test('to test the property `emailDefaultReceiptStatus`', () async {
      // TODO
    });

    // User foreign key
    // int emailDefaultReceiptPaidById
    test('to test the property `emailDefaultReceiptPaidById`', () async {
      // TODO
    });

    // Prompt foreign key
    // int promptId
    test('to test the property `promptId`', () async {
      // TODO
    });

    // Fallback prompt foreign key
    // int fallbackPromptId
    test('to test the property `fallbackPromptId`', () async {
      // TODO
    });

  });
}
