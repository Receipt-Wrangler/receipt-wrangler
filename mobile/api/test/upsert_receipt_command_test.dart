import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for UpsertReceiptCommand
void main() {
  final instance = UpsertReceiptCommandBuilder();
  // TODO add properties to the builder and call build()

  group(UpsertReceiptCommand, () {
    // Receipt name
    // String name
    test('to test the property `name`', () async {
      // TODO
    });

    // Receipt total amount
    // String amount
    test('to test the property `amount`', () async {
      // TODO
    });

    // Receipt date
    // String date
    test('to test the property `date`', () async {
      // TODO
    });

    // Group foreign key
    // int groupId
    test('to test the property `groupId`', () async {
      // TODO
    });

    // User paid foreign key
    // int paidByUserId
    test('to test the property `paidByUserId`', () async {
      // TODO
    });

    // ReceiptStatus status
    test('to test the property `status`', () async {
      // TODO
    });

    // Categories associated to receipt
    // BuiltList<UpsertCategoryCommand> categories
    test('to test the property `categories`', () async {
      // TODO
    });

    // Tags associated to receipt
    // BuiltList<UpsertTagCommand> tags
    test('to test the property `tags`', () async {
      // TODO
    });

    // Items associated to receipt
    // BuiltList<UpsertItemCommand> receiptItems
    test('to test the property `receiptItems`', () async {
      // TODO
    });

    // Comments associated to receipt
    // BuiltList<UpsertCommentCommand> comments
    test('to test the property `comments`', () async {
      // TODO
    });

  });
}
