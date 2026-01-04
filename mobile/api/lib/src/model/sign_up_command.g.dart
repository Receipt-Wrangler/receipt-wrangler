// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'sign_up_command.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$SignUpCommand extends SignUpCommand {
  @override
  final String username;
  @override
  final String password;
  @override
  final String? displayName;
  @override
  final bool? isDummyUser;
  @override
  final UserRole? userRole;

  factory _$SignUpCommand([void Function(SignUpCommandBuilder)? updates]) =>
      (new SignUpCommandBuilder()..update(updates))._build();

  _$SignUpCommand._(
      {required this.username,
      required this.password,
      this.displayName,
      this.isDummyUser,
      this.userRole})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        username, r'SignUpCommand', 'username');
    BuiltValueNullFieldError.checkNotNull(
        password, r'SignUpCommand', 'password');
  }

  @override
  SignUpCommand rebuild(void Function(SignUpCommandBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  SignUpCommandBuilder toBuilder() => new SignUpCommandBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is SignUpCommand &&
        username == other.username &&
        password == other.password &&
        displayName == other.displayName &&
        isDummyUser == other.isDummyUser &&
        userRole == other.userRole;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, username.hashCode);
    _$hash = $jc(_$hash, password.hashCode);
    _$hash = $jc(_$hash, displayName.hashCode);
    _$hash = $jc(_$hash, isDummyUser.hashCode);
    _$hash = $jc(_$hash, userRole.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'SignUpCommand')
          ..add('username', username)
          ..add('password', password)
          ..add('displayName', displayName)
          ..add('isDummyUser', isDummyUser)
          ..add('userRole', userRole))
        .toString();
  }
}

class SignUpCommandBuilder
    implements Builder<SignUpCommand, SignUpCommandBuilder> {
  _$SignUpCommand? _$v;

  String? _username;
  String? get username => _$this._username;
  set username(String? username) => _$this._username = username;

  String? _password;
  String? get password => _$this._password;
  set password(String? password) => _$this._password = password;

  String? _displayName;
  String? get displayName => _$this._displayName;
  set displayName(String? displayName) => _$this._displayName = displayName;

  bool? _isDummyUser;
  bool? get isDummyUser => _$this._isDummyUser;
  set isDummyUser(bool? isDummyUser) => _$this._isDummyUser = isDummyUser;

  UserRole? _userRole;
  UserRole? get userRole => _$this._userRole;
  set userRole(UserRole? userRole) => _$this._userRole = userRole;

  SignUpCommandBuilder() {
    SignUpCommand._defaults(this);
  }

  SignUpCommandBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _username = $v.username;
      _password = $v.password;
      _displayName = $v.displayName;
      _isDummyUser = $v.isDummyUser;
      _userRole = $v.userRole;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(SignUpCommand other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$SignUpCommand;
  }

  @override
  void update(void Function(SignUpCommandBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  SignUpCommand build() => _build();

  _$SignUpCommand _build() {
    final _$result = _$v ??
        new _$SignUpCommand._(
            username: BuiltValueNullFieldError.checkNotNull(
                username, r'SignUpCommand', 'username'),
            password: BuiltValueNullFieldError.checkNotNull(
                password, r'SignUpCommand', 'password'),
            displayName: displayName,
            isDummyUser: isDummyUser,
            userRole: userRole);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
