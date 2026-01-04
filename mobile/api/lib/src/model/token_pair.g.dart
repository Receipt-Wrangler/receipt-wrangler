// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'token_pair.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$TokenPair extends TokenPair {
  @override
  final String jwt;
  @override
  final String refreshToken;

  factory _$TokenPair([void Function(TokenPairBuilder)? updates]) =>
      (new TokenPairBuilder()..update(updates))._build();

  _$TokenPair._({required this.jwt, required this.refreshToken}) : super._() {
    BuiltValueNullFieldError.checkNotNull(jwt, r'TokenPair', 'jwt');
    BuiltValueNullFieldError.checkNotNull(
        refreshToken, r'TokenPair', 'refreshToken');
  }

  @override
  TokenPair rebuild(void Function(TokenPairBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  TokenPairBuilder toBuilder() => new TokenPairBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is TokenPair &&
        jwt == other.jwt &&
        refreshToken == other.refreshToken;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, jwt.hashCode);
    _$hash = $jc(_$hash, refreshToken.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'TokenPair')
          ..add('jwt', jwt)
          ..add('refreshToken', refreshToken))
        .toString();
  }
}

class TokenPairBuilder implements Builder<TokenPair, TokenPairBuilder> {
  _$TokenPair? _$v;

  String? _jwt;
  String? get jwt => _$this._jwt;
  set jwt(String? jwt) => _$this._jwt = jwt;

  String? _refreshToken;
  String? get refreshToken => _$this._refreshToken;
  set refreshToken(String? refreshToken) => _$this._refreshToken = refreshToken;

  TokenPairBuilder() {
    TokenPair._defaults(this);
  }

  TokenPairBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _jwt = $v.jwt;
      _refreshToken = $v.refreshToken;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(TokenPair other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$TokenPair;
  }

  @override
  void update(void Function(TokenPairBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  TokenPair build() => _build();

  _$TokenPair _build() {
    final _$result = _$v ??
        new _$TokenPair._(
            jwt:
                BuiltValueNullFieldError.checkNotNull(jwt, r'TokenPair', 'jwt'),
            refreshToken: BuiltValueNullFieldError.checkNotNull(
                refreshToken, r'TokenPair', 'refreshToken'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
