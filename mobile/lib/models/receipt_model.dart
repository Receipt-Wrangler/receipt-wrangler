import 'package:flutter/material.dart';
import 'package:flutter/widgets.dart';
import 'package:flutter_form_builder/flutter_form_builder.dart';
import 'package:infinite_carousel/infinite_carousel.dart';
import 'package:openapi/openapi.dart';
import 'package:receipt_wrangler_mobile/interfaces/form_item.dart';
import 'package:receipt_wrangler_mobile/interfaces/upload_multipart_file_data.dart';
import 'package:receipt_wrangler_mobile/utils/receipts.dart';
import 'package:rxdart/rxdart.dart';

class ReceiptModel extends ChangeNotifier {
  Receipt _receipt = getDefaultReceipt();

  Receipt get receipt => _receipt;

  Receipt _modifiedReceipt = getDefaultReceipt();

  Receipt get modifiedReceipt => _modifiedReceipt;

  List<Comment> _comments = [];

  List<Comment> get comments => _comments;

  List<FormItem> _items = [];

  List<FormItem> get items => _items;

  BehaviorSubject<List<FileDataView?>> _imageBehaviorSubject =
      BehaviorSubject<List<FileDataView?>>.seeded([]);

  BehaviorSubject<List<FileDataView?>> get imageBehaviorSubject =>
      _imageBehaviorSubject;

  BehaviorSubject<List<UploadMultipartFileData>>
      _imagesToUploadBehaviorSubject =
      BehaviorSubject<List<UploadMultipartFileData>>.seeded([]);

  BehaviorSubject<List<UploadMultipartFileData>>
      get imagesToUploadBehaviorSubject => _imagesToUploadBehaviorSubject;

  InfiniteScrollController _infiniteScrollController =
      InfiniteScrollController();

  InfiniteScrollController get infiniteScrollController =>
      _infiniteScrollController;

  var _receiptFormKey = GlobalKey<FormBuilderState>();

  GlobalKey<FormBuilderState> get receiptFormKey => _receiptFormKey;

  var _quickActionsFormKey = GlobalKey<FormBuilderState>();

  GlobalKey<FormBuilderState> get quickActionsFormKey => _quickActionsFormKey;

  void setReceipt(Receipt receipt, bool notify) {
    _receipt = receipt;

    _modifiedReceipt = receipt;

    _comments = (receipt.comments)?.toList() ?? [];

    _items = FormItem.fromItems((receipt.receiptItems)?.toList() ?? []);

    _imageBehaviorSubject = BehaviorSubject<List<FileDataView?>>.seeded([]);

    _imagesToUploadBehaviorSubject =
        BehaviorSubject<List<UploadMultipartFileData>>.seeded([]);

    _receiptFormKey = GlobalKey<FormBuilderState>();

    if (notify) {
      notifyListeners();
    }
  }

  void setComments(List<Comment> comments) {
    _comments = comments;
    notifyListeners();
  }

  void setItems(List<FormItem> items) {
    _items = items;
    notifyListeners();
  }

  void setModifiedReceipt(Receipt receipt) {
    _modifiedReceipt = receipt;
    notifyListeners();
  }

  void resetQuickActionsFormKey() {
    _quickActionsFormKey = GlobalKey<FormBuilderState>();
  }

  void resetModel() {
    _receipt = getDefaultReceipt();
    _modifiedReceipt = getDefaultReceipt();
    _comments = [];
    _items = [];
    _imageBehaviorSubject = BehaviorSubject<List<FileDataView?>>.seeded([]);
    _imagesToUploadBehaviorSubject =
        BehaviorSubject<List<UploadMultipartFileData>>.seeded([]);
    _infiniteScrollController = InfiniteScrollController();
    _receiptFormKey = GlobalKey<FormBuilderState>();
  }
}
