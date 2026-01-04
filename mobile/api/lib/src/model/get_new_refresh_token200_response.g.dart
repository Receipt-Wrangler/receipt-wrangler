// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'get_new_refresh_token200_response.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$GetNewRefreshToken200Response extends GetNewRefreshToken200Response {
  @override
  final AnyOf anyOf;

  factory _$GetNewRefreshToken200Response(
          [void Function(GetNewRefreshToken200ResponseBuilder)? updates]) =>
      (new GetNewRefreshToken200ResponseBuilder()..update(updates))._build();

  _$GetNewRefreshToken200Response._({required this.anyOf}) : super._() {
    BuiltValueNullFieldError.checkNotNull(
        anyOf, r'GetNewRefreshToken200Response', 'anyOf');
  }

  @override
  GetNewRefreshToken200Response rebuild(
          void Function(GetNewRefreshToken200ResponseBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  GetNewRefreshToken200ResponseBuilder toBuilder() =>
      new GetNewRefreshToken200ResponseBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is GetNewRefreshToken200Response && anyOf == other.anyOf;
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
    return (newBuiltValueToStringHelper(r'GetNewRefreshToken200Response')
          ..add('anyOf', anyOf))
        .toString();
  }
}

class GetNewRefreshToken200ResponseBuilder
    implements
        Builder<GetNewRefreshToken200Response,
            GetNewRefreshToken200ResponseBuilder> {
  _$GetNewRefreshToken200Response? _$v;

  AnyOf? _anyOf;
  AnyOf? get anyOf => _$this._anyOf;
  set anyOf(AnyOf? anyOf) => _$this._anyOf = anyOf;

  GetNewRefreshToken200ResponseBuilder() {
    GetNewRefreshToken200Response._defaults(this);
  }

  GetNewRefreshToken200ResponseBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _anyOf = $v.anyOf;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(GetNewRefreshToken200Response other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$GetNewRefreshToken200Response;
  }

  @override
  void update(void Function(GetNewRefreshToken200ResponseBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  GetNewRefreshToken200Response build() => _build();

  _$GetNewRefreshToken200Response _build() {
    final _$result = _$v ??
        new _$GetNewRefreshToken200Response._(
            anyOf: BuiltValueNullFieldError.checkNotNull(
                anyOf, r'GetNewRefreshToken200Response', 'anyOf'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
