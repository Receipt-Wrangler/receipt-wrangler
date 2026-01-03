import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';

class TagModel extends ChangeNotifier {
  List<Tag> _tags = [];

  List<Tag> get tags => _tags;

  void setTags(List<Tag> tags) {
    _tags = tags;

    notifyListeners();
  }
}
