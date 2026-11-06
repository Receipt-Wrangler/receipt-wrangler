import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';

class UserModel extends ChangeNotifier {
  List<UserView> _users = [];

  List<UserView> get users => _users;

  void setUsers(List<UserView> users) {
    _users = users;

    notifyListeners();
  }

  /// Null when no visible user has [id] -- a deleted user, or one member
  /// isolation hides. It used to throw instead, despite the nullable return.
  UserView? getUserById(String id) {
    for (final user in _users) {
      if (user.id.toString() == id) {
        return user;
      }
    }
    return null;
  }
}
