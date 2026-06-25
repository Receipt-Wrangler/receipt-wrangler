import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';

class TagModel extends ChangeNotifier {
  List<Tag> _tags = [];

  // Per-group catalogs the caller may use, keyed by group id. Delivered on
  // AppData (`groupTags`) filtered to the user's group-role grants (the full
  // pool when unrestricted). Non-admins receive tags only here — the flat
  // `_tags` list is admin-only — so receipt pickers must source from this map,
  // scoped to the receipt's group.
  Map<int, List<Tag>> _groupTags = {};

  List<Tag> get tags => _tags;

  void setTags(List<Tag> tags) {
    _tags = tags;

    notifyListeners();
  }

  void setGroupTags(Map<int, List<Tag>> groupTags) {
    _groupTags = groupTags;

    notifyListeners();
  }

  List<Tag> tagsForGroup(int groupId) {
    return _groupTags[groupId] ?? const [];
  }
}
