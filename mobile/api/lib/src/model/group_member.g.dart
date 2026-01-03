// GENERATED CODE - DO NOT MODIFY BY HAND

part of 'group_member.dart';

// **************************************************************************
// BuiltValueGenerator
// **************************************************************************

class _$GroupMember extends GroupMember {
  @override
  final String? createdAt;
  @override
  final int groupId;
  @override
  final GroupRole groupRole;
  @override
  final String? updatedAt;
  @override
  final int userId;

  factory _$GroupMember([void Function(GroupMemberBuilder)? updates]) =>
      (new GroupMemberBuilder()..update(updates))._build();

  _$GroupMember._(
      {this.createdAt,
      required this.groupId,
      required this.groupRole,
      this.updatedAt,
      required this.userId})
      : super._() {
    BuiltValueNullFieldError.checkNotNull(groupId, r'GroupMember', 'groupId');
    BuiltValueNullFieldError.checkNotNull(
        groupRole, r'GroupMember', 'groupRole');
    BuiltValueNullFieldError.checkNotNull(userId, r'GroupMember', 'userId');
  }

  @override
  GroupMember rebuild(void Function(GroupMemberBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  GroupMemberBuilder toBuilder() => new GroupMemberBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is GroupMember &&
        createdAt == other.createdAt &&
        groupId == other.groupId &&
        groupRole == other.groupRole &&
        updatedAt == other.updatedAt &&
        userId == other.userId;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, groupRole.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, userId.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'GroupMember')
          ..add('createdAt', createdAt)
          ..add('groupId', groupId)
          ..add('groupRole', groupRole)
          ..add('updatedAt', updatedAt)
          ..add('userId', userId))
        .toString();
  }
}

class GroupMemberBuilder implements Builder<GroupMember, GroupMemberBuilder> {
  _$GroupMember? _$v;

  String? _createdAt;
  String? get createdAt => _$this._createdAt;
  set createdAt(String? createdAt) => _$this._createdAt = createdAt;

  int? _groupId;
  int? get groupId => _$this._groupId;
  set groupId(int? groupId) => _$this._groupId = groupId;

  GroupRole? _groupRole;
  GroupRole? get groupRole => _$this._groupRole;
  set groupRole(GroupRole? groupRole) => _$this._groupRole = groupRole;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  int? _userId;
  int? get userId => _$this._userId;
  set userId(int? userId) => _$this._userId = userId;

  GroupMemberBuilder() {
    GroupMember._defaults(this);
  }

  GroupMemberBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _createdAt = $v.createdAt;
      _groupId = $v.groupId;
      _groupRole = $v.groupRole;
      _updatedAt = $v.updatedAt;
      _userId = $v.userId;
      _$v = null;
    }
    return this;
  }

  @override
  void replace(GroupMember other) {
    ArgumentError.checkNotNull(other, 'other');
    _$v = other as _$GroupMember;
  }

  @override
  void update(void Function(GroupMemberBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  GroupMember build() => _build();

  _$GroupMember _build() {
    final _$result = _$v ??
        new _$GroupMember._(
            createdAt: createdAt,
            groupId: BuiltValueNullFieldError.checkNotNull(
                groupId, r'GroupMember', 'groupId'),
            groupRole: BuiltValueNullFieldError.checkNotNull(
                groupRole, r'GroupMember', 'groupRole'),
            updatedAt: updatedAt,
            userId: BuiltValueNullFieldError.checkNotNull(
                userId, r'GroupMember', 'userId'));
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
