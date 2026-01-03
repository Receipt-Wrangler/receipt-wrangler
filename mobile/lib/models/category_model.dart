import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';

class CategoryModel extends ChangeNotifier {
  List<Category> _categories = [];

  List<Category> get categories => _categories;

  void setCategories(List<Category> categories) {
    _categories = categories;

    notifyListeners();
  }
}
