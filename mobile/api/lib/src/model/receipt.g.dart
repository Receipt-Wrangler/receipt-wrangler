// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$Receipt extends Receipt {
  @override
  final String amount;
  @override
  final BuiltList<Category> categories;
  @override
  final BuiltList<Comment> comments;
  @override
  final BuiltList<CustomFieldValue> customFields;
  @override
  final String? createdAt;
  @override
  final int? createdBy;
  @override
  final String date;
  @override
  final int groupId;
  @override
  final int id;
  @override
  final BuiltList<FileData>? imageFiles;
  @override
  final String name;
  @override
  final int paidByUserId;
  @override
  final BuiltList<Item> receiptItems;
  @override
  final String? resolvedDate;
  @override
  final ReceiptStatus status;
  @override
  final BuiltList<Tag> tags;
  @override
  final String? updatedAt;
  @override
  final String? createdByString;

  factory _$Receipt([void Function(ReceiptBuilder)? updates]) =>
      (new ReceiptBuilder()..update(updates))._build();

  _$Receipt._(
      {required this.amount,
      required this.categories,
      required this.comments,
      required this.customFields,
      this.createdAt,
      this.createdBy,
      required this.date,
      required this.groupId,
      required this.id,
      this.imageFiles,
      required this.name,
      required this.paidByUserId,
      required this.receiptItems,
      this.resolvedDate,
      required this.status,
      required this.tags,
      this.updatedAt,
      this.createdByString})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(amount, r'Receipt', 'amount');
    BuiltValueNullFieldError.checkNotNull(categories, r'Receipt', 'categories');
    BuiltValueNullFieldError.checkNotNull(comments, r'Receipt', 'comments');
    BuiltValueNullFieldError.checkNotNull(
        customFields, r'Receipt', 'customFields');
    BuiltValueNullFieldError.checkNotNull(date, r'Receipt', 'date');
    BuiltValueNullFieldError.checkNotNull(groupId, r'Receipt', 'groupId');
    BuiltValueNullFieldError.checkNotNull(id, r'Receipt', 'id');
    BuiltValueNullFieldError.checkNotNull(name, r'Receipt', 'name');
    BuiltValueNullFieldError.checkNotNull(
        paidByUserId, r'Receipt', 'paidByUserId');
    BuiltValueNullFieldError.checkNotNull(
        receiptItems, r'Receipt', 'receiptItems');
    BuiltValueNullFieldError.checkNotNull(status, r'Receipt', 'status');
    BuiltValueNullFieldError.checkNotNull(tags, r'Receipt', 'tags');
  }

  @override
  Receipt rebuild(void Function(ReceiptBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptBuilder toBuilder() => new ReceiptBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Receipt &&
        amount == other.amount &&
        categories == other.categories &&
        comments == other.comments &&
        customFields == other.customFields &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        date == other.date &&
        groupId == other.groupId &&
        id == other.id &&
        imageFiles == other.imageFiles &&
        name == other.name &&
        paidByUserId == other.paidByUserId &&
        receiptItems == other.receiptItems &&
        resolvedDate == other.resolvedDate &&
        status == other.status &&
        tags == other.tags &&
        updatedAt == other.updatedAt &&
        createdByString == other.createdByString;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, amount.hashCode);
    _$hash = $jc(_$hash, categories.hashCode);
    _$hash = $jc(_$hash, comments.hashCode);
    _$hash = $jc(_$hash, customFields.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, date.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, imageFiles.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, paidByUserId.hashCode);
    _$hash = $jc(_$hash, receiptItems.hashCode);
    _$hash = $jc(_$hash, resolvedDate.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, tags.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, createdByString.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Receipt')
          ..add('amount', amount)
          ..add('categories', categories)
          ..add('comments', comments)
          ..add('customFields', customFields)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('date', date)
          ..add('groupId', groupId)
          ..add('id', id)
          ..add('imageFiles', imageFiles)
          ..add('name', name)
          ..add('paidByUserId', paidByUserId)
          ..add('receiptItems', receiptItems)
          ..add('resolvedDate', resolvedDate)
          ..add('status', status)
          ..add('tags', tags)
          ..add('updatedAt', updatedAt)
          ..add('createdByString', createdByString))
        .toString();
  }
}

class ReceiptBuilder implements Builder<Receipt, ReceiptBuilder> {
  _$Receipt? _$v;

  String? _amount;
  String? get amount => _$this._amount;
  set amount(String? amount) => _$this._amount = amount;

  ListBuilder<Category>? _categories;
  ListBuilder<Category> get categories =>
      _$this._categories ??= new ListBuilder<Category>();
  set categories(ListBuilder<Category>? categories) =>
      _$this._categories = categories;

