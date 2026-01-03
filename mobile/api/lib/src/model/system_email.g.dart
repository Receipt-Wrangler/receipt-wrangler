// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'system_email.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$SystemEmail extends SystemEmail {
  @override
  final String? password;
  @override
  final String? port;
  @override
  final String? host;
  @override
  final bool? useStartTLS;
  @override
  final String? username;
  @override
  final int id;
  @override
  final String createdAt;
  @override
  final int? createdBy;
  @override
  final String? createdByString;
  @override
  final String? updatedAt;

  factory _$SystemEmail([void Function(SystemEmailBuilder)? updates]) =>
      (new SystemEmailBuilder()..update(updates))._build();

  _$SystemEmail._(
      {this.password,
      this.port,
      this.host,
      this.useStartTLS,
      this.username,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(id, r'SystemEmail', 'id');
    BuiltValueNullFieldError.checkNotNull(
        createdAt, r'SystemEmail', 'createdAt');
  }

  @override
  SystemEmail rebuild(void Function(SystemEmailBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  SystemEmailBuilder toBuilder() => new SystemEmailBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is SystemEmail &&
        password == other.password &&
        port == other.port &&
        host == other.host &&
        useStartTLS == other.useStartTLS &&
        username == other.username &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, password.hashCode);
    _$hash = $jc(_$hash, port.hashCode);
    _$hash = $jc(_$hash, host.hashCode);
    _$hash = $jc(_$hash, useStartTLS.hashCode);
    _$hash = $jc(_$hash, username.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, createdByString.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'SystemEmail')
          ..add('password', password)
          ..add('port', port)
          ..add('host', host)
          ..add('useStartTLS', useStartTLS)
          ..add('username', username)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class SystemEmailBuilder
    implements Builder<SystemEmail, SystemEmailBuilder>, BaseModelBuilder {
  _$SystemEmail? _$v;

  String? _password;
  String? get password => _$this._password;
  set password(covariant String? password) => _$this._password = password;

  String? _port;
  String? get port => _$this._port;
  set port(covariant String? port) => _$this._port = port;

  String? _host;
  String? get host => _$this._host;
  set host(covariant String? host) => _$this._host = host;

  bool? _useStartTLS;
  bool? get useStartTLS => _$this._useStartTLS;
  set useStartTLS(covariant bool? useStartTLS) =>
      _$this._useStartTLS = useStartTLS;

  String? _username;
  String? get username => _$this._username;
  set username(covariant String? username) => _$this._username = username;

  int? _id;
  int? get id => _$this._id;
  set id(covariant int? id) => _$this._id = id;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(covariant String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(covariant int? createdBy) => _$this._createdBy = createdBy;

  String? _createdByString;
  String? get createdByString => _$this._createdByString;
  set createdByString(covariant String? createdByString) =>
      _$this._createdByString = createdByString;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(covariant String? updatedAt) => _$this._updatedAt = updatedAt;

  SystemEmailBuilder() {
    SystemEmail._defaults(this);
  }

  SystemEmailBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _password = $v.password;
      _port = $v.port;
      _host = $v.host;
      _useStartTLS = $v.useStartTLS;
      _username = $v.username;
      _id = $v.id;
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _createdByString = $v.createdByString;
      _updatedAt = $v.updatedAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(covariant SystemEmail other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$SystemEmail;
  }

  @override
  void update(void Function(SystemEmailBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  SystemEmail build() => _build();

  _$SystemEmail _build() {
    final _$result = _$v ??
        new _$SystemEmail._(
            password: password,
            port: port,
            host: host,
            useStartTLS: useStartTLS,
            username: username,
            id: BuiltValueNullFieldError.checkNotNull(id, r'SystemEmail', 'id'),
            createdAt: BuiltValueNullFieldError.checkNotNull(
                createdAt, r'SystemEmail', 'createdAt'),
            createdBy: createdBy,
            createdByString: createdByString,
            updatedAt: updatedAt);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
