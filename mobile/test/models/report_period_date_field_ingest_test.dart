import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart';

/// Guards the API boundary for `ReportPeriod.dateField`.
///
/// The field names which receipt date a report period covers (`date`,
/// `resolvedDate` or `createdAt`, the receipts quick date filter's keys). It
/// rides inside `ReportTemplate.configuration`, so it arrives on the report
/// templates list, and mobile's preview sends that configuration straight back
/// to `POST /report/preview` (`fetchReportPreview`). Two things follow:
///
/// - It must stay a plain `String?`, never an enum. The quick date filter can
///   gain a field, and a closed built_value `EnumClass` throws on a value it
///   does not know, failing the whole template payload on every
///   already-released build (the mechanism behind the two `Permission` login
///   outages). The last test below fails if a regen ever closes it.
/// - It must survive the round trip, or a previewed template would silently
///   fall back to the receipt date on the server.
///
/// Mobile never reads the field itself, which is why it needs a guard: nothing
/// in the app would notice it going missing.
Map<String, Object?> _templateJson({Object? dateField}) => {
      // ReportTemplate implements BaseModel, so id/createdAt are required.
      'id': 1,
      'createdAt': '2026-01-01',
      'name': 'Monthly Spend',
      'configurationVersion': 1,
      'configuration': {
        'groupIds': ['1'],
        'period': {
          'preset': 'custom',
          'startDate': '2026-06-01',
          'endDate': '2026-06-30',
          if (dateField != null) 'dateField': dateField,
        },
        'detail': {'mode': 'records'},
        'columns': [
          {'kind': 'dimension', 'name': 'Name', 'field': 'name'},
        ],
        'formats': ['csv'],
      },
    };

ReportTemplate _deserialize(Map<String, Object?> json) =>
    standardSerializers.deserializeWith(ReportTemplate.serializer, json)!;

/// Serializes the configuration exactly as `fetchReportPreview` sends it.
Map<Object?, Object?> _previewPeriod(ReportTemplate template) {
  final command = standardSerializers.serializeWith(
    ReportRequestCommand.serializer,
    template.configuration,
  ) as Map<Object?, Object?>;
  return command['period'] as Map<Object?, Object?>;
}

void main() {
  group('ReportPeriod.dateField ingest', () {
    // A template saved before the field existed carries no key. It must still
    // parse, and preview must not invent one: the server reads an absent key
    // as the receipt date, which is what that template always covered.
    test('an absent key deserializes to null and is not re-sent', () {
      final template = _deserialize(_templateJson());

      expect(template.configuration.period.dateField, isNull);
      expect(_previewPeriod(template).containsKey('dateField'), isFalse);
    });

    for (final dateField in ['date', 'resolvedDate', 'createdAt']) {
      test('$dateField deserializes and survives the preview round trip', () {
        final template = _deserialize(_templateJson(dateField: dateField));

        expect(template.configuration.period.dateField, dateField);
        expect(_previewPeriod(template)['dateField'], dateField);
      });
    }

    // The regression this file exists for. A date field added to the desktop
    // picker after this build shipped must not fail the template - it is
    // carried through untouched, and the server decides what it means.
    test('an unknown value deserializes without throwing and is kept', () {
      final template = _deserialize(_templateJson(dateField: 'paidAt'));

      expect(template.name, 'Monthly Spend');
      expect(template.configuration.period.dateField, 'paidAt');
      expect(_previewPeriod(template)['dateField'], 'paidAt');
    });
  });
}
