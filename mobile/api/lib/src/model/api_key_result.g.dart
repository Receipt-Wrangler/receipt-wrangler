// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'api_key_result.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ApiKeyResult extends ApiKeyResult {
  @override
  final String key;

  factory _$ApiKeyResult([void Function(ApiKeyResultBuilder)? updates]) =>
      (new ApiKeyResultBuilder()..update(updates))._build();

  _$ApiKeyResult._({required this.key}) : super._() {
    BuiltValueNullFieldError.checkNotNull(key, r'ApiKeyResult', 'key');
  }

  @override
  ApiKeyResult rebuild(void Function(ApiKeyResultBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ApiKeyResultBuilder toBuilder() => new ApiKeyResultBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ApiKeyResult && key == other.key;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, key.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ApiKeyResult')..add('key', key))
        .toString();
  }
}

class ApiKeyResultBuilder
    implements Builder<ApiKeyResult, ApiKeyResultBuilder> {
  _$ApiKeyResult? _$v;

  String? _key;
  String? get key => _$this._key;
  set key(String? key) => _$this._key = key;

  ApiKeyResultBuilder() {
    ApiKeyResult._defaults(this);
  }

  ApiKeyResultBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _key = $v.key;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ApiKeyResult other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$ApiKeyResult;
  }

  @override
  void update(void Function(ApiKeyResultBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ApiKeyResult build() => _build();

  _$ApiKeyResult _build() {
    final _$result = _$v ??
        new _$ApiKeyResult._(
            key: BuiltValueNullFieldError.checkNotNull(
                key, r'ApiKeyResult', 'key'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
