import 'package:flutter/material.dart';

class LoadingModel extends ChangeNotifier {
  bool _isLoading = false;

  bool get isLoading => _isLoading;

  void setIsLoading(bool isLoading) {
    _isLoading = isLoading;

    notifyListeners();
  }
}
