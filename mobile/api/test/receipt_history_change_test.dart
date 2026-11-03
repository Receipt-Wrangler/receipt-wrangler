import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for ReceiptHistoryChange
void main() {
  final instance = ReceiptHistoryChangeBuilder();
  // TODO add properties to the builder and call build()

  group(ReceiptHistoryChange, () {
    // Which field changed. One of NAME, AMOUNT, DATE, GROUP, PAID_BY, STATUS, RESOLVED_DATE, CATEGORIES, TAGS, CUSTOM_FIELD, ITEM, SHARE, ITEM_CATEGORIES, ITEM_TAGS, COMMENT, IMAGE, EMAIL_INBOX, EMAIL_SENDER. Plain string.
    // String field
    test('to test the property `field`', () async {
      // TODO
    });

    // One of CHANGED, ADDED, REMOVED. Plain string.
    // String kind
    test('to test the property `kind`', () async {
      // TODO
    });

    // What the change is about, where the field alone does not say: the item's name (ITEM, ITEM_CATEGORIES, ITEM_TAGS) or the custom field's name (CUSTOM_FIELD)
    // String label
    test('to test the property `label`', () async {
      // TODO
    });

    // The share's charged-to user (SHARE) or the comment's author (COMMENT)
    // int userId
    test('to test the property `userId`', () async {
      // TODO
    });

    // ReceiptHistoryValue from
    test('to test the property `from`', () async {
      // TODO
    });

    // ReceiptHistoryValue to
    test('to test the property `to`', () async {
      // TODO
    });

  });
}
