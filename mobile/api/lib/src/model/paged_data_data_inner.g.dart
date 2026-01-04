// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'paged_data_data_inner.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$PagedDataDataInner extends PagedDataDataInner {
  @override
  final AnyOf anyOf;

  factory _$PagedDataDataInner(
          [void Function(PagedDataDataInnerBuilder)? updates]) =>
      (new PagedDataDataInnerBuilder()..update(updates))._build();

  _$PagedDataDataInner._({required this.anyOf}) : super._() {
    BuiltValueNullFieldError.checkNotNull(
        anyOf, r'PagedDataDataInner', 'anyOf');
  }

  @override
  PagedDataDataInner rebuild(
          void Function(PagedDataDataInnerBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  PagedDataDataInnerBuilder toBuilder() =>
      new PagedDataDataInnerBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is PagedDataDataInner && anyOf == other.anyOf;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, anyOf.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'PagedDataDataInner')
          ..add('anyOf', anyOf))
        .toString();
  }
}

class PagedDataDataInnerBuilder
    implements Builder<PagedDataDataInner, PagedDataDataInnerBuilder> {
  _$PagedDataDataInner? _$v;

  AnyOf? _anyOf;
  AnyOf? get anyOf => _$this._anyOf;
  set anyOf(AnyOf? anyOf) => _$this._anyOf = anyOf;

  PagedDataDataInnerBuilder() {
    PagedDataDataInner._defaults(this);
  }

  PagedDataDataInnerBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _anyOf = $v.anyOf;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(PagedDataDataInner other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$PagedDataDataInner;
  }

  @override
  void update(void Function(PagedDataDataInnerBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  PagedDataDataInner build() => _build();

  _$PagedDataDataInner _build() {
    final _$result = _$v ??
        new _$PagedDataDataInner._(
            anyOf: BuiltValueNullFieldError.checkNotNull(
                anyOf, r'PagedDataDataInner', 'anyOf'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
