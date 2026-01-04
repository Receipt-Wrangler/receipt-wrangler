// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'reset_password_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$ResetPasswordCommand extends ResetPasswordCommand {
  @override
  final String password;

  factory _$ResetPasswordCommand(
          [void Function(ResetPasswordCommandBuilder)? updates]) =>
      (new ResetPasswordCommandBuilder()..update(updates))._build();

  _$ResetPasswordCommand._({required this.password}) : super._() {
    BuiltValueNullFieldError.checkNotNull(
        password, r'ResetPasswordCommand', 'password');
  }

  @override
  ResetPasswordCommand rebuild(
          void Function(ResetPasswordCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  ResetPasswordCommandBuilder toBuilder() =>
      new ResetPasswordCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is ResetPasswordCommand && password == other.password;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, password.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'ResetPasswordCommand')
          ..add('password', password))
        .toString();
  }
}

class ResetPasswordCommandBuilder
    implements Builder<ResetPasswordCommand, ResetPasswordCommandBuilder> {
  _$ResetPasswordCommand? _$v;

  String? _password;
  String? get password => _$this._password;
  set password(String? password) => _$this._password = password;

  ResetPasswordCommandBuilder() {
    ResetPasswordCommand._defaults(this);
  }

  ResetPasswordCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _password = $v.password;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(ResetPasswordCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$ResetPasswordCommand;
  }

  @override
  void update(void Function(ResetPasswordCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  ResetPasswordCommand build() => _build();

  _$ResetPasswordCommand _build() {
    final _$result = _$v ??
        new _$ResetPasswordCommand._(
            password: BuiltValueNullFieldError.checkNotNull(
                password, r'ResetPasswordCommand', 'password'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
