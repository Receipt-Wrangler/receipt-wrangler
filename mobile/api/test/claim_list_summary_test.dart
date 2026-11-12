import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ClaimListSummary
void main() {
  final instance = ClaimListSummaryBuilder();
  // TODO add properties to the builder and call build()

  group(ClaimListSummary, () {
    // Today in the app time zone (yyyy-MM-dd), the day daysWaiting is counted to
    // String asOfDate
    test('to test the property `asOfDate`', () async {
      // TODO
    });

    // Number of claims in the list
    // int claimCount
    test('to test the property `claimCount`', () async {
      // TODO
    });

    // Number of claims with awaitingSettlement true
    // int awaitingCount
    test('to test the property `awaitingCount`', () async {
      // TODO
    });

    // Sum of outstandingAmount over the claims with awaitingSettlement true
    // String outstandingTotal
    test('to test the property `outstandingTotal`', () async {
      // TODO
    });

    // The largest daysWaiting in the list; absent when no claim has one
    // int oldestDaysWaiting
    test('to test the property `oldestDaysWaiting`', () async {
      // TODO
    });

  });
}
