// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'item.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$Item extends Item {
  @override
  final bool? isTaxed;
  @override
  final String amount;
  @override
  final int? chargedToUserId;
  @override
  final String? createdAt;
  @override
  final int? createdBy;
  @override
  final int? id;
  @override
  final String name;
  @override
  final int receiptId;
  @override
  final ItemStatus status;
  @override
  final BuiltList<Item>? linkedItems;
  @override
  final BuiltList<Category>? categories;
  @override
  final BuiltList<Tag>? tags;
  @override
  final String? updatedAt;

  factory _$Item([void Function(ItemBuilder)? updates]) =>
      (new ItemBuilder()..update(updates))._build();

  _$Item._(
      {this.isTaxed,
      required this.amount,
      this.chargedToUserId,
      this.createdAt,
      this.createdBy,
      this.id,
      required this.name,
      required this.receiptId,
      required this.status,
      this.linkedItems,
      this.categories,
      this.tags,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(amount, r'Item', 'amount');
    BuiltValueNullFieldError.checkNotNull(name, r'Item', 'name');
    BuiltValueNullFieldError.checkNotNull(receiptId, r'Item', 'receiptId');
    BuiltValueNullFieldError.checkNotNull(status, r'Item', 'status');
  }

  @override
  Item rebuild(void Function(ItemBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ItemBuilder toBuilder() => new ItemBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is Item &&
        isTaxed == other.isTaxed &&
        amount == other.amount &&
        chargedToUserId == other.chargedToUserId &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        id == other.id &&
        name == other.name &&
        receiptId == other.receiptId &&
        status == other.status &&
        linkedItems == other.linkedItems &&
        categories == other.categories &&
        tags == other.tags &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, isTaxed.hashCode);
    _$hash = $jc(_$hash, amount.hashCode);
    _$hash = $jc(_$hash, chargedToUserId.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, receiptId.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, linkedItems.hashCode);
    _$hash = $jc(_$hash, categories.hashCode);
    _$hash = $jc(_$hash, tags.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'Item')
          ..add('isTaxed', isTaxed)
          ..add('amount', amount)
          ..add('chargedToUserId', chargedToUserId)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('id', id)
          ..add('name', name)
          ..add('receiptId', receiptId)
          ..add('status', status)
          ..add('linkedItems', linkedItems)
          ..add('categories', categories)
          ..add('tags', tags)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class ItemBuilder implements Builder<Item, ItemBuilder> {
  _$Item? _$v;

  bool? _isTaxed;
  bool? get isTaxed => _$this._isTaxed;
  set isTaxed(bool? isTaxed) => _$this._isTaxed = isTaxed;

  String? _amount;
  String? get amount => _$this._amount;
  set amount(String? amount) => _$this._amount = amount;

  int? _chargedToUserId;
  int? get chargedToUserId => _$this._chargedToUserId;
  set chargedToUserId(int? chargedToUserId) =>
      _$this._chargedToUserId = chargedToUserId;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  int? _receiptId;
  int? get receiptId => _$this._receiptId;
  set receiptId(int? receiptId) => _$this._receiptId = receiptId;

  ItemStatus? _status;
  ItemStatus? get status => _$this._status;
  set status(ItemStatus? status) => _$this._status = status;

  ListBuilder<Item>? _linkedItems;
  ListBuilder<Item> get linkedItems =>
      _$this._linkedItems ??= new ListBuilder<Item>();
  set linkedItems(ListBuilder<Item>? linkedItems) =>
      _$this._linkedItems = linkedItems;

  ListBuilder<Category>? _categories;
  ListBuilder<Category> get categories =>
      _$this._categories ??= new ListBuilder<Category>();
  set categories(ListBuilder<Category>? categories) =>
      _$this._categories = categories;

  ListBuilder<Tag>? _tags;
  ListBuilder<Tag> get tags => _$this._tags ??= new ListBuilder<Tag>();
  set tags(ListBuilder<Tag>? tags) => _$this._tags = tags;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  ItemBuilder() {
    Item._defaults(this);
  }

  ItemBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _isTaxed = $v.isTaxed;
      _amount = $v.amount;
      _chargedToUserId = $v.chargedToUserId;
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _id = $v.id;
      _name = $v.name;
      _receiptId = $v.receiptId;
      _status = $v.status;
      _linkedItems = $v.linkedItems?.toBuilder();
      _categories = $v.categories?.toBuilder();
      _tags = $v.tags?.toBuilder();
      _updatedAt = $v.updatedAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(Item other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$Item;
  }

  @override
  void update(void Function(ItemBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  Item build() => _build();

  _$Item _build() {
    _$Item _$result;
    try {
      _$result = _$v ??
          new _$Item._(
              isTaxed: isTaxed,
              amount: BuiltValueNullFieldError.checkNotNull(
                  amount, r'Item', 'amount'),
              chargedToUserId: chargedToUserId,
              createdAt: createdAt,
              createdBy: createdBy,
              id: id,
              name:
                  BuiltValueNullFieldError.checkNotNull(name, r'Item', 'name'),
              receiptId: BuiltValueNullFieldError.checkNotNull(
                  receiptId, r'Item', 'receiptId'),
              status: BuiltValueNullFieldError.checkNotNull(
                  status, r'Item', 'status'),
              linkedItems: _linkedItems?.build(),
              categories: _categories?.build(),
              tags: _tags?.build(),
              updatedAt: updatedAt);
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'linkedItems';
        _linkedItems?.build();
        _$failedField = 'categories';
        _categories?.build();
        _$failedField = 'tags';
        _tags?.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'Item', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
