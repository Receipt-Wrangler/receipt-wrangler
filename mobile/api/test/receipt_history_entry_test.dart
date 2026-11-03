import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ReceiptHistoryEntry
void main() {
  final instance = ReceiptHistoryEntryBuilder();
  // TODO add properties to the builder and call build()

  group(ReceiptHistoryEntry, () {
    // The system task the entry was built from
    // int systemTaskId
    test('to test the property `systemTaskId`', () async {
      // TODO
    });

    // What the entry did. One of CREATED, UPDATED, STATUS_CHANGED, RESOLVED, IMAGES_COMMENTS_UPDATED, FAILED. A plain string, not an enum, so a client must tolerate values it does not know.
    // String action
    test('to test the property `action`', () async {
      // TODO
    });

    // When the entry was recorded (RFC 3339 instant)
    // String recordedAt
    test('to test the property `recordedAt`', () async {
      // TODO
    });

    // ReceiptHistoryActor actor
    test('to test the property `actor`', () async {
      // TODO
    });

    // Whether getReceiptHistorySnapshots has receipt copies for this entry
    // bool hasSnapshots
    test('to test the property `hasSnapshots`', () async {
      // TODO
    });

    // The error a FAILED entry recorded
    // String errorMessage
    test('to test the property `errorMessage`', () async {
      // TODO
    });

    // How far the compared \"before\" can be trusted, on update entries. One of COMPLETE, REBUILT (an older row compared against an earlier complete copy) or INCOMPLETE (an older row with no such copy). Plain string.
    // String beforeState
    test('to test the property `beforeState`', () async {
      // TODO
    });

    // On a REBUILT entry, the system task type the \"before\" copy came from
    // String beforeSourceType
    test('to test the property `beforeSourceType`', () async {
      // TODO
    });

    // On a REBUILT entry, when the \"before\" copy was recorded (RFC 3339 instant)
    // String beforeSourceRecordedAt
    test('to test the property `beforeSourceRecordedAt`', () async {
      // TODO
    });

    // BuiltList<ReceiptHistoryChange> changes
    test('to test the property `changes`', () async {
      // TODO
    });

  });
}
