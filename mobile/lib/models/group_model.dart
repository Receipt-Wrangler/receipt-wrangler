import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';

class GroupModel extends ChangeNotifier {
  List<Group> _groups = [];

  List<Group> get groups => _groups;

  List<Group> get groupsWithoutAllGroup =>
      _groups.where((group) => !group.isAllGroup).toList();

  void setGroups(List<Group> groups) {
    _groups = groups;

    notifyListeners();
  }

  Group? getGroupById(String id) {
    try {
      return _groups.firstWhere((group) => group.id == int.tryParse(id));
    } catch (e) {
      return null;
    }
  }

  GroupReceiptSettings? getGroupReceiptSettings(int groupId) {
    if (groupId == 0) {
      return null;
    }

    return getGroupById(groupId.toString())?.groupReceiptSettings;
  }
}
