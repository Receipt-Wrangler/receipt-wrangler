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
  final int? groupRoleId;
  @override
  final String? updatedAt;
  @override
  final int userId;
  @override
  final BuiltList<int>? categoryGrants;
  @override
  final BuiltList<int>? tagGrants;

  factory _$GroupMember([void Function(GroupMemberBuilder)? updates]) =>
      (GroupMemberBuilder()..update(updates))._build();

  _$GroupMember._(
      {this.createdAt,
      required this.groupId,
      this.groupRoleId,
      this.updatedAt,
      required this.userId,
      this.categoryGrants,
      this.tagGrants})
      : super._();
  @override
  GroupMember rebuild(void Function(GroupMemberBuilder) updates) =>
      (toBuilder()..update(updates)).build();

  @override
  GroupMemberBuilder toBuilder() => GroupMemberBuilder()..replace(this);

  @override
  bool operator ==(Object other) {
    if (identical(other, this)) return true;
    return other is GroupMember &&
        createdAt == other.createdAt &&
        groupId == other.groupId &&
        groupRoleId == other.groupRoleId &&
        updatedAt == other.updatedAt &&
        userId == other.userId &&
        categoryGrants == other.categoryGrants &&
        tagGrants == other.tagGrants;
  }

  @override
  int get hashCode {
    var _$hash = 0;
    _$hash = $jc(_$hash, createdAt.hashCode);
    _$hash = $jc(_$hash, groupId.hashCode);
    _$hash = $jc(_$hash, groupRoleId.hashCode);
    _$hash = $jc(_$hash, updatedAt.hashCode);
    _$hash = $jc(_$hash, userId.hashCode);
    _$hash = $jc(_$hash, categoryGrants.hashCode);
    _$hash = $jc(_$hash, tagGrants.hashCode);
    _$hash = $jf(_$hash);
    return _$hash;
  }

  @override
  String toString() {
    return (newBuiltValueToStringHelper(r'GroupMember')
          ..add('createdAt', createdAt)
          ..add('groupId', groupId)
          ..add('groupRoleId', groupRoleId)
          ..add('updatedAt', updatedAt)
          ..add('userId', userId)
          ..add('categoryGrants', categoryGrants)
          ..add('tagGrants', tagGrants))
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

  int? _groupRoleId;
  int? get groupRoleId => _$this._groupRoleId;
  set groupRoleId(int? groupRoleId) => _$this._groupRoleId = groupRoleId;

  String? _updatedAt;
  String? get updatedAt => _$this._updatedAt;
  set updatedAt(String? updatedAt) => _$this._updatedAt = updatedAt;

  int? _userId;
  int? get userId => _$this._userId;
  set userId(int? userId) => _$this._userId = userId;

  ListBuilder<int>? _categoryGrants;
  ListBuilder<int> get categoryGrants =>
      _$this._categoryGrants ??= ListBuilder<int>();
  set categoryGrants(ListBuilder<int>? categoryGrants) =>
      _$this._categoryGrants = categoryGrants;

  ListBuilder<int>? _tagGrants;
  ListBuilder<int> get tagGrants => _$this._tagGrants ??= ListBuilder<int>();
  set tagGrants(ListBuilder<int>? tagGrants) => _$this._tagGrants = tagGrants;

  GroupMemberBuilder() {
    GroupMember._defaults(this);
  }

  GroupMemberBuilder get _$this {
    final $v = _$v;
    if ($v != null) {
      _createdAt = $v.createdAt;
      _groupId = $v.groupId;
      _groupRoleId = $v.groupRoleId;
      _updatedAt = $v.updatedAt;
      _userId = $v.userId;
      _categoryGrants = $v.categoryGrants?.toBuilder();
      _tagGrants = $v.tagGrants?.toBuilder();
      _$v = null;
    }
    return this;
  }

  @override
  void replace(GroupMember other) {
    _$v = other as _$GroupMember;
  }

  @override
  void update(void Function(GroupMemberBuilder)? updates) {
    if (updates != null) updates(this);
  }

  @override
  GroupMember build() => _build();

  _$GroupMember _build() {
    _$GroupMember _$result;
    try {
      _$result = _$v ??
          _$GroupMember._(
            createdAt: createdAt,
            groupId: BuiltValueNullFieldError.checkNotNull(
                groupId, r'GroupMember', 'groupId'),
            groupRoleId: groupRoleId,
            updatedAt: updatedAt,
            userId: BuiltValueNullFieldError.checkNotNull(
                userId, r'GroupMember', 'userId'),
            categoryGrants: _categoryGrants?.build(),
            tagGrants: _tagGrants?.build(),
          );
    } catch (_) {
      late String _$failedField;
      try {
        _$failedField = 'categoryGrants';
        _categoryGrants?.build();
        _$failedField = 'tagGrants';
        _tagGrants?.build();
      } catch (e) {
        throw BuiltValueNestedFieldError(
            r'GroupMember', _$failedField, e.toString());
      }
      rethrow;
    }
    replace(_$result);
    return _$result;
  }
}

// ignore_for_file: deprecated_member_use_from_same_package,type=lint
