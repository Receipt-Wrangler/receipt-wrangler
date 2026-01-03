// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'paged_request_field_value.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$PagedRequestFieldValue extends PagedRequestFieldValue {
  @override
  final OneOf oneOf;

  factory _$PagedRequestFieldValue(
          [void Function(PagedRequestFieldValueBuilder)? updates]) =>
      (new PagedRequestFieldValueBuilder()..update(updates))._build();

  _$PagedRequestFieldValue._({required this.oneOf}) : super._() {
    BuiltValueNullFieldError.checkNotNull(
        oneOf, r'PagedRequestFieldValue', 'oneOf');
  }

  @override
  PagedRequestFieldValue rebuild(
          void Function(PagedRequestFieldValueBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  PagedRequestFieldValueBuilder toBuilder() =>
      new PagedRequestFieldValueBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is PagedRequestFieldValue && oneOf == other.oneOf;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, oneOf.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'PagedRequestFieldValue')
          ..add('oneOf', oneOf))
        .toString();
  }
}

class PagedRequestFieldValueBuilder
    implements Builder<PagedRequestFieldValue, PagedRequestFieldValueBuilder> {
  _$PagedRequestFieldValue? _$v;

  OneOf? _oneOf;
  OneOf? get oneOf => _$this._oneOf;
  set oneOf(OneOf? oneOf) => _$this._oneOf = oneOf;

  PagedRequestFieldValueBuilder() {
    PagedRequestFieldValue._defaults(this);
  }

  PagedRequestFieldValueBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _oneOf = $v.oneOf;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(PagedRequestFieldValue other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$PagedRequestFieldValue;
  }

  @override
  void update(void Function(PagedRequestFieldValueBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  PagedRequestFieldValue build() => _build();

  _$PagedRequestFieldValue _build() {
    final _$result = _$v ??
        new _$PagedRequestFieldValue._(
            oneOf: BuiltValueNullFieldError.checkNotNull(
                oneOf, r'PagedRequestFieldValue', 'oneOf'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
