// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'paged_data.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$PagedData extends PagedData {
  @override
  final BuiltList<PagedDataDataInner> data;
  @override
  final int totalCount;

  factory _$PagedData([void Function(PagedDataBuilder)? updates]) =>
      (new PagedDataBuilder()..update(updates))._build();

  _$PagedData._({required this.data, required this.totalCount}) : super._() {
    BuiltValueNullFieldError.checkNotNull(data, r'PagedData', 'data');
    BuiltValueNullFieldError.checkNotNull(
        totalCount, r'PagedData', 'totalCount');
  }

  @override
  PagedData rebuild(void Function(PagedDataBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  PagedDataBuilder toBuilder() => new PagedDataBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is PagedData &&
        data == other.data &&
        totalCount == other.totalCount;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, data.hashCode);
    _$hash = $jc(_$hash, totalCount.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'PagedData')
          ..add('data', data)
          ..add('totalCount', totalCount))
        .toString();
  }
}

class PagedDataBuilder implements Builder<PagedData, PagedDataBuilder> {
  _$PagedData? _$v;

  ListBuilder<PagedDataDataInner>? _data;
  ListBuilder<PagedDataDataInner> get data =>
      _$this._data ??= new ListBuilder<PagedDataDataInner>();
  set data(ListBuilder<PagedDataDataInner>? data) => _$this._data = data;

  int? _totalCount;
  int? get totalCount => _$this._totalCount;
  set totalCount(int? totalCount) => _$this._totalCount = totalCount;

  PagedDataBuilder() {
    PagedData._defaults(this);
  }

  PagedDataBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _data = $v.data.toBuilder();
      _totalCount = $v.totalCount;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(PagedData other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$PagedData;
  }

  @override
  void update(void Function(PagedDataBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  PagedData build() => _build();

  _$PagedData _build() {
    _$PagedData _$result;
    try {
      _$result = _$v ??
          new _$PagedData._(
              data: data.build(),
              totalCount: BuiltValueNullFieldError.checkNotNull(
                  totalCount, r'PagedData', 'totalCount'));
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'data';
        data.build();
      } catch (e) {
        throw new BuiltValueNestedFieldError(
            r'PagedData', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
