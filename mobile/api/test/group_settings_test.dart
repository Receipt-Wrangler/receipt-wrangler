import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for GroupSettings
void main() {
  final instance = GroupSettingsBuilder();
  // TODO add properties to the builder and call build()

  group(GroupSettings, () {
    // Group settings id
    // int id
    test('to test the property `id`', () async {
      // TODO
    });

    // Group foreign key
    // int groupId
    test('to test the property `groupId`', () async {
      // TODO
    });

    // Whether email integration is enabled
    // bool emailIntegrationEnabled
    test('to test the property `emailIntegrationEnabled`', () async {
      // TODO
    });

    // System email foreign key
    // int systemEmailId
    test('to test the property `systemEmailId`', () async {
      // TODO
    });

    // SystemEmail systemEmail
    test('to test the property `systemEmail`', () async {
      // TODO
    });

    // Email to read
    // String emailToRead
    test('to test the property `emailToRead`', () async {
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

    // Prompt prompt
    test('to test the property `prompt`', () async {
      // TODO
    });

    // Prompt foreign key
    // int promptId
    test('to test the property `promptId`', () async {
      // TODO
    });

    // Prompt fallbackPrompt
    test('to test the property `fallbackPrompt`', () async {
      // TODO
    });

    // Fallback prompt foreign key
    // int fallbackPromptId
    test('to test the property `fallbackPromptId`', () async {
      // TODO
    });

    // String createdAt
    test('to test the property `createdAt`', () async {
      // TODO
    });

    // int createdBy
    test('to test the property `createdBy`', () async {
      // TODO
    });

    // String updatedAt
    test('to test the property `updatedAt`', () async {
      // TODO
    });

  });
}
