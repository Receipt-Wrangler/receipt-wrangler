import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ReportColumn
void main() {
  final instance = ReportColumnBuilder();
  // TODO add properties to the builder and call build()

  group(ReportColumn, () {
    // String kind
    test('to test the property `kind`', () async {
      // TODO
    });

    // Machine identifier; formula expressions reference this
    // String name
    test('to test the property `name`', () async {
      // TODO
    });

    // String label
    test('to test the property `label`', () async {
      // TODO
    });

    // Field key the column displays (dimension columns)
    // String field
    test('to test the property `field`', () async {
      // TODO
    });

    // Aggregate function (aggregate columns)
    // String aggFunc
    test('to test the property `aggFunc`', () async {
      // TODO
    });

    // Measure field key (aggregate columns, omitted for COUNT)
    // String measure
    test('to test the property `measure`', () async {
      // TODO
    });

    // Expression over other column names (formula columns)
    // String expr
    test('to test the property `expr`', () async {
      // TODO
    });

  });
}
