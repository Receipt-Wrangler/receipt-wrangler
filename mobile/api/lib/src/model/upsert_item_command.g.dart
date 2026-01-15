// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_item_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertItemCommand extends UpsertItemCommand {
  @override
  final String amount;
  @override
  final int? chargedToUserId;
  @override
  final String name;
  @override
  final int receiptId;
  @override
  final ItemStatus status;
  @override
  final BuiltList<UpsertCategoryCommand>? categories;
  @override
  final BuiltList<UpsertTagCommand>? tags;
  @override
  final BuiltList<UpsertItemCommand>? linkedItems;

  factory _$UpsertItemCommand(
          [void Function(UpsertItemCommandBuilder)? updates]) =>
      (new UpsertItemCommandBuilder()..update(updates))._build();

  _$UpsertItemCommand._(
      {required this.amount,
      this.chargedToUserId,
      required this.name,
      required this.receiptId,
      required this.status,
      this.categories,
      this.tags,
      this.linkedItems})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        amount, r'UpsertItemCommand', 'amount');
    BuiltValueNullFieldError.checkNotNull(name, r'UpsertItemCommand', 'name');
    BuiltValueNullFieldError.checkNotNull(
        receiptId, r'UpsertItemCommand', 'receiptId');
    BuiltValueNullFieldError.checkNotNull(
        status, r'UpsertItemCommand', 'status');
  }

  @override
  UpsertItemCommand rebuild(void Function(UpsertItemCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertItemCommandBuilder toBuilder() =>
      new UpsertItemCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertItemCommand &&
        amount == other.amount &&
        chargedToUserId == other.chargedToUserId &&
        name == other.name &&
        receiptId == other.receiptId &&
        status == other.status &&
        categories == other.categories &&
        tags == other.tags &&
        linkedItems == other.linkedItems;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, amount.hashCode);
    _$hash = $jc(_$hash, chargedToUserId.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, receiptId.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, categories.hashCode);
    _$hash = $jc(_$hash, tags.hashCode);
    _$hash = $jc(_$hash, linkedItems.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertItemCommand')
          ..add('amount', amount)
          ..add('chargedToUserId', chargedToUserId)
          ..add('name', name)
          ..add('receiptId', receiptId)
          ..add('status', status)
          ..add('categories', categories)
          ..add('tags', tags)
          ..add('linkedItems', linkedItems))
        .toString();
  }
}

class UpsertItemCommandBuilder
    implements Builder<UpsertItemCommand, UpsertItemCommandBuilder> {
  _$UpsertItemCommand? _$v;

  String? _amount;
  String? get amount => _$this._amount;
  set amount(String? amount) => _$this._amount = amount;

  int? _chargedToUserId;
  int? get chargedToUserId => _$this._chargedToUserId;
  set chargedToUserId(int? chargedToUserId) =>
      _$this._chargedToUserId = chargedToUserId;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  int? _receiptId;
  int? get receiptId => _$this._receiptId;
  set receiptId(int? receiptId) => _$this._receiptId = receiptId;

  ItemStatus? _status;
  ItemStatus? get status => _$this._status;
  set status(ItemStatus? status) => _$this._status = status;

  ListBuilder<UpsertCategoryCommand>? _categories;
  ListBuilder<UpsertCategoryCommand> get categories =>
      _$this._categories ??= new ListBuilder<UpsertCategoryCommand>();
  set categories(ListBuilder<UpsertCategoryCommand>? categories) =>
      _$this._categories = categories;

  ListBuilder<UpsertTagCommand>? _tags;
  ListBuilder<UpsertTagCommand> get tags =>
      _$this._tags ??= new ListBuilder<UpsertTagCommand>();
  set tags(ListBuilder<UpsertTagCommand>? tags) => _$this._tags = tags;

  ListBuilder<UpsertItemCommand>? _linkedItems;
  ListBuilder<UpsertItemCommand> get linkedItems =>
      _$this._linkedItems ??= new ListBuilder<UpsertItemCommand>();
  set linkedItems(ListBuilder<UpsertItemCommand>? linkedItems) =>
      _$this._linkedItems = linkedItems;

  UpsertItemCommandBuilder() {
    UpsertItemCommand._defaults(this);
  }

  UpsertItemCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _amount = $v.amount;
      _chargedToUserId = $v.chargedToUserId;
      _name = $v.name;
      _receiptId = $v.receiptId;
      _status = $v.status;
      _categories = $v.categories?.toBuilder();
      _tags = $v.tags?.toBuilder();
      _linkedItems = $v.linkedItems?.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertItemCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertItemCommand;
  }

  @override
  void update(void Function(UpsertItemCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertItemCommand build() => _build();

  _$UpsertItemCommand _build() {
    _$UpsertItemCommand _$result;
    try {
      _$result = _$v ??
          new _$UpsertItemCommand._(
              amount: BuiltValueNullFieldError.checkNotNull(
                  amount, r'UpsertItemCommand', 'amount'),
              chargedToUserId: chargedToUserId,
              name: BuiltValueNullFieldError.checkNotNull(
                  name, r'UpsertItemCommand', 'name'),
              receiptId: BuiltValueNullFieldError.checkNotNull(
                  receiptId, r'UpsertItemCommand', 'receiptId'),
              status: BuiltValueNullFieldError.checkNotNull(
                  status, r'UpsertItemCommand', 'status'),
              categories: _categories?.build(),
              tags: _tags?.build(),
              linkedItems: _linkedItems?.build());
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'categories';
        _categories?.build();
        _$failedField = 'tags';
        _tags?.build();
        _$failedField = 'linkedItems';
        _linkedItems?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'UpsertItemCommand', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
