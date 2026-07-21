import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart' as api;
import 'package:receipt_wrangler_mobile/reports/functions/report_filename.dart';

import '../helpers/report_test_helpers.dart';

void main() {
  group('reportFilename', () {
    test('single format uses that extension', () {
      final command = buildReportRequestCommand(
        name: 'Q1 Report',
        formats: [api.ReportRequestCommandFormatsEnum.pdf],
      );
      expect(reportFilename(command), 'Q1_Report.pdf');
    });

    test('xlsx single format', () {
      final command = buildReportRequestCommand(
        name: 'Expenses',
        formats: [api.ReportRequestCommandFormatsEnum.xlsx],
      );
      expect(reportFilename(command), 'Expenses.xlsx');
    });

    test('multiple formats zip', () {
      final command = buildReportRequestCommand(
        name: 'Everything',
        formats: [
          api.ReportRequestCommandFormatsEnum.csv,
          api.ReportRequestCommandFormatsEnum.pdf,
        ],
      );
      expect(reportFilename(command), 'Everything.zip');
    });

    test('empty formats fall back to csv', () {
      final command = buildReportRequestCommand(name: 'NoFormat', formats: []);
      expect(reportFilename(command), 'NoFormat.csv');
    });

    test('sanitizes non-word characters and collapses runs', () {
      final command = buildReportRequestCommand(
        name: 'Q1 / 2024: Report!!',
        formats: [api.ReportRequestCommandFormatsEnum.csv],
      );
      expect(reportFilename(command), 'Q1_2024_Report.csv');
    });

    test('empty/blank name falls back to "report"', () {
      final command = buildReportRequestCommand(
        name: '   ',
        formats: [api.ReportRequestCommandFormatsEnum.csv],
      );
      expect(reportFilename(command), 'report.csv');
    });

    test('null name falls back to "report"', () {
      final command = buildReportRequestCommand(
        name: null,
        formats: [api.ReportRequestCommandFormatsEnum.pdf],
      );
      expect(reportFilename(command), 'report.pdf');
    });
  });
}
