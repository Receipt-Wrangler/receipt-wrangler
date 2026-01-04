// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'internal_error_response.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$InternalErrorResponse extends InternalErrorResponse {
  @override
  final String errorMsg;

  factory _$InternalErrorResponse(
          [void Function(InternalErrorResponseBuilder)? updates]) =>
      (new InternalErrorResponseBuilder()..update(updates))._build();

  _$InternalErrorResponse._({required this.errorMsg}) : super._() {
    BuiltValueNullFieldError.checkNotNull(
        errorMsg, r'InternalErrorResponse', 'errorMsg');
  }

  @override
  InternalErrorResponse rebuild(
          void Function(InternalErrorResponseBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  InternalErrorResponseBuilder toBuilder() =>
      new InternalErrorResponseBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is InternalErrorResponse && errorMsg == other.errorMsg;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, errorMsg.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'InternalErrorResponse')
          ..add('errorMsg', errorMsg))
        .toString();
  }
}

class InternalErrorResponseBuilder
    implements Builder<InternalErrorResponse, InternalErrorResponseBuilder> {
  _$InternalErrorResponse? _$v;

  String? _errorMsg;
  String? get errorMsg => _$this._errorMsg;
  set errorMsg(String? errorMsg) => _$this._errorMsg = errorMsg;

  InternalErrorResponseBuilder() {
    InternalErrorResponse._defaults(this);
  }

  InternalErrorResponseBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _errorMsg = $v.errorMsg;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(InternalErrorResponse other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$InternalErrorResponse;
  }

  @override
  void update(void Function(InternalErrorResponseBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  InternalErrorResponse build() => _build();

  _$InternalErrorResponse _build() {
    final _$result = _$v ??
        new _$InternalErrorResponse._(
            errorMsg: BuiltValueNullFieldError.checkNotNull(
                errorMsg, r'InternalErrorResponse', 'errorMsg'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
