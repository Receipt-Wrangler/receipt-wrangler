// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'upsert_system_email_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UpsertSystemEmailCommand extends UpsertSystemEmailCommand {
  @override
  final String host;
  @override
  final String port;
  @override
  final String username;
  @override
  final String password;
  @override
  final bool? useStartTLS;

  factory _$UpsertSystemEmailCommand(
          [void Function(UpsertSystemEmailCommandBuilder)? updates]) =>
      (new UpsertSystemEmailCommandBuilder()..update(updates))._build();

  _$UpsertSystemEmailCommand._(
      {required this.host,
      required this.port,
      required this.username,
      required this.password,
      this.useStartTLS})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        host, r'UpsertSystemEmailCommand', 'host');
    BuiltValueNullFieldError.checkNotNull(
        port, r'UpsertSystemEmailCommand', 'port');
    BuiltValueNullFieldError.checkNotNull(
        username, r'UpsertSystemEmailCommand', 'username');
    BuiltValueNullFieldError.checkNotNull(
        password, r'UpsertSystemEmailCommand', 'password');
  }

  @override
  UpsertSystemEmailCommand rebuild(
          void Function(UpsertSystemEmailCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UpsertSystemEmailCommandBuilder toBuilder() =>
      new UpsertSystemEmailCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UpsertSystemEmailCommand &&
        host == other.host &&
        port == other.port &&
        username == other.username &&
        password == other.password &&
        useStartTLS == other.useStartTLS;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, host.hashCode);
    _$hash = $jc(_$hash, port.hashCode);
    _$hash = $jc(_$hash, username.hashCode);
    _$hash = $jc(_$hash, password.hashCode);
    _$hash = $jc(_$hash, useStartTLS.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UpsertSystemEmailCommand')
          ..add('host', host)
          ..add('port', port)
          ..add('username', username)
          ..add('password', password)
          ..add('useStartTLS', useStartTLS))
        .toString();
  }
}

class UpsertSystemEmailCommandBuilder
    implements
        Builder<UpsertSystemEmailCommand, UpsertSystemEmailCommandBuilder> {
  _$UpsertSystemEmailCommand? _$v;

  String? _host;
  String? get host => _$this._host;
  set host(String? host) => _$this._host = host;

  String? _port;
  String? get port => _$this._port;
  set port(String? port) => _$this._port = port;

  String? _username;
  String? get username => _$this._username;
  set username(String? username) => _$this._username = username;

  String? _password;
  String? get password => _$this._password;
  set password(String? password) => _$this._password = password;

  bool? _useStartTLS;
  bool? get useStartTLS => _$this._useStartTLS;
  set useStartTLS(bool? useStartTLS) => _$this._useStartTLS = useStartTLS;

  UpsertSystemEmailCommandBuilder() {
    UpsertSystemEmailCommand._defaults(this);
  }

  UpsertSystemEmailCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _host = $v.host;
      _port = $v.port;
      _username = $v.username;
      _password = $v.password;
      _useStartTLS = $v.useStartTLS;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UpsertSystemEmailCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UpsertSystemEmailCommand;
  }

  @override
  void update(void Function(UpsertSystemEmailCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UpsertSystemEmailCommand build() => _build();

  _$UpsertSystemEmailCommand _build() {
    final _$result = _$v ??
        new _$UpsertSystemEmailCommand._(
            host: BuiltValueNullFieldError.checkNotNull(
                host, r'UpsertSystemEmailCommand', 'host'),
            port: BuiltValueNullFieldError.checkNotNull(
                port, r'UpsertSystemEmailCommand', 'port'),
            username: BuiltValueNullFieldError.checkNotNull(
                username, r'UpsertSystemEmailCommand', 'username'),
            password: BuiltValueNullFieldError.checkNotNull(
                password, r'UpsertSystemEmailCommand', 'password'),
            useStartTLS: useStartTLS);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
