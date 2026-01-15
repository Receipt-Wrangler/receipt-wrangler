// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'api_key_filter.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ApiKeyFilter extends ApiKeyFilter {
  @override
  final AssociatedApiKeys? associatedApiKeys;

  factory _$ApiKeyFilter([void Function(ApiKeyFilterBuilder)? updates]) =>
      (new ApiKeyFilterBuilder()..update(updates))._build();

  _$ApiKeyFilter._({this.associatedApiKeys}) : super._();

  @override
  ApiKeyFilter rebuild(void Function(ApiKeyFilterBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ApiKeyFilterBuilder toBuilder() => new ApiKeyFilterBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ApiKeyFilter &&
        associatedApiKeys == other.associatedApiKeys;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, associatedApiKeys.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ApiKeyFilter')
          ..add('associatedApiKeys', associatedApiKeys))
        .toString();
  }
}

class ApiKeyFilterBuilder
    implements Builder<ApiKeyFilter, ApiKeyFilterBuilder> {
  _$ApiKeyFilter? _$v;

  AssociatedApiKeys? _associatedApiKeys;
  AssociatedApiKeys? get associatedApiKeys => _$this._associatedApiKeys;
  set associatedApiKeys(AssociatedApiKeys? associatedApiKeys) =>
      _$this._associatedApiKeys = associatedApiKeys;

  ApiKeyFilterBuilder() {
    ApiKeyFilter._defaults(this);
  }

  ApiKeyFilterBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _associatedApiKeys = $v.associatedApiKeys;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ApiKeyFilter other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$ApiKeyFilter;
  }

  @override
  void update(void Function(ApiKeyFilterBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ApiKeyFilter build() => _build();

  _$ApiKeyFilter _build() {
    final _$result =
        _$v ?? new _$ApiKeyFilter._(associatedApiKeys: associatedApiKeys);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
