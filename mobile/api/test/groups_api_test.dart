import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for GroupsApi
void main() {
  final instance = Openapi().getGroupsApi();

  group(GroupsApi, () {
    // Create group
    //
    // This will create a group
    //
    //Future createGroup(Group group) async
    test('test createGroup', () async {
      // TODO
    });

    // Delete group
    //
    // This will delete a group by id
    //
    //Future deleteGroup(int groupId) async
    test('test deleteGroup', () async {
      // TODO
    });

    // Gets a group by Id
    //
    // This will get a group by Id
    //
    //Future getGroupById(int groupId) async
    test('test getGroupById', () async {
      // TODO
    });

    // Get groups for user
    //
    // This will get groups for the currently logged in user
    //
    //Future<BuiltList<Group>> getGroupsForuser() async
    test('test getGroupsForuser', () async {
      // TODO
    });

    // Reads each image in a group and returns the zipped read text
    //
    // This will get the ocr text, zipped, for each image in a group and one text file per image
    //
    //Future getOcrTextForGroup(int groupId) async
    test('test getOcrTextForGroup', () async {
      // TODO
    });

    // Get paged groups
    //
    // This will return paged groups
    //
    //Future<PagedData> getPagedGroups(PagedGroupRequestCommand pagedGroupRequestCommand) async
    test('test getPagedGroups', () async {
      // TODO
    });

    // Poll group email
    //
    // This will poll the group email for new receipts and add them to the group
    //
    //Future pollGroupEmail(int groupId) async
    test('test pollGroupEmail', () async {
      // TODO
    });

    // Update a group
    //
    // This will update a group
    //
    //Future updateGroup(int groupId, Group group) async
    test('test updateGroup', () async {
      // TODO
    });

    // Update group settings
    //
    // This will update the group settings for a group
    //
    //Future<GroupSettings> updateGroupSettings(int groupId, UpdateGroupSettingsCommand updateGroupSettingsCommand) async
    test('test updateGroupSettings', () async {
      // TODO
    });

  });
}
