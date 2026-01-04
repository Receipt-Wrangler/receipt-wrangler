// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'logout_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$LogoutCommand extends LogoutCommand {
  @override
  final String refreshToken;

  factory _$LogoutCommand([void Function(LogoutCommandBuilder)? updates]) =>
      (new LogoutCommandBuilder()..update(updates))._build();

  _$LogoutCommand._({required this.refreshToken}) : super._() {
    BuiltValueNullFieldError.checkNotNull(
        refreshToken, r'LogoutCommand', 'refreshToken');
  }

  @override
  LogoutCommand rebuild(void Function(LogoutCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  LogoutCommandBuilder toBuilder() => new LogoutCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is LogoutCommand && refreshToken == other.refreshToken;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, refreshToken.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'LogoutCommand')
          ..add('refreshToken', refreshToken))
        .toString();
  }
}

class LogoutCommandBuilder
    implements Builder<LogoutCommand, LogoutCommandBuilder> {
  _$LogoutCommand? _$v;

  String? _refreshToken;
  String? get refreshToken => _$this._refreshToken;
  set refreshToken(String? refreshToken) => _$this._refreshToken = refreshToken;

  LogoutCommandBuilder() {
    LogoutCommand._defaults(this);
  }

  LogoutCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _refreshToken = $v.refreshToken;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(LogoutCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$LogoutCommand;
  }

  @override
  void update(void Function(LogoutCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  LogoutCommand build() => _build();

  _$LogoutCommand _build() {
    final _$result = _$v ??
        new _$LogoutCommand._(
            refreshToken: BuiltValueNullFieldError.checkNotNull(
                refreshToken, r'LogoutCommand', 'refreshToken'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
