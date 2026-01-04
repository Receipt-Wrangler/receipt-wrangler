// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$User extends User {
  @override
  final String? password;
  @override
  final String username;
  @override
  final String? createdAt;
  @override
  final int? createdBy;
  @override
  final String? defaultAvatarColor;
  @override
  final String displayName;
  @override
  final int id;
  @override
  final bool isDummyUser;
  @override
  final String? updatedAt;
  @override
  final UserRole userRole;
  @override
  final String? lastLoginDate;

  factory _$User([void Function(UserBuilder)? updates]) =>
      (new UserBuilder()..update(updates))._build();

  _$User._(
      {this.password,
      required this.username,
      this.createdAt,
      this.createdBy,
      this.defaultAvatarColor,
      required this.displayName,
      required this.id,
      required this.isDummyUser,
      this.updatedAt,
      required this.userRole,
      this.lastLoginDate})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(username, r'User', 'username');
    BuiltValueNullFieldError.checkNotNull(displayName, r'User', 'displayName');
    BuiltValueNullFieldError.checkNotNull(id, r'User', 'id');
    BuiltValueNullFieldError.checkNotNull(isDummyUser, r'User', 'isDummyUser');
    BuiltValueNullFieldError.checkNotNull(userRole, r'User', 'userRole');
  }

  @override
  User rebuild(void Function(UserBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UserBuilder toBuilder() => new UserBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is User &&
        password == other.password &&
        username == other.username &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        defaultAvatarColor == other.defaultAvatarColor &&
        displayName == other.displayName &&
        id == other.id &&
        isDummyUser == other.isDummyUser &&
        updatedAt == other.updatedAt &&
        userRole == other.userRole &&
        lastLoginDate == other.lastLoginDate;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, password.hashCode);
    _$hash = $jc(_$hash, username.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, defaultAvatarColor.hashCode);
    _$hash = $jc(_$hash, displayName.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, isDummyUser.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, userRole.hashCode);
    _$hash = $jc(_$hash, lastLoginDate.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'User')
          ..add('password', password)
          ..add('username', username)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('defaultAvatarColor', defaultAvatarColor)
          ..add('displayName', displayName)
          ..add('id', id)
          ..add('isDummyUser', isDummyUser)
          ..add('updatedAt', updatedAt)
          ..add('userRole', userRole)
          ..add('lastLoginDate', lastLoginDate))
        .toString();
  }
}

class UserBuilder implements Builder<User, UserBuilder> {
  _$User? _$v;

  String? _password;
  String? get password => _$this._password;
  set password(String? password) => _$this._password = password;

  String? _username;
  String? get username => _$this._username;
  set username(String? username) => _$this._username = username;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  String? _defaultAvatarColor;
  String? get defaultAvatarColor => _$this._defaultAvatarColor;
  set defaultAvatarColor(String? defaultAvatarColor) =>
      _$this._defaultAvatarColor = defaultAvatarColor;

  String? _displayName;
  String? get displayName => _$this._displayName;
  set displayName(String? displayName) => _$this._displayName = displayName;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  bool? _isDummyUser;
  bool? get isDummyUser => _$this._isDummyUser;
  set isDummyUser(bool? isDummyUser) => _$this._isDummyUser = isDummyUser;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  UserRole? _userRole;
  UserRole? get userRole => _$this._userRole;
  set userRole(UserRole? userRole) => _$this._userRole = userRole;

  String? _lastLoginDate;
  String? get lastLoginDate => _$this._lastLoginDate;
  set lastLoginDate(String? lastLoginDate) =>
      _$this._lastLoginDate = lastLoginDate;

  UserBuilder() {
    User._defaults(this);
  }

  UserBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _password = $v.password;
      _username = $v.username;
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _defaultAvatarColor = $v.defaultAvatarColor;
      _displayName = $v.displayName;
      _id = $v.id;
      _isDummyUser = $v.isDummyUser;
      _updatedAt = $v.updatedAt;
      _userRole = $v.userRole;
      _lastLoginDate = $v.lastLoginDate;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(User other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$User;
  }

  @override
  void update(void Function(UserBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  User build() => _build();

  _$User _build() {
    final _$result = _$v ??
        new _$User._(
            password: password,
            username: BuiltValueNullFieldError.checkNotNull(
                username, r'User', 'username'),
            createdAt: createdAt,
            createdBy: createdBy,
            defaultAvatarColor: defaultAvatarColor,
            displayName: BuiltValueNullFieldError.checkNotNull(
                displayName, r'User', 'displayName'),
            id: BuiltValueNullFieldError.checkNotNull(id, r'User', 'id'),
            isDummyUser: BuiltValueNullFieldError.checkNotNull(
                isDummyUser, r'User', 'isDummyUser'),
            updatedAt: updatedAt,
            userRole: BuiltValueNullFieldError.checkNotNull(
                userRole, r'User', 'userRole'),
            lastLoginDate: lastLoginDate);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
