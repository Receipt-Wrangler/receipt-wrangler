// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'check_email_connectivity_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$CheckEmailConnectivityCommand extends CheckEmailConnectivityCommand {
  @override
  final int? id;
  @override
  final String? host;
  @override
  final int? port;
  @override
  final String? username;
  @override
  final String? password;

  factory _$CheckEmailConnectivityCommand(
          [void Function(CheckEmailConnectivityCommandBuilder)? updates]) =>
      (new CheckEmailConnectivityCommandBuilder()..update(updates))._build();

  _$CheckEmailConnectivityCommand._(
      {this.id, this.host, this.port, this.username, this.password})
      : super._();

  @override
  CheckEmailConnectivityCommand rebuild(
          void Function(CheckEmailConnectivityCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  CheckEmailConnectivityCommandBuilder toBuilder() =>
      new CheckEmailConnectivityCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is CheckEmailConnectivityCommand &&
        id == other.id &&
        host == other.host &&
        port == other.port &&
        username == other.username &&
        password == other.password;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, host.hashCode);
    _$hash = $jc(_$hash, port.hashCode);
    _$hash = $jc(_$hash, username.hashCode);
    _$hash = $jc(_$hash, password.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'CheckEmailConnectivityCommand')
          ..add('id', id)
          ..add('host', host)
          ..add('port', port)
          ..add('username', username)
          ..add('password', password))
        .toString();
  }
}

class CheckEmailConnectivityCommandBuilder
    implements
        Builder<CheckEmailConnectivityCommand,
            CheckEmailConnectivityCommandBuilder> {
  _$CheckEmailConnectivityCommand? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  String? _host;
  String? get host => _$this._host;
  set host(String? host) => _$this._host = host;

  int? _port;
  int? get port => _$this._port;
  set port(int? port) => _$this._port = port;

  String? _username;
  String? get username => _$this._username;
  set username(String? username) => _$this._username = username;

  String? _password;
  String? get password => _$this._password;
  set password(String? password) => _$this._password = password;

  CheckEmailConnectivityCommandBuilder() {
    CheckEmailConnectivityCommand._defaults(this);
  }

  CheckEmailConnectivityCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _host = $v.host;
      _port = $v.port;
      _username = $v.username;
      _password = $v.password;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(CheckEmailConnectivityCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$CheckEmailConnectivityCommand;
  }

  @override
  void update(void Function(CheckEmailConnectivityCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  CheckEmailConnectivityCommand build() => _build();

  _$CheckEmailConnectivityCommand _build() {
    final _$result = _$v ??
        new _$CheckEmailConnectivityCommand._(
            id: id,
            host: host,
            port: port,
            username: username,
            password: password);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