  ListBuilder<Comment>? _comments;
  ListBuilder<Comment> get comments =>
      _$this._comments ??= new ListBuilder<Comment>();
  set comments(ListBuilder<Comment>? comments) => _$this._comments = comments;

  ListBuilder<CustomFieldValue>? _customFields;
  ListBuilder<CustomFieldValue> get customFields =>
      _$this._customFields ??= new ListBuilder<CustomFieldValue>();
  set customFields(ListBuilder<CustomFieldValue>? customFields) =>
      _$this._customFields = customFields;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  String? _date;
  String? get date => _$this._date;
  set date(String? date) => _$this._date = date;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(int? groupId) => _$this._groupId = groupId;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  ListBuilder<FileData>? _imageFiles;
  ListBuilder<FileData> get imageFiles =>
      _$this._imageFiles ??= new ListBuilder<FileData>();
  set imageFiles(ListBuilder<FileData>? imageFiles) =>
      _$this._imageFiles = imageFiles;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  int? _paidByUserId;
  int? get paidByUserId => _$this._paidByUserId;
  set paidByUserId(int? paidByUserId) => _$this._paidByUserId = paidByUserId;

  ListBuilder<Item>? _receiptItems;
  ListBuilder<Item> get receiptItems =>
      _$this._receiptItems ??= new ListBuilder<Item>();
  set receiptItems(ListBuilder<Item>? receiptItems) =>
      _$this._receiptItems = receiptItems;

  String? _resolvedDate;
  String? get resolvedDate => _$this._resolvedDate;
  set resolvedDate(String? resolvedDate) => _$this._resolvedDate = resolvedDate;

  ReceiptStatus? _status;
  ReceiptStatus? get status => _$this._status;
  set status(ReceiptStatus? status) => _$this._status = status;

  ListBuilder<Tag>? _tags;
  ListBuilder<Tag> get tags => _$this._tags ??= new ListBuilder<Tag>();
  set tags(ListBuilder<Tag>? tags) => _$this._tags = tags;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  String? _createdByString;
  String? get createdByString => _$this._createdByString;
  set createdByString(String? createdByString) =>
      _$this._createdByString = createdByString;

  ReceiptBuilder() {
    Receipt._defaults(this);
  }

  ReceiptBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _amount = $v.amount;
      _categories = $v.categories.toBuilder();
      _comments = $v.comments.toBuilder();
      _customFields = $v.customFields.toBuilder();
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _date = $v.date;
      _groupId = $v.groupId;
      _id = $v.id;
      _imageFiles = $v.imageFiles?.toBuilder();
      _name = $v.name;
      _paidByUserId = $v.paidByUserId;
      _receiptItems = $v.receiptItems.toBuilder();
      _resolvedDate = $v.resolvedDate;
      _status = $v.status;
      _tags = $v.tags.toBuilder();
      _updatedAt = $v.updatedAt;
      _createdByString = $v.createdByString;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Receipt other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$Receipt;
  }

  @override
  void update(void Function(ReceiptBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Receipt build() => _build();

  _$Receipt _build() {
    _$Receipt _$result;
    try {
      _$result = _$v ??
          new _$Receipt._(
              amount: BuiltValueNullFieldError.checkNotNull(
                  amount, r'Receipt', 'amount'),
              categories: categories.build(),
              comments: comments.build(),
              customFields: customFields.build(),
              createdAt: createdAt,
              createdBy: createdBy,
              date: BuiltValueNullFieldError.checkNotNull(
                  date, r'Receipt', 'date'),
              groupId: BuiltValueNullFieldError.checkNotNull(
                  groupId, r'Receipt', 'groupId'),
              id: BuiltValueNullFieldError.checkNotNull(id, r'Receipt', 'id'),
              imageFiles: _imageFiles?.build(),
              name: BuiltValueNullFieldError.checkNotNull(
                  name, r'Receipt', 'name'),
              paidByUserId: BuiltValueNullFieldError.checkNotNull(
                  paidByUserId, r'Receipt', 'paidByUserId'),
              receiptItems: receiptItems.build(),
              resolvedDate: resolvedDate,
              status: BuiltValueNullFieldError.checkNotNull(
                  status, r'Receipt', 'status'),
              tags: tags.build(),
              updatedAt: updatedAt,
              createdByString: createdByString);
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'categories';
        categories.build();
        _$failedField = 'comments';
        comments.build();
        _$failedField = 'customFields';
        customFields.build();

        _$failedField = 'imageFiles';
        _imageFiles?.build();

        _$failedField = 'receiptItems';
        receiptItems.build();

        _$failedField = 'tags';
        tags.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'Receipt', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
