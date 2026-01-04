import 'package:flutter/material.dart';
import 'package:openapi/openapi.dart';
import 'package:rxdart/rxdart.dart';

class SearchModel extends ChangeNotifier {
  BehaviorSubject<String> _searchTermBehaviorSubject =
      BehaviorSubject<String>();

  BehaviorSubject<String> get searchTermBehaviorSubject =>
      _searchTermBehaviorSubject;

  List<SearchResult> _searchResults = [];

  List<SearchResult> get searchResults => _searchResults;

  void setSearchResults(List<SearchResult> searchResults, {notify = true}) {
    _searchResults = searchResults;

    if (notify) {
      notifyListeners();
    }
  }
}
