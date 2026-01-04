import 'package:flutter/material.dart';

class ContextModel extends ChangeNotifier {
  BuildContext? _shellContext;

  BuildContext? get shellContext => _shellContext;

  void setShellContext(BuildContext context) {
    _shellContext = context;
  }
}
