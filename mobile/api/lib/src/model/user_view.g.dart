// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user_view.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UserView extends UserView {
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
  final int? appRoleId;

  factory _$UserView([void Function(UserViewBuilder)? updates]) =>
      (UserViewBuilder()..update(updates))._build();

  _$UserView._(
      {required this.username,
      this.createdAt,
      this.createdBy,
      this.defaultAvatarColor,
      required this.displayName,
      required this.id,
      required this.isDummyUser,
      this.updatedAt,
      this.appRoleId})
      : super._();
  @override
  UserView rebuild(void Function(UserViewBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UserViewBuilder toBuilder() => UserViewBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UserView &&
        username == other.username &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        defaultAvatarColor == other.defaultAvatarColor &&
        displayName == other.displayName &&
        id == other.id &&
        isDummyUser == other.isDummyUser &&
        updatedAt == other.updatedAt &&
        appRoleId == other.appRoleId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, username.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, defaultAvatarColor.hashCode);
    _$hash = $jc(_$hash, displayName.hashCode);
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, isDummyUser.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, appRoleId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'UserView')
          ..add('username', username)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('defaultAvatarColor', defaultAvatarColor)
          ..add('displayName', displayName)
          ..add('id', id)
          ..add('isDummyUser', isDummyUser)
          ..add('updatedAt', updatedAt)
          ..add('appRoleId', appRoleId))
        .toString();
  }
}

class UserViewBuilder implements Builder<UserView, UserViewBuilder> {
  _$UserView? _$v;

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

  int? _appRoleId;
  int? get appRoleId => _$this._appRoleId;
  set appRoleId(int? appRoleId) => _$this._appRoleId = appRoleId;

  UserViewBuilder() {
    UserView._defaults(this);
  }

  UserViewBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _username = $v.username;
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _defaultAvatarColor = $v.defaultAvatarColor;
      _displayName = $v.displayName;
      _id = $v.id;
      _isDummyUser = $v.isDummyUser;
      _updatedAt = $v.updatedAt;
      _appRoleId = $v.appRoleId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(UserView other) {
    _$v = other as _$UserView;
  }

  @override
  void update(void Function(UserViewBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UserView build() => _build();

  _$UserView _build() {
    final _$result = _$v ??
        _$UserView._(
          username: BuiltValueNullFieldError.checkNotNull(
              username, r'UserView', 'username'),
          createdAt: createdAt,
          createdBy: createdBy,
          defaultAvatarColor: defaultAvatarColor,
          displayName: BuiltValueNullFieldError.checkNotNull(
              displayName, r'UserView', 'displayName'),
          id: BuiltValueNullFieldError.checkNotNull(id, r'UserView', 'id'),
          isDummyUser: BuiltValueNullFieldError.checkNotNull(
              isDummyUser, r'UserView', 'isDummyUser'),
          updatedAt: updatedAt,
          appRoleId: appRoleId,
        );
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
