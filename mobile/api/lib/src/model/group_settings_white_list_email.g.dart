// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'group_settings_white_list_email.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$GroupSettingsWhiteListEmail extends GroupSettingsWhiteListEmail {
  @override
  final int id;
  @override
  final int groupSettingsId;
  @override
  final String email;
  @override
  final String? createdAt;
  @override
  final int? createdBy;
  @override
  final String? updatedAt;

  factory _$GroupSettingsWhiteListEmail(
          [void Function(GroupSettingsWhiteListEmailBuilder)? updates]) =>
      (new GroupSettingsWhiteListEmailBuilder()..update(updates))._build();

  _$GroupSettingsWhiteListEmail._(
      {required this.id,
      required this.groupSettingsId,
      required this.email,
      this.createdAt,
      this.createdBy,
      this.updatedAt})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(
        id, r'GroupSettingsWhiteListEmail', 'id');
    BuiltValueNullFieldError.checkNotNull(
        groupSettingsId, r'GroupSettingsWhiteListEmail', 'groupSettingsId');
    BuiltValueNullFieldError.checkNotNull(
        email, r'GroupSettingsWhiteListEmail', 'email');
  }

  @override
  GroupSettingsWhiteListEmail rebuild(
          void Function(GroupSettingsWhiteListEmailBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  GroupSettingsWhiteListEmailBuilder toBuilder() =>
      new GroupSettingsWhiteListEmailBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is GroupSettingsWhiteListEmail &&
        id == other.id &&
        groupSettingsId == other.groupSettingsId &&
        email == other.email &&
        createdAt == other.createdAt &&
        createdBy == other.createdBy &&
        updatedAt == other.updatedAt;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, id.hashCode);
    _$hash = $jc(_$hash, groupSettingsId.hashCode);
    _$hash = $jc(_$hash, email.hashCode);
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, createdBy.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'GroupSettingsWhiteListEmail')
          ..add('id', id)
          ..add('groupSettingsId', groupSettingsId)
          ..add('email', email)
          ..add('createdAt', createdAt)
          ..add('createdBy', createdBy)
          ..add('updatedAt', updatedAt))
        .toString();
  }
}

class GroupSettingsWhiteListEmailBuilder
    implements
        Builder<GroupSettingsWhiteListEmail,
            GroupSettingsWhiteListEmailBuilder> {
  _$GroupSettingsWhiteListEmail? _$v;

  int? _id;
  int? get id => _$this._id;
  set id(int? id) => _$this._id = id;

  int? _groupSettingsId;
  int? get groupSettingsId => _$this._groupSettingsId;
  set groupSettingsId(int? groupSettingsId) =>
      _$this._groupSettingsId = groupSettingsId;

  String? _email;
  String? get email => _$this._email;
  set email(String? email) => _$this._email = email;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  int? _createdBy;
  int? get createdBy => _$this._createdBy;
  set createdBy(int? createdBy) => _$this._createdBy = createdBy;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  GroupSettingsWhiteListEmailBuilder() {
    GroupSettingsWhiteListEmail._defaults(this);
  }

  GroupSettingsWhiteListEmailBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _id = $v.id;
      _groupSettingsId = $v.groupSettingsId;
      _email = $v.email;
      _createdAt = $v.createdAt;
      _createdBy = $v.createdBy;
      _updatedAt = $v.updatedAt;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(GroupSettingsWhiteListEmail other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$GroupSettingsWhiteListEmail;
  }

  @override
  void update(void Function(GroupSettingsWhiteListEmailBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  GroupSettingsWhiteListEmail build() => _build();

  _$GroupSettingsWhiteListEmail _build() {
    final _$result = _$v ??
        new _$GroupSettingsWhiteListEmail._(
            id: BuiltValueNullFieldError.checkNotNull(
                id, r'GroupSettingsWhiteListEmail', 'id'),
            groupSettingsId: BuiltValueNullFieldError.checkNotNull(
                groupSettingsId,
                r'GroupSettingsWhiteListEmail',
                'groupSettingsId'),
            email: BuiltValueNullFieldError.checkNotNull(
                email, r'GroupSettingsWhiteListEmail', 'email'),
            createdAt: createdAt,
            createdBy: createdBy,
            updatedAt: updatedAt);
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
