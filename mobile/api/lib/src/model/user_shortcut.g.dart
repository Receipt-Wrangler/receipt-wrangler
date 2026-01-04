// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'user_shortcut.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$UserShortcut extends UserShortcut {
  @override
  final String name;
  @override
  final String? icon;
  @override
  final int? userPreferncesId;
  @override
  final String? url;
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

  factory _$UserShortcut([void Function(UserShortcutBuilder)? updates]) =>
      (new UserShortcutBuilder()..update(updates))._build();

  _$UserShortcut._(
      {required this.name,
      this.icon,
      this.userPreferncesId,
      this.url,
      required this.id,
      required this.createdAt,
      this.createdBy,
      this.createdByString,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(name, r'UserShortcut', 'name');
    BuiltValueNullFieldError.checkNotNull(id, r'UserShortcut', 'id');
    BuiltValueNullFieldError.checkNotNull(
        createdAt, r'UserShortcut', 'createdAt');
  }

  @override
  UserShortcut rebuild(void Function(UserShortcutBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  UserShortcutBuilder toBuilder() => new UserShortcutBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is UserShortcut &&
        name == other.name &&
        icon == other.icon &&
        userPreferncesId == other.userPreferncesId &&
        url == other.url &&
        id == other.id &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        createdByString == other.createdByString &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, name.hashCode);
    _$hash = $jc(_$hash, icon.hashCode);
    _$hash = $jc(_$hash, userPreferncesId.hashCode);
    _$hash = $jc(_$hash, url.hashCode);
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
    return (newBuiltValueToStringHelper(r'UserShortcut')
          ..add('name', name)
          ..add('icon', icon)
          ..add('userPreferncesId', userPreferncesId)
          ..add('url', url)
          ..add('id', id)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('createdByString', createdByString)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class UserShortcutBuilder
    implements Builder<UserShortcut, UserShortcutBuilder>, BaseModelBuilder {
  _$UserShortcut? _$v;

  String? _name;
  String? get name => _$this._name;
  set name(covariant String? name) => _$this._name = name;

  String? _icon;
  String? get icon => _$this._icon;
  set icon(covariant String? icon) => _$this._icon = icon;

  int? _userPreferncesId;
  int? get userPreferncesId => _$this._userPreferncesId;
  set userPreferncesId(covariant int? userPreferncesId) =>
      _$this._userPreferncesId = userPreferncesId;

  String? _url;
  String? get url => _$this._url;
  set url(covariant String? url) => _$this._url = url;

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

  UserShortcutBuilder() {
    UserShortcut._defaults(this);
  }

  UserShortcutBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _name = $v.name;
      _icon = $v.icon;
      _userPreferncesId = $v.userPreferncesId;
      _url = $v.url;
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
  void replace(covariant UserShortcut other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$UserShortcut;
  }

  @override
  void update(void Function(UserShortcutBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  UserShortcut build() => _build();

  _$UserShortcut _build() {
    final _$result = _$v ??
        new _$UserShortcut._(
            name: BuiltValueNullFieldError.checkNotNull(
                name, r'UserShortcut', 'name'),
            icon: icon,
            userPreferncesId: userPreferncesId,
            url: url,
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'UserShortcut', 'id'),
            createdAt: BuiltValueNullFieldError.checkNotNull(
                createdAt, r'UserShortcut', 'createdAt'),
            createdBy: createdBy,
            createdByString: createdByString,
            updatedAt: updatedAt);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
