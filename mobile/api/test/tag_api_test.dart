import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for TagApi
void main() {
  final instance = Openapi().getTagApi();

  group(TagApi, () {
    // Create tag
    //
    // This will create a tag
    //
    //Future<Tag> createTag(UpsertTagCommand upsertTagCommand) async
    test('test createTag', () async {
      // TODO
    });

    // Delete tag
    //
    // This will delete a tag by id
    //
    //Future deleteTag(int tagId) async
    test('test deleteTag', () async {
      // TODO
    });

    // Get all tags
    //
    // This will return all tags in the system
    //
    //Future<BuiltList<Tag>> getAllTags() async
    test('test getAllTags', () async {
      // TODO
    });

    // Get paged tags
    //
    // This will return paged tags
    //
    //Future<PagedData> getPagedTags(PagedRequestCommand pagedRequestCommand) async
    test('test getPagedTags', () async {
      // TODO
    });

    // Get tag count by name
    //
    // This will count of names with the same name
    //
    //Future<int> getTagCountByName(String tagName) async
    test('test getTagCountByName', () async {
      // TODO
    });

    // Update tag
    //
    // This will update a tag
    //
    //Future<Tag> updateTag(int tagId, UpsertTagCommand upsertTagCommand) async
    test('test updateTag', () async {
      // TODO
    });

  });
}
