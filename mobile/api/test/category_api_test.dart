import 'package:test/test.dart';
import 'package:openapi/openapi.dart';


/// tests for CategoryApi
void main() {
  final instance = Openapi().getCategoryApi();

  group(CategoryApi, () {
    // Create category
    //
    // This will create a category
    //
    //Future<Category> createCategory(Category category) async
    test('test createCategory', () async {
      // TODO
    });

    // Delete category
    //
    // This will delete a category by id
    //
    //Future deleteCategory(int categoryId) async
    test('test deleteCategory', () async {
      // TODO
    });

    // Get all categories
    //
    // This will return all categories in the system
    //
    //Future<BuiltList<Category>> getAllCategories() async
    test('test getAllCategories', () async {
      // TODO
    });

    // Get category count by name
    //
    // This will return a count of categories with the same name
    //
    //Future<int> getCategoryCountByName(String categoryName) async
    test('test getCategoryCountByName', () async {
      // TODO
    });

    // Get paged categories
    //
    // This will return paged categories
    //
    //Future<PagedData> getPagedCategories(PagedRequestCommand pagedRequestCommand) async
    test('test getPagedCategories', () async {
      // TODO
    });

    // Update category
    //
    // This will update a category
    //
    //Future updateCategory(int categoryId, Category category) async
    test('test updateCategory', () async {
      // TODO
    });

  });
}
