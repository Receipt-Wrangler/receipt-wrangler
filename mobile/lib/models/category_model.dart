import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';

class CategoryModel extends ChangeNotifier {
  List<Category> _categories = [];

  // Per-group catalogs the caller may use, keyed by group id. Delivered on
  // AppData (`groupCategories`) filtered to the user's group-role grants (the
  // full pool when unrestricted). Non-admins receive categories only here — the
  // flat `_categories` list is admin-only — so receipt pickers must source from
  // this map, scoped to the receipt's group.
  Map<int, List<Category>> _groupCategories = {};

  List<Category> get categories => _categories;

  void setCategories(List<Category> categories) {
    _categories = categories;

    notifyListeners();
  }

  void setGroupCategories(Map<int, List<Category>> groupCategories) {
    _groupCategories = groupCategories;

    notifyListeners();
  }

  List<Category> categoriesForGroup(int groupId) {
    return _groupCategories[groupId] ?? const [];
  }
}
