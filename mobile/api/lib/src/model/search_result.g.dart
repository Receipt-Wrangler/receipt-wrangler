// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'search_result.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$SearchResult extends SearchResult {
  @override
  final int id;
  @override
  final String name;
  @override
  final String type;
  @override
  final int groupId;
  @override
  final String date;
  @override
  final String? amount;
  @override
  final ReceiptStatus? receiptStatus;
  @override
  final int? paidByUserId;
  @override
  final String createdAt;

  factory _$SearchResult([void Function(SearchResultBuilder)? updates]) =>
      (new SearchResultBuilder()..update(updates))._build();

  _$SearchResult._(
      {required this.id,
      required this.name,
      required this.type,
      required this.groupId,
      required this.date,
      this.amount,
      this.receiptStatus,
      this.paidByUserId,
      required this.createdAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'SearchResult', 'id');
    BuiltValueNullFieldError.checkNotNull(name, r'SearchResult', 'name');
    BuiltValueNullFieldError.checkNotNull(type, r'SearchResult', 'type');
    BuiltValueNullFieldError.checkNotNull(groupId, r'SearchResult', 'groupId');
    BuiltValueNullFieldError.checkNotNull(date, r'SearchResult', 'date');
    BuiltValueNullFieldError.checkNotNull(
        createdAt, r'SearchResult', 'createdAt');
  }

  @override
  SearchResult rebuild(void Function(SearchResultBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  SearchResultBuilder toBuilder() => new SearchResultBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is SearchResult &&
        id == other.id &&
        name == other.name &&
        type == other.type &&
        groupId == other.groupId &&
        date == other.date &&
        amount == other.amount &&
        receiptStatus == other.receiptStatus &&
        paidByUserId == other.paidByUserId &&
        createdAt == other.createdAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, type.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, date.hashCode);
    _$hash = $jc(_$hash, amount.hashCode);
    _$hash = $jc(_$hash, receiptStatus.hashCode);
    _$hash = $jc(_$hash, paidByUserId.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'SearchResult')
          ..add('id', id)
          ..add('name', name)
          ..add('type', type)
          ..add('groupId', groupId)
          ..add('date', date)
          ..add('amount', amount)
          ..add('receiptStatus', receiptStatus)
          ..add('paidByUserId', paidByUserId)
          ..add('createdAt', createdAt))
        .toString();
  }
}

class SearchResultBuilder
    implements Builder<SearchResult, SearchResultBuilder> {
  _$SearchResult? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _name;
  String? get name => _$this._name;
  set name(String? name) => _$this._name = name;

  String? _type;
  String? get type => _$this._type;
  set type(String? type) => _$this._type = type;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(int? groupId) => _$this._groupId = groupId;

  String? _date;
  String? get date => _$this._date;
  set date(String? date) => _$this._date = date;

  String? _amount;
  String? get amount => _$this._amount;
  set amount(String? amount) => _$this._amount = amount;

  ReceiptStatus? _receiptStatus;
  ReceiptStatus? get receiptStatus => _$this._receiptStatus;
  set receiptStatus(ReceiptStatus? receiptStatus) =>
      _$this._receiptStatus = receiptStatus;

  int? _paidByUserId;
  int? get paidByUserId => _$this._paidByUserId;
  set paidByUserId(int? paidByUserId) => _$this._paidByUserId = paidByUserId;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  SearchResultBuilder() {
    SearchResult._defaults(this);
  }

  SearchResultBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _name = $v.name;
      _type = $v.type;
      _groupId = $v.groupId;
      _date = $v.date;
      _amount = $v.amount;
      _receiptStatus = $v.receiptStatus;
      _paidByUserId = $v.paidByUserId;
      _createdAt = $v.createdAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(SearchResult other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$SearchResult;
  }

  @override
  void update(void Function(SearchResultBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  SearchResult build() => _build();

  _$SearchResult _build() {
    final _$result = _$v ??
        new _$SearchResult._(
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'SearchResult', 'id'),
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'SearchResult', 'name'),
            type: BuiltValueNullFieldError.checkNotNull(
                type, r'SearchResult', 'type'),
            groupId: BuiltValueNullFieldError.checkNotNull(
                groupId, r'SearchResult', 'groupId'),
            date: BuiltValueNullFieldError.checkNotNull(
                date, r'SearchResult', 'date'),
            amount: amount,
            receiptStatus: receiptStatus,
            paidByUserId: paidByUserId,
            createdAt: BuiltValueNullFieldError.checkNotNull(
                createdAt, r'SearchResult', 'createdAt'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
