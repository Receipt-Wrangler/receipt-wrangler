import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart';

/// Guards the API boundary for the report split options
/// (`splitCategoriesEqually`, `splitTagsEqually`, `splitExcludedFields`).
///
/// They ride inside `ReportTemplate.configuration`, so they arrive on the report
/// templates list, and mobile's preview sends that configuration straight back
/// to `POST /report/preview` (`fetchReportPreview`). Mobile never reads them, so
/// nothing in the app would notice them going missing:
///
/// - A template saved before they existed omits every key. It must still parse,
///   and preview must not invent the keys.
/// - A template that uses them must survive the round trip, or its preview
///   would silently attribute every receipt in full to each category.
/// - `splitExcludedFields` must stay a list of plain strings, never an enum: the
///   entries are `custom_<id>` keys, an open set.
Map<String, Object?> _templateJson({Map<String, Object?> split = const {}}) => {
      // ReportTemplate implements BaseModel, so id/createdAt are required.
      'id': 1,
      'createdAt': '2026-01-01',
      'name': 'Split Spend',
      'configurationVersion': 1,
      'configuration': {
        'groupIds': ['1'],
        'period': {'preset': 'this_month'},
        'detail': {'mode': 'aggregate', 'by': 'category'},
        'columns': [
          {'kind': 'dimension', 'name': 'Category', 'field': 'category'},
        ],
        'formats': ['csv'],
        ...split,
      },
    };

ReportTemplate _deserialize(Map<String, Object?> json) =>
    standardSerializers.deserializeWith(ReportTemplate.serializer, json)!;

/// Serializes the configuration exactly as `fetchReportPreview` sends it.
Map<Object?, Object?> _previewCommand(ReportTemplate template) =>
    standardSerializers.serializeWith(
      ReportRequestCommand.serializer,
      template.configuration,
    ) as Map<Object?, Object?>;

void main() {
  group('Report split options ingest', () {
    test('absent keys deserialize to null and are not re-sent', () {
      final template = _deserialize(_templateJson());
      final configuration = template.configuration;

      expect(configuration.splitCategoriesEqually, isNull);
      expect(configuration.splitTagsEqually, isNull);
      expect(configuration.splitExcludedFields, isNull);

      final command = _previewCommand(template);
      expect(command.containsKey('splitCategoriesEqually'), isFalse);
      expect(command.containsKey('splitTagsEqually'), isFalse);
      expect(command.containsKey('splitExcludedFields'), isFalse);
    });

    test('set options deserialize and survive the preview round trip', () {
      final template = _deserialize(_templateJson(split: {
        'splitCategoriesEqually': true,
        'splitTagsEqually': true,
        'splitExcludedFields': ['custom_3', 'custom_12'],
      }));
      final configuration = template.configuration;

      expect(configuration.splitCategoriesEqually, isTrue);
      expect(configuration.splitTagsEqually, isTrue);
      expect(configuration.splitExcludedFields?.toList(),
          ['custom_3', 'custom_12']);

      final command = _previewCommand(template);
      expect(command['splitCategoriesEqually'], isTrue);
      expect(command['splitTagsEqually'], isTrue);
      expect(command['splitExcludedFields'], ['custom_3', 'custom_12']);
    });

    // An excluded key the client has never seen is carried through untouched;
    // the server decides whether it names a field.
    test('an arbitrary excluded key deserializes without throwing', () {
      final template = _deserialize(_templateJson(split: {
        'splitExcludedFields': ['custom_999999'],
      }));

      expect(template.name, 'Split Spend');
      expect(template.configuration.splitExcludedFields?.toList(),
          ['custom_999999']);
    });
  });
}
