import 'package:test/test.dart';
import 'package:openapi/openapi.dart';

// tests for UpdateGroupMemberGrantsCommand
void main() {
  final instance = UpdateGroupMemberGrantsCommandBuilder();
  // TODO add properties to the builder and call build()

  group(UpdateGroupMemberGrantsCommand, () {
    // Category ids to assign to this member. Every id must sit within the ceiling set by the member's group role, or the request is rejected with 400. An empty array clears the member's category restriction, handing them back to their role's set.
    // BuiltList<int> categoryGrants
    test('to test the property `categoryGrants`', () async {
      // TODO
    });

    // Tag counterpart of categoryGrants.
    // BuiltList<int> tagGrants
    test('to test the property `tagGrants`', () async {
      // TODO
    });

  });
}
