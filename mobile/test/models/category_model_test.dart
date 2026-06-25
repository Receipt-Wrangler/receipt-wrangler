import 'package:flutter_test/flutter_test.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/models/category_model.dart';

Category _category(int id, String name) =>
    Category((b) => b
      ..id = id
      ..name = name);

void main() {
  group('CategoryModel.categoriesForGroup', () {
    test('returns the per-group catalog for a known group', () {
      final model = CategoryModel();
      final food = _category(1, 'Food');
      model.setGroupCategories({
        7: [food],
      });

      expect(model.categoriesForGroup(7), [food]);
    });

    test('isolates catalogs per group', () {
      final model = CategoryModel();
      final food = _category(1, 'Food');
      final travel = _category(2, 'Travel');
      model.setGroupCategories({
        7: [food],
        9: [travel],
      });

      expect(model.categoriesForGroup(7), [food]);
      expect(model.categoriesForGroup(9), [travel]);
    });

    test('returns an empty list for an unknown group', () {
      final model = CategoryModel();
      model.setGroupCategories({
        7: [_category(1, 'Food')],
      });

      expect(model.categoriesForGroup(99), isEmpty);
    });

    test('returns an empty list before hydration', () {
      expect(CategoryModel().categoriesForGroup(7), isEmpty);
    });

    test('setGroupCategories notifies listeners', () {
      final model = CategoryModel();
      var notified = 0;
      model.addListener(() => notified++);
      model.setGroupCategories({});
      expect(notified, 1);
    });
  });
}
