import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ReportRequestCommand
void main() {
  final instance = ReportRequestCommandBuilder();
  // TODO add properties to the builder and call build()

  group(ReportRequestCommand, () {
    // Report name; used for the download filename
    // String name
    test('to test the property `name`', () async {
      // TODO
    });

    // Ids of the groups the report covers
    // BuiltList<String> groupIds
    test('to test the property `groupIds`', () async {
      // TODO
    });

    // ReportPeriod period
    test('to test the property `period`', () async {
      // TODO
    });

    // Which receipts go into the report
    // ReceiptPagedRequestFilter filter
    test('to test the property `filter`', () async {
      // TODO
    });

    // Ordered engine field keys to nest the report by
    // BuiltList<String> groupBy
    test('to test the property `groupBy`', () async {
      // TODO
    });

    // ReportDetail detail
    test('to test the property `detail`', () async {
      // TODO
    });

    // BuiltList<ReportColumn> columns
    test('to test the property `columns`', () async {
      // TODO
    });

    // Emit a subtotal row at each grouping level
    // bool subtotals
    test('to test the property `subtotals`', () async {
      // TODO
    });

    // Emit one grand-total row across everything
    // bool grandTotals
    test('to test the property `grandTotals`', () async {
      // TODO
    });

    // ReportDocument document
    test('to test the property `document`', () async {
      // TODO
    });

    // One or more output formats; multiple are zipped together
    // BuiltList<String> formats
    test('to test the property `formats`', () async {
      // TODO
    });

  });
}
