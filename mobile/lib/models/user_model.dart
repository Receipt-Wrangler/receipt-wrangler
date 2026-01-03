import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';

class UserModel extends ChangeNotifier {
  List<UserView> _users = [];

  List<UserView> get users => _users;

  void setUsers(List<UserView> users) {
    _users = users;

    notifyListeners();
  }

  UserView? getUserById(String id) {
    return _users.firstWhere((user) => user.id.toString() == id);
  }
}
