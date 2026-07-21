import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ReportTemplate
void main() {
  final instance = ReportTemplateBuilder();
  // TODO add properties to the builder and call build()

  group(ReportTemplate, () {
    // int id
    test('to test the property `id`', () async {
      // TODO
    });

    // String createdAt
    test('to test the property `createdAt`', () async {
      // TODO
    });

    // int createdBy (default value: 0)
    test('to test the property `createdBy`', () async {
      // TODO
    });

    // Created by entity's name
    // String createdByString (default value: '')
    test('to test the property `createdByString`', () async {
      // TODO
    });

    // String updatedAt (default value: '')
    test('to test the property `updatedAt`', () async {
      // TODO
    });

    // The template name (mirrors the saved report's name).
    // String name
    test('to test the property `name`', () async {
      // TODO
    });

    // ReportRequestCommand configuration
    test('to test the property `configuration`', () async {
      // TODO
    });

    // Schema version the stored configuration was written under.
    // int configurationVersion
    test('to test the property `configurationVersion`', () async {
      // TODO
    });

    // The actions the requesting user may perform on this template (read, generate, update, delete, duplicate), resolved per user and populated only on the list response. Drives the row action buttons.
    // BuiltList<String> allowedActions
    test('to test the property `allowedActions`', () async {
      // TODO
    });

  });
}
