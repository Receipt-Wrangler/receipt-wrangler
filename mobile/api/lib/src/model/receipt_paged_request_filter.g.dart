// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'receipt_paged_request_filter.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ReceiptPagedRequestFilter extends ReceiptPagedRequestFilter {
  @override
  final JsonObject? date;
  @override
  final JsonObject? amount;
  @override
  final JsonObject? name;
  @override
  final JsonObject? paidBy;
  @override
  final JsonObject? categories;
  @override
  final JsonObject? tags;
  @override
  final JsonObject? status;
  @override
  final JsonObject? resolvedDate;
  @override
  final JsonObject? createdAt;

  factory _$ReceiptPagedRequestFilter(
          [void Function(ReceiptPagedRequestFilterBuilder)? updates]) =>
      (new ReceiptPagedRequestFilterBuilder()..update(updates))._build();

  _$ReceiptPagedRequestFilter._(
      {this.date,
      this.amount,
      this.name,
      this.paidBy,
      this.categories,
      this.tags,
      this.status,
      this.resolvedDate,
      this.createdAt})
      : super._();

  @override
  ReceiptPagedRequestFilter rebuild(
          void Function(ReceiptPagedRequestFilterBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ReceiptPagedRequestFilterBuilder toBuilder() =>
      new ReceiptPagedRequestFilterBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ReceiptPagedRequestFilter &&
        date == other.date &&
        amount == other.amount &&
        name == other.name &&
        paidBy == other.paidBy &&
        categories == other.categories &&
        tags == other.tags &&
        status == other.status &&
        resolvedDate == other.resolvedDate &&
        createdAt == other.createdAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, date.hashCode);
    _$hash = $jc(_$hash, amount.hashCode);
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, paidBy.hashCode);
    _$hash = $jc(_$hash, categories.hashCode);
    _$hash = $jc(_$hash, tags.hashCode);
    _$hash = $jc(_$hash, status.hashCode);
    _$hash = $jc(_$hash, resolvedDate.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ReceiptPagedRequestFilter')
          ..add('date', date)
          ..add('amount', amount)
          ..add('name', name)
          ..add('paidBy', paidBy)
          ..add('categories', categories)
          ..add('tags', tags)
          ..add('status', status)
          ..add('resolvedDate', resolvedDate)
          ..add('createdAt', createdAt))
        .toString();
  }
}

class ReceiptPagedRequestFilterBuilder
    implements
        Builder<ReceiptPagedRequestFilter, ReceiptPagedRequestFilterBuilder> {
  _$ReceiptPagedRequestFilter? _$v;

  JsonObject? _date;
  JsonObject? get date => _$this._date;
  set date(JsonObject? date) => _$this._date = date;

  JsonObject? _amount;
  JsonObject? get amount => _$this._amount;
  set amount(JsonObject? amount) => _$this._amount = amount;

  JsonObject? _name;
  JsonObject? get name => _$this._name;
  set name(JsonObject? name) => _$this._name = name;

  JsonObject? _paidBy;
  JsonObject? get paidBy => _$this._paidBy;
  set paidBy(JsonObject? paidBy) => _$this._paidBy = paidBy;

  JsonObject? _categories;
  JsonObject? get categories => _$this._categories;
  set categories(JsonObject? categories) => _$this._categories = categories;

  JsonObject? _tags;
  JsonObject? get tags => _$this._tags;
  set tags(JsonObject? tags) => _$this._tags = tags;

  JsonObject? _status;
  JsonObject? get status => _$this._status;
  set status(JsonObject? status) => _$this._status = status;

  JsonObject? _resolvedDate;
  JsonObject? get resolvedDate => _$this._resolvedDate;
  set resolvedDate(JsonObject? resolvedDate) =>
      _$this._resolvedDate = resolvedDate;

  JsonObject? _createdAt;
  JsonObject? get createdAt => _$this._createdAt;
  set createdAt(JsonObject? createdAt) => _$this._createdAt = createdAt;

  ReceiptPagedRequestFilterBuilder() {
    ReceiptPagedRequestFilter._defaults(this);
  }

  ReceiptPagedRequestFilterBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _date = $v.date;
      _amount = $v.amount;
      _name = $v.name;
      _paidBy = $v.paidBy;
      _categories = $v.categories;
      _tags = $v.tags;
      _status = $v.status;
      _resolvedDate = $v.resolvedDate;
      _createdAt = $v.createdAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ReceiptPagedRequestFilter other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$ReceiptPagedRequestFilter;
  }

  @override
  void update(void Function(ReceiptPagedRequestFilterBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ReceiptPagedRequestFilter build() => _build();

  _$ReceiptPagedRequestFilter _build() {
    final _$result = _$v ??
        new _$ReceiptPagedRequestFilter._(
            date: date,
            amount: amount,
            name: name,
            paidBy: paidBy,
            categories: categories,
            tags: tags,
            status: status,
            resolvedDate: resolvedDate,
            createdAt: createdAt);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
