// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_receipt_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertReceiptCommand extends UpsertReceiptCommand {
  @override
  final String name;
  @override
  final String amount;
  @override
  final String date;
  @override
  final int groupId;
  @override
  final int paidByUserId;
  @override
  final ReceiptStatus status;
  @override
  final BuiltList<UpsertCategoryCommand>? categories;
  @override
  final BuiltList<UpsertTagCommand>? tags;
  @override
  final BuiltList<UpsertItemCommand>? receiptItems;
  @override
  final BuiltList<UpsertCommentCommand>? comments;
  @override
  final BuiltList<UpsertCustomFieldValueCommand>? customFields;

  factory _$UpsertReceiptCommand(
          [void Function(UpsertReceiptCommandBuilder)? updates]) =>
      (new UpsertReceiptCommandBuilder()..update(updates))._build();

  _$UpsertReceiptCommand._(
      {required this.name,
      required this.amount,
      required this.date,
      required this.groupId,
      required this.paidByUserId,
      required this.status,
      this.categories,
      this.tags,
      this.receiptItems,
      this.comments,
      this.customFields})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        name, r'UpsertReceiptCommand', 'name');
    BuiltValueNullFieldError.checkNotNull(
        amount, r'UpsertReceiptCommand', 'amount');
    BuiltValueNullFieldError.checkNotNull(
        date, r'UpsertReceiptCommand', 'date');
    BuiltValueNullFieldError.checkNotNull(
        groupId, r'UpsertReceiptCommand', 'groupId');
    BuiltValueNullFieldError.checkNotNull(
        paidByUserId, r'UpsertReceiptCommand', 'paidByUserId');
    BuiltValueNullFieldError.checkNotNull(
        status, r'UpsertReceiptCommand', 'status');
  }

  @override
  UpsertReceiptCommand rebuild(
          void Function(UpsertReceiptCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertReceiptCommandBuilder toBuilder() =>
      new UpsertReceiptCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertReceiptCommand &&
        name == other.name &&
        amount == other.amount &&
        date == other.date &&
        groupId == other.groupId &&
        paidByUserId == other.paidByUserId &&
        status == other.status &&
        categories == other.categories &&
        tags == other.tags &&
        receiptItems == other.receiptItems &&
        comments == other.comments &&
        customFields == other.customFields;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, amount.hashCode);
    _$hash = $jc(_$hash, date.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, paidByUserId.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, categories.hashCode);
    _$hash = $jc(_$hash, tags.hashCode);
    _$hash = $jc(_$hash, receiptItems.hashCode);
    _$hash = $jc(_$hash, comments.hashCode);
    _$hash = $jc(_$hash, customFields.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertReceiptCommand')
          ..add('name', name)
          ..add('amount', amount)
          ..add('date', date)
          ..add('groupId', groupId)
          ..add('paidByUserId', paidByUserId)
          ..add('status', status)
          ..add('categories', categories)
          ..add('tags', tags)
          ..add('receiptItems', receiptItems)
          ..add('comments', comments)
          ..add('customFields', customFields))
        .toString();
  }
}

class UpsertReceiptCommandBuilder
    implements Builder<UpsertReceiptCommand, UpsertReceiptCommandBuilder> {
  _$UpsertReceiptCommand? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _amount;
  String? get amount => _$this._amount;
  set amount(String? amount) => _$this._amount = amount;

  String? _date;
  String? get date => _$this._date;
  set date(String? date) => _$this._date = date;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(int? groupId) => _$this._groupId = groupId;

  int? _paidByUserId;
  int? get paidByUserId => _$this._paidByUserId;
  set paidByUserId(int? paidByUserId) => _$this._paidByUserId = paidByUserId;

  ReceiptStatus? _status;
  ReceiptStatus? get status => _$this._status;
  set status(ReceiptStatus? status) => _$this._status = status;

  ListBuilder<UpsertCategoryCommand>? _categories;
  ListBuilder<UpsertCategoryCommand> get categories =>
      _$this._categories ??= new ListBuilder<UpsertCategoryCommand>();
  set categories(ListBuilder<UpsertCategoryCommand>? categories) =>
      _$this._categories = categories;

  ListBuilder<UpsertTagCommand>? _tags;
  ListBuilder<UpsertTagCommand> get tags =>
      _$this._tags ??= new ListBuilder<UpsertTagCommand>();
  set tags(ListBuilder<UpsertTagCommand>? tags) => _$this._tags = tags;

  ListBuilder<UpsertItemCommand>? _receiptItems;
  ListBuilder<UpsertItemCommand> get receiptItems =>
      _$this._receiptItems ??= new ListBuilder<UpsertItemCommand>();
  set receiptItems(ListBuilder<UpsertItemCommand>? receiptItems) =>
      _$this._receiptItems = receiptItems;

  ListBuilder<UpsertCommentCommand>? _comments;
  ListBuilder<UpsertCommentCommand> get comments =>
      _$this._comments ??= new ListBuilder<UpsertCommentCommand>();
  set comments(ListBuilder<UpsertCommentCommand>? comments) =>
      _$this._comments = comments;

  ListBuilder<UpsertCustomFieldValueCommand>? _customFields;
  ListBuilder<UpsertCustomFieldValueCommand> get customFields =>
      _$this._customFields ??= new ListBuilder<UpsertCustomFieldValueCommand>();
  set customFields(ListBuilder<UpsertCustomFieldValueCommand>? customFields) =>
      _$this._customFields = customFields;

  UpsertReceiptCommandBuilder() {
    UpsertReceiptCommand._defaults(this);
  }

  UpsertReceiptCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _amount = $v.amount;
      _date = $v.date;
      _groupId = $v.groupId;
      _paidByUserId = $v.paidByUserId;
      _status = $v.status;
      _categories = $v.categories?.toBuilder();
      _tags = $v.tags?.toBuilder();
      _receiptItems = $v.receiptItems?.toBuilder();
      _comments = $v.comments?.toBuilder();
      _customFields = $v.customFields?.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertReceiptCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertReceiptCommand;
  }

  @override
  void update(void Function(UpsertReceiptCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertReceiptCommand build() => _build();

  _$UpsertReceiptCommand _build() {
    _$UpsertReceiptCommand _$result;
    try {
      _$result = _$v ??
          new _$UpsertReceiptCommand._(
              name: BuiltValueNullFieldError.checkNotNull(
                  name, r'UpsertReceiptCommand', 'name'),
              amount: BuiltValueNullFieldError.checkNotNull(
                  amount, r'UpsertReceiptCommand', 'amount'),
              date: BuiltValueNullFieldError.checkNotNull(
                  date, r'UpsertReceiptCommand', 'date'),
              groupId: BuiltValueNullFieldError.checkNotNull(
                  groupId, r'UpsertReceiptCommand', 'groupId'),
              paidByUserId: BuiltValueNullFieldError.checkNotNull(
                  paidByUserId, r'UpsertReceiptCommand', 'paidByUserId'),
              status: BuiltValueNullFieldError.checkNotNull(
                  status, r'UpsertReceiptCommand', 'status'),
              categories: _categories?.build(),
              tags: _tags?.build(),
              receiptItems: _receiptItems?.build(),
              comments: _comments?.build(),
              customFields: _customFields?.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'categories';
        _categories?.build();
        _$failedField = 'tags';
        _tags?.build();
        _$failedField = 'receiptItems';
        _receiptItems?.build();
        _$failedField = 'comments';
        _comments?.build();
        _$failedField = 'customFields';
        _customFields?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'UpsertReceiptCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
